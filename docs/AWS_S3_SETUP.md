# AWS S3 Configuration for Dona Tutti API

## Required S3 Bucket Configuration

### 1. Create S3 Bucket
```bash
# Replace 'donatutti' with your bucket name
aws s3 mb s3://donatutti --region us-east-1
```

### 2. Configure Bucket Policy for Public Read Access

Go to AWS S3 Console → Your Bucket → Permissions → Bucket Policy

**Add this policy** (replace `donatutti` with your bucket name):

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::donatutti/*"
        }
    ]
}
```

### 3. Configure CORS (Optional, for web frontend)

Go to AWS S3 Console → Your Bucket → Permissions → CORS

```json
[
    {
        "AllowedHeaders": ["*"],
        "AllowedMethods": ["GET", "PUT", "POST"],
        "AllowedOrigins": ["*"],
        "ExposeHeaders": []
    }
]
```

### 4. Environment Variables

Ensure these are set in your `.env` file:

```bash
AWS_REGION=us-east-1
AWS_S3_BUCKET=donatutti
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
```

### 5. IAM User Permissions

Your AWS user needs these permissions:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:DeleteObject",
                "s3:GetObject"
            ],
            "Resource": "arn:aws:s3:::donatutti/*"
        }
    ]
}
```

## Testing Upload

After configuration, test with:

```bash
curl -X POST \
  http://localhost:9999/api/campaigns/{campaign-id}/upload \
  -H 'Authorization: Bearer {jwt-token}' \
  -F 'file=@image.jpg'
```

## File Access URLs

Uploaded files will be accessible at:
```
https://donatutti.s3.amazonaws.com/campaign/{campaign-id}/{timestamp}-{uuid}.jpg
```

## Security Notes

- Files are publicly readable via the bucket policy
- Only authenticated admin users can upload files
- File types are validated (images and PDFs only)
- Unique file names prevent conflicts