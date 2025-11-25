---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "Auto Cloudflare Tunnel"
  text: "在Docker环境中自动化部署 Cloudflare Tunnel 规则的工具"
  tagline: 简单、高效、易用
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/deploy

features:
  - title: 自动化
    details: 无需手动在Cloudflare Zero Trust中配置隧道，节省时间和精力
  - title: 秒级同步
    details: 默认每 10 秒执行一次，可配置更低的时间间隔
  - title: 低资源占用
    details: Golang 编写，单文件运行；无额外数据库，仅一个历史 JSON 文件
  - title: 安全合并
    details: 仅覆盖由标签驱动的主机名，Web 新增的规则自动保留
---

