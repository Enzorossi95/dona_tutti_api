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
    
    echo "✅ LocalStack is ready!"
    echo "Endpoint: ${LOCALSTACK_ENDPOINT}"
    echo "Bucket: ${BUCKET_NAME}"
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