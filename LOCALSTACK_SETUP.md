# LocalStack Setup for Dona Tutti API

Simple setup for simulating AWS S3 locally using LocalStack.

## Quick Start

```bash
# Start LocalStack and create S3 bucket
make localstack

# Start development with LocalStack
make dev-local
```

That's it! Your API will now use LocalStack S3 instead of AWS.

## Commands

### Using Make

```bash
make localstack        # Start LocalStack and create S3 bucket
make localstack-stop   # Stop LocalStack
make localstack-status # Check if LocalStack is running
make dev-local         # Start full development environment with LocalStack
```

### Using Script Directly

```bash
./localstack.sh setup   # Start LocalStack and create bucket
./localstack.sh stop    # Stop LocalStack
./localstack.sh status  # Check status
```

## Configuration

The S3 bucket name is: `dona-tutti-s3`

To use LocalStack, the API needs:
- `LOCALSTACK_ENDPOINT=http://localhost:4566`
- `AWS_S3_BUCKET=dona-tutti-s3`
- `AWS_ACCESS_KEY_ID=test`
- `AWS_SECRET_ACCESS_KEY=test`

These are automatically set when using `make dev-local`.

## Testing

Test file upload with LocalStack running:

```bash
# Upload a file via API (requires authentication)
curl -X POST http://localhost:9999/api/campaigns/{id}/upload \
  -H "Authorization: Bearer {token}" \
  -F "file=@test.jpg"

# Check uploaded files
aws --endpoint-url=http://localhost:4566 s3 ls s3://dona-tutti-s3/ --recursive
```

## Troubleshooting

If LocalStack won't start:
```bash
# Stop everything and try again
docker-compose --profile localstack down
docker volume rm dona_tutti_api_localstack_data
make localstack
```