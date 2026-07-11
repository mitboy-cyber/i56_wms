/**
 * I56 UI AI — Barrel import for BDL 2.0 AI Web Components.
 *
 * Load this single file to register all AI components:
 *   <script type="module" src="i56-ui-ai.js"></script>
 *   <script type="module">import './i56-ui-ai.js';</script>
 *
 * Or as a classic script:
 *   <script src="i56-ui-ai.js"></script>
 *
 * Components registered:
 *   <i56-ai-workspace>      — Full-screen/half-screen AI console
 *   <i56-ai-bar>            — AI Command Bar (Ctrl+K/Ctrl+J)
 *   <i56-ai-chat>           — Chat panel with SSE streaming
 *   <i56-ai-panel>          — Side drawer proactive insights
 *   <i56-ai-agent-monitor>  — Agent runtime dashboard
 *   <i56-ai-security-guardrail> — Security guardrail indicator
 *   <i56-ai-context-menu>   — Right-click AI actions
 *
 * All components use Shadow DOM, CSS custom properties (BDL 2.0 tokens),
 * custom events (i56:* prefix), and keyboard accessibility.
 *
 * BDL 2.0 Design Tokens (override on :root):
 *   --i56-brand:              #1D4ED8
 *   --i56-brand-hover:        #1E40AF
 *   --i56-brand-light:        #DBEAFE
 *   --i56-success:            #059669
 *   --i56-warning:            #D97706
 *   --i56-danger:             #DC2626
 *   --i56-info:               #2563EB
 *   --i56-bg-base:            #F8F9FB
 *   --i56-bg-surface:         #FFFFFF
 *   --i56-bg-secondary:       #F1F5F9
 *   --i56-border:             #E2E8F0
 *   --i56-text-primary:       #111827
 *   --i56-text-secondary:     #64748B
 *   --i56-text-tertiary:      #94A3B8
 *   --i56-text-inverse:       #FFFFFF
 *
 * Version: 2.0.0
 * License: MIT
 */

(function () {
  'use strict';

  // The components are self-registering via customElements.define() in their IIFEs.
  // This barrel file is a simple loader that imports all of them.

  // In a module context, each import side-effects the registration.
  // In a classic script context, they are loaded in order.

  // Note: This file is designed to be used as both ESM (import side-effects)
  // and classic script (loaded via <script src>).
  // Each sub-component self-registers, so the order below is the load order.

  const COMPONENTS = [
    './i56-ai-workspace.js',
    './i56-ai-bar.js',
    './i56-ai-chat.js',
    './i56-ai-panel.js',
    './i56-ai-agent-monitor.js',
    './i56-ai-security-guardrail.js',
    './i56-ai-context-menu.js',
  ];

  const REGISTRY = {};

  /**
   * Load all AI components as classic scripts (non-module).
   * Each component self-registers via customElements.define().
   */
  function loadAll() {
    const basePath = (document.currentScript && document.currentScript.src)
      ? document.currentScript.src.replace(/\/[^/]+$/, '/')
      : './';

    COMPONENTS.forEach(path => {
      const script = document.createElement('script');
      script.src = basePath + path;
      script.async = false; // maintain insertion order
      document.head.appendChild(script);
    });
  }

  /**
   * Check which components are registered.
   * @returns {Object} map of tag name -> boolean
   */
  function getRegistered() {
    const tags = [
      'i56-ai-workspace',
      'i56-ai-bar',
      'i56-ai-chat',
      'i56-ai-panel',
      'i56-ai-agent-monitor',
      'i56-ai-security-guardrail',
      'i56-ai-context-menu',
    ];
    const result = {};
    tags.forEach(tag => {
      result[tag] = !!customElements.get(tag);
    });
    return result;
  }

  /**
   * Check if all AI components are loaded.
   * @returns {boolean}
   */
  function isReady() {
    return Object.values(getRegistered()).every(Boolean);
  }

  /**
   * Wait for all AI components to be defined.
   * @returns {Promise<void>}
   */
  function whenReady() {
    const tags = [
      'i56-ai-workspace',
      'i56-ai-bar',
      'i56-ai-chat',
      'i56-ai-panel',
      'i56-ai-agent-monitor',
      'i56-ai-security-guardrail',
      'i56-ai-context-menu',
    ];
    return Promise.all(
      tags.map(tag =>
        customElements.whenDefined(tag).catch(() => {})
      )
    );
  }

  // If the document is already interactive/complete, load immediately.
  // Otherwise defer to DOMContentLoaded.
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', loadAll);
  } else {
    loadAll();
  }

  // Expose API on window for debugging/testing
  if (typeof window !== 'undefined') {
    window.I56AI = window.I56AI || {};
    window.I56AI.COMPONENTS = COMPONENTS;
    window.I56AI.getRegistered = getRegistered;
    window.I56AI.isReady = isReady;
    window.I56AI.whenReady = whenReady;
    window.I56AI.REGISTRY = REGISTRY;
  }
})();
