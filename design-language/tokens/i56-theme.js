/**
 * I56 Design Language (BDL) 1.0 — Theme Manager
 *
 * Manages dark/light theme toggling and tenant brand color injection.
 * Uses [data-theme] attribute on <html> for instant (CSS custom property)
 * theme switching — no classList manipulation, no layout thrash.
 *
 * Usage:
 *   I56Theme.toggle()           — toggle between dark and light
 *   I56Theme.set('dark')        — set explicit theme
 *   I56Theme.set('light')       — set explicit theme
 *   I56Theme.set('system')      — follow OS preference
 *   I56Theme.get()              — returns 'dark' | 'light'
 *   I56Theme.setBrand('#ff6600') — override brand color at runtime
 *
 * HTML integration:
 *   <html data-brand-color="#ff6600">
 *   Or: ?brand=ff6600 in URL search params.
 *
 * @license MIT
 * @version 2.0.0
 */

(function () {
  'use strict';

  const STORAGE_KEY = 'i56-theme';
  const THEME_DATA_ATTR = 'data-theme';
  const BRAND_DATA_ATTR = 'data-brand-color';
  const BRAND_URL_PARAM = 'brand';

  // ── Private helpers ──────────────────────────────────────────────────────

  /**
   * Parse a hex color into RGB components.
   * Accepts #rgb, #rrggbb, and raw hex strings.
   * Returns { r, g, b } or null.
   */
  function parseHex(hex) {
    if (!hex) return null;

    // Strip leading #
    hex = hex.replace(/^#/, '');

    // Expand shorthand (#rgb → #rrggbb)
    if (hex.length === 3) {
      hex = hex[0] + hex[0] + hex[1] + hex[1] + hex[2] + hex[2];
    }

    if (hex.length !== 6) return null;

    var r = parseInt(hex.substring(0, 2), 16);
    var g = parseInt(hex.substring(2, 4), 16);
    var b = parseInt(hex.substring(4, 6), 16);

    if (isNaN(r) || isNaN(g) || isNaN(b)) return null;

    return { r: r, g: g, b: b };
  }

  /**
   * Convert RGB to HSL.
   * Returns { h, s, l } with h in [0, 360], s/l in [0, 100].
   */
  function rgbToHsl(r, g, b) {
    r /= 255;
    g /= 255;
    b /= 255;

    var max = Math.max(r, g, b);
    var min = Math.min(r, g, b);
    var h, s, l = (max + min) / 2;

    if (max === min) {
      h = s = 0;
    } else {
      var d = max - min;
      s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

      switch (max) {
        case r: h = ((g - b) / d + (g < b ? 6 : 0)) / 6; break;
        case g: h = ((b - r) / d + 2) / 6; break;
        case b: h = ((r - g) / d + 4) / 6; break;
      }
    }

    return {
      h: Math.round(h * 360),
      s: Math.round(s * 100),
      l: Math.round(l * 100)
    };
  }

  /**
   * Generate a color scale (50–900) from a base hex color.
   * Returns an object with keys 50, 100, 200, ..., 900.
   */
  function generateScale(hex) {
    var rgb = parseHex(hex);
    if (!rgb) return null;

    var hsl = rgbToHsl(rgb.r, rgb.g, rgb.b);
    var scale = {};

    // Lighten for lower numbers, darken for higher
    var stops = [
      { key: 50,  dl: 40 },
      { key: 100, dl: 32 },
      { key: 200, dl: 22 },
      { key: 300, dl: 12 },
      { key: 400, dl: 4  },
      { key: 500, dl: 0  },
      { key: 600, dl: -8 },
      { key: 700, dl: -18 },
      { key: 800, dl: -26 },
      { key: 900, dl: -34 }
    ];

    stops.forEach(function (stop) {
      var l = Math.max(0, Math.min(100, hsl.l + stop.dl));
      scale[stop.key] = 'hsl(' + hsl.h + ', ' + hsl.s + '%, ' + l + '%)';
    });

    return scale;
  }

  /**
   * Determine the OS-level color scheme preference.
   * Returns 'dark' or 'light'.
   */
  function getSystemPreference() {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches) {
      return 'light';
    }
    return 'dark';
  }

  /**
   * Get the stored preference, falling back to system.
   * Returns 'dark' or 'light'.
   */
  function getStoredPreference() {
    try {
      var stored = localStorage.getItem(STORAGE_KEY);
      if (stored === 'dark' || stored === 'light') {
        return stored;
      }
      if (stored === 'system') {
        return getSystemPreference();
      }
    } catch (e) {
      // localStorage unavailable (private browsing, etc.) — fall through
    }
    return getSystemPreference();
  }

  // ── Public API ───────────────────────────────────────────────────────────

  var I56Theme = {
    /**
     * Set the theme by writing the [data-theme] attribute on <html>.
     * CSS custom properties cascade instantly — sub-1ms performance.
     * @param {'dark'|'light'|'system'} theme
     */
    set: function (theme) {
      var html = document.documentElement;
      var resolved;

      if (theme === 'dark' || theme === 'light') {
        resolved = theme;
      } else if (theme === 'system') {
        resolved = getSystemPreference();
      } else {
        // Unknown value — default to dark
        resolved = 'dark';
      }

      // Single attribute write — CSS custom properties recalc instantly
      html.setAttribute(THEME_DATA_ATTR, resolved);

      // Persist the raw preference (not resolved)
      try {
        localStorage.setItem(STORAGE_KEY, theme);
      } catch (e) {
        // Silently fail
      }

      // Dispatch custom event for other scripts to react
      try {
        window.dispatchEvent(new CustomEvent('i56-theme-change', {
          detail: { theme: resolved, raw: theme }
        }));
      } catch (e) {
        // Silently fail in older browsers
      }

      return resolved;
    },

    /**
     * Toggle between dark and light.
     * @returns {string} The new theme ('dark' | 'light')
     */
    toggle: function () {
      var current = I56Theme.get();
      var next = current === 'dark' ? 'light' : 'dark';
      // Store the toggled value explicitly (not 'system')
      return I56Theme.set(next);
    },

    /**
     * Get the currently active theme.
     * Reads [data-theme] attribute on <html>.
     * @returns {'dark'|'light'}
     */
    get: function () {
      var html = document.documentElement;
      var attr = html.getAttribute(THEME_DATA_ATTR);
      if (attr === 'light' || attr === 'dark') return attr;

      // No attribute set — derive from computed system preference
      if (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches) {
        return 'light';
      }
      return 'dark';
    },

    /**
     * Override the brand color at runtime.
     * Accepts hex strings like '#ff6600' or 'ff6600' or '#f60'.
     * Regenerates the full 50–900 scale and sets CSS custom properties.
     * @param {string} hex - The brand color in hex format.
     */
    setBrand: function (hex) {
      var scale = generateScale(hex);
      if (!scale) {
        console.warn('I56Theme: Invalid brand color "' + hex + '". Expected hex like #ff6600.');
        return false;
      }

      var root = document.documentElement;
      var keys = Object.keys(scale);
      for (var i = 0; i < keys.length; i++) {
        var key = keys[i];
        root.style.setProperty('--i56-brand-' + key, scale[key]);
      }

      // Also set the shorthand
      root.style.setProperty('--i56-brand', scale[500]);

      // Dispatch custom event
      try {
        window.dispatchEvent(new CustomEvent('i56-brand-change', {
          detail: { hex: hex, scale: scale }
        }));
      } catch (e) {
        // Silently fail in older browsers
      }

      return true;
    },

    /**
     * Listen for OS-level preference changes (e.g., system auto-switch at sunset).
     * Only fires if the stored preference is 'system'.
     */
    _listenSystem: function () {
      if (!window.matchMedia) return;

      try {
        var mq = window.matchMedia('(prefers-color-scheme: dark)');
        var handler = function () {
          try {
            var stored = localStorage.getItem(STORAGE_KEY);
            if (stored === 'system') {
              I56Theme.set('system');
            }
          } catch (e) {
            // Ignore
          }
        };

        if (mq.addEventListener) {
          mq.addEventListener('change', handler);
        } else if (mq.addListener) {
          // Safari < 14 fallback
          mq.addListener(handler);
        }
      } catch (e) {
        // Ignore
      }
    }
  };

  // ── Initialization ───────────────────────────────────────────────────────

  function init() {
    var html = document.documentElement;

    // 1. Apply stored or system preference
    var stored;
    try {
      stored = localStorage.getItem(STORAGE_KEY);
    } catch (e) {
      stored = null;
    }

    if (stored === 'dark' || stored === 'light' || stored === 'system') {
      I56Theme.set(stored);
    } else {
      // First visit — use system preference
      I56Theme.set('system');
    }

    // 2. Brand color — check data attribute then URL param
    var brandHex = html.getAttribute(BRAND_DATA_ATTR);
    if (!brandHex) {
      try {
        var params = new URLSearchParams(window.location.search);
        brandHex = params.get(BRAND_URL_PARAM);
      } catch (e) {
        brandHex = null;
      }
    }

    if (brandHex) {
      I56Theme.setBrand(brandHex);
    }

    // 3. Watch for OS preference changes
    I56Theme._listenSystem();
  }

  // ── Expose globally ──────────────────────────────────────────────────────

  window.I56Theme = I56Theme;

  // Run on DOM ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
