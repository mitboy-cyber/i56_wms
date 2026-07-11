/**
 * I56 AI Workspace — Full-screen/half-screen AI console with multi-Agent collaboration.
 *
 * BDL 2.0 AI Web Component: native Custom Element, Shadow DOM, zero dependencies.
 *
 * Usage:
 *   <i56-ai-workspace mode="split" agents='[{"name":"Analyst","status":"idle"},{"name":"Coder","status":"active"}]'>
 *     <i56-ai-bar slot="ai-bar"></i56-ai-bar>
 *     <i56-ai-chat slot="ai-chat"></i56-ai-chat>
 *     <i56-ai-agent-monitor slot="ai-agent-monitor"></i56-ai-agent-monitor>
 *   </i56-ai-workspace>
 *
 * Attributes:
 *   mode: dock | split | fullscreen (default: dock)
 *   agents: JSON array of {name, status, icon?}
 *   expanded: present when open/expanded
 *
 * Slots:
 *   ai-bar — AI command bar
 *   ai-chat — chat panel
 *   ai-agent-monitor — agent runtime dashboard
 *   (default) — additional content below panels
 *
 * Events:
 *   i56:workspace-open  — detail: { mode }
 *   i56:workspace-close — detail: {}
 *   i56:mode-change     — detail: { mode, previous }
 *
 * Keyboard:
 *   Ctrl+Shift+K — toggle workspace
 *   Escape — close (when expanded)
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
          --i56-radius-xl: 16px;
          --i56-shadow-lg: 0 10px 15px -3px rgba(0,0,0,0.08), 0 4px 6px -4px rgba(0,0,0,0.06);
          --i56-shadow-xl: 0 20px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1);
          --i56-font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
          --i56-font-size-xs: 0.75rem;
          --i56-font-size-sm: 0.875rem;
          --i56-font-size-base: 1rem;
          --i56-font-size-lg: 1.125rem;
          --i56-font-size-xl: 1.25rem;
          --i56-transition: 150ms ease;
          --i56-transition-slow: 300ms ease;
          --i56-line-height: 1.5;
        }

        /* Dark theme */
        :host(.theme-dark), .theme-dark :host {
          --i56-color-bg: #0F172A;
          --i56-color-bg-surface: #1E293B;
          --i56-color-bg-secondary: #334155;
          --i56-color-border: #334155;
          --i56-color-border-hover: #475569;
          --i56-color-text: #F1F5F9;
          --i56-color-text-secondary: #94A3B8;
          --i56-color-text-tertiary: #64748B;
        }
      `);
    }
    return _sharedSheet;
  }

  function emit(el, name, detail = {}) {
    el.dispatchEvent(new CustomEvent(`i56:${name}`, { bubbles: true, composed: true, detail }));
  }

  class I56AiWorkspace extends HTMLElement {
    static get observedAttributes() {
      return ['mode', 'agents', 'expanded'];
    }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
      this._mode = 'dock';
      this._attached = false;
      this._boundKeydown = this._onGlobalKeydown.bind(this);
      this._boundClickOutside = this._onClickOutside.bind(this);
    }

    get mode() { return this._mode; }
    set mode(v) {
      const prev = this._mode;
      this._mode = v;
      this.setAttribute('mode', v);
      emit(this, 'mode-change', { mode: v, previous: prev });
    }

    get expanded() { return this.hasAttribute('expanded'); }

    get agents() {
      try { return JSON.parse(this.getAttribute('agents') || '[]'); } catch { return []; }
    }
    set agents(arr) {
      this.setAttribute('agents', JSON.stringify(arr));
    }

    connectedCallback() {
      if (this._attached) return;
      this._attached = true;
      this._mode = this.getAttribute('mode') || 'dock';
      document.addEventListener('keydown', this._boundKeydown);
      this.render();
    }

    disconnectedCallback() {
      document.removeEventListener('keydown', this._boundKeydown);
      document.removeEventListener('click', this._boundClickOutside);
      this._attached = false;
    }

    attributeChangedCallback(name, oldVal, newVal) {
      if (name === 'mode') {
        this._mode = newVal || 'dock';
        this.render();
      }
      if (name === 'agents') this.render();
    }

    // -- Public API --
    open(mode) {
      if (mode) this.mode = mode;
      this.setAttribute('expanded', '');
      this.render();
      emit(this, 'workspace-open', { mode: this._mode });
      document.addEventListener('click', this._boundClickOutside);
    }

    close() {
      this.removeAttribute('expanded');
      this.render();
      emit(this, 'workspace-close', {});
      document.removeEventListener('click', this._boundClickOutside);
    }

    toggle() {
      this.expanded ? this.close() : this.open();
    }

    // -- Internals --
    _onGlobalKeydown(e) {
      if ((e.ctrlKey || e.metaKey) && e.shiftKey && (e.key === 'K' || e.key === 'k')) {
        e.preventDefault();
        this.toggle();
      }
      if (e.key === 'Escape' && this.expanded) {
        this.close();
      }
    }

    _onClickOutside(e) {
      const el = this._root.querySelector('.ai-workspace-container');
      if (el && !el.contains(e.target) && !this.contains(e.target)) {
        if (this._mode === 'dock') return; // dock stays open
        this.close();
      }
    }

    _renderAgentHeader() {
      const agents = this.agents;
      if (!agents.length) return '';

      const badges = agents.slice(0, 4).map(a => {
        const statusClass = a.status === 'active' ? 'agent-active' :
                           a.status === 'busy' ? 'agent-busy' :
                           a.status === 'error' ? 'agent-error' : 'agent-idle';
        const icon = a.icon || (a.status === 'active' ? '●' : a.status === 'busy' ? '◉' : '○');
        return `<span class="agent-badge ${statusClass}" title="${a.name}: ${a.status}">
          <span class="agent-dot">${icon}</span>${a.name}
        </span>`;
      }).join('');

      const more = agents.length > 4 ? `<span class="agent-more">+${agents.length - 4}</span>` : '';

      return `<div class="workspace-agents">${badges}${more}</div>`;
    }

    render() {
      const mode = this._mode;
      const isExpanded = this.expanded;
      const isDock = mode === 'dock';
      const isSplit = mode === 'split';
      const isFull = mode === 'fullscreen';

      const containerClass = isFull ? 'fullscreen' : isSplit ? 'split' : 'dock';

      this._root.innerHTML = `
        <style>
          :host {
            display: contents;
            position: relative;
          }

          .ai-workspace-container {
            display: ${isExpanded ? 'flex' : 'none'};
            font-family: var(--i56-font-family);
            background: var(--i56-color-bg);
            border: 1px solid var(--i56-color-border);
            color: var(--i56-color-text);
            transition: all var(--i56-transition-slow) cubic-bezier(0.4, 0, 0.2, 1);
            overflow: hidden;
          }

          /* Dock mode: bottom panel */
          .ai-workspace-container.dock {
            position: fixed;
            bottom: 0;
            left: 0;
            right: 0;
            height: 45vh;
            max-height: 600px;
            z-index: 1000;
            border-radius: var(--i56-radius-xl) var(--i56-radius-xl) 0 0;
            box-shadow: var(--i56-shadow-xl);
            flex-direction: column;
            animation: i56-ws-slide-up 300ms cubic-bezier(0.4, 0, 0.2, 1);
          }

          @keyframes i56-ws-slide-up {
            from { transform: translateY(100%); }
            to { transform: translateY(0); }
          }

          @keyframes i56-ws-scale-in {
            from { transform: scale(0.95); opacity: 0; }
            to { transform: scale(1); opacity: 1; }
          }

          /* Split mode: right-side panel */
          .ai-workspace-container.split {
            position: fixed;
            top: 0;
            right: 0;
            bottom: 0;
            width: 40vw;
            min-width: 400px;
            max-width: 700px;
            z-index: 1000;
            flex-direction: column;
            border-radius: var(--i56-radius-xl) 0 0 var(--i56-radius-xl);
            box-shadow: var(--i56-shadow-xl);
            animation: i56-ws-slide-right 300ms cubic-bezier(0.4, 0, 0.2, 1);
          }

          @keyframes i56-ws-slide-right {
            from { transform: translateX(100%); }
            to { transform: translateX(0); }
          }

          /* Fullscreen mode */
          .ai-workspace-container.fullscreen {
            position: fixed;
            inset: 0;
            z-index: 1050;
            flex-direction: column;
            border-radius: 0;
            border: none;
            animation: i56-ws-scale-in 250ms cubic-bezier(0.4, 0, 0.2, 1);
          }

          /* Header */
          .workspace-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0.75rem 1.25rem;
            border-bottom: 1px solid var(--i56-color-border);
            background: var(--i56-color-bg-surface);
            flex-shrink: 0;
            user-select: none;
          }

          .workspace-title {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-weight: 700;
            font-size: var(--i56-font-size-base);
            color: var(--i56-color-text);
          }

          .workspace-title-icon {
            font-size: var(--i56-font-size-lg);
          }

          .workspace-agents {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            flex-wrap: wrap;
          }

          .agent-badge {
            display: inline-flex;
            align-items: center;
            gap: 0.25rem;
            padding: 0.125rem 0.5rem;
            border-radius: var(--i56-radius-full);
            font-size: var(--i56-font-size-xs);
            font-weight: 500;
            background: var(--i56-color-bg-secondary);
            color: var(--i56-color-text-secondary);
          }

          .agent-dot { font-size: 0.5rem; }
          .agent-active { color: #059669; background: #ECFDF5; }
          .agent-busy { color: #D97706; background: #FFFBEB; }
          .agent-error { color: #DC2626; background: #FEF2F2; }
          .agent-idle { color: var(--i56-color-text-tertiary); }

          .agent-more {
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
          }

          .workspace-actions {
            display: flex;
            align-items: center;
            gap: 0.375rem;
          }

          .ws-action-btn {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            width: 2rem;
            height: 2rem;
            border: none;
            border-radius: var(--i56-radius);
            background: transparent;
            color: var(--i56-color-text-secondary);
            cursor: pointer;
            font-size: var(--i56-font-size-base);
            transition: all var(--i56-transition);
          }
          .ws-action-btn:hover {
            background: var(--i56-color-bg-secondary);
            color: var(--i56-color-text);
          }
          .ws-action-btn:focus-visible {
            box-shadow: 0 0 0 2px var(--i56-color-brand-light);
            outline: none;
          }
          .ws-action-btn.active {
            background: var(--i56-color-brand-light);
            color: var(--i56-color-brand);
          }

          /* Body: main content area with slots */
          .workspace-body {
            flex: 1;
            display: flex;
            flex-direction: column;
            overflow: hidden;
            min-height: 0;
          }

          .slot-ai-bar {
            padding: 0.75rem 1.25rem;
            border-bottom: 1px solid var(--i56-color-border);
            flex-shrink: 0;
          }

          .slot-ai-chat {
            flex: 1;
            overflow: hidden;
            min-height: 0;
          }

          .slot-ai-agent-monitor {
            padding: 0.75rem 1.25rem;
            border-top: 1px solid var(--i56-color-border);
            flex-shrink: 0;
            max-height: 120px;
            overflow-y: auto;
          }

          .slot-default {
            padding: 1rem 1.25rem;
          }

          /* Empty slot placeholders */
          .slot-empty {
            color: var(--i56-color-text-tertiary);
            font-size: var(--i56-font-size-xs);
            text-align: center;
            padding: 1rem;
          }

          /* Backdrop for fullscreen/split */
          .workspace-backdrop {
            display: ${isFull || isSplit ? 'block' : 'none'};
            position: fixed;
            inset: 0;
            z-index: 1049;
            background: rgba(15, 23, 42, 0.3);
            backdrop-filter: blur(4px);
            animation: i56-fade-in 200ms ease;
          }

          @keyframes i56-fade-in {
            from { opacity: 0; }
            to { opacity: 1; }
          }

          /* Resize handle for dock mode */
          .resize-handle {
            display: ${isDock ? 'block' : 'none'};
            position: absolute;
            top: -8px;
            left: 0;
            right: 0;
            height: 12px;
            cursor: ns-resize;
            z-index: 2;
          }

          .resize-handle::after {
            content: '';
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 48px;
            height: 4px;
            border-radius: var(--i56-radius-full);
            background: var(--i56-color-border-hover);
            transition: all var(--i56-transition);
          }
          .resize-handle:hover::after {
            background: var(--i56-color-brand);
            width: 64px;
          }

          /* Mode indicator pills */
          .mode-indicator {
            display: flex;
            gap: 0.25rem;
            padding: 0.25rem;
            background: var(--i56-color-bg-secondary);
            border-radius: var(--i56-radius);
          }

          .mode-pill {
            padding: 0.25rem 0.625rem;
            border-radius: var(--i56-radius-sm);
            border: none;
            background: transparent;
            color: var(--i56-color-text-secondary);
            font-size: var(--i56-font-size-xs);
            font-weight: 500;
            cursor: pointer;
            transition: all var(--i56-transition);
            font-family: var(--i56-font-family);
          }
          .mode-pill:hover { color: var(--i56-color-text); }
          .mode-pill.active {
            background: var(--i56-color-brand);
            color: var(--i56-color-text-inverse);
          }
          .mode-pill:focus-visible {
            box-shadow: 0 0 0 2px var(--i56-color-brand-light);
            outline: none;
          }
        </style>

        ${isFull || isSplit ? '<div class="workspace-backdrop"></div>' : ''}

        <div class="ai-workspace-container ${containerClass}" role="dialog" aria-label="AI Workspace">
          ${isDock ? '<div class="resize-handle" aria-hidden="true"></div>' : ''}

          <div class="workspace-header">
            <div class="workspace-title">
              <span class="workspace-title-icon" aria-hidden="true">🤖</span>
              <span>AI Workspace</span>
              ${this._renderAgentHeader()}
            </div>

            <div class="workspace-actions">
              <div class="mode-indicator">
                <button class="mode-pill ${isDock ? 'active' : ''}" data-mode="dock" aria-pressed="${isDock}" title="Dock to bottom">
                  ⬇
                </button>
                <button class="mode-pill ${isSplit ? 'active' : ''}" data-mode="split" aria-pressed="${isSplit}" title="Split right">
                  ⬌
                </button>
                <button class="mode-pill ${isFull ? 'active' : ''}" data-mode="fullscreen" aria-pressed="${isFull}" title="Fullscreen">
                  ⛶
                </button>
              </div>
              <button class="ws-action-btn" title="Close workspace" aria-label="Close workspace">
                ✕
              </button>
            </div>
          </div>

          <div class="workspace-body">
            <div class="slot-ai-bar">
              <slot name="ai-bar"></slot>
            </div>
            <div class="slot-ai-chat">
              <slot name="ai-chat"></slot>
            </div>
            <div class="slot-ai-agent-monitor">
              <slot name="ai-agent-monitor"></slot>
            </div>
            <div class="slot-default">
              <slot></slot>
            </div>
          </div>
        </div>
      `;

      // Wire mode pills
      this._root.querySelectorAll('.mode-pill').forEach(pill => {
        pill.addEventListener('click', () => {
          this.mode = pill.dataset.mode;
          this.render();
          this.open(this._mode);
        });
      });

      // Wire close button
      const closeBtn = this._root.querySelector('.ws-action-btn[title="Close workspace"]');
      if (closeBtn) {
        closeBtn.addEventListener('click', () => this.close());
      }

      // Wire resize handle (drag to resize in dock mode)
      const handle = this._root.querySelector('.resize-handle');
      if (handle) {
        let startY, startHeight;
        const onMouseMove = (e) => {
          const delta = startY - e.clientY;
          const newHeight = Math.max(200, Math.min(800, startHeight + delta));
          const container = this._root.querySelector('.ai-workspace-container');
          if (container) {
            container.style.height = newHeight + 'px';
          }
        };
        const onMouseUp = () => {
          document.removeEventListener('mousemove', onMouseMove);
          document.removeEventListener('mouseup', onMouseUp);
        };
        handle.addEventListener('mousedown', (e) => {
          startY = e.clientY;
          const container = this._root.querySelector('.ai-workspace-container');
          startHeight = container ? container.offsetHeight : 400;
          document.addEventListener('mousemove', onMouseMove);
          document.addEventListener('mouseup', onMouseUp);
          e.preventDefault();
        });
      }

      // Wire backdrop click
      const backdrop = this._root.querySelector('.workspace-backdrop');
      if (backdrop) {
        backdrop.addEventListener('click', () => this.close());
      }
    }
  }

  customElements.define('i56-ai-workspace', I56AiWorkspace);
})();
