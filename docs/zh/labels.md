---
outline: deep
---

# 容器标签参考
以 `autocft.` 开头的 Docker 标签控制建立的 Ingress 规则。

## 必需标签
| 标签 | 说明                                                                        |
|------|---------------------------------------------------------------------------|
| autocft.enabled | 设为 `true` 才会被管理。                                                          |
| autocft.hostname | 暴露的公网 Hostname。                                                           |
| autocft.service | 上游服务地址，例如 `http://app:8080`、`tcp://redis:6379` 等，指从cloudflared访问到目标服务的地址。 |

## 可选标签
| 标签 | 类型 | 说明                | 回退 |
|------|------|-------------------|------|
| autocft.path | string | URL 路径前缀（默认 `/`）。 | `/` |
| autocft.origin.connect-timeout | int | 连接超时，单位：ms。       | ENV |
| autocft.origin.disable-chunked-encoding | bool | 禁用分块。             | ENV |
| autocft.origin.http2-origin | bool | 强制 HTTP/2。        | ENV |
| autocft.origin.http-host-header | string | 覆盖 Host 头。        | ENV |
| autocft.origin.keep-alive-connections | int | 最大空闲连接。           | ENV |
| autocft.origin.keep-alive-timeout | int | 空闲超时，单位：s。        | ENV |
| autocft.origin.no-happy-eyeballs | bool | 禁用双栈竞速。           | ENV |
| autocft.origin.no-tls-verify | bool | 跳过 TLS 验证。        | ENV |
| autocft.origin.origin-server-name | string | 证书主机名。            | ENV |
| autocft.origin.proxy-type | string | 空（http）或 socks。   | ENV |
| autocft.origin.tcp-keep-alive | int | TCP keepalive，单位：s。  | ENV |
| autocft.origin.tls-timeout | int | TLS 超时，单位：s。         | ENV |

## 示例
```yaml
services:
  api:
    image: ghcr.io/example/api:latest
    labels:
      - autocft.enabled=true
      - autocft.hostname=api.example.com
      - autocft.service=https://api:443
  web:
    image: ghcr.io/example/web:latest
    labels:
      - autocft.enabled=true
      - autocft.hostname=app.example.com
      - autocft.service=http://web:3000
      - autocft.origin.no-tls-verify=true
```

## 冲突规则
同一 hostname 只保留一个条目；请避免重复。若与“网页手动添加”的条目冲突，手动条目会在首次同步被识别为 web‑managed 并被保留。

## 删除
移除或将 `autocft.enabled=false`，下次同步即删除该条（若不是 web‑managed）。

