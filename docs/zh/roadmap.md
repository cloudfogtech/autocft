---
outline: deep
---
# 发展路线图

## 🚧 Roadmap
- Web UI 管理视图（PocketBase 支撑）
  - 当前 Ingress 列表可视化（区分 Web 管理 vs 标签托管）。
  - 身份认证与权限管理
- 模板化 Hostname / 变量替换
- 简易通知/告警功能（如同步失败、配置变更数量异常）（Webhook, Email）
- Diff 历史浏览与回滚。
- 手动同步按钮（支持 Dry Run 预览）。
- 变更审计与操作日志。
- 文档
  - 引入 Lighthouse CI 做文档性能检测。

## 📊 PocketBase 与未来 Web UI
AutoCFT 内置了一个 PocketBase 实例（在入口 `cmd/autocft/main.go` 中创建），当前主要作用：
- 预留未来的 Web UI / 管理界面能力（例如查看当前 Ingress、手动触发同步、查看日志/差异）。
- 存储潜在的扩展数据（事件、操作审计、告警规则）。
- 加入未来的身份认证功能。

默认数据目录：`${AUTOCFT_BASEDIR}/pb_data`（已通过 Docker 卷挂载示例）。

如需在现在就访问 PocketBase（初始 UI 为其默认后台），请在 `docker-compose` 中加入：
```
ports:
  - 8090:8090
```
访问：`http://<你的宿主IP>:8090`。

> 注意：当前版本未正式提供图形化管理页面；端口暴露仅用于验证与未来扩展，生产环境可暂不映射。

---

> PocketBase 默认管理员账号：
> - 邮箱：admin@example.com
> - 密码：autocft@admin#123
