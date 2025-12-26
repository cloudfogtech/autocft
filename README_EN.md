# Auto Cloudflare Tunnel (AutoCFT)

[![CI](https://github.com/cloudfogtech/autocft/actions/workflows/ci.yml/badge.svg)](https://github.com/cloudfogtech/autocft/actions) [![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](go.mod) [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE) [![Release](https://img.shields.io/github/v/release/cloudfogtech/autocft?color=green)](https://github.com/cloudfogtech/autocft/releases) ![Status](https://img.shields.io/badge/status-alpha-orange) [![Docker Pulls](https://img.shields.io/docker/pulls/cloudfogtech/autocft)](https://hub.docker.com/r/cloudfogtech/autocft) [![Commit Activity](https://img.shields.io/github/commit-activity/m/cloudfogtech/autocft)](https://github.com/cloudfogtech/autocft) [![Docs](https://img.shields.io/badge/Docs-Cloudflare%20Pages-blue)](https://<your-pages-subdomain>.pages.dev)

> Lightweight tool that auto-manages Cloudflare Tunnel ingress rules from Docker container labels.
>
> Add `autocft.*` labels to your application containers. A scheduler periodically parses, diffs, merges and updates the Tunnel configuration. Entries you add manually in the Cloudflare Zero Trust dashboard are detected as "web managed" and always preserved.

---

Please see more details in：https://autocft.cloudfogtech.ltd/en

## ✨ Features
- **Zero manual maintenance**: No more editing Tunnel ingress repeatedly in the dashboard.
- **Second-level sync**: Runs every 10s by default (configurable via `AUTOCFT_CRON`).
- **Low footprint**: Written in Go; single binary; no external DB (just one history JSON file).
- **Safe merge**: Only overwrites label-driven hostnames; preserves web‑added routes.
- **Dry run**: Manually preview the diff before applying changes.

## 🚀 Quick Start
### 1. Gather Cloudflare values
| Item | How to Obtain |
|------|---------------|
| CF_API_TOKEN | Dashboard -> My Profile -> API Tokens (custom token) |
| CF_ACCOUNT_ID | Dashboard URL `/accounts/<id>/` or API Tokens page |
| CF_TUNNEL_ID | `cloudflared tunnel list` or Tunnel detail page UUID |

### 2. docker-compose example
```yaml
version: '3.8'
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
    # Optional: expose PocketBase internal UI/API port for future Web UI
    ports:
      - 8090:8090
    restart: unless-stopped
volumes:
  autocft_data: {}
```
`.env` file:
```
AUTOCFT_CF_API_TOKEN=your_token
AUTOCFT_CF_ACCOUNT_ID=account_id
AUTOCFT_CF_TUNNEL_ID=tunnel_uuid
```
Start:
```
docker compose up -d autocft
```

### 3. Label application containers
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
After roughly one scheduling period (default 10s) visiting `https://app.example.com` should route through the Tunnel.

### 4. Manual dry run
```
docker exec -it autocft ./autocft run --dry
```
Prints the pending diff without performing the update.

## ⚙️ Environment Variables
System-level:
| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| AUTOCFT_CF_API_TOKEN | Yes | - | Cloudflare API Token (Tunnel read/write). |
| AUTOCFT_CF_ACCOUNT_ID | Yes | - | Cloudflare Account ID. |
| AUTOCFT_CF_TUNNEL_ID | Yes | - | Tunnel UUID. |
| AUTOCFT_BASEDIR | No | /app/data | Working directory (history / PocketBase data). |
| AUTOCFT_CRON | No | */10 * * * * * | Sync cron (with seconds) default 10s. |

Global origin fallbacks (used if a container omits the corresponding label):
| Variable | Type | Description |
|----------|------|-------------|
| AUTOCFT_ORIGIN_CONNECT_TIMEOUT | int(ms) | TCP connect timeout. |
| AUTOCFT_ORIGIN_DISABLE_CHUNKED_ENCODING | bool | Disable chunked encoding. |
| AUTOCFT_ORIGIN_HTTP2_ORIGIN | bool | Attempt HTTP/2. |
| AUTOCFT_ORIGIN_HTTP_HEADER | string | Override Host header. |
| AUTOCFT_ORIGIN_KEEP_ALIVE_CONNECTIONS | int | Max idle keepalive connections. |
| AUTOCFT_ORIGIN_KEEP_ALIVE_TIME | int(s) | Keepalive idle timeout. |
| AUTOCFT_ORIGIN_NO_HAPPY_EYEBALLS | bool | Disable Happy Eyeballs. |
| AUTOCFT_ORIGIN_NO_TLS_VERIFY | bool | Skip TLS verification. |
| AUTOCFT_ORIGIN_ORIGIN_SERVER_NAME | string | Expected certificate hostname. |
| AUTOCFT_ORIGIN_PROXY_TYPE | string | "" or socks. |
| AUTOCFT_ORIGIN_TCP_KEEP_ALIVE | int(s) | TCP keepalive interval. |
| AUTOCFT_ORIGIN_TLS_TIMEOUT | int(s) | TLS handshake timeout. |

## 🏷️ Container Labels
Required:
| Label | Description |
|-------|-------------|
| autocft.enabled | Set to `true` to manage this container. |
| autocft.hostname | Public hostname to expose. |
| autocft.service | Upstream service address (http://, https://, tcp://, etc.). |

Optional:
| Label | Description |
|-------|-------------|
| autocft.path | URL path prefix (default `/`). |
| autocft.origin.connect-timeout | Connect timeout ms. |
| autocft.origin.disable-chunked-encoding | Disable chunked transfer. |
| autocft.origin.http2-origin | Force HTTP/2. |
| autocft.origin.http-host-header | Override Host header. |
| autocft.origin.keep-alive-connections | Max idle connections. |
| autocft.origin.keep-alive-timeout | Idle timeout s. |
| autocft.origin.no-happy-eyeballs | Disable dual-stack race. |
| autocft.origin.no-tls-verify | Skip TLS verification. |
| autocft.origin.origin-server-name | Certificate hostname. |
| autocft.origin.proxy-type | Blank or socks. |
| autocft.origin.tcp-keep-alive | TCP keepalive s. |
| autocft.origin.tls-timeout | TLS handshake timeout s. |

Removal: remove the label or set `autocft.enabled=false`; next sync deletes the rule (if not web‑managed).

## 🔄 Merge Logic Summary
1. Fetch Cloudflare ingress list (skip required 404 fallback entry).
2. Read history file `latest.json` (if absent treat all existing ingress entries as web‑managed).
3. Parse container labels for new desired ingress entries.
4. Merge: web‑managed + label entries (container overrides same hostname).
5. If new set (excluding fallback) equals history -> skip update.
6. Write label-managed set to `latest.json`.

## 🧪 Build & Run
Build binary locally:
```
go build -o autocft ./cmd/autocft
./autocft
```
Build image:
```
docker build -t cloudfogtech/autocft:dev .
```
Run directly:
```
AUTOCFT_CF_API_TOKEN=xxx \
AUTOCFT_CF_ACCOUNT_ID=xxx \
AUTOCFT_CF_TUNNEL_ID=xxx \
./autocft
```

## 🛠️ Troubleshooting Quick Table
| Issue | Hint |
|-------|------|
| No ingress generated | Look for `Container skipped`; verify required labels. |
| 404 response | Check `autocft.path`; verify upstream port/service. |
| TLS errors | Use label or env to set `no-tls-verify`. |
| Changes not applied | Delete `latest.json`; ensure API token permissions valid. |
| Cron not firing | Must have 6 segments (with seconds). |
| Multiple containers same host | Only one kept; avoid duplicates. |

See detailed docs under `docs/en/*` (VitePress site) for more.

## 📊 PocketBase & Future Web UI
A PocketBase instance is started internally (see `cmd/autocft/main.go`). Current purposes:
- Foundation for an upcoming Web UI (visual ingress list, manual sync trigger, diff viewer).
- Persistence for future features (events, audit logs, alert rules).
- Unified auth & access control (planned RBAC / API tokens).

Default data dir: `${AUTOCFT_BASEDIR}/pb_data`.

If you want to inspect it early, map the port:
```
ports:
  - 8090:8090
```
Visit: `http://<host-ip>:8090`.

> Note: No official UI yet; exposing the port is optional and mainly for experimentation.

### Planned UI Features (Draft)
- Ingress list visualization (Web managed vs Label managed).
- Diff history & rollback.
- Manual sync (with Dry Run preview panel).
- Change audit & operation logs export.
- Basic alerting (sync failure, abnormal change volume).
- Templated hostname / variable substitution editor.

## 🌐 Automated Docs Deployment (Cloudflare Pages)
The `docs/` directory is a VitePress site. A GitHub Actions workflow builds and deploys docs to Cloudflare Pages whenever documentation changes land on `main`.

### Workflow Overview
- Trigger: push to `main` affecting `docs/**`.
- Build: `pnpm install` + `pnpm run docs:build`.
- Deploy: `wrangler pages deploy docs/.vitepress/dist`.
- Concurrency: latest run for a branch cancels previous in-flight run.

### Required Secrets
| Secret | Description | Where to Get |
|--------|-------------|--------------|
| CLOUDFLARE_API_TOKEN | API token with Cloudflare Pages: Edit permission | Dashboard -> API Tokens (custom token) |
| CLOUDFLARE_ACCOUNT_ID | Your Cloudflare Account ID | Dashboard URL `/accounts/<id>/` |
| CLOUDFLARE_PAGES_PROJECT_NAME | Pages project name (slug) | Pages project settings |

### Paths & Output
- Source: `docs/`
- Output: `docs/.vitepress/dist`
- Home pages: `docs/en/index.md` / `docs/zh/index.md`

### Workflow Snippet
```yaml
on:
  push:
    branches: [ main ]
    paths:
      - 'docs/**'
jobs:
  deploy-docs:
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'pnpm'
      - run: corepack enable
      - run: pnpm install --frozen-lockfile
        working-directory: docs
      - run: pnpm run docs:build
        working-directory: docs
      - run: npx wrangler pages deploy docs/.vitepress/dist --project-name "$PAGES_PROJECT_NAME" --branch "${GITHUB_REF_NAME}" --commit-dirty=true
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
          CLOUDFLARE_ACCOUNT_ID: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
```
> If you only need production deploys, fix the branch to `main` or disable branch previews in Pages.

### FAQ
| Issue | Resolution |
|-------|------------|
| 403 / permission error | Ensure token has Cloudflare Pages: Edit |
| 404 after deploy | Confirm output path `docs/.vitepress/dist` and project slug correctness |
| Slow builds | Use Node 20 + pnpm caching; avoid unnecessary deps |
| Need preview envs | Keep dynamic `--branch ${GITHUB_REF_NAME}` for branch‑specific previews |
| Mermaid diagrams missing | Verify `vitepress-plugin-diagrams` compatibility with installed VitePress version |

### Potential Enhancements
- Versioned docs based on tags.
- Cache purge for changed assets.
- Lighthouse CI for performance metrics.

## 🗺️ Roadmap
- Web UI management view (PocketBase powered)
- Templated hostnames / variable substitution
- Notifications / alerts (Webhook, Email)
- Fine-grained permissions (RBAC / tokens)
- Diff history & rollback
- Audit log & export

## 🤝 Contributing
1. Fork & create feature branch.
2. Keep code style consistent.
3. Add comments / doc updates where helpful.
4. Build & self‑test before PR.

## 📄 License
MIT License © 2025 CloudFogTech

## 💬 Contact / Support
- GitHub: https://github.com/cloudfogtech/autocft
- When filing issues include labels, concise logs, version (redact sensitive data).

---
**Disclaimer**: Early 0.x releases may introduce breaking changes; evaluate before production use.
