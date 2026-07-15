# I56 WMS — 企业级仓库管理系统

> I56 Framework 1.0 LTS 的首个业务产品

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.22+ |
| 前端 | React 18 + TypeScript + Vite |
| 样式 | Tailwind CSS 4 |
| 状态 | TanStack Query + Zustand |
| 图表 | Recharts |
| 图标 | Lucide React |
| 部署 | Docker + Nginx + systemd |

## 快速启动

```bash
# 后端
cd apps/wms
go run ./cmd/server/

# 前端 (开发模式)
cd apps/wms/frontend
npm install
npm run dev

# 构建
npm run build
```

## 访问

```
管理后台: https://wms.mikaplay.com/admin
客户端:   https://wms.mikaplay.com/client
```

## 文档

- [产品需求文档](docs/PRD-I56-Framework-1.0-LTS.md)
- [领域词汇表](docs/CONTEXT.md)
- [核心架构](docs/ARCHITECTURE.md)
- [部署架构](docs/DEPLOYMENT.md)
- [API 参考](docs/API-REFERENCE.md)
- [架构决策记录](docs/adr/)

## 项目结构

```
apps/wms/
├── cmd/server/         # 主程序入口
├── internal/
│   ├── adminapi/       # 管理 API
│   ├── clientapi/      # 客户端 API
│   ├── pdaapi/         # PDA API
│   ├── domain/         # 领域模型 + 种子数据
│   └── server/         # HTTP 服务器
├── frontend/
│   ├── src/
│   │   ├── pages/admin/    # 76 管理页面
│   │   ├── pages/client/   # 8 客户端页面
│   │   ├── components/     # 通用组件
│   │   └── api/            # API 客户端
│   └── package.json
├── docker/
└── docs/
```

## 版本历史

| 版本 | 日期 | 里程碑 |
|------|------|--------|
| v1-v78 | 2026-07 | 管理后台 MVP + 中文化 |
| v79-v98 | 2026-07 | BFT56对齐 + 运营看板 + 报表 |
| v99-v114 | 2026-07 | 客户端门户 + 完整业务逻辑 + 文档 |
