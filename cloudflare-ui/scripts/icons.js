/* ==========================================================================
   Cloudflare Kumo SVG Icon Helpers
   Maps to extracted 244-symbol spritesheet in assets/combined-icons.svg
   ========================================================================== */

window.KumoIcons = {
  render(symbolId, size = 16, className = "") {
    return `<svg width="${size}" height="${size}" viewBox="0 0 ${size} ${size}" class="fill-current shrink-0 ${className}" role="presentation" aria-hidden="true">
      <use href="assets/combined-icons.svg#${symbolId}"></use>
    </svg>`;
  },

  home(size = 16) { return this.render('icon-home', size); },
  workers(size = 16) { return this.render('icon-workers-pages', size); },
  d1(size = 16) { return this.render('icon-d1', size); },
  r2(size = 16) { return this.render('icon-r2', size); },
  analytics(size = 16) { return this.render('icon-cis-analytics', size); },
  zeroTrust(size = 16) { return this.render('icon-cloudflare-zero-trust', size); },
  bell(size = 16) { return this.render('icon-cis-bell', size); },
  chevron(size = 12) { return this.render('icon-cis-chevron-right-medium', size, 'chevron-icon'); },
  bolt(size = 16) { return this.render('icon-bolt', size); },
  settings(size = 16) { return this.render('icon-cis-settings-gear-2', size); },
  waf(size = 16) { return this.render('icon-security-waf', size); },
  kv(size = 16) { return this.render('icon-workers-kv', size); },

  search() {
    return `<svg width="15" height="15" viewBox="0 0 256 256" fill="currentColor">
      <path d="M229.66,218.34l-50.07-50.06a88.11,88.11,0,1,0-11.31,11.31l50.06,50.07a8,8,0,0,0,11.32-11.32ZM40,112a72,72,0,1,1,72,72A72.08,72.08,0,0,1,40,112Z"/>
    </svg>`;
  }
};
