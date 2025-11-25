---
outline: deep
---

# 部署指南

Auto Cloudflare Tunnel（AutoCFT）用于自动将 Docker 容器上的标签同步为 Cloudflare Tunnel 的 Ingress 配置。你只需要在容器上添加 `autocft.*` 标签，AutoCFT 会按计划任务自动对比并更新 Tunnel 配置。

## 工作流程概述
1. 从环境变量（前缀 `AUTOCFT_`）读取系统默认配置。
2. 枚举当前运行中的容器，解析所有以 `autocft.` 开头的标签。
3. 校验必须标签：`autocft.enabled`、`autocft.hostname`、`autocft.service`。
4. 调用 Cloudflare API 读取现有 Tunnel Ingress，保留“非本工具管理”（通过网页配置）的条目。
5. 合并 + 生成更新列表，如有变化则调用 API 更新。
6. 将本次由标签管理的 Ingress 写入历史文件 `latest.json`，用于下次 diff。

## 前置条件
- 已激活并托管在 Cloudflare 上的域名。
- 已创建的 Cloudflare Tunnel（可用 `cloudflared tunnel create <name>` 或 Zero Trust 控制台）。需要 Tunnel ID。
- 一个 Cloudflare API Token，权限至少包含：
  - 账号级：Cloudflare Tunnel Read & Edit
  - （多 Zone Hostname 场景）Zone 级：DNS Read（可选，推荐）
- AutoCFT 容器能够访问宿主的 `/var/run/docker.sock`（只读即可）。

## 需要收集的信息
| 项目 | 获取方式 |
|------|----------|
| CF_API_TOKEN | Dashboard -> My Profile -> API Tokens 创建自定义 Token |
| CF_ACCOUNT_ID | Dashboard URL `/accounts/<id>/` 或 API Tokens 页面 |
| CF_TUNNEL_ID | `cloudflared tunnel list` 或 Tunnel 详情页 UUID |

详情见 [Cloudflare 配置指南](/zh/cloudflare)。

## Image
镜像源一共有两种选择：
- [GitHub Container Registry](https://github.com/cloudfogtech/autocft/pkgs/container/autocft): `ghcr.io/cloudfogtech/autocft:latest`
- [Docker Hub](https://hub.docker.com/r/cloudfogtech/autocft): `github.com/cloudfogtech/autocft:latest`

## docker-compose 示例
```yaml
services:
  autocft:
    image: github.com/cloudfogtech/autocft:latest
    # image: ghcr.io/cloudfogtech/autocft:latest
    container_name: autocft
    environment:
      - AUTOCFT_CF_API_TOKEN=${AUTOCFT_CF_API_TOKEN}
      - AUTOCFT_CF_ACCOUNT_ID=${AUTOCFT_CF_ACCOUNT_ID}
      - AUTOCFT_CF_TUNNEL_ID=${AUTOCFT_CF_TUNNEL_ID}
      # 可选
      - AUTOCFT_CRON=*/30 * * * * *    # 默认 */10 * * * * * 每 10 秒
      - AUTOCFT_BASEDIR=/app/autocft
      #- AUTOCFT_ORIGIN_NO_TLS_VERIFY=true
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - autocft_data:/app/autocft
    restart: unless-stopped
volumes:
  autocft_data: {}
```

启动：
```
docker compose up -d autocft
```

环境变量详见：[环境变量](/zh/environment)。

## 目标业务容器打标签

标签详见：[容器标签](/zh/labels)。

```yaml
  nginx:
    image: github.com/nginx:latest
    container_name: nginx
    labels:
      - autocft.enabled=true
      - autocft.hostname=web.example.com
      - autocft.service=http://nginx:80
      # 可选
      - autocft.path=/
      # 如果 autocft.service 是 http, 则可以跳过 TLS 验证
      - autocft.origin.no-tls-verify=true
```

大约一个调度周期后（默认 10 秒）访问 `https://web.example.com` 即可通过 Tunnel 访问内部容器。

## 手动测试（Dry Run）
```
docker exec -it autocft ./autocft run --dry
```
只输出将要更新的差异，不真正提交。

## 历史文件
`${AUTOCFT_BASEDIR}/latest.json` 仅保存“由容器标签管理”的 ingress 项。删除该文件会在下次执行时视作首次运行并重新对比。
