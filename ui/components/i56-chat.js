/**
 * I56 Chat — Chat component with SSE streaming text rendering (Claude-style).
 *
 * Usage:
 *   <i56-chat placeholder="Ask anything…" empty-state="Start a conversation">
 *   </i56-chat>
 *
 *   const chat = document.querySelector('i56-chat');
 *   chat.addMessage({ role: 'user', content: 'Hello!' });
 *   chat.addMessage({ role: 'assistant', content: '' });  // placeholder for streaming
 *   chat.connectSSE('/api/chat/stream', { body: JSON.stringify({ messages: [...] }) });
 *
 * Attributes:
 *   placeholder: input placeholder text
 *   empty-state: message shown when no messages
 *   streaming: read-only, present while streaming
 *
 * Properties:
 *   messages: get/set array of {role, content, timestamp?}
 *
 * Methods:
 *   addMessage({role, content, timestamp?}): append a message
 *   updateLastMessage(content): update the last assistant message content
 *   connectSSE(url, options?): connect to SSE endpoint, auto-render streaming
 *   clear(): clear all messages
 *   focus(): focus the input
 *
 * Events:
 *   i56:send — detail: { content }
 *   i56:stream-end — detail: { fullContent }
 *   i56:message-click — detail: { index, message }
 *
 * Keyboard:
 *   Enter — send message
 *   Shift+Enter — newline
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
          --i56-color-text: var(--i56-text, #111827);
          --i56-color-text-secondary: var(--i56-text-secondary, #6B7280);
          --i56-color-text-tertiary: var(--i56-text-tertiary, #9CA3AF);
          --i56-color-text-inverse: var(--i56-text-inverse, #FFFFFF);
          --i56-radius: 6px;
          --i56-radius-lg: 12px;
          --i56-radius-full: 9999px;
          --i56-shadow-sm: 0 1px 2px 0 rgba(0,0,0,0.05);
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

  function escapeHtml(str) {
    return String(str ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }

  /** Simple markdown-to-HTML: bold, italic, code blocks, inline code, links */
  function renderMarkdown(text) {
    let html = escapeHtml(text);
    // Code blocks ```
    html = html.replace(/```(\w*)\n?([\s\S]*?)```/g, '<pre><code class="language-$1">$2</code></pre>');
    // Inline code `
    html = html.replace(/`([^`]+)`/g, '<code>$1</code>');
    // Bold **
    html = html.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
    // Italic *
    html = html.replace(/\*([^*]+)\*/g, '<em>$1</em>');
    // Links [text](url)
    html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener">$1</a>');
    // Line breaks
    html = html.replace(/\n\n/g, '</p><p>');
    html = html.replace(/\n/g, '<br>');
    return '<p>' + html + '</p>';
  }

  class I56Chat extends HTMLElement {
    static get observedAttributes() { return ['placeholder', 'empty-state']; }

    constructor() {
      super();
      this._root = this.attachShadow({ mode: 'open' });
      this._root.adoptedStyleSheets = [getSharedSheet()];
      this._messages = [];
      this._streaming = false;
      this._streamContent = '';
      this._abortController = null;
    }

    get messages() { return [...this._messages]; }
    set messages(arr) {
      this._messages = [...arr];
      this._streaming = false;
      this._streamContent = '';
      this.render();
    }

    connectedCallback() { this.render(); }
    attributeChangedCallback() { this.render(); }

    // -- Public API --
    addMessage(msg) {
      this._messages.push({
        role: msg.role || 'user',
        content: msg.content || '',
        timestamp: msg.timestamp || Date.now(),
      });
      this.render();
      this._scrollToBottom();
      return this._messages.length - 1;
    }

    updateLastMessage(content) {
      if (this._messages.length > 0) {
        this._messages[this._messages.length - 1].content = content;
        this._updateLastBubble(content);
        this._scrollToBottom();
      }
    }

    connectSSE(url, options = {}) {
      // Abort any existing connection
      if (this._abortController) {
        this._abortController.abort();
      }

      this._streaming = true;
      this._streamContent = '';
      this.setAttribute('streaming', '');

      // Add empty assistant message if the last isn't an empty assistant
      const last = this._messages[this._messages.length - 1];
      if (!last || last.role !== 'assistant' || last.content !== '') {
        this._messages.push({ role: 'assistant', content: '', timestamp: Date.now() });
      }

      this._abortController = new AbortController();
      const { signal } = this._abortController;

      const fetchOptions = {
        method: options.method || 'POST',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'text/event-stream',
          ...options.headers,
        },
        body: options.body || null,
        signal,
      };

      fetch(url, fetchOptions)
        .then(async (response) => {
          if (!response.ok) {
            throw new Error(`SSE connection failed: ${response.status}`);
          }
          const reader = response.body.getReader();
          const decoder = new TextDecoder();
          let buffer = '';

          while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            buffer += decoder.decode(value, { stream: true });
            const lines = buffer.split('\n');
            buffer = lines.pop() || '';

            for (const line of lines) {
              if (line.startsWith('data: ')) {
                const data = line.slice(6).trim();
                if (data === '[DONE]') {
                  this._streaming = false;
                  this.removeAttribute('streaming');
                  emit(this, 'stream-end', { fullContent: this._streamContent });
                  return;
                }
                try {
                  const parsed = JSON.parse(data);
                  const delta = parsed.delta || parsed.content || parsed.text || '';
                  this._streamContent += delta;
                  this._messages[this._messages.length - 1].content = this._streamContent;
                  this._updateLastBubble(this._streamContent);
                  this._scrollToBottom();
                } catch {
                  // Plain text delta
                  this._streamContent += data;
                  this._messages[this._messages.length - 1].content = this._streamContent;
                  this._updateLastBubble(this._streamContent);
                  this._scrollToBottom();
                }
              }
            }
          }

          this._streaming = false;
          this.removeAttribute('streaming');
          emit(this, 'stream-end', { fullContent: this._streamContent });
        })
        .catch((err) => {
          if (err.name !== 'AbortError') {
            console.error('i56-chat SSE error:', err);
            this._streaming = false;
            this.removeAttribute('streaming');
            this._messages[this._messages.length - 1].content =
              this._streamContent + `\n\n⚠️ **Error:** ${escapeHtml(err.message)}`;
            this.render();
          }
        });

      this.render();
    }

    clear() {
      this._messages = [];
      this._streaming = false;
      this._streamContent = '';
      if (this._abortController) this._abortController.abort();
      this.removeAttribute('streaming');
      this.render();
    }

    focus() {
      const textarea = this._root.querySelector('.chat-input');
      if (textarea) textarea.focus();
    }

    // -- Internals --
    _onSend = () => {
      const textarea = this._root.querySelector('.chat-input');
      if (!textarea) return;
      const content = textarea.value.trim();
      if (!content || this._streaming) return;

      textarea.value = '';
      this._resizeTextarea(textarea);

      emit(this, 'send', { content });

      // Also add the user message locally
      this.addMessage({ role: 'user', content });
    };

    _onKeydown = (e) => {
      if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        this._onSend();
      }
    };

    _resizeTextarea(el) {
      el.style.height = 'auto';
      el.style.height = Math.min(el.scrollHeight, 150) + 'px';
    }

    _scrollToBottom() {
      requestAnimationFrame(() => {
        const container = this._root.querySelector('.messages-container');
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    }

    _updateLastBubble(content) {
      const bubbles = this._root.querySelectorAll('.msg-bubble.assistant');
      const last = bubbles[bubbles.length - 1];
      if (last) {
        last.innerHTML = renderMarkdown(content);
        // Add streaming cursor if still streaming
        if (this._streaming) {
          last.innerHTML += '<span class="streaming-cursor">|</span>';
        }
      }
    }

    _formatTime(ts) {
      const d = new Date(ts);
      return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }

    render() {
      const placeholder = this.getAttribute('placeholder') || 'Type a message…';
      const emptyState = this.getAttribute('empty-state') || 'Start a conversation';
      const streaming = this._streaming;

      const messagesHtml = this._messages.map((msg, i) => {
        const isUser = msg.role === 'user';
        const isAssistant = msg.role === 'assistant';
        const isSystem = msg.role === 'system';
        const content = msg.content || '';

        if (isSystem) {
          return `<div class="system-msg">${renderMarkdown(content)}</div>`;
        }

        return `
          <div class="msg-row ${isUser ? 'user-row' : 'assistant-row'}">
            <div class="msg-avatar" aria-hidden="true">${isUser ? '👤' : '🤖'}</div>
            <div class="msg-body">
              <div class="msg-meta">
                <span class="msg-role">${isUser ? 'You' : 'Assistant'}</span>
                <span class="msg-time">${this._formatTime(msg.timestamp)}</span>
              </div>
              <div class="msg-bubble ${isUser ? 'user' : 'assistant'}">
                ${renderMarkdown(content)}${streaming && isAssistant && i === this._messages.length - 1 ? '<span class="streaming-cursor">|</span>' : ''}
              </div>
            </div>
          </div>
        `;
      }).join('');

      this._root.innerHTML = `
        <style>
          :host {
            display: flex;
            flex-direction: column;
            font-family: var(--i56-font-family);
            height: 100%;
            min-height: 300px;
            background: var(--i56-color-bg);
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-lg);
            overflow: hidden;
          }

          .messages-container {
            flex: 1;
            overflow-y: auto;
            padding: 1rem;
            display: flex;
            flex-direction: column;
            gap: 0.75rem;
            -webkit-overflow-scrolling: touch;
          }

          .empty-state {
            flex: 1;
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--i56-color-text-tertiary);
            font-size: var(--i56-font-size-sm);
            text-align: center;
            padding: 2rem;
          }

          .system-msg {
            text-align: center;
            font-size: var(--i56-font-size-xs);
            color: var(--i56-color-text-tertiary);
            padding: 0.25rem 0;
          }
          .system-msg p { margin: 0; display: inline; }

          .msg-row {
            display: flex;
            gap: 0.625rem;
            max-width: 85%;
          }
          .user-row {
            align-self: flex-end;
            flex-direction: row-reverse;
          }
          .assistant-row {
            align-self: flex-start;
          }

          .msg-avatar {
            flex-shrink: 0;
            width: 2rem;
            height: 2rem;
            border-radius: var(--i56-radius-full);
            font-size: var(--i56-font-size-base);
            display: flex;
            align-items: center;
            justify-content: center;
            background: var(--i56-color-bg-tertiary);
          }

          .msg-body {
            display: flex;
            flex-direction: column;
            gap: 0.25rem;
            min-width: 0;
          }
          .msg-meta {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-size: var(--i56-font-size-xs);
          }
          .msg-role {
            font-weight: 600;
            color: var(--i56-color-text-secondary);
          }
          .msg-time {
            color: var(--i56-color-text-tertiary);
          }
          .user-row .msg-meta { justify-content: flex-end; }

          .msg-bubble {
            padding: 0.625rem 0.875rem;
            border-radius: var(--i56-radius-lg);
            font-size: var(--i56-font-size-sm);
            line-height: var(--i56-line-height);
            word-break: break-word;
          }
          .msg-bubble.user {
            background: var(--i56-color-brand);
            color: var(--i56-color-text-inverse);
            border-bottom-right-radius: var(--i56-radius);
          }
          .msg-bubble.assistant {
            background: var(--i56-color-bg-secondary);
            color: var(--i56-color-text);
            border: 1px solid var(--i56-color-border);
            border-bottom-left-radius: var(--i56-radius);
          }

          .msg-bubble p { margin: 0; }
          .msg-bubble p + p { margin-top: 0.5rem; }
          .msg-bubble code {
            background: rgba(0,0,0,0.06);
            padding: 0.125rem 0.375rem;
            border-radius: 4px;
            font-size: 0.8125rem;
            font-family: 'SF Mono', 'Fira Code', monospace;
          }
          .msg-bubble pre {
            background: #1e1e2e;
            color: #cdd6f4;
            padding: 0.75rem 1rem;
            border-radius: var(--i56-radius);
            overflow-x: auto;
            margin: 0.5rem 0;
            font-size: var(--i56-font-size-xs);
          }
          .msg-bubble pre code {
            background: none;
            padding: 0;
            font-size: inherit;
          }
          .msg-bubble a {
            color: var(--i56-color-brand);
            text-decoration: underline;
          }
          .msg-bubble.user a { color: var(--i56-color-text-inverse); }

          .streaming-cursor {
            display: inline-block;
            animation: i56-blink 1s step-end infinite;
            color: var(--i56-color-brand);
            font-weight: bold;
            margin-left: 1px;
          }
          @keyframes i56-blink {
            0%, 100% { opacity: 1; }
            50% { opacity: 0; }
          }

          .input-area {
            display: flex;
            align-items: flex-end;
            gap: 0.5rem;
            padding: 0.75rem 1rem;
            border-top: 1px solid var(--i56-color-border);
            background: var(--i56-color-bg);
          }
          .chat-input {
            flex: 1;
            border: 1px solid var(--i56-color-border);
            border-radius: var(--i56-radius-lg);
            padding: 0.5rem 0.75rem;
            font-size: var(--i56-font-size-sm);
            font-family: var(--i56-font-family);
            color: var(--i56-color-text);
            background: var(--i56-color-bg);
            resize: none;
            outline: none;
            min-height: 2.25rem;
            max-height: 150px;
            line-height: var(--i56-line-height);
            transition: border-color var(--i56-transition);
          }
          .chat-input:focus {
            border-color: var(--i56-color-brand);
            box-shadow: 0 0 0 3px var(--i56-color-brand-light);
          }
          .chat-input::placeholder { color: var(--i56-color-text-tertiary); }
          .chat-input:disabled { opacity: 0.5; background: var(--i56-color-bg-secondary); }

          .send-btn {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 2.25rem;
            height: 2.25rem;
            border: none;
            border-radius: var(--i56-radius);
            background: var(--i56-color-brand);
            color: var(--i56-color-text-inverse);
            cursor: pointer;
            font-size: var(--i56-font-size-lg);
            flex-shrink: 0;
            transition: all var(--i56-transition);
          }
          .send-btn:hover { background: var(--i56-color-brand-hover); }
          .send-btn:disabled { opacity: 0.4; cursor: not-allowed; }
          .send-btn:focus-visible { box-shadow: 0 0 0 3px var(--i56-color-brand-light); outline: none; }
        </style>

        <div class="messages-container" role="log" aria-live="polite">
          ${this._messages.length === 0 ? `<div class="empty-state">${emptyState}</div>` : messagesHtml}
        </div>
        <div class="input-area">
          <textarea
            class="chat-input"
            rows="1"
            placeholder="${placeholder}"
            ?disabled="${streaming}"
            aria-label="Chat message input"
          ></textarea>
          <button class="send-btn" ?disabled="${streaming}" aria-label="Send message" title="Send (Enter)">↑</button>
        </div>
      `;

      // Wire events
      const textarea = this._root.querySelector('.chat-input');
      const sendBtn = this._root.querySelector('.send-btn');

      if (textarea) {
        textarea.addEventListener('keydown', this._onKeydown);
        textarea.addEventListener('input', () => this._resizeTextarea(textarea));
      }
      if (sendBtn) {
        sendBtn.addEventListener('click', this._onSend);
      }
    }
  }

  customElements.define('i56-chat', I56Chat);
})();
