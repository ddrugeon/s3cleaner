package pkg

import (
	"encoding/json"
	"fmt"
	"runtime"
)

var (
	// GitCommit holds short commit hash of source tree
	GitCommit string

	// GitBranch holds current branch name the code is built off
	GitBranch string

	// GitState shows whether there are uncommitted changes
	GitState string

	// GitSummary holds output of git describe --tags --dirty --always
	GitSummary string

	// BuildDate holds RFC3339 formatted UTC date (build time)
	BuildDate string

	// Version holds contents of ./VERSION file, if exists, or the value passed via the -version option
	Version string
)

// JSONVersion returns a json formatted Version information
func JSONVersion() string {
	version := map[string]string{
		"Version":   Version,
		"GitCommit": GitCommit,
		"GitState":  GitState,
		"BuildDate": BuildDate,
		"GoVersion": runtime.Version(),
		"Compiler":  runtime.Compiler,
		"GOARCH":    runtime.GOARCH,
		"GOOS":      runtime.GOOS,
	}

	jsonVersion, err := json.Marshal(version)
	if err != nil {
		fmt.Println("Can't serialize")
	}
	return string(jsonVersion)
}
