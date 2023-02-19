package shared

import (
	"bytes"
	"fmt"
	"github.com/zakyyudha/tig/config"
	"github.com/zakyyudha/tig/src/commit/domain"
	"os/exec"
	"runtime"
	"strings"
)

func ParseGitCommitFlags(args []string) (flags []string, values []string, err error) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			if i == len(args)-1 {
				err = fmt.Errorf("Flag %s requires a value", arg)
				return
			}
			flags = append(flags, arg)
			values = append(values, args[i+1])
			i++
		} else {
			values = append(values, arg)
		}
	}
	return
}

func BuildGitCommitCommand(flags []string, values []string, params *domain.Params) string {
	cmdStr := "git commit"
	for i, f := range flags {
		if strings.HasPrefix(f, "-m") || strings.HasPrefix(f, "--message") {
			params.Activity = values[i]
		}
		if strings.HasPrefix(f, "-tN") || strings.HasPrefix(f, "--tigNote") {
			params.Note = values[i]
			continue
		}
		if strings.HasPrefix(f, "-tJC") || strings.HasPrefix(f, "--tigJiraCode") {
			params.ReferTo = fmt.Sprintf(config.JiraURL, values[i])
			continue
		}
		cmdStr += fmt.Sprintf(" %s '%s'", f, values[i])
	}
	return cmdStr
}

func RunCommand(cmdStr string) (string, string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", cmdStr)
	default:
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	return outb.String(), errb.String(), err
}
