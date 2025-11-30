package s3client

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// UploadBytesRequest represents a request to upload raw bytes to S3
type UploadBytesRequest struct {
	Data         []byte
	FileName     string
	ResourceType string // "campaign", "activity", "receipt", "donation/receipts"
	ResourceID   string // UUID of the resource
}

// UploadBytes uploads raw bytes to S3 and returns the public URL
func (c *Client) UploadBytes(ctx context.Context, req UploadBytesRequest) (*UploadResponse, error) {
	// Validate file type
	if err := validateFileType(req.FileName); err != nil {
		return nil, err
	}

	// Generate unique key
	key := generateS3Key(req.ResourceType, req.ResourceID, req.FileName)

	// Create a reader from bytes
	reader := bytes.NewReader(req.Data)

	// Prepare upload input
	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(getContentType(req.FileName)),
	}

	// Upload to S3
	_, err := c.s3Client.PutObject(ctx, uploadInput)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Generate public URL based on environment
	var url string
	if c.endpoint != "" {
		// LocalStack URL format: http://localhost:4566/bucket/key
		url = fmt.Sprintf("%s/%s/%s", c.endpoint, c.bucketName, key)
	} else {
		// AWS URL format: https://bucket.s3.amazonaws.com/key
		url = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.bucketName, key)
	}

	return &UploadResponse{
		URL:      url,
		Key:      key,
		FileName: req.FileName,
		Size:     int64(len(req.Data)),
	}, nil
}

// generateS3KeyForBytes creates a unique S3 key (duplicated to avoid import issues)
func generateS3KeyForBytes(resourceType, resourceID, filename string) string {
	ext := filepath.Ext(filename)
	timestamp := time.Now().Unix()
	uniqueID := uuid.New().String()[:8]
	return fmt.Sprintf("%s/%s/%d-%s%s", resourceType, resourceID, timestamp, uniqueID, ext)
}

// validateFileTypeForBytes checks if the file type is allowed (duplicated to avoid import issues)
func validateFileTypeForBytes(filename string) error {
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".pdf":  true,
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("file type %s not allowed. Allowed types: jpg, jpeg, png, gif, webp, pdf", ext)
	}

	return nil
}

// getContentTypeForBytes returns the MIME type for the file (duplicated to avoid import issues)
func getContentTypeForBytes(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".pdf":  "application/pdf",
	}

	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}

	return "application/octet-stream"
}

