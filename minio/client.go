// Package minio provides a wrapper around MinIO Go client with simplified operations
// for object storage management including bucket operations, file uploads, and object manipulation.
package minio

import (
	"github.com/minio/minio-go/v7"
	credentialsv7 "github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	// MINIO_DEFAULT_REGION is the default AWS region used for MinIO buckets
	MINIO_DEFAULT_REGION = "ap-southeast-1"
)

// Client wraps the MinIO client with additional configuration and helper methods.
// It maintains connection details and provides convenient methods for common operations.
type Client struct {
	MinioClient    *minio.Client // The underlying MinIO client instance
	MinioEndPoint  string        // MinIO server endpoint (e.g., "localhost:9000")
	MinioAccessKey string        // Access key for authentication
	MinioSecretKey string        // Secret key for authentication
	MinioSSL       bool          // Whether to use SSL/TLS for connections
	Region         string        // AWS region for the MinIO server
}

// NewMinio creates and initializes a new MinIO client with the provided credentials.
// Returns an error if the connection cannot be established.
//
// Parameters:
//   - endpoint: MinIO server endpoint (e.g., "localhost:9000")
//   - access: Access key ID for authentication
//   - secret: Secret access key for authentication
//   - ssl: Whether to use HTTPS (true) or HTTP (false)
//   - region: AWS region (use MINIO_DEFAULT_REGION if empty)
//
// Example:
//
//	client, err := minio.NewMinio("localhost:9000", "minioadmin", "minioadmin", false, "ap-southeast-1")
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewMinio(endpoint string, access string, secret string, ssl bool, region string) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentialsv7.NewStaticV4(access, secret, ""),
		Secure: ssl,
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	c := &Client{
		MinioClient:    minioClient,
		MinioEndPoint:  endpoint,
		MinioAccessKey: access,
		MinioSecretKey: secret,
		MinioSSL:       ssl,
		Region:         region,
	}

	return c, nil
}

// GetMinioURI constructs and returns the full URI for the MinIO server.
// Returns HTTP or HTTPS URI based on SSL configuration.
//
// Example:
//
//	uri := client.GetMinioURI()  // Returns "https://localhost:9000" or "http://localhost:9000"
func (c *Client) GetMinioURI() string {
	var minioEndURI string
	if c.MinioSSL {
		minioEndURI = "https://"
	} else {
		minioEndURI = "http://"
	}
	minioEndURI = minioEndURI + c.MinioEndPoint
	return minioEndURI
}

// GetClient returns the underlying MinIO client instance.
// Use this to access advanced MinIO operations not wrapped by this package.
//
// Example:
//
//	minioClient := client.GetClient()
//	// Use minioClient for advanced operations
func (c *Client) GetClient() *minio.Client {
	return c.MinioClient
}

// GetEndPoint returns the configured MinIO server endpoint.
//
// Example:
//
//	endpoint := client.GetEndPoint()  // Returns "localhost:9000"
func (c *Client) GetEndPoint() string {
	return c.MinioEndPoint
}
