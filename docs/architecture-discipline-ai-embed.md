# I56 Framework — 架构约束与回归防治分析

> 2026-07-12 | 深度诊断

---

## 一、问题诊断：Bootstrap 回归根因

### 现象
每次迭代新增页面/功能，自动出现 Bootstrap 残留，没有遵循 BDL (Business Design Language)。

### 根因分析

#### 1. 无单一样式入口
```
当前样式加载方式（混乱）:
├── templates/base.html       → 引用 i56-bdl.css ✅
├── templates/pda/base.html   → 独立CSS
├── templates/client/base.html→ 独立CSS
├── admin_layout.html         → 继承 base.html ✅
└── 内联页面                  → 各自内嵌<style> ❌
```
**任何新页面如果不经过 `base.html` 或 `admin_layout.html`，就会丢失 BDL。**

#### 2. 三种渲染路径并存
```
路径A: Go Template (base.html/admin_layout.html) → BDL ✅
路径B: common.RenderAdminPage() → 内联HTML → 可能丢失BDL ⚠️
路径C: fmt.Fprintf 直接写 HTML → 几乎必然丢失BDL ❌
```
**开发者看到路径B/C可用，就会复制粘贴，导致Bootstrap回归。**

#### 3. 无 CSS 命名空间保护
```
Bootstrap: .btn .table .card .modal .navbar → 全局污染
BDL:       .i56-btn .i56-table .i56-card → 有前缀保护
```
**没有 BDL CSS 的强制前缀策略，Bootstrap 类名可以直接使用。**

#### 4. 无构建时检查
```
go build → 通过（Go不检查CSS）
无 linter → CSS 没有规范检查
无 template 审查 → HTML 没有模板完整性检查
```

---

## 二、根本解决方案

### 架构约束（强制执行）

#### 约束 1：单一路径法则
```
❌ 删除:    fmt.Fprintf 输出 HTML
❌ 删除:    common.RenderAdminPage()
✅ 唯一:    Go Template (base.html → admin_layout.html → 具体页面)
```

#### 约束 2：BDL 命名空间强制
```
所有CSS类名加 .i56- 前缀:
  .i56-btn .i56-table .i56-card .i56-modal .i56-badge .i56-form

检查脚本: grep -r 'class="[^"]*\b(btn|table|card|modal|navbar|row|col)\b[^"]*"' templates/
→ 如果匹配但无 .i56- 前缀 → 构建失败
```

#### 约束 3：模板继承链
```
所有页面必须继承以下之一:
├── templates/base.html           → 全局基础 (字体/主题/Footer)
│   └── admin_layout.html         → 管理后台 (侧边栏/面包屑)
│       └── data_table.html       → 数据列表
└── templates/client_base.html    → 客户门户
└── templates/pda_base.html       → PDA终端
```

#### 约束 4：CSS 审计钩子
```bash
# pre-commit hook (放入 .git/hooks/pre-commit)
if grep -rq 'class=".*\(btn\|table\|card\|modal\|navbar\)\([^-]\|$\)' templates/; then
    echo "❌ Bootstrap class detected! Use .i56-* instead."
    exit 1
fi
```

---

## 三、AI 助手嵌入方案深度研判

### 需求
AI 助手应当出现在 admin/client **每个页面**，而非独立页面。

### 方案对比

| 方案 | 可见性 | 侵入性 | 实现复杂度 | 用户体验 |
|:--|:--|:--|:--|:--|
| A: 浮动聊天气泡 | ⭐⭐⭐⭐⭐ | 低 | 中 | 最佳(不遮挡) |
| B: 侧边栏面板 | ⭐⭐⭐ | 高 | 中 | 一般(占空间) |
| C: 底部状态栏 | ⭐⭐ | 低 | 低 | 差(小屏幕) |
| D: 全局快捷键弹窗 | ⭐⭐ | 低 | 低 | 不可发现 |
| E: 页面内嵌 HTMX 组件 | ⭐⭐⭐ | 中 | 高 | 好(需改每页) |

### 推荐：**方案 A — 浮动聊天气泡 + HTMX**

```
┌─────────────────────────────────────┐
│ I56 Admin                        🔔 │
├─────────────────────────────────────┤
│ 侧边栏 │         页面内容           │
│        │                            │
│        │     ┌──────────┐           │
│        │     │  💬 AI   │  ← 浮动   │
│        │     │  助手    │   气泡     │
│        │     └──────────┘           │
│        │                            │
├─────────────────────────────────────┤
│ © 2026 I56 Framework v2.4 LTS       │
└─────────────────────────────────────┘
```

### 实现架构

```html
<!-- 在 base.html 底部添加，所有页面自动继承 -->
<div id="i56-ai-bubble" class="i56-ai-bubble">
  <button id="i56-ai-toggle" class="i56-ai-toggle" 
          onclick="I56AI.toggle()" title="AI 助手">
    🤖
  </button>
  <div id="i56-ai-panel" class="i56-ai-panel" style="display:none">
    <div class="i56-ai-header">
      <span>🤖 I56 AI 助手</span>
      <button onclick="I56AI.close()">✕</button>
    </div>
    <div id="i56-ai-messages" class="i56-ai-messages"
         hx-ext="sse" sse-connect="/api/v1/ai/chat-stream"
         sse-swap="aiResponse">
      <div class="i56-ai-welcome">您好！我可以帮您查询订单、包裹、客户等信息。</div>
    </div>
    <form class="i56-ai-input" onsubmit="I56AI.send(event)">
      <input type="text" id="i56-ai-query" placeholder="输入问题..." />
      <button type="submit">➤</button>
    </form>
  </div>
</div>
```

### JavaScript (I56AI 对象)

```javascript
var I56AI = {
  toggle: function() {
    var p = document.getElementById('i56-ai-panel');
    p.style.display = p.style.display === 'none' ? 'flex' : 'none';
  },
  close: function() {
    document.getElementById('i56-ai-panel').style.display = 'none';
  },
  send: function(e) {
    e.preventDefault();
    var q = document.getElementById('i56-ai-query').value;
    if (!q) return;
    var msgs = document.getElementById('i56-ai-messages');
    msgs.innerHTML += '<div class="i56-ai-msg-user">' + q + '</div>';
    document.getElementById('i56-ai-query').value = '';
    fetch('/api/v1/ai/chat', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({message: q})
    }).then(function(r) { return r.text(); })
      .then(function(t) {
        msgs.innerHTML += '<div class="i56-ai-msg-ai">' + t + '</div>';
        msgs.scrollTop = msgs.scrollHeight;
      });
  }
};
```

### CSS (嵌入 base.html)

```css
.i56-ai-bubble { position:fixed; bottom:20px; right:20px; z-index:9998; }
.i56-ai-toggle { width:48px; height:48px; border-radius:50%; background:var(--i56-brand,#1D4ED8); color:#fff; border:none; font-size:22px; cursor:pointer; box-shadow:0 2px 12px rgba(0,0,0,0.2); }
.i56-ai-panel { position:fixed; bottom:80px; right:20px; width:380px; max-height:500px; background:var(--i56-bg-surface,#fff); border-radius:12px; box-shadow:0 4px 24px rgba(0,0,0,0.15); display:flex; flex-direction:column; overflow:hidden; }
.i56-ai-header { padding:12px 16px; border-bottom:1px solid var(--i56-border); font-weight:600; display:flex; justify-content:space-between; }
.i56-ai-messages { flex:1; overflow-y:auto; padding:12px; max-height:350px; }
.i56-ai-input { display:flex; padding:8px; border-top:1px solid var(--i56-border); }
.i56-ai-input input { flex:1; border:1px solid var(--i56-border); border-radius:6px; padding:8px 12px; }
.i56-ai-input button { background:var(--i56-brand,#1D4ED8); color:#fff; border:none; border-radius:6px; padding:8px 14px; margin-left:6px; cursor:pointer; }
```

### 关键优势

| 特性 | 说明 |
|:--|:--|
| **全局可用** | base.html 加载，所有页面自动继承 |
| **零页面改动** | 不需要修改任何业务页面 |
| **BDL 兼容** | 使用 .i56-* 前缀，不冲突 |
| **可关闭/打开** | 浮动气泡，不遮挡内容 |
| **SSE 流式** | 聊天消息实时推送 |
| **响应式** | 移动端自动适配宽度 |

---

## 四、实施计划

### Phase 1: 约束强制执行
1. 在 `.git/hooks/pre-commit` 添加 CSS 检查
2. 标记 renderAdminPage() 为 DEPRECATED
3. 所有新页面强制使用 admin_layout 模板

### Phase 2: AI 气泡嵌入
1. 在 base.html 添加浮动气泡 HTML + CSS + JS
2. 在 base.html 添加 I56AI JS 对象
3. 确保三端（admin/client/pda）风格一致

### Phase 3: 清理回归
1. grep 全项目 Bootstrap 残留
2. 逐个迁移到 BDL
3. 移除 Bootstrap CDN 引用
