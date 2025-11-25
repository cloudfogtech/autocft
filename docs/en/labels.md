---
outline: deep
---

# Container Labels Reference
Docker labels starting with `autocft.` control creation of ingress rules.

## Required Labels
| Label | Description |
|-------|-------------|
| autocft.enabled | Set to `true` to be managed. |
| autocft.hostname | Public hostname to expose. |
| autocft.service | Upstream service address, e.g. `http://app:8080`, `tcp://redis:6379`. |

## Optional Labels
| Label | Type | Description | Fallback |
|-------|------|-------------|----------|
| autocft.path | string | URL path prefix (default `/`). | `/` |
| autocft.origin.connect-timeout | int | Connect timeout ms. | ENV |
| autocft.origin.disable-chunked-encoding | bool | Disable chunked encoding. | ENV |
| autocft.origin.http2-origin | bool | Force HTTP/2. | ENV |
| autocft.origin.http-host-header | string | Override Host header. | ENV |
| autocft.origin.keep-alive-connections | int | Max idle connections. | ENV |
| autocft.origin.keep-alive-timeout | int | Keepalive idle timeout s. | ENV |
| autocft.origin.no-happy-eyeballs | bool | Disable dual‑stack race. | ENV |
| autocft.origin.no-tls-verify | bool | Skip TLS verification. | ENV |
| autocft.origin.origin-server-name | string | Certificate host name. | ENV |
| autocft.origin.proxy-type | string | Blank or socks. | ENV |
| autocft.origin.tcp-keep-alive | int | TCP keepalive s. | ENV |
| autocft.origin.tls-timeout | int | TLS handshake timeout s. | ENV |

## Example
```yaml
services:
  api:
    image: ghcr.io/example/api:latest
    labels:
      - autocft.enabled=true
      - autocft.hostname=api.example.com
      - autocft.service=http://api:8080
  web:
    image: ghcr.io/example/web:latest
    labels:
      - autocft.enabled=true
      - autocft.hostname=app.example.com
      - autocft.service=http://web:3000
      - autocft.origin.no-tls-verify=true
```

## Collision Rules
Only one entry per hostname is kept; avoid duplicates. If a hostname conflicts with a web‑added entry it will be recognized as web‑managed on first sync and preserved.

## Removal
Remove or set `autocft.enabled=false`; on next sync the entry is deleted (if not web‑managed).
