---
outline: deep
---
# Roadmap

## 🚧 Planned Features
- Web UI management view (PocketBase backed)
  - Visualize current Ingress list (distinguish web-managed vs label-managed)
  - Authentication & RBAC
- Templated hostname / variable substitution
- Lightweight notification & alerting (sync failure, abnormal change counts) (Webhook, Email)
- Diff history browsing & rollback
- Manual sync button (supports Dry Run preview)
- Change audit & operation logs
- Documentation
  - Introduce Lighthouse CI for documentation performance checks

## 📊 PocketBase & Future Web UI
AutoCFT embeds a PocketBase instance (created in entry `cmd/autocft/main.go`). Current purposes:
- Reserve future Web UI / admin capabilities (view ingress, trigger manual sync, view logs/diffs).
- Store potential extension data (events, audit trail, alerting rules).
- Foundation for future authentication.

Default data directory: `${AUTOCFT_BASEDIR}/pb_data` (already included in Docker volume examples).

To access PocketBase now (initial UI is the default admin backend), add to `docker-compose`:
```
ports:
  - 8090:8090
```
Visit: `http://<host-ip>:8090`

> Note: Current version does not ship a formal graphical admin page; exposing the port is optional and mainly for validation or future extension. In production you may omit it.

---

> PocketBase default admin credentials:
> - Email：autocft@cloudfogtech.ltd
> - Password：autocft@cloudfogtech#123
