package s3client

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploadRequest struct {
	File         multipart.File
	Header       *multipart.FileHeader
	ResourceType string // "campaign", "activity", "receipt"
	ResourceID   string // UUID of the resource
}

type UploadResponse struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
}

// Upload uploads a file to S3 and returns the public URL
func (c *Client) Upload(ctx context.Context, req UploadRequest) (*UploadResponse, error) {
	// Validate file type
	if err := validateFileType(req.Header.Filename); err != nil {
		return nil, err
	}

	// Generate unique key
	key := generateS3Key(req.ResourceType, req.ResourceID, req.Header.Filename)

	// Prepare upload input
	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		Body:        req.File,
		ContentType: aws.String(getContentType(req.Header.Filename)),
		// Note: Public access should be configured via bucket policy instead of ACL
	}

	// Upload to S3
	_, err := c.s3Client.PutObject(ctx, uploadInput)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Generate public URL
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.bucketName, key)

	return &UploadResponse{
		URL:      url,
		Key:      key,
		FileName: req.Header.Filename,
		Size:     req.Header.Size,
	}, nil
}

// generateS3Key creates a unique S3 key for the file
func generateS3Key(resourceType, resourceID, filename string) string {
	// Extract file extension
	ext := filepath.Ext(filename)

	// Generate timestamp for uniqueness
	timestamp := time.Now().Unix()

	// Generate unique identifier
	uniqueID := uuid.New().String()[:8]

	// Create key: resourceType/resourceID/timestamp-uniqueID.ext
	return fmt.Sprintf("%s/%s/%d-%s%s", resourceType, resourceID, timestamp, uniqueID, ext)
}

// validateFileType checks if the file type is allowed
func validateFileType(filename string) error {
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".pdf":  true, // For receipts
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("file type %s not allowed. Allowed types: jpg, jpeg, png, gif, webp, pdf", ext)
	}

	return nil
}

// getContentType returns the MIME type for the file
func getContentType(filename string) string {
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
