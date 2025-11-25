import {defineConfig} from 'vitepress'
import { configureDiagramsPlugin } from "vitepress-plugin-diagrams";

// https://vitepress.dev/reference/site-config
export default defineConfig({
    base: '/autocft/',
    title: "Auto Cloudflare Tunnel",
    description: "A tool to auto deploy applications for Cloudflare Tunnel in Docker with Docker Compose",
    themeConfig: {
        // https://vitepress.dev/reference/default-theme-config
        socialLinks: [
            {icon: 'github', link: 'https://github.com/cloudfogtech/autocft'}
        ]
    },
    locales: {
        root: {
            label: 'English',
            lang: 'en',
            link: '/en',
            themeConfig: {
                nav: [
                    {text: 'Home', link: '/en'},
                    {text: 'Deploy', link: '/en/deploy'},
                    {text: 'FAQ', link: '/en/faq'}
                ],
                sidebar: [
                    {
                        text: 'Getting Started',
                        items: [
                            {text: 'Deployment', link: '/en/deploy'},
                            {text: 'Environment Variables', link: '/en/environment'},
                            {text: 'Container Labels', link: '/en/labels'},
                            {text: 'Cloudflare', link: '/en/cloudflare'},
                        ]
                    },
                    {
                        text: 'Deep Dive',
                        items: [
                            {text: 'Architecture', link: '/en/architecture'},
                            {text: 'Troubleshooting', link: '/en/troubleshooting'},
                            {text: 'FAQ', link: '/en/faq'},
                            {text: 'Roadmap', link: '/en/roadmap'},
                        ]
                    },
                    {
                        text: 'Project',
                        items: [
                            {text: 'Changelog', link: '/en/changelog'}
                        ]
                    }
                ],
                footer: {
                    message: 'Released under the MIT License.',
                    copyright: 'Copyright © 2025 CloudfogTech'
                }
            }
        },
        zh: {
            label: '简体中文',
            lang: 'zh',
            link: '/zh',
            themeConfig: {
                nav: [
                    {text: '首页', link: '/zh'},
                    {text: '部署', link: '/zh/deploy'},
                    {text: '常见问题', link: '/zh/faq'}
                ],
                sidebar: [
                    {
                        text: '快速开始',
                        items: [
                            {text: '部署指南', link: '/zh/deploy'},
                            {text: '环境变量', link: '/zh/environment'},
                            {text: '容器标签', link: '/zh/labels'},
                            {text: 'Cloudflare参数获取', link: '/zh/cloudflare'},
                        ]
                    },
                    {
                        text: '深入理解',
                        items: [
                            {text: '架构设计', link: '/zh/architecture'},
                            {text: '故障排查', link: '/zh/troubleshooting'},
                            {text: '常见问题', link: '/zh/faq'},
                            {text: '发展路线图', link: '/zh/roadmap'},
                        ]
                    },
                    {
                        text: '项目',
                        items: [
                            {text: '更新日志', link: '/zh/changelog'}
                        ]
                    }
                ],
                footer: {
                    message: 'Released under the MIT License.',
                    copyright: 'Copyright © 2025 CloudfogTech'
                }
            }
        }
    },
    markdown:{
        config: (md) => {
            configureDiagramsPlugin(md, {
                diagramsDir: "assets/diagrams", // 可选：自定义 SVG 文件目录
                publicPath: "/autocft/assets/diagrams", // 可选：自定义公共路径
                krokiServerUrl: "https://kroki.io", // 可选：自定义 Kroki 服务器地址
                excludedDiagramTypes: [], // 可选：排除特定图表类型
            });
        },
    }
})
