---
outline: deep
---

# 环境变量配置

所有运行期配置通过环境变量传入，系统类变量统一前缀 `AUTOCFT_`。

## 系统变量
| 变量                     | 必填 | 默认 | 说明                               |
|------------------------|------|------|----------------------------------|
| AUTOCFT_CF_API_TOKEN   | 是 | - | Cloudflare API Token（Tunnel 读写）。 |
| AUTOCFT_CF_ACCOUNT_ID  | 是 | - | Cloudflare Account ID。           |
| AUTOCFT_CF_TUNNEL_ID   | 是 | - | 目标 Cloudflare Tunnel UUID。       |
| AUTOCFT_BASEDIR        | 否 | /app/autocft | 工作目录（历史文件 / PocketBase 数据）。      |
| AUTOCFT_CRON           | 否 | */10 * * * * * | 同步Cron（含秒），默认10秒触发一次。            |
| AUTOCFT_ADMIN_EMAIL    | 否 | admin@example.com | 设置 admin账户的邮箱（为PocketBase预留） |
| AUTOCFT_ADMIN_PASSWORD | 否 | autocft@admin#123 | 设置 admin账户的密码（为PocketBase预留） |

## 默认 Origin 覆盖项
作为全局回退，当容器未设置相应标签时使用，这部分参数可以不设置，默认值由Cloudflare提供。

| 变量 | 类型 | 说明 |
|------|------|------|
| AUTOCFT_ORIGIN_CONNECT_TIMEOUT | int(ms) | TCP 连接超时。 |
| AUTOCFT_ORIGIN_DISABLE_CHUNKED_ENCODING | bool | 禁用分块传输。 |
| AUTOCFT_ORIGIN_HTTP2_ORIGIN | bool | 尝试 HTTP/2。 |
| AUTOCFT_ORIGIN_HTTP_HEADER | string | 覆盖 Host 头。 |
| AUTOCFT_ORIGIN_KEEP_ALIVE_CONNECTIONS | int | 最大空闲连接数。 |
| AUTOCFT_ORIGIN_KEEP_ALIVE_TIME | int(s) | 空闲连接超时。 |
| AUTOCFT_ORIGIN_NO_HAPPY_EYEBALLS | bool | 禁用 Happy Eyeballs。 |
| AUTOCFT_ORIGIN_NO_TLS_VERIFY | bool | 跳过 TLS 校验。 |
| AUTOCFT_ORIGIN_ORIGIN_SERVER_NAME | string | 预期证书主机名。 |
| AUTOCFT_ORIGIN_PROXY_TYPE | string | "" 或 socks。 |
| AUTOCFT_ORIGIN_TCP_KEEP_ALIVE | int(s) | TCP keepalive 间隔。 |
| AUTOCFT_ORIGIN_TLS_TIMEOUT | int(s) | TLS 握手超时。 |

## Cron 表达式
使用 `robfig/cron`（含秒）格式：`秒 分 时 日 月 周`。

## 安全建议
- 限制 API Token 权限最小化。
- docker.sock 挂载建议只读。
- 监控日志异常改动。
- 定期轮换 API Token。
