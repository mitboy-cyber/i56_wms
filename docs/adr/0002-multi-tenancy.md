# ADR-0002: Shared Table Multi-Tenancy for 1.0 LTS

## 状态
Accepted (2026-07-15)

## 背景
I56 Framework 是 SaaS 平台，需要支持多个租户（企业客户）数据隔离。BFT56 对标竞品采用 Shared Table 方案。

## 决策
**1.0 LTS 采用 Shared Table 方案**，所有租户共享同一组数据库表，通过 `tenant_id` 字段隔离。

## 实现
```go
// TenantProvider 接口
type TenantProvider interface {
    CurrentTenantID(ctx context.Context) int64
}

// 所有仓储查询自动注入 tenant_id
func (r *OrderRepo) List(ctx context.Context) ([]Order, error) {
    return r.db.Where("tenant_id = ?", tenant.FromContext(ctx)).Find(&orders)
}
```

## 替代方案

| 方案 | 适用场景 | 1.0 评估 |
|------|---------|---------|
| Shared Table | 中小租户 (≤1000) | ✅ 简单，快速启动 |
| Schema Per Tenant | 中型租户 (≤10000) | 需要 PostgreSQL，Go 1.24 生态不成熟 |
| Database Per Tenant | 大型企业客户 | 运维成本高，1.0 不需要 |

## 数据隔离规则

```sql
-- 所有表必须包含
tenant_id BIGINT NOT NULL INDEX

-- 所有查询必须带租户条件
WHERE tenant_id = $1

-- 中间件自动注入
ctx = tenant.WithTenantID(ctx, claims.TenantID)
```

## 演进路径
- 1.0 LTS: Shared Table
- 2.0 LTS: 支持 Schema Per Tenant (PG 专用)
- 3.0 LTS: 完整多策略 Provider 切换

## 后果
- 开发简单，不需要额外的数据库连接管理
- 数据隔离依赖于应用层（所有查询必须带 tenant_id）
- 大租户可能影响其他租户性能（可通过读写分离缓解）
