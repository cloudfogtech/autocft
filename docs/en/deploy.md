---
outline: deep
---

# Deployment Guide

Auto Cloudflare Tunnel (AutoCFT) automatically synchronizes Docker container labels into Cloudflare Tunnel ingress rules. You only need to add `autocft.*` labels to containers; AutoCFT periodically diffs and updates the Tunnel configuration.

## Workflow Overview
1. Read system default config from environment variables (prefix `AUTOCFT_`).
2. Enumerate running containers and parse labels beginning with `autocft.`.
3. Validate required labels: `autocft.enabled`, `autocft.hostname`, `autocft.service`.
4. Call Cloudflare API to load current Tunnel ingress; preserve entries not managed by labels (web‑configured entries).
5. Merge + generate update list; if changes exist call API to update.
6. Write current label‑managed ingress list to history file `latest.json` for next diff.

## Prerequisites
- Domain active and managed in Cloudflare (zone present).
- A created Cloudflare Tunnel (`cloudflared tunnel create <name>` or via Zero Trust dashboard). Need the Tunnel UUID.
- One Cloudflare API Token with minimum permissions:
  - Account level: Cloudflare Tunnel Read & Edit
  - (Multi‑zone hostname scenario) Zone level: DNS Read (optional, recommended)
- AutoCFT container can access host `/var/run/docker.sock` (read‑only is sufficient).

## Required Values
| Item | How to Obtain |
|------|---------------|
| CF_API_TOKEN | Dashboard -> My Profile -> API Tokens (create custom token) |
| CF_ACCOUNT_ID | Dashboard URL `/accounts/<id>/` or API Tokens page |
| CF_TUNNEL_ID | `cloudflared tunnel list` or Tunnel detail page UUID |

## docker-compose Example
```yaml
services:
  autocft:
    image: cloudfogtech/autocft:latest
    container_name: autocft
    environment:
      - AUTOCFT_CF_API_TOKEN=${AUTOCFT_CF_API_TOKEN}
      - AUTOCFT_CF_ACCOUNT_ID=${AUTOCFT_CF_ACCOUNT_ID}
      - AUTOCFT_CF_TUNNEL_ID=${AUTOCFT_CF_TUNNEL_ID}
      # Optional
      - AUTOCFT_CRON=*/30 * * * * *    # default */10 * * * * * (every 10s)
      - AUTOCFT_BASEDIR=/app/data
      #- AUTOCFT_ORIGIN_NO_TLS_VERIFY=true
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - autocft_data:/app/data
    restart: unless-stopped
volumes:
  autocft_data: {}
```

`.env` file:
```
AUTOCFT_CF_API_TOKEN=your_token_here
AUTOCFT_CF_ACCOUNT_ID=account_id_here
AUTOCFT_CF_TUNNEL_ID=tunnel_uuid_here
```

Start:
```
docker compose up -d autocft
```

## Label Target Business Containers
```yaml
  myapp:
    image: ghcr.io/example/myapp:latest
    labels:
      - autocft.enabled=true
      - autocft.hostname=app.example.com
      - autocft.service=http://myapp:8080
      # Optional
      - autocft.path=/
      - autocft.origin.no-tls-verify=true
```

After roughly one scheduling period (default 10s) visiting `https://app.example.com` should route through the Tunnel to the internal container.

## Manual Sync (Dry Run)
```
docker exec -it autocft ./autocft run --dry
```
Outputs the diff that would be applied; does not submit updates.

## History File
`${AUTOCFT_BASEDIR}/latest.json` stores only ingress entries managed by labels. Deleting it causes next run to be treated as a first run and re‑reconcile.

## Next Steps
- [Environment variable details](/en/environment)
- [Container labels](/en/labels)
- [Troubleshooting](/en/troubleshooting)
