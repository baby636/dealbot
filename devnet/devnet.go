package devnet

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

func runCmdsWithLog(ctx context.Context, name string, commands [][]string) {
	logFile, err := os.Create(name + ".log")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	for _, cmdArgs := range commands {
		log.Printf("command for %s: %s", name, strings.Join(cmdArgs, " "))
		cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = logFile
		cmd.Stderr = logFile
		// If ctx.Err()!=nil, we cancelled the command via SIGINT.
		if err := cmd.Run(); err != nil && ctx.Err() == nil {
			log.Printf("%s; check %s for derunLotusNodetails", err, logFile.Name())
			break
		}
	}
}

func runLotusDaemon(ctx context.Context, home string) {
	cmds := [][]string{
		{"lotus-seed", "genesis", "new", "localnet.json"},
		{"lotus-seed", "pre-seal", "--sector-size=2048", "--num-sectors=10"},
		{"lotus-seed", "genesis", "add-miner", "localnet.json",
			filepath.Join(home, ".genesis-sectors", "pre-seal-t01000.json")},
		{"lotus", "daemon", "--lotus-make-genesis=dev.gen",
			"--genesis-template=localnet.json", "--bootstrap=false"},
	}

	runCmdsWithLog(ctx, "lotus-daemon", cmds)
}

func runLotusMiner(ctx context.Context, home string) {
	cmds := [][]string{
		{"lotus", "wait-api"}, // wait for lotus node to run
		{"lotus", "wallet", "import",
			filepath.Join(home, ".genesis-sectors", "pre-seal-t01000.key")},
		{"lotus-miner", "init", "--genesis-miner", "--actor=t01000", "--sector-size=2048",
			"--pre-sealed-sectors=" + filepath.Join(home, ".genesis-sectors"),
			"--pre-sealed-metadata=" + filepath.Join(home, ".genesis-sectors", "pre-seal-t01000.json"),
			"--nosync"},
		{"lotus-miner", "run", "--nosync"},
	}

	runCmdsWithLog(ctx, "lotus-miner", cmds)
}

func publishDealsPeriodicallyCmd(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
		}

		cmd := exec.CommandContext(ctx, "lotus-miner",
			"storage-deals", "pending-publish", "--publish-now")
		cmd.Run() // we ignore errors
	}
}

func setDefaultWalletCmd(ctx context.Context) {
	// TODO: do this without a shell
	setDefaultWalletCmd := "lotus wallet list | grep t3 | awk '{print $1}' | xargs lotus wallet set-default"

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
		}

		cmd := exec.CommandContext(ctx, "sh", "-c", setDefaultWalletCmd)
		_, err := cmd.CombinedOutput()
		if err != nil {
			continue
		}
		// TODO: stop once we've set the default wallet once.
	}
}

func Main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(4)
	go func() {
		runLotusDaemon(ctx, home)
		wg.Done()
	}()

	go func() {
		runLotusMiner(ctx, home)
		wg.Done()
	}()

	go func() {
		publishDealsPeriodicallyCmd(ctx)
		wg.Done()
	}()

	go func() {
		setDefaultWalletCmd(ctx)
		wg.Done()
	}()

	// setup a signal handler to cancel the context
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-interrupt:
		log.Println("closing as we got interrupt")
		cancel()
	case <-ctx.Done():
	}

	wg.Wait()
}
