/**
 * I56 AI Bar — AI command bar with natural language input.
 *
 * A floating command bar at the bottom of the page that accepts natural language
 * prompts. Inspired by the design patterns of modern AI coding assistants.
 *
 * Usage:
 *   <i56-ai-bar placeholder="Ask anything…" suggestions='["Summarize this page","Fix this bug","Add feature"]'>
 *   </i56-ai-bar>
 *
 *   // Or activate programmatically:
 *   const bar = document.querySelector('i56-ai-bar');
 *   bar.open();
 *   bar.close();
 *   bar.focus();
 *
 * Attributes:
 *   placeholder: input placeholder (default: "Ask AI anything…")
 *   suggestions: JSON array of suggestion labels
 *   loading: present during AI processing
 *   expanded: present when the bar is open/visible
 *   position: bottom | top (default: bottom)
 *   persistent: if present, bar stays visible at all times
 *
 * Properties:
 *   value: get/set current input text
 *   loading: get/set loading state
 *
 * Events:
 *   i56:submit — detail: { prompt }
 *   i56:suggestion-click — detail: { suggestion, index }
 *   i56:close
 *   i56:open
 *
 * Keyboard:
 *   Ctrl+J / ⌘J — toggle the bar
 *   Enter (in textarea) — submit
 *   Shift+Enter — newline
 *   Escape — close (if not persistent)
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
          --i56-color-bg: var(--i56-bg, #FFFFFF);
          --i56-color-bg-secondary: var(--i56-bg-secondary, #F9FAFB);
          --i56-color-bg-tertiary: var(--i56-bg-tertiary, #F3F4F6);
          --i56-color-border: var(--i56-border, #E5E7EB);
          --i56-color-border-hover: var(--i56-border-hover, #D1D5DB);
          --i56-color-text: var(--i56-text, #111827);
          --i56-color-text-secondary: var(--i56-text-secondary, #6B7280);
          --i56-color-text-tertiary: var(--i56-text-tertiary, #9CA3AF);
          --i56-color-text-inverse: var(--i56-text-inverse, #FFFFFF);
          --i56-radius: 6px;
          --i56-radius-md: 8px;
          --i56-radius-lg: 12px;
          --i56-radius-full: 9999px;
          --i56-shadow-lg: 0 10px 15px -3px rgba(0,0,0,0.1), 0 4px 6px -4px rgba(0,0,0,0.1);
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
      `);
    }
    return _sharedSheet;
  }

  function emit(el, name, detail = {}) {
    el.dispatchEvent(new CustomEvent(`i56:${name}`, { bubbles: true, composed: true, detail }));
  }

  class I56AiBar extends HTMLElement {
    static get observedAttributes() {
      return ['placeholder', 'suggestions', 'loading', 'expanded', 'position', 'persistent'];
    }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
      this._attached = false;
      this._boundCloseEsc = this._onEscKey.bind(this);
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
    get persistent() { return this.hasAttribute('persistent'); }

    connectedCallback() {
      if (this._attached) return;
      this._attached = true;
      this.render();

      // Global keyboard shortcut
      document.addEventListener('keydown', this._globalShortcut.bind(this));

      // Auto-open if persistent
      if (this.persistent) {
        this.setAttribute('expanded', '');
      }
    }

    disconnectedCallback() {
      document.removeEventListener('keydown', this._globalShortcut.bind(this));
      document.removeEventListener('keydown', this._boundCloseEsc);
    }

    attributeChangedCallback(name, oldVal, newVal) {
      if (name === 'expanded') {
        if (newVal !== null) document.addEventListener('keydown', this._boundCloseEsc);
        else document.removeEventListener('keydown', this._boundCloseEsc);
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
    _globalShortcut(e) {
      if ((e.ctrlKey || e.metaKey) && e.key === 'j') {
        e.preventDefault();
        this.toggle();
      }
    }

    _onEscKey(e) {
      if (e.key === 'Escape' && !this.persistent) {
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
      emit(this, 'submit', { prompt });
    }

    _onKeydown = (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        this._onSubmit();
      } else if (e.key === 'Escape' && !this.persistent) {
        this.close();
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

    _renderSuggestions() {
      let suggestions = [];
      try { suggestions = JSON.parse(this.getAttribute('suggestions') || '[]'); } catch { /* keep empty */ }

      if (!suggestions.length) return '';

      const chips = suggestions.map((s, i) => {
        // Support {label, icon} shape or plain string
        const label = typeof s === 'object' ? s.label : s;
        const icon = typeof s === 'object' ? s.icon || '' : '';
        return `
          <button class="suggestion-chip" data-idx="${i}" type="button" ?disabled="${this.loading}">
            ${icon ? `<span class="chip-icon">${icon}</span>` : ''}
            <span>${label}</span>
          </button>
        `;
      }).join('');

      return `<div class="suggestions-row">${chips}</div>`;
    }

    render() {
      const placeholder = this.getAttribute('placeholder') || 'Ask AI anything…';
      const isExpanded = this.expanded;
      const isLoading = this.loading;
      const position = this.getAttribute('position') || 'bottom';
      const isBottom = position === 'bottom';

      this._root.innerHTML = `
        <style>
          :host {
            display: contents;
          }

          .ai-bar-overlay {
            display: ${isExpanded || this.persistent ? 'flex' : 'none'};
            position: fixed;
            ${isBottom ? 'bottom' : 'top'}: 0;
            left: 50%;
            transform: translateX(-50%);
            z-index: 1100;
            width: 100%;
            max-width: 48rem;
            padding: 1rem;
            flex-direction: column;
            align-items: ${isExpanded ? 'stretch' : 'center'};
            pointer-events: ${isExpanded || this.persistent ? 'auto' : 'none'};
            animation: i56-ai-bar-in var(--i56-transition-slow) cubic-bezier(0.4, 0, 0.2, 1);
          }

          @keyframes i56-ai-bar-in {
            from {
              ${isBottom ? 'transform: translateX(-50%) translateY(1rem);' : 'transform: translateX(-50%) translateY(-1rem);'}
              opacity: 0;
            }
            to {
              ${isBottom ? 'transform: translateX(-50%) translateY(0);' : 'transform: translateX(-50%) translateY(0);'}
              opacity: 1;
            }
          }

          .ai-bar-container {
            background: var(--i56-color-bg);
            border: 1px solid var(--i56-color-border);
            border-radius: ${isExpanded ? 'var(--i56-radius-lg)' : 'var(--i56-radius-full)'};
            box-shadow: ${isExpanded ? 'var(--i56-shadow-xl)' : 'var(--i56-shadow-lg)'};
            transition: border-radius var(--i56-transition-slow), box-shadow var(--i56-transition);
            overflow: hidden;
          }
          ${!isExpanded && !this.persistent ? `
            .ai-bar-container {
              cursor: pointer;
            }
            .ai-bar-container:hover {
              border-color: var(--i56-color-border-hover);
            }
          ` : ''}

          .input-row {
            display: flex;
            align-items: flex-start;
            gap: 0.5rem;
            padding: ${isExpanded ? '0.75rem 1rem' : '0.625rem 1rem'};
          }

          .ai-icon {
            flex-shrink: 0;
            font-size: var(--i56-font-size-lg);
            display: flex;
            align-items: center;
            padding-top: ${isExpanded ? '0.25rem' : '0.125rem'};
            ${isLoading ? 'animation: i56-ai-pulse 1.5s ease-in-out infinite;' : ''}
          }
          @keyframes i56-ai-pulse {
            0%, 100% { opacity: 1; transform: scale(1); }
            50% { opacity: 0.7; transform: scale(0.95); }
          }

          .ai-input {
            flex: 1;
            border: none;
            outline: none;
            font-size: var(--i56-font-size-sm);
            font-family: var(--i56-font-family);
            color: var(--i56-color-text);
            background: transparent;
            resize: none;
            min-height: 1.5rem;
            max-height: 200px;
            line-height: var(--i56-line-height);
            padding: 0;
            cursor: ${!isExpanded && !this.persistent ? 'pointer' : 'text'};
          }
          .ai-input::placeholder { color: var(--i56-color-text-tertiary); }
          .ai-input:disabled { opacity: 0.6; }

          .submit-btn {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 2rem;
            height: 2rem;
            border: none;
            border-radius: var(--i56-radius);
            background: ${isLoading ? 'var(--i56-color-bg-tertiary)' : 'var(--i56-color-brand)'};
            color: var(--i56-color-text-inverse);
            cursor: ${isLoading ? 'wait' : 'pointer'};
            font-size: var(--i56-font-size-base);
            flex-shrink: 0;
            transition: all var(--i56-transition);
          }
          .submit-btn:hover:not(:disabled) { background: var(--i56-color-brand-hover); }
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
            display: ${isExpanded ? 'flex' : 'none'};
            flex-wrap: wrap;
            gap: 0.375rem;
            padding: 0 1rem 0.75rem;
          }
          .suggestion-chip {
            display: inline-flex;
            align-items: center;
            gap: 0.25rem;
            padding: 0.25rem 0.625rem;
            font-size: var(--i56-font-size-xs);
            font-family: var(--i56-font-family);
            color: var(--i56-color-text-secondary);
            background: var(--i56-color-bg-secondary);
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

          /* Status bar */
          .status-bar {
            display: ${isExpanded ? 'flex' : 'none'};
            align-items: center;
            justify-content: space-between;
            padding: 0.375rem 1rem;
            border-top: 1px solid var(--i56-color-border);
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
          }
          .status-hint { display: flex; align-items: center; gap: 0.75rem; }
          .status-hint kbd {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            min-width: 1rem;
            height: 1rem;
            padding: 0 0.25rem;
            background: var(--i56-color-bg-secondary);
            border: 1px solid var(--i56-color-border);
            border-radius: 3px;
            font-family: monospace;
            font-size: 0.625rem;
          }

          /* Backdrop when expanded */
          .backdrop {
            display: ${isExpanded ? 'block' : 'none'};
            position: fixed;
            inset: 0;
            z-index: 1099;
            background: rgba(0,0,0,0.2);
            backdrop-filter: blur(2px);
            animation: i56-fade-in 150ms ease;
          }
          @keyframes i56-fade-in {
            from { opacity: 0; }
            to { opacity: 1; }
          }
        </style>

        ${isExpanded ? '<div class="backdrop"></div>' : ''}
        <div class="ai-bar-overlay">
          <div class="ai-bar-container">
            <div class="input-row">
              <span class="ai-icon" aria-hidden="true">${isLoading ? '⏳' : '✨'}</span>
              <textarea
                class="ai-input"
                rows="1"
                placeholder="${placeholder}"
                ?disabled="${isLoading}"
                aria-label="AI prompt input"
              ></textarea>
              <button class="submit-btn" ?disabled="${isLoading}" aria-label="${isLoading ? 'Processing…' : 'Submit prompt'}" title="${isLoading ? 'Processing…' : 'Submit'}">
                ${isLoading ? '<span class="spinner"></span>' : '↑'}
              </button>
            </div>
            ${this._renderSuggestions()}
            <div class="status-bar">
              <span class="status-hint">
                <span><kbd>↵</kbd> Submit</span>
                <span><kbd>Shift</kbd>+<kbd>↵</kbd> New line</span>
                ${!this.persistent ? '<span><kbd>Esc</kbd> Close</span>' : ''}
              </span>
              <span><kbd>⌘</kbd>+<kbd>J</kbd> Toggle</span>
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

      // Expand on focus if minimized
      if (textarea && !isExpanded && !this.persistent) {
        textarea.addEventListener('focus', () => this.open(), { once: true });
      }

      // Click to open when minimized
      if (!isExpanded && !this.persistent) {
        const container = this._root.querySelector('.ai-bar-container');
        if (container) {
          container.addEventListener('click', () => this.open());
        }
      }

      // Suggestion clicks
      this._root.querySelectorAll('.suggestion-chip').forEach(chip => {
        chip.addEventListener('click', (e) => {
          e.stopPropagation();
          const idx = parseInt(chip.dataset.idx);
          const suggestions = JSON.parse(this.getAttribute('suggestions') || '[]');
          const suggestion = suggestions[idx] || '';
          const label = typeof suggestion === 'object' ? suggestion.label : suggestion;
          this._onSuggestionClick(label, idx);
        });
      });

      // Backdrop click to close
      if (backdrop && !this.persistent) {
        backdrop.addEventListener('click', () => this.close());
      }
    }
  }

  customElements.define('i56-ai-bar', I56AiBar);
})();
