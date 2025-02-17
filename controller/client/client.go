package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/filecoin-project/dealbot/tasks"
	"github.com/ipld/go-car"
	dagjson "github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/urfave/cli/v2"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("controller")

type ErrRequestFailed struct {
	Code int
}

func (e ErrRequestFailed) Error() string {
	return fmt.Sprintf("Request failed: %s", http.StatusText(e.Code))
}

// Client is the API client that performs all operations
// against the dealbot controller.
type Client struct {
	client   *http.Client
	endpoint string
}

// New initializes a new API client
func New(ctx *cli.Context) *Client {
	endpoint := ctx.String("endpoint")

	log.Infow("dealbot controller client initialized", "addr", endpoint)

	return NewFromEndpoint(endpoint)
}

// NewFromEndpoint returns an API client at the given endpoint
func NewFromEndpoint(endpoint string) *Client {
	return &Client{
		client: &http.Client{
			// As a fallback, never take more than a minute.
			// Most client API calls should use a context.
			Timeout: time.Minute,
		},
		endpoint: endpoint,
	}
}

// Close the transport used by the client
func (c *Client) Close() error {
	if t, ok := c.client.Transport.(*http.Transport); ok {
		t.CloseIdleConnections()
	}
	return nil
}

func (c *Client) ListTasks(ctx context.Context) ([]tasks.Task, error) {
	resp, err := c.request(ctx, "GET", "/tasks", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tlp := tasks.Type.Tasks.NewBuilder()
	err = dagjson.Decoder(tlp, resp.Body)
	if err != nil {
		return nil, err
	}
	tsks := tlp.Build().(tasks.Tasks)

	return tsks.List(), nil
}

func (c *Client) PopTask(ctx context.Context, r tasks.PopTask) (tasks.Task, error) {
	var body bytes.Buffer
	err := dagjson.Encoder(r.Representation(), &body)
	if err != nil {
		return nil, err
	}

	resp, err := c.request(ctx, "POST", "/pop-task", &body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrRequestFailed{resp.StatusCode}
	}

	rp := tasks.Type.Task.NewBuilder()
	if err := dagjson.Decoder(rp, resp.Body); err != nil {
		return nil, err
	}
	return rp.Build().(tasks.Task), nil
}

func (c *Client) UpdateTask(ctx context.Context, uuid string, r tasks.UpdateTask) (tasks.Task, error) {
	var body bytes.Buffer
	err := dagjson.Encoder(r.Representation(), &body)
	if err != nil {
		return nil, err
	}

	resp, err := c.request(ctx, "PATCH", "/tasks/"+uuid, &body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrRequestFailed{resp.StatusCode}
	}

	rp := tasks.Type.Task.NewBuilder()
	if err := dagjson.Decoder(rp, resp.Body); err != nil {
		return nil, err
	}
	return rp.Build().(tasks.Task), nil
}

func (c *Client) GetTask(ctx context.Context, uuid string) (tasks.Task, error) {
	resp, err := c.request(ctx, "GET", "/tasks/"+uuid, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrRequestFailed{resp.StatusCode}
	}

	tp := tasks.Type.Task.NewBuilder()
	err = dagjson.Decoder(tp, resp.Body)
	if err != nil {
		return nil, err
	}
	return tp.Build().(tasks.Task), nil
}

func (c *Client) Drain(ctx context.Context, worker string) error {
	resp, err := c.request(ctx, "POST", "/drain/"+worker, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrRequestFailed{resp.StatusCode}
	}
	return nil
}

func (c *Client) Complete(ctx context.Context, worker string) error {
	resp, err := c.request(ctx, "POST", "/complete/"+worker, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrRequestFailed{resp.StatusCode}
	}
	return nil
}

func (c *Client) CARExport(ctx context.Context) (*car.CarReader, func() error, error) {
	resp, err := c.request(ctx, "GET", "/car", nil)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, ErrRequestFailed{resp.StatusCode}
	}
	rdr, err := car.NewCarReader(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, nil, err
	}
	return rdr, resp.Body.Close, nil
}

func (c *Client) CreateStorageTask(ctx context.Context, storageTask tasks.StorageTask) (tasks.Task, error) {
	var body bytes.Buffer
	err := dagjson.Encoder(storageTask.Representation(), &body)
	if err != nil {
		return nil, err
	}

	resp, err := c.request(ctx, "POST", "/tasks/storage", &body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, ErrRequestFailed{resp.StatusCode}
	}

	tp := tasks.Type.Task.NewBuilder()
	err = dagjson.Decoder(tp, resp.Body)
	if err != nil {
		return nil, err
	}
	return tp.Build().(tasks.Task), nil
}

func (c *Client) CreateRetrievalTask(ctx context.Context, retrievalTask tasks.RetrievalTask) (tasks.Task, error) {
	var body bytes.Buffer
	err := dagjson.Encoder(retrievalTask.Representation(), &body)
	if err != nil {
		return nil, err
	}

	resp, err := c.request(ctx, "POST", "/tasks/retrieval", &body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, ErrRequestFailed{resp.StatusCode}
	}

	tp := tasks.Type.Task.NewBuilder()
	err = dagjson.Decoder(tp, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not parse /tasks/retrieval response: %w", err)
	}
	return tp.Build().(tasks.Task), nil
}

func (c *Client) request(ctx context.Context, method string, path string, body io.Reader, headers ...string) (*http.Response, error) {
	if len(headers)%2 != 0 {
		return nil, fmt.Errorf("headers must be tuples: key1, value1, key2, value2")
	}
	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, body)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(headers); i = i + 2 {
		req.Header.Add(headers[i], headers[i+1])
	}

	return c.client.Do(req)
}
