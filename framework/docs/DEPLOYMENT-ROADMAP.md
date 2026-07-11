# I56 Framework 1.0 LTS — 部署与开发路线图

---

## 一、三仓库初始化

```bash
# 1. Framework Core
mkdir i56-framework && cd i56-framework
git init && go mod init github.com/i56/framework
# → https://github.com/i56/framework

# 2. Admin Shell
mkdir i56-admin && cd i56-admin
git init && go mod init github.com/i56/admin
# → https://github.com/i56/admin

# 3. Apps
mkdir i56-apps && cd i56-apps
git init
# → https://github.com/i56/apps
#    ├── i56-wms/
#    ├── i56-oms/
#    ├── i56-tms/
#    └── ...
```

---

## 二、当前 I56 Framework 代码统计

```
目录结构：44 个目录
Go 源文件：17 个核心模块
代码行数：~2,500 行 Go 代码
构建状态：✅ go build ./... 零错误通过
运行时验证：
  ✅ GET /api/health → 200 {"data":{"status":"ok","name":"I56 Framework","version":"1.1.0"}}
  ✅ GET /api/v1/me → 200 {"data":{"message":"current user"}}
  ✅ GET /api/nonexistent → 404
```

### 已实现的 17 个 Core 模块

| # | 模块 | 文件 | 状态 |
|---|------|------|------|
| 1 | config | core/config/config.go | ✅ |
| 2 | logger | core/logger/logger.go | ✅ |
| 3 | errors | core/errors/errors.go | ✅ |
| 4 | response | core/response/response.go | ✅ |
| 5 | validator | core/validator/validator.go | ✅ |
| 6 | middleware | core/middleware/middleware.go | ✅ |
| 7 | router | core/router/router.go | ✅ |
| 8 | tenant | core/tenant/tenant.go | ✅ |
| 9 | auth | core/auth/auth.go | ✅ |
| 10 | rbac | core/rbac/rbac.go | ✅ |
| 11 | eventbus | core/eventbus/eventbus.go | ✅ |
| 12 | scheduler | core/scheduler/scheduler.go | ✅ |
| 13 | cache | core/cache/cache.go | ✅ |
| 14 | storage | core/storage/storage.go | ✅ |
| 15 | notification | core/notification/notification.go | ✅ |
| 16 | audit | core/audit/audit.go | ✅ |
| 17 | workflow | core/workflow/workflow.go | ✅ |

---

## 三、开发 Phase 执行计划

### Phase 1: Core 完善（4 周）

| 周 | 任务 | 产出 |
|----|------|------|
| W1 | config 完整实现（YAML/ENV/etcd） | config 模块 |
| W2 | auth JWT 完整实现（ed25519/RSA） | auth 模块 |
| W3 | rbac 数据库存储 + 中间件集成 | rbac 模块 |
| W4 | 全量单元测试（覆盖率 > 80%） | tests/unit/ |

### Phase 2: Admin Shell（4 周）

| 周 | 任务 |
|----|------|
| W5-6 | Go HTML 模板引擎（Bootstrap 5 + HTMX + Alpine.js）|
| W7-8 | 登录/角色/用户管理 UI + API 对接 |

### Phase 3: 业务模块（对标 BFT56）（12 周）

| 周 | 模块 |
|----|------|
| W9-10 | Customer（客户/会员/账号/申报人/地址）|
| W11-12 | Warehouse（仓库/库位/区域/集装柜）|
| W13-14 | Parcel（包裹/入库/状态/异常）|
| W15-16 | Order（集运订单/状态机）|
| W17-18 | Transport（线路/承运商/清关/物流）|
| W19-20 | Finance（充值/流水/账单/报表）|

### Phase 4: 增强（8 周）

- 打印模板引擎
- 通知中心全渠道
- Webhook 事件推送
- OpenAPI/SDK
- BI 报表引擎

---

## 四、Docker Compose 部署

```yaml
# deployments/compose/docker-compose.yml
version: "3.9"
services:
  server:
    build:
      context: ../..
      dockerfile: deployments/docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DATABASE_URL=postgres://i56:i56@postgres:5432/i56?sslmode=disable
      - REDIS_URL=redis://redis:6379/0
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: i56
      POSTGRES_PASSWORD: i56
      POSTGRES_DB: i56
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U i56"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:
```

---

## 五、Kubernetes Helm 部署（规划）

```
deployments/helm/i56/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   └── hpa.yaml
```

---

## 六、文档索引

| 文档 | 路径 | 说明 |
|------|------|------|
| BFT56 反向工程 PRD | `docs/BFT56-REVERSE-ENGINEERING-PRD.md` | 竞品完整分析 |
| I56 Framework 架构设计 | `docs/I56-FRAMEWORK-ARCHITECTURE.md` | 技术架构规格 |
| BFT56 vs I56 覆盖矩阵 | `docs/BFT56-I56-COVERAGE-MATRIX.md` | 功能对标记 |
| 部署与路线图 | `docs/DEPLOYMENT-ROADMAP.md` | 本文档 |

---

*生成时间：2026-07-10 | I56 Framework 1.0 LTS*
