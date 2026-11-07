# Railway 部署说明

本项目已适配 [Railway](https://railway.com/)：`railway.json` 使用 Dockerfile 构建，
健康检查走 `/healthz`（超时 300s，启动时 v8go 初始化较慢属正常）。

## 两种部署方式

**A. GitHub 仓库部署**（代码推送后自动重新部署）
Railway → New Project → Deploy from GitHub repo → 选择本仓库。

**B. Railway CLI 本地部署**
```bash
railway login
railway init --name nvidia-proxy
railway up
railway domain
```

## 必设环境变量

| 变量 | 说明 |
|------|------|
| `API_KEY` | **必设**。入站鉴权密钥；客户端用 `Authorization: Bearer <key>` 或 `x-api-key: <key>` 调用。不设则服务对所有人生效（仅限内网时才可接受） |

## 可选环境变量

| 变量 | 说明 |
|------|------|
| `PORT` | 监听端口，Railway 自动注入，无需手动设置 |
| `CHROME_PROXY` | PoW/模型目录/WAF token 走的代理（Railway 出口 IP 被上游风控时用，如 `socks5://host:port`） |

注意：上游是 NVIDIA 匿名 Playground 接口，随时可能变动；数据中心 IP 更容易触发
WAF/验证码。`-catalog-cookie` 等参数仅 CLI 可用，镜像入口为 `serve -auto`。
