package minio

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
)

// CreateBucket creates a new bucket with the specified name and region.
// Uses MINIO_DEFAULT_REGION if region parameter is empty.
//
// Parameters:
//   - bucketName: Name of the bucket to create
//   - region: AWS region for the bucket (use empty string for default)
//
// Example:
//
//	err := client.CreateBucket("my-bucket", "ap-southeast-1")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) CreateBucket(bucketName string, region string) error {
	if region == "" {
		region = MINIO_DEFAULT_REGION
	}
	if err := c.GetClient().MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
		Region: region,
	}); err != nil {
		return err
	}
	return nil
}

// CreateBucketWithContext creates a new bucket with custom context for cancellation and timeout control.
// Uses MINIO_DEFAULT_REGION if region parameter is empty.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - bucketName: Name of the bucket to create
//   - region: AWS region for the bucket (use empty string for default)
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	err := client.CreateBucketWithContext(ctx, "my-bucket", "")
func (c *Client) CreateBucketWithContext(ctx context.Context, bucketName string, region string) error {
	if region == "" {
		region = MINIO_DEFAULT_REGION
	}
	if err := c.GetClient().MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: region,
	}); err != nil {
		return err
	}
	return nil
}

// ExistBucket checks if a bucket exists.
// Returns true if the bucket exists, false otherwise.
//
// Example:
//
//	exists, err := client.ExistBucket("my-bucket")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if exists {
//	    fmt.Println("Bucket exists")
//	}
func (c *Client) ExistBucket(bucketName string) (bool, error) {
	exists, err := c.GetClient().BucketExists(context.Background(), bucketName)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ExistBucketWithContext checks if a bucket exists with custom context.
// Returns true if the bucket exists, false otherwise.
//
// Example:
//
//	ctx := context.Background()
//	exists, err := client.ExistBucketWithContext(ctx, "my-bucket")
func (c *Client) ExistBucketWithContext(ctx context.Context, bucketName string) (bool, error) {
	exists, err := c.GetClient().BucketExists(ctx, bucketName)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SetBucketPublicPolicy sets a public read policy for the bucket.
// Requires a policy template file at "./policy/policy_public.json".
// The policy allows public read access to all objects in the bucket.
//
// Parameters:
//   - bucketName: Name of the bucket to make public
//
// Example:
//
//	err := client.SetBucketPublicPolicy("my-bucket")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) SetBucketPublicPolicy(bucketName string) error {
	var buf bytes.Buffer
	bu, _ := os.ReadFile("./policy/policy_public.json")
	t, err := template.New("policy").Parse(string(bu))
	if err != nil {
		return err
	}

	if err := t.Execute(&buf, bucketName); err != nil {
		return err
	}

	policy := buf.String()
	if err := c.GetClient().SetBucketPolicy(context.Background(), bucketName, policy); err != nil {
		return err
	}
	log.Println("create bucket with policy success")
	fmt.Println(policy)
	return nil
}
