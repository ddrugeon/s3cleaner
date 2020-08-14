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
		---------------------------------------------------------------------------------------------
		The following ls command lists objects from specified bucket. In this 
		example, the user owns the bucket mybucket with the objects test.txt and prefix/test.txt.
	  
		s3cleanup ls -b mybucket

		Output:
			Name:          prefix/
			Last modified: 2020-08-13 14:33:19 +0000 UTC

			Name:          prefix/test.txt
			Last modified: 2020-08-13 14:33:29 +0000 UTC

			Name:          test.txt
			Last modified: 2020-08-13 14:33:04 +0000 UTC

		Found 3 items in bucket s3-ftp-sowee-prod

		---------------------------------------------------------------------------------------------
		The following ls command lists objects and all versions from specified bucket. In this 
		example, the user owns the bucket mybucket with the object test.txt with two different versions.

		s3cleanup ls -a -b mybucket

		Output:
			Name:          test.txt  (Latest Version) - Version ID:  a0RyXDUUC1qbrDzsZFyUhUJ8mxTiBEPb
			Last modified: 2020-08-14 09:34:36 +0000 UTC

			Name:          test.txt  - Version ID:  0AaGnbi3925aNKiP0pHXmPIuuiWqcEEm
			Last modified: 2020-08-14 09:34:02 +0000 UTC

			Found 2 versions in bucket test-emptybucket-ddh
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
