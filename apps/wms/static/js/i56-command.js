/**
 * I56 Command Palette
 * Linear.app-style command palette overlay with fuzzy search.
 * Activated by Ctrl+K / ⌘K globally.
 *
 * Usage:
 *   I56Command.register({ id: 'dashboard', name: 'Dashboard', icon: '📊', shortcut: 'G D', group: 'Navigation', action: () => location.href = '/' });
 *   I56Command.open();      // Open the palette
 *   I56Command.toggle();    // Toggle open/close
 */

(function () {
  'use strict';

  // =============================================================================
  // Shared sheet (reuse from i56-components if loaded, otherwise create)
  // =============================================================================
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
          --i56-color-text: var(--i56-text, #111827);
          --i56-color-text-secondary: var(--i56-text-secondary, #6B7280);
          --i56-color-text-tertiary: var(--i56-text-tertiary, #9CA3AF);
          --i56-color-text-inverse: var(--i56-text-inverse, #FFFFFF);
          --i56-radius: 6px;
          --i56-radius-lg: 12px;
          --i56-shadow-xl: 0 20px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1);
          --i56-font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
          --i56-font-size-xs: 0.75rem;
          --i56-font-size-sm: 0.875rem;
          --i56-font-size-base: 1rem;
        }
      `);
    }
    return _sharedSheet;
  }

  // =============================================================================
  // Command registry
  // =============================================================================
  const commands = [];
  const listeners = [];

  const CommandAPI = {
    /**
     * Register a command.
     * @param {{id: string, name: string, icon?: string, shortcut?: string, group?: string, action: Function}} cmd
     */
    register(cmd) {
      if (!cmd.id || !cmd.name) throw new Error('Command must have id and name');
      // Remove existing with same id
      const idx = commands.findIndex(c => c.id === cmd.id);
      if (idx >= 0) commands.splice(idx, 1);
      commands.push({
        id: cmd.id,
        name: cmd.name,
        icon: cmd.icon || '',
        shortcut: cmd.shortcut || '',
        group: cmd.group || 'Actions',
        action: cmd.action || (() => {})
      });
      // Notify listeners
      listeners.forEach(fn => fn(commands));
    },

    /** Remove a command by id */
    unregister(id) {
      const idx = commands.findIndex(c => c.id === id);
      if (idx >= 0) {
        commands.splice(idx, 1);
        listeners.forEach(fn => fn(commands));
      }
    },

    /** Subscribe to command registry changes */
    onChange(fn) { listeners.push(fn); },

    /** Get all registered commands */
    getAll() { return [...commands]; },

    /** Open the palette */
    open() {
      let palette = document.querySelector('i56-command-palette');
      if (!palette) {
        palette = document.createElement('i56-command-palette');
        document.body.appendChild(palette);
      }
      palette.show();
    },

    /** Toggle the palette */
    toggle() {
      const palette = document.querySelector('i56-command-palette');
      if (palette && palette.hasAttribute('open')) {
        palette.hide();
      } else {
        this.open();
      }
    }
  };

  // =============================================================================
  // Fuzzy search
  // =============================================================================

  /**
   * Simple fuzzy-match scoring.
   * Returns a score (higher is better) or 0 for no match.
   */
  function fuzzyScore(text, query) {
    if (!query) return 1;
    text = text.toLowerCase();
    query = query.toLowerCase();

    let score = 0;
    let qi = 0;
    let consecutiveBonus = 0;
    let prevMatchIdx = -1;

    for (let i = 0; i < text.length && qi < query.length; i++) {
      if (text[i] === query[qi]) {
        score += 1;
        // Bonus for consecutive characters and start-of-word matches
        if (prevMatchIdx >= 0 && i === prevMatchIdx + 1) {
          consecutiveBonus += 2;
        }
        if (i === 0 || text[i - 1] === ' ' || text[i - 1] === '-' || text[i - 1] === '_') {
          score += 3; // Word boundary bonus
        }
        score += consecutiveBonus;
        prevMatchIdx = i;
        qi++;
      } else {
        consecutiveBonus = 0;
      }
    }

    if (qi < query.length) return 0; // Not all characters matched
    // Penalty for longer text with same match count
    return score / (text.length * 0.1 + 1);
  }

  /**
   * Highlight matching characters in text.
   */
  function highlightMatch(text, query) {
    if (!query) return escapeHtml(text);
    text = escapeHtml(text);
    const chars = query.toLowerCase().split('');
    let result = '';
    let qi = 0;
    for (let i = 0; i < text.length && qi < chars.length; i++) {
      if (text[i].toLowerCase() === chars[qi]) {
        result += '<mark>' + text[i] + '</mark>';
        qi++;
      } else {
        result += text[i];
      }
    }
    result += text.slice(result.replace(/<[^>]+>/g, '').length);
    return result;
  }

  function escapeHtml(str) {
    return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }

  // =============================================================================
  // <i56-command-palette>
  // =============================================================================

  class I56CommandPalette extends HTMLElement {
    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
      this._query = '';
      this._selectedIdx = 0;
      this._filtered = [];
      this._grouped = {};
    }

    connectedCallback() {
      // Listen for registry changes
      CommandAPI.onChange(() => { if (this.hasAttribute('open')) this._filter(); });
      this.render();
    }

    show() {
      if (this.hasAttribute('open')) return;
      this.setAttribute('open', '');
      this._query = '';
      this._selectedIdx = 0;
      this._filter();
      this.render();
      document.addEventListener('keydown', this._globalKeydown);

      // Focus input after render
      requestAnimationFrame(() => {
        const input = this._root.querySelector('.search-input');
        if (input) input.focus();
      });
    }

    hide() {
      this.removeAttribute('open');
      document.removeEventListener('keydown', this._globalKeydown);
      emit(this, 'close');
    }

    _onBackdropClick = (e) => {
      if (e.target.classList.contains('backdrop')) this.hide();
    };

    _onInput = (e) => {
      this._query = e.target.value;
      this._selectedIdx = 0;
      this._filter();
      this.render();
    };

    _onKeydown = (e) => {
      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          this._selectedIdx = Math.min(this._selectedIdx + 1, this._filtered.length - 1);
          this._updateSelection();
          break;
        case 'ArrowUp':
          e.preventDefault();
          this._selectedIdx = Math.max(this._selectedIdx - 1, 0);
          this._updateSelection();
          break;
        case 'Enter':
          e.preventDefault();
          this._executeCurrent();
          break;
        case 'Escape':
          e.preventDefault();
          this.hide();
          break;
      }
    };

    _filter() {
      const all = commands;
      if (!this._query) {
        this._filtered = [...all];
      } else {
        const scored = all.map(cmd => ({
          cmd,
          score: fuzzyScore(cmd.name, this._query) * 2 + fuzzyScore(cmd.group, this._query)
        }));
        this._filtered = scored
          .filter(s => s.score > 0)
          .sort((a, b) => b.score - a.score)
          .map(s => s.cmd);
      }

      // Group results
      this._grouped = {};
      this._filtered.forEach(cmd => {
        const group = cmd.group || 'Actions';
        if (!this._grouped[group]) this._grouped[group] = [];
        this._grouped[group].push(cmd);
      });
    }

    _updateSelection() {
      const items = this._root.querySelectorAll('.cmd-item');
      items.forEach((item, i) => {
        item.classList.toggle('selected', i === this._selectedIdx);
      });
      // Scroll into view
      const selected = items[this._selectedIdx];
      if (selected) selected.scrollIntoView({ block: 'nearest' });
    }

    _executeCurrent() {
      const cmd = this._filtered[this._selectedIdx];
      if (cmd) {
        this.hide();
        try { cmd.action(); } catch (err) { console.error('I56Command action error:', err); }
      }
    }

    _onItemClick(idx) {
      this._selectedIdx = idx;
      this._executeCurrent();
    }

    render() {
      const isOpen = this.hasAttribute('open');

      this._root.innerHTML = `
        <style>
          .backdrop {
            display: ${isOpen ? 'flex' : 'none'};
            position: fixed;
            inset: 0;
            z-index: 2000;
            align-items: flex-start;
            justify-content: center;
            padding-top: 15vh;
            background: rgba(0, 0, 0, 0.4);
            backdrop-filter: blur(4px);
            animation: i56-cmd-fade-in 120ms ease;
          }
          .palette {
            background: var(--i56-color-bg);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-lg);
            box-shadow: var(--i56-shadow-xl);
            width: 100%;
            max-width: 36rem;
            max-height: 60vh;
            display: flex;
            flex-direction: column;
            overflow: hidden;
            animation: i56-cmd-scale-in 120ms ease;
            font-family: var(--i56-font-family);
          }
          .search-area {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.75rem 1rem;
            border-bottom: 1px solid var(--i56-color-border);
          }
          .search-icon {
            color: var(--i56-color-text-tertiary);
            font-size: var(--i56-font-size-base);
            flex-shrink: 0;
          }
          .search-input {
            flex: 1;
            border: none;
            outline: none;
            font-size: var(--i56-font-size-base);
            font-family: var(--i56-font-family);
            color: var(--i56-color-text);
            background: transparent;
          }
          .search-input::placeholder { color: var(--i56-color-text-tertiary); }
          .hint {
            padding: 0.125rem 0.375rem;
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
            border: 1px solid var(--i56-color-border);
            border-radius: 4px;
            white-space: nowrap;
            flex-shrink: 0;
          }
          .results {
            flex: 1;
            overflow-y: auto;
            padding: 0.375rem;
            scroll-behavior: smooth;
          }
          .group-header {
            padding: 0.375rem 0.75rem;
            font-size: var(--i56-font-size-xs);
            font-weight: 600;
            color: var(--i56-color-text-tertiary);
            text-transform: uppercase;
            letter-spacing: 0.05em;
          }
          .cmd-item {
            display: flex;
            align-items: center;
            gap: 0.625rem;
            padding: 0.5rem 0.75rem;
            border-radius: var(--i56-radius);
            cursor: pointer;
            transition: background 80ms ease;
            font-size: var(--i56-font-size-sm);
            color: var(--i56-color-text);
          }
          .cmd-item:hover, .cmd-item.selected {
            background: var(--i56-color-brand-light);
            color: var(--i56-color-brand);
          }
          .cmd-item mark {
            background: transparent;
            color: var(--i56-color-brand);
            font-weight: 600;
          }
          .cmd-item.selected mark { color: var(--i56-color-brand-hover); }
          .cmd-icon {
            flex-shrink: 0;
            width: 1.5rem;
            text-align: center;
            font-size: var(--i56-font-size-base);
          }
          .cmd-name {
            flex: 1;
            min-width: 0;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
          }
          .cmd-shortcut {
            flex-shrink: 0;
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
            font-family: monospace;
            letter-spacing: 0.05em;
          }
          .cmd-item.selected .cmd-shortcut { color: var(--i56-color-brand); opacity: 0.7; }
          .empty {
            padding: 2rem 1rem;
            text-align: center;
            color: var(--i56-color-text-tertiary);
            font-size: var(--i56-font-size-sm);
          }
          .footer {
            display: flex;
            gap: 1rem;
            padding: 0.5rem 1rem;
            border-top: 1px solid var(--i56-color-border);
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
          }
          .footer span {
            display: inline-flex;
            align-items: center;
            gap: 0.25rem;
          }
          .footer kbd {
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
          @keyframes i56-cmd-fade-in { from { opacity: 0; } to { opacity: 1; } }
          @keyframes i56-cmd-scale-in { from { opacity: 0; transform: scale(0.96) translateY(-0.5rem); } to { opacity: 1; transform: scale(1) translateY(0); } }
        </style>
        <div class="backdrop">
          <div class="palette" role="dialog" aria-modal="true" aria-label="Command Palette">
            <div class="search-area">
              <span class="search-icon" aria-hidden="true">⌘</span>
              <input class="search-input" type="text" placeholder="Type a command or search..." value="${escapeHtml(this._query)}" aria-label="Search commands" />
              <span class="hint">ESC</span>
            </div>
            <div class="results">
              ${this._renderResults()}
            </div>
            <div class="footer">
              <span><kbd>↑↓</kbd> Navigate</span>
              <span><kbd>↵</kbd> Select</span>
              <span><kbd>Esc</kbd> Close</span>
            </div>
          </div>
        </div>
      `;

      if (isOpen) {
        this._root.querySelector('.backdrop').addEventListener('click', this._onBackdropClick);
        const input = this._root.querySelector('.search-input');
        if (input) {
          input.addEventListener('input', this._onInput);
          input.addEventListener('keydown', this._onKeydown);
        }

        // Click handlers on items
        this._root.querySelectorAll('.cmd-item').forEach((item, i) => {
          item.addEventListener('click', () => this._onItemClick(i));
        });
      }
    }

    _renderResults() {
      const groups = Object.keys(this._grouped);
      if (groups.length === 0) {
        return '<div class="empty">No matching commands</div>';
      }

      let idx = 0;
      return groups.map(group => {
        const cmds = this._grouped[group];
        const items = cmds.map(cmd => {
          const isSelected = idx === this._selectedIdx;
          const item = `
            <div class="cmd-item${isSelected ? ' selected' : ''}" role="option" aria-selected="${isSelected}">
              <span class="cmd-icon">${escapeHtml(cmd.icon || '')}</span>
              <span class="cmd-name">${highlightMatch(cmd.name, this._query)}</span>
              ${cmd.shortcut ? `<span class="cmd-shortcut">${escapeHtml(cmd.shortcut)}</span>` : ''}
            </div>
          `;
          idx++;
          return item;
        }).join('');

        return `<div class="group-header">${escapeHtml(group)}</div>${items}`;
      }).join('');
    }
  }

  // =============================================================================
  // Global keyboard listener
  // =============================================================================

  document.addEventListener('keydown', (e) => {
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
      e.preventDefault();
      CommandAPI.toggle();
    }
  });

  // =============================================================================
  // Pre-register core commands
  // =============================================================================

  CommandAPI.register({
    id: 'dashboard',
    name: 'Dashboard',
    icon: '📊',
    shortcut: 'G D',
    group: 'Navigation',
    action: () => console.log('Navigate to Dashboard')
  });

  CommandAPI.register({
    id: 'orders',
    name: 'Orders',
    icon: '📦',
    shortcut: 'G O',
    group: 'Navigation',
    action: () => console.log('Navigate to Orders')
  });

  CommandAPI.register({
    id: 'parcels',
    name: 'Parcels',
    icon: '🚚',
    shortcut: 'G P',
    group: 'Navigation',
    action: () => console.log('Navigate to Parcels')
  });

  CommandAPI.register({
    id: 'settings',
    name: 'Settings',
    icon: '⚙️',
    shortcut: 'G S',
    group: 'Settings',
    action: () => console.log('Open Settings')
  });

  CommandAPI.register({
    id: 'theme-toggle',
    name: 'Theme Toggle',
    icon: '🌓',
    shortcut: '⌘⇧T',
    group: 'Settings',
    action: () => {
      document.documentElement.classList.toggle('theme-dark');
      console.log('Theme toggled');
    }
  });

  CommandAPI.register({
    id: 'help',
    name: 'Help',
    icon: '❓',
    shortcut: '?',
    group: 'Actions',
    action: () => console.log('Open Help')
  });

  // =============================================================================
  // Helpers
  // =============================================================================

  function emit(el, name, detail = {}) {
    el.dispatchEvent(new CustomEvent(`i56:${name}`, { bubbles: true, composed: true, detail }));
  }

  // =============================================================================
  // Register
  // =============================================================================

  customElements.define('i56-command-palette', I56CommandPalette);

  // Expose global API
  window.I56Command = CommandAPI;

})();
