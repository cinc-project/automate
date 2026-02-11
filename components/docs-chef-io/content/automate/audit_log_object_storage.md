+++
title = "Store Audit Logs in S3 or MinIO"

draft = false

gh_repo = "automate"

[menu]
  [menu.automate]
    title = "Audit Logs in S3/MinIO"
    parent = "automate/configuring_automate"
    identifier = "automate/configuring_automate/Audit Logs in S3 or MinIO"
    weight = 43
+++

You can configure Chef Automate to use an object storage backend (AWS S3 or MinIO) for audit log storage.
This is helpful when you want audit log data to live outside the local filesystem and to integrate with object-storage-based retention and access controls.

## Prerequisites

- Network connectivity from Chef Automate services to your S3/MinIO endpoint
- An existing bucket for audit logs
- Credentials for the bucket (or an IAM role / instance profile when using AWS)

## Configure audit log uploads

To start uploading audit logs, patch the Chef Automate configuration.

`[global.v1.audit.logging].enabled` controls whether the audit logging pipeline is enabled:

- When set to `true`, Chef Automate deploys and runs the `automate-fluent-bit` service.
- When set to `false`, Chef Automate removes that service (audit logs are not collected/uploaded).

1. Create a TOML file with the following content on the node running Chef Automate in a standalone deployment or on the bastion host in an Automate HA cluster.

### AWS S3 example

```toml
[global.v1.audit.logging]
enabled = true

[global.v1.audit.storage]
storage_type = "s3"
endpoint = "https://s3.amazonaws.com"
bucket = "<BUCKET_NAME>"
storage_region = "<AWS_REGION>"
path_prefix = "audit-logs/"

# If you are using an IAM role/instance profile, omit access_key/secret_key.
# If you are using static credentials, set both.
# access_key = "<ACCESS_KEY>"
# secret_key = "<SECRET_KEY>"

  [global.v1.audit.storage.ssl]
  enabled = true
  verify_ssl = true
```

Set the following values:

- `bucket`: The S3 bucket where audit logs are stored.
- `storage_region`: AWS region for the bucket (for example, `"us-east-1"`).
- `path_prefix`: Optional prefix inside the bucket.
- `access_key` and `secret_key`: Optional static credentials.

### MinIO example

```toml
[global.v1.audit.logging]
enabled = true

[global.v1.audit.storage]
storage_type = "minio"
endpoint = "http://minio.example.com:9000"
bucket = "<BUCKET_NAME>"
storage_region = "us-east-1"
path_prefix = "audit-logs/"
access_key = "<ACCESS_KEY>"
secret_key = "<SECRET_KEY>"

  [global.v1.audit.storage.ssl]
  # For http:// endpoints, set enabled=false.
  # For https:// endpoints, set enabled=true.
  enabled = false
  verify_ssl = false
  # For private CAs/self-signed certs, set enabled=true and provide a PEM-encoded CA certificate.
  # root_cert = """-----BEGIN CERTIFICATE-----
  # ...
  # -----END CERTIFICATE-----"""
```

Set the following values:

- `endpoint`: MinIO endpoint URL.
- `bucket`: The bucket where audit logs are stored.
- `access_key` and `secret_key`: MinIO credentials.
- `ssl.enabled` / `ssl.verify_ssl` / `ssl.root_cert`: TLS settings.

1. Patch the Chef Automate configuration:

```bash
sudo chef-automate config patch </PATH/TO/TOML/FILE>
```

After you patch the Automate configuration, Chef Automate starts running the audit log collector and uploads audit log data to the configured bucket.

If you do not set `[global.v1.audit.output]`, Chef Automate uses the defaults shown in Configure upload behavior.

## Configure local audit log rotation

Chef Automate rotates the local audit log file written by the load balancer.
To change the rotation size threshold, patch your Automate configuration.

1. Create a TOML file with the following content:

```toml
[global.v1.audit.input]
max_file_size = "10MB"
```

Set the following values:

- `max_file_size`: Maximum size of the local audit log file before rotation occurs.

Rotation behavior:

- Audit entries are written to `/hab/svc/automate-load-balancer/data/audit.log`.
- When `audit.log` exceeds `max_file_size`, it is rotated to `audit.1.log`.
- Older rotated files are shifted up (`audit.1.log` → `audit.2.log`, etc.).
- Chef Automate keeps up to 10 rotated files (`audit.1.log` through `audit.10.log`).

1. Patch the Chef Automate configuration:

```bash
sudo chef-automate config patch </PATH/TO/TOML/FILE>
```

## Configure upload behavior

These settings control how the collector uploads audit logs to S3/MinIO (object size splitting, multipart chunk size, and upload timeouts).

The collector (`automate-fluent-bit`) retries failed uploads indefinitely by default. This retry limit is not configurable via `[global.v1.audit]`.

If you do not set `[global.v1.audit.output]`, Chef Automate uses these defaults:

- `total_file_size = "100M"`
- `upload_chunk_size = "6M"`
- `upload_timeout = "10m"`

1. Create a TOML file with the following content:

```toml
[global.v1.audit.output]
total_file_size = "100M"
upload_chunk_size = "6M"
upload_timeout = "10m"

```

Set the following values:

- `total_file_size`: Total size threshold before the output is split into additional objects.
- `upload_chunk_size`: Multipart upload part size.
- `upload_timeout`: Upload timeout (minutes or hours).

1. Patch the Chef Automate configuration:

```bash
sudo chef-automate config patch </PATH/TO/TOML/FILE>
```

## Verify

To verify audit log uploads:

1. Confirm the configuration was applied:

    ```shell
    chef-automate config show | sed -n '/\[global\.v1\.audit\]/,/^\[/p'
    ```

1. Confirm new objects are being written to the configured bucket/prefix.

## Troubleshooting

- If uploads fail to MinIO with TLS enabled, verify the endpoint scheme (`http://` vs `https://`) matches the `ssl.enabled` setting.
- If you use a private CA for MinIO, provide `root_cert` and set `ssl.enabled = true`.
- If you use AWS IAM roles, omit `access_key` and `secret_key` to use the default AWS credential chain.

## How it works

When audit logging is enabled, Chef Automate runs `automate-fluent-bit` as the audit log collector.

- The load balancer writes audit entries to `/hab/svc/automate-load-balancer/data/audit.log` (and rotates the file based on `[global.v1.audit.input]`).
- Fluent Bit tails the active `audit.log` file and uploads matching log entries to S3/MinIO when `[global.v1.audit.storage]` is configured.
- Only mutating HTTP operations are written to the audit log and therefore stored in S3/MinIO: `POST`, `PUT`, `PATCH`, and `DELETE`.

## Global audit configuration reference

This section explains the fields under `[global.v1.audit]`.

### `[global.v1.audit.logging]`

- `enabled` (boolean)
  - Enables or disables the audit logging pipeline.
  - Default: `false`.
  - When `true`, Chef Automate deploys and runs the audit log collector (`automate-fluent-bit`).
  - When `false`, Chef Automate removes that service (audit logs are not collected/uploaded).

### `[global.v1.audit.async]`

These settings tune the asynchronous worker pool used by services that generate audit log files on-demand (for example, for the audit log retrieval APIs).

- `max_concurrent_workers` (integer)
  - Maximum number of concurrent workers processing audit-log generation/filter tasks.
  - Default: `4`.
  - Higher values can increase throughput, but also increase CPU/memory usage.

- `queue_size` (integer)
  - Maximum number of queued tasks waiting to be processed.
  - Default: `100`.
  - If the queue is full, new requests may be rejected until capacity is available.

- `multipart_chunk_size` (string)
  - Chunk size (part size) used for multipart operations during large file generation/uploads.
  - Default: `"10MB"`.
  - Format: `KB`, `MB`, or `GB` suffixes (for example, `"20MB"`).
  - Note: Use `MB/GB/KB` (for example, `"20MB"`), not `"20M"`.

### `[global.v1.audit.input]`

These settings control how audit logs are written/rotated locally on the Automate node and how the collector tails/buffers them.

- `max_file_size` (string)
  - Maximum size of the local audit log file before rotation occurs.
  - Default: `"10MB"`.
  - Supported formats include: `K`/`KB`, `M`/`MB`, `G`/`GB` (for example, `"100M"`, `"100MB"`, `"1G"`).
  - Rotation behavior:
    - The load balancer writes audit entries to `/hab/svc/automate-load-balancer/data/audit.log`.
    - When `audit.log` exceeds `max_file_size`, it is rotated to `audit.1.log`.
    - Older rotated files are shifted up (`audit.1.log` → `audit.2.log`, etc.).
    - Chef Automate keeps up to 10 rotated files (`audit.1.log` through `audit.10.log`).

- `refresh_interval` (string)
  - Polling interval (in seconds) for the `automate-fluent-bit` Tail input to detect new data and file rotation.
  - This value is used as Fluent Bit `Refresh_Interval`.
  - Default: `"5"`.
  - Example: `"10"`.

- `mem_buf_limit` (string)
  - In-memory buffer limit for the `automate-fluent-bit` Tail input while reading logs.
  - This value is used as Fluent Bit `Mem_Buf_Limit`.
  - Default: `"5MB"`.
  - Example: `"20MB"`.

### `[global.v1.audit.storage]`

These settings configure the object storage destination for audit logs (S3 or S3-compatible storage such as MinIO).

Chef Automate services use this configuration to create the S3/MinIO connection (endpoint/region), select the target bucket/prefix, and authenticate (static keys or AWS credential chain).

- `storage_type` (string)
  - Storage backend type. Must be `"s3"` or `"minio"`.
  - Default: `"s3"`.

- `endpoint` (string)
  - S3/MinIO endpoint URL.
  - Examples: `"https://s3.amazonaws.com"` (AWS), `"http://localhost:9000"` (local MinIO).

- `bucket` (string)
  - Bucket name where audit logs are stored.
  - Required: yes. Uploads (collector) do not run unless a bucket is configured.

- `storage_region` (string)
  - Region to use for S3 API calls.
  - For MinIO, a region value is still required by many S3 clients; `"us-east-1"` is commonly used.

- `path_prefix` (string)
  - Optional key prefix inside the bucket.
  - Use this to group audit logs under a specific path (for example, `"audit-logs/"`).

- `access_key` (string)
  - Access key for S3/MinIO.
  - For MinIO, this is typically required.
  - For AWS, you may omit static credentials if using the AWS default credential chain (IAM role / environment).
  - Default: empty.

- `secret_key` (string)
  - Secret key for S3/MinIO.
  - If you set `access_key`, you must also set `secret_key`.

#### `[global.v1.audit.storage.ssl]`

- `enabled` (boolean)
  - Enables TLS for the storage connection.
  - Set to `true` for `https://` endpoints; set to `false` for `http://` endpoints.
  - Default: `true`.

- `verify_ssl` (boolean)
  - Controls whether TLS certificates are verified.
  - Keep `true` in production when possible.
  - Default: `true`.

- `root_cert` (string)
  - Optional PEM-encoded CA certificate used to trust a private CA / self-signed certificate on your S3-compatible endpoint.
  - Not typically required for AWS S3.
  - Provide the certificate contents directly (not a file path) using TOML triple quotes (`"""`).
  - Example:

    ```toml
      [global.v1.audit.storage.ssl]
      root_cert = """-----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----"""
    ```
  - Default: empty.

### `[global.v1.audit.output]`

These settings control how the audit log collector uploads data to S3/MinIO (for example, how large each uploaded object should be, multipart chunk sizing, and upload timeouts).

- `total_file_size` (string)
  - Total size threshold before the S3/MinIO output is split into additional objects.
  - Default: `"100M"`.
  - Supported units: `M` or `G` only (MiB/GiB).
  - Minimum: `"1M"`.

- `upload_chunk_size` (string)
  - Multipart upload chunk size (part size) used when uploading.
  - Default: `"6M"`.
  - Supported units: `M` or `G` only (MiB/GiB).
  - Minimum: `"6M"`.

- `upload_timeout` (string)
  - Upload timeout.
  - Default: `"10m"`.
  - Supported units: minutes (`m`) or hours (`h`). Seconds are not supported.

## Audit log retrieval APIs

Chef Automate provides APIs (via the `user-settings-service`) to request and track asynchronous generation of audit logs for the current user.

### Request self audit logs (async)

- Method/Path: `GET /api/v0/audit/self/request`
- Authentication: Required (Chef Automate session)
- Access control: Self only (the service uses the authenticated user identity)
- Query parameters:
  - `from` (optional): Start time (RFC 3339 timestamp). Default: 3 hours ago.
  - `to` (optional): End time (RFC 3339 timestamp). Default: now.
  - `order` (optional): Sort order, `asc` or `desc`. Default: `desc`.
- Constraints:
  - The requested time range must be 30 days or less.

Example:

```shell
curl -sS \
  -H "api-token: $TOKEN" \
  "https://$FQDN/api/v0/audit/self/request?from=2025-11-10T00:00:00Z&to=2025-11-11T00:00:00Z&order=desc"
```

Response:

```json
{
  "request_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "status": "processing",
  "message": "Self audit log generation started. Use the request ID to check status and download when ready."
}
```

### Check request status

- Method/Path: `GET /api/v0/audit/status`
- Authentication: Required
- Access control: Self only (a user can only view their own requested logs)
- Query parameters:
  - `request_id` (optional):
    - If omitted, returns the status of the latest request for the current user.
    - If provided, returns the status for that specific request ID.

Examples:

Get the latest request status for the current user:

```shell
curl -sS \
  -H "api-token: $TOKEN" \
  "https://$FQDN/api/v0/audit/status"
```

Get status for a specific request ID:

```shell
curl -sS \
  -H "api-token: $TOKEN" \
  "https://$FQDN/api/v0/audit/status?request_id=f47ac10b-58cc-4372-a567-0e02b2c3d479"
```

Response fields:

- `request_id`: The UUID for the request.
- `status`: One of `processing`, `completed`, `error`, `not_found`.
- `file_size`: Present when `status` is `completed`. Human-readable size string (for example, `"10.7 KB"`).
- `download_url`: Present when `status` is `completed`. Use the returned URL/path to download the generated file.
  - The returned value may be an absolute URL (when an external FQDN is configured) or a relative path.
  - The download URL is formatted like `/api/v1/audit/download?request_id=<request_id>`.
- `error`: Present when `status` indicates an error. Possible values include `audit_disabled`, `request_not_found`, and `no_requested_logs`.
- `message`: Human-readable status/error details.

Example (completed):

```json
{
  "request_id": "a1c977e1-96a1-4a09-85f8-364721ff9f11",
  "status": "completed",
  "file_size": "10.7 KB",
  "download_url": "https://$FQDN/api/v1/audit/download?request_id=a1c977e1-96a1-4a09-85f8-364721ff9f11",
  "error": "",
  "message": ""
}
```
