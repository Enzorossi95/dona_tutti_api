#!/bin/bash

# Simple LocalStack Setup for Dona Tutti API

LOCALSTACK_ENDPOINT="http://localhost:4566"
BUCKET_NAME="donatutti"

setup() {
    echo "Starting LocalStack..."
    docker-compose --profile localstack up -d localstack
    
    echo "Waiting for LocalStack to be ready..."
    sleep 5
    
    # Check if LocalStack is running
    if ! curl -s "${LOCALSTACK_ENDPOINT}/_localstack/health" > /dev/null 2>&1; then
        echo "Error: LocalStack failed to start"
        exit 1
    fi
    
    echo "Creating S3 bucket: ${BUCKET_NAME}..."
    aws --endpoint-url="${LOCALSTACK_ENDPOINT}" \
        s3api create-bucket \
        --bucket "${BUCKET_NAME}" \
        --region us-east-1 2>/dev/null || echo "Bucket already exists"
    
    echo "Configuring bucket policy for public read access..."
    aws --endpoint-url="${LOCALSTACK_ENDPOINT}" \
        s3api put-bucket-policy \
        --bucket "${BUCKET_NAME}" \
        --policy '{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Sid": "PublicReadGetObject",
                    "Effect": "Allow",
                    "Principal": "*",
                    "Action": "s3:GetObject",
                    "Resource": "arn:aws:s3:::'"${BUCKET_NAME}"'/*"
                }
            ]
        }' 2>/dev/null || echo "Failed to set bucket policy"
    
    echo "Configuring CORS for bucket..."
    aws --endpoint-url="${LOCALSTACK_ENDPOINT}" \
        s3api put-bucket-cors \
        --bucket "${BUCKET_NAME}" \
        --cors-configuration '{
            "CORSRules": [
                {
                    "AllowedOrigins": ["*"],
                    "AllowedMethods": ["GET", "HEAD"],
                    "AllowedHeaders": ["*"],
                    "MaxAgeSeconds": 3000
                }
            ]
        }' 2>/dev/null || echo "Failed to set CORS configuration"
    
    echo "✅ LocalStack is ready!"
    echo "Endpoint: ${LOCALSTACK_ENDPOINT}"
    echo "Bucket: ${BUCKET_NAME}"
    echo "Bucket Policy: Public read access enabled"
    echo "CORS: Configured for GET/HEAD methods"
}

stop() {
    echo "Stopping LocalStack..."
    docker-compose --profile localstack down
    echo "LocalStack stopped"
}

status() {
    if curl -s "${LOCALSTACK_ENDPOINT}/_localstack/health" > /dev/null 2>&1; then
        echo "✅ LocalStack is running"
        echo "Buckets:"
        aws --endpoint-url="${LOCALSTACK_ENDPOINT}" s3 ls 2>/dev/null || echo "No buckets"
    else
        echo "❌ LocalStack is not running"
    fi
}

case "$1" in
    setup)
        setup
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {setup|stop|status}"
        echo ""
        echo "  setup  - Start LocalStack and create S3 bucket"
        echo "  stop   - Stop LocalStack"
        echo "  status - Check if LocalStack is running"
        exit 1
        ;;
esac