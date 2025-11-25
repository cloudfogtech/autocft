---
outline: deep
---

# 常见问题

## 为什么不直接在 Web 控制台改？
可以。`AutoCFT` 只管理历史中出现过的“标签驱动”记录，其他保留为 `web-managed`。

## 同步频率？
默认 10 秒，可用 `AUTOCFT_CRON` 调整。

## 会自动创建 DNS 记录吗？
会，更新 `Tunnel ingress` 后由 `Cloudflare` 自动创建。

## 可多实例同时运行同一 Tunnel 吗？
不推荐，会导致竞争写，造成不可预知的错误，所以请保持单实例运行。

## 路径路由支持？
支持，通过 `autocft.path` 设置，多条按 `hostname + path` 排序。

## 如何删除规则？
移除 `autocft.enabled` 或将其设为 `false` 即可。

## 模板化 Hostname 支持？
暂未支持，需要显式写出。后续计划加入模板 / 变量功能。
