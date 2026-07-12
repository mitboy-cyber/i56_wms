#!/bin/sh
# I56 Framework — MinIO Init Script
# Creates required buckets on startup (runs via docker-entrypoint)

echo "[minio-init] Waiting for MinIO to be ready..."

# Wait for MinIO to be ready
until mc ready local 2>/dev/null; do
  sleep 1
done

echo "[minio-init] MinIO is ready. Creating buckets..."

# Create buckets
mc mb local/i56-uploads --ignore-existing
mc mb local/i56-backups --ignore-existing
mc mb local/i56-logs --ignore-existing
mc mb local/i56-static --ignore-existing

# Set bucket policies
mc anonymous set download local/i56-uploads
mc anonymous set download local/i56-static

# Set lifecycle rules for logs bucket (expire after 90 days)
mc ilm import local/i56-logs <<EOF
{
  "Rules": [
    {
      "ID": "expire-after-90-days",
      "Status": "Enabled",
      "Expiration": {
        "Days": 90
      }
    }
  ]
}
EOF

echo "[minio-init] Buckets created: i56-uploads, i56-backups, i56-logs, i56-static"
echo "[minio-init] Initialization complete."
