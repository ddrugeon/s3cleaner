package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ddrugeon/s3cleaner/pkg/common"
)

// Client represents an AWS Client
type Client struct {
	session *session.Session
	svc     *s3.S3
}

// NewClient returns a new AWS client instance with current configuration
func NewClient(profile string, region string) *Client {
	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState:       session.SharedConfigEnable,
				AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
				Profile:                 string(profile),
				Config:                  aws.Config{Region: aws.String(region)},
			},
		),
	)

	client := Client{
		session: sess,
		svc:     s3.New(sess),
	}

	return &client
}

func (client *Client) GetBucketLists() []string {
	result, err := client.svc.ListBuckets(nil)
	if err != nil {
		common.ExitWithError("Unable to list buckets:\n\t %v", err)
	}

	var buckets = []string{}
	for _, b := range result.Buckets {
		buckets = append(buckets, aws.StringValue(b.Name))
	}

	return buckets
}

func (client *Client) ListObjects(bucket string) {
	resp, err := client.svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		common.ExitWithError("Unable to list objects from bucket %s:\n\t %v", bucket, err)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("")
	}
	fmt.Println("Found", len(resp.Contents), "items in bucket", bucket)
	fmt.Println("")
}

func (client *Client) ListObjectVersions(bucket string) {
	resp, err := client.svc.ListObjectVersions(&s3.ListObjectVersionsInput{Bucket: aws.String(bucket)})
	if err != nil {
		common.ExitWithError("Unable to list objects from bucket %s:\n\t %v", bucket, err)
	}

	for _, item := range resp.Versions {
		if *item.IsLatest {
			fmt.Println("Name:         ", *item.Key, " (Latest Version) - Version ID: ", *item.VersionId)
		} else {
			fmt.Println("Name:         ", *item.Key, " - Version ID: ", *item.VersionId)
		}
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("")
	}
	fmt.Println("Found", len(resp.Versions), "versions in bucket", bucket)
	fmt.Println("")
}
