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

var rootCmd = &cobra.Command{
	Use:   "s3cleaner [OPTIONS] COMMAND",
	Short: "Tool to list / delete objects from a bucket",
	Long:  `s3cleaner helps to list and remove all object versions of specified bucket.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("s3cleaner %s \n", pkg.Version)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("profile", "p", "", "Use a specific profile from your AWS credential file.")
	rootCmd.PersistentFlags().StringP("region", "r", "eu-west-1", "The region to use. Overrides config/env settings.")
	rootCmd.PersistentFlags().StringP("bucket", "b", "", "The Bucket name to use.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
