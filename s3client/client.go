package s3client

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3Client   *s3.Client
	bucketName string
}

type Config struct {
	Region     string
	BucketName string
	AccessKey  string
	SecretKey  string
}

func NewClient() (*Client, error) {
	cfg := Config{
		Region:     getEnvOrDefault("AWS_REGION", "us-east-1"),
		BucketName: getEnvOrDefault("AWS_S3_BUCKET", "dona-tutti-files"),
		AccessKey:  os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:  os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("AWS credentials not provided")
	}

	if cfg.BucketName == "" {
		return nil, fmt.Errorf("AWS S3 bucket name not provided")
	}

	// Load AWS config
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with optional LocalStack endpoint
	s3Options := func(o *s3.Options) {
		// Check for LocalStack endpoint
		if endpoint := os.Getenv("LOCALSTACK_ENDPOINT"); endpoint != "" {
			o.BaseEndpoint = &endpoint
			o.UsePathStyle = true // LocalStack requires path-style URLs
		}
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(awsConfig, s3Options)

	return &Client{
		s3Client:   s3Client,
		bucketName: cfg.BucketName,
	}, nil
}

func (c *Client) GetBucketName() string {
	return c.bucketName
}

func (c *Client) GetS3Client() *s3.Client {
	return c.s3Client
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
