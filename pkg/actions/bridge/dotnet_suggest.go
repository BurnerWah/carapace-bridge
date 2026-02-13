package bridge

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/carapace-sh/carapace"
	shlex "github.com/carapace-sh/carapace-shlex"
)

// ActionDotnetSuggest bridges https://github.com/dotnet/command-line-api
//
// This does not provide descriptions for completions
func ActionDotnetSuggest(command ...string) carapace.Action {
	return actionCommand(command...)(func(command ...string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			path, err := exec.LookPath(command[0])
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}

			fullLine := append(command, c.Args...)
			fullLine = append(fullLine, c.Value)

			input := shlex.Join(fullLine)

			// TODO: add --position (if possible)
			args := []string{"get", "--executable", path, "--", input}

			return carapace.ActionExecCommand("dotnet-suggest", args...)(func(output []byte) carapace.Action {
				separator := "\n"
				if runtime.GOOS == "windows" {
					separator = "\r\n"
				}
				lines := strings.Split(string(output), separator)

				switch len(lines) {
				case 0:
					return carapace.ActionFiles()
				default:
					return carapace.ActionValues(lines...)
				}
			})
		})
	})
}
