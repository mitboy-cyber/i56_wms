/**
 * I56 AI Context Menu — Right-click AI actions with dynamic menu items
 * based on selected entity context.
 *
 * BDL 2.0 AI Web Component: native Custom Element, Shadow DOM, zero dependencies.
 *
 * Usage:
 *   <i56-ai-context-menu id="ai-menu"
 *     data-context='{"entity":"order","id":"ORD-1234","type":"sales"}'>
 *   </i56-ai-context-menu>
 *
 *   // Show on right-click:
 *   document.addEventListener('contextmenu', e => {
 *     const menu = document.querySelector('#ai-menu');
 *     menu.context = { entity: 'order', id: 'ORD-1234' };
 *     menu.showAtEvent(e);
 *   });
 *
 * Attributes:
 *   data-context: JSON string with current entity context
 *   items: JSON array for static menu items (overrides auto-generated)
 *
 * Properties:
 *   context: get/set context object
 *   items: get/set menu items
 *
 * Events:
 *   i56:ai-action — detail: { action, context, entity }
 *   i56:close
 *   i56:open
 *
 * Keyboard:
 *   Arrow keys to navigate, Enter to select, Escape to close
 *   Submenus: ArrowRight to open, ArrowLeft to close
 */

(function () {
  'use strict';

  let _sharedSheet = null;
  function getSharedSheet() {
    if (!_sharedSheet) {
      _sharedSheet = new CSSStyleSheet();
      _sharedSheet.replaceSync(`
        :host {
          --i56-color-brand: var(--i56-brand, #1D4ED8);
          --i56-color-brand-hover: var(--i56-brand-hover, #1E40AF);
          --i56-color-brand-light: var(--i56-brand-light, #DBEAFE);
          --i56-color-danger: var(--i56-danger, #DC2626);
          --i56-color-danger-light: var(--i56-danger-light, #FEF2F2);
          --i56-color-bg: var(--i56-bg-base, #F8F9FB);
          --i56-color-bg-surface: var(--i56-bg-surface, #FFFFFF);
          --i56-color-bg-secondary: var(--i56-bg-secondary, #F1F5F9);
          --i56-color-border: var(--i56-border, #E2E8F0);
          --i56-color-border-hover: var(--i56-border-hover, #CBD5E1);
          --i56-color-text: var(--i56-text-primary, #111827);
          --i56-color-text-secondary: var(--i56-text-secondary, #64748B);
          --i56-color-text-tertiary: var(--i56-text-tertiary, #94A3B8);
          --i56-color-text-inverse: var(--i56-text-inverse, #FFFFFF);
          --i56-radius-sm: 4px;
          --i56-radius: 6px;
          --i56-radius-md: 8px;
          --i56-radius-lg: 12px;
          --i56-shadow-lg: 0 10px 15px -3px rgba(0,0,0,0.08), 0 4px 6px -4px rgba(0,0,0,0.06);
          --i56-shadow-xl: 0 20px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1);
          --i56-font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
          --i56-font-size-xs: 0.75rem;
          --i56-font-size-sm: 0.875rem;
          --i56-font-size-base: 1rem;
          --i56-transition: 150ms ease;
        }
      `);
    }
    return _sharedSheet;
  }

  function emit(el, name, detail = {}) {
    el.dispatchEvent(new CustomEvent(`i56:${name}`, { bubbles: true, composed: true, detail }));
  }

  /**
   * Auto-generate AI context menu items based on entity type.
   * Can be overridden by providing the `items` attribute.
   */
  function getDefaultItems(context) {
    const entity = (context.entity || '').toLowerCase();
    const items = [];

    if (!entity) return items;

    switch (entity) {
      case 'order':
        items.push(
          { label: '🤖 Analyze Order', action: 'ai:analyze-order', icon: '🔍', shortcut: 'Ctrl+Shift+A' },
          { label: '📋 Summarize', action: 'ai:summarize', icon: '📝' },
          { type: 'divider' },
          { label: '🔮 Predict Delivery', action: 'ai:predict-delivery', icon: '📦' },
          { label: '⚠️ Fraud Check', action: 'ai:fraud-check', icon: '🛡️', danger: true },
          { type: 'divider' },
          { label: '💬 Ask AI About...', action: 'ai:ask', icon: '✨', children: [
            { label: 'Shipping History', action: 'ai:ask-shipping' },
            { label: 'Payment Issues', action: 'ai:ask-payment' },
            { label: 'Customer Notes', action: 'ai:ask-notes' },
          ]},
        );
        break;

      case 'user':
      case 'customer':
        items.push(
          { label: '👤 Profile Summary', action: 'ai:profile-summary', icon: '📋' },
          { label: '📊 Activity Analysis', action: 'ai:activity-analysis', icon: '📈' },
          { type: 'divider' },
          { label: '🚨 Risk Assessment', action: 'ai:risk-assessment', icon: '⚠️', danger: true },
          { label: '💬 Ask AI About...', action: 'ai:ask', icon: '✨', children: [
            { label: 'Order History', action: 'ai:ask-orders' },
            { label: 'Support Tickets', action: 'ai:ask-tickets' },
            { label: 'Payment Methods', action: 'ai:ask-payment' },
          ]},
        );
        break;

      case 'product':
      case 'inventory':
        items.push(
          { label: '📦 Stock Analysis', action: 'ai:stock-analysis', icon: '📊' },
          { label: '💰 Price Optimization', action: 'ai:price-optimize', icon: '💡' },
          { type: 'divider' },
          { label: '🔍 Quality Review', action: 'ai:quality-review', icon: '✅' },
          { label: '💬 Ask AI About...', action: 'ai:ask', icon: '✨' },
        );
        break;

      case 'code':
      case 'file':
        items.push(
          { label: '🔍 Code Review', action: 'ai:code-review', icon: '👀', shortcut: 'Ctrl+Shift+R' },
          { label: '🐛 Find Bugs', action: 'ai:find-bugs', icon: '🔎' },
          { type: 'divider' },
          { label: '📝 Explain This', action: 'ai:explain', icon: '💡' },
          { label: '⚡ Optimize', action: 'ai:optimize', icon: '🚀' },
          { label: '🧪 Generate Tests', action: 'ai:generate-tests', icon: '🧪' },
          { type: 'divider' },
          { label: '💬 More AI Actions...', action: 'ai:ask', icon: '✨' },
        );
        break;

      case 'document':
        items.push(
          { label: '📝 Summarize', action: 'ai:summarize', icon: '📋', shortcut: 'Ctrl+Shift+S' },
          { label: '🌐 Translate', action: 'ai:translate', icon: '🔤' },
          { type: 'divider' },
          { label: '📊 Extract Data', action: 'ai:extract-data', icon: '📈' },
          { label: '💬 Ask AI About...', action: 'ai:ask', icon: '✨' },
        );
        break;

      default:
        items.push(
          { label: '✨ Ask AI', action: 'ai:ask', icon: '💬' },
          { label: '🔍 Analyze', action: 'ai:analyze', icon: '🔎' },
          { label: '📝 Summarize', action: 'ai:summarize', icon: '📋' },
        );
    }

    // Common items appended at the end
    items.push(
      { type: 'divider' },
      { label: '⚙️ AI Settings', action: 'ai:settings', icon: '🔧' },
    );

    return items;
  }

  class I56AiContextMenu extends HTMLElement {
    static get observedAttributes() { return ['data-context', 'items']; }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
      this._open = false;
      this._items = [];
      this._activeIdx = -1;
      this._submenuOpen = null;
      this._context = {};
      this._boundGlobalClick = this._onGlobalClick.bind(this);
      this._boundGlobalKey = this._onGlobalKey.bind(this);
      this._boundScroll = this._close.bind(this);
    }

    get context() {
      try { return JSON.parse(this.getAttribute('data-context') || '{}'); } catch { return this._context; }
    }
    set context(obj) {
      this._context = obj;
      this.setAttribute('data-context', JSON.stringify(obj));
      // Only regenerate items if no explicit items attribute
      if (!this.hasAttribute('items')) {
        this._items = getDefaultItems(obj);
      }
    }

    get items() { return [...this._items]; }
    set items(arr) {
      this._items = [...arr];
      if (this._open) this._renderItems();
    }

    connectedCallback() {
      this._context = this.context;
      if (this.hasAttribute('items')) {
        try { this._items = JSON.parse(this.getAttribute('items') || '[]'); } catch { this._items = []; }
      } else {
        this._items = getDefaultItems(this._context);
      }
    }

    disconnectedCallback() { this._close(); }

    attributeChangedCallback(name) {
      if (name === 'data-context') {
        this._context = this.context;
        if (!this.hasAttribute('items')) {
          this._items = getDefaultItems(this._context);
        }
      }
      if (name === 'items') {
        try { this._items = JSON.parse(this.getAttribute('items') || '[]'); } catch { this._items = []; }
      }
      if (this._open) {
        this.render();
        this._renderItems();
      }
    }

    // -- Public API --
    showAt(x, y) {
      this._close();
      this._open = true;
      this._activeIdx = -1;
      this._submenuOpen = null;
      this.setAttribute('data-open', '');

      document.addEventListener('click', this._boundGlobalClick, true);
      document.addEventListener('keydown', this._boundGlobalKey);
      window.addEventListener('scroll', this._boundScroll, true);

      this.render();
      this._position(x, y);
      emit(this, 'open', { context: this._context });
    }

    showAtEvent(e) {
      e.preventDefault();
      e.stopPropagation();
      this.showAt(e.clientX, e.clientY);
    }

    hide() { this._close(); }

    // -- Internals --
    _close() {
      if (!this._open) return;
      this._open = false;
      this._submenuOpen = null;
      this.removeAttribute('data-open');

      document.removeEventListener('click', this._boundGlobalClick, true);
      document.removeEventListener('keydown', this._boundGlobalKey);
      window.removeEventListener('scroll', this._boundScroll, true);

      this._root.querySelectorAll('.submenu-panel').forEach(s => s.remove());
      emit(this, 'close');
    }

    _onGlobalClick(e) {
      if (!this._root.contains(e.target)) {
        this._close();
      }
    }

    _onGlobalKey(e) {
      if (!this._open) return;

      const items = this._root.querySelectorAll('.ctx-item:not(.disabled):not(.divider)');
      const subItems = this._submenuOpen
        ? this._submenuOpen.querySelectorAll('.ctx-item:not(.disabled):not(.divider)')
        : [];

      switch (e.key) {
        case 'Escape':
          e.preventDefault();
          if (this._submenuOpen) { this._closeSubmenu(); } else { this._close(); }
          break;
        case 'ArrowDown':
          e.preventDefault();
          if (this._submenuOpen) { this._subNavigate(subItems, 1); }
          else { this._activeIdx = (this._activeIdx + 1) % items.length; this._highlightActive(items); }
          break;
        case 'ArrowUp':
          e.preventDefault();
          if (this._submenuOpen) { this._subNavigate(subItems, -1); }
          else { this._activeIdx = (this._activeIdx - 1 + items.length) % items.length; this._highlightActive(items); }
          break;
        case 'ArrowRight':
          e.preventDefault();
          if (!this._submenuOpen && this._activeIdx >= 0 && items[this._activeIdx]) {
            const sub = items[this._activeIdx].querySelector('.submenu-panel');
            if (sub) this._openSubmenu(items[this._activeIdx]);
          }
          break;
        case 'ArrowLeft':
          if (this._submenuOpen) { e.preventDefault(); this._closeSubmenu(); }
          break;
        case 'Enter':
        case ' ':
          e.preventDefault();
          if (this._submenuOpen && subItems.length > 0) {
            const subActive = this._submenuOpen.querySelector('.ctx-item.active');
            if (subActive) this._selectItem(subActive);
          } else if (this._activeIdx >= 0 && items.length > 0) {
            const active = items[this._activeIdx];
            if (!active.classList.contains('has-submenu')) this._selectItem(active);
          }
          break;
      }
    }

    _subNavigate(subItems, delta) {
      const currentIdx = Array.from(subItems).findIndex(el => el.classList.contains('active'));
      const nextIdx = ((currentIdx + delta) + subItems.length) % subItems.length;
      subItems.forEach(el => el.classList.remove('active'));
      if (subItems[nextIdx]) subItems[nextIdx].classList.add('active');
    }

    _highlightActive(items) {
      items.forEach(el => el.classList.remove('active'));
      if (items[this._activeIdx]) {
        items[this._activeIdx].classList.add('active');
        items[this._activeIdx].scrollIntoView({ block: 'nearest' });
      }
    }

    _selectItem(el) {
      const action = el.dataset.action;
      if (action) {
        emit(this, 'ai-action', { action, context: this._context, entity: this._context.entity });
        this._close();
      }
    }

    _openSubmenu(parentEl) {
      this._closeSubmenu();
      let sub = parentEl.querySelector('.submenu-panel');
      if (sub) {
        sub.style.display = 'block';
        this._submenuOpen = sub;
        const parentRect = parentEl.getBoundingClientRect();
        const menuRect = this._root.querySelector('.ctx-menu').getBoundingClientRect();
        sub.style.left = (parentRect.right - menuRect.left + 4) + 'px';
        sub.style.top = (parentRect.top - menuRect.top) + 'px';

        const subRect = sub.getBoundingClientRect();
        if (subRect.right > window.innerWidth - 8) {
          sub.style.left = (parentRect.left - menuRect.left - subRect.width - 4) + 'px';
        }
        if (subRect.bottom > window.innerHeight - 8) {
          sub.style.top = (parentRect.bottom - menuRect.top - subRect.height) + 'px';
        }

        const firstItem = sub.querySelector('.ctx-item:not(.disabled):not(.divider)');
        if (firstItem) {
          sub.querySelectorAll('.ctx-item').forEach(i => i.classList.remove('active'));
          firstItem.classList.add('active');
        }
      }
    }

    _closeSubmenu() {
      if (this._submenuOpen) {
        this._submenuOpen.style.display = 'none';
        this._submenuOpen.querySelectorAll('.ctx-item').forEach(i => i.classList.remove('active'));
        this._submenuOpen = null;
      }
    }

    _onItemClick(el) {
      if (el.classList.contains('disabled') || el.classList.contains('divider')) return;
      if (el.classList.contains('has-submenu')) {
        this._openSubmenu(el);
        return;
      }
      this._selectItem(el);
    }

    _onItemHover(el) {
      this._root.querySelectorAll('.ctx-item').forEach(i => i.classList.remove('active'));
      el.classList.add('active');
      if (this._submenuOpen) this._closeSubmenu();
      if (el.classList.contains('has-submenu')) this._openSubmenu(el);
    }

    _position(x, y) {
      const menu = this._root.querySelector('.ctx-menu');
      if (!menu) return;

      requestAnimationFrame(() => {
        const rect = menu.getBoundingClientRect();
        let left = x;
        let top = y;

        if (left + rect.width > window.innerWidth - 8) left = window.innerWidth - rect.width - 8;
        if (top + rect.height > window.innerHeight - 8) top = window.innerHeight - rect.height - 8;
        if (left < 8) left = 8;
        if (top < 8) top = 8;

        menu.style.left = left + 'px';
        menu.style.top = top + 'px';
        menu.style.display = 'block';
      });
    }

    _renderItem(item, idx) {
      if (item.type === 'divider') {
        return '<div class="ctx-divider" role="separator"></div>';
      }

      const hasChildren = item.children && item.children.length > 0;
      const dangerClass = item.danger ? ' danger' : '';
      const disabledClass = item.disabled ? ' disabled' : '';
      const subClass = hasChildren ? ' has-submenu' : '';

      return `
        <div class="ctx-item${dangerClass}${disabledClass}${subClass}"
             data-action="${item.action || ''}"
             data-idx="${idx}"
             role="menuitem"
             tabindex="-1"
             aria-disabled="${item.disabled ? 'true' : 'false'}">
          <span class="ctx-icon">${item.icon || ''}</span>
          <span class="ctx-label">${item.label || ''}</span>
          <span class="ctx-spacer"></span>
          ${item.shortcut ? '<span class="ctx-shortcut">' + item.shortcut + '</span>' : ''}
          ${hasChildren ? '<span class="ctx-arrow">▸</span>' : ''}
          ${hasChildren ? this._renderSubmenu(item.children) : ''}
        </div>
      `;
    }

    _renderSubmenu(children) {
      const items = children.map(child => {
        if (child.type === 'divider') return '<div class="ctx-divider" role="separator"></div>';
        const danger = child.danger ? ' danger' : '';
        const disabled = child.disabled ? ' disabled' : '';
        return `
          <div class="ctx-item${danger}${disabled}" data-action="${child.action || ''}" role="menuitem" tabindex="-1">
            <span class="ctx-icon">${child.icon || ''}</span>
            <span class="ctx-label">${child.label || ''}</span>
          </div>
        `;
      }).join('');
      return '<div class="submenu-panel" role="menu">' + items + '</div>';
    }

    _renderItems() {
      const list = this._root.querySelector('.ctx-list');
      if (!list) return;
      list.innerHTML = this._items.map((item, i) => this._renderItem(item, i)).join('');
      this._bindItemEvents();
    }

    _bindItemEvents() {
      this._root.querySelectorAll('.ctx-list > .ctx-item').forEach(el => {
        el.addEventListener('click', (e) => {
          e.stopPropagation();
          this._onItemClick(el);
        });
        el.addEventListener('mouseenter', () => this._onItemHover(el));
      });

      this._root.querySelectorAll('.submenu-panel .ctx-item').forEach(el => {
        el.addEventListener('click', (e) => {
          e.stopPropagation();
          if (!el.classList.contains('disabled')) this._selectItem(el);
        });
        el.addEventListener('mouseenter', () => {
          this._root.querySelectorAll('.submenu-panel .ctx-item').forEach(i => i.classList.remove('active'));
          el.classList.add('active');
        });
      });
    }

    render() {
      const entity = this._context.entity || '';
      const entityLabel = entity ? entity.charAt(0).toUpperCase() + entity.slice(1) : '';
      const entityId = this._context.id || '';

      this._root.innerHTML = `
        <style>
          :host {
            display: contents;
            position: relative;
          }

          .ctx-menu {
            display: ${this._open ? 'block' : 'none'};
            position: fixed;
            z-index: 9999;
            min-width: 14rem;
            max-width: 20rem;
            background: var(--i56-color-bg-surface);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-lg);
            box-shadow: var(--i56-shadow-xl);
            font-family: var(--i56-font-family);
            font-size: var(--i56-font-size-sm);
            overflow: visible;
            animation: i56-ctx-fade-in 120ms ease;
          }

          @keyframes i56-ctx-fade-in {
            from { opacity: 0; transform: scale(0.96); }
            to { opacity: 1; transform: scale(1); }
          }

          .ctx-header {
            padding: 0.5rem 0.75rem;
            border-bottom: 1px solid var(--i56-color-border);
            display: ${entity ? 'flex' : 'none'};
            align-items: center;
            gap: 0.375rem;
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-secondary);
            font-weight: 500;
          }
          .ctx-header-icon {
            font-size: var(--i56-font-size-base);
          }
          .ctx-header-id {
            font-family: 'SF Mono', 'Fira Code', monospace;
            color: var(--i56-color-brand);
            background: var(--i56-color-brand-light);
            padding: 0.125rem 0.375rem;
            border-radius: var(--i56-radius-sm);
          }

          .ctx-list {
            display: flex;
            flex-direction: column;
            padding: 0.25rem;
          }

          .ctx-item {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.375rem 0.5rem;
            border-radius: var(--i56-radius);
            cursor: pointer;
            color: var(--i56-color-text);
            transition: background 80ms ease, color 80ms ease;
            user-select: none;
            white-space: nowrap;
            position: relative;
          }
          .ctx-item:hover,
          .ctx-item.active {
            background: var(--i56-color-brand-light);
            color: var(--i56-color-brand);
          }
          .ctx-item.danger { color: var(--i56-color-danger); }
          .ctx-item.danger:hover,
          .ctx-item.danger.active {
            background: var(--i56-color-danger-light);
            color: var(--i56-color-danger);
          }
          .ctx-item.disabled {
            opacity: 0.4;
            cursor: not-allowed;
            pointer-events: none;
          }

          .ctx-divider {
            height: 1px;
            background: var(--i56-color-border);
            margin: 0.25rem 0.5rem;
          }

          .ctx-icon {
            flex-shrink: 0;
            width: 1.25rem;
            text-align: center;
            font-size: var(--i56-font-size-sm);
          }
          .ctx-label { flex: 1; min-width: 0; }
          .ctx-spacer { flex: 1; }
          .ctx-shortcut {
            flex-shrink: 0;
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
            font-family: monospace;
            letter-spacing: 0.03em;
          }
          .ctx-item.active .ctx-shortcut { color: inherit; opacity: 0.7; }
          .ctx-arrow {
            flex-shrink: 0;
            margin-left: 0.25rem;
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
          }

          .submenu-panel {
            display: none;
            position: absolute;
            top: -0.25rem;
            min-width: 12rem;
            background: var(--i56-color-bg-surface);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-md);
            box-shadow: var(--i56-shadow-lg);
            padding: 0.25rem;
            z-index: 1;
          }
        </style>

        <div class="ctx-menu" role="menu">
          ${entity ? '<div class="ctx-header"><span class="ctx-header-icon">🎯</span>' + entityLabel + (entityId ? ' <span class="ctx-header-id">' + entityId + '</span>' : '') + '</div>' : ''}
          <div class="ctx-list">
            ${this._items.map((item, i) => this._renderItem(item, i)).join('')}
          </div>
        </div>
      `;

      this._bindItemEvents();
    }
  }

  customElements.define('i56-ai-context-menu', I56AiContextMenu);
})();
