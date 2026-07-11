/**
 * I56 AI Panel — Side drawer for proactive AI insights and diagnostics.
 *
 * BDL 2.0 AI Web Component: native Custom Element, Shadow DOM, zero dependencies.
 * Provides auto-analysis triggers, diagnostic reports, and entity-aware insights.
 *
 * Usage:
 *   <i56-ai-panel data-context='{"entity":"order","id":"ORD-1234","section":"detail"}'
 *                 auto-analyze>
 *   </i56-ai-panel>
 *
 * Attributes:
 *   data-context: JSON string with current entity context
 *   auto-analyze: present to trigger analysis on context change
 *   expanded: present when panel is open
 *   position: left | right (default: right)
 *   width: panel width (default: 380px)
 *
 * Properties:
 *   context: get/set context object
 *   expanded: get/set expanded state
 *   insights: get/set insights array
 *
 * Events:
 *   i56:panel-open — detail: { context }
 *   i56:panel-close — detail: {}
 *   i56:analyze — detail: { context } — request analysis from parent
 *   i56:insight-click — detail: { insight, index }
 *
 * Keyboard:
 *   Escape — close panel
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

  class I56AiPanel extends HTMLElement {
    static get observedAttributes() {
      return ['data-context', 'auto-analyze', 'expanded', 'position', 'width'];
    }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
      this._insights = [];
      this._boundKeydown = this._onKeydown.bind(this);
      this._prevContext = '';
    }

    get context() {
      try { return JSON.parse(this.getAttribute('data-context') || '{}'); } catch { return {}; }
    }
    set context(obj) {
      this.setAttribute('data-context', JSON.stringify(obj));
    }

    get expanded() { return this.hasAttribute('expanded'); }
    get autoAnalyze() { return this.hasAttribute('auto-analyze'); }
    get position() { return this.getAttribute('position') || 'right'; }

    get insights() { return [...this._insights]; }
    set insights(arr) {
      this._insights = [...arr];
      this.render();
    }

    connectedCallback() {
      document.addEventListener('keydown', this._boundKeydown);
      this.render();
    }

    disconnectedCallback() {
      document.removeEventListener('keydown', this._boundKeydown);
    }

    attributeChangedCallback(name, oldVal, newVal) {
      if (name === 'data-context' && newVal !== this._prevContext) {
        this._prevContext = newVal;
        if (this.autoAnalyze && this.expanded) {
          this._runAnalysis();
        }
      }
      if (name === 'expanded') this.render();
    }

    // -- Public API --
    open() {
      if (this.expanded) return;
      this.setAttribute('expanded', '');
      if (this.autoAnalyze) this._runAnalysis();
      emit(this, 'panel-open', { context: this.context });
    }

    close() {
      if (!this.expanded) return;
      this.removeAttribute('expanded');
      emit(this, 'panel-close', {});
    }

    toggle() { this.expanded ? this.close() : this.open(); }

    addInsight(insight) {
      this._insights.push({
        id: insight.id || ('ins-' + Math.random().toString(36).slice(2, 7)),
        type: insight.type || 'info',
        title: insight.title || '',
        description: insight.description || '',
        timestamp: insight.timestamp || Date.now(),
        action: insight.action || null,
        severity: insight.severity || 'low',
      });
      this.render();
      return this._insights.length - 1;
    }

    clearInsights() {
      this._insights = [];
      this.render();
    }

    // -- Internals --
    _runAnalysis() {
      emit(this, 'analyze', { context: this.context });
      // Auto-generate placeholder insights based on context
      if (this._insights.length === 0) {
        const ctx = this.context;
        if (ctx.entity) {
          this.addInsight({
            type: 'info',
            title: `Analyzing ${ctx.entity}${ctx.id ? ' ' + ctx.id : ''}`,
            description: `AI is reviewing this ${ctx.entity} for patterns, anomalies, and recommendations.`,
          });
        }
        this.addInsight({
          type: 'suggestion',
          title: 'Quick Actions',
          description: 'Click to see suggested AI actions for this context.',
          action: { label: 'View Suggestions', event: 'suggestions' },
        });
      }
    }

    _onKeydown(e) {
      if (e.key === 'Escape' && this.expanded) {
        this.close();
      }
    }

    _onInsightClick(insight, index) {
      emit(this, 'insight-click', { insight, index });
      if (insight.action) {
        emit(this, 'insight-action', { action: insight.action, insight });
      }
    }

    _onManualAnalyze() {
      this._insights = [];
      this._runAnalysis();
    }

    _renderInsightCard(insight, index) {
      const typeIcons = {
        info: 'ℹ️',
        warning: '⚠️',
        success: '✅',
        danger: '🔴',
        suggestion: '💡',
        diagnostic: '🔍',
      };

      const typeClasses = {
        info: 'insight-info',
        warning: 'insight-warning',
        success: 'insight-success',
        danger: 'insight-danger',
        suggestion: 'insight-suggestion',
        diagnostic: 'insight-diagnostic',
      };

      const icon = typeIcons[insight.type] || 'ℹ️';
      const typeClass = typeClasses[insight.type] || 'insight-info';

      const time = new Date(insight.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

      return `
        <div class="insight-card ${typeClass}" tabindex="0" role="button"
             aria-label="${insight.title}"
             data-index="${index}">
          <div class="insight-header">
            <span class="insight-icon" aria-hidden="true">${icon}</span>
            <span class="insight-title">${insight.title}</span>
            <span class="insight-time">${time}</span>
          </div>
          ${insight.description ? `<div class="insight-desc">${insight.description}</div>` : ''}
          ${insight.action ? `
            <button class="insight-action-btn" type="button">
              ${insight.action.label || 'Take Action'}
              <span aria-hidden="true">→</span>
            </button>
          ` : ''}
        </div>
      `;
    }

    render() {
      const isExpanded = this.expanded;
      const pos = this.position;
      const isRight = pos === 'right';
      const width = this.getAttribute('width') || '380px';

      this._root.innerHTML = `
        <style>
          :host {
            display: contents;
          }

          .ai-panel-container {
            position: fixed;
            top: 0;
            bottom: 0;
            ${isRight ? 'right' : 'left'}: 0;
            width: ${width};
            max-width: 90vw;
            z-index: 1050;
            transform: translateX(${isExpanded ? '0' : (isRight ? '100%' : '-100%')});
            transition: transform var(--i56-transition-slow) cubic-bezier(0.4, 0, 0.2, 1);
            display: flex;
            flex-direction: column;
            background: var(--i56-color-bg-surface);
            border-${isRight ? 'left' : 'right'}: 1px solid var(--i56-color-border);
            box-shadow: var(--i56-shadow-xl);
            font-family: var(--i56-font-family);
          }

          /* Toggle tab */
          .panel-tab {
            position: absolute;
            ${isRight ? 'left' : 'right'}: 100%;
            top: 50%;
            transform: translateY(-50%);
            width: 2rem;
            height: 5rem;
            display: flex;
            align-items: center;
            justify-content: center;
            background: var(--i56-color-brand);
            color: var(--i56-color-text-inverse);
            border: none;
            border-radius: ${isRight ? 'var(--i56-radius-md) 0 0 var(--i56-radius-md)' : '0 var(--i56-radius-md) var(--i56-radius-md) 0'};
            cursor: pointer;
            font-size: var(--i56-font-size-lg);
            writing-mode: vertical-rl;
            letter-spacing: 0.1em;
            font-weight: 600;
            font-family: var(--i56-font-family);
            transition: all var(--i56-transition);
            box-shadow: var(--i56-shadow-lg);
            z-index: -1;
          }
          .panel-tab:hover {
            background: var(--i56-color-brand-hover);
            ${isRight ? 'left' : 'right'}: calc(100% + 4px);
          }
          .panel-tab:focus-visible {
            outline: none;
            box-shadow: 0 0 0 3px var(--i56-color-brand-light);
          }

          /* Header */
          .panel-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 1rem 1.25rem;
            border-bottom: 1px solid var(--i56-color-border);
            flex-shrink: 0;
          }
          .panel-title {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-weight: 700;
            font-size: var(--i56-font-size-base);
            color: var(--i56-color-text);
          }
          .panel-title-icon { font-size: var(--i56-font-size-lg); }

          .context-badge {
            display: inline-flex;
            align-items: center;
            gap: 0.25rem;
            padding: 0.125rem 0.5rem;
            border-radius: var(--i56-radius-full);
            font-size: var(--i56-font-size-xs);
            font-weight: 500;
            background: var(--i56-color-brand-light);
            color: var(--i56-color-brand);
          }

          .panel-actions {
            display: flex;
            align-items: center;
            gap: 0.375rem;
          }
          .panel-btn {
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
          .panel-btn:hover {
            background: var(--i56-color-bg-secondary);
            color: var(--i56-color-text);
          }
          .panel-btn:focus-visible {
            box-shadow: 0 0 0 2px var(--i56-color-brand-light);
            outline: none;
          }

          /* Body */
          .panel-body {
            flex: 1;
            overflow-y: auto;
            padding: 1rem 1.25rem;
            display: flex;
            flex-direction: column;
            gap: 0.75rem;
          }

          .empty-state {
            flex: 1;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            color: var(--i56-color-text-tertiary);
            font-size: var(--i56-font-size-sm);
            text-align: center;
            gap: 0.75rem;
          }
          .empty-state::before {
            content: '🔮';
            font-size: 2.5rem;
            opacity: 0.3;
          }

          .analyze-btn {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.5rem 1rem;
            border: 1px solid var(--i56-color-brand);
            border-radius: var(--i56-radius-full);
            background: var(--i56-color-brand-light);
            color: var(--i56-color-brand);
            font-size: var(--i56-font-size-sm);
            font-weight: 500;
            cursor: pointer;
            font-family: var(--i56-font-family);
            transition: all var(--i56-transition);
          }
          .analyze-btn:hover {
            background: var(--i56-color-brand);
            color: var(--i56-color-text-inverse);
          }

          /* Insight cards */
          .insight-card {
            padding: 0.75rem 1rem;
            border-radius: var(--i56-radius-lg);
            border: 1px solid var(--i56-color-border);
            background: var(--i56-color-bg-surface);
            cursor: pointer;
            transition: all var(--i56-transition);
          }
          .insight-card:hover {
            border-color: var(--i56-color-border-hover);
            box-shadow: var(--i56-shadow-lg);
            transform: translateY(-1px);
          }
          .insight-card:focus-visible {
            box-shadow: 0 0 0 2px var(--i56-color-brand);
            outline: none;
          }

          .insight-info { border-left: 3px solid var(--i56-color-info); }
          .insight-warning { border-left: 3px solid var(--i56-color-warning); background: var(--i56-color-warning-light); }
          .insight-success { border-left: 3px solid var(--i56-color-success); }
          .insight-danger { border-left: 3px solid var(--i56-color-danger); background: var(--i56-color-danger-light); }
          .insight-suggestion { border-left: 3px solid var(--i56-color-brand); background: var(--i56-color-brand-light); }
          .insight-diagnostic { border-left: 3px solid var(--i56-color-info); }

          .insight-header {
            display: flex;
            align-items: center;
            gap: 0.5rem;
          }
          .insight-icon { font-size: var(--i56-font-size-base); flex-shrink: 0; }
          .insight-title {
            font-weight: 600;
            font-size: var(--i56-font-size-sm);
            color: var(--i56-color-text);
            flex: 1;
            min-width: 0;
          }
          .insight-time {
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
            flex-shrink: 0;
          }
          .insight-desc {
            margin-top: 0.375rem;
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-secondary);
            line-height: var(--i56-line-height);
          }
          .insight-action-btn {
            display: inline-flex;
            align-items: center;
            gap: 0.25rem;
            margin-top: 0.5rem;
            padding: 0.25rem 0.75rem;
            border: none;
            border-radius: var(--i56-radius-full);
            background: var(--i56-color-brand);
            color: var(--i56-color-text-inverse);
            font-size: var(--i56-font-size-xs);
            font-weight: 500;
            cursor: pointer;
            font-family: var(--i56-font-family);
            transition: all var(--i56-transition);
          }
          .insight-action-btn:hover { background: var(--i56-color-brand-hover); }

          /* Footer */
          .panel-footer {
            padding: 0.75rem 1.25rem;
            border-top: 1px solid var(--i56-color-border);
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
            text-align: center;
            flex-shrink: 0;
          }
        </style>

        <div class="ai-panel-container" role="complementary" aria-label="AI Insights Panel"
             aria-expanded="${isExpanded}">
          <button class="panel-tab"
                  aria-label="${isExpanded ? 'Close AI Panel' : 'Open AI Panel'}"
                  title="${isExpanded ? 'Close' : 'AI Insights'}">
            ${isExpanded ? '✕' : 'AI'}
          </button>

          <div class="panel-header">
            <div class="panel-title">
              <span class="panel-title-icon" aria-hidden="true">🔮</span>
              <span>AI Insights</span>
              ${this.context.entity ? '<span class="context-badge">' + this.context.entity + (this.context.id ? ' ' + this.context.id : '') + '</span>' : ''}
            </div>
            <div class="panel-actions">
              <button class="panel-btn" title="Refresh analysis" aria-label="Refresh analysis">🔄</button>
              <button class="panel-btn" title="Close panel" aria-label="Close panel">✕</button>
            </div>
          </div>

          <div class="panel-body">
            ${this._insights.length === 0 ? `
              <div class="empty-state">
                <span>No insights yet</span>
                <span style="font-size: var(--i56-font-size-xs); color: var(--i56-color-text-tertiary);">
                  AI will analyze your context and surface actionable insights here.
                </span>
                <button class="analyze-btn">🔍 Analyze Current Context</button>
              </div>
            ` : this._insights.map((ins, i) => this._renderInsightCard(ins, i)).join('')}
          </div>

          <div class="panel-footer">
            ${this._insights.length} insight${this._insights.length !== 1 ? 's' : ''}
            ${this.autoAnalyze ? ' · Auto-analyze on' : ''}
          </div>
        </div>
      `;

      // Wire tab button
      const tabBtn = this._root.querySelector('.panel-tab');
      if (tabBtn) {
        tabBtn.addEventListener('click', () => this.toggle());
      }

      // Wire close button
      const closeBtn = this._root.querySelector('.panel-btn[aria-label="Close panel"]');
      if (closeBtn) {
        closeBtn.addEventListener('click', () => this.close());
      }

      // Wire refresh button
      const refreshBtn = this._root.querySelector('.panel-btn[aria-label="Refresh analysis"]');
      if (refreshBtn) {
        refreshBtn.addEventListener('click', () => this._onManualAnalyze());
      }

      // Wire analyze button
      const analyzeBtn = this._root.querySelector('.analyze-btn');
      if (analyzeBtn) {
        analyzeBtn.addEventListener('click', () => this._onManualAnalyze());
      }

      // Wire insight card clicks
      this._root.querySelectorAll('.insight-card').forEach(card => {
        card.addEventListener('click', () => {
          const idx = parseInt(card.dataset.index);
          this._onInsightClick(this._insights[idx], idx);
        });
        card.addEventListener('keydown', (e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            const idx = parseInt(card.dataset.index);
            this._onInsightClick(this._insights[idx], idx);
          }
        });
      });

      // Wire insight action buttons
      this._root.querySelectorAll('.insight-action-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
          e.stopPropagation();
          const card = btn.closest('.insight-card');
          if (card) {
            const idx = parseInt(card.dataset.index);
            const insight = this._insights[idx];
            emit(this, 'insight-action', { action: insight.action, insight });
          }
        });
      });
    }
  }

  customElements.define('i56-ai-panel', I56AiPanel);
})();
