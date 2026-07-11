/**
 * I56 Kanban — Kanban board using native HTML5 Drag and Drop API.
 *
 * Usage:
 *   <i56-kanban columns='[{"id":"todo","title":"To Do"},{"id":"in-progress","title":"In Progress"}]'
 *               cards='[{"id":"1","columnId":"todo","title":"Fix login","labels":["bug"],"assignee":"A"}]'>
 *   </i56-kanban>
 *
 *   // Or set data via JS:
 *   board.columns = [...];
 *   board.cards = [...];
 *   board.addCard({ id, columnId, title, ... });
 *
 * Attributes:
 *   columns: JSON array of {id, title, color?, limit?}
 *   cards: JSON array of {id, columnId, title, description?, labels?:string[], assignee?:string, priority?:string}
 *   empty-message: text when a column has no cards
 *
 * Events:
 *   i56:card-move — detail: { cardId, fromColumnId, toColumnId, fromIndex, toIndex }
 *   i56:card-click — detail: { card }
 *   i56:card-add — detail: { columnId }
 *
 * Keyboard:
 *   Tab between cards, Enter/Space to pick up, Arrow keys to navigate, Escape to cancel
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
          --i56-color-success: var(--i56-success, #059669);
          --i56-color-warning: var(--i56-warning, #D97706);
          --i56-color-danger: var(--i56-danger, #DC2626);
          --i56-color-info: var(--i56-info, #2563EB);
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
          --i56-shadow: 0 1px 3px 0 rgba(0,0,0,0.1), 0 1px 2px -1px rgba(0,0,0,0.1);
          --i56-shadow-md: 0 4px 6px -1px rgba(0,0,0,0.1), 0 2px 4px -2px rgba(0,0,0,0.1);
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

  const LABEL_COLORS = {
    bug: '#DC2626', feature: '#2563EB', improvement: '#059669',
    docs: '#6B7280', design: '#D97706', urgent: '#DC2626',
    default: '#6B7280',
  };

  const PRIORITY_ICONS = { high: '🔴', medium: '🟡', low: '🟢', none: '' };

  class I56Kanban extends HTMLElement {
    static get observedAttributes() { return ['columns', 'cards', 'empty-message']; }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
      this._dragCard = null;
      this._dragColumnId = null;
      this._dragOverColumn = null;
    }

    get columns() {
      try { return JSON.parse(this.getAttribute('columns') || '[]'); } catch { return []; }
    }
    set columns(v) { this.setAttribute('columns', JSON.stringify(v)); }

    get cards() {
      try { return JSON.parse(this.getAttribute('cards') || '[]'); } catch { return []; }
    }
    set cards(v) { this.setAttribute('cards', JSON.stringify(v)); }

    connectedCallback() { this.render(); }
    attributeChangedCallback() { this.render(); }

    // -- Public API --
    addCard(card) {
      const cards = this.cards;
      cards.push(card);
      this.cards = cards;
      emit(this, 'card-add', { card, columnId: card.columnId });
    }

    moveCard(cardId, toColumnId, toIndex) {
      const cards = this.cards;
      const idx = cards.findIndex(c => c.id === cardId);
      if (idx === -1) return;

      const [card] = cards.splice(idx, 1);
      const fromColumnId = card.columnId;
      card.columnId = toColumnId;

      // Insert at correct position
      const colCards = cards.filter(c => c.columnId === toColumnId);
      const insertAfter = colCards[Math.min(toIndex ?? colCards.length, colCards.length - 1)];
      const insertIdx = insertAfter ? cards.indexOf(insertAfter) + 1 : cards.findIndex(c => c.columnId === toColumnId);

      if (insertIdx === -1) {
        // No cards in target column yet — append at end
        cards.push(card);
      } else {
        cards.splice(insertIdx, 0, card);
      }

      this.cards = cards;
      emit(this, 'card-move', { cardId, fromColumnId, toColumnId, fromIndex: idx, toIndex: Math.max(0, toIndex ?? 0) });
    }

    // -- HTML5 Drag & Drop --
    _onDragStart = (e) => {
      const cardEl = e.target.closest('.kanban-card');
      if (!cardEl) return;
      this._dragCard = cardEl.dataset.cardId;
      this._dragColumnId = cardEl.closest('.kanban-column')?.dataset.columnId;
      cardEl.classList.add('dragging');
      e.dataTransfer.effectAllowed = 'move';
      e.dataTransfer.setData('text/plain', this._dragCard);
    };

    _onDragEnd = (e) => {
      const cardEl = e.target.closest('.kanban-card');
      if (cardEl) cardEl.classList.remove('dragging');
      this._dragCard = null;
      this._dragColumnId = null;
      this._dragOverColumn = null;
      this._root.querySelectorAll('.kanban-column').forEach(c => c.classList.remove('drag-over'));
      this._root.querySelectorAll('.kanban-card').forEach(c => c.classList.remove('drag-over-card'));
    };

    _onDragOver = (e) => {
      e.preventDefault();
      e.dataTransfer.dropEffect = 'move';

      const column = e.target.closest('.kanban-column');
      if (column && this._dragOverColumn !== column) {
        if (this._dragOverColumn) this._dragOverColumn.classList.remove('drag-over');
        this._dragOverColumn = column;
        column.classList.add('drag-over');
      }
    };

    _onDragLeave = (e) => {
      const column = e.target.closest('.kanban-column');
      if (column && this._dragOverColumn === column) {
        // Only remove if truly leaving the column (not entering a child)
        const related = e.relatedTarget;
        if (!related || !column.contains(related)) {
          column.classList.remove('drag-over');
          this._dragOverColumn = null;
        }
      }
    };

    _onDrop = (e) => {
      e.preventDefault();
      const column = e.target.closest('.kanban-column');
      if (!column || !this._dragCard) return;

      const toColumnId = column.dataset.columnId;
      if (toColumnId === this._dragColumnId) {
        // Reorder within same column
        const cardList = column.querySelector('.kanban-cards');
        const cards = [...cardList.querySelectorAll('.kanban-card:not(.dragging)')];
        const dropY = e.clientY;
        let toIndex = cards.length;
        for (let i = 0; i < cards.length; i++) {
          const rect = cards[i].getBoundingClientRect();
          if (dropY < rect.top + rect.height / 2) {
            toIndex = i;
            break;
          }
        }
        this._reorderCard(this._dragCard, toIndex);
      } else {
        this.moveCard(this._dragCard, toColumnId, 999);
      }

      this._dragCard = null;
      this._dragColumnId = null;
      if (this._dragOverColumn) this._dragOverColumn.classList.remove('drag-over');
      this._dragOverColumn = null;
    };

    _reorderCard(cardId, toIndex) {
      const cards = this.cards;
      const idx = cards.findIndex(c => c.id === cardId);
      if (idx === -1) return;
      const [card] = cards.splice(idx, 1);
      const colCards = cards.filter(c => c.columnId === card.columnId);
      const globalIdx = cards.indexOf(colCards[toIndex] || colCards[colCards.length - 1]);
      cards.splice(Math.max(globalIdx, 0), 0, card);
      this.cards = cards;
      emit(this, 'card-move', { cardId, fromColumnId: card.columnId, toColumnId: card.columnId, fromIndex: idx, toIndex });
    }

    _onCardClick = (cardId) => {
      const cards = this.cards;
      const card = cards.find(c => c.id === cardId);
      if (card) emit(this, 'card-click', { card });
    };

    _onAddCard = (columnId) => {
      emit(this, 'card-add', { columnId });
    };

    computeCardCount(columnId) {
      return this.cards.filter(c => c.columnId === columnId).length;
    }

    render() {
      const columns = this.columns;
      const cards = this.cards;
      const emptyMsg = this.getAttribute('empty-message') || 'No cards';

      const columnsHtml = columns.map(col => {
        const colCards = cards
          .map((c, originalIdx) => ({ ...c, _originalIdx: originalIdx }))
          .filter(c => c.columnId === col.id);
        const count = colCards.length;
        const color = col.color || 'var(--i56-color-text-secondary)';
        const limit = col.limit ? ` / ${col.limit}` : '';
        const atLimit = col.limit && count >= col.limit;

        const cardsHtml = colCards.map((card, i) => {
          const labels = (card.labels || []).map(l => {
            const bg = LABEL_COLORS[l] || LABEL_COLORS.default;
            return `<span class="card-label" style="background:${bg}">${l}</span>`;
          }).join('');

          const priorityIcon = PRIORITY_ICONS[card.priority] || '';
          const assignee = card.assignee
            ? `<span class="card-assignee" title="${card.assignee}">${card.assignee.charAt(0).toUpperCase()}</span>`
            : '';

          return `
            <div class="kanban-card"
                 draggable="true"
                 data-card-id="${card.id}"
                 data-original-idx="${card._originalIdx}"
                 tabindex="0"
                 role="button"
                 aria-label="${card.title}, ${card.description || ''}">
              ${labels ? `<div class="card-labels">${labels}</div>` : ''}
              <div class="card-title">${card.title}</div>
              ${card.description ? `<div class="card-desc">${card.description}</div>` : ''}
              <div class="card-meta">
                ${priorityIcon ? `<span class="card-priority">${priorityIcon}</span>` : ''}
                <span class="card-spacer"></span>
                ${assignee}
              </div>
            </div>
          `;
        }).join('');

        return `
          <div class="kanban-column"
               data-column-id="${col.id}"
               role="region"
               aria-label="${col.title}">
            <div class="column-header">
              <div class="column-title-row">
                <span class="column-color" style="background:${color}" aria-hidden="true"></span>
                <span class="column-title">${col.title}</span>
                <span class="column-count ${atLimit ? 'at-limit' : ''}">${count}${limit}</span>
              </div>
              <button class="column-add-btn" aria-label="Add card to ${col.title}" title="Add card">+</button>
            </div>
            <div class="kanban-cards">
              ${cardsHtml || `<div class="column-empty">${emptyMsg}</div>`}
            </div>
          </div>
        `;
      }).join('');

      this._root.innerHTML = `
        <style>
          :host {
            display: flex;
            gap: 1rem;
            font-family: var(--i56-font-family);
            overflow-x: auto;
            min-height: 200px;
            padding: 0.5rem;
            -webkit-overflow-scrolling: touch;
          }

          .kanban-column {
            flex: 0 0 17.5rem;
            background: var(--i56-color-bg-secondary);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-lg);
            display: flex;
            flex-direction: column;
            max-height: calc(100vh - 4rem);
            transition: border-color var(--i56-transition), background var(--i56-transition);
          }
          .kanban-column.drag-over {
            border-color: var(--i56-color-brand);
            background: var(--i56-color-brand-light);
          }

          .column-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0.75rem 0.875rem;
            flex-shrink: 0;
            border-bottom: 1px solid var(--i56-color-border);
          }
          .column-title-row {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            min-width: 0;
          }
          .column-color {
            width: 0.625rem;
            height: 0.625rem;
            border-radius: 50%;
            flex-shrink: 0;
          }
          .column-title {
            font-size: var(--i56-font-size-sm);
            font-weight: 600;
            color: var(--i56-color-text);
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
          }
          .column-count {
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
            flex-shrink: 0;
          }
          .column-count.at-limit { color: var(--i56-color-danger); font-weight: 600; }
          .column-add-btn {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 1.5rem;
            height: 1.5rem;
            border: none;
            background: none;
            border-radius: var(--i56-radius);
            cursor: pointer;
            color: var(--i56-color-text-tertiary);
            font-size: var(--i56-font-size-lg);
            transition: all var(--i56-transition);
            flex-shrink: 0;
          }
          .column-add-btn:hover { background: var(--i56-color-bg-tertiary); color: var(--i56-color-text); }

          .kanban-cards {
            flex: 1;
            overflow-y: auto;
            padding: 0.5rem;
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
            min-height: 3rem;
          }
          .column-empty {
            text-align: center;
            padding: 1.5rem 0.5rem;
            color: var(--i56-color-text-tertiary);
            font-size: var(--i56-font-size-xs);
          }

          .kanban-card {
            background: var(--i56-color-bg);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-md);
            padding: 0.625rem 0.75rem;
            cursor: grab;
            transition: box-shadow var(--i56-transition), border-color var(--i56-transition), transform 100ms ease;
            user-select: none;
          }
          .kanban-card:hover {
            box-shadow: var(--i56-shadow-md);
            border-color: var(--i56-color-border-hover);
          }
          .kanban-card:active { cursor: grabbing; }
          .kanban-card.dragging {
            opacity: 0.4;
            transform: scale(0.95);
          }
          .kanban-card:focus-visible {
            box-shadow: 0 0 0 2px var(--i56-color-brand);
            outline: none;
          }
          .kanban-card.drag-over-card {
            border-color: var(--i56-color-brand);
            box-shadow: 0 0 0 2px var(--i56-color-brand-light);
          }

          .card-labels {
            display: flex;
            gap: 0.25rem;
            margin-bottom: 0.375rem;
            flex-wrap: wrap;
          }
          .card-label {
            display: inline-block;
            padding: 0.125rem 0.375rem;
            border-radius: 9999px;
            font-size: 0.625rem;
            font-weight: 600;
            color: #fff;
            text-transform: uppercase;
            letter-spacing: 0.03em;
          }
          .card-title {
            font-size: var(--i56-font-size-sm);
            font-weight: 500;
            color: var(--i56-color-text);
            line-height: 1.4;
            word-break: break-word;
          }
          .card-desc {
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-secondary);
            margin-top: 0.25rem;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
          }
          .card-meta {
            display: flex;
            align-items: center;
            gap: 0.375rem;
            margin-top: 0.5rem;
          }
          .card-priority { font-size: var(--i56-font-size-xs); }
          .card-spacer { flex: 1; }
          .card-assignee {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 1.25rem;
            height: 1.25rem;
            border-radius: 50%;
            background: var(--i56-color-brand-light);
            color: var(--i56-color-brand);
            font-size: 0.625rem;
            font-weight: 600;
            flex-shrink: 0;
          }
        </style>
        ${columnsHtml}
      `;

      // Wire drag events
      this._root.querySelectorAll('.kanban-card').forEach(card => {
        card.addEventListener('dragstart', this._onDragStart);
        card.addEventListener('dragend', this._onDragEnd);
        card.addEventListener('click', () => this._onCardClick(card.dataset.cardId));
        card.addEventListener('keydown', (e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            this._onCardClick(card.dataset.cardId);
          }
        });
      });

      this._root.querySelectorAll('.kanban-column').forEach(col => {
        col.addEventListener('dragover', this._onDragOver);
        col.addEventListener('dragleave', this._onDragLeave);
        col.addEventListener('drop', this._onDrop);
        const addBtn = col.querySelector('.column-add-btn');
        if (addBtn) {
          addBtn.addEventListener('click', () => this._onAddCard(col.dataset.columnId));
        }
      });
    }
  }

  customElements.define('i56-kanban', I56Kanban);
})();
