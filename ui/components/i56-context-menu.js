/**
 * I56 Context Menu — Right-click context menu with keyboard support.
 *
 * Usage:
 *   <i56-context-menu id="my-menu" items='[
 *     {"label":"Edit","icon":"✏️","shortcut":"Ctrl+E","action":"edit"},
 *     {"label":"Copy","icon":"📋","shortcut":"Ctrl+C","action":"copy"},
 *     {"type":"divider"},
 *     {"label":"Delete","icon":"🗑️","danger":true,"action":"delete"}
 *   ]'></i56-context-menu>
 *
 *   // Programmatic open:
 *   const menu = document.querySelector('#my-menu');
 *   menu.showAt(x, y);
 *   menu.showAtEvent(event);  // from right-click or keyboard
 *
 *   // Listen for actions:
 *   menu.addEventListener('i56:action', e => console.log(e.detail.action));
 *
 * Attributes:
 *   items: JSON array for static items
 *
 * Properties:
 *   items: get/set array of menu items
 *
 * Menu item shape:
 *   {type?: 'divider', label?: string, icon?: string, shortcut?: string,
 *    danger?: boolean, disabled?: boolean, action?: string, children?: items[]}
 *
 * Events:
 *   i56:action — detail: { action, item }
 *   i56:close
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
          --i56-color-brand: var(--i56-brand, #4F46E5);
          --i56-color-brand-hover: var(--i56-brand-hover, #4338CA);
          --i56-color-brand-light: var(--i56-brand-light, #EEF2FF);
          --i56-color-danger: var(--i56-danger, #DC2626);
          --i56-color-danger-light: var(--i56-danger-light, #FEF2F2);
          --i56-color-bg: var(--i56-bg, #FFFFFF);
          --i56-color-bg-secondary: var(--i56-bg-secondary, #F9FAFB);
          --i56-color-border: var(--i56-border, #E5E7EB);
          --i56-color-text: var(--i56-text, #111827);
          --i56-color-text-secondary: var(--i56-text-secondary, #6B7280);
          --i56-color-text-tertiary: var(--i56-text-tertiary, #9CA3AF);
          --i56-radius: 6px;
          --i56-radius-md: 8px;
          --i56-shadow-lg: 0 10px 15px -3px rgba(0,0,0,0.1), 0 4px 6px -4px rgba(0,0,0,0.1);
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

  class I56ContextMenu extends HTMLElement {
    static get observedAttributes() { return ['items']; }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
      this._open = false;
      this._items = [];
      this._activeIdx = -1;
      this._submenuOpen = null;
      this._boundGlobalClick = this._onGlobalClick.bind(this);
      this._boundGlobalKey = this._onGlobalKey.bind(this);
      this._boundScroll = this._close.bind(this);
    }

    get items() { return this._items; }
    set items(arr) {
      this._items = arr;
      if (this._open) this._renderItems();
    }

    connectedCallback() {
      this._items = this._parseItems();
    }

    disconnectedCallback() {
      this._close();
    }

    attributeChangedCallback(name) {
      if (name === 'items') {
        this._items = this._parseItems();
        if (this._open) this._renderItems();
      }
    }

    _parseItems() {
      try { return JSON.parse(this.getAttribute('items') || '[]'); } catch { return []; }
    }

    // -- Public API --
    showAt(x, y) {
      this._close(); // close any existing instance
      this._open = true;
      this._activeIdx = -1;
      this._submenuOpen = null;
      this.setAttribute('data-open', '');

      document.addEventListener('click', this._boundGlobalClick, true);
      document.addEventListener('keydown', this._boundGlobalKey);
      window.addEventListener('scroll', this._boundScroll, true);

      this.render();
      this._position(x, y);
      emit(this, 'open');
    }

    showAtEvent(e) {
      e.preventDefault();
      this.showAt(e.clientX, e.clientY);
    }

    hide() {
      this._close();
    }

    // -- Internals --
    _close() {
      if (!this._open) return;
      this._open = false;
      this._submenuOpen = null;
      this.removeAttribute('data-open');

      document.removeEventListener('click', this._boundGlobalClick, true);
      document.removeEventListener('keydown', this._boundGlobalKey);
      window.removeEventListener('scroll', this._boundScroll, true);

      // Remove submenus
      this._root.querySelectorAll('.submenu-panel').forEach(s => s.remove());

      emit(this, 'close');
    }

    _onGlobalClick(e) {
      // Close if clicking outside the menu
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
          if (this._submenuOpen) {
            this._closeSubmenu();
          } else {
            this._close();
          }
          break;

        case 'ArrowDown':
          e.preventDefault();
          if (this._submenuOpen) {
            this._subNavigate(subItems, 1);
          } else {
            this._activeIdx = (this._activeIdx + 1) % items.length;
            this._highlightActive(items);
          }
          break;

        case 'ArrowUp':
          e.preventDefault();
          if (this._submenuOpen) {
            this._subNavigate(subItems, -1);
          } else {
            this._activeIdx = (this._activeIdx - 1 + items.length) % items.length;
            this._highlightActive(items);
          }
          break;

        case 'ArrowRight':
          e.preventDefault();
          if (!this._submenuOpen && this._activeIdx >= 0) {
            const item = items[this._activeIdx];
            const sub = item?.querySelector('.submenu-panel');
            if (sub) this._openSubmenu(item);
          }
          break;

        case 'ArrowLeft':
          if (this._submenuOpen) {
            e.preventDefault();
            this._closeSubmenu();
          }
          break;

        case 'Enter':
        case ' ':
          e.preventDefault();
          if (this._submenuOpen && subItems.length > 0) {
            const subActive = this._submenuOpen.querySelector('.ctx-item.active');
            if (subActive) this._selectItem(subActive);
          } else if (this._activeIdx >= 0 && items.length > 0) {
            const active = items[this._activeIdx];
            if (!active.classList.contains('has-submenu')) {
              this._selectItem(active);
            }
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
      if (items[this._activeIdx]) items[this._activeIdx].classList.add('active');
      // Scroll into view
      if (items[this._activeIdx]) {
        items[this._activeIdx].scrollIntoView({ block: 'nearest' });
      }
    }

    _selectItem(el) {
      const action = el.dataset.action;
      if (action) {
        const item = this._items.find(i => String(i.action) === action);
        emit(this, 'action', { action, item: item || {} });
        this._close();
      }
    }

    _openSubmenu(parentEl) {
      this._closeSubmenu();
      let sub = parentEl.querySelector('.submenu-panel');
      if (sub) {
        sub.style.display = 'block';
        this._submenuOpen = sub;
        // Position
        const parentRect = parentEl.getBoundingClientRect();
        const menuRect = this._root.querySelector('.ctx-menu').getBoundingClientRect();
        sub.style.left = (parentRect.right - menuRect.left + 4) + 'px';
        sub.style.top = (parentRect.top - menuRect.top) + 'px';

        // Keep in viewport
        const subRect = sub.getBoundingClientRect();
        if (subRect.right > window.innerWidth - 8) {
          sub.style.left = (parentRect.left - menuRect.left - subRect.width - 4) + 'px';
        }
        if (subRect.bottom > window.innerHeight - 8) {
          sub.style.top = (parentRect.bottom - menuRect.top - subRect.height) + 'px';
        }

        // Highlight first
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
      if (el.classList.contains('has-submenu')) {
        this._openSubmenu(el);
      }
    }

    _position(x, y) {
      const menu = this._root.querySelector('.ctx-menu');
      if (!menu) return;

      // Use requestAnimationFrame to ensure DOM is laid out
      requestAnimationFrame(() => {
        const rect = menu.getBoundingClientRect();
        let left = x;
        let top = y;

        // Keep in viewport
        if (left + rect.width > window.innerWidth - 8) {
          left = window.innerWidth - rect.width - 8;
        }
        if (top + rect.height > window.innerHeight - 8) {
          top = window.innerHeight - rect.height - 8;
        }
        if (left < 8) left = 8;
        if (top < 8) top = 8;

        menu.style.left = left + 'px';
        menu.style.top = top + 'px';
        menu.style.display = 'block';
        menu.style.animation = 'i56-ctx-fade-in 120ms ease';
      });
    }

    _renderItem(item, idx) {
      if (item.type === 'divider') {
        return `<div class="ctx-divider" role="separator"></div>`;
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
          ${item.shortcut ? `<span class="ctx-shortcut">${item.shortcut}</span>` : ''}
          ${hasChildren ? '<span class="ctx-arrow">▸</span>' : ''}
          ${hasChildren ? this._renderSubmenu(item.children) : ''}
        </div>
      `;
    }

    _renderSubmenu(children) {
      const items = children.map((child, i) => {
        if (child.type === 'divider') return `<div class="ctx-divider" role="separator"></div>`;
        const danger = child.danger ? ' danger' : '';
        const disabled = child.disabled ? ' disabled' : '';
        return `
          <div class="ctx-item${danger}${disabled}" data-action="${child.action || ''}" role="menuitem" tabindex="-1">
            <span class="ctx-icon">${child.icon || ''}</span>
            <span class="ctx-label">${child.label || ''}</span>
            <span class="ctx-spacer"></span>
            ${child.shortcut ? `<span class="ctx-shortcut">${child.shortcut}</span>` : ''}
          </div>
        `;
      }).join('');

      return `<div class="submenu-panel" role="menu">${items}</div>`;
    }

    _renderItems() {
      const list = this._root.querySelector('.ctx-list');
      if (!list) return;
      list.innerHTML = this._items.map((item, i) => this._renderItem(item, i)).join('');
      this._bindItemEvents();
    }

    _bindItemEvents() {
      this._root.querySelectorAll('.ctx-item').forEach(el => {
        el.addEventListener('click', (e) => {
          e.stopPropagation();
          this._onItemClick(el);
        });
        el.addEventListener('mouseenter', () => this._onItemHover(el));
      });

      // Submenu items
      this._root.querySelectorAll('.submenu-panel .ctx-item').forEach(el => {
        el.addEventListener('click', (e) => {
          e.stopPropagation();
          if (!el.classList.contains('disabled')) {
            this._selectItem(el);
          }
        });
        el.addEventListener('mouseenter', () => {
          this._root.querySelectorAll('.submenu-panel .ctx-item').forEach(i => i.classList.remove('active'));
          el.classList.add('active');
        });
      });
    }

    render() {
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
            min-width: 12rem;
            max-width: 18rem;
            background: var(--i56-color-bg);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-md);
            box-shadow: var(--i56-shadow-lg);
            padding: 0.25rem;
            font-family: var(--i56-font-family);
            font-size: var(--i56-font-size-sm);
            overflow: visible;
          }

          .ctx-list {
            display: flex;
            flex-direction: column;
          }

          .ctx-item {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.375rem 0.5rem;
            border-radius: var(--i56-radius);
            cursor: pointer;
            color: var(--i56-color-text);
            transition: background 80ms ease;
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

          .ctx-divider {
            height: 1px;
            background: var(--i56-color-border);
            margin: 0.25rem 0.5rem;
          }

          /* Submenu */
          .submenu-panel {
            display: none;
            position: absolute;
            z-index: 10000;
            min-width: 10rem;
            background: var(--i56-color-bg);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-md);
            box-shadow: var(--i56-shadow-lg);
            padding: 0.25rem;
            font-size: var(--i56-font-size-sm);
          }

          @keyframes i56-ctx-fade-in {
            from {
              opacity: 0;
              transform: scale(0.95);
            }
            to {
              opacity: 1;
              transform: scale(1);
            }
          }
        </style>
        <div class="ctx-menu" role="menu" aria-orientation="vertical">
          <div class="ctx-list">
            ${this._items.map((item, i) => this._renderItem(item, i)).join('')}
          </div>
        </div>
      `;

      this._bindItemEvents();

      if (this._open) {
        // Auto-focus for keyboard nav
        const firstItem = this._root.querySelector('.ctx-item:not(.disabled):not(.divider)');
        if (firstItem) firstItem.focus();
      }
    }
  }

  customElements.define('i56-context-menu', I56ContextMenu);
})();
