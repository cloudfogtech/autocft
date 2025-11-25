---
outline: deep
---

# 故障排查

## 没有生成 Ingress
- 查看日志是否有 `Container skipped`。
- 确认打了必要标签。
- 手动触发：`docker exec -it autocft ./autocft run`。

## 访问 404
- `autocft.path` 是否与请求路径匹配。
- `autocft.service` 端口是否正确、服务是否启动。

## TLS 证书问题
自签名证书可加：`autocft.origin.no-tls-verify=true` 或全局环境变量。

## 修改不生效
- 删除 `latest.json` 强制全量。
- 检查 API Token 权限是否仍有效。

## Cron 不执行
确认表达式 6 段。例如每分钟一次：`0 * * * * *`。

## 多容器同 Host
只会保留一个，避免重复。

## 仍有问题
提交 issue，附：相关标签、精简日志、版本号（脱敏）。

