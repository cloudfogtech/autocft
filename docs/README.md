# AutoCFT 开发与部署指南

本文件为 AutoCFT 项目的开发者提供开发和部署相关的说明，帮助你快速上手项目的本地开发、构建与文档维护流程。

## 目录结构说明

```
├── assets/           # 静态资源（图片、图表等）
│   ├── diagrams/
│   └── images/
├── en/               # 英文文档
├── zh/               # 中文文档
├── index.md          # 文档入口
├── package.json      # Node.js 项目配置
├── pnpm-lock.yaml    # pnpm 依赖锁定文件
└── README.zh.md      # 中文开发文档（本文件）
```

## 开发环境准备

1. **克隆仓库**
   ```bash
   git clone <你的仓库地址>
   cd autocft/docs
   ```

2. **安装依赖**
   推荐使用 pnpm 进行依赖管理：
   ```bash
   pnpm install
   ```
   如未安装 pnpm，可先通过 npm 安装：
   ```bash
   npm install -g pnpm
   ```

3. **文档本地预览与开发**
   项目文档建议使用 [VitePress](https://vitepress.dev/) 或类似工具进行本地预览和构建。请根据实际配置调整命令。

   启动本地文档服务：
   ```bash
   pnpm docs:dev
   ```

   构建静态文档：
   ```bash
   pnpm docs:build
   ```

## 目录与文档维护

- 所有文档内容请分别维护于 `en/`（英文）和 `zh/`（中文）目录下。
- 静态资源统一放置于 `assets/` 目录。
- 文档入口为 `index.md`，请保持中英文文档结构一致，便于维护。

## 依赖管理

- 推荐使用 pnpm 进行依赖管理，确保依赖一致性。
- 如需新增依赖，请同步更新 `package.json` 并提交 `pnpm-lock.yaml`。

## 贡献指南

- 欢迎通过 Pull Request 贡献文档或改进开发流程。
- 提交前请确保文档格式规范、内容准确。
- 如有疑问请先查阅 `zh/faq.md` 或提交 Issue。



## 许可证

MIT License © 2025 CloudfogTech
