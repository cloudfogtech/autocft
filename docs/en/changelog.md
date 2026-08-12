# Changelog

## 0.2.3
- Support the `matchSNItoHost` origin parameter (preserve SNI settings configured on wildcard rules in the Cloudflare dashboard)
- Upgrade cloudflare-go to v6.10.0

## 0.2.2
- Sort wildcard hostnames after exact hostnames in ingress rules

## 0.2.1
- Update Go version to 1.26 in Docker build arguments

## 0.2.0
- Adapt to the new Cloudflare API
- Upgrade to Go version 1.26
- Upgrade dependencies

## 0.1.5
- Remove default user in Dockerfile

## 0.1.4
- Fix cron job logic
- Upgrade dependencies

## 0.1.3
- Fix startup issue in Dockerfile
- Add default mount directory permissions

## 0.1.2
- Change docker mounted permission

## 0.1.1
- Change default base directory to `/app/data` to avoid docker mount error

## 0.1.0
- Change default admin account username and password
- Release first preview version

## 0.0.3
- Core function is usable

## 0.0.2
- Improve CI/CD workflow

## 0.0.1
- Initial prototype
- Docker label discovery
- Cloudflare Tunnel ingress sync (merge + diff)
- Fallback 404 ingress auto-append
- History tracking file
- Basic cron scheduling
