/**
 * I56 AI Agent Monitor — Agent runtime dashboard with grid card layout.
 *
 * BDL 2.0 AI Web Component: native Custom Element, Shadow DOM, zero dependencies.
 * Shows agent name, status, current task, last action time in a card grid.
 *
 * Usage:
 *   <i56-ai-agent-monitor agents='[
 *     {"name":"Analyst","status":"active","task":"Scanning orders","lastAction":"2s ago"},
 *     {"name":"Coder","status":"idle","task":"-","lastAction":"5m ago"},
 *     {"name":"Reviewer","status":"error","task":"Failed lint check","lastAction":"30s ago"}
 *   ]' poll-interval="5000">
 *   </i56-ai-agent-monitor>
 *
 * Attributes:
 *   agents: JSON array of {name, status, task?, lastAction?, icon?}
 *   poll-interval: auto-refresh interval in ms (0 = disabled)
 *   compact: present for compact mode
 *   empty-message: message when no agents (default: "No agents running")
 *
 * Properties:
 *   agents: get/set agents array
 *
 * Events:
 *   i56:agent-click — detail: { agent, index }
 *   i56:agent-status-change — detail: { agent, previousStatus }
 *   i56:poll — detail: { agents }
 *
 * Keyboard:
 *   Tab between agent cards, Enter/Space to select
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
          --i56-radius-xl: 16px;
          --i56-radius-full: 9999px;
          --i56-shadow-sm: 0 1px 2px 0 rgba(0,0,0,0.05);
          --i56-shadow: 0 1px 3px 0 rgba(0,0,0,0.06), 0 1px 2px -1px rgba(0,0,0,0.06);
          --i56-shadow-md: 0 4px 6px -1px rgba(0,0,0,0.06), 0 2px 4px -2px rgba(0,0,0,0.06);
          --i56-font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
          --i56-font-size-xs: 0.75rem;
          --i56-font-size-sm: 0.875rem;
          --i56-font-size-base: 1rem;
          --i56-font-size-lg: 1.125rem;
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

  class I56AiAgentMonitor extends HTMLElement {
    static get observedAttributes() {
      return ['agents', 'poll-interval', 'compact', 'empty-message'];
    }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
      this._pollTimer = null;
      this._prevStatuses = {};
    }

    get agents() {
      try { return JSON.parse(this.getAttribute('agents') || '[]'); } catch { return []; }
    }
    set agents(arr) {
      // Track status changes
      const prev = this.agents;
      prev.forEach(a => { this._prevStatuses[a.name] = a.status; });
      arr.forEach(a => {
        const old = this._prevStatuses[a.name];
        if (old && old !== a.status) {
          emit(this, 'agent-status-change', { agent: a, previousStatus: old });
        }
      });
      this.setAttribute('agents', JSON.stringify(arr));
    }

    get pollInterval() {
      const v = parseInt(this.getAttribute('poll-interval'));
      return v > 0 ? v : 0;
    }
    get compact() { return this.hasAttribute('compact'); }
    get emptyMessage() { return this.getAttribute('empty-message') || 'No agents running'; }

    connectedCallback() {
      this.render();
      this._startPolling();
    }

    disconnectedCallback() {
      this._stopPolling();
    }

    attributeChangedCallback(name) {
      if (name === 'poll-interval') {
        this._stopPolling();
        this._startPolling();
      }
      this.render();
    }

    // -- Public API --
    updateAgent(name, update) {
      const agents = this.agents;
      const idx = agents.findIndex(a => a.name === name);
      if (idx >= 0) {
        agents[idx] = { ...agents[idx], ...update, lastAction: 'just now' };
        this.agents = agents;
      }
    }

    refresh() {
      emit(this, 'poll', { agents: this.agents });
    }

    // -- Internals --
    _startPolling() {
      const interval = this.pollInterval;
      if (interval > 0) {
        this._pollTimer = setInterval(() => this.refresh(), interval);
      }
    }

    _stopPolling() {
      if (this._pollTimer) {
        clearInterval(this._pollTimer);
        this._pollTimer = null;
      }
    }

    _onAgentClick(agent, index) {
      emit(this, 'agent-click', { agent, index });
    }

    _statusIcon(status) {
      switch (status) {
        case 'active': return '<span class="pulse-dot active"></span>';
        case 'busy': return '<span class="pulse-dot busy"></span>';
        case 'error': return '<span class="pulse-dot error"></span>';
        case 'completed': return '<span class="pulse-dot completed"></span>';
        default: return '<span class="pulse-dot idle"></span>';
      }
    }

    _statusLabel(status) {
      return status ? status.charAt(0).toUpperCase() + status.slice(1) : 'Idle';
    }

    _renderAgentCard(agent, index) {
      const status = agent.status || 'idle';
      const icon = this._statusIcon(status);
      const task = agent.task || '-';
      const lastAction = agent.lastAction || '-';
      const agentIcon = agent.icon || '🤖';

      return `
        <div class="agent-card" tabindex="0" role="button"
             aria-label="${agent.name}: ${this._statusLabel(status)}, ${task}"
             data-index="${index}">
          <div class="card-header">
            <span class="agent-icon" aria-hidden="true">${agentIcon}</span>
            <span class="agent-name">${agent.name}</span>
            ${icon}
          </div>
          <div class="card-body">
            <div class="agent-status-label status-${status}">${this._statusLabel(status)}</div>
            <div class="agent-task" title="${task}">${task}</div>
          </div>
          <div class="card-footer">
            <span class="agent-last-action" title="Last action: ${lastAction}">
              ${lastAction !== '-' ? '🕐 ' + lastAction : '-'}
            </span>
          </div>
        </div>
      `;
    }

    render() {
      const agents = this.agents;
      const compact = this.compact;

      const cardsHtml = agents.length === 0
        ? `<div class="empty-state">${this.emptyMessage}</div>`
        : agents.map((a, i) => this._renderAgentCard(a, i)).join('');

      this._root.innerHTML = `
        <style>
          :host {
            display: block;
            font-family: var(--i56-font-family);
          }

          .agent-grid {
            display: grid;
            grid-template-columns: ${compact ? 'repeat(auto-fill, minmax(160px, 1fr))' : 'repeat(auto-fill, minmax(200px, 1fr))'};
            gap: ${compact ? '0.5rem' : '0.75rem'};
          }

          .empty-state {
            text-align: center;
            padding: 1.5rem;
            color: var(--i56-color-text-tertiary);
            font-size: var(--i56-font-size-sm);
          }

          .agent-card {
            display: flex;
            flex-direction: column;
            background: var(--i56-color-bg-surface);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-lg);
            padding: ${compact ? '0.5rem 0.625rem' : '0.75rem 0.875rem'};
            cursor: pointer;
            transition: all var(--i56-transition);
            outline: none;
          }
          .agent-card:hover {
            border-color: var(--i56-color-brand);
            box-shadow: var(--i56-shadow-md);
            transform: translateY(-2px);
          }
          .agent-card:focus-visible {
            box-shadow: 0 0 0 2px var(--i56-color-brand);
          }

          .card-header {
            display: flex;
            align-items: center;
            gap: 0.5rem;
          }
          .agent-icon {
            font-size: ${compact ? '1rem' : '1.25rem'};
            flex-shrink: 0;
          }
          .agent-name {
            font-weight: 600;
            font-size: var(--i56-font-size-sm);
            color: var(--i56-color-text);
            flex: 1;
            min-width: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          /* Pulse dots */
          .pulse-dot {
            display: inline-block;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            flex-shrink: 0;
          }
          .pulse-dot.active {
            background: var(--i56-color-success);
            box-shadow: 0 0 0 0 var(--i56-color-success);
            animation: i56-pulse 2s infinite;
          }
          .pulse-dot.busy {
            background: var(--i56-color-warning);
            box-shadow: 0 0 0 0 var(--i56-color-warning);
            animation: i56-pulse 1.5s infinite;
          }
          .pulse-dot.error {
            background: var(--i56-color-danger);
            box-shadow: 0 0 0 0 var(--i56-color-danger);
            animation: i56-pulse-error 1s infinite;
          }
          .pulse-dot.completed {
            background: var(--i56-color-brand);
          }
          .pulse-dot.idle {
            background: var(--i56-color-text-tertiary);
          }

          @keyframes i56-pulse {
            0% { box-shadow: 0 0 0 0 currentColor; }
            70% { box-shadow: 0 0 0 6px transparent; }
            100% { box-shadow: 0 0 0 0 transparent; }
          }
          @keyframes i56-pulse-error {
            0%, 100% { box-shadow: 0 0 0 0 currentColor; }
            50% { box-shadow: 0 0 0 4px transparent; }
          }

          .card-body {
            margin-top: 0.5rem;
            display: flex;
            flex-direction: column;
            gap: 0.25rem;
          }
          .agent-status-label {
            font-size: var(--i56-font-size-xs);
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.05em;
          }
          .status-active { color: var(--i56-color-success); }
          .status-busy { color: var(--i56-color-warning); }
          .status-error { color: var(--i56-color-danger); }
          .status-completed { color: var(--i56-color-brand); }
          .status-idle { color: var(--i56-color-text-tertiary); }

          .agent-task {
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-secondary);
            line-height: 1.3;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .card-footer {
            margin-top: 0.5rem;
            padding-top: 0.375rem;
            border-top: 1px solid var(--i56-color-border);
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
          }
          .agent-last-action {
            display: flex;
            align-items: center;
            gap: 0.25rem;
          }
        </style>

        <div class="agent-grid" role="list" aria-label="AI Agents">
          ${cardsHtml}
        </div>
      `;

      // Wire card events
      this._root.querySelectorAll('.agent-card').forEach(card => {
        card.addEventListener('click', () => {
          const idx = parseInt(card.dataset.index);
          this._onAgentClick(this.agents[idx], idx);
        });
        card.addEventListener('keydown', (e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            const idx = parseInt(card.dataset.index);
            this._onAgentClick(this.agents[idx], idx);
          }
        });
      });
    }
  }

  customElements.define('i56-ai-agent-monitor', I56AiAgentMonitor);
})();
