/**
 * I56 AI Bar — AI Command Bar (Ctrl+K / Ctrl+J shortcut).
 *
 * BDL 2.0 AI Web Component: native Custom Element, Shadow DOM, zero dependencies.
 *
 * A floating command bar at the bottom of the page for natural-language AI input.
 * Features: command completion, natural language input, tool hints, backdrop blur.
 *
 * Usage:
 *   <i56-ai-bar placeholder="Ask anything…" model="sonnet" tools-enabled
 *              suggestions='["Summarize","Fix bug","Add feature","Explain code"]'>
 *   </i56-ai-bar>
 *
 * Attributes:
 *   placeholder: input placeholder (default: "Ask AI anything…")
 *   model: AI model identifier
 *   tools-enabled: present if tool use is enabled
 *   suggestions: JSON array of suggestion labels (or {label, icon} objects)
 *   loading: present during AI processing
 *   expanded: present when the bar is open/visible
 *
 * Properties:
 *   value: get/set current input text
 *   loading: get/set loading state
 *
 * Events:
 *   i56:command — detail: { prompt, model?, toolsEnabled? }
 *   i56:submit — detail: { prompt }
 *   i56:suggestion-click — detail: { suggestion, index }
 *   i56:close
 *   i56:open
 *
 * Keyboard:
 *   Ctrl+K / ⌘K — toggle the bar
 *   Ctrl+J / ⌘J — alternate toggle
 *   Enter — submit
 *   Shift+Enter — newline
 *   Escape — close
 *   Tab — cycle through tool hints
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
          --i56-radius-full: 9999px;
          --i56-shadow-lg: 0 10px 15px -3px rgba(0,0,0,0.08), 0 4px 6px -4px rgba(0,0,0,0.06);
          --i56-shadow-xl: 0 20px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1);
          --i56-font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
          --i56-font-size-xs: 0.75rem;
          --i56-font-size-sm: 0.875rem;
          --i56-font-size-base: 1rem;
          --i56-font-size-lg: 1.125rem;
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

  class I56AiBar extends HTMLElement {
    static get observedAttributes() {
      return ['placeholder', 'model', 'tools-enabled', 'loading', 'expanded', 'suggestions'];
    }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
      this._attached = false;
      this._toolHintIdx = -1;
      this._boundGlobalKey = this._onGlobalKeydown.bind(this);
      this._boundEsc = this._onEscKey.bind(this);
    }

    get value() {
      const textarea = this._root.querySelector('.ai-input');
      return textarea ? textarea.value : '';
    }
    set value(v) {
      const textarea = this._root.querySelector('.ai-input');
      if (textarea) { textarea.value = v; this._resize(textarea); }
    }

    get loading() { return this.hasAttribute('loading'); }
    set loading(v) {
      if (v) this.setAttribute('loading', '');
      else this.removeAttribute('loading');
    }

    get expanded() { return this.hasAttribute('expanded'); }
    get toolsEnabled() { return this.hasAttribute('tools-enabled'); }
    get model() { return this.getAttribute('model') || ''; }

    connectedCallback() {
      if (this._attached) return;
      this._attached = true;
      document.addEventListener('keydown', this._boundGlobalKey);
      this.render();
    }

    disconnectedCallback() {
      document.removeEventListener('keydown', this._boundGlobalKey);
      document.removeEventListener('keydown', this._boundEsc);
      this._attached = false;
    }

    attributeChangedCallback(name, oldVal, newVal) {
      if (name === 'expanded') {
        if (newVal !== null) document.addEventListener('keydown', this._boundEsc);
        else document.removeEventListener('keydown', this._boundEsc);
      }
      this.render();
    }

    // -- Public API --
    open() {
      if (this.expanded) return;
      this.setAttribute('expanded', '');
      this.focus();
      emit(this, 'open');
    }

    close() {
      if (!this.expanded) return;
      this.removeAttribute('expanded');
      emit(this, 'close');
    }

    toggle() {
      this.expanded ? this.close() : this.open();
    }

    focus() {
      requestAnimationFrame(() => {
        const textarea = this._root.querySelector('.ai-input');
        if (textarea) textarea.focus();
      });
    }

    // -- Internals --
    _onGlobalKeydown(e) {
      // Ctrl+K or Ctrl+J toggle
      if ((e.ctrlKey || e.metaKey) && (e.key === 'k' || e.key === 'K' || e.key === 'j' || e.key === 'J')) {
        // Don't steal from focused inputs that may use Ctrl+K (e.g. code editors)
        if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA' || e.target.isContentEditable) return;
        e.preventDefault();
        this.toggle();
      }
    }

    _onEscKey(e) {
      if (e.key === 'Escape') {
        this.close();
      }
    }

    _onSubmit() {
      const textarea = this._root.querySelector('.ai-input');
      if (!textarea) return;
      const prompt = textarea.value.trim();
      if (!prompt || this.loading) return;

      textarea.value = '';
      this._resize(textarea);

      emit(this, 'command', {
        prompt,
        model: this.model || undefined,
        toolsEnabled: this.toolsEnabled,
      });
      emit(this, 'submit', { prompt });
    }

    _onKeydown = (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        this._onSubmit();
      } else if (e.key === 'Escape') {
        this.close();
      } else if (e.key === 'Tab') {
        e.preventDefault();
        this._cycleToolHint();
      }
    };

    _onSuggestionClick(suggestion, index) {
      const textarea = this._root.querySelector('.ai-input');
      if (textarea) {
        textarea.value = suggestion;
        this._resize(textarea);
        textarea.focus();
      }
      emit(this, 'suggestion-click', { suggestion, index });
    }

    _resize(textarea) {
      textarea.style.height = 'auto';
      textarea.style.height = Math.min(textarea.scrollHeight, 200) + 'px';
    }

    _cycleToolHint() {
      const hints = this._root.querySelectorAll('.tool-hint');
      if (!hints.length) return;
      this._toolHintIdx = (this._toolHintIdx + 1) % hints.length;
      hints.forEach((h, i) => {
        h.classList.toggle('active', i === this._toolHintIdx);
      });
    }

    _renderSuggestions() {
      let suggestions = [];
      try { suggestions = JSON.parse(this.getAttribute('suggestions') || '[]'); } catch { /* keep empty */ }
      if (!suggestions.length) return '';

      const chips = suggestions.map((s, i) => {
        const label = typeof s === 'object' ? s.label : s;
        const icon = typeof s === 'object' ? s.icon || '' : '';
        return `<button class="suggestion-chip" data-idx="${i}" type="button" ?disabled="${this.loading}">
          ${icon ? '<span class="chip-icon">' + icon + '</span>' : ''}
          <span>${label}</span>
        </button>`;
      }).join('');
      return `<div class="suggestions-row">${chips}</div>`;
    }

    _renderToolHints() {
      if (!this.toolsEnabled) return '';
      const tools = ['read_file', 'search_code', 'edit_file', 'run_command', 'web_fetch'];
      const hints = tools.map(t =>
        `<span class="tool-hint">/tool ${t}</span>`
      ).join('');
      return `<div class="tool-hints">${hints}</div>`;
    }

    render() {
      const placeholder = this.getAttribute('placeholder') || 'Ask AI anything…';
      const isExpanded = this.expanded;
      const isLoading = this.loading;
      const isTools = this.toolsEnabled;
      const model = this.model;

      this._root.innerHTML = `
        <style>
          :host {
            display: contents;
          }

          .ai-bar-overlay {
            display: ${isExpanded ? 'flex' : 'none'};
            position: fixed;
            bottom: 0;
            left: 50%;
            transform: translateX(-50%);
            z-index: 1100;
            width: 100%;
            max-width: 52rem;
            padding: 1rem;
            flex-direction: column;
            align-items: stretch;
            pointer-events: auto;
            animation: i56-ai-bar-in var(--i56-transition-slow) cubic-bezier(0.4, 0, 0.2, 1);
          }

          @keyframes i56-ai-bar-in {
            from { transform: translateX(-50%) translateY(1rem); opacity: 0; }
            to { transform: translateX(-50%) translateY(0); opacity: 1; }
          }

          .ai-bar-container {
            background: var(--i56-color-bg-surface);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-xl);
            box-shadow: var(--i56-shadow-xl);
            overflow: hidden;
            backdrop-filter: blur(12px);
            -webkit-backdrop-filter: blur(12px);
          }

          .input-row {
            display: flex;
            align-items: flex-start;
            gap: 0.625rem;
            padding: 0.875rem 1.125rem;
          }

          .ai-icon {
            flex-shrink: 0;
            font-size: var(--i56-font-size-lg);
            display: flex;
            align-items: center;
            padding-top: 0.25rem;
            transition: all var(--i56-transition);
          }
          ${isLoading ? '.ai-icon { animation: i56-ai-pulse 1.5s ease-in-out infinite; }' : ''}
          @keyframes i56-ai-pulse {
            0%, 100% { opacity: 1; transform: scale(1); }
            50% { opacity: 0.6; transform: scale(0.9); }
          }

          .ai-input {
            flex: 1;
            border: none;
            outline: none;
            font-size: var(--i56-font-size-base);
            font-family: var(--i56-font-family);
            color: var(--i56-color-text);
            background: transparent;
            resize: none;
            min-height: 1.75rem;
            max-height: 200px;
            line-height: var(--i56-line-height);
            padding: 0;
          }
          .ai-input::placeholder { color: var(--i56-color-text-tertiary); font-size: var(--i56-font-size-sm); }
          .ai-input:disabled { opacity: 0.5; }

          .submit-btn {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 2.25rem;
            height: 2.25rem;
            border: none;
            border-radius: var(--i56-radius);
            background: ${isLoading ? 'var(--i56-color-bg-secondary)' : 'var(--i56-color-brand)'};
            color: var(--i56-color-text-inverse);
            cursor: ${isLoading ? 'wait' : 'pointer'};
            font-size: var(--i56-font-size-base);
            flex-shrink: 0;
            transition: all var(--i56-transition);
          }
          .submit-btn:hover:not(:disabled) { background: var(--i56-color-brand-hover); transform: translateY(-1px); }
          .submit-btn:disabled { opacity: 0.5; cursor: not-allowed; }
          .submit-btn:focus-visible { box-shadow: 0 0 0 3px var(--i56-color-brand-light); outline: none; }

          .submit-btn .spinner {
            width: 0.875rem;
            height: 0.875rem;
            border: 2px solid var(--i56-color-text-inverse);
            border-right-color: transparent;
            border-radius: 50%;
            animation: i56-spin 0.6s linear infinite;
          }
          @keyframes i56-spin {
            to { transform: rotate(360deg); }
          }

          /* Suggestions */
          .suggestions-row {
            display: flex;
            flex-wrap: wrap;
            gap: 0.375rem;
            padding: 0 1.125rem 0.75rem;
          }
          .suggestion-chip {
            display: inline-flex;
            align-items: center;
            gap: 0.25rem;
            padding: 0.25rem 0.75rem;
            font-size: var(--i56-font-size-xs);
            font-family: var(--i56-font-family);
            color: var(--i56-color-text-secondary);
            background: var(--i56-color-bg);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-full);
            cursor: pointer;
            transition: all var(--i56-transition);
          }
          .suggestion-chip:hover {
            color: var(--i56-color-brand);
            border-color: var(--i56-color-brand);
            background: var(--i56-color-brand-light);
          }
          .suggestion-chip:focus-visible {
            box-shadow: 0 0 0 2px var(--i56-color-brand-light);
            outline: none;
          }
          .chip-icon { font-size: var(--i56-font-size-xs); }

          /* Tool hints */
          .tool-hints {
            display: ${isTools ? 'flex' : 'none'};
            flex-wrap: wrap;
            gap: 0.375rem;
            padding: 0 1.125rem 0.75rem;
          }
          .tool-hint {
            display: inline-flex;
            align-items: center;
            padding: 0.25rem 0.625rem;
            font-size: var(--i56-font-size-xs);
            font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', monospace;
            color: var(--i56-color-text-tertiary);
            background: var(--i56-color-bg-secondary);
            border-radius: var(--i56-radius);
            transition: all var(--i56-transition);
          }
          .tool-hint.active {
            color: var(--i56-color-brand);
            background: var(--i56-color-brand-light);
            border: 1px solid var(--i56-color-brand);
          }

          /* Status bar */
          .status-bar {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0.5rem 1.125rem;
            border-top: 1px solid var(--i56-color-border);
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
          }
          .status-hint { display: flex; align-items: center; gap: 0.75rem; }
          .status-hint kbd {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            min-width: 1.125rem;
            height: 1.125rem;
            padding: 0 0.25rem;
            background: var(--i56-color-bg-secondary);
            border: 1px solid var(--i56-color-border);
            border-radius: 3px;
            font-family: monospace;
            font-size: 0.625rem;
          }

          .model-badge {
            display: flex;
            align-items: center;
            gap: 0.25rem;
            padding: 0.125rem 0.5rem;
            border-radius: var(--i56-radius-full);
            font-size: var(--i56-font-size-xs);
            font-weight: 500;
          }
          ${isTools ? '.model-badge { color: var(--i56-color-brand); background: var(--i56-color-brand-light); }' : ''}

          /* Backdrop */
          .backdrop {
            display: block;
            position: fixed;
            inset: 0;
            z-index: 1099;
            background: rgba(15, 23, 42, 0.2);
            backdrop-filter: blur(6px);
            -webkit-backdrop-filter: blur(6px);
            animation: i56-fade-in 150ms ease;
          }
          @keyframes i56-fade-in {
            from { opacity: 0; }
            to { opacity: 1; }
          }
        </style>

        <div class="backdrop"></div>
        <div class="ai-bar-overlay" role="dialog" aria-label="AI Command Bar">
          <div class="ai-bar-container">
            <div class="input-row">
              <span class="ai-icon" aria-hidden="true">${isLoading ? '⏳' : '✦'}</span>
              <textarea
                class="ai-input"
                rows="1"
                placeholder="${placeholder}"
                ?disabled="${isLoading}"
                aria-label="AI prompt input"
              ></textarea>
              <button class="submit-btn" ?disabled="${isLoading}"
                      aria-label="${isLoading ? 'Processing\u2026' : 'Submit'}"
                      title="${isLoading ? 'Processing\u2026' : 'Submit'}">
                ${isLoading ? '<span class="spinner"></span>' : '↵'}
              </button>
            </div>
            ${this._renderSuggestions()}
            ${this._renderToolHints()}
            <div class="status-bar">
              <span class="status-hint">
                <span><kbd>↵</kbd> Submit</span>
                <span><kbd>Shift</kbd>+<kbd>↵</kbd> New line</span>
                <span><kbd>Esc</kbd> Close</span>
              </span>
              <span>
                <kbd>${navigator.platform.includes('Mac') ? '⌘' : 'Ctrl'}</kbd>+<kbd>K</kbd> Toggle
                ${model ? '<span class="model-badge">' + model + '</span>' : ''}
                ${isTools ? '<span class="model-badge">🛠 Tools on</span>' : ''}
              </span>
            </div>
          </div>
        </div>
      `;

      // Wire events
      const textarea = this._root.querySelector('.ai-input');
      const submitBtn = this._root.querySelector('.submit-btn');
      const backdrop = this._root.querySelector('.backdrop');

      if (textarea) {
        textarea.addEventListener('keydown', this._onKeydown);
        textarea.addEventListener('input', () => this._resize(textarea));
      }
      if (submitBtn && !isLoading) {
        submitBtn.addEventListener('click', () => this._onSubmit());
      }
      if (backdrop) {
        backdrop.addEventListener('click', () => this.close());
      }

      // Suggestion clicks
      this._root.querySelectorAll('.suggestion-chip').forEach(chip => {
        chip.addEventListener('click', (e) => {
          e.stopPropagation();
          const idx = parseInt(chip.dataset.idx);
          const suggestions = JSON.parse(this.getAttribute('suggestions') || '[]');
          const suggestion = typeof suggestions[idx] === 'object' ? suggestions[idx].label : suggestions[idx];
          this._onSuggestionClick(suggestion, idx);
        });
      });

      // Tool hint clicks
      this._root.querySelectorAll('.tool-hint').forEach(hint => {
        hint.addEventListener('click', () => {
          const textarea = this._root.querySelector('.ai-input');
          if (textarea) {
            textarea.value += ' ' + hint.textContent.trim();
            this._resize(textarea);
            textarea.focus();
          }
        });
      });
    }
  }

  customElements.define('i56-ai-bar', I56AiBar);
})();
