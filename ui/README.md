# I56 UI — Web Component Library

> HTML5 + ES Modules + Web Components — 零框架依赖。

## 定位

`i56-ui` 是 I56 平台的 UI 层实现。基于 BDL 1.0 设计规范，提供原生 Web Components、页面模板和交互组件。

```
i56-design-language/  (设计标准)
        ↓
i56-ui/               ← 这个仓库 (组件实现)
        ↓
i56-apps/             (产品引用)
```

## 技术栈

```
HTML5 Custom Elements + Shadow DOM + ES Modules
CSS Custom Properties (BDL Tokens)
Zero dependencies (no React, Vue, Bootstrap)
```

## 组件清单 — 15 个 Web Components

| 组件 | 标签 | 功能 |
|------|------|------|
| Button | `<i56-button>` | primary/secondary/ghost/danger, sm/md/lg, loading, icon, 快捷键提示 |
| Card | `<i56-card>` | header/body/footer slots, hover, clickable, compact |
| Table | `<i56-table>` | JSON data, sortable columns, row selection, empty state |
| FormGroup | `<i56-form-group>` | label + input wrapper, error, hint, required |
| Input | `<i56-input>` | icon prefix, clearable, error state |
| Select | `<i56-select>` | custom dropdown, search filter, keyboard nav |
| Modal | `<i56-modal>` | 4 sizes, ESC close, backdrop click, focus trap |
| Toast | `<i56-toast>` | success/error/warning/info, 5s auto-dismiss, stacking |
| Badge | `<i56-badge>` | 6 colors, pill-shaped |
| Tabs | `<i56-tabs>` | underline style, keyboard nav (←→) |
| Avatar | `<i56-avatar>` | image or initials, deterministic color, 5 sizes |
| Spinner | `<i56-spinner>` | CSS animation, brand color, 4 sizes |
| Timeline | `<i56-timeline>` | completed/current/pending, connected line |
| Command | `<i56-command-palette>` | Ctrl+K, fuzzy search, keyboard nav |
| — | `I56Toast` / `I56Command` | global JS API |

## 设计模式

```
Shadow DOM (open)               — 样式隔离
CSS Custom Properties            — BDL 令牌驱动
CustomEvent 'i56:' prefix        — 事件体系
ARIA attributes                  — 无障碍
Keyboard accessibility           — Tab/Enter/Escape/Arrows
```

## 模板

```
templates/admin/
├── base_new.html                — BDL 布局（顶栏 + 侧栏 + 内容区 + AI Bar）
├── permissions_new.html         — 权限管理示例页
├── client_new.html              — 客户端布局
└── pda_new.html                 — PDA 布局
```

## 使用

```html
<link rel="stylesheet" href="/static/css/i56-bdl.css">
<script type="module" src="/static/js/i56-theme.js"></script>
<script type="module" src="/static/js/i56-components.js"></script>
<script type="module" src="/static/js/i56-command.js"></script>

<i56-table data='[{"id":1,"name":"权限一"}]' columns='["id","name"]'></i56-table>
<i56-button variant="primary" onclick="I56Toast.show({type:'success',message:'Done'})">
  保存
</i56-button>
```

## 许可证

MIT
