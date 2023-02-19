package cli

import (
	"fmt"
	"github.com/zakyyudha/tig/config"
	"github.com/zakyyudha/tig/pkg/shared"
	"github.com/zakyyudha/tig/src/commit/domain"
	"github.com/zakyyudha/tig/src/commit/usecase"
	"os"
)

type CLIHandler struct {
	commitUsecase usecase.CommitUsecase
}

func NewCliHandler(commitUsecase usecase.CommitUsecase) *CLIHandler {
	return &CLIHandler{commitUsecase: commitUsecase}
}

func (h *CLIHandler) Commit(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Not enough arguments for commit command")
		os.Exit(1)
	}

	command := args[1]

	switch command {
	case "commit":
		flags, values, err := shared.ParseGitCommitFlags(args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing commit flags: %v\n", err)
			os.Exit(1)
		}

		params := &domain.Params{}
		cmdStr := shared.BuildGitCommitCommand(flags, values, params)

		fmt.Printf("Running command: %s\n", cmdStr)

		outb, errb, err := shared.RunCommand(cmdStr)
		if len(outb) > 0 {
			fmt.Fprintln(os.Stdout, outb)
		} else {
			fmt.Fprintln(os.Stdout, errb)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
			os.Exit(1)
		}

		err = h.commitUsecase.ToSpreadsheet(params)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error pushing data to spreadsheet: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unrecognized command: %s\n%s", command, config.HelpUsage)
		os.Exit(1)
	}
}
