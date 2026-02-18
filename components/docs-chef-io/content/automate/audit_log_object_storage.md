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

## Quick start

Use this section if you want the minimum configuration to start uploading audit logs to your bucket.

For field descriptions and defaults, see [Global audit configuration reference](#global-audit-configuration-reference).

### Minimum required configuration

At a minimum, you must:

- Enable audit logging (`[global.v1.audit.logging].enabled = true`)
- Configure storage (`[global.v1.audit.storage]`) with:
  - `storage_type` (`"s3"` for AWS S3 or `"minio"` for MinIO)
  - `endpoint`
  - `bucket`
  - `storage_region`
- Provide credentials:
  - For AWS S3, you can omit `access_key`/`secret_key` when using an IAM role/instance profile (AWS default credential chain)
  - For MinIO, set `access_key` and `secret_key`

For default values and optional snippets (including `root_cert`), see [Defaults and validation reference](#defaults-and-validation-reference).

AWS S3 (IAM role / instance profile):

```toml
[global.v1.audit]

  [global.v1.audit.logging]
    enabled = true

  [global.v1.audit.storage]
    storage_type = "s3"
    endpoint = "https://s3.amazonaws.com"
    bucket = "<BUCKET_NAME>"
    storage_region = "<AWS_REGION>"
```

MinIO (static credentials):

```toml
[global.v1.audit]

  [global.v1.audit.logging]
    enabled = true

  [global.v1.audit.storage]
    storage_type = "minio"
    endpoint = "http://minio.example.com:9000"
    bucket = "<BUCKET_NAME>"
    storage_region = "us-east-1"
    access_key = "<ACCESS_KEY>"
    secret_key = "<SECRET_KEY>"

    [global.v1.audit.storage.ssl]
      enabled = false
      verify_ssl = false
```

Patch the Chef Automate configuration:

```bash
sudo chef-automate config patch </PATH/TO/TOML/FILE>
```

{{< note >}}
In a high availability (HA) deployment, patch the config for the correct product:

```bash
chef-automate config patch automate.toml --automate
```

Shorthands for `--automate` are `--a2` and `-a`.

```bash
chef-automate config patch chefserver.toml --chef_server
```

Shorthands for `--chef_server` are `--cs` and `-c`.
{{< /note >}}

### Optional/advanced configuration

Once uploads work, you can optionally configure:

- `path_prefix` to group objects under a prefix inside the bucket
- TLS options under `[global.v1.audit.storage.ssl]` (enable/verify and `root_cert` for private CAs)
- Local log rotation and Tail input tuning under `[global.v1.audit.input]` (see Configure local audit log rotation)
- Upload batching and timeouts under `[global.v1.audit.output]` (see Configure upload behavior)
- Worker/queue tuning under `[global.v1.audit.async]`

For details on each field, see [Global audit configuration reference](#global-audit-configuration-reference).

## Configure audit log uploads

To start uploading audit logs, patch the Chef Automate configuration.

`[global.v1.audit.logging].enabled` controls whether the audit logging pipeline is enabled:

- When set to `true`, Chef Automate deploys and runs the `automate-fluent-bit` service.
- When set to `false`, Chef Automate removes that service (no audit log data is collected or uploaded).

{{< note >}}
If you set `[global.v1.audit.logging].enabled = true` but do not configure `[global.v1.audit.storage]`, Chef Automate will still write audit events locally to `/hab/svc/automate-load-balancer/data/audit.log` (and rotate them based on `[global.v1.audit.input]`), but nothing will be uploaded to S3/MinIO until storage is configured.
{{< /note >}}

1. Use one of the TOML examples in [Quick start](#quick-start) as a starting point, then adjust values for your environment.

{{< note >}}
For MinIO, the `endpoint` scheme must match the TLS settings:

- Use `http://...` with `ssl.enabled = false`
- Use `https://...` with `ssl.enabled = true`

If MinIO uses a private CA or self-signed certificate, set `ssl.enabled = true` and provide `ssl.root_cert` as PEM contents (not a file path).
{{< /note >}}

If you need details about a specific field, see [Global audit configuration reference](#global-audit-configuration-reference).

1. Patch the Chef Automate configuration (see Quick start for the command).

After you patch the Automate configuration, Chef Automate starts running the audit log collector and uploads audit log data to the configured bucket.

If you do not set `[global.v1.audit.output]`, Chef Automate uses the defaults shown in Configure upload behavior.

## Configure local audit log rotation

Chef Automate rotates the local audit log file written by the load balancer.
To change the rotation size threshold, patch your Automate configuration.

1. Create a TOML file with the following content:

```toml
[global.v1.audit]

  [global.v1.audit.input]
    max_file_size = "10MB"
```

Set the following values:

- `max_file_size`: Maximum size of the local audit log file before rotation occurs.

{{< note >}}
If you also set `[global.v1.audit.input].refresh_interval` or `[global.v1.audit.input].mem_buf_limit`, those values are passed through to Fluent Bit's Tail input (`Refresh_Interval` and `Mem_Buf_Limit`).

See the [Fluent Bit Tail input documentation](https://docs.fluentbit.io/manual/pipeline/inputs/tail).
{{< /note >}}

Rotation behavior:

- Audit entries are written to `/hab/svc/automate-load-balancer/data/audit.log`.
- When `audit.log` exceeds `max_file_size`, it is rotated to `audit.1.log`.
- Older rotated files are shifted up (`audit.1.log` → `audit.2.log`, etc.).
- Chef Automate keeps up to 10 rotated files (`audit.1.log` through `audit.10.log`).

1. Patch the Chef Automate configuration (see Quick start for the command).

## Configure upload behavior

These settings control how the collector uploads audit logs to S3/MinIO (object size splitting, multipart chunk size, and upload timeouts).

If you do not set `[global.v1.audit.output]`, Chef Automate uses these defaults:

- `total_file_size = "100M"`
- `upload_chunk_size = "6M"`
- `upload_timeout = "10m"`

1. Create a TOML file with the following content:

```toml
[global.v1.audit]

  [global.v1.audit.output]
    total_file_size = "100M"
    upload_chunk_size = "6M"
    upload_timeout = "10m"
```

{{< note >}}
These `[global.v1.audit.output]` settings control the Fluent Bit S3 output plugin behavior.
See the [Fluent Bit S3 output plugin documentation](https://docs.fluentbit.io/manual/pipeline/outputs/s3) for details and constraints.
{{< /note >}}

Set the following values:

- `total_file_size`: Total size threshold before the output is split into additional objects.
- `upload_chunk_size`: Multipart upload part size.
- `upload_timeout`: Upload timeout (minutes or hours).

1. Patch the Chef Automate configuration (see Quick start for the command).

## Complete patch example

You can patch all audit log settings in a single TOML file (logging, local rotation, S3/MinIO storage, TLS settings, and upload behavior).

For field descriptions and defaults, see [Global audit configuration reference](#global-audit-configuration-reference).

```toml
[global.v1.audit]

  [global.v1.audit.logging]
    enabled = true

  [global.v1.audit.input]
    max_file_size = "10MB"
    refresh_interval = "5"
    mem_buf_limit = "5MB"

  [global.v1.audit.async]
    # max_concurrent_workers = 4
    # queue_size = 100
    # multipart_chunk_size = "10MB"

  [global.v1.audit.storage]
    # Use "s3" for AWS S3 or "minio" for MinIO.
    storage_type = "minio"
    endpoint = "http://fqdn:9000"
    bucket = "audit-logs"
    storage_region = "us-east-1"
    path_prefix = "audit-logs/"
    access_key = "<ACCESS_KEY>"
    secret_key = "<SECRET_KEY>"

    [global.v1.audit.storage.ssl]
      # Set enabled=true for https:// endpoints; enabled=false for http:// endpoints.
      enabled = false
      verify_ssl = false
      # For private CAs/self-signed certs, provide a PEM-encoded CA certificate.
      # root_cert = """-----BEGIN CERTIFICATE-----
      # ...
      # -----END CERTIFICATE-----"""

  [global.v1.audit.output]
    total_file_size = "100M"
    upload_chunk_size = "6M"
    upload_timeout = "10m"
```

Patch the Chef Automate configuration (see Quick start for the command).

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

For copy/paste examples, see [Quick start](#quick-start) (minimum required configuration) and [Complete patch example](#complete-patch-example) (all common settings in one file).

### Defaults and validation reference

**View defaults in TOML format:**

```toml
[global.v1.audit]

  [global.v1.audit.logging]
    enabled = false

  [global.v1.audit.input]
    max_file_size = "10MB"
    refresh_interval = "60"
    mem_buf_limit = "5M"

  [global.v1.audit.async]
    max_concurrent_workers = 4
    queue_size = 100
    multipart_chunk_size = "10MB"

  [global.v1.audit.storage]
    storage_type = "s3"
    endpoint = "https://s3.amazonaws.com"
    storage_region = "us-east-1"
    path_prefix = ""

    [global.v1.audit.storage.ssl]
      enabled = false
      verify_ssl = false
      root_cert = ""

  [global.v1.audit.output]
    total_file_size = "12M"
    upload_chunk_size = "6M"
    upload_timeout = "10m"
```

**Configuration tables by section:**

|Field|Default|Validation|
|---|---|---|
|`enabled`|`false`|Must be `true` or `false`.|

#### `[global.v1.audit.async]`

|Field|Default|Validation|
|---|---|---|
|`max_concurrent_workers`|`4`|Higher values increase throughput but also CPU/memory usage.|
|`queue_size`|`100`|If full, new requests may be rejected.|
|`multipart_chunk_size`|`"10MB"`|Format: `KB`, `MB`, or `GB` suffixes (use `"20MB"`, not `"20M"`).|

#### `[global.v1.audit.input]`

|Field|Default|Validation|
|---|---|---|
|`max_file_size`|`"10MB"`|If set, must be a positive size with `K`, `M`, or `G` units (optional `B`, no spaces; for example, `10MB`) and must be ≥ 1 MiB.|
|`refresh_interval`|`"60"`|If set, must be a positive integer number of seconds.|
|`mem_buf_limit`|`"5M"`|If set, must be a positive value matching `^\d+M$` (capital `M`, no spaces, no `B` suffix).|

#### `[global.v1.audit.storage]`

|Field|Default|Validation|
|---|---|---|
|`storage_type`|`"s3"`|If set, must be `"s3"` or `"minio"` (cannot be empty).|
|`endpoint`|`"https://s3.amazonaws.com"`|Required when `bucket` is set.|
|`bucket`|—|If set, enables uploads. If omitted, storage is treated as not configured and `endpoint`/`storage_region` are not required.|
|`storage_region`|`"us-east-1"`|Required for `"s3"` when `bucket` is set. Optional for MinIO.|
|`path_prefix`|`""`|Optional; if set, must be non-empty.|
|`access_key`|`""`|For MinIO: typically required. For AWS: optional if using IAM role.|
|`secret_key`|`""`|Required if `access_key` is set.|

#### `[global.v1.audit.storage.ssl]`

|Field|Default|Validation|
|---|---|---|
|`enabled`|`false`|If `true`, `root_cert` must be set and non-empty.|
|`verify_ssl`|`false`|Allowed only when `enabled = true` and `root_cert` is set.|
|`root_cert`|`""`|Required when `enabled = true`; must be PEM-encoded and non-empty.|

#### `[global.v1.audit.output]`

|Field|Default|Min|Max|Validation|
|---|---|---|---|---|
|`total_file_size`|`"12M"`|`"1M"`|`"50G"`|If set, units must be `M` or `G` only and the value must be between 1M and 50G. Must also be ≥ 2× `upload_chunk_size`.|
|`upload_chunk_size`|`"6M"`|`"6M"`|`"50M"`|If set, units must be `M` or `G` only and the value must be between 6M and 50M.|
|`upload_timeout`|`"10m"`|—|—|If set, must be a positive duration with `s`, `m`, or `h` suffix (for example, `30s`, `10m`, `1h`).|

## Audit log retrieval APIs

Chef Automate provides APIs (via the `user-settings-service`) to request and track asynchronous generation of audit logs.

### Authentication

All audit APIs require authentication.

Bearer tokens (JWT) in the `Authorization` header work for all endpoints:

```bash
curl -sS \
  -H "Authorization: Bearer $TOKEN" \
  "https://$FQDN/api/v0/audit/..."
```

Alternatively, you may use the `api-token` header in some API contexts.

{{< note >}}
`api-token` authentication works with the Admin request API, Status API, and Download API.
The Self request API does not currently accept `api-token` authentication; use a bearer (JWT) token for `GET /api/v0/audit/self/request`.
{{< /note >}}

Example using an API token header:

```bash
curl -sS \
  -H "api-token: $API_TOKEN" \
  "https://$FQDN/api/v0/audit/admin/request"
```

### Request admin audit logs (async)

- Method/Path: `GET /api/v0/audit/admin/request`
- Authentication: Required (Chef Automate admin session)
- Access control: Admin only
- Query parameters (all optional):
  - `usernames` (optional): Comma-separated list of users to filter logs. If omitted, returns logs for all users.
  - `from` (optional): Start time (RFC 3339 timestamp). Default: 3 hours ago.
  - `to` (optional): End time (RFC 3339 timestamp). Default: now.
  - `order` (optional): Sort order, `asc` or `desc`. Default: `desc`.
- Constraints:
  - The requested time range (`to` - `from`) must be 30 days or less.
  - If you omit both `from` and `to`, the request defaults to the last 3 hours.

Example:

```shell
curl -sS \
  -H "Authorization: Bearer $TOKEN" \
  "https://$FQDN/api/v0/audit/admin/request?usernames=user1,user2&from=2025-11-10T00:00:00Z&to=2025-11-11T00:00:00Z&order=desc"
```

Response:

```json
{
  "request_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "status": "processing",
  "message": "Admin audit log generation started. Use the request ID to check status and download when ready."
}
```

### Request self audit logs (async)

- Method/Path: `GET /api/v0/audit/self/request`
- Authentication: Required (Chef Automate session)
- Access control: Self only (the service uses the authenticated user identity)
- Query parameters (all optional):
  - `from` (optional): Start time (RFC 3339 timestamp). Default: 3 hours ago.
  - `to` (optional): End time (RFC 3339 timestamp). Default: now.
  - `order` (optional): Sort order, `asc` or `desc`. Default: `desc`.
- Constraints:
  - The requested time range (`to` - `from`) must be 30 days or less.
  - If you omit both `from` and `to`, the request defaults to the last 3 hours.

Example:

```shell
curl -sS \
  -H "Authorization: Bearer $TOKEN" \
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

- Method/Path: `GET /api/v1/audit/status`
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
  -H "Authorization: Bearer $TOKEN" \
  "https://$FQDN/api/v1/audit/status"
```

Get status for a specific request ID:

```shell
curl -sS \
  -H "Authorization: Bearer $TOKEN" \
  "https://$FQDN/api/v1/audit/status?request_id=f47ac10b-58cc-4372-a567-0e02b2c3d479"
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

### Download audit logs

- Method/Path: `GET /api/v1/audit/download`
- Authentication: Required
- Access control: Self only (a user can only download audit logs they requested)
- Query parameters:
  - `request_id` (optional):
    - If omitted, returns the last requested audit log file for the current user.
    - If provided, returns the audit log file for that specific request ID.
- Returns: The audit log file in the format specified by the original request (typically JSON or CSV).

Each request generates a single downloadable file containing the audit logs for the full requested time range (up to 30 days).

Examples:

Download the latest audit log for the current user:

```shell
curl -sS \
  -H "Authorization: Bearer $TOKEN" \
  "https://$FQDN/api/v1/audit/download" > audit_logs.json
```

Download audit logs for a specific request ID:

```shell
curl -sS \
  -H "Authorization: Bearer $TOKEN" \
  "https://$FQDN/api/v1/audit/download?request_id=a1c977e1-96a1-4a09-85f8-364721ff9f11" > audit_logs.json
```
