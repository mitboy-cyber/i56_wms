# ADR-0003: React SPA + TanStack Query for Admin Frontend

## 状态
Accepted (2026-07-15)

## 背景
管理后台需要支持 CRUD、实时数据、报表图表等交互密集操作。BFT56 采用的是 SSR + jQuery 方案，I56 需要更现代的体验。

## 决策
**管理后台采用 React SPA + TanStack Query + Tailwind CSS**。

```
技术栈:
  React 18+          — 组件化 UI
  TypeScript         — 类型安全
  TanStack Query     — 服务端状态管理
  Zustand            — 客户端状态
  Tailwind CSS 4     — 实用优先样式
  Vite               — 构建工具
  Recharts/Chart.js  — 图表库
```

## 架构
```
/src
├── api/client.ts          — API 客户端 (axios)
├── components/            — 通用组件
│   └── GenericListPage    — 通用 CRUD 列表 (76页共享)
├── pages/
│   ├── admin/             — 管理后台 76 页
│   └── client/            — 客户端 19 页
├── stores/                — Zustand stores
└── hooks/                 — 自定义 hooks
```

## 核心模式

### GenericListPage
所有 CRUD 页面通过 `GenericListPage` 组件统一实现：
```tsx
<GenericListPage
  title="集运订单"
  queryKey={['admin-orders']}
  queryFn={() => client.get('/admin/api/orders')}
  apiBase="/admin/api/orders"
  columns={[...]}
/>
```

自动提供：中文状态徽章、货币格式化、日期格式化、增删改查。

### 自动 API 基路径 (autoApiBase)
当 `apiBase` 未显式指定时，从 `queryFn` 中自动提取 API 路径。

## 替代方案

| 方案 | 评估 |
|------|------|
| SSR (Go Template + HTMX) | 适合简单页面，不适合复杂报表交互 |
| Vue + Nuxt | 生态不如 React 丰富 |
| Angular | 过重，学习曲线陡 |
| **React SPA** | ✅ 生态最丰富，类型安全，组件复用 |

## 后果
- 前端独立部署，与后端通过 API 通信
- 需要处理 CORS 和认证 token 管理
- 前端构建产物 ~74 JS chunks，首屏加载需优化
