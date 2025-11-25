---
outline: deep
---

# Troubleshooting

## No Ingress Generated
- Look for log entries like `Container skipped`.
- Confirm required labels are set.
- Manual trigger: `docker exec -it autocft ./autocft run`.

## 404 When Accessing Hostname
- Check whether `autocft.path` matches the request path.
- Verify `autocft.service` port correctness and that upstream service is started.

## TLS Certificate Issues
Add label `autocft.origin.no-tls-verify=true` or set global env variable if using self‑signed certificates.

## Changes Not Taking Effect
- Delete `latest.json` to force full diff.
- Check if API Token permissions are still valid.

## Cron Not Executing
Ensure expression has 6 fields. Example every minute: `0 * * * * *`.

## Multiple Containers Same Host
Only one will be kept; avoid duplication.

## Still Having Problems
Open an issue; include relevant labels, concise log excerpt, version (redact sensitive data).
