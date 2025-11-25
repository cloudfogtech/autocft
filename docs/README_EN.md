# AutoCFT Development & Deployment Guide (English)

This document provides development and operational guidance for contributors working on AutoCFT. It mirrors the structure of the Chinese developer README.

## Directory Layout
```
assets/            # Static assets (images, diagrams)
en/                # English documentation
zh/                # Chinese documentation
index.md           # Documentation entry
package.json       # Node.js project config (docs tooling)
pnpm-lock.yaml     # pnpm dependency lockfile
README.md          # Chinese dev guide
README_EN.md       # English dev guide (this file)
```

## Prerequisites
1. Clone repository:
```bash
git clone <REPO_URL>
cd autocft/docs
```
2. Install dependencies (prefer pnpm):
```bash
pnpm install
```
If pnpm is not installed:
```bash
npm install -g pnpm
```

## Local Docs Development
Project docs can be served with VitePress (or similar). Adjust commands if package scripts differ.
```bash
pnpm docs:dev    # Start local dev server
pnpm docs:build  # Build static site
```

## Documentation Maintenance
- Keep parallel structure between `en/` and `zh/` for easier synchronization.
- Place shared assets under `assets/`.
- Maintain `index.md` as the landing/entry page.
- When adding a new Chinese guide, immediately create its English counterpart (even a placeholder) to track parity.

## Adding / Updating Content
1. Create or edit Markdown under the appropriate locale directory.
2. Reference images using relative paths from `/assets/...`.
3. Run `pnpm docs:dev` to preview changes before committing.
4. Keep headings consistent (English <-> Chinese) so automated TOC / link checks remain stable.

## Dependency Management
- Use `pnpm` to ensure deterministic installs.
- Commit updates to both `package.json` and `pnpm-lock.yaml` when adding or bumping tooling.
- Avoid adding heavy tooling unless necessary (prefer lightweight markdown linting).

## Suggested Conventions
- File names: lowercase, hyphen-separated (e.g., `cloudflare.md`).
- Front matter: include `outline: deep` where long pages benefit from a sidebar outline.
- Use fenced code blocks with a language tag (`bash`, `powershell`, `yaml`).
- Prefer absolute doc links like `/en/deploy` for internal navigation.

## Sync Checklist (CN -> EN)
Use this list whenever Chinese docs are updated:
- Created/updated English equivalent file.
- Verified tables translated correctly (column alignment, units, default values).
- Confirmed any images exist and are language-neutral; if language-specific, add an English variant.
- Added cross-links (e.g., environment variables -> labels page).
- Ran build to ensure no broken links.

## Contribution Workflow
1. Fork / branch: `feat/<topic>` or `docs/<topic>`.
2. Make changes & preview locally.
3. Run optional lint (if configured): `pnpm run lint`.
4. Commit with conventional message (`docs: add english roadmap`).
5. Open Pull Request; describe scope & parity status.

## Common Scripts (if defined in package.json)
```bash
pnpm docs:dev     # Live preview
pnpm docs:build   # Static build
pnpm docs:preview # Preview built site (if configured)
```

## Security & Secrets
- Never commit real tokens; use placeholders in examples (`xxxxxxxx`).
- If adding CI examples, rely on platform secrets (e.g., GitHub Actions secrets).

## Roadmap & Future UI
See `/en/roadmap` for planned Web UI & PocketBase integration items.

## Cloudflare Credentials Guide
See `/en/cloudflare` for obtaining `CF_API_TOKEN`, `CF_ACCOUNT_ID`, `CF_TUNNEL_ID`.

## Environment Variables
See `/en/configuration` for runtime variables and overrides.

## Troubleshooting
Consult `/en/troubleshooting` and compare with Chinese `/zh/troubleshooting` for any divergence.

## License
MIT License © 2025 CloudfogTech
---
outline: deep
---
# Cloudflare Credentials Acquisition Guide

This document explains how to obtain and configure the following three core parameters in the Cloudflare Dashboard:

- `CF_API_TOKEN` (API Token, recommend least privilege)
- `CF_ACCOUNT_ID` (Account ID)
- `CF_TUNNEL_ID` (Cloudflare Tunnel ID)

---
## Table of Contents
1. [Concept Overview](#1-concept-overview)
2. [Preparation](#2-preparation)
3. [Obtain CF_ACCOUNT_ID](#3-obtain-cf_account_id)
4. [Create Least-Privilege CF_API_TOKEN](#4-create-least-privilege-cf_api_token)
5. [Get or Create Cloudflare Tunnel and CF_TUNNEL_ID](#5-get-or-create-cloudflare-tunnel-and-cf_tunnel_id)
6. [Verify Configuration](#6-verify-configuration)
7. [Common Issues & Troubleshooting](#7-common-issues--troubleshooting)
8. [Security Notes](#8-security-notes)
9. [Quick Checklist](#9-quick-checklist)

---
## 1. Concept Overview
- `CF_ACCOUNT_ID`: Unique identifier of your Cloudflare account; required for some Account-level API calls.
- `CF_API_TOKEN`: Modern scoped API token. Prefer this over the deprecated Global API Key.
- `CF_TUNNEL_ID`: UUID of a Cloudflare Tunnel instance (formerly Argo Tunnel) used for API-based management.

> Avoid using the "Global API Key" in automation. If temporarily used, rotate it ASAP and migrate to a scoped token.

---
## 2. Preparation
Ensure you have:
- Registered and logged into Cloudflare.
- At least one domain (Zone) active in Cloudflare (not strictly required, but helpful for DNS linkage).
- Decided required permissions (Tunnel read-only vs. read/write vs. DNS Edit needs).

> This guide uses only the REST API. No need to install the `cloudflared` CLI locally.

---
## 3. Obtain CF_ACCOUNT_ID

![Account ID path](/assets/images/cloudflare-01.png)

### Steps
1. Log in: https://dash.cloudflare.com/
2. Enter any Zone (domain) detail page.
3. Inspect browser URL: `https://dash.cloudflare.com/<ACCOUNT_ID>/home/domains`.
4. Copy `<ACCOUNT_ID>` as `CF_ACCOUNT_ID`.

### Alternative: List accounts via API
Requires a token with Account read permissions:
```
GET https://api.cloudflare.com/client/v4/accounts
Authorization: Bearer <CF_API_TOKEN>
```
Response JSON `result[0].id` is the Account ID (choose the needed one if multiple).

---
## 4. Create Least-Privilege CF_API_TOKEN

### Target Permissions (adjust as needed):
- Account -> Cloudflare Tunnel: Edit (allows reading & updating Tunnel ingress configuration)

> If you only need to read Tunnel status, you can restrict to Read instead of Edit.

### Steps
1. Top-right avatar dropdown -> "My Profile".
2. Left navigation: "API Tokens".
3. Click "Create Token".
4. Choose "Create Custom Token".
5. Fill `Token Name`, e.g., "AutoCFT Token".
6. Add permission:
   - `Account` / `Cloudflare Tunnel` / `Edit` (or Read if sufficient).
7. In `Account Resources`, select the specific account (or all if only one).
8. Optional: restrict Client IP Addresses, set TTL if required.
9. Continue to summary and create.
10. Copy the generated token (only shown once) as `CF_API_TOKEN`.

---
## 5. Get or Create Cloudflare Tunnel and CF_TUNNEL_ID
You can use Dashboard or API. Either is fine.

### A. Dashboard
1. In left navigation search or expand "Zero Trust".
2. Enter "Tunnels" (Path: Zero Trust > Network > Tunnels).
3. Click the required Tunnel row.
4. The UUID shown in URL or details is `CF_TUNNEL_ID`.

### B. API List Existing Tunnels
```
curl -X GET "https://api.cloudflare.com/client/v4/accounts/<CF_ACCOUNT_ID>/cfd_tunnel" \
  -H "Authorization: Bearer <CF_API_TOKEN>" \
  -H "Content-Type: application/json"
```
Example partial response:
```
{
  "success": true,
  "result": [
    {
      "id": "3a1b2c3d-4e5f-6789-abcd-0123456789ef",
      "name": "my-tunnel",
      "created_at": "2025-10-01T02:33:11Z",
      "connections": [...]
    }
  ]
}
```
Each object's `id` is a `CF_TUNNEL_ID`. Empty array means no tunnels yet.

---
## 6. Verify Configuration

### Environment Variable Example (Windows CMD)
```
set CF_API_TOKEN=xxxxxxxxxxxxxxxx
set CF_ACCOUNT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
set CF_TUNNEL_ID=3a1b2c3d-4e5f-6789-abcd-0123456789ef
```

### Test Token & Tunnel via API
```
curl -X GET "https://api.cloudflare.com/client/v4/accounts/%CF_ACCOUNT_ID%/cfd_tunnel/%CF_TUNNEL_ID%" ^
  -H "Authorization: Bearer %CF_API_TOKEN%" ^
  -H "Content-Type: application/json"
```

### List All (reconfirm permissions)
```
curl -X GET "https://api.cloudflare.com/client/v4/accounts/%CF_ACCOUNT_ID%/cfd_tunnel" ^
  -H "Authorization: Bearer %CF_API_TOKEN%" ^
  -H "Content-Type: application/json"
```

### Common HTTP Codes
- 200: Success
- 403: Insufficient permissions (missing Cloudflare Tunnel scope)
- 401: Invalid or revoked token

---
## 7. Common Issues & Troubleshooting
| Issue | Possible Cause | Recommendation |
|-------|----------------|----------------|
| 403 Forbidden | Token lacks Tunnel scope | Recreate token with Account -> Cloudflare Tunnel -> Read/Edit |
| Empty tunnel list | No tunnel created yet | Use POST to create a new tunnel |
| Cannot find expected tunnel | Wrong Account ID | Confirm ACCOUNT_ID in dashboard URL matches token scope |
| DNS record not auto-updated | Missing Zone DNS Edit permission | Add Zone / DNS / Edit or configure record manually |
| Wrong UUID copied | Extra characters in paste | Copy only 36-char UUID (with 4 dashes) |
| Global API Key used | Confused token types | Migrate to scoped API Token, revoke unused global key |

---
## 8. Security Notes
- Principle of least privilege: scope and resource restriction.
- Never commit tokens to public repos; use env vars or secret managers.
- Rotate tokens periodically; revoke unused ones.
- For CI/CD: store secrets using platform secret storage (e.g., GitHub Actions Secrets).
- Avoid printing full tokens in logs; mask except leading/trailing chars.

---
## 9. Quick Checklist
- Correct `CF_ACCOUNT_ID` obtained (URL / API list)
- Least-privilege `CF_API_TOKEN` created & stored
- `CF_TUNNEL_ID` listed or created via dashboard/API
- Environment variables exported for AutoCFT
- Standalone GET verified for tunnel access

> After completing the above, you can safely use Cloudflare features in AutoCFT.

