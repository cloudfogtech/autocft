---
outline: deep
---

# FAQ

## Why not just edit in the Web dashboard?
You can. AutoCFT only manages the ingress records that previously appeared as "label‑driven" in history; all others are preserved as web‑managed.

## Sync frequency?
Default 10 seconds; adjustable via `AUTOCFT_CRON`.

## Does it create DNS records automatically?
No. It only updates Tunnel ingress. Required CNAME must be created manually on first hostname addition (or automatically by Cloudflare if applicable).

## Multiple instances against one Tunnel?
Not recommended; would cause competing writes. Keep a single instance.

## Path routing supported?
Yes. Use `autocft.path`. Multiple entries are sorted by hostname + path.

## How to delete a rule?
Remove `autocft.enabled` or set it to false.

## Templated hostnames supported?
Not yet; must be explicit. Planned for future (templates / variables).
