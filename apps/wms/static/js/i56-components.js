/**
 * I56 Web Components Library
 * Vanilla Web Components with Shadow DOM, CSS custom properties, keyboard accessibility.
 * No external dependencies.
 *
 * Design tokens are set on :host via CSS custom properties, inheriting from :root or dark/light themes.
 *
 * Components:
 *   i56-button, i56-card, i56-table, i56-form-group, i56-input,
 *   i56-select, i56-modal, i56-toast, i56-badge, i56-tabs,
 *   i56-avatar, i56-spinner, i56-timeline
 */

// =============================================================================
// Shared helpers
// =============================================================================

/**
 * Adopt the global I56 design-system stylesheet into a shadow root.
 * The stylesheet is created once and shared across all component instances.
 */
let _sharedSheet = null;

function getSharedSheet() {
  if (!_sharedSheet) {
    _sharedSheet = new CSSStyleSheet();
    _sharedSheet.replaceSync(`
      /* I56 Design Tokens — override these on :root, .theme-dark, .theme-light */
      :host {
        /* Colors */
        --i56-color-brand: var(--i56-brand, #4F46E5);
        --i56-color-brand-hover: var(--i56-brand-hover, #4338CA);
        --i56-color-brand-light: var(--i56-brand-light, #EEF2FF);
        --i56-color-success: var(--i56-success, #059669);
        --i56-color-success-light: var(--i56-success-light, #ECFDF5);
        --i56-color-warning: var(--i56-warning, #D97706);
        --i56-color-warning-light: var(--i56-warning-light, #FFFBEB);
        --i56-color-danger: var(--i56-danger, #DC2626);
        --i56-color-danger-light: var(--i56-danger-light, #FEF2F2);
        --i56-color-info: var(--i56-info, #2563EB);
        --i56-color-info-light: var(--i56-info-light, #EFF6FF);
        --i56-color-neutral: var(--i56-neutral, #6B7280);
        --i56-color-neutral-light: var(--i56-neutral-light, #F9FAFB);

        /* Surfaces */
        --i56-color-bg: var(--i56-bg, #FFFFFF);
        --i56-color-bg-secondary: var(--i56-bg-secondary, #F9FAFB);
        --i56-color-bg-tertiary: var(--i56-bg-tertiary, #F3F4F6);
        --i56-color-border: var(--i56-border, #E5E7EB);
        --i56-color-border-hover: var(--i56-border-hover, #D1D5DB);

        /* Text */
        --i56-color-text: var(--i56-text, #111827);
        --i56-color-text-secondary: var(--i56-text-secondary, #6B7280);
        --i56-color-text-tertiary: var(--i56-text-tertiary, #9CA3AF);
        --i56-color-text-inverse: var(--i56-text-inverse, #FFFFFF);

        /* Sizing */
        --i56-radius-sm: 4px;
        --i56-radius: 6px;
        --i56-radius-md: 8px;
        --i56-radius-lg: 12px;
        --i56-radius-full: 9999px;

        /* Shadows */
        --i56-shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
        --i56-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1);
        --i56-shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1);
        --i56-shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1);
        --i56-shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);

        /* Transitions */
        --i56-transition: 150ms ease;
        --i56-transition-slow: 300ms ease;

        /* Typography */
        --i56-font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
        --i56-font-size-xs: 0.75rem;
        --i56-font-size-sm: 0.875rem;
        --i56-font-size-base: 1rem;
        --i56-font-size-lg: 1.125rem;
        --i56-font-size-xl: 1.25rem;
        --i56-font-size-2xl: 1.5rem;

        --i56-line-height: 1.5;
      }

      /* Dark theme overrides */
      :host(.theme-dark),
      .theme-dark :host {
        --i56-color-bg: #1F2937;
        --i56-color-bg-secondary: #111827;
        --i56-color-bg-tertiary: #374151;
        --i56-color-border: #374151;
        --i56-color-border-hover: #4B5563;
        --i56-color-text: #F9FAFB;
        --i56-color-text-secondary: #D1D5DB;
        --i56-color-text-tertiary: #9CA3AF;
        --i56-color-neutral-light: #374151;
      }
    `);
  }
  return _sharedSheet;
}

/**
 * Emit a custom event with the 'i56:' prefix.
 */
function emit(el, name, detail = {}) {
  el.dispatchEvent(new CustomEvent(`i56:${name}`, { bubbles: true, composed: true, detail }));
}

/**
 * Fill a slot-based container with plain HTML.
 */
function fillSlot(el, slotName, html) {
  const existing = el.querySelector(`[slot="${slotName}"]`);
  if (existing) existing.remove();
  if (html != null) {
    const wrapper = document.createElement('span');
    wrapper.setAttribute('slot', slotName);
    wrapper.innerHTML = html;
    el.appendChild(wrapper);
  }
}

// =============================================================================
// 1. <i56-button>
// =============================================================================

class I56Button extends HTMLElement {
  static get observedAttributes() {
    return ['variant', 'size', 'disabled', 'loading', 'icon-prefix', 'icon-suffix', 'shortcut'];
  }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this.render();
  }

  attributeChangedCallback() { this.render(); }

  connectedCallback() {
    this.setAttribute('role', 'button');
    if (!this.hasAttribute('tabindex')) this.setAttribute('tabindex', '0');
    this.addEventListener('click', this._onClick);
    this.addEventListener('keydown', this._onKeydown);
  }

  disconnectedCallback() {
    this.removeEventListener('click', this._onClick);
    this.removeEventListener('keydown', this._onKeydown);
  }

  _onClick = (e) => {
    if (this.hasAttribute('disabled') || this.hasAttribute('loading')) {
      e.preventDefault();
      e.stopPropagation();
      return;
    }
    emit(this, 'click');
  };

  _onKeydown = (e) => {
    if ((e.key === 'Enter' || e.key === ' ') && !this.hasAttribute('disabled') && !this.hasAttribute('loading')) {
      e.preventDefault();
      emit(this, 'click');
      this.click();
    }
  };

  render() {
    const variant = this.getAttribute('variant') || 'primary';
    const size = this.getAttribute('size') || 'md';
    const disabled = this.hasAttribute('disabled');
    const loading = this.hasAttribute('loading');
    const iconPrefix = this.getAttribute('icon-prefix') || '';
    const iconSuffix = this.getAttribute('icon-suffix') || '';
    const shortcut = this.getAttribute('shortcut') || '';

    const isDisabled = disabled || loading;

    this._root.innerHTML = `
      <style>
        :host {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          gap: 0.5rem;
          font-family: var(--i56-font-family);
          font-weight: 500;
          line-height: 1;
          border: 1px solid transparent;
          border-radius: var(--i56-radius);
          cursor: pointer;
          user-select: none;
          white-space: nowrap;
          transition: all var(--i56-transition);
          outline: none;
          position: relative;
        }
        :host(:focus-visible) {
          box-shadow: 0 0 0 3px var(--i56-color-brand-light);
        }
        :host([disabled]), :host([loading]) {
          opacity: 0.5;
          cursor: not-allowed;
          pointer-events: none;
        }

        /* Sizes */
        :host([size="sm"]) { padding: 0.375rem 0.75rem; font-size: var(--i56-font-size-sm); border-radius: var(--i56-radius-sm); }
        :host([size="md"]) { padding: 0.5rem 1rem; font-size: var(--i56-font-size-sm); }
        :host([size="lg"]) { padding: 0.625rem 1.25rem; font-size: var(--i56-font-size-base); }

        /* Variants */
        :host([variant="primary"]) {
          background: var(--i56-color-brand);
          color: var(--i56-color-text-inverse);
        }
        :host([variant="primary"]:hover:not([disabled]):not([loading])) {
          background: var(--i56-color-brand-hover);
        }

        :host([variant="secondary"]) {
          background: var(--i56-color-bg);
          color: var(--i56-color-text);
          border-color: var(--i56-color-border);
        }
        :host([variant="secondary"]:hover:not([disabled]):not([loading])) {
          background: var(--i56-color-bg-secondary);
          border-color: var(--i56-color-border-hover);
        }

        :host([variant="ghost"]) {
          background: transparent;
          color: var(--i56-color-text-secondary);
        }
        :host([variant="ghost"]:hover:not([disabled]):not([loading])) {
          background: var(--i56-color-bg-secondary);
          color: var(--i56-color-text);
        }

        :host([variant="danger"]) {
          background: var(--i56-color-danger);
          color: var(--i56-color-text-inverse);
        }
        :host([variant="danger"]:hover:not([disabled]):not([loading])) {
          background: #B91C1C;
        }

        .spinner {
          display: none;
          width: 1em;
          height: 1em;
          border: 2px solid currentColor;
          border-right-color: transparent;
          border-radius: 50%;
          animation: i56-spin 0.6s linear infinite;
          flex-shrink: 0;
        }
        :host([loading]) .spinner { display: block; }
        :host([loading]) .content { visibility: hidden; }

        .content {
          display: inline-flex;
          align-items: center;
          gap: 0.5rem;
        }
        .shortcut {
          margin-left: 0.5rem;
          padding: 0.125rem 0.375rem;
          font-size: var(--i56-font-size-xs);
          background: rgba(0,0,0,0.08);
          border-radius: var(--i56-radius-sm);
          opacity: 0.7;
        }
        :host([variant="primary"]) .shortcut,
        :host([variant="danger"]) .shortcut {
          background: rgba(255,255,255,0.2);
        }

        @keyframes i56-spin {
          to { transform: rotate(360deg); }
        }
      </style>
      <span class="spinner"></span>
      <span class="content">
        ${iconPrefix ? `<span class="icon-prefix">${iconPrefix}</span>` : ''}
        <slot></slot>
        ${iconSuffix ? `<span class="icon-suffix">${iconSuffix}</span>` : ''}
      </span>
      ${shortcut ? `<span class="shortcut">${shortcut}</span>` : ''}
    `;
  }
}

// =============================================================================
// 2. <i56-card>
// =============================================================================

class I56Card extends HTMLElement {
  static get observedAttributes() { return ['hover', 'clickable', 'compact']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this.render();
  }

  attributeChangedCallback() { this.render(); }

  connectedCallback() {
    if (this.hasAttribute('clickable')) {
      this.setAttribute('role', 'button');
      this.setAttribute('tabindex', '0');
      this.addEventListener('click', this._onClick);
      this.addEventListener('keydown', this._onKeydown);
    }
  }

  disconnectedCallback() {
    this.removeEventListener('click', this._onClick);
    this.removeEventListener('keydown', this._onKeydown);
  }

  _onClick = () => { emit(this, 'click'); };
  _onKeydown = (e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); emit(this, 'click'); } };

  render() {
    const compact = this.hasAttribute('compact');
    this._root.innerHTML = `
      <style>
        :host {
          display: block;
          background: var(--i56-color-bg);
          border: 1px solid var(--i56-color-border);
          border-radius: var(--i56-radius-lg);
          overflow: hidden;
          transition: box-shadow var(--i56-transition), border-color var(--i56-transition);
        }
        :host([hover]:hover) {
          box-shadow: var(--i56-shadow-md);
          border-color: var(--i56-color-border-hover);
        }
        :host([clickable]) { cursor: pointer; }
        :host([clickable]:hover) { border-color: var(--i56-color-brand); }
        :host(:focus-visible) { box-shadow: 0 0 0 3px var(--i56-color-brand-light); outline: none; }

        .card-header {
          padding: ${compact ? '0.75rem 1rem' : '1rem 1.25rem'};
          border-bottom: 1px solid var(--i56-color-border);
          font-weight: 600;
          font-size: var(--i56-font-size-base);
          color: var(--i56-color-text);
          font-family: var(--i56-font-family);
        }
        .card-header:empty { display: none; }
        .card-body {
          padding: ${compact ? '0.75rem 1rem' : '1rem 1.25rem'};
          font-family: var(--i56-font-family);
          color: var(--i56-color-text-secondary);
          font-size: var(--i56-font-size-sm);
          line-height: var(--i56-line-height);
        }
        .card-body:empty { display: none; }
        .card-footer {
          padding: ${compact ? '0.75rem 1rem' : '1rem 1.25rem'};
          border-top: 1px solid var(--i56-color-border);
          background: var(--i56-color-bg-secondary);
          font-family: var(--i56-font-family);
        }
        .card-footer:empty { display: none; }
      </style>
      <div class="card-header"><slot name="header"></slot></div>
      <div class="card-body"><slot></slot></div>
      <div class="card-footer"><slot name="footer"></slot></div>
    `;
  }
}

// =============================================================================
// 3. <i56-table>
// =============================================================================

class I56Table extends HTMLElement {
  static get observedAttributes() { return ['data', 'columns', 'striped', 'hoverable', 'sortable', 'selectable', 'empty-message']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this._sortKey = null;
    this._sortDir = 'asc';
    this._selectedRows = new Set();
    this.render();
  }

  get _data() {
    try { return JSON.parse(this.getAttribute('data') || '[]'); } catch { return []; }
  }

  get _columns() {
    try { return JSON.parse(this.getAttribute('columns') || '[]'); } catch { return []; }
  }

  attributeChangedCallback(name) {
    if (name === 'data' || name === 'columns') {
      this._sortKey = null;
      this._sortDir = 'asc';
      this._selectedRows.clear();
    }
    this.render();
  }

  _onHeaderClick(col) {
    if (!this.hasAttribute('sortable')) return;
    if (this._sortKey === col) {
      this._sortDir = this._sortDir === 'asc' ? 'desc' : 'asc';
    } else {
      this._sortKey = col;
      this._sortDir = 'asc';
    }
    this.render();
    emit(this, 'sort', { key: this._sortKey, direction: this._sortDir });
  }

  _onSelectAll(e) {
    const checked = e.target.checked;
    const data = this._sortedData();
    if (checked) {
      data.forEach((_, i) => this._selectedRows.add(i));
    } else {
      this._selectedRows.clear();
    }
    this.render();
    emit(this, 'select', { selected: [...this._selectedRows] });
  }

  _onSelectRow(idx, e) {
    if (e.target.checked) this._selectedRows.add(idx);
    else this._selectedRows.delete(idx);
    emit(this, 'select', { selected: [...this._selectedRows] });
    this._updateSelectAll();
  }

  _updateSelectAll() {
    const checkbox = this._root.querySelector('.select-all-checkbox');
    if (!checkbox) return;
    const data = this._sortedData();
    checkbox.checked = data.length > 0 && this._selectedRows.size === data.length;
    checkbox.indeterminate = this._selectedRows.size > 0 && this._selectedRows.size < data.length;
  }

  _sortedData() {
    const data = [...this._data];
    if (this._sortKey && this.hasAttribute('sortable')) {
      data.sort((a, b) => {
        const va = a[this._sortKey] ?? '', vb = b[this._sortKey] ?? '';
        const cmp = typeof va === 'number' && typeof vb === 'number' ? va - vb : String(va).localeCompare(String(vb));
        return this._sortDir === 'asc' ? cmp : -cmp;
      });
    }
    return data;
  }

  render() {
    const data = this._sortedData();
    const columns = this._columns;
    const emptyMessage = this.getAttribute('empty-message') || 'No data available';
    const selectable = this.hasAttribute('selectable');
    const striped = this.hasAttribute('striped');
    const hoverable = this.hasAttribute('hoverable');
    const sortable = this.hasAttribute('sortable');

    const colKeys = columns.length > 0 ? columns : (data.length > 0 ? Object.keys(data[0]) : []);

    this._root.innerHTML = `
      <style>
        :host { display: block; font-family: var(--i56-font-family); }
        .table-wrapper { overflow-x: auto; border: 1px solid var(--i56-color-border); border-radius: var(--i56-radius-lg); }
        table { width: 100%; border-collapse: collapse; font-size: var(--i56-font-size-sm); }
        th, td { padding: 0.625rem 1rem; text-align: left; border-bottom: 1px solid var(--i56-color-border); white-space: nowrap; }
        th {
          font-weight: 600;
          color: var(--i56-color-text-secondary);
          background: var(--i56-color-bg-secondary);
          font-size: var(--i56-font-size-xs);
          text-transform: uppercase;
          letter-spacing: 0.05em;
          user-select: none;
        }
        th.sortable { cursor: pointer; }
        th.sortable:hover { color: var(--i56-color-text); background: var(--i56-color-bg-tertiary); }
        th .sort-icon { margin-left: 0.25rem; opacity: 0.4; display: inline-block; width: 0.75rem; }
        th .sort-icon.active { opacity: 1; color: var(--i56-color-brand); }
        td { color: var(--i56-color-text); }
        tbody tr:last-child td { border-bottom: none; }
        :host([striped]) tbody tr:nth-child(even) { background: var(--i56-color-bg-secondary); }
        :host([hoverable]) tbody tr:hover { background: var(--i56-color-brand-light); }
        .empty-state {
          text-align: center;
          padding: 2rem 1rem;
          color: var(--i56-color-text-tertiary);
          font-size: var(--i56-font-size-sm);
        }
        .select-col { width: 2.5rem; text-align: center; }
        input[type="checkbox"] {
          accent-color: var(--i56-color-brand);
          width: 1rem; height: 1rem; cursor: pointer;
        }
      </style>
      <div class="table-wrapper">
        ${data.length === 0 && colKeys.length === 0 ? `
          <div class="empty-state">${emptyMessage}</div>
        ` : `
          <table role="grid">
            <thead>
              <tr>
                ${selectable ? `<th class="select-col"><input type="checkbox" class="select-all-checkbox" aria-label="Select all rows" /></th>` : ''}
                ${colKeys.map(k => {
                  const key = typeof k === 'object' ? k.key : k;
                  const label = typeof k === 'object' ? k.label || key : key;
                  const isSorted = this._sortKey === key;
                  const icon = isSorted ? (this._sortDir === 'asc' ? '▲' : '▼') : '⇅';
                  return `<th class="${sortable ? 'sortable' : ''}" data-col="${key}">
                    ${label}<span class="sort-icon${isSorted ? ' active' : ''}">${sortable ? icon : ''}</span>
                  </th>`;
                }).join('')}
              </tr>
            </thead>
            <tbody>
              ${data.map((row, idx) => `
                <tr>
                  ${selectable ? `<td class="select-col"><input type="checkbox" class="row-checkbox" data-idx="${idx}" ${this._selectedRows.has(idx) ? 'checked' : ''} aria-label="Select row ${idx + 1}" /></td>` : ''}
                  ${colKeys.map(k => {
                    const key = typeof k === 'object' ? k.key : k;
                    return `<td>${row[key] ?? ''}</td>`;
                  }).join('')}
                </tr>
              `).join('')}
              ${data.length === 0 ? `<tr><td colspan="${colKeys.length + (selectable ? 1 : 0)}" class="empty-state">${emptyMessage}</td></tr>` : ''}
            </tbody>
          </table>
        `}
      </div>
    `;

    // Bind events
    this._root.querySelectorAll('th.sortable').forEach(th => {
      th.addEventListener('click', () => this._onHeaderClick(th.dataset.col));
    });
    const selectAll = this._root.querySelector('.select-all-checkbox');
    if (selectAll) selectAll.addEventListener('change', e => this._onSelectAll(e));
    this._root.querySelectorAll('.row-checkbox').forEach(cb => {
      cb.addEventListener('change', e => this._onSelectRow(parseInt(cb.dataset.idx), e));
    });
  }
}

// =============================================================================
// 4. <i56-form-group>
// =============================================================================

class I56FormGroup extends HTMLElement {
  static get observedAttributes() { return ['label', 'error', 'hint', 'required', 'label-for']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
  }

  attributeChangedCallback() { this.render(); }
  connectedCallback() { this.render(); }

  render() {
    const label = this.getAttribute('label') || '';
    const error = this.getAttribute('error') || '';
    const hint = this.getAttribute('hint') || '';
    const required = this.hasAttribute('required');
    const labelFor = this.getAttribute('label-for') || '';

    this._root.innerHTML = `
      <style>
        :host { display: block; font-family: var(--i56-font-family); margin-bottom: 1rem; }
        label {
          display: block;
          font-size: var(--i56-font-size-sm);
          font-weight: 500;
          color: var(--i56-color-text);
          margin-bottom: 0.375rem;
        }
        .required { color: var(--i56-color-danger); margin-left: 0.125rem; }
        .hint { font-size: var(--i56-font-size-xs); color: var(--i56-color-text-tertiary); margin-top: 0.25rem; }
        .error { font-size: var(--i56-font-size-xs); color: var(--i56-color-danger); margin-top: 0.25rem; }
        slot { display: block; }
      </style>
      ${label ? `<label for="${labelFor}">${label}${required ? '<span class="required" aria-hidden="true">*</span>' : ''}</label>` : ''}
      <slot></slot>
      ${error ? `<div class="error" role="alert">${error}</div>` : ''}
      ${hint && !error ? `<div class="hint">${hint}</div>` : ''}
    `;
  }
}

// =============================================================================
// 5. <i56-input>
// =============================================================================

class I56Input extends HTMLElement {
  static get observedAttributes() { return ['placeholder', 'value', 'type', 'disabled', 'error', 'icon-prefix', 'clearable', 'name']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this.render();
  }

  attributeChangedCallback(name, oldVal, newVal) {
    if (name === 'value' && newVal !== this._inputValue) {
      this._inputValue = newVal;
      const input = this._root.querySelector('input');
      if (input) input.value = newVal || '';
      this._updateClearButton();
    }
    this.render();
  }

  get value() { return this._root.querySelector('input')?.value || ''; }
  set value(v) { this.setAttribute('value', v); }

  connectedCallback() {
    this._inputValue = this.getAttribute('value') || '';
    this.render();
  }

  _onInput = (e) => {
    this._inputValue = e.target.value;
    this._updateClearButton();
    emit(this, 'change', { value: this._inputValue });
  };

  _onClear = () => {
    const input = this._root.querySelector('input');
    if (input) { input.value = ''; input.focus(); }
    this._inputValue = '';
    this._updateClearButton();
    emit(this, 'change', { value: '' });
  };

  _updateClearButton() {
    const btn = this._root.querySelector('.clear-btn');
    if (btn) btn.style.display = this._inputValue ? 'flex' : 'none';
  }

  render() {
    const placeholder = this.getAttribute('placeholder') || '';
    const type = this.getAttribute('type') || 'text';
    const disabled = this.hasAttribute('disabled');
    const error = this.hasAttribute('error');
    const iconPrefix = this.getAttribute('icon-prefix') || '';
    const clearable = this.hasAttribute('clearable');
    const name = this.getAttribute('name') || '';
    const val = (this._inputValue !== undefined ? this._inputValue : this.getAttribute('value')) || '';

    this._root.innerHTML = `
      <style>
        :host { display: block; font-family: var(--i56-font-family); }
        .wrapper {
          display: flex;
          align-items: center;
          border: 1px solid ${error ? 'var(--i56-color-danger)' : 'var(--i56-color-border)'};
          border-radius: var(--i56-radius);
          background: var(--i56-color-bg);
          transition: border-color var(--i56-transition), box-shadow var(--i56-transition);
        }
        .wrapper:focus-within {
          border-color: ${error ? 'var(--i56-color-danger)' : 'var(--i56-color-brand)'};
          box-shadow: 0 0 0 3px ${error ? 'var(--i56-color-danger-light)' : 'var(--i56-color-brand-light)'};
        }
        :host([disabled]) .wrapper { opacity: 0.5; background: var(--i56-color-bg-secondary); }
        .icon-prefix {
          display: flex;
          align-items: center;
          padding-left: 0.75rem;
          color: var(--i56-color-text-tertiary);
          font-size: var(--i56-font-size-base);
          flex-shrink: 0;
        }
        input {
          flex: 1;
          border: none;
          outline: none;
          padding: 0.5rem 0.75rem;
          font-size: var(--i56-font-size-sm);
          font-family: var(--i56-font-family);
          color: var(--i56-color-text);
          background: transparent;
          min-width: 0;
        }
        input::placeholder { color: var(--i56-color-text-tertiary); }
        input:disabled { cursor: not-allowed; }
        .clear-btn {
          display: ${clearable && val ? 'flex' : 'none'};
          align-items: center;
          justify-content: center;
          padding: 0 0.5rem;
          cursor: pointer;
          color: var(--i56-color-text-tertiary);
          border: none;
          background: none;
          font-size: var(--i56-font-size-lg);
          line-height: 1;
          flex-shrink: 0;
        }
        .clear-btn:hover { color: var(--i56-color-text); }
      </style>
      <div class="wrapper">
        ${iconPrefix ? `<span class="icon-prefix" aria-hidden="true">${iconPrefix}</span>` : ''}
        <input
          type="${type}"
          placeholder="${placeholder}"
          ?disabled="${disabled}"
          name="${name}"
          value="${val}"
          aria-invalid="${error ? 'true' : 'false'}"
        />
        ${clearable ? `<button class="clear-btn" type="button" aria-label="Clear input" tabindex="-1">&times;</button>` : ''}
      </div>
    `;

    const input = this._root.querySelector('input');
    if (input) input.addEventListener('input', this._onInput);
    const clear = this._root.querySelector('.clear-btn');
    if (clear) clear.addEventListener('click', this._onClear);
  }
}

// =============================================================================
// 6. <i56-select>
// =============================================================================

class I56Select extends HTMLElement {
  static get observedAttributes() { return ['placeholder', 'value', 'disabled', 'error', 'options']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this._open = false;
    this._searchText = '';
    this._selectedValue = null;
    this._selectedLabel = '';
    this.render();
  }

  attributeChangedCallback() { this.render(); }

  get value() { return this._selectedValue; }
  set value(v) { this._selectedValue = v; this.setAttribute('value', v); this.render(); }

  get _options() {
    try { return JSON.parse(this.getAttribute('options') || '[]'); } catch { return []; }
  }

  connectedCallback() {
    this._selectedValue = this.getAttribute('value') || null;
    document.addEventListener('click', this._onOutsideClick);
    this.render();
  }

  disconnectedCallback() {
    document.removeEventListener('click', this._onOutsideClick);
  }

  _onOutsideClick = (e) => {
    if (!this.contains(e.target)) this._close();
  };

  _toggle = () => {
    if (this.hasAttribute('disabled')) return;
    this._open ? this._close() : this._openDropdown();
  };

  _openDropdown() {
    this._open = true;
    this._searchText = '';
    this.render();
    const input = this._root.querySelector('.search-input');
    if (input) setTimeout(() => input.focus(), 50);
  }

  _close() {
    this._open = false;
    this._searchText = '';
    this.render();
  }

  _select(opt) {
    this._selectedValue = opt.value;
    this._selectedLabel = opt.label;
    this.setAttribute('value', opt.value);
    this._close();
    emit(this, 'change', { value: opt.value, label: opt.label });
  }

  _onKeydown = (e) => {
    if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); this._toggle(); return; }
    if (e.key === 'ArrowDown' && this._open) { e.preventDefault(); this._moveFocus(1); return; }
    if (e.key === 'ArrowUp' && this._open) { e.preventDefault(); this._moveFocus(-1); return; }
    if (e.key === 'Escape') { this._close(); return; }
  };

  _moveFocus(dir) {
    const items = this._root.querySelectorAll('.option:not(.hidden)');
    if (items.length === 0) return;
    const focused = this._root.querySelector('.option.focused');
    let idx = focused ? [...items].indexOf(focused) + dir : 0;
    if (idx < 0) idx = items.length - 1;
    if (idx >= items.length) idx = 0;
    items.forEach(i => i.classList.remove('focused'));
    items[idx].classList.add('focused');
    items[idx].scrollIntoView({ block: 'nearest' });
  }

  _onSearchInput = (e) => {
    this._searchText = e.target.value.toLowerCase();
    this._filterOptions();
  };

  _filterOptions() {
    const items = this._root.querySelectorAll('.option');
    items.forEach(item => {
      const label = (item.dataset.label || '').toLowerCase();
      item.classList.toggle('hidden', this._searchText && !label.includes(this._searchText));
    });
  }

  render() {
    const placeholder = this.getAttribute('placeholder') || 'Select...';
    const disabled = this.hasAttribute('disabled');
    const error = this.hasAttribute('error');
    const options = this._options;

    // Resolve label for current value
    if (this._selectedValue && !this._selectedLabel) {
      const found = options.find(o => o.value === this._selectedValue);
      if (found) this._selectedLabel = found.label;
    }

    const displayText = this._selectedLabel || placeholder;
    const isPlaceholder = !this._selectedLabel;

    this._root.innerHTML = `
      <style>
        :host { display: block; font-family: var(--i56-font-family); position: relative; }
        .trigger {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 0.5rem 0.75rem;
          border: 1px solid ${error ? 'var(--i56-color-danger)' : 'var(--i56-color-border)'};
          border-radius: var(--i56-radius);
          background: var(--i56-color-bg);
          cursor: pointer;
          font-size: var(--i56-font-size-sm);
          color: ${isPlaceholder ? 'var(--i56-color-text-tertiary)' : 'var(--i56-color-text)'};
          min-height: 2.25rem;
          transition: border-color var(--i56-transition), box-shadow var(--i56-transition);
          user-select: none;
        }
        .trigger:focus-visible { outline: none; box-shadow: 0 0 0 3px var(--i56-color-brand-light); border-color: var(--i56-color-brand); }
        :host([disabled]) .trigger { opacity: 0.5; background: var(--i56-color-bg-secondary); cursor: not-allowed; }
        .trigger.open { border-color: var(--i56-color-brand); box-shadow: 0 0 0 3px var(--i56-color-brand-light); }
        .arrow { transition: transform var(--i56-transition); color: var(--i56-color-text-tertiary); font-size: 0.625rem; }
        .arrow.open { transform: rotate(180deg); }
        .dropdown {
          display: ${this._open ? 'block' : 'none'};
          position: absolute;
          top: 100%;
          left: 0;
          right: 0;
          margin-top: 0.25rem;
          background: var(--i56-color-bg);
          border: 1px solid var(--i56-color-border);
          border-radius: var(--i56-radius);
          box-shadow: var(--i56-shadow-lg);
          z-index: 100;
          overflow: hidden;
        }
        .search-wrapper { padding: 0.5rem; border-bottom: 1px solid var(--i56-color-border); }
        .search-input {
          width: 100%;
          box-sizing: border-box;
          padding: 0.375rem 0.5rem;
          border: 1px solid var(--i56-color-border);
          border-radius: var(--i56-radius-sm);
          font-size: var(--i56-font-size-sm);
          font-family: var(--i56-font-family);
          background: var(--i56-color-bg);
          color: var(--i56-color-text);
          outline: none;
        }
        .search-input:focus { border-color: var(--i56-color-brand); }
        .options-list { max-height: 15rem; overflow-y: auto; }
        .option {
          padding: 0.5rem 0.75rem;
          cursor: pointer;
          font-size: var(--i56-font-size-sm);
          color: var(--i56-color-text);
        }
        .option:hover, .option.focused { background: var(--i56-color-brand-light); color: var(--i56-color-brand); }
        .option.hidden { display: none; }
        .option.selected { font-weight: 600; color: var(--i56-color-brand); }
        .empty { padding: 1rem; text-align: center; color: var(--i56-color-text-tertiary); font-size: var(--i56-font-size-sm); }
      </style>
      <div class="trigger ${this._open ? 'open' : ''}" tabindex="0" role="combobox" aria-expanded="${this._open}" aria-haspopup="listbox">
        <span>${displayText}</span>
        <span class="arrow ${this._open ? 'open' : ''}" aria-hidden="true">▼</span>
      </div>
      <div class="dropdown" role="listbox">
        <div class="search-wrapper"><input class="search-input" type="text" placeholder="Search..." /></div>
        <div class="options-list">
          ${options.length === 0 ? '<div class="empty">No options</div>' : options.map(opt => `
            <div class="option${opt.value === this._selectedValue ? ' selected' : ''}" data-value="${opt.value}" data-label="${opt.label}" role="option" aria-selected="${opt.value === this._selectedValue}">${opt.label}</div>
          `).join('')}
        </div>
      </div>
    `;

    // Events
    const trigger = this._root.querySelector('.trigger');
    trigger.addEventListener('click', this._toggle);
    trigger.addEventListener('keydown', this._onKeydown);

    this._root.querySelectorAll('.option').forEach(opt => {
      opt.addEventListener('click', () => {
        this._select({ value: opt.dataset.value, label: opt.dataset.label });
      });
    });
    const search = this._root.querySelector('.search-input');
    if (search) search.addEventListener('input', this._onSearchInput);
  }
}

// =============================================================================
// 7. <i56-modal>
// =============================================================================

class I56Modal extends HTMLElement {
  static get observedAttributes() { return ['open', 'size', 'title', 'close-button']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this.render();
  }

  attributeChangedCallback(name, oldVal, newVal) {
    if (name === 'open') {
      if (newVal !== null) this._openModal();
      else this._closeModal();
    }
    this.render();
  }

  connectedCallback() {
    if (this.hasAttribute('open')) this._openModal();
  }

  disconnectedCallback() {
    this._closeModal();
  }

  _openModal() {
    document.body.style.overflow = 'hidden';
    const backdrop = this._root.querySelector('.backdrop');
    if (backdrop) backdrop.style.display = 'flex';
    // Focus trap
    setTimeout(() => {
      const close = this._root.querySelector('.close-btn');
      if (close) close.focus();
    }, 100);
    emit(this, 'open');
  }

  _closeModal() {
    document.body.style.overflow = '';
    const backdrop = this._root.querySelector('.backdrop');
    if (backdrop) backdrop.style.display = 'none';
    emit(this, 'close');
  }

  _onBackdropClick = (e) => {
    if (e.target.classList.contains('backdrop')) {
      this.removeAttribute('open');
      this._closeModal();
    }
  };

  _onKeydown = (e) => {
    if (e.key === 'Escape') {
      this.removeAttribute('open');
      this._closeModal();
    }
  };

  show() { this.setAttribute('open', ''); }
  hide() { this.removeAttribute('open'); }

  render() {
    const size = this.getAttribute('size') || 'md';
    const title = this.getAttribute('title') || '';
    const closeButton = this.getAttribute('close-button') !== 'false';
    const isOpen = this.hasAttribute('open');

    const sizeMap = { sm: '24rem', md: '32rem', lg: '40rem', xl: '56rem' };
    const maxWidth = sizeMap[size] || sizeMap.md;

    this._root.innerHTML = `
      <style>
        .backdrop {
          display: ${isOpen ? 'flex' : 'none'};
          position: fixed;
          inset: 0;
          z-index: 1000;
          align-items: flex-start;
          justify-content: center;
          padding: 2rem 1rem;
          background: rgba(0, 0, 0, 0.5);
          backdrop-filter: blur(4px);
          animation: i56-fade-in 150ms ease;
        }
        .dialog {
          background: var(--i56-color-bg);
          border-radius: var(--i56-radius-lg);
          box-shadow: var(--i56-shadow-xl);
          width: 100%;
          max-width: ${maxWidth};
          max-height: calc(100vh - 4rem);
          display: flex;
          flex-direction: column;
          animation: i56-scale-in 150ms ease;
          font-family: var(--i56-font-family);
        }
        .header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 1rem 1.25rem;
          border-bottom: 1px solid var(--i56-color-border);
        }
        .title {
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
        .body {
          flex: 1;
          overflow-y: auto;
          padding: 1.25rem;
        }
        @keyframes i56-fade-in { from { opacity: 0; } to { opacity: 1; } }
        @keyframes i56-scale-in { from { opacity: 0; transform: scale(0.95) translateY(0.5rem); } to { opacity: 1; transform: scale(1) translateY(0); } }
      </style>
      <div class="backdrop" role="dialog" aria-modal="true" aria-label="${title}">
        <div class="dialog">
          ${title || closeButton ? `
            <div class="header">
              <span class="title">${title}</span>
              ${closeButton ? '<button class="close-btn" aria-label="Close dialog">&times;</button>' : ''}
            </div>
          ` : ''}
          <div class="body"><slot></slot></div>
        </div>
      </div>
    `;

    if (isOpen) {
      this._root.querySelector('.backdrop').addEventListener('click', this._onBackdropClick);
      document.addEventListener('keydown', this._onKeydown);
      const close = this._root.querySelector('.close-btn');
      if (close) close.addEventListener('click', () => { this.removeAttribute('open'); this._closeModal(); });
    }
  }
}

// =============================================================================
// 8. <i56-toast> — integrated toast system
// =============================================================================

class I56ToastContainer {
  constructor() {
    this._container = null;
    this._toasts = [];
    this._position = 'top-right';
    this._init();
  }

  _init() {
    if (document.getElementById('i56-toast-container')) return;
    this._container = document.createElement('div');
    this._container.id = 'i56-toast-container';
    this._updatePosition();
    document.body.appendChild(this._container);
  }

  _updatePosition() {
    if (!this._container) return;
    const isBottom = this._position === 'bottom-right';
    this._container.style.cssText = `
      position: fixed;
      z-index: 9999;
      ${isBottom ? 'bottom' : 'top'}: 1rem;
      right: 1rem;
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      pointer-events: none;
    `;
  }

  show({ type = 'info', title = '', message = '', duration = 5000, position } = {}) {
    if (position && position !== this._position) {
      this._position = position;
      this._updatePosition();
    }

    const toast = document.createElement('i56-toast');
    toast.setAttribute('type', type);
    toast.setAttribute('title', title);
    toast.setAttribute('message', message);
    toast.setAttribute('duration', duration);
    this._container.appendChild(toast);

    const idx = this._toasts.length;
    this._toasts.push(toast);

    toast.addEventListener('i56:close', () => {
      toast.remove();
      this._toasts = this._toasts.filter(t => t !== toast);
    });

    // Auto-dismiss
    if (duration > 0) {
      setTimeout(() => {
        toast.dismiss();
      }, duration);
    }

    return toast;
  }
}

// Global toast API
window.I56Toast = new I56ToastContainer();

class I56ToastEl extends HTMLElement {
  static get observedAttributes() { return ['type', 'title', 'message']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this.render();
  }

  attributeChangedCallback() { this.render(); }

  connectedCallback() {
    // Trigger enter animation
    requestAnimationFrame(() => {
      const el = this._root.querySelector('.toast');
      if (el) el.classList.add('visible');
    });
  }

  dismiss() {
    const el = this._root.querySelector('.toast');
    if (el) {
      el.classList.remove('visible');
      el.classList.add('hiding');
      setTimeout(() => emit(this, 'close'), 200);
    }
  }

  render() {
    const type = this.getAttribute('type') || 'info';
    const title = this.getAttribute('title') || '';
    const message = this.getAttribute('message') || '';

    const iconMap = { success: '✓', error: '✕', warning: '⚠', info: 'ℹ' };
    const icon = iconMap[type] || 'ℹ';

    this._root.innerHTML = `
      <style>
        .toast {
          display: flex;
          align-items: flex-start;
          gap: 0.75rem;
          padding: 0.75rem 1rem;
          background: var(--i56-color-bg);
          border: 1px solid var(--i56-color-border);
          border-radius: var(--i56-radius);
          box-shadow: var(--i56-shadow-lg);
          min-width: 18rem;
          max-width: 24rem;
          font-family: var(--i56-font-family);
          pointer-events: auto;
          transform: translateX(120%);
          opacity: 0;
          transition: all 200ms ease;
        }
        .toast.visible { transform: translateX(0); opacity: 1; }
        .toast.hiding { transform: translateX(120%); opacity: 0; }
        .icon {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 1.5rem;
          height: 1.5rem;
          border-radius: 50%;
          font-size: var(--i56-font-size-xs);
          font-weight: 700;
          flex-shrink: 0;
          color: #fff;
        }
        .toast.success .icon { background: var(--i56-color-success); }
        .toast.error .icon { background: var(--i56-color-danger); }
        .toast.warning .icon { background: var(--i56-color-warning); }
        .toast.info .icon { background: var(--i56-color-info); }
        .content { flex: 1; min-width: 0; }
        .title { font-size: var(--i56-font-size-sm); font-weight: 600; color: var(--i56-color-text); margin-bottom: 0.125rem; }
        .message { font-size: var(--i56-font-size-xs); color: var(--i56-color-text-secondary); }
        .close {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 1.25rem;
          height: 1.25rem;
          border: none;
          background: none;
          cursor: pointer;
          color: var(--i56-color-text-tertiary);
          font-size: var(--i56-font-size-base);
          line-height: 1;
          border-radius: var(--i56-radius-sm);
          flex-shrink: 0;
          padding: 0;
        }
        .close:hover { background: var(--i56-color-bg-secondary); color: var(--i56-color-text); }
      </style>
      <div class="toast ${type}">
        <span class="icon">${icon}</span>
        <div class="content">
          ${title ? `<div class="title">${title}</div>` : ''}
          <div class="message">${message}</div>
        </div>
        <button class="close" aria-label="Close notification">&times;</button>
      </div>
    `;

    const closeBtn = this._root.querySelector('.close');
    if (closeBtn) closeBtn.addEventListener('click', () => this.dismiss());
  }
}

// =============================================================================
// 9. <i56-badge>
// =============================================================================

class I56Badge extends HTMLElement {
  static get observedAttributes() { return ['color']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this.render();
  }

  attributeChangedCallback() { this.render(); }

  render() {
    const color = this.getAttribute('color') || 'neutral';

    this._root.innerHTML = `
      <style>
        :host {
          display: inline-flex;
          align-items: center;
          padding: 0.125rem 0.5rem;
          font-size: var(--i56-font-size-xs);
          font-weight: 500;
          line-height: 1.5;
          border-radius: var(--i56-radius-full);
          font-family: var(--i56-font-family);
          white-space: nowrap;
          user-select: none;
        }
        :host([color="brand"])  { background: var(--i56-color-brand-light);  color: var(--i56-color-brand); }
        :host([color="success"]){ background: var(--i56-color-success-light); color: var(--i56-color-success); }
        :host([color="warning"]){ background: var(--i56-color-warning-light); color: var(--i56-color-warning); }
        :host([color="danger"]) { background: var(--i56-color-danger-light);  color: var(--i56-color-danger); }
        :host([color="neutral"]){ background: var(--i56-color-neutral-light); color: var(--i56-color-neutral); }
        :host([color="info"])   { background: var(--i56-color-info-light);    color: var(--i56-color-info); }
      </style>
      <slot></slot>
    `;
  }
}

// =============================================================================
// 10. <i56-tabs> and <i56-tab-panel>
// =============================================================================

class I56TabPanel extends HTMLElement {
  static get observedAttributes() { return ['label', 'active', 'disabled']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
  }

  connectedCallback() {
    this.setAttribute('role', 'tabpanel');
    if (!this.hasAttribute('active')) this.style.display = 'none';
  }

  attributeChangedCallback(name, oldVal, newVal) {
    if (name === 'active') {
      this.style.display = newVal !== null ? '' : 'none';
    }
  }
}

class I56Tabs extends HTMLElement {
  static get observedAttributes() { return ['active']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this._activeIndex = 0;
  }

  connectedCallback() {
    this.render();
    this.addEventListener('keydown', this._onKeydown);
  }

  disconnectedCallback() {
    this.removeEventListener('keydown', this._onKeydown);
  }

  _getPanels() {
    return [...this.querySelectorAll('i56-tab-panel')];
  }

  _onTabClick(idx) {
    this._activate(idx);
    emit(this, 'change', { index: idx, panel: this._getPanels()[idx] });
  }

  _activate(idx) {
    const panels = this._getPanels();
    panels.forEach((p, i) => {
      if (i === idx) p.setAttribute('active', '');
      else p.removeAttribute('active');
    });
    this._activeIndex = idx;
    this.setAttribute('active', idx);
    this.render();
  }

  _onKeydown = (e) => {
    const panels = this._getPanels();
    let idx = this._activeIndex;
    if (e.key === 'ArrowRight') { e.preventDefault(); idx = (idx + 1) % panels.length; }
    else if (e.key === 'ArrowLeft') { e.preventDefault(); idx = (idx - 1 + panels.length) % panels.length; }
    else if (e.key === 'Home') { e.preventDefault(); idx = 0; }
    else if (e.key === 'End') { e.preventDefault(); idx = panels.length - 1; }
    else return;
    this._activate(idx);
    emit(this, 'change', { index: idx, panel: panels[idx] });
    const tab = this._root.querySelectorAll('.tab')[idx];
    if (tab) tab.focus();
  };

  render() {
    const panels = this._getPanels();

    this._root.innerHTML = `
      <style>
        :host { display: block; font-family: var(--i56-font-family); }
        .tab-list {
          display: flex;
          border-bottom: 2px solid var(--i56-color-border);
          gap: 0;
          overflow-x: auto;
          scrollbar-width: none;
        }
        .tab-list::-webkit-scrollbar { display: none; }
        .tab {
          position: relative;
          padding: 0.625rem 1rem;
          font-size: var(--i56-font-size-sm);
          font-weight: 500;
          color: var(--i56-color-text-secondary);
          background: none;
          border: none;
          cursor: pointer;
          font-family: var(--i56-font-family);
          white-space: nowrap;
          transition: color var(--i56-transition);
          outline: none;
          flex-shrink: 0;
        }
        .tab:hover { color: var(--i56-color-text); }
        .tab.active { color: var(--i56-color-brand); }
        .tab.active::after {
          content: '';
          position: absolute;
          bottom: -2px;
          left: 0;
          right: 0;
          height: 2px;
          background: var(--i56-color-brand);
          border-radius: 1px 1px 0 0;
        }
        .tab:focus-visible { box-shadow: 0 0 0 2px var(--i56-color-brand-light); border-radius: var(--i56-radius-sm); }
        .tab[disabled] { opacity: 0.4; cursor: not-allowed; pointer-events: none; }
        .panels { padding: 1rem 0; }
      </style>
      <div class="tab-list" role="tablist">
        ${panels.map((p, i) => {
          const label = p.getAttribute('label') || `Tab ${i + 1}`;
          const disabled = p.hasAttribute('disabled');
          const active = i === this._activeIndex;
          return `<button class="tab${active ? ' active' : ''}" role="tab" aria-selected="${active}"
            aria-controls="panel-${i}" tabindex="${active ? '0' : '-1'}"
            ?disabled="${disabled}" data-index="${i}">${label}</button>`;
        }).join('')}
      </div>
      <div class="panels">
        ${panels.map((p, i) => `<div role="tabpanel" id="panel-${i}" aria-labelledby="tab-${i}" style="display:${i === this._activeIndex ? '' : 'none'}"><slot name="panel-${i}"></slot></div>`).join('')}
      </div>
    `;

    this._root.querySelectorAll('.tab').forEach(tab => {
      tab.addEventListener('click', () => {
        if (!tab.disabled) this._onTabClick(parseInt(tab.dataset.index));
      });
    });

    // Sync panel slots — render panels inside their tabpanel containers
    panels.forEach((panel, i) => {
      const slotName = `panel-${i}`;
      if (!panel.hasAttribute('slot')) panel.setAttribute('slot', slotName);
    });
  }
}

// =============================================================================
// 11. <i56-avatar>
// =============================================================================

class I56Avatar extends HTMLElement {
  static get observedAttributes() { return ['src', 'name', 'size', 'alt']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
  }

  attributeChangedCallback() { this.render(); }
  connectedCallback() { this.render(); }

  _getInitials(name) {
    if (!name) return '?';
    const parts = name.trim().split(/\s+/);
    return parts.length > 1 ? (parts[0][0] + parts[parts.length - 1][0]).toUpperCase() : parts[0].substring(0, 2).toUpperCase();
  }

  _getColor(name) {
    const colors = [
      '#4F46E5', '#059669', '#D97706', '#DC2626', '#2563EB',
      '#7C3AED', '#DB2777', '#0891B2', '#65A30D', '#9333EA'
    ];
    let hash = 0;
    for (let i = 0; i < (name || '').length; i++) hash = (hash * 31 + name.charCodeAt(i)) & 0xffffffff;
    return colors[Math.abs(hash) % colors.length];
  }

  render() {
    const src = this.getAttribute('src') || '';
    const name = this.getAttribute('name') || '';
    const size = this.getAttribute('size') || 'md';
    const alt = this.getAttribute('alt') || name || 'Avatar';
    const initials = this._getInitials(name);
    const bgColor = this._getColor(name);

    const sizeMap = { xs: '1.5rem', sm: '2rem', md: '2.5rem', lg: '3rem', xl: '4rem' };
    const dim = sizeMap[size] || sizeMap.md;
    const fontSizeMap = { xs: '0.625rem', sm: '0.75rem', md: '0.875rem', lg: '1rem', xl: '1.25rem' };

    this._root.innerHTML = `
      <style>
        :host {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          width: ${dim};
          height: ${dim};
          border-radius: 50%;
          overflow: hidden;
          font-family: var(--i56-font-family);
          font-weight: 600;
          font-size: ${fontSizeMap[size] || '0.875rem'};
          color: #fff;
          background: ${bgColor};
          flex-shrink: 0;
          user-select: none;
        }
        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          display: ${src ? 'block' : 'none'};
        }
        .initials { display: ${src ? 'none' : 'flex'}; align-items: center; justify-content: center; width: 100%; height: 100%; }
        img[src=""] { display: none; }
      </style>
      ${src ? `<img src="${src}" alt="${alt}" onerror="this.style.display='none';this.nextElementSibling.style.display='flex'" />` : ''}
      <span class="initials" aria-hidden="true">${initials}</span>
    `;
  }
}

// =============================================================================
// 12. <i56-spinner>
// =============================================================================

class I56Spinner extends HTMLElement {
  static get observedAttributes() { return ['size']; }

  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this.render();
  }

  attributeChangedCallback() { this.render(); }

  render() {
    const size = this.getAttribute('size') || 'md';
    const sizeMap = { sm: '1rem', md: '1.5rem', lg: '2.5rem', xl: '4rem' };
    const dim = sizeMap[size] || sizeMap.md;

    this._root.innerHTML = `
      <style>
        :host {
          display: inline-block;
          width: ${dim};
          height: ${dim};
          flex-shrink: 0;
        }
        .spinner {
          width: 100%;
          height: 100%;
          border: 2px solid var(--i56-color-border);
          border-top-color: var(--i56-color-brand);
          border-radius: 50%;
          animation: i56-spin 0.7s linear infinite;
          box-sizing: border-box;
        }
        @keyframes i56-spin {
          to { transform: rotate(360deg); }
        }
      </style>
      <div class="spinner" role="status" aria-label="Loading"></div>
    `;
  }
}

// =============================================================================
// 13. <i56-timeline> and <i56-timeline-item>
// =============================================================================

class I56TimelineItem extends HTMLElement {
  static get observedAttributes() { return ['status', 'time', 'title']; }
  // status: completed, current, pending
}

class I56Timeline extends HTMLElement {
  constructor() {
    super();
    this._root = this.attachShadow({ mode: 'open' });
    this._root.adoptedStyleSheets = [getSharedSheet()];
    this.render();
  }

  connectedCallback() { this.render(); }

  render() {
    const items = [...this.querySelectorAll('i56-timeline-item')];

    this._root.innerHTML = `
      <style>
        :host { display: block; font-family: var(--i56-font-family); }
        .timeline {
          position: relative;
          padding-left: 1.75rem;
        }
        .item {
          position: relative;
          padding-bottom: 1.5rem;
        }
        .item:last-child { padding-bottom: 0; }
        .item::before {
          content: '';
          position: absolute;
          left: -1.25rem;
          top: 0.5rem;
          width: 0.625rem;
          height: 0.625rem;
          border-radius: 50%;
          background: var(--i56-color-border);
          z-index: 1;
        }
        .item.completed::before { background: var(--i56-color-success); }
        .item.current::before { background: var(--i56-color-brand); box-shadow: 0 0 0 3px var(--i56-color-brand-light); }
        .item.pending::before { background: var(--i56-color-border); }
        .item::after {
          content: '';
          position: absolute;
          left: -1rem;
          top: 0.75rem;
          bottom: 0;
          width: 1px;
          background: var(--i56-color-border);
        }
        .item:last-child::after { display: none; }
        .item.completed + .item::after { background: var(--i56-color-success); }
        .time {
          font-size: var(--i56-font-size-xs);
          color: var(--i56-color-text-tertiary);
          margin-bottom: 0.125rem;
        }
        .item-title {
          font-size: var(--i56-font-size-sm);
          font-weight: 500;
          color: var(--i56-color-text);
          margin-bottom: 0.125rem;
        }
        .item-content {
          font-size: var(--i56-font-size-sm);
          color: var(--i56-color-text-secondary);
          line-height: var(--i56-line-height);
        }
      </style>
      <div class="timeline">
        ${items.length === 0 ? '<div style="color:var(--i56-color-text-tertiary);font-size:var(--i56-font-size-sm)">No timeline items</div>' : ''}
        ${items.map((item, i) => {
          const status = item.getAttribute('status') || 'pending';
          const time = item.getAttribute('time') || '';
          const title = item.getAttribute('title') || '';
          return `
            <div class="item ${status}">
              ${time ? `<div class="time">${time}</div>` : ''}
              ${title ? `<div class="item-title">${title}</div>` : ''}
              <div class="item-content"><slot name="item-${i}"></slot></div>
            </div>
          `;
        }).join('')}
      </div>
    `;

    items.forEach((item, i) => {
      if (!item.hasAttribute('slot')) item.setAttribute('slot', `item-${i}`);
    });
  }
}

// =============================================================================
// Registration
// =============================================================================

customElements.define('i56-button', I56Button);
customElements.define('i56-card', I56Card);
customElements.define('i56-table', I56Table);
customElements.define('i56-form-group', I56FormGroup);
customElements.define('i56-input', I56Input);
customElements.define('i56-select', I56Select);
customElements.define('i56-modal', I56Modal);
customElements.define('i56-toast', I56ToastEl);
customElements.define('i56-badge', I56Badge);
customElements.define('i56-tabs', I56Tabs);
customElements.define('i56-tab-panel', I56TabPanel);
customElements.define('i56-avatar', I56Avatar);
customElements.define('i56-spinner', I56Spinner);
customElements.define('i56-timeline', I56Timeline);
customElements.define('i56-timeline-item', I56TimelineItem);
