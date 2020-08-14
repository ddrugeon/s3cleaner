package aws

import (
	"time"

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

// Object represents a S3 Object
type Object struct {
	Name          string
	LastModified  time.Time
	IsLastVersion bool
	VersionID     string
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

// GetBucketList returns list of buckets present in current region and profile
func (client *Client) GetBucketList() []string {
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

// ListObjects returns list of objects present in given bucket
func (client *Client) ListObjects(bucket string) ([]Object, error) {
	resp, err := client.svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		return nil, err
	}

	objects := make([]Object, len(resp.Contents))
	for index, item := range resp.Contents {
		objects[index] = Object{
			Name:         *item.Key,
			LastModified: *item.LastModified,
		}
	}

	return objects, nil
}

// ListObjectWithVersions returns list of objects present (with each version) in given bucket
func (client *Client) ListObjectWithVersions(bucket string) ([]Object, error) {
	resp, err := client.svc.ListObjectVersions(&s3.ListObjectVersionsInput{Bucket: aws.String(bucket)})
	if err != nil {
		return nil, err
	}

	objects := make([]Object, len(resp.Versions))
	for index, item := range resp.Versions {
		objects[index] = Object{
			Name:          *item.Key,
			LastModified:  *item.LastModified,
			IsLastVersion: *item.IsLatest,
			VersionID:     *item.VersionId,
		}
	}

	return objects, nil
}

// DeleteObjects delete specified objects from given bucket
func (client *Client) DeleteObjects(items []Object, bucket string) error {

	var err error
	for _, item := range items {
		if item.VersionID != "" {
			_, err = client.svc.DeleteObject(&s3.DeleteObjectInput{Bucket: &bucket, Key: &item.Name, VersionId: &item.VersionID})
		} else {
			_, err = client.svc.DeleteObject(&s3.DeleteObjectInput{Bucket: &bucket, Key: &item.Name})
		}
	}

	return err
}
