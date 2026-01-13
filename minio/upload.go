package minio

import (
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go/v7"
)

// CreateImageFile saves an image.Image to disk in the specified format.
// Supports JPEG and PNG formats.
//
// Parameters:
//   - img: Image data to save
//   - format: Image format ("jpg", "jpeg", or "png")
//   - filename: Name for the saved file (unused parameter)
//   - desc: Destination file path
//
// Example:
//
//	err := minio.CreateImageFile(img, "jpg", "photo", "/tmp/photo.jpg")
func CreateImageFile(img image.Image, format string, filename string, desc string) error {
	file, err := os.Create(desc)
	if err != nil {
		return err
	}
	switch format {
	case "jpg":
		if err := jpeg.Encode(file, img, nil); err != nil {
			return err
		}
	case "jpeg":
		if err := jpeg.Encode(file, img, nil); err != nil {
			return err
		}
	case "png":
		if err := png.Encode(file, img); err != nil {
			return err
		}
	}
	defer file.Close()

	return nil
}

// UploadMultipartFile uploads a file from HTTP multipart form data to MinIO.
// Automatically extracts content type and size from the file header.
//
// Parameters:
//   - bucketName: Target bucket name
//   - objectName: Destination object path
//   - file: Multipart file header from HTTP request
//
// Example:
//
//	// In Gin handler
//	file, _ := c.FormFile("upload")
//	err := client.UploadMultipartFile("my-bucket", "uploads/file.jpg", file)
func (c *Client) UploadMultipartFile(bucketName string, objectName string, file *multipart.FileHeader) (err error) {
	contentType := file.Header.Get("Content-Type")
	size := file.Size

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	if _, err = c.GetClient().PutObject(context.Background(), bucketName, objectName, src, size, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return err
	}
	return nil
}

// UploadMultipartFileWithContext uploads a file from HTTP multipart form data with custom context.
//
// Example:
//
//	ctx := context.Background()
//	err := client.UploadMultipartFileWithContext(ctx, "my-bucket", "uploads/file.jpg", file)
func (c *Client) UploadMultipartFileWithContext(ctx context.Context, bucketName string, objectName string, file *multipart.FileHeader) (err error) {
	contentType := file.Header.Get("Content-Type")
	size := file.Size

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	if _, err = c.GetClient().PutObject(ctx, bucketName, objectName, src, size, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return err
	}
	return nil
}

// UploadFileWithReader uploads data from an io.Reader to MinIO.
// Useful for uploading generated content or streaming data.
//
// Parameters:
//   - bucketName: Target bucket name
//   - objectName: Destination object path
//   - reader: Data source (io.Reader)
//   - size: Total size of data in bytes
//   - contentType: MIME type (e.g., "image/jpeg", "application/pdf")
//   - contentEncoding: Content encoding (e.g., "UTF-8", "gzip")
//
// Example:
//
//	data := bytes.NewReader([]byte("file content"))
//	err := client.UploadFileWithReader("my-bucket", "file.txt", data, int64(len("file content")), "text/plain", "UTF-8")
func (c *Client) UploadFileWithReader(bucketName string, objectName string, reader io.Reader, size int64, contentType string, contentEncoding string) (err error) {
	if _, err = c.GetClient().PutObject(context.Background(), bucketName, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType, ContentEncoding: contentEncoding}); err != nil {
		return err
	}
	return nil
}

// UploadFileWithReaderWithContext uploads data from an io.Reader with custom context.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	err := client.UploadFileWithReaderWithContext(ctx, "my-bucket", "file.txt", reader, size, "text/plain", "UTF-8")
func (c *Client) UploadFileWithReaderWithContext(ctx context.Context, bucketName string, objectName string, reader io.Reader, size int64, contentType string, contentEncoding string) (err error) {
	if _, err = c.GetClient().PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType, ContentEncoding: contentEncoding}); err != nil {
		return err
	}
	return nil
}

// UploadFromFile uploads a file from local filesystem to MinIO.
// Combines folder path and filename to create the object path.
//
// Parameters:
//   - bucketName: Target bucket name
//   - foldername: Folder path in the bucket
//   - pathFile: Local file path to upload
//   - filename: Destination filename
//
// Example:
//
//	err := client.UploadFromFile("my-bucket", "uploads", "/tmp/photo.jpg", "photo.jpg")
func (c *Client) UploadFromFile(bucketName string, foldername string, pathFile string, filename string) error {
	objectName := foldername + "/" + filename

	src, err := os.Open(pathFile)
	if err != nil {
		return err
	}

	defer src.Close()

	if _, err := c.GetClient().FPutObject(context.Background(), bucketName, objectName, pathFile, minio.PutObjectOptions{}); err != nil {
		return err
	}
	return nil
}

// UploadFromFileWithContact uploads a file from local filesystem with custom context.
// Note: Function name has typo "Contact" instead of "Context" - kept for backward compatibility.
//
// Example:
//
//	ctx := context.Background()
//	err := client.UploadFromFileWithContact(ctx, "my-bucket", "uploads", "/tmp/file.jpg", "file.jpg")
func (c *Client) UploadFromFileWithContact(ctx context.Context, bucketName string, foldername string, pathFile string, filename string) error {
	objectName := foldername + "/" + filename

	src, err := os.Open(pathFile)
	if err != nil {
		return err
	}

	defer src.Close()

	if _, err := c.GetClient().FPutObject(ctx, bucketName, objectName, pathFile, minio.PutObjectOptions{}); err != nil {
		return err
	}
	return nil
}

// UploadFromFilePDF uploads a PDF file from local filesystem with proper content type.
// Automatically sets content type to "application/pdf" and encoding to "UTF-8".
//
// Parameters:
//   - bucketName: Target bucket name
//   - foldername: Folder path in the bucket
//   - pathFile: Local PDF file path to upload
//   - filename: Destination filename
//
// Example:
//
//	err := client.UploadFromFilePDF("my-bucket", "documents", "/tmp/report.pdf", "report.pdf")
func (c *Client) UploadFromFilePDF(bucketName string, foldername string, pathFile string, filename string) error {
	objectName := foldername + "/" + filename

	src, err := os.Open(pathFile)
	if err != nil {
		return err
	}

	defer src.Close()

	if _, err := c.GetClient().FPutObject(context.Background(), bucketName, objectName, pathFile, minio.PutObjectOptions{ContentType: "application/pdf", ContentEncoding: "UTF-8"}); err != nil {
		return err
	}
	return nil
}

// UploadFromFilePDFWithContext uploads a PDF file with custom context.
// Automatically sets content type to "application/pdf" and encoding to "UTF-8".
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
//	defer cancel()
//	err := client.UploadFromFilePDFWithContext(ctx, "my-bucket", "documents", "/tmp/report.pdf", "report.pdf")
func (c *Client) UploadFromFilePDFWithContext(ctx context.Context, bucketName string, foldername string, pathFile string, filename string) error {
	objectName := foldername + "/" + filename

	src, err := os.Open(pathFile)
	if err != nil {
		return err
	}

	defer src.Close()

	if _, err := c.GetClient().FPutObject(ctx, bucketName, objectName, pathFile, minio.PutObjectOptions{ContentType: "application/pdf", ContentEncoding: "UTF-8"}); err != nil {
		return err
	}
	return nil
}
