# I56 Framework 2.0 LTS — 统一平台

```
I56 Framework 2.0 LTS
Enterprise Application Development Platform
Go 1.24+ | Modular Monolith | Cloud Native Ready
```

---

## 架构布局

```
i56/
├── framework/          # 核心框架 (17 packages)
│   ├── core/           # auth, rbac, router, cache, logger...
│   ├── ai/             # AI Runtime (Model Agnostic)
│   ├── cmd/            # server, migrate, worker, cli
│   ├── configs/        # YAML 配置
│   ├── pkg/            # 公开 SDK + 工具
│   ├── themes/         # 主题系统
│   ├── tests/          # e2e, integration, unit
│   └── deployments/    # docker, compose, helm, k8s
├── modules/            # 业务模块 (14 domains)
├── apps/wms/           # WMS 应用（主产品）
├── ui/                 # BDL 2.0 Web Components (Shadow DOM)
├── design-language/    # 设计令牌 / Token / Spec / Icon
├── deployments/        # 顶级部署配置
├── sdk/                # API SDK
├── docs/               # 文档
├── scripts/            # 构建/部署脚本
├── go.work             # Go Workspace
└── README.md
```

## 三端产品

| 端 | 路由 | 页面数 | 认证 |
|:--|:--|:--|:--|
| **Admin** | `/admin/*` | 53子菜单 | JWT + RBAC |
| **Client** | `/client/*` | 17菜单 | JWT (plat_*) |
| **PDA** | `/pda/*` | 15操作屏 | 工号+PIN |

## 核心功能

| 模块 | 内容 | 状态 |
|:--|:--|:--|
| 系统API配置 | 物流/清关/通知/打印/存储 | ✅ |
| 定价模型 | 5维 × 43记录 | ✅ |
| 工作流引擎 | 2流程12步骤 | ✅ |
| 客户端 | HMAC凭证 + Webhook | ✅ |
| PDA | 8任务抢单池 + 能力匹配 | ✅ |
| 包裹 | 14状态流转 | ✅ |

## 编译

```bash
cd i56/apps/wms
CGO_ENABLED=0 go build -ldflags="-s -w" -o i56-server ./cmd/server/
```

## 部署

```bash
scp i56-server ubuntu@106.52.164.139:/tmp/
ssh ubuntu@106.52.164.139 'sudo systemctl restart i56'
```

## 版本

- 当前: **I56 Framework 2.0 LTS**
- 生产: https://wms.mikaplay.com
- AI Runtime: active ✅
