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
	"os"

	"github.com/ddrugeon/s3cleaner/pkg/aws"
	"github.com/ddrugeon/s3cleaner/pkg/common"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List S3 objects of a specified Bucket",
	Long: `List S3 objects of a specified Bucket
	
	EXAMPLES:
		The following ls command lists objects from specified bucket. In this 
		example, the use owns the bucket mybucket with the objects test.txt and prefix/test.txt.
	  
		s3cleanup ls s3://mybucket

		Output:
			PRE prefix/
			test.txt

	`,
	Run: func(cmd *cobra.Command, args []string) {
		currentProfile, _ := cmd.Flags().GetString("profile")
		currentRegion, _ := cmd.Flags().GetString("region")
		currentBucket, _ := cmd.Flags().GetString("bucket")
		showAllVersions, _ := cmd.Flags().GetBool("all")

		if currentProfile == "" {
			if os.Getenv("AWS_PROFILE") != "" {
				currentProfile = os.Getenv("AWS_PROFILE")
			} else if os.Getenv("AWS_DEFAULT_PROFILE") != "" {
				currentProfile = os.Getenv("AWS_DEFAULT_PROFILE")
			} else {
				profiles := common.ListProfilesFromAWSConfig()

				currentProfile = common.SelectProfile(profiles)
			}
		}

		client := aws.NewClient(currentProfile, currentRegion)

		if currentBucket == "" {
			currentBucket = common.SelectBucket(client.GetBucketLists())
		}

		if showAllVersions {
			client.ListObjectVersions(currentBucket)
		} else {
			client.ListObjects(currentBucket)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolP("all", "a", false, "All object versions are also included.")
}
