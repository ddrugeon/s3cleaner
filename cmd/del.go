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

	"github.com/ddrugeon/s3cleaner/pkg/aws"
	"github.com/ddrugeon/s3cleaner/pkg/common"
	"github.com/spf13/cobra"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete S3 objects of a specified Bucket",
	Long:  `"Delete S3 objects of a specified Bucket`,
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
			currentBucket = common.SelectBucket(client.GetBucketList())
		}

		var objects []aws.Object
		var err error

		if showAllVersions {
			objects, err = client.ListObjectWithVersions(currentBucket)
		} else {
			objects, err = client.ListObjects(currentBucket)
		}

		if err != nil {
			common.ExitWithError("Unable to list objects from bucket %s:\n%v", currentBucket, err)
		}

		for _, item := range objects {
			if item.VersionID != "" {
				if item.IsLastVersion {
					fmt.Println("Name:         ", item.Name, " (Latest Version) - Version ID: ", item.VersionID)
				} else {
					fmt.Println("Name:         ", item.Name, " - Version ID: ", item.VersionID)
				}
			} else {
				fmt.Println("Name:         ", item.Name)
			}
		}
		fmt.Println("Delete", len(objects), "items in bucket", currentBucket)
		fmt.Println("")

		if common.ConfirmDeletion() {
			err := client.DeleteObjects(objects, currentBucket)
			if err != nil {
				common.ExitWithError("Unable to delete objects from bucket %s:\n%v", currentBucket, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(delCmd)

	lsCmd.Flags().BoolP("dry-run", "d", false, "Do not delete files or version")
}
