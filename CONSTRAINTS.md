# I56 Framework — 开发约束 (Constraints)

> 最后更新: 2026-07-12  
> 适用范围: `/home/ubuntu/i56/` 下所有项目

---

## 1. Git 工作流 (强制)

### 1.1 每次部署前必须推送

```
部署流程:
  1. git add -A
  2. git commit -m "描述本次变更"
  3. git push origin main
  4. 创建 Pull Request
  5. 构建 & 部署
```

**绝对禁止**: 部署后补push。部署前代码必须在GitHub上。

### 1.2 GitHub 仓库

| 项目 | 仓库 |
|:--|:--|
| I56 Framework | https://github.com/mitboy-cyber/i56_wms.git |
| 分支 | `main` |

### 1.3 PR 规范

- 每次部署前从 `main` 创建 PR
- PR 标题: `v{版本号}: {简述}`
- PR 描述列出: 新增功能、修复Bug、文件变更数

---

## 2. 代码规范

### 2.1 语言 & 框架

- 后端: Go 1.22+
- 前端: HTML + CSS (BDL) + HTMX + Vanilla JS
- 模板: Go `html/template`
- 禁止: Bootstrap, jQuery, React, Vue

### 2.2 目录结构

```
i56/
├── framework/          # 核心框架 (独立)
│   └── core/
│       ├── scheduler/
│       ├── audit/
│       ├── storage/
│       ├── cache/
│       ├── router/
│       └── report/
├── apps/
│   └── wms/           # WMS 应用
│       ├── cmd/server/
│       ├── internal/
│       │   ├── common/
│       │   ├── crmroute/
│       │   ├── ezway/
│       │   ├── omsroute/
│       │   ├── sysroute/
│       │   ├── tmsroute/
│       │   └── wmsroute/
│       ├── static/
│       ├── templates/
│       └── modules/
├── deployments/
│   ├── compose/
│   └── init/
└── docs/
```

### 2.3 命名

- 路由: `/admin/{module}/{action}` RESTful
- 模板: `{module}_{action}.html`
- 函数: PascalCase (导出) / camelCase (内部)

---

## 3. UI/UX 规范

### 3.1 BDL (Business Design Language)

```
主题:       Light (#F8FAFC)
主色:       Inter Blue #1D4ED8
字体:       Inter / system-ui
禁止:       暗色主题
```

### 3.2 三端统一

- Admin: `/admin` — 管理后台
- Client: `/client` — 客户门户  
- PDA: `/pda` — 手持终端

**三端数据必须互通，共享同一后端。**

### 3.3 移动端

- 所有页面必须响应式
- PDA 页面: 320-428px 移动优先
- 汉堡菜单: 移动端自动折叠

---

## 4. 功能完整性

### 4.1 CRUD

所有列表页必须具备:
- ✅ 添加 (Add)
- ✅ 编辑 (Edit) — 弹窗模态框
- ✅ 删除 (Delete) — 确认提示
- ✅ 搜索 (Search)
- ✅ 排序 (Sort)
- ✅ 批量操作 (Batch)

### 4.2 编辑功能

- 使用 HTMX `hx-post` + `HX-Redirect`
- 保存后自动刷新列表
- 编辑表单预填现有数据
- 失败显示错误提示

---

## 5. 部署

### 5.1 生产环境

```
服务器:    wms.mikaplay.com (106.52.164.139)
二进制:    /opt/i56/i56-server
模板:      /opt/i56/templates/
服务:      systemctl restart i56
```

### 5.2 部署命令

```bash
cd /home/ubuntu/i56/apps/wms
CGO_ENABLED=0 go build -ldflags="-s -w" -o /tmp/i56-server ./cmd/server/
tar czf /tmp/tpl.tar.gz templates/
sshpass -p 'PASSWORD' scp /tmp/i56-server ubuntu@106.52.164.139:/tmp/
sshpass -p 'PASSWORD' scp /tmp/tpl.tar.gz ubuntu@106.52.164.139:/tmp/
sshpass -p 'PASSWORD' ssh ubuntu@106.52.164.139 'sudo systemctl stop i56; sudo cp /tmp/i56-server /opt/i56/; cd /opt/i56 && sudo tar xzf /tmp/tpl.tar.gz; sudo systemctl restart i56'
```

---

## 6. 禁止事项

- ❌ 部署前不push代码
- ❌ 跳过PR直接部署
- ❌ 使用暗色主题
- ❌ 引入Bootstrap/jQuery
- ❌ 硬编码中文enums（必须用映射函数）
- ❌ 编辑/保存不走真实数据库操作
- ❌ 三端数据不一致

---

## 7. 待补充 (用户指定)

<!-- 用户将在此处添加更多约束 -->
