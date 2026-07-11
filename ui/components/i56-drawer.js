/**
 * I56 Drawer — Slide-in panel using native <dialog> with backdrop blur.
 *
 * Usage:
 *   <i56-drawer position="right" title="Settings">
 *     <div slot="body">Drawer content here</div>
 *   </i56-drawer>
 *
 *   drawer.show() / drawer.hide() / drawer.toggle()
 *
 * Attributes:
 *   position: left | right | top | bottom (default: right)
 *   open: open state (present = open)
 *   title: header title text
 *   size: sm | md | lg (default: md)
 *   close-button: "false" to hide close button
 *
 * Events:
 *   i56:open, i56:close
 *
 * Keyboard: Escape to close
 */

(function () {
  'use strict';

  // ---- shared design-token sheet (mirrors i56-components) ----
  let _sharedSheet = null;
  function getSharedSheet() {
    if (!_sharedSheet) {
      _sharedSheet = new CSSStyleSheet();
      _sharedSheet.replaceSync(`
        :host {
          --i56-color-brand: var(--i56-brand, #4F46E5);
          --i56-color-brand-hover: var(--i56-brand-hover, #4338CA);
          --i56-color-bg: var(--i56-bg, #FFFFFF);
          --i56-color-bg-secondary: var(--i56-bg-secondary, #F9FAFB);
          --i56-color-border: var(--i56-border, #E5E7EB);
          --i56-color-text: var(--i56-text, #111827);
          --i56-color-text-secondary: var(--i56-text-secondary, #6B7280);
          --i56-color-text-tertiary: var(--i56-text-tertiary, #9CA3AF);
          --i56-radius: 6px;
          --i56-radius-lg: 12px;
          --i56-shadow-xl: 0 20px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1);
          --i56-font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
          --i56-font-size-sm: 0.875rem;
          --i56-font-size-base: 1rem;
          --i56-font-size-lg: 1.125rem;
          --i56-font-size-xl: 1.25rem;
          --i56-transition: 150ms ease;
          --i56-transition-slow: 300ms ease;
        }
      `);
    }
    return _sharedSheet;
  }

  function emit(el, name, detail = {}) {
    el.dispatchEvent(new CustomEvent(`i56:${name}`, { bubbles: true, composed: true, detail }));
  }

  class I56Drawer extends HTMLElement {
    static get observedAttributes() {
      return ['position', 'open', 'title', 'size', 'close-button'];
    }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
    }

    get position() { return this.getAttribute('position') || 'right'; }
    get size() { return this.getAttribute('size') || 'md'; }

    connectedCallback() {
      this.render();
      if (this.hasAttribute('open')) this._open();
    }

    disconnectedCallback() {
      this._close(true);
    }

    attributeChangedCallback(name, oldVal, newVal) {
      if (name === 'open') {
        if (newVal !== null) this._open();
        else this._close();
      }
      if (this._dialog) this.render();
    }

    // -- public API --
    show() { this.setAttribute('open', ''); }
    hide() { this.removeAttribute('open'); }
    toggle() {
      this.hasAttribute('open') ? this.hide() : this.show();
    }

    // -- internals --
    _open() {
      if (!this._dialog) { this.render(); }
      const dialog = this._dialog;
      if (!dialog) return;

      document.body.style.overflow = 'hidden';

      // Use native <dialog> showModal() for backdrop + focus trap
      if (typeof dialog.showModal === 'function' && !dialog.open) {
        dialog.showModal();
      }
      dialog.setAttribute('data-open', '');

      this._bindEsc();
      // Animate in
      requestAnimationFrame(() => {
        const panel = this._root.querySelector('.drawer-panel');
        if (panel) panel.classList.add('open');
      });

      emit(this, 'open');
    }

    _close(disconnected = false) {
      if (disconnected) document.body.style.overflow = '';
      const dialog = this._dialog;
      if (!dialog) return;

      const panel = this._root.querySelector('.drawer-panel');
      if (panel) panel.classList.remove('open');

      // Close native dialog after animation
      setTimeout(() => {
        if (dialog.open && typeof dialog.close === 'function') dialog.close();
        dialog.removeAttribute('data-open');
        if (!disconnected) document.body.style.overflow = '';
        emit(this, 'close');
      }, 250);
    }

    _onOverlayClick = (e) => {
      // Native <dialog> fires click on the backdrop itself
      const rect = this._dialog?.getBoundingClientRect();
      if (!rect) return;
      if (
        e.clientX < rect.left || e.clientX > rect.right ||
        e.clientY < rect.top || e.clientY > rect.bottom
      ) {
        this.hide();
      }
    };

    _onKeydown = (e) => {
      if (e.key === 'Escape') {
        this.hide();
      }
    };

    _bindEsc() {
      document.addEventListener('keydown', this._onKeydown);
    }

    _unbindEsc() {
      document.removeEventListener('keydown', this._onKeydown);
    }

    render() {
      const position = this.position;
      const title = this.getAttribute('title') || '';
      const size = this.size;
      const closeButton = this.getAttribute('close-button') !== 'false';
      const isOpen = this.hasAttribute('open');

      const sizeMap = { sm: '20rem', md: '28rem', lg: '36rem' };
      const dimension = sizeMap[size] || sizeMap.md;

      // Position-specific slide transforms & dimensions
      const isVertical = position === 'top' || position === 'bottom';
      const widthOrHeight = isVertical ? 'height' : 'width';

      const slideFromMap = {
        left: 'translateX(-100%)',
        right: 'translateX(100%)',
        top: 'translateY(-100%)',
        bottom: 'translateY(100%)',
      };
      const slideFrom = slideFromMap[position] || slideFromMap.right;

      const dimensionStyle = isVertical
        ? `height: ${dimension}; width: 100%;`
        : `width: ${dimension}; height: 100%;`;

      const edgeStyleMap = {
        left: 'top: 0; left: 0; bottom: 0;',
        right: 'top: 0; right: 0; bottom: 0;',
        top: 'top: 0; left: 0; right: 0;',
        bottom: 'bottom: 0; left: 0; right: 0;',
      };
      const edgeStyle = edgeStyleMap[position] || edgeStyleMap.right;

      this._root.innerHTML = `
        <style>
          :host { display: contents; }

          dialog {
            position: fixed;
            inset: 0;
            border: none;
            padding: 0;
            margin: 0;
            background: transparent;
            max-width: none;
            max-height: none;
            width: 100%;
            height: 100%;
            overflow: hidden;
            z-index: 1050;
            opacity: 0;
            pointer-events: none;
          }
          dialog[data-open] {
            opacity: 1;
            pointer-events: auto;
          }
          dialog::backdrop {
            background: rgba(0, 0, 0, 0.4);
            backdrop-filter: blur(4px);
            animation: i56-drawer-fade-in 200ms ease;
          }

          .drawer-panel {
            position: absolute;
            ${edgeStyle}
            ${dimensionStyle}
            background: var(--i56-color-bg);
            box-shadow: var(--i56-shadow-xl);
            display: flex;
            flex-direction: column;
            font-family: var(--i56-font-family);
            transform: ${slideFrom};
            transition: transform var(--i56-transition-slow) cubic-bezier(0.4, 0, 0.2, 1);
          }
          .drawer-panel.open {
            transform: translate(0, 0);
          }

          .drawer-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 1rem 1.25rem;
            border-bottom: 1px solid var(--i56-color-border);
            flex-shrink: 0;
          }
          .drawer-title {
            font-size: var(--i56-font-size-lg);
            font-weight: 600;
            color: var(--i56-color-text);
            margin: 0;
          }
          .close-btn {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 2rem;
            height: 2rem;
            border: none;
            background: none;
            border-radius: var(--i56-radius);
            cursor: pointer;
            color: var(--i56-color-text-tertiary);
            font-size: var(--i56-font-size-xl);
            line-height: 1;
            transition: all var(--i56-transition);
            flex-shrink: 0;
          }
          .close-btn:hover { background: var(--i56-color-bg-secondary); color: var(--i56-color-text); }
          .close-btn:focus-visible { box-shadow: 0 0 0 2px var(--i56-color-brand); outline: none; }

          .drawer-body {
            flex: 1;
            overflow-y: auto;
            padding: 1.25rem;
            -webkit-overflow-scrolling: touch;
          }

          .drawer-footer {
            padding: 1rem 1.25rem;
            border-top: 1px solid var(--i56-color-border);
            flex-shrink: 0;
          }
          .drawer-footer:empty { display: none; }

          @keyframes i56-drawer-fade-in {
            from { opacity: 0; }
            to { opacity: 1; }
          }
        </style>

        <dialog>
          <div class="drawer-panel${isOpen ? ' open' : ''}">
            ${title || closeButton ? `
              <div class="drawer-header">
                <span class="drawer-title">${title}</span>
                ${closeButton ? `<button class="close-btn" aria-label="Close drawer">&times;</button>` : ''}
              </div>
            ` : ''}
            <div class="drawer-body">
              <slot name="body"><slot></slot></slot>
            </div>
            <div class="drawer-footer">
              <slot name="footer"></slot>
            </div>
          </div>
        </dialog>
      `;

      this._dialog = this._root.querySelector('dialog');

      // Wire close button
      const closeBtn = this._root.querySelector('.close-btn');
      if (closeBtn) {
        closeBtn.addEventListener('click', () => this.hide());
      }

      // Wire overlay click (native dialog backdrop click)
      if (this._dialog) {
        this._dialog.addEventListener('click', this._onOverlayClick);
        this._dialog.addEventListener('close', () => {
          if (this.hasAttribute('open')) this.hide();
        });
      }

      if (isOpen) this._open();
    }
  }

  customElements.define('i56-drawer', I56Drawer);
})();
