package s3client

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Delete removes a file from S3 using the full URL or key
func (c *Client) Delete(ctx context.Context, urlOrKey string) error {
	key := extractKeyFromURL(urlOrKey, c.bucketName)
	if key == "" {
		return fmt.Errorf("invalid URL or key provided")
	}

	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	_, err := c.s3Client.DeleteObject(ctx, deleteInput)
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// DeleteByKey removes a file from S3 using the S3 key directly
func (c *Client) DeleteByKey(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	_, err := c.s3Client.DeleteObject(ctx, deleteInput)
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// extractKeyFromURL extracts the S3 key from a full S3 URL
// Supports both formats:
// - https://bucket-name.s3.amazonaws.com/key/path/file.jpg
// - https://s3.amazonaws.com/bucket-name/key/path/file.jpg
func extractKeyFromURL(urlOrKey string, bucketName string) string {
	// If it doesn't look like a URL, assume it's already a key
	if !strings.HasPrefix(urlOrKey, "http") {
		return urlOrKey
	}

	parsedURL, err := url.Parse(urlOrKey)
	if err != nil {
		return ""
	}

	// Format: https://bucket-name.s3.amazonaws.com/key/path/file.jpg
	if strings.HasPrefix(parsedURL.Host, bucketName+".s3") {
		// Remove leading slash from path
		return strings.TrimPrefix(parsedURL.Path, "/")
	}

	// Format: https://s3.amazonaws.com/bucket-name/key/path/file.jpg
	if strings.Contains(parsedURL.Host, "s3.amazonaws.com") {
		pathParts := strings.SplitN(strings.TrimPrefix(parsedURL.Path, "/"), "/", 2)
		if len(pathParts) == 2 && pathParts[0] == bucketName {
			return pathParts[1]
		}
	}

	return ""
}

// GetKeyFromCampaignImageURL is a helper function to extract S3 key from campaign image URL
// This is useful when updating campaign images to delete the old one
func (c *Client) GetKeyFromCampaignImageURL(imageURL string) string {
	return extractKeyFromURL(imageURL, c.bucketName)
}