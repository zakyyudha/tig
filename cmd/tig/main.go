package main

import (
	"fmt"
	"github.com/zakyyudha/tig/config"
	"github.com/zakyyudha/tig/src/commit/delivery/cli"
	"github.com/zakyyudha/tig/src/commit/usecase"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, config.HelpUsage)
		os.Exit(1)
	}

	cfg := config.New()
	commitUsecase := usecase.NewCommitUsecaseImpl(cfg)
	commitCliDelivery := cli.NewCliHandler(commitUsecase)
	commitCliDelivery.Commit(os.Args)
}
