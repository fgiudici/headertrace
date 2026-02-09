package cmd

import (
	"fmt"
)

var (
	version   = "v0.0.0"
	gitCommit = ""
)

func getVersion() string {
	commit := ""
	if len(gitCommit) > 7 {
		commit = fmt.Sprintf(" (commit: %s)", gitCommit[:7])
	}
	return fmt.Sprintf("%s%s", version, commit)
}
