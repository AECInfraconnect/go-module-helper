package minio

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

// generateObjectName generates a unique object name with timestamp and random number.
// Format: {foldername}/{YYYYMMDD}_{id}_{random}.{extension}
//
// Internal helper function used by GenerateObjectName methods.
func generateObjectName(foldername string, id string, extension string) string {
	date := time.Now().Format("20060102")
	generateNumber := fmt.Sprintf("%010d", rand.Intn(10000000000))
	extension = strings.TrimPrefix(extension, ".")

	if foldername != "" && !strings.HasSuffix(foldername, "/") {
		foldername += "/"
	}

	return fmt.Sprintf("%s%s_%s_%s.%s", foldername, date, id, generateNumber, extension)
}

// GenerateObjectName generates a unique object name for file storage.
// Combines folder path, current date, ID, and random number for uniqueness.
//
// Parameters:
//   - foldername: Folder path (automatically adds trailing slash if missing)
//   - id: Identifier (e.g., user ID, document ID)
//   - filename: File extension (with or without dot)
//
// Example:
//
//	objectName := minio.GenerateObjectName("uploads", "user123", ".jpg")
//	// Returns: "uploads/20260113_user123_1234567890.jpg"
func GenerateObjectName(foldername string, id string, filename string) string {
	return generateObjectName(foldername, id, filename)
}

// GenerateObjectName generates a unique object name for file storage (method version).
//
// Example:
//
//	objectName := client.GenerateObjectName("uploads", "user123", "jpg")
func (c *Client) GenerateObjectName(foldername string, id string, filename string) string {
	return generateObjectName(foldername, id, filename)
}

// GetObjectnameFromURL extracts bucket name and object path from a full URL.
// Parses MinIO URLs and returns the bucket and object name components.
//
// Parameters:
//   - link: Full URL to the object (e.g., "https://minio.example.com/my-bucket/path/to/file.jpg")
//
// Returns:
//   - bucket: Extracted bucket name
//   - objectName: Extracted object path
//
// Example:
//
//	bucket, objectName := client.GetObjectnameFromURL("https://minio.example.com/my-bucket/uploads/file.jpg")
//	// bucket = "my-bucket", objectName = "uploads/file.jpg"
func (c *Client) GetObjectnameFromURL(link string) (string, string) {
	var bucket string
	var pathImageRegex = regexp.MustCompile(`\/[a-z\-]+\/`)

	uri, err := url.Parse(link)
	if err != nil {
		return bucket, link
	}

	objectName := uri.Path
	if objectName != "" {
		if pathImageRegex.MatchString(objectName) && uri.Scheme != "" {
			loc := pathImageRegex.FindStringIndex(objectName)
			if len(loc) > 0 {
				bucket = strings.Trim(objectName[loc[0]:loc[1]], "/")
				objectName = objectName[loc[1]:]
			}
		}
	}
	if strings.Contains(objectName, "?") {
		spl := strings.Split(objectName, "?")
		objectName = spl[0]
	}
	if objectName == "" {
		return bucket, link
	}

	return bucket, strings.Trim(objectName, "/")
}

// RemoveObject deletes an object from the specified bucket.
//
// Parameters:
//   - bucketName: Name of the bucket containing the object
//   - objectName: Path to the object to delete
//
// Example:
//
//	err := client.RemoveObject("my-bucket", "uploads/file.jpg")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) RemoveObject(bucketName string, objectName string) error {
	if err := c.GetClient().RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	return nil
}

// RemoveObjectWithContext deletes an object with custom context for cancellation control.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - bucketName: Name of the bucket containing the object
//   - objectName: Path to the object to delete
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	err := client.RemoveObjectWithContext(ctx, "my-bucket", "uploads/file.jpg")
func (c *Client) RemoveObjectWithContext(ctx context.Context, bucketName string, objectName string) error {
	if err := c.GetClient().RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	return nil
}
