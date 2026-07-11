/**
 * I56 AI Security Guardrail — Security indicator with PII masking, prompt injection,
 * and RBAC guard status display.
 *
 * BDL 2.0 AI Web Component: native Custom Element, Shadow DOM, zero dependencies.
 *
 * Usage:
 *   <i56-ai-security-guardrail
 *     pii-masking="enabled"
 *     prompt-injection="active"
 *     rbac-guard="restricted"
 *     data-policy='{"level":"strict","piiFields":["email","phone","ssn"]}'>
 *   </i56-ai-security-guardrail>
 *
 * Attributes:
 *   pii-masking: enabled | disabled | partial
 *   prompt-injection: active | watching | blocked | idle
 *   rbac-guard: full | restricted | read-only | none
 *   data-policy: JSON with policy configuration
 *   compact: present for a single-line compact display
 *   expanded: present to show details
 *
 * Events:
 *   i56:guardrail-click — detail: { guard, status }
 *   i56:policy-violation — detail: { type, details }
 *   i56:guardrail-toggle — detail: { expanded }
 *
 * Keyboard:
 *   Enter/Space — toggle details
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
          --i56-color-brand-light: var(--i56-brand-light, #DBEAFE);
          --i56-color-success: var(--i56-success, #059669);
          --i56-color-success-light: var(--i56-success-light, #ECFDF5);
          --i56-color-warning: var(--i56-warning, #D97706);
          --i56-color-warning-light: var(--i56-warning-light, #FFFBEB);
          --i56-color-danger: var(--i56-danger, #DC2626);
          --i56-color-danger-light: var(--i56-danger-light, #FEF2F2);
          --i56-color-info: var(--i56-info, #2563EB);
          --i56-color-info-light: var(--i56-info-light, #EFF6FF);
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
          --i56-radius-full: 9999px;
          --i56-shadow-sm: 0 1px 2px 0 rgba(0,0,0,0.05);
          --i56-font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
          --i56-font-size-xs: 0.75rem;
          --i56-font-size-sm: 0.875rem;
          --i56-font-size-base: 1rem;
          --i56-transition: 150ms ease;
          --i56-line-height: 1.5;
        }
      `);
    }
    return _sharedSheet;
  }

  function emit(el, name, detail = {}) {
    el.dispatchEvent(new CustomEvent(`i56:${name}`, { bubbles: true, composed: true, detail }));
  }

  class I56AiSecurityGuardrail extends HTMLElement {
    static get observedAttributes() {
      return ['pii-masking', 'prompt-injection', 'rbac-guard', 'data-policy', 'compact', 'expanded'];
    }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
    }

    get piiMasking() { return this.getAttribute('pii-masking') || 'disabled'; }
    get promptInjection() { return this.getAttribute('prompt-injection') || 'idle'; }
    get rbacGuard() { return this.getAttribute('rbac-guard') || 'none'; }
    get compact() { return this.hasAttribute('compact'); }
    get expanded() { return this.hasAttribute('expanded'); }

    get dataPolicy() {
      try { return JSON.parse(this.getAttribute('data-policy') || '{}'); } catch { return {}; }
    }
    set dataPolicy(obj) {
      this.setAttribute('data-policy', JSON.stringify(obj));
    }

    connectedCallback() { this.render(); }
    attributeChangedCallback() { this.render(); }

    // -- Public API --
    toggle() {
      this.expanded ? this.removeAttribute('expanded') : this.setAttribute('expanded', '');
      emit(this, 'guardrail-toggle', { expanded: this.hasAttribute('expanded') });
    }

    reportViolation(type, details = {}) {
      emit(this, 'policy-violation', { type, details });
    }

    // -- Internals --
    _statusIndicator(status) {
      const map = {
        enabled: { icon: '✅', class: 'status-success', label: 'Active' },
        active: { icon: '🛡️', class: 'status-success', label: 'Active' },
        full: { icon: '🔐', class: 'status-success', label: 'Full' },
        disabled: { icon: '⭕', class: 'status-idle', label: 'Disabled' },
        idle: { icon: '⭕', class: 'status-idle', label: 'Idle' },
        none: { icon: '⭕', class: 'status-idle', label: 'None' },
        partial: { icon: '⚠️', class: 'status-warning', label: 'Partial' },
        watching: { icon: '👁️', class: 'status-warning', label: 'Watching' },
        blocked: { icon: '🚫', class: 'status-danger', label: 'Blocked' },
        restricted: { icon: '🔒', class: 'status-warning', label: 'Restricted' },
        'read-only': { icon: '📖', class: 'status-info', label: 'Read-only' },
      };
      return map[status] || map.idle;
    }

    _renderGuard(name, status, label) {
      const indicator = this._statusIndicator(status);
      return `
        <div class="guard-item" tabindex="0" role="button"
             aria-label="${label}: ${indicator.label}"
             data-guard="${name}">
          <span class="guard-icon" aria-hidden="true">${indicator.icon}</span>
          <span class="guard-label">${label}</span>
          <span class="guard-status ${indicator.class}">${indicator.label}</span>
        </div>
      `;
    }

    render() {
      const isCompact = this.compact;
      const isExpanded = this.expanded;
      const policy = this.dataPolicy;

      const guards = [
        { name: 'pii', status: this.piiMasking, label: 'PII Masking' },
        { name: 'injection', status: this.promptInjection, label: 'Prompt Guard' },
        { name: 'rbac', status: this.rbacGuard, label: 'RBAC' },
      ];

      this._root.innerHTML = `
        <style>
          :host {
            display: block;
            font-family: var(--i56-font-family);
          }

          .guardrail-container {
            background: var(--i56-color-bg-surface);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-lg);
            overflow: hidden;
            transition: all var(--i56-transition);
          }

          .guardrail-bar {
            display: flex;
            align-items: center;
            gap: 0.75rem;
            padding: ${isCompact ? '0.375rem 0.75rem' : '0.5rem 1rem'};
            cursor: pointer;
            transition: background var(--i56-transition);
          }
          .guardrail-bar:hover {
            background: var(--i56-color-bg-secondary);
          }
          .guardrail-bar:focus-visible {
            box-shadow: 0 0 0 2px var(--i56-color-brand-light);
            outline: none;
            border-radius: var(--i56-radius-lg);
          }

          .guardrail-icon {
            font-size: ${isCompact ? '1rem' : '1.25rem'};
            flex-shrink: 0;
          }

          .guardrail-title {
            font-weight: 600;
            font-size: var(--i56-font-size-sm);
            color: var(--i56-color-text);
            flex: 1;
            min-width: 0;
          }

          /* Guard indicators row */
          .guard-indicators {
            display: flex;
            align-items: center;
            gap: 0.5rem;
          }

          .guard-dot {
            width: 8px;
            height: 8px;
            border-radius: 50%;
          }
          .guard-dot.success { background: var(--i56-color-success); }
          .guard-dot.warning { background: var(--i56-color-warning); }
          .guard-dot.danger { background: var(--i56-color-danger); }
          .guard-dot.info { background: var(--i56-color-info); }
          .guard-dot.idle { background: var(--i56-color-text-tertiary); }

          .guard-summary {
            display: flex;
            align-items: center;
            gap: 0.25rem;
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-secondary);
          }

          /* Detail panel */
          .guard-details {
            display: ${isExpanded ? 'block' : 'none'};
            border-top: 1px solid var(--i56-color-border);
            padding: 0.75rem 1rem;
            animation: i56-guard-details-in 200ms ease;
          }

          @keyframes i56-guard-details-in {
            from { opacity: 0; max-height: 0; }
            to { opacity: 1; max-height: 300px; }
          }

          .guard-item {
            display: flex;
            align-items: center;
            gap: 0.625rem;
            padding: 0.5rem 0.75rem;
            border-radius: var(--i56-radius);
            cursor: pointer;
            transition: background var(--i56-transition);
            outline: none;
          }
          .guard-item:hover {
            background: var(--i56-color-bg-secondary);
          }
          .guard-item:focus-visible {
            box-shadow: 0 0 0 2px var(--i56-color-brand-light);
            outline: none;
          }
          .guard-item + .guard-item {
            border-top: 1px solid var(--i56-color-border);
          }

          .guard-icon {
            font-size: var(--i56-font-size-sm);
            flex-shrink: 0;
            width: 1.25rem;
            text-align: center;
          }
          .guard-label {
            font-size: var(--i56-font-size-xs);
            font-weight: 500;
            color: var(--i56-color-text-secondary);
            flex: 1;
            min-width: 0;
          }
          .guard-status {
            font-size: var(--i56-font-size-xs);
            font-weight: 600;
            padding: 0.125rem 0.5rem;
            border-radius: var(--i56-radius-full);
            flex-shrink: 0;
          }

          .status-success { color: #065F46; background: #ECFDF5; }
          .status-warning { color: #92400E; background: #FFFBEB; }
          .status-danger { color: #991B1B; background: #FEF2F2; }
          .status-info { color: #1E40AF; background: #EFF6FF; }
          .status-idle { color: var(--i56-color-text-tertiary); background: var(--i56-color-bg-secondary); }

          /* Policy info */
          .policy-info {
            margin-top: 0.75rem;
            padding: 0.5rem 0.75rem;
            border-radius: var(--i56-radius);
            background: var(--i56-color-bg-secondary);
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-secondary);
            display: ${policy.level ? 'block' : 'none'};
          }
          .policy-level {
            font-weight: 600;
            color: var(--i56-color-text);
            text-transform: uppercase;
            letter-spacing: 0.05em;
          }
          .policy-fields {
            margin-top: 0.25rem;
            display: flex;
            flex-wrap: wrap;
            gap: 0.25rem;
          }
          .policy-field {
            padding: 0.125rem 0.5rem;
            border-radius: var(--i56-radius-sm);
            background: var(--i56-color-bg-surface);
            border: 1px solid var(--i56-color-border);
            font-family: 'SF Mono', 'Fira Code', monospace;
            font-size: 0.6875rem;
          }
        </style>

        <div class="guardrail-container">
          <div class="guardrail-bar" tabindex="0" role="button"
               aria-label="Security Guardrail - ${guards.filter(g => g.status !== 'disabled' && g.status !== 'idle' && g.status !== 'none').length} active"
               aria-expanded="${isExpanded}">
            <span class="guardrail-icon" aria-hidden="true">🛡️</span>
            <span class="guardrail-title">Security Guard</span>
            <span class="guard-summary">
              ${guards.map(g => {
                const dot = this._statusIndicator(g.status);
                const dotClass = dot.class.replace('status-', '');
                return `<span class="guard-dot ${dotClass}" title="${g.label}: ${dot.label}" aria-hidden="true"></span>`;
              }).join('')}
            </span>
          </div>

          <div class="guard-details" role="region" aria-label="Guard details">
            ${guards.map(g => this._renderGuard(g.name, g.status, g.label)).join('')}

            ${policy.level ? `
              <div class="policy-info">
                <div class="policy-level">Policy: ${policy.level}</div>
                ${policy.piiFields && policy.piiFields.length ? `
                  <div class="policy-fields">
                    ${policy.piiFields.map(f => `<span class="policy-field">${f}</span>`).join('')}
                  </div>
                ` : ''}
              </div>
            ` : ''}
          </div>
        </div>
      `;

      // Wire toggle
      const bar = this._root.querySelector('.guardrail-bar');
      if (bar) {
        bar.addEventListener('click', () => this.toggle());
        bar.addEventListener('keydown', (e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            this.toggle();
          }
        });
      }

      // Wire guard items
      this._root.querySelectorAll('.guard-item').forEach(item => {
        item.addEventListener('click', (e) => {
          e.stopPropagation();
          emit(this, 'guardrail-click', {
            guard: item.dataset.guard,
            status: this.getAttribute(
              item.dataset.guard === 'pii' ? 'pii-masking' :
              item.dataset.guard === 'injection' ? 'prompt-injection' : 'rbac-guard'
            ) || 'idle',
          });
        });
      });
    }
  }

  customElements.define('i56-ai-security-guardrail', I56AiSecurityGuardrail);
})();
