package lotus

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/node/repo"
	lru "github.com/hashicorp/golang-lru"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/mitchellh/go-homedir"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

type APICloser func()

type API interface {
	StartDeal(ctx context.Context, params *api.StartDealParams) (*cid.Cid, error)
	Import(ctx context.Context, ref api.FileRef) (*api.ImportRes, error)
	QueryAsk(ctx context.Context, p peer.ID, miner address.Address) (*storagemarket.StorageAsk, error)
	ChainHead(ctx context.Context) (*types.TipSet, error)
	MinerInfo(ctx context.Context, a address.Address, tsk types.TipSetKey) (miner.MinerInfo, error)
	DealPieceCID(ctx context.Context, root cid.Cid) (api.DataCIDSize, error)
	GetDealUpdates(ctx context.Context) (<-chan api.DealInfo, error)
}

type APIOpener struct {
	cache   *lru.ARCCache // cache shared across all instances of the api
	addr    string
	headers http.Header
}

func NewAPIOpener(cctx *cli.Context, cacheSize int) (*APIOpener, APICloser, error) {
	ac, err := lru.NewARC(cacheSize)
	if err != nil {
		return nil, nil, xerrors.Errorf("new arc cache: %w", err)
	}

	var rawaddr, rawtoken string

	if cctx.IsSet("api") {
		tokenMaddr := cctx.String("api")
		toks := strings.Split(tokenMaddr, ":")
		if len(toks) != 2 {
			return nil, nil, fmt.Errorf("invalid api tokens, expected <token>:<maddr>, got: %s", tokenMaddr)
		}

		rawtoken = toks[0]
		rawaddr = toks[1]
	} else if cctx.IsSet("repo") {
		repoPath := cctx.String("repo")
		p, err := homedir.Expand(repoPath)
		if err != nil {
			return nil, nil, xerrors.Errorf("expand home dir (%s): %w", repoPath, err)
		}

		r, err := repo.NewFS(p)
		if err != nil {
			return nil, nil, xerrors.Errorf("open repo at path: %s; %w", p, err)
		}

		ma, err := r.APIEndpoint()
		if err != nil {
			return nil, nil, xerrors.Errorf("api endpoint: %w", err)
		}

		token, err := r.APIToken()
		if err != nil {
			return nil, nil, xerrors.Errorf("api token: %w", err)
		}

		rawaddr = ma.String()
		rawtoken = string(token)
	} else {
		return nil, nil, xerrors.Errorf("cannot connect to lotus api: missing --api or --repo flags")
	}

	parsedAddr, err := ma.NewMultiaddr(rawaddr)
	if err != nil {
		return nil, nil, xerrors.Errorf("parse listen address: %w", err)
	}

	_, addr, err := manet.DialArgs(parsedAddr)
	if err != nil {
		return nil, nil, xerrors.Errorf("dial multiaddress: %w", err)
	}

	o := &APIOpener{
		cache:   ac,
		addr:    apiURI(addr),
		headers: apiHeaders(rawtoken),
	}

	return o, APICloser(func() {}), nil
}

func (o *APIOpener) Open(ctx context.Context) (API, APICloser, error) {
	api, closer, err := client.NewFullNodeRPC(ctx, o.addr, o.headers)
	if err != nil {
		return nil, nil, xerrors.Errorf("new full node rpc: %w", err)
	}

	cacheStore, err := NewCacheCtxStore(ctx, api, o.cache)
	if err != nil {
		return nil, nil, xerrors.Errorf("new cache store: %w", err)
	}

	lotusAPI := NewAPIWrapper(api, cacheStore)

	return lotusAPI, APICloser(closer), nil
}

func apiURI(addr string) string {
	return "ws://" + addr + "/rpc/v0"
}

func apiHeaders(token string) http.Header {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)
	return headers
}
