# Auto Cloudflare Tunnel (AutoCFT)

[![CI](https://github.com/cloudfogtech/autocft/actions/workflows/ci.yml/badge.svg)](https://github.com/cloudfogtech/autocft/actions) [![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](go.mod) [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE) [![Release](https://img.shields.io/github/v/release/cloudfogtech/autocft?color=green)](https://github.com/cloudfogtech/autocft/releases) ![Status](https://img.shields.io/badge/status-alpha-orange) [![Docker Pulls](https://img.shields.io/docker/pulls/cloudfogtech/autocft)](https://hub.docker.com/r/cloudfogtech/autocft) [![Docs](https://img.shields.io/badge/Docs-Cloudflare%20Pages-blue)](https://<your-pages-subdomain>.pages.dev)

> 在 Docker 环境中基于容器标签自动同步并管理 Cloudflare Tunnel Ingress 规则的轻量工具。
>
> 通过给业务容器添加 `autocft.*` 标签，定时任务会自动解析、对比、合并并更新 Tunnel 配置；人工在 Cloudflare Zero Trust / Dashboard 上新增的条目会被识别为“Web 管理”并始终保留。

---

详细文档请访问：https://autocft.cloudfogtech.ltd/zh

## ✨ 特性
- **自动化**：无需手动在Cloudflare Zero Trust中配置隧道，节省时间和精力
- **秒级同步**：默认每 10 秒执行一次（可通过 `AUTOCFT_CRON` 配置）。
- **低资源占用**：Golang 编写，单文件运行；无额外数据库，仅一个历史 JSON 文件。
- **安全合并**：仅覆盖由标签驱动的主机名，Web 新增的规则自动保留。
- **Dry Run 支持**：手动触发并查看将要应用的差异，避免误操作。

## 📦 快速开始
### 1. 准备 Cloudflare 信息
收集以下 3 个必需值：
| 项目 | 获取方式 |
|------|----------|
| AUTOCFT_CF_API_TOKEN | Dashboard -> My Profile -> API Tokens 自定义创建 |
| CF_ACCOUNT_ID | Dashboard URL `/accounts/<id>/` 或 API Tokens 页 |
| CF_TUNNEL_ID | `cloudflared tunnel list` 或 Zero Trust Tunnel 详情页 UUID |

详情见 [Cloudflare 配置指南](https://autocft.cloudfogtech.ltd/zh/cloudflare)。

### 2. docker-compose 示例

**重要**
`./autocft_data` 文件夹需要将权限设置为 `10001:10001`
```bash
mkdir ./autocft_data
chown -R 10001:10001 ./autocft_data
```

```yaml
services:
  autocft:
    image: cloudfogtech/autocft:latest
    # image: ghcr.io/cloudfogtech/autocft:latest
    container_name: autocft
    users: 10001:10001
    environment:
      - AUTOCFT_CF_API_TOKEN=${AUTOCFT_CF_API_TOKEN}
      - AUTOCFT_CF_ACCOUNT_ID=${AUTOCFT_CF_ACCOUNT_ID}
      - AUTOCFT_CF_TUNNEL_ID=${AUTOCFT_CF_TUNNEL_ID}
      # 可选
      - AUTOCFT_CRON=*/30 * * * * *    # 默认 */10 * * * * * 每 10 秒
      - AUTOCFT_BASEDIR=/app/data
      #- AUTOCFT_ORIGIN_NO_TLS_VERIFY=true
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./autocft_data:/app/data
    restart: unless-stopped
```

启动：
```
docker compose up -d autocft
```

详情见 [环境变量](https://autocft.cloudfogtech.ltd/zh/environment)。

### 3. 给业务容器打标签
```yaml
  myapp:
    image: example/myapp:latest
    labels:
      # 必需
      - autocft.enabled=true
      - autocft.hostname=app.example.com
      # 可以解释为从cloudflared服务访问到该服务的地址
      - autocft.service=http://myapp:8080
      # 可选
      - autocft.path=/
      - autocft.origin.no-tls-verify=true
```
等待一个调度周期（默认 10 秒）后访问 `https://app.example.com` 即可通过 Tunnel 路由到内部服务。

详情见 [容器标签](https://autocft.cloudfogtech.ltd/zh/labels)。

## 🧪 构建与运行
本地构建二进制：
```
go build -o autocft ./cmd/autocft
./autocft
```

本地构建镜像：
```
docker build -t cloudfogtech/autocft:dev .
```
示例运行（直接二进制）：
```
AUTOCFT_CF_API_TOKEN=xxx \
AUTOCFT_CF_ACCOUNT_ID=xxx \
AUTOCFT_CF_TUNNEL_ID=xxx \
./autocft
```

## 🤝 贡献
欢迎提交 Issue 与 Pull Request：
1. Fork 仓库并创建分支
2. 保持代码风格一致
3. 添加必要注释 / 文档
4. 提交前本地构建并自测

### 架构简述
- 后端：Golang + PocketBase
- 前端：预留
- 文档：Cloudflare Pages + VitePress

## 📄 许可证
MIT License © 2025 CloudFogTech

---
**免责声明**：当前版本处于早期阶段（0.x），接口与行为可能出现不兼容变更，请在生产使用前充分评估风险。
