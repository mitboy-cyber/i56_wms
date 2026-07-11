# I56 Design Language (BDL) 1.0

> 设计标准 — 不是代码，是规范。

## 定位

BDL（I56 Design Language）是 I56 平台的设计标准，定义了所有 I56 产品的视觉语言、组件规范和交互模式。

```
i56-design-language/          ← 这个仓库
    ↓ 规范驱动
i56-ui/                       ← 组件实现
    ↓ 组件引用
i56-apps/                     ← 产品交付
```

## 版本

| 版本 | 状态 | 说明 |
|------|------|------|
| BDL 1.0 | ✅ 发布 | 首个版本，定义了完整的令牌体系 |

## 核心原则

```
Minimal     — 极简布局，减少视觉噪音
Consistent  — 全局统一间距、字体、图标和交互
AI Native   — AI 助手是一级能力，不是外挂
Keyboard First — 全局快捷键与 Command Palette
Fast        — 首屏快、动画轻、响应快
Modular     — 每个页面由可组合的卡片和模块构成
Responsive  — 同一套 HTML5 适配 PC、平板和移动端
Accessible  — 高对比度、键盘可操作、语义化 HTML
```

## 视觉风格

**参考**: Claude × Linear × Stripe Dashboard

```
暗色默认:
  Background: #0a0a0f
  Surface:    #14141f
  Border:     #1e1e2e
  Brand:      #6366f1 (indigo-500)

亮色模式:
  Background: #ffffff
  Surface:    #f8f9fa
  Border:     #e5e7eb
  Brand:      #4f46e5 (indigo-600)
```

## 目录

```
i56-design-language/
├── tokens/
│   ├── i56-bdl-1.0.css     — 完整 CSS 自定义属性（556 个令牌）
│   └── i56-theme.js        — 主题切换引擎
├── icons/                   — SVG 图标库
├── specs/                   — 组件规范（MDX）
│   ├── button.md
│   ├── card.md
│   ├── table.md
│   ├── modal.md
│   ├── toast.md
│   └── ...
└── docs/
    ├── ARCHITECTURE.md      — 设计架构
    └── CONTRIBUTING.md      — 贡献指南
```

## 设计令牌（Token）

BDL 1.0 使用 CSS `@layer` 分层：

```
@layer i56-reset      — 重置样式
@layer i56-tokens     — 设计令牌定义
@layer i56-base       — 基础样式
@layer i56-utilities  — 工具类
@layer i56-components — 组件样式
@layer i56-pages      — 页面样式
```

### 品牌色变量

每个租户可以覆盖这些变量实现品牌定制：

```css
--i56-brand-50  ~ --i56-brand-950  (10 阶色阶)
--i56-brand       (主品牌色)
```

## 许可证

MIT
