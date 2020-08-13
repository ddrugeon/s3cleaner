/*
Copyright © 2020 David Drugeon-Hamon <zebeurton@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package cmd regroups available commands to this tool
package cmd

import (
	"fmt"
	"os"

	"github.com/ddrugeon/s3cleaner/pkg"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of this tool",
	Long: `Give version information of this tool as JSON String
	{
		"Version":   Semantic version of this tool,
		"GitCommit": Git SHA1 commit,
		"GitState":  Git state,
		"BuildDate": Date when tool has been built,
		"GoVersion": Go tool version,
		"Compiler":  Compiler type,
		"GOARCH":    Architecture type (amd64, i386...),
		"GOOS":      OS type,
	}`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(pkg.JSONVersion())
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
