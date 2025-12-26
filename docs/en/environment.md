---
outline: deep
---

# Configuration (Environment Variables)
All runtime configuration is provided via environment variables. System-level variables share the common prefix `AUTOCFT_`.

## System Variables
| Variable | Required | Default | Description                                         |
|----------|----------|---------|-----------------------------------------------------|
| AUTOCFT_CF_API_TOKEN | Yes      | - | Cloudflare API Token (Tunnel read/write).           |
| AUTOCFT_CF_ACCOUNT_ID | Yes      | - | Cloudflare Account ID.                              |
| AUTOCFT_CF_TUNNEL_ID | Yes      | - | Target Cloudflare Tunnel UUID.                      |
| AUTOCFT_BASEDIR | No       | /app/autocft | Working directory (history file / PocketBase data). |
| AUTOCFT_CRON | No       | */10 * * * * * | Sync cron (seconds enabled). Default: every 10s.    |
| AUTOCFT_ADMIN_EMAIL    | No       | autocft@cloudfogtech.ltd | Set admin account email(prepare for PocketBase)     |
| AUTOCFT_ADMIN_PASSWORD | No       | autocft@cloudfogtech#123 | Set admin account password(prepare for PocketBase)  |


## Default Origin Override Variables
Used as global fallbacks when a container does not specify the corresponding label.

| Variable | Type | Description |
|----------|------|-------------|
| AUTOCFT_ORIGIN_CONNECT_TIMEOUT | int(ms) | TCP connect timeout. |
| AUTOCFT_ORIGIN_DISABLE_CHUNKED_ENCODING | bool | Disable chunked transfer encoding. |
| AUTOCFT_ORIGIN_HTTP2_ORIGIN | bool | Attempt HTTP/2. |
| AUTOCFT_ORIGIN_HTTP_HEADER | string | Override Host header. |
| AUTOCFT_ORIGIN_KEEP_ALIVE_CONNECTIONS | int | Max idle connections. |
| AUTOCFT_ORIGIN_KEEP_ALIVE_TIME | int(s) | Idle connection timeout. |
| AUTOCFT_ORIGIN_NO_HAPPY_EYEBALLS | bool | Disable Happy Eyeballs (IPv4/IPv6 race). |
| AUTOCFT_ORIGIN_NO_TLS_VERIFY | bool | Skip TLS verification. |
| AUTOCFT_ORIGIN_ORIGIN_SERVER_NAME | string | Expected certificate hostname. |
| AUTOCFT_ORIGIN_PROXY_TYPE | string | "" or socks. |
| AUTOCFT_ORIGIN_TCP_KEEP_ALIVE | int(s) | TCP keepalive interval. |
| AUTOCFT_ORIGIN_TLS_TIMEOUT | int(s) | TLS handshake timeout. |

## Cron Expression
Uses robfig/cron (seconds field enabled) format: `sec min hour day month dow`.

## Security Recommendations
- Minimize API Token permissions (least privilege).
- Mount docker.sock read-only where possible.
- Monitor logs for unexpected changes.

## Next
Container labels: [/en/labels](/en/labels)

