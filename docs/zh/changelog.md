# 更新日志

## 0.2.3
- 支持 `matchSNItoHost` Origin 参数（保留 Web 通配符规则中设置的 SNI 配置）
- 升级 cloudflare-go 至 v6.10.0

## 0.2.2
- 修复通配符 hostname 排序，精确 hostname 优先匹配

## 0.2.1
- 更新 Docker 构建参数中的 Go 版本为 1.26

## 0.2.0
- 适配新的Cloudflare API
- 升级到 Go 1.26 版本
- 升级依赖

## 0.1.5
- 移除Dockerfile中默认的user

## 0.1.4
- 修复cron任务逻辑
- 升级依赖

## 0.1.3
- 修复Dockerfile中的启动问题
- 增加挂载目录的权限设置

## 0.1.2
- 调整容器的挂载权限设置

## 0.1.1
- 为了避免docker挂载错误将默认的目录修改为`/app/data`

## 0.1.0
- 修改默认admin账户的用户名和密码
- 发布第一个预览版本

## 0.0.3 版本
- 功能基本可用

## 0.0.2 版本
- 无业务逻辑修改
- 优化CI/CD流程

## 0.0.1 版本
- 初始原型
- Docker 标签发现
- Cloudflare Tunnel ingress 同步（合并 + diff）
- 自动追加 404 fallback
- 历史文件追踪
- 秒级 Cron 调度
