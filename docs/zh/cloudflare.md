# Cloudflare 凭据获取指南

本文介绍如何在 Cloudflare Dashboard 中获取并配置以下三个核心参数：

- `CF_API_TOKEN` （API 令牌，推荐最小权限）
- `CF_ACCOUNT_ID` （账号 ID）
- `CF_TUNNEL_ID` （Cloudflare Tunnel ID）

---
## 目录
1. [概念速览](#1-概念速览)
2. [准备工作](#2-准备工作)
3. [获取 CF_ACCOUNT_ID](#3-获取-cf_account_id)
4. [创建最小权限 CF_API_TOKEN](#4-创建最小权限-cf_api_token)
5. [获取或创建 Cloudflare Tunnel 并获得 CF_TUNNEL_ID](#5-获取或创建-cloudflare-tunnel-并获得-cf_tunnel_id)
6. [验证配置](#6-验证配置)
7. [常见问题与排错](#7-常见问题与排错)
8. [安全注意事项](#8-安全注意事项)

---
## 1. 概念速览
- `CF_ACCOUNT_ID`：Cloudflare 账户的唯一标识，用于调用部分 Account 级别 API。
- `CF_API_TOKEN`：基于权限范围（Scopes）的现代 API 令牌，优先使用它而不是 Global API Key。
- `CF_TUNNEL_ID`：Cloudflare Tunnel（原 Argo Tunnel）实例的唯一标识，用于通过 API 管控该隧道。

> 强烈不建议在自动化中使用 "Global API Key"。如必须临时使用，请限制环境并尽快改用最小权限 Token。

---
## 2. 准备工作
确保：
- 已注册并登录 Cloudflare(登录 https://dash.cloudflare.com/)。
- 已至少在一个域名上启用 Cloudflare（不强制，但后续若要关联 DNS 记录需要域名）。
- 已计划好需要的权限（Tunnel 只读 / 可创建与管理 Tunnel / 是否需要 DNS Edit）。

> 本教程完全使用 API 方式，无需安装 `cloudflared` CLI。

---
## 3. 获取 CF_ACCOUNT_ID

![账号ID路径](/assets/images/cloudflare-01.png)

### 步骤
1. 登录 Cloudflare Dashboard： https://dash.cloudflare.com/
2. 在左侧侧边栏或主页中，进入任意一个 Zone（域名）详情页。
3. 浏览器地址栏 URL 形如：`https://dash.cloudflare.com/<ACCOUNT_ID>/home/domains`。
4. 复制 `<ACCOUNT_ID>` 作为 `CF_ACCOUNT_ID`。


### 备用方式：API 列出账号
使用 API Token（需具备 Account 查看权限）调用：
```
GET https://api.cloudflare.com/client/v4/accounts
Authorization: Bearer <CF_API_TOKEN>
```
返回 JSON 中 `result[0].id` 即为账号 ID（多个账户时选择需要的）。

---
## 4. 创建最小权限 CF_API_TOKEN

### 目标权限（根据需求增减）：
- 账户 -> Cloudflare Tunnel:  编辑（需要读取/修改Tunnel Ingress的参数）

> 视你的自动化需求可只保留 Tunnel 相关权限。若只读取 Tunnel 状态，可全部设为 Read。

### 步骤
1. 右上角头像下拉，进入 `配置文件`。
   ![进入配置文件](/assets/images/zh/cloudflare-02.png)
2. 左侧导航选择 `API令牌`。
3. 点击 `创建令牌`。
   ![进入API令牌](/assets/images/zh/cloudflare-03.png)
4. 选择 `创建自定义令牌`。
   ![创建自定义令牌](/assets/images/zh/cloudflare-04.png)
5. 填入`令牌名称`，如 "AutoCFT Token"。
6. 在`权限`中加入：
   - `账户` / `Cloudflare Tunnel` / `编辑`
7. 在 `账户资源` 选择条件为`包括`，账户现在需要使用的账户（如果仅有一个账户，选择`所有账户`也没有问题）。
8. `客户端IP地址筛选` 和 `TTL` 可根据需要限制来源 IP，默认不填。
   ![填入信息](/assets/images/zh/cloudflare-05.png)
9. 点击 `继续以显示摘要` 到最后确认，创建并复制 Token。
    ![显示摘要](/assets/images/zh/cloudflare-06.png)
10. 复制后生成的 Token 即 `CF_API_TOKEN`（只显示一次，务必保存）。
    ![显示摘要](/assets/images/zh/cloudflare-07.png)

---
## 5. 获取或创建 Cloudflare Tunnel 并获得 CF_TUNNEL_ID

可以通过 Dashboard 或 API。两者选其一即可。

### A. Dashboard 查看
1. 主页左侧导航搜索或展开 "Zero Trust"。
   ![Network -> Connectors](/assets/images/zh/cloudflare-08.png)
2. 进入 "Tunnels"（路径：Zero Trust > 网络 > Connector(Tunnel)）。
3. 列表中每个 Tunnel 一行，点击需要的 Tunnel。
   ![Tunnel列表](/assets/images/zh/cloudflare-09.png)
4. 详情页中或地址栏可见 Tunnel 的 UUID，即 `CF_TUNNEL_ID`。
   ![Tunnel详情](/assets/images/zh/cloudflare-10.png)

### B. 使用 API 列出已有 Tunnel
请求Tunnel信息，CF_API_TOKEN为上一步获取：
```
curl -X GET "https://api.cloudflare.com/client/v4/accounts/<CF_ACCOUNT_ID>/cfd_tunnel" \
  -H "Authorization: Bearer <CF_API_TOKEN>" \
  -H "Content-Type: application/json"
```
示例响应片段：
```
{
  "success": true,
  "result": [
    {
      "id": "3a1b2c3d-4e5f-6789-abcd-0123456789ef",
      "name": "my-tunnel",
      "created_at": "2025-10-01T02:33:11Z",
      "connections": [...]
    }
  ]
}
```
`result` 数组中每个对象的 `id` 即为 `CF_TUNNEL_ID`。若数组为空表示当前账号尚无 Tunnel。

---
## 6. 验证配置

### 环境变量示例（Windows CMD）
```
set CF_API_TOKEN=xxxxxxxxxxxxxxxx
set CF_ACCOUNT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
set CF_TUNNEL_ID=3a1b2c3d-4e5f-6789-abcd-0123456789ef
```

### 使用 API 测试 Token 与 Tunnel
```
curl -X GET "https://api.cloudflare.com/client/v4/accounts/%CF_ACCOUNT_ID%/cfd_tunnel/%CF_TUNNEL_ID%" ^
  -H "Authorization: Bearer %CF_API_TOKEN%" ^
  -H "Content-Type: application/json"
```
返回中应包含对应的 `id` 与名称。

### 列出全部（再次确认权限）
```
curl -X GET "https://api.cloudflare.com/client/v4/accounts/%CF_ACCOUNT_ID%/cfd_tunnel" ^
  -H "Authorization: Bearer %CF_API_TOKEN%" ^
  -H "Content-Type: application/json"
```

### 常见返回码
- 200：成功
- 403：Token 权限不足（检查是否缺少 Cloudflare Tunnel 权限）
- 401：Token 无效或已撤销

---
## 7. 常见问题与排错

| 问题 | 可能原因 | 解决建议 |
|------|----------|----------|
| 403 Forbidden | Token 缺少 Tunnel 权限 | 重新创建 Token，添加 Account -> Cloudflare Tunnel -> Read/Edit |
| 返回空列表 | 尚未创建 Tunnel | 使用 POST 创建新 Tunnel |
| 查不到期望的 Tunnel | 使用了错误的 Account ID | 确认 URL 中的 ACCOUNT_ID 与 Token 权限范围一致 |
| DNS 记录未自动更新 | 没有 Zone DNS Edit 权限 | 添加 Zone / DNS / Edit 或手动配置记录 |
| UUID 复制错误 | 粘贴时包含多余字符 | 仅复制纯 36 字符 UUID（包含 4 个短横线） |
| "Global API Key" 被误用 | 没有区分令牌类型 | 改用最小权限 API Token 并撤销不必要的全局密钥 |

---
## 8. 安全注意事项
- 最小权限原则：只授予实际需要的权限范围与资源范围。
- 不要把 Token 写入公开仓库；使用环境变量或秘密管理（Secrets）。
- 定期轮换 Token，并在用途结束后撤销无用 Token。
- 对 CI/CD：使用平台提供的秘密管理（如 GitHub Actions Secrets）。
- 日志中避免打印完整 Token；必要时只显示前后若干字符。

---
## 快速检查清单
- 已获取正确的 `CF_ACCOUNT_ID`（URL / API 列表）
- 已创建并保存最小权限 `CF_API_TOKEN`
- 已通过 API 列出或创建 Tunnel 获得 `CF_TUNNEL_ID`
- 已设置环境变量供`AutoCFT`使用
- 已用 API 单独 GET 验证 Tunnel 可访问

> 完成以上步骤即可在`AutoCFT`中安全使用 Cloudflare 相关功能。
