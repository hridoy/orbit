/* ==========================================================================
   Cloudflare Dashboard (Kumo Design System) — Master Single Page Application Controller
   1:1 Authentic Mirror of Production Dashboard (Account: 330b08634cd698ef2ff8587d66ac4b2c)
   Pure Native Vanilla JS — Zero External Runtime Dependencies, 100% Envato Commercial Ready
   ========================================================================== */

(function () {
  'use strict';

  // Master Reactive Application State
  const state = {
    currentTab: 'home',
    subTab: 'overview',
    selectedDomain: 'northlab.dev',
    selectedWorker: 'code-collection',
    selectedDb: '670002a1-53c5-4f67-b8c9-6aff880f8052',
    selectedBucket: 'code-collection',
    searchOpen: false,
    searchQuery: '',
    theme: localStorage.getItem('cf_theme') || 'light',
    accent: localStorage.getItem('cf_accent') || 'orange',
    sidebarMini: localStorage.getItem('cf_sidebar_mini') === 'true',
    mobileSidebarOpen: false,
    customizerOpen: false,
    openGroups: {
      recents: true,
      compute: true,
      storage: true,
      networking: true,
      protect: true,
      observability: false,
      account: true,
    },
    activeModal: null,
    copiedKey: null,
    analyticsTimeRange: '24h',
    data: window.CLOUDFLARE_MOCK_DATA || {},
    sqlQuery: 'SELECT id, title, lang, author, downloads FROM code_snippets WHERE status = \'published\' LIMIT 10;',
    sqlResults: null,
    d1ConsoleRunning: false,
    r2FolderPrefix: '',
  };

  // DOM Elements Cache
  const appContainer = document.getElementById('cf-app-viewport');
  const sidebarContainer = document.getElementById('cf-sidebar-container');
  const headerContainer = document.getElementById('cf-header-container');
  const modalContainer = document.getElementById('cf-modal-container');
  const quickSearchBackdrop = document.getElementById('cf-quick-search-backdrop');
  const quickSearchInput = document.getElementById('cf-quick-search-input');
  const quickSearchResults = document.getElementById('cf-quick-search-results');
  const customizerDrawer = document.getElementById('cf-customizer-drawer');
  const customizerBackdrop = document.getElementById('cf-customizer-backdrop');
  const mobileBackdrop = document.getElementById('cf-mobile-backdrop');

  // Format Helper
  function formatBytes(bytes) {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }

  // 1-Click Clipboard Copy
  window.copyToClipboard = function (text, key) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(text);
    } else {
      const ta = document.createElement('textarea');
      ta.value = text;
      ta.style.position = 'fixed';
      ta.style.opacity = '0';
      document.body.appendChild(ta);
      ta.select();
      document.execCommand('copy');
      document.body.removeChild(ta);
    }
    state.copiedKey = key;
    render();
    setTimeout(() => {
      state.copiedKey = null;
      render();
    }, 2000);
  };

  // Tab & Sub-Tab Switcher
  window.switchTab = function (tab, subTab = 'overview') {
    state.currentTab = tab;
    state.subTab = subTab;
    window.location.hash = tab;
    if (state.mobileSidebarOpen) toggleMobileSidebar(false);
    render();
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  window.switchSubTab = function (subTab) {
    state.subTab = subTab;
    render();
  };

  window.toggleGroup = function (groupName) {
    state.openGroups[groupName] = !state.openGroups[groupName];
    renderSidebar();
  };

  // Modals Manager
  window.openModal = function (modalName) {
    state.activeModal = modalName;
    renderModals();
  };

  window.closeModal = function () {
    state.activeModal = null;
    renderModals();
  };

  // Live Theme Customizer Controls
  window.toggleCustomizer = function (open) {
    state.customizerOpen = (open !== undefined) ? open : !state.customizerOpen;
    if (customizerDrawer && customizerBackdrop) {
      if (state.customizerOpen) {
        customizerDrawer.classList.add('open');
        customizerBackdrop.classList.add('open');
      } else {
        customizerDrawer.classList.remove('open');
        customizerBackdrop.classList.remove('open');
      }
    }
  };

  window.setAccentColor = function (color) {
    state.accent = color;
    document.documentElement.setAttribute('data-accent', color);
    localStorage.setItem('cf_accent', color);
    document.querySelectorAll('.color-swatch-btn').forEach(btn => {
      if (btn.getAttribute('data-color') === color) {
        btn.classList.add('active');
        btn.textContent = '✓';
      } else {
        btn.classList.remove('active');
        btn.textContent = '';
      }
    });
  };

  window.setThemeMode = function (theme) {
    state.theme = theme;
    localStorage.setItem('cf_theme', theme);
    applyTheme();
    renderHeader();
  };

  window.toggleTheme = function () {
    const nextTheme = state.theme === 'dark' ? 'light' : 'dark';
    window.setThemeMode(nextTheme);
  };

  function applyTheme() {
    if (state.theme === 'dark') {
      document.documentElement.classList.add('dark');
      document.body.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
      document.body.classList.remove('dark');
    }
    const lightBtn = document.getElementById('theme-light-btn');
    const darkBtn = document.getElementById('theme-dark-btn');
    if (lightBtn && darkBtn) {
      if (state.theme === 'dark') {
        darkBtn.classList.add('btn-cf-primary');
        darkBtn.classList.remove('btn-cf-secondary');
        lightBtn.classList.add('btn-cf-secondary');
        lightBtn.classList.remove('btn-cf-primary');
      } else {
        lightBtn.classList.add('btn-cf-primary');
        lightBtn.classList.remove('btn-cf-secondary');
        darkBtn.classList.add('btn-cf-secondary');
        darkBtn.classList.remove('btn-cf-primary');
      }
    }
  }

  window.setSidebarMode = function (mode) {
    state.sidebarMini = (mode === 'mini');
    localStorage.setItem('cf_sidebar_mini', state.sidebarMini);
    if (state.sidebarMini) {
      document.body.classList.add('sidebar-mini');
    } else {
      document.body.classList.remove('sidebar-mini');
    }
    const fullBtn = document.getElementById('sidebar-full-btn');
    const miniBtn = document.getElementById('sidebar-mini-btn');
    if (fullBtn && miniBtn) {
      if (state.sidebarMini) {
        miniBtn.classList.add('btn-cf-primary');
        miniBtn.classList.remove('btn-cf-secondary');
        fullBtn.classList.add('btn-cf-secondary');
        fullBtn.classList.remove('btn-cf-primary');
      } else {
        fullBtn.classList.add('btn-cf-primary');
        fullBtn.classList.remove('btn-cf-secondary');
        miniBtn.classList.add('btn-cf-secondary');
        miniBtn.classList.remove('btn-cf-primary');
      }
    }
  };

  window.toggleMobileSidebar = function (open) {
    state.mobileSidebarOpen = (open !== undefined) ? open : !state.mobileSidebarOpen;
    const sb = document.getElementById('cf-sidebar-container');
    if (sb && mobileBackdrop) {
      if (state.mobileSidebarOpen) {
        sb.classList.add('mobile-open');
        mobileBackdrop.classList.add('active');
      } else {
        sb.classList.remove('mobile-open');
        mobileBackdrop.classList.remove('active');
      }
    }
  };

  window.resetCustomizer = function () {
    window.setAccentColor('orange');
    window.setThemeMode('light');
    window.setSidebarMode('full');
    window.toggleCustomizer(false);
  };

  // --------------------------------------------------------------------------
  // Master Render Engine
  // --------------------------------------------------------------------------
  function render() {
    renderSidebar();
    renderHeader();
    renderMainContent();
    renderModals();
  }

  // --------------------------------------------------------------------------
  // Component 1: Sidebar Navigation (Exact Match to Production Inspection)
  // --------------------------------------------------------------------------
  function renderSidebar() {
    if (!sidebarContainer) return;
    const data = state.data;
    const open = state.openGroups;
    const tab = state.currentTab;

    sidebarContainer.innerHTML = `
      <div class="sidebar-header">
        <a href="#" class="sidebar-logo" onclick="switchTab('home'); return false;">
          <img src="assets/cloudflare-logo.svg" alt="Cloudflare logo" style="height: 22px; width: auto;" />
          <span style="margin-left: 0.25rem;">Cloudflare</span>
        </a>
      </div>

      <div class="account-selector" onclick="alert('Account: ${data.account.name}\\nAccount ID: ${data.account.id}')">
        <div class="account-name-group">
          <span class="account-avatar-dot"></span>
          <span class="account-name-text">${data.account.name}</span>
        </div>
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m6 9 6 6 6-6"/></svg>
      </div>

      <div class="sidebar-nav">
        <!-- Quick Search Button -->
        <button class="sidebar-search-btn" onclick="openQuickSearch()">
          ${window.KumoIcons.search()}
          <span>Quick search...</span>
          <kbd>⌘K</kbd>
        </button>

        <!-- Top Level Nav: Account Home -->
        <button class="kumo-nav-item ${tab === 'home' ? 'active' : ''}" onclick="switchTab('home')">
          <svg width="16" height="16" viewBox="0 0 16 16" class="fill-current opacity-60 shrink-0">
            <path d="M8 1L1 7h2v7h4v-4h2v4h4V7h2L8 1z"/>
          </svg>
          <span>Account home</span>
        </button>

        <!-- Recents Accordion -->
        <button class="kumo-nav-group-btn ${open.recents ? 'expanded' : ''}" onclick="toggleGroup('recents')">
          <div class="kumo-nav-group-left">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
            </svg>
            <span>Recents</span>
          </div>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        ${open.recents ? `
          <div class="kumo-nav-sub-menu">
            <button class="kumo-nav-sub-item ${tab === 'worker-detail' ? 'active' : ''}" onclick="switchTab('worker-detail')">
              <span class="truncate">⚡ code-collection</span>
              <span class="badge-tag-beta">Worker</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'd1-detail' ? 'active' : ''}" onclick="switchTab('d1-detail')">
              <span class="truncate">🗄️ code_collection</span>
              <span class="badge-tag-beta">D1</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'r2-detail' ? 'active' : ''}" onclick="switchTab('r2-detail')">
              <span class="truncate">📦 code-collection</span>
              <span class="badge-tag-beta">R2</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'zone-overview' ? 'active' : ''}" onclick="switchTab('zone-overview')">
              <span class="truncate">🌐 northlab.dev</span>
              <span class="badge-tag-beta">Zone</span>
            </button>
          </div>
        ` : ''}

        <!-- SECTION: BUILD -->
        <div class="sidebar-section-divider">Build</div>

        <!-- Compute Accordion (Workers & Pages) -->
        <button class="kumo-nav-group-btn ${open.compute ? 'expanded' : ''}" onclick="toggleGroup('compute')">
          <div class="kumo-nav-group-left">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
            </svg>
            <span>Compute</span>
          </div>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        ${open.compute ? `
          <div class="kumo-nav-sub-menu">
            <button class="kumo-nav-sub-item ${tab === 'workers' ? 'active' : ''}" onclick="switchTab('workers')">
              <span>Workers & Pages</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'worker-detail' ? 'active' : ''}" onclick="switchTab('worker-detail')">
              <span>Worker Service (code-collection)</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'ai' ? 'active' : ''}" onclick="switchTab('ai')">
              <span>Workers AI & AI Gateway</span>
              <span class="badge-tag-beta">AI</span>
            </button>
            <button class="kumo-nav-sub-item" onclick="alert('Containers & Sandboxes isolation is active under Linux cgroups v2.')">
              <span>Containers & Sandboxes</span>
            </button>
            <button class="kumo-nav-sub-item" onclick="alert('Queues & Workflows: code-collection-imports queue active.')">
              <span>Queues & Workflows</span>
              <span class="badge-tag-beta">Beta</span>
            </button>
          </div>
        ` : ''}

        <!-- Storage & Databases Accordion (D1, R2, KV, Hyperdrive) -->
        <button class="kumo-nav-group-btn ${open.storage ? 'expanded' : ''}" onclick="toggleGroup('storage')">
          <div class="kumo-nav-group-left">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/><path d="M3 12c0 1.66 4 3 9 3s9-1.34 9-3"/>
            </svg>
            <span>Storage & databases</span>
          </div>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        ${open.storage ? `
          <div class="kumo-nav-sub-menu">
            <button class="kumo-nav-sub-item ${tab === 'd1' ? 'active' : ''}" onclick="switchTab('d1')">
              <span>D1 SQLite Databases</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'd1-detail' ? 'active' : ''}" onclick="switchTab('d1-detail')">
              <span>D1 Database (code_collection)</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'r2' ? 'active' : ''}" onclick="switchTab('r2')">
              <span>R2 Object Storage</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'r2-detail' ? 'active' : ''}" onclick="switchTab('r2-detail')">
              <span>R2 Bucket (code-collection)</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'kv' ? 'active' : ''}" onclick="switchTab('kv')">
              <span>Workers KV Store</span>
              <span class="badge-tag-beta">Global</span>
            </button>
            <button class="kumo-nav-sub-item" onclick="alert('Hyperdrive Postgres Connection Pool is bound to code-collection.')">
              <span>Hyperdrive (Postgres/MySQL)</span>
            </button>
          </div>
        ` : ''}

        <!-- SECTION: PROTECT & CONNECT -->
        <div class="sidebar-section-divider">Protect & connect</div>

        <!-- Networking & DNS Accordion -->
        <button class="kumo-nav-group-btn ${open.networking ? 'expanded' : ''}" onclick="toggleGroup('networking')">
          <div class="kumo-nav-group-left">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
            </svg>
            <span>Domains & Networking</span>
          </div>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        ${open.networking ? `
          <div class="kumo-nav-sub-menu">
            <button class="kumo-nav-sub-item ${tab === 'zone-overview' ? 'active' : ''}" onclick="switchTab('zone-overview')">
              <span>Zone Overview (northlab.dev)</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'dns' ? 'active' : ''}" onclick="switchTab('dns')">
              <span>DNS Records & Management</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'ssl' ? 'active' : ''}" onclick="switchTab('ssl')">
              <span>SSL/TLS Certificates</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'caching' ? 'active' : ''}" onclick="switchTab('caching')">
              <span>Caching & Optimization</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'rules' ? 'active' : ''}" onclick="switchTab('rules')">
              <span>Transform & Page Rules</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'network' ? 'active' : ''}" onclick="switchTab('network')">
              <span>Network (HTTP/3 & gRPC)</span>
            </button>
          </div>
        ` : ''}

        <!-- Zero Trust & Security -->
        <button class="kumo-nav-group-btn ${open.protect ? 'expanded' : ''}" onclick="toggleGroup('protect')">
          <div class="kumo-nav-group-left">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
            <span>Zero Trust & Security</span>
          </div>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        ${open.protect ? `
          <div class="kumo-nav-sub-menu">
            <button class="kumo-nav-sub-item ${tab === 'zero-trust' ? 'active' : ''}" onclick="switchTab('zero-trust')">
              <span>Access & Tunnels</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'waf' ? 'active' : ''}" onclick="switchTab('waf')">
              <span>WAF & Security Rules</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'turnstile' ? 'active' : ''}" onclick="switchTab('turnstile')">
              <span>Turnstile CAPTCHA</span>
              <span class="badge-tag-beta">Privacy</span>
            </button>
          </div>
        ` : ''}

        <!-- SECTION: OBSERVABILITY -->
        <div class="sidebar-section-divider">Observability</div>

        <button class="kumo-nav-group-btn ${open.observability ? 'expanded' : ''}" onclick="toggleGroup('observability')">
          <div class="kumo-nav-group-left">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
            </svg>
            <span>Telemetry & Logs</span>
          </div>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        ${open.observability ? `
          <div class="kumo-nav-sub-menu">
            <button class="kumo-nav-sub-item ${tab === 'analytics' ? 'active' : ''}" onclick="switchTab('analytics')">
              <span>Account Analytics</span>
            </button>
            <button class="kumo-nav-sub-item" onclick="switchTab('analytics')">
              <span>Edge Proxy Telemetry</span>
            </button>
          </div>
        ` : ''}

        <!-- SECTION: MANAGE ACCOUNT -->
        <div class="sidebar-section-divider">Manage account</div>

        <button class="kumo-nav-group-btn ${open.account ? 'expanded' : ''}" onclick="toggleGroup('account')">
          <div class="kumo-nav-group-left">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <circle cx="12" cy="7" r="4"/><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            </svg>
            <span>Account settings</span>
          </div>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        ${open.account ? `
          <div class="kumo-nav-sub-menu">
            <button class="kumo-nav-sub-item ${tab === 'team' ? 'active' : ''}" onclick="switchTab('team')">
              <span>Members & RBAC</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'billing' ? 'active' : ''}" onclick="switchTab('billing')">
              <span>Billing & Subscriptions</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'tokens' ? 'active' : ''}" onclick="switchTab('tokens')">
              <span>API Tokens & Keys</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'audit-logs' ? 'active' : ''}" onclick="switchTab('audit-logs')">
              <span>Audit Logs Stream</span>
            </button>
            <button class="kumo-nav-sub-item ${tab === 'settings' ? 'active' : ''}" onclick="switchTab('settings')">
              <span>Profile & Preferences</span>
            </button>
          </div>
        ` : ''}

        <!-- SECTION: ENVATO UI KIT -->
        <div class="sidebar-section-divider">Template kit</div>

        <button class="kumo-nav-item ${tab === 'components' ? 'active' : ''}" onclick="switchTab('components')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="fill-current opacity-60 shrink-0">
            <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/>
          </svg>
          <span>UI Components Kit</span>
        </button>
      </div>

      <!-- Sidebar Footer -->
      <div class="sidebar-footer">
        <div class="operational-status-bar">
          <div style="display: flex; align-items: center; gap: 0.45rem;">
            <div class="status-dot"></div>
            <span>All Systems Normal</span>
          </div>
          <span style="font-family: var(--font-mono); font-size: 0.72rem;">v2026.11</span>
        </div>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // Component 2: Top Header Bar
  // --------------------------------------------------------------------------
  function renderHeader() {
    if (!headerContainer) return;
    const data = state.data;

    headerContainer.innerHTML = `
      <div class="header-left">
        <!-- Mobile Hamburger Button -->
        <button class="mobile-menu-btn" onclick="toggleMobileSidebar(true)" aria-label="Toggle Mobile Sidebar">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="18" x2="21" y2="18"/>
          </svg>
        </button>

        <div class="header-search-box" onclick="openQuickSearch()">
          ${window.KumoIcons.search()}
          <input type="text" placeholder="Search resources, apps, databases... (⌘K)" readonly />
          <kbd>⌘K</kbd>
        </div>
      </div>

      <div class="header-right">
        <!-- Status Pill -->
        <div class="cf-pill active" style="font-size: 0.8rem; padding: 0.25rem 0.75rem;">
          <span class="status-dot" style="width: 6px; height: 6px;"></span>
          <span>100% Operational</span>
        </div>

        <!-- Customizer Button -->
        <button class="header-action-btn" onclick="toggleCustomizer(true)" title="Live Theme Customizer">
          ⚙️
        </button>

        <!-- Dark/Light Theme Switcher -->
        <button class="header-action-btn" onclick="toggleTheme()" title="Switch Light/Dark Mode">
          ${state.theme === 'dark' ? '☀️' : '🌙'}
        </button>

        <!-- Notification Bell -->
        <button class="header-action-btn" onclick="alert('No unread alerts. Your cluster is completely healthy.')" title="Notifications">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>
          <span class="header-badge-dot"></span>
        </button>

        <!-- User Profile Pill -->
        <div class="header-user-pill" onclick="switchTab('settings')">
          <div class="user-avatar-circle">H</div>
          <span style="font-weight: 500;">${data.account.email}</span>
        </div>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // Main Router Engine
  // --------------------------------------------------------------------------
  function renderMainContent() {
    if (!appContainer) return;

    switch (state.currentTab) {
      case 'home':
        appContainer.innerHTML = renderHomeView();
        break;
      case 'zone-overview':
        appContainer.innerHTML = renderZoneOverviewView();
        break;
      case 'dns':
        appContainer.innerHTML = renderDnsView();
        break;
      case 'waf':
        appContainer.innerHTML = renderWafView();
        break;
      case 'ssl':
        appContainer.innerHTML = renderSslView();
        break;
      case 'caching':
        appContainer.innerHTML = renderCachingView();
        break;
      case 'rules':
        appContainer.innerHTML = renderRulesView();
        break;
      case 'network':
        appContainer.innerHTML = renderNetworkView();
        break;
      case 'workers':
        appContainer.innerHTML = renderWorkersView();
        break;
      case 'worker-detail':
        appContainer.innerHTML = renderWorkerDetailView();
        break;
      case 'd1':
        appContainer.innerHTML = renderD1View();
        break;
      case 'd1-detail':
        appContainer.innerHTML = renderD1DetailView();
        break;
      case 'r2':
        appContainer.innerHTML = renderR2View();
        break;
      case 'r2-detail':
        appContainer.innerHTML = renderR2DetailView();
        break;
      case 'analytics':
        appContainer.innerHTML = renderAnalyticsView();
        break;
      case 'zero-trust':
        appContainer.innerHTML = renderZeroTrustView();
        break;
      case 'team':
        appContainer.innerHTML = renderTeamView();
        break;
      case 'billing':
        appContainer.innerHTML = renderBillingView();
        break;
      case 'tokens':
        appContainer.innerHTML = renderTokensView();
        break;
      case 'audit-logs':
        appContainer.innerHTML = renderAuditLogsView();
        break;
      case 'settings':
        appContainer.innerHTML = renderSettingsView();
        break;
      case 'components':
        appContainer.innerHTML = renderComponentsView();
        break;
      case 'ai':
        appContainer.innerHTML = renderAiView();
        break;
      case 'kv':
        appContainer.innerHTML = renderKvView();
        break;
      case 'turnstile':
        appContainer.innerHTML = renderTurnstileView();
        break;
      default:
        appContainer.innerHTML = renderHomeView();
    }
  }

  // --------------------------------------------------------------------------
  // View 1: Account Home (Exact Match to Production Inspection)
  // --------------------------------------------------------------------------
  function renderHomeView() {
    const data = state.data;
    return `
      <div class="breadcrumbs">
        <span>Account</span>
        <span class="divider">/</span>
        <span class="current">Home</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">A new day ahead.</h1>
          <p class="page-subtitle">Unified overview of your domains, serverless workers, databases, and global edge security.</p>
        </div>
        <div style="display: flex; gap: 0.5rem;">
          <button class="btn-cf btn-cf-primary" onclick="openModal('create_app')">+ Create Worker</button>
          <button class="btn-cf btn-cf-secondary" onclick="openModal('create_db')">+ Add Database</button>
        </div>
      </div>

      <!-- Quick Action Tiles Banner -->
      <div class="quick-actions-banner">
        <div class="action-tile" onclick="switchTab('dns')">
          <div class="action-tile-icon">🌐</div>
          <div>
            <div class="action-tile-title">Manage DNS Zones</div>
            <div class="action-tile-desc">Configure Anycast DNS records for northlab.dev</div>
          </div>
        </div>
        <div class="action-tile" onclick="switchTab('worker-detail')">
          <div class="action-tile-icon">⚡</div>
          <div>
            <div class="action-tile-title">Deploy Worker</div>
            <div class="action-tile-desc">Push zero-downtime updates to code-collection</div>
          </div>
        </div>
        <div class="action-tile" onclick="switchTab('d1-detail')">
          <div class="action-tile-icon">🗄️</div>
          <div>
            <div class="action-tile-title">Query SQLite D1</div>
            <div class="action-tile-desc">Run SQL console queries on code_collection</div>
          </div>
        </div>
        <div class="action-tile" onclick="switchTab('r2-detail')">
          <div class="action-tile-icon">📦</div>
          <div>
            <div class="action-tile-title">R2 Object Storage</div>
            <div class="action-tile-desc">Browse files with zero egress fees worldwide</div>
          </div>
        </div>
      </div>

      <!-- Live 24h Metrics Tiles Strip (Exact Match to Cloudflare Dashboard) -->
      <div class="stats-overview-grid">
        <div class="metric-card">
          <div class="metric-header">
            <span class="metric-label">Total Requests (24h)</span>
            <span class="metric-delta negative">${data.analytics.totalRequestsDelta}</span>
          </div>
          <div class="metric-value">${data.analytics.totalRequests}</div>
          <svg width="100%" height="24" viewBox="0 0 100 24" fill="none" style="margin: 4px 0;">
            <path d="M0 18 Q 20 6, 45 15 T 100 6" stroke="var(--cf-brand-orange)" stroke-width="2" fill="none" stroke-linecap="round"/>
          </svg>
          <div class="metric-footer">Across ${data.domains.length} active zones</div>
        </div>

        <div class="metric-card">
          <div class="metric-header">
            <span class="metric-label">Worker Invocations</span>
            <span class="metric-delta positive">${data.analytics.workerInvocationsDelta}</span>
          </div>
          <div class="metric-value">${data.analytics.workerInvocations}</div>
          <svg width="100%" height="24" viewBox="0 0 100 24" fill="none" style="margin: 4px 0;">
            <path d="M0 20 L 25 16 L 45 18 L 70 8 L 100 4" stroke="#059669" stroke-width="2" fill="none" stroke-linecap="round"/>
          </svg>
          <div class="metric-footer">Errors: ${data.analytics.workerErrors} (0.0%)</div>
        </div>

        <div class="metric-card">
          <div class="metric-header">
            <span class="metric-label">Cache Hit Rate</span>
            <span class="metric-delta positive">${data.analytics.cacheHitRateDelta}</span>
          </div>
          <div class="metric-value">${data.analytics.cacheHitRate}</div>
          <svg width="100%" height="24" viewBox="0 0 100 24" fill="none" style="margin: 4px 0;">
            <path d="M0 16 Q 30 22, 60 10 T 100 2" stroke="#0051c3" stroke-width="2" fill="none" stroke-linecap="round"/>
          </svg>
          <div class="metric-footer">Edge global caching</div>
        </div>

        <div class="metric-card">
          <div class="metric-header">
            <span class="metric-label">CPU Time P90</span>
            <span class="metric-delta positive">${data.analytics.cpuTimeDelta}</span>
          </div>
          <div class="metric-value">${data.analytics.cpuTimeP90}</div>
          <svg width="100%" height="24" viewBox="0 0 100 24" fill="none" style="margin: 4px 0;">
            <path d="M0 8 L 30 16 L 55 12 L 80 20 L 100 15" stroke="#7c3aed" stroke-width="2" fill="none" stroke-linecap="round"/>
          </svg>
          <div class="metric-footer">Limit: 50 ms budget</div>
        </div>
      </div>

      <!-- Active Domains / Zones Section -->
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; margin-top: 2rem;">
        <h2 style="font-size: 1.15rem; font-weight: 600;">Active Domains & Zones (${data.domains.length})</h2>
        <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="alert('Add Domain wizard is available for commercial buyers.')">+ Add Domain</button>
      </div>

      <div class="domain-cards-grid">
        ${data.domains.map(d => `
          <div class="domain-card" onclick="switchTab('zone-overview')">
            <div class="domain-card-head">
              <div>
                <div class="domain-card-name">${d.name}</div>
                <span class="cf-pill active" style="margin-top: 0.35rem;">
                  <span class="status-dot"></span>
                  <span>${d.status}</span>
                </span>
              </div>
              <span class="badge-tag-beta">${d.plan}</span>
            </div>
            <div class="domain-card-stats">
              <div><strong>${d.requests_24h}</strong> requests (24h)</div>
              <div><strong>${d.bandwidth_24h}</strong> bandwidth</div>
              <div>SSL: <strong>${d.ssl}</strong></div>
            </div>
          </div>
        `).join('')}
      </div>

      <!-- Active Compute Workers Section -->
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; margin-top: 2rem;">
        <h2 style="font-size: 1.15rem; font-weight: 600;">Serverless Workers & Compute (${data.workers.length})</h2>
        <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="switchTab('workers')">View All Workers →</button>
      </div>

      <div class="item-cards-grid">
        ${data.workers.map(w => `
          <div class="kumo-item-card" onclick="switchTab('worker-detail')">
            <div class="kumo-card-primary">
              <div class="item-title-row">
                <div class="item-name">⚡ ${w.name}</div>
                <span class="cf-pill active"><span class="status-dot"></span><span>${w.status}</span></span>
              </div>
              <div class="item-meta-row">
                <span style="font-family: var(--font-mono); font-size: 0.8rem; color: var(--cf-brand-orange);">${w.custom_domain}</span>
                <span>•</span>
                <span>${w.repo} (${w.branch})</span>
              </div>
              <div style="font-size: 0.82rem; color: var(--text-color-kumo-subtle); margin-top: 0.5rem;">
                Deployed: ${w.last_deployed} (${w.current_version})
              </div>
            </div>
            <div class="kumo-card-secondary">
              <div class="secondary-stat-col">
                <span class="sec-label">Invocations</span>
                <span class="sec-val">${w.invocations_24h}</span>
              </div>
              <div class="secondary-stat-col">
                <span class="sec-label">CPU P90</span>
                <span class="sec-val">${w.median_cpu_time}</span>
              </div>
              <div class="secondary-stat-col">
                <span class="sec-label">Errors</span>
                <span class="sec-val">${w.errors_24h}</span>
              </div>
            </div>
          </div>
        `).join('')}
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 2: Zone Overview (northlab.dev)
  // --------------------------------------------------------------------------
  function renderZoneOverviewView() {
    const zone = state.data.activeZone;
    return `
      <div class="breadcrumbs">
        <a href="#" onclick="switchTab('home'); return false;">Account</a>
        <span class="divider">/</span>
        <span class="current">${zone.name}</span>
      </div>

      <div class="page-title-row">
        <div>
          <div style="display: flex; align-items: center; gap: 0.75rem;">
            <h1 class="page-title">${zone.name}</h1>
            <span class="cf-pill active"><span class="status-dot"></span><span>${zone.status}</span></span>
            <span class="badge-tag-beta">${zone.plan}</span>
          </div>
          <p class="page-subtitle">Zone-scoped edge configurations, Anycast DNS, SSL/TLS certificates, and edge caching.</p>
        </div>
        <div style="display: flex; gap: 0.5rem;">
          <button class="btn-cf btn-cf-primary" onclick="switchTab('dns')">Manage DNS Records</button>
          <button class="btn-cf btn-cf-secondary" onclick="openModal('purge_cache')">Purge Cache</button>
        </div>
      </div>

      <!-- Domain Sub-Navigation Tabs Bar -->
      <div class="cf-tab-nav">
        <button class="cf-tab-btn ${state.subTab === 'overview' ? 'active' : ''}" onclick="switchSubTab('overview')">Overview</button>
        <button class="cf-tab-btn" onclick="switchTab('dns')">DNS Records</button>
        <button class="cf-tab-btn" onclick="switchTab('ssl')">SSL/TLS</button>
        <button class="cf-tab-btn" onclick="switchTab('waf')">Security (WAF)</button>
        <button class="cf-tab-btn" onclick="switchTab('caching')">Caching</button>
        <button class="cf-tab-btn" onclick="switchTab('rules')">Rules</button>
        <button class="cf-tab-btn" onclick="switchTab('network')">Network</button>
      </div>

      <!-- Domain Overview Content Cards -->
      <div class="stats-overview-grid">
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">SSL Encryption Mode</span></div>
          <div class="metric-value" style="font-size: 1.4rem;">${zone.ssl.mode}</div>
          <div class="metric-footer">Universal SSL Active</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">DNS Records Count</span></div>
          <div class="metric-value">${zone.dnsRecords.length}</div>
          <div class="metric-footer">Anycast Global DNS</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Cache Level</span></div>
          <div class="metric-value" style="font-size: 1.4rem;">${zone.caching.cachingLevel}</div>
          <div class="metric-footer">TTL: ${zone.caching.browserCacheTtl}</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">HTTP/3 (QUIC)</span></div>
          <div class="metric-value" style="color: #059669;">Active</div>
          <div class="metric-footer">0-RTT Resumption On</div>
        </div>
      </div>

      <!-- Quick DNS Records Preview Table -->
      <div class="kumo-table-container" style="margin-top: 1.5rem;">
        <div class="table-header-title">
          <span>Configured DNS Records for ${zone.name}</span>
          <button class="btn-cf btn-cf-primary btn-cf-sm" onclick="switchTab('dns')">Open Full DNS Manager →</button>
        </div>
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Type</th>
              <th>Name</th>
              <th>Content</th>
              <th>Proxy Status</th>
              <th>TTL</th>
            </tr>
          </thead>
          <tbody>
            ${zone.dnsRecords.slice(0, 4).map(r => `
              <tr>
                <td><span class="dns-type-badge">${r.type}</span></td>
                <td><strong>${r.name}</strong></td>
                <td><code>${r.content}</code></td>
                <td>
                  <span class="dns-status-cloud ${r.proxied ? 'proxied' : 'dns-only'}">
                    <span>☁️</span>
                    <span>${r.proxied ? 'Proxied' : 'DNS only'}</span>
                  </span>
                </td>
                <td>${r.ttl}</td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 3: DNS Records Management (northlab.dev)
  // --------------------------------------------------------------------------
  function renderDnsView() {
    const zone = state.data.activeZone;
    return `
      <div class="breadcrumbs">
        <a href="#" onclick="switchTab('home'); return false;">Account</a>
        <span class="divider">/</span>
        <a href="#" onclick="switchTab('zone-overview'); return false;">${zone.name}</a>
        <span class="divider">/</span>
        <span class="current">DNS Records</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">DNS Management for ${zone.name}</h1>
          <p class="page-subtitle">Configure authoritative DNS records with global Anycast routing and automatic DDoS protection.</p>
        </div>
        <div style="display: flex; gap: 0.5rem;">
          <button class="btn-cf btn-cf-primary" onclick="openModal('add_dns_record')">+ Add Record</button>
          <button class="btn-cf btn-cf-secondary" onclick="alert('BIND Zone exported successfully.')">Export BIND Zone</button>
        </div>
      </div>

      <!-- Domain Sub-Navigation Tabs Bar -->
      <div class="cf-tab-nav">
        <button class="cf-tab-btn" onclick="switchTab('zone-overview')">Overview</button>
        <button class="cf-tab-btn active" onclick="switchTab('dns')">DNS Records</button>
        <button class="cf-tab-btn" onclick="switchTab('ssl')">SSL/TLS</button>
        <button class="cf-tab-btn" onclick="switchTab('waf')">Security (WAF)</button>
        <button class="cf-tab-btn" onclick="switchTab('caching')">Caching</button>
        <button class="cf-tab-btn" onclick="switchTab('rules')">Rules</button>
        <button class="cf-tab-btn" onclick="switchTab('network')">Network</button>
      </div>

      <!-- Cloudflare DNS Notice Banner -->
      <div class="ui-alert-banner info">
        <div style="font-size: 1.15rem;">ℹ️</div>
        <div>
          <strong>Cloudflare Anycast Nameservers are Active</strong>
          <div style="font-size: 0.8rem; margin-top: 0.2rem;">Your DNS records are globally resolved from 330+ edge locations with 100% SLA and zero propagation delay.</div>
        </div>
      </div>

      <!-- DNS Records Table -->
      <div class="kumo-table-container">
        <div class="table-header-title">
          <span>DNS Records (${zone.dnsRecords.length})</span>
          <input type="text" class="cf-input" placeholder="Search DNS records..." style="width: 240px; padding: 0.35rem 0.65rem; font-size: 0.82rem;" />
        </div>
        <table class="kumo-table">
          <thead>
            <tr>
              <th style="width: 80px;">Type</th>
              <th>Name</th>
              <th>Content / Value</th>
              <th>Proxy Status</th>
              <th>TTL</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            ${zone.dnsRecords.map(r => `
              <tr>
                <td><span class="dns-type-badge">${r.type}</span></td>
                <td><strong style="color: var(--text-color-kumo-strong);">${r.name}</strong></td>
                <td><code>${r.content}</code></td>
                <td>
                  <span class="dns-status-cloud ${r.proxied ? 'proxied' : 'dns-only'}" onclick="alert('Toggle Proxy status for ${r.name}')" title="Click to toggle proxy">
                    <span>☁️</span>
                    <span>${r.proxied ? 'Proxied' : 'DNS only'}</span>
                  </span>
                </td>
                <td>${r.ttl}</td>
                <td>
                  <div style="display: flex; gap: 0.35rem;">
                    <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="alert('Editing record: ${r.name}')">Edit</button>
                    <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="alert('Delete record: ${r.name}')" style="color: #dc2626;">Delete</button>
                  </div>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 4: Worker Service Detail (code-collection)
  // --------------------------------------------------------------------------
  function renderWorkerDetailView() {
    const worker = state.data.workers[0]; // code-collection
    const subTab = state.subTab;

    return `
      <div class="breadcrumbs">
        <a href="#" onclick="switchTab('workers'); return false;">Compute</a>
        <span class="divider">/</span>
        <span class="current">${worker.name}</span>
      </div>

      <div class="page-title-row">
        <div>
          <div style="display: flex; align-items: center; gap: 0.75rem;">
            <h1 class="page-title">⚡ ${worker.name}</h1>
            <span class="cf-pill active"><span class="status-dot"></span><span>${worker.status}</span></span>
            <span class="badge-tag-beta">Production</span>
          </div>
          <p class="page-subtitle">Serverless Worker running on V8 Isolates with instant global cold starts & zero cold start latency.</p>
        </div>
        <div style="display: flex; gap: 0.5rem;">
          <a href="https://${worker.custom_domain}" target="_blank" class="btn-cf btn-cf-secondary">Visit Site ↗</a>
          <button class="btn-cf btn-cf-primary" onclick="alert('Triggering zero-downtime blue/green deployment via wrangler...')">Deploy Version</button>
        </div>
      </div>

      <!-- Worker Top Tabs Navigation (Exact Match to Production Inspection) -->
      <div class="cf-tab-nav">
        <button class="cf-tab-btn ${subTab === 'overview' ? 'active' : ''}" onclick="switchSubTab('overview')">Overview</button>
        <button class="cf-tab-btn ${subTab === 'metrics' ? 'active' : ''}" onclick="switchSubTab('metrics')">Metrics</button>
        <button class="cf-tab-btn ${subTab === 'deployments' ? 'active' : ''}" onclick="switchSubTab('deployments')">Deployments</button>
        <button class="cf-tab-btn ${subTab === 'logs' ? 'active' : ''}" onclick="switchSubTab('logs')">Observability & Live Logs</button>
        <button class="cf-tab-btn ${subTab === 'settings' ? 'active' : ''}" onclick="switchSubTab('settings')">Settings & Bindings</button>
      </div>

      ${subTab === 'overview' ? `
        <!-- Overview Tab -->
        <div class="stats-overview-grid">
          <div class="metric-card">
            <div class="metric-header"><span class="metric-label">24h Invocations</span></div>
            <div class="metric-value">${worker.invocations_24h}</div>
            <div class="metric-footer">${worker.invocations_delta} vs yesterday</div>
          </div>
          <div class="metric-card">
            <div class="metric-header"><span class="metric-label">Median CPU Time</span></div>
            <div class="metric-value">${worker.median_cpu_time}</div>
            <div class="metric-footer">${worker.median_cpu_delta} execution</div>
          </div>
          <div class="metric-card">
            <div class="metric-header"><span class="metric-label">Active Version</span></div>
            <div class="metric-value" style="font-size: 1.4rem;">${worker.current_version.split(' ')[0]}</div>
            <div class="metric-footer">Git commit: ${worker.current_version.split(' ')[1]}</div>
          </div>
          <div class="metric-card">
            <div class="metric-header"><span class="metric-label">Active Routes</span></div>
            <div class="metric-value" style="font-size: 1.15rem; color: var(--cf-brand-orange);">${worker.custom_domain}</div>
            <div class="metric-footer">Subdomain: ${worker.subdomain}</div>
          </div>
        </div>

        <!-- Connected Bindings Summary Card -->
        <div class="kumo-table-container" style="margin-top: 1.5rem;">
          <div class="table-header-title">
            <span>Connected Resources & Bindings (D1, R2, Queues)</span>
            <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="switchSubTab('settings')">Configure Bindings →</button>
          </div>
          <table class="kumo-table">
            <thead>
              <tr>
                <th>Binding Variable</th>
                <th>Resource Type</th>
                <th>Connected Target</th>
              </tr>
            </thead>
            <tbody>
              ${worker.bindings.map(b => `
                <tr>
                  <td><code>env.${b.name}</code></td>
                  <td><span class="badge-tag-beta">${b.type}</span></td>
                  <td><strong>${b.target}</strong></td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      ` : subTab === 'deployments' ? `
        <!-- Deployments Tab -->
        <div class="kumo-table-container">
          <div class="table-header-title">
            <span>Deployment Release History (${worker.deployments.length})</span>
            <button class="btn-cf btn-cf-primary btn-cf-sm" onclick="alert('Rollback triggered.')">Rollback to v2.8.3</button>
          </div>
          <table class="kumo-table">
            <thead>
              <tr>
                <th>Version</th>
                <th>Git Commit</th>
                <th>Triggered By</th>
                <th>Timestamp (UTC)</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              ${worker.deployments.map(dep => `
                <tr>
                  <td><strong>${dep.version}</strong></td>
                  <td><code>${dep.commit}</code></td>
                  <td>${dep.author}</td>
                  <td>${dep.date}</td>
                  <td>
                    <span class="cf-pill ${dep.status.includes('Active') ? 'active' : 'stopped'}">
                      <span class="status-dot"></span>
                      <span>${dep.status}</span>
                    </span>
                  </td>
                  <td>
                    ${dep.status.includes('Active') ? '<span style="color: var(--text-color-kumo-subtle); font-size: 0.8rem;">Current Live</span>' : `<button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="alert('Rolling back to ${dep.version}')">Rollback</button>`}
                  </td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      ` : subTab === 'logs' ? `
        <!-- Observability & Live Logs Tab -->
        <div class="d1-sql-console" style="height: 380px; display: flex; flex-direction: column;">
          <div class="d1-sql-console-head">
            <span style="color: #f8fafc; font-weight: 600;">⚡ Realtime Log Stream — ${worker.name}</span>
            <div style="display: flex; gap: 0.5rem; align-items: center;">
              <span class="status-dot" style="background: #10b981;"></span>
              <span style="color: #94a3b8; font-size: 0.75rem;">Listening on Anycast edges...</span>
            </div>
          </div>
          <div style="flex: 1; overflow-y: auto; font-family: var(--font-mono); font-size: 0.82rem; line-height: 1.6; color: #94a3b8; padding-top: 0.5rem;">
            <div>[2026-11-22 08:52:10 UTC] [INFO] GET https://codecollection.northlab.dev/ HTTP/2 200 OK (11.2ms) [IAD Edge]</div>
            <div>[2026-11-22 08:52:14 UTC] [INFO] D1 SQL Execution: SELECT * FROM code_snippets WHERE id = 1 -> 1 row returned (0.8ms)</div>
            <div>[2026-11-22 08:52:18 UTC] [INFO] R2 Object Stream: GET /thumbnails/app-hero-light.png -> 200 OK (Cache HIT)</div>
            <div>[2026-11-22 08:52:25 UTC] [INFO] GET https://codecollection.northlab.dev/api/snippets HTTP/2 200 OK (1.4ms) [FRA Edge]</div>
            <div style="color: #38bdf8;">[2026-11-22 08:52:32 UTC] [METRIC] Memory: 14.8 MB / 128 MB | CPU: 1.1ms | Isolate healthy</div>
          </div>
        </div>
      ` : `
        <!-- Settings & Bindings Tab -->
        <div class="ui-kit-section">
          <h3 class="ui-kit-title">Runtime Environment Variables & Secrets</h3>
          <p class="ui-kit-desc">Encrypted variables injected directly into isolate memory at boot time.</p>
          <table class="kumo-table">
            <thead>
              <tr>
                <th>Variable Name</th>
                <th>Type</th>
                <th>Value</th>
              </tr>
            </thead>
            <tbody>
              ${worker.secrets.map(s => `
                <tr>
                  <td><strong>${s.key}</strong></td>
                  <td><span class="badge-tag-beta">${s.type}</span></td>
                  <td><code>${s.value}</code></td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      `}
    `;
  }

  // --------------------------------------------------------------------------
  // View 5: D1 SQLite Databases List
  // --------------------------------------------------------------------------
  function renderD1View() {
    const dbs = state.data.databases;
    return `
      <div class="breadcrumbs">
        <span>Storage & Databases</span>
        <span class="divider">/</span>
        <span class="current">D1 SQL Databases</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">D1 SQL Databases</h1>
          <p class="page-subtitle">Serverless relational SQLite databases with instant query replication to Cloudflare's global edge.</p>
        </div>
        <div style="display: flex; gap: 0.5rem;">
          <button class="btn-cf btn-cf-primary" onclick="openModal('create_db')">+ Create Database</button>
        </div>
      </div>

      <div class="kumo-table-container">
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Database Name</th>
              <th>Engine</th>
              <th>Location</th>
              <th>Size</th>
              <th>Tables</th>
              <th>24h Reads / Writes</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            ${dbs.map(d => `
              <tr onclick="switchTab('d1-detail')" style="cursor: pointer;">
                <td><strong style="color: var(--text-color-kumo-strong);">🗄️ ${d.name}</strong></td>
                <td><span class="badge-tag-beta">${d.engine}</span></td>
                <td>${d.location}</td>
                <td>${d.size}</td>
                <td>${d.tables_count} tables</td>
                <td>${d.read_units_24h} / ${d.write_units_24h}</td>
                <td>
                  <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="event.stopPropagation(); switchTab('d1-detail')">Open Console →</button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 6: D1 Database Detail (code_collection)
  // --------------------------------------------------------------------------
  function renderD1DetailView() {
    const db = state.data.databases[0]; // code_collection
    const subTab = state.subTab;

    return `
      <div class="breadcrumbs">
        <a href="#" onclick="switchTab('d1'); return false;">D1 Databases</a>
        <span class="divider">/</span>
        <span class="current">${db.name}</span>
      </div>

      <div class="page-title-row">
        <div>
          <div style="display: flex; align-items: center; gap: 0.75rem;">
            <h1 class="page-title">🗄️ ${db.name}</h1>
            <span class="cf-pill active"><span class="status-dot"></span><span>Active</span></span>
            <span class="badge-tag-beta">SQLite 3.45</span>
          </div>
          <p class="page-subtitle">UUID: <code>${db.id}</code> (Region: ${db.location})</p>
        </div>
        <div style="display: flex; gap: 0.5rem;">
          <button class="btn-cf btn-cf-secondary" onclick="copyToClipboard('${db.id}', 'db_uuid')">
            ${state.copiedKey === 'db_uuid' ? '✓ Copied UUID' : 'Copy Database UUID'}
          </button>
          <button class="btn-cf btn-cf-primary" onclick="executeSqlQuery()">Run SQL Query</button>
        </div>
      </div>

      <!-- D1 Top Tabs Navigation -->
      <div class="cf-tab-nav">
        <button class="cf-tab-btn ${subTab === 'overview' ? 'active' : ''}" onclick="switchSubTab('overview')">SQL Console</button>
        <button class="cf-tab-btn ${subTab === 'tables' ? 'active' : ''}" onclick="switchSubTab('tables')">Tables & Schema (${db.tables.length})</button>
        <button class="cf-tab-btn ${subTab === 'metrics' ? 'active' : ''}" onclick="switchSubTab('metrics')">Metrics & Units</button>
        <button class="cf-tab-btn ${subTab === 'time-travel' ? 'active' : ''}" onclick="switchSubTab('time-travel')">Time Travel Backups</button>
      </div>

      ${subTab === 'overview' ? `
        <!-- SQL Console View -->
        <div class="d1-sql-console">
          <div class="d1-sql-console-head">
            <span style="color: #f8fafc; font-weight: 600;">Interactive D1 SQL Runner</span>
            <div style="display: flex; gap: 0.5rem;">
              <button class="btn-cf btn-cf-primary btn-cf-sm" onclick="executeSqlQuery()">▶ Execute (Ctrl+Enter)</button>
            </div>
          </div>
          <textarea id="cf-sql-input" class="d1-sql-editor" onkeydown="if(event.ctrlKey && event.key==='Enter') executeSqlQuery()">${state.sqlQuery}</textarea>
        </div>

        <!-- Query Execution Results -->
        <div class="kumo-table-container" style="margin-top: 1.5rem;">
          <div class="table-header-title">
            <span>Query Results (5 rows returned in 0.8ms)</span>
            <span class="cf-pill active"><span class="status-dot"></span><span>200 OK</span></span>
          </div>
          <table class="kumo-table">
            <thead>
              <tr>
                <th>id</th>
                <th>title</th>
                <th>lang</th>
                <th>author</th>
                <th>downloads</th>
                <th>status</th>
              </tr>
            </thead>
            <tbody>
              ${db.sampleQueryResults.map(r => `
                <tr>
                  <td><code>${r.id}</code></td>
                  <td><strong>${r.title}</strong></td>
                  <td><span class="badge-tag-beta">${r.lang}</span></td>
                  <td>${r.author}</td>
                  <td>${r.downloads.toLocaleString()}</td>
                  <td><span class="cf-pill active"><span class="status-dot"></span><span>${r.status}</span></span></td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      ` : subTab === 'tables' ? `
        <!-- Tables & Schema View -->
        <div class="kumo-table-container">
          <table class="kumo-table">
            <thead>
              <tr>
                <th>Table Name</th>
                <th>Rows Count</th>
                <th>Columns</th>
                <th>Table Size</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              ${db.tables.map(t => `
                <tr>
                  <td><strong>${t.name}</strong></td>
                  <td>${t.rows.toLocaleString()}</td>
                  <td>${t.columns}</td>
                  <td>${t.size}</td>
                  <td><button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="queryTable('${t.name}')">SELECT *</button></td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      ` : `
        <!-- Metrics / Time Travel -->
        <div class="stats-overview-grid">
          <div class="metric-card">
            <div class="metric-header"><span class="metric-label">Row Reads (24h)</span></div>
            <div class="metric-value">${db.read_units_24h}</div>
          </div>
          <div class="metric-card">
            <div class="metric-header"><span class="metric-label">Row Writes (24h)</span></div>
            <div class="metric-value">${db.write_units_24h}</div>
          </div>
          <div class="metric-card">
            <div class="metric-header"><span class="metric-label">Total Queries</span></div>
            <div class="metric-value">${db.queries_count}</div>
          </div>
          <div class="metric-card">
            <div class="metric-header"><span class="metric-label">Database Size</span></div>
            <div class="metric-value">${db.size}</div>
          </div>
        </div>
      `}
    `;
  }

  // --------------------------------------------------------------------------
  // View 7: R2 Object Storage & Bucket Detail
  // --------------------------------------------------------------------------
  function renderR2View() {
    const buckets = state.data.buckets;
    return `
      <div class="breadcrumbs">
        <span>Storage & Databases</span>
        <span class="divider">/</span>
        <span class="current">R2 Object Storage</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">R2 Object Storage</h1>
          <p class="page-subtitle">S3-compatible distributed object storage with zero egress fees worldwide.</p>
        </div>
        <div style="display: flex; gap: 0.5rem;">
          <button class="btn-cf btn-cf-primary" onclick="openModal('create_bucket')">+ Create Bucket</button>
        </div>
      </div>

      <div class="kumo-table-container">
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Bucket Name</th>
              <th>Region</th>
              <th>Storage Class</th>
              <th>Size</th>
              <th>Objects</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            ${buckets.map(b => `
              <tr onclick="switchTab('r2-detail')" style="cursor: pointer;">
                <td><strong style="color: var(--text-color-kumo-strong);">📦 ${b.name}</strong></td>
                <td>${b.region}</td>
                <td><span class="badge-tag-beta">${b.storage_class}</span></td>
                <td>${b.size}</td>
                <td>${b.objects_count} objects</td>
                <td><button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="event.stopPropagation(); switchTab('r2-detail')">Browse Files →</button></td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  function renderR2DetailView() {
    const bucket = state.data.buckets[0]; // code-collection
    return `
      <div class="breadcrumbs">
        <a href="#" onclick="switchTab('r2'); return false;">R2 Storage</a>
        <span class="divider">/</span>
        <span class="current">${bucket.name}</span>
      </div>

      <div class="page-title-row">
        <div>
          <div style="display: flex; align-items: center; gap: 0.75rem;">
            <h1 class="page-title">📦 ${bucket.name}</h1>
            <span class="cf-pill active"><span class="status-dot"></span><span>Zero Egress</span></span>
          </div>
          <p class="page-subtitle">Endpoint: <code>${bucket.s3_endpoint}</code></p>
        </div>
        <div style="display: flex; gap: 0.5rem;">
          <button class="btn-cf btn-cf-primary" onclick="alert('Upload file modal.')">+ Upload Object</button>
          <button class="btn-cf btn-cf-secondary" onclick="alert('New folder created.')">+ Add Folder</button>
        </div>
      </div>

      <div class="kumo-table-container">
        <div class="table-header-title">
          <span>Files & Objects in ${bucket.name} (${bucket.objects.length})</span>
          <input type="text" class="cf-input" placeholder="Search objects by prefix..." style="width: 240px; padding: 0.35rem 0.65rem; font-size: 0.82rem;" />
        </div>
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Object Key</th>
              <th>Content Type</th>
              <th>Size</th>
              <th>Last Modified</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            ${bucket.objects.map(obj => `
              <tr>
                <td><strong>${obj.key}</strong></td>
                <td><span class="badge-tag-beta">${obj.type}</span></td>
                <td>${obj.size}</td>
                <td>${obj.modified}</td>
                <td>
                  <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="alert('Downloading ${obj.key}')">Download</button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 8: SSL / TLS Encryption
  // --------------------------------------------------------------------------
  function renderSslView() {
    const ssl = state.data.activeZone.ssl;
    return `
      <div class="breadcrumbs">
        <a href="#" onclick="switchTab('home'); return false;">Account</a>
        <span class="divider">/</span>
        <span class="current">SSL / TLS Encryption</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">SSL / TLS Encryption</h1>
          <p class="page-subtitle">Manage end-to-end edge cryptography, universal certificates, and minimum TLS protocols.</p>
        </div>
      </div>

      <div class="ssl-modes-grid">
        <div class="ssl-mode-card ${ssl.mode === 'Off' ? 'active' : ''}" onclick="setSslMode('Off')">
          <div class="ssl-mode-head">🔓 Off</div>
          <div class="ssl-mode-desc">No encryption applied. Visitors connect over plain HTTP.</div>
        </div>
        <div class="ssl-mode-card ${ssl.mode === 'Flexible' ? 'active' : ''}" onclick="setSslMode('Flexible')">
          <div class="ssl-mode-head">🔐 Flexible</div>
          <div class="ssl-mode-desc">Encrypted between visitor and Cloudflare; unencrypted to origin.</div>
        </div>
        <div class="ssl-mode-card ${ssl.mode === 'Full' ? 'active' : ''}" onclick="setSslMode('Full')">
          <div class="ssl-mode-head">🔒 Full</div>
          <div class="ssl-mode-desc">End-to-end encrypted to origin with self-signed certificate.</div>
        </div>
        <div class="ssl-mode-card ${ssl.mode === 'Full (Strict)' ? 'active' : ''}" onclick="setSslMode('Full (Strict)')">
          <div class="ssl-mode-head">🛡️ Full (Strict)</div>
          <div class="ssl-mode-desc">Enforces trusted CA certificates on origin server. Highly recommended.</div>
        </div>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 9: WAF & Security Rules
  // --------------------------------------------------------------------------
  function renderWafView() {
    return `
      <div class="breadcrumbs">
        <span>Security</span>
        <span class="divider">/</span>
        <span class="current">WAF & Rules</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Web Application Firewall (WAF)</h1>
          <p class="page-subtitle">Inspect incoming HTTP traffic with OWASP Core Rulesets, bot mitigation, and rate limiting.</p>
        </div>
        <button class="btn-cf btn-cf-primary" onclick="alert('Create Rule wizard opened.')">+ Create Firewall Rule</button>
      </div>

      <div class="stats-overview-grid">
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Threats Blocked (24h)</span></div>
          <div class="metric-value">12</div>
          <div class="metric-footer">SQLi & XSS mitigated</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Managed Ruleset</span></div>
          <div class="metric-value" style="color: #059669;">OWASP v3.3.4</div>
          <div class="metric-footer">Sensitivity: Normal</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Bot Fight Mode</span></div>
          <div class="metric-value" style="color: #059669;">Active</div>
          <div class="metric-footer">Automated bot challenges</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">DDoS Mitigation</span></div>
          <div class="metric-value" style="color: #059669;">Autonomous</div>
          <div class="metric-footer">Layer 7 mitigation</div>
        </div>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 10: Caching & Optimization
  // --------------------------------------------------------------------------
  function renderCachingView() {
    const c = state.data.activeZone.caching;
    return `
      <div class="breadcrumbs">
        <span>Optimization</span>
        <span class="divider">/</span>
        <span class="current">Caching & CDN</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Caching Configuration</h1>
          <p class="page-subtitle">Accelerate static content delivery and reduce origin load across 330+ edge data centers.</p>
        </div>
        <button class="btn-cf btn-cf-primary" onclick="openModal('purge_cache')">Purge Everything</button>
      </div>

      <div class="ui-kit-section">
        <h3 class="ui-kit-title">Tiered Cache & Argo Smart Routing</h3>
        <p class="ui-kit-desc">Uses Cloudflare's regional data centers to minimize requests to your origin server.</p>
        <span class="cf-pill active"><span class="status-dot"></span><span>${c.tieredCache}</span></span>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 11: Transform & Page Rules
  // --------------------------------------------------------------------------
  function renderRulesView() {
    return `
      <div class="breadcrumbs">
        <span>Rules</span>
        <span class="divider">/</span>
        <span class="current">Transform & Page Rules</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Edge Transform Rules</h1>
          <p class="page-subtitle">Customize HTTP request and response headers, URL rewriting, and redirect rules.</p>
        </div>
        <button class="btn-cf btn-cf-primary" onclick="alert('Create Rule modal opened.')">+ Create Rule</button>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 12: Network (HTTP/3, WebSockets, gRPC)
  // --------------------------------------------------------------------------
  function renderNetworkView() {
    const n = state.data.activeZone.network;
    return `
      <div class="breadcrumbs">
        <span>Networking</span>
        <span class="divider">/</span>
        <span class="current">Network Protocols</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Edge Network Protocols</h1>
          <p class="page-subtitle">Manage modern transport protocols including HTTP/3 (QUIC), gRPC, and WebSockets.</p>
        </div>
      </div>

      <div class="stats-overview-grid">
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">HTTP/3 (with QUIC)</span></div>
          <div class="metric-value" style="color: #059669;">Enabled</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">0-RTT Connection Resumption</span></div>
          <div class="metric-value" style="color: #059669;">Enabled</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">IPv6 Compatibility</span></div>
          <div class="metric-value" style="color: #059669;">Full Global IPv6</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">WebSockets</span></div>
          <div class="metric-value" style="color: #059669;">Unlimited</div>
        </div>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 13: Workers & Pages List
  // --------------------------------------------------------------------------
  function renderWorkersView() {
    const workers = state.data.workers;
    return `
      <div class="breadcrumbs">
        <span>Compute</span>
        <span class="divider">/</span>
        <span class="current">Workers & Pages</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Workers & Pages</h1>
          <p class="page-subtitle">Deploy serverless code instantly to Cloudflare's global network with zero cold start latency.</p>
        </div>
        <button class="btn-cf btn-cf-primary" onclick="openModal('create_app')">+ Create Application</button>
      </div>

      <div class="item-cards-grid">
        ${workers.map(w => `
          <div class="kumo-item-card" onclick="switchTab('worker-detail')">
            <div class="kumo-card-primary">
              <div class="item-title-row">
                <div class="item-name">⚡ ${w.name}</div>
                <span class="cf-pill active"><span class="status-dot"></span><span>${w.status}</span></span>
              </div>
              <div class="item-meta-row">
                <span style="font-family: var(--font-mono); font-size: 0.8rem; color: var(--cf-brand-orange);">${w.custom_domain}</span>
                <span>•</span>
                <span>${w.repo}</span>
              </div>
            </div>
            <div class="kumo-card-secondary">
              <div class="secondary-stat-col">
                <span class="sec-label">Invocations</span>
                <span class="sec-val">${w.invocations_24h}</span>
              </div>
              <div class="secondary-stat-col">
                <span class="sec-label">CPU P90</span>
                <span class="sec-val">${w.median_cpu_time}</span>
              </div>
            </div>
          </div>
        `).join('')}
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 14: Account Analytics
  // --------------------------------------------------------------------------
  function renderAnalyticsView() {
    const a = state.data.analytics;
    return `
      <div class="breadcrumbs">
        <span>Observability</span>
        <span class="divider">/</span>
        <span class="current">Account Analytics</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Account Telemetry & Edge Analytics</h1>
          <p class="page-subtitle">Real-time telemetry and resource quotas across 330+ Anycast edge data centers.</p>
        </div>
      </div>

      <div class="stats-overview-grid">
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Total Requests</span></div>
          <div class="metric-value">${a.totalRequests}</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Cache Hit Ratio</span></div>
          <div class="metric-value">${a.cacheHitRate}</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Median CPU Time</span></div>
          <div class="metric-value">${a.cpuTimeP90}</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Build Minutes</span></div>
          <div class="metric-value">${a.buildMinutes}</div>
        </div>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 15: Zero Trust & Security Center
  // --------------------------------------------------------------------------
  function renderZeroTrustView() {
    return `
      <div class="breadcrumbs">
        <span>Security</span>
        <span class="divider">/</span>
        <span class="current">Zero Trust Center</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Cloudflare Zero Trust & Tunnels</h1>
          <p class="page-subtitle">Identity-aware access proxy, sovereign Cloudflare Tunnels, and private subnet routing.</p>
        </div>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 16: Team Members & RBAC
  // --------------------------------------------------------------------------
  function renderTeamView() {
    const members = state.data.teamMembers;
    return `
      <div class="breadcrumbs">
        <span>Manage Account</span>
        <span class="divider">/</span>
        <span class="current">Members & Permissions</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Account Members & Roles (${members.length})</h1>
          <p class="page-subtitle">Manage member access, role-based access control (RBAC), and 2FA enforcement.</p>
        </div>
        <button class="btn-cf btn-cf-primary" onclick="openModal('invite_member')">+ Invite Members</button>
      </div>

      <div class="kumo-table-container">
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Member Name</th>
              <th>Email</th>
              <th>Role</th>
              <th>Groups</th>
              <th>2FA Status</th>
              <th>Last Active</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            ${members.map(m => `
              <tr>
                <td>
                  <div style="display: flex; align-items: center; gap: 0.65rem;">
                    <div class="user-avatar-initials" style="background: ${m.color};">${m.initials}</div>
                    <strong>${m.name}</strong>
                  </div>
                </td>
                <td>${m.email}</td>
                <td><span class="badge-tag-beta">${m.role}</span></td>
                <td>${m.groups}</td>
                <td>
                  <span class="cf-pill ${m.twoFactor ? 'active' : 'stopped'}">
                    <span class="status-dot"></span>
                    <span>${m.twoFactor ? 'Enabled' : 'Disabled'}</span>
                  </span>
                </td>
                <td>${m.lastActive}</td>
                <td>
                  <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="alert('Manage permissions for ${m.email}')">Manage</button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 17: Billing & Subscriptions
  // --------------------------------------------------------------------------
  function renderBillingView() {
    const b = state.data.billing;
    return `
      <div class="breadcrumbs">
        <span>Manage Account</span>
        <span class="divider">/</span>
        <span class="current">Billing & Subscriptions</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Billing & Subscriptions</h1>
          <p class="page-subtitle">Current Plan: <strong>${b.currentPlan}</strong> (${b.monthlyPrice}). Next invoice on ${b.nextPayment}.</p>
        </div>
        <button class="btn-cf btn-cf-primary" onclick="openModal('upgrade_plan')">Change Plan</button>
      </div>

      <!-- Plan Cards Grid -->
      <div class="plan-cards-grid">
        <div class="plan-card">
          <div>
            <h3>Free Plan</h3>
            <div class="plan-price">$0 <span>/ month</span></div>
            <ul class="plan-features-list">
              <li>✓ Fast Anycast DNS</li>
              <li>✓ Universal SSL Certificate</li>
              <li>✓ 100k Workers Invocations/day</li>
              <li>✓ Community Support</li>
            </ul>
          </div>
          <button class="btn-cf btn-cf-secondary" disabled>Default</button>
        </div>

        <div class="plan-card current">
          <span class="plan-badge-current">Active Plan</span>
          <div>
            <h3>Pro Plan</h3>
            <div class="plan-price">$20 <span>/ month</span></div>
            <ul class="plan-features-list">
              <li>✓ Everything in Free</li>
              <li>✓ Web Application Firewall (WAF)</li>
              <li>✓ Lossless Image Optimization</li>
              <li>✓ 20 Page Rules</li>
            </ul>
          </div>
          <button class="btn-cf btn-cf-primary" disabled>Current Subscription</button>
        </div>

        <div class="plan-card">
          <div>
            <h3>Business Plan</h3>
            <div class="plan-price">$200 <span>/ month</span></div>
            <ul class="plan-features-list">
              <li>✓ Everything in Pro</li>
              <li>✓ 100% Uptime SLA</li>
              <li>✓ Custom Nameservers</li>
              <li>✓ PCI DSS Compliance</li>
            </ul>
          </div>
          <button class="btn-cf btn-cf-secondary" onclick="openModal('upgrade_plan')">Upgrade to Business</button>
        </div>
      </div>

      <!-- Invoices History Table -->
      <div class="kumo-table-container">
        <div class="table-header-title">
          <span>Invoices & Statements (${b.invoices.length})</span>
          <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">Payment method: ${b.paymentMethods[0].brand} ending in ${b.paymentMethods[0].last4}</span>
        </div>
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Invoice Number</th>
              <th>Date</th>
              <th>Description</th>
              <th>Amount</th>
              <th>Status</th>
              <th>Receipt</th>
            </tr>
          </thead>
          <tbody>
            ${b.invoices.map(inv => `
              <tr>
                <td><code>${inv.id}</code></td>
                <td>${inv.date}</td>
                <td>${inv.description}</td>
                <td><strong>${inv.amount}</strong></td>
                <td><span class="cf-pill active"><span class="status-dot"></span><span>${inv.status}</span></span></td>
                <td><button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="alert('Downloading invoice PDF...')">Download PDF ↗</button></td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 18: API Tokens & Keys
  // --------------------------------------------------------------------------
  function renderTokensView() {
    const tokens = state.data.apiTokens;
    return `
      <div class="breadcrumbs">
        <span>Manage Account</span>
        <span class="divider">/</span>
        <span class="current">API Tokens</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">API Tokens & Access Credentials</h1>
          <p class="page-subtitle">Generate scoped authorization tokens with granular read/write permissions for CI/CD and Terraform.</p>
        </div>
        <button class="btn-cf btn-cf-primary" onclick="openModal('create_token')">+ Create Token</button>
      </div>

      <div class="kumo-table-container">
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Token Name</th>
              <th>Permissions</th>
              <th>Status</th>
              <th>Created</th>
              <th>Last Used</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            ${tokens.map(t => `
              <tr>
                <td><strong>${t.name}</strong></td>
                <td>${t.permissions.map(p => `<span class="badge-tag-beta" style="margin-right: 0.25rem;">${p}</span>`).join('')}</td>
                <td>
                  <span class="cf-pill ${t.status === 'Active' ? 'active' : 'stopped'}">
                    <span class="status-dot"></span>
                    <span>${t.status}</span>
                  </span>
                </td>
                <td>${t.created}</td>
                <td>${t.lastUsed}</td>
                <td>
                  <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="copyToClipboard('${t.full}', 'tok_${t.id}')">
                    ${state.copiedKey === 'tok_' + t.id ? '✓ Copied' : 'Copy'}
                  </button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 19: Audit Logs Stream
  // --------------------------------------------------------------------------
  function renderAuditLogsView() {
    const logs = state.data.auditLogs;
    return `
      <div class="breadcrumbs">
        <span>Manage Account</span>
        <span class="divider">/</span>
        <span class="current">Audit Logs</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Audit Logs & Activity Stream</h1>
          <p class="page-subtitle">Immutable security trail tracking all account configurations, DNS mutations, and worker deployments.</p>
        </div>
        <button class="btn-cf btn-cf-secondary" onclick="alert('Audit logs exported as JSON.')">Export JSON</button>
      </div>

      <div class="kumo-table-container">
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Timestamp (UTC)</th>
              <th>Action / Event</th>
              <th>Subject / Resource</th>
              <th>Actor</th>
              <th>Zone Scope</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            ${logs.map(l => `
              <tr>
                <td><code style="font-size: 0.78rem;">${l.timestamp}</code></td>
                <td><strong>${l.event}</strong></td>
                <td><code>${l.subject}</code></td>
                <td>${l.actor}</td>
                <td>${l.zone}</td>
                <td><span class="cf-pill active"><span class="status-dot"></span><span>${l.status}</span></span></td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 20: Profile & Preferences Settings
  // --------------------------------------------------------------------------
  function renderSettingsView() {
    const p = state.data.settings.profile;
    const s = state.data.settings.security;
    return `
      <div class="breadcrumbs">
        <span>Manage Account</span>
        <span class="divider">/</span>
        <span class="current">Profile & Preferences</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">User Profile & Account Preferences</h1>
          <p class="page-subtitle">Manage personal account credentials, sessions, and multi-factor authentication.</p>
        </div>
      </div>

      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(360px, 1fr)); gap: 1.5rem;">
        <div class="ui-kit-section">
          <h3 class="ui-kit-title">User Profile</h3>
          <div class="cf-form-group" style="margin-top: 1rem;">
            <label>Full Name</label>
            <input type="text" class="cf-input" value="${p.firstName} ${p.lastName}" />
          </div>
          <div class="cf-form-group">
            <label>Primary Email</label>
            <input type="email" class="cf-input" value="${p.email}" readonly />
          </div>
          <div class="cf-form-group">
            <label>Preferred Timezone</label>
            <input type="text" class="cf-input" value="${p.timezone}" />
          </div>
          <button class="btn-cf btn-cf-primary" onclick="alert('Profile changes saved.')">Save Changes</button>
        </div>

        <div class="ui-kit-section">
          <h3 class="ui-kit-title">Active Security Sessions</h3>
          <p class="ui-kit-desc">Review connected devices and revoke unauthorized browser sessions.</p>
          <div style="display: flex; flex-direction: column; gap: 0.75rem;">
            ${s.activeSessions.map(sess => `
              <div style="display: flex; justify-content: space-between; align-items: center; padding: 0.75rem; border: 1px solid var(--color-kumo-line); border-radius: var(--radius-md);">
                <div>
                  <div style="font-weight: 600; font-size: 0.88rem;">${sess.device}</div>
                  <div style="font-size: 0.78rem; color: var(--text-color-kumo-subtle);">${sess.browser} • ${sess.location}</div>
                </div>
                ${sess.current ? '<span class="cf-pill active"><span class="status-dot"></span><span>Current Session</span></span>' : '<button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="alert(\'Session revoked.\')">Revoke</button>'}
              </div>
            `).join('')}
          </div>
        </div>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 21: UI Components Kit Showcase (Envato Commercial Requirement)
  // --------------------------------------------------------------------------
  function renderComponentsView() {
    return `
      <div class="breadcrumbs">
        <span>Template Kit</span>
        <span class="divider">/</span>
        <span class="current">UI Components Showcase</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Kumo Design System UI Kit</h1>
          <p class="page-subtitle">Commercial component library including buttons, status pills, inputs, cards, tables, and alerts.</p>
        </div>
      </div>

      <!-- Buttons Showcase -->
      <div class="ui-kit-section">
        <h3 class="ui-kit-title">Button Variants & Sizes</h3>
        <p class="ui-kit-desc">Primary brand buttons, secondary neutrals, small action buttons, and loading states.</p>
        <div class="ui-kit-row">
          <button class="btn-cf btn-cf-primary">Primary Button</button>
          <button class="btn-cf btn-cf-secondary">Secondary Button</button>
          <button class="btn-cf btn-cf-primary btn-cf-sm">Small Primary</button>
          <button class="btn-cf btn-cf-secondary btn-cf-sm">Small Secondary</button>
          <button class="btn-cf btn-cf-primary" disabled>Disabled State</button>
        </div>
      </div>

      <!-- Badges & Status Pills -->
      <div class="ui-kit-section">
        <h3 class="ui-kit-title">Badges & Status Pills</h3>
        <div class="ui-kit-row">
          <span class="cf-pill active"><span class="status-dot"></span><span>Active</span></span>
          <span class="cf-pill stopped"><span class="status-dot" style="background: #94a3b8;"></span><span>Inactive</span></span>
          <span class="badge-tag-beta">Beta Tag</span>
          <span class="dns-type-badge">A Record</span>
          <span class="dns-type-badge">CNAME</span>
          <span class="dns-status-cloud proxied"><span>☁️</span><span>Proxied</span></span>
          <span class="dns-status-cloud dns-only"><span>☁️</span><span>DNS Only</span></span>
        </div>
      </div>

      <!-- Alert Banners -->
      <div class="ui-kit-section">
        <h3 class="ui-kit-title">Notification & Alert Banners</h3>
        <div class="ui-alert-banner info">
          <div>ℹ️</div>
          <div><strong>Informational Notice:</strong> Anycast edge routing is operating with zero packet loss worldwide.</div>
        </div>
        <div class="ui-alert-banner success">
          <div>✓</div>
          <div><strong>Success Event:</strong> Worker application deployed successfully in 335ms.</div>
        </div>
        <div class="ui-alert-banner warning">
          <div>⚠️</div>
          <div><strong>Warning Notice:</strong> DNS record propagation may take up to 60 seconds.</div>
        </div>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 22: Workers AI & AI Gateway
  // --------------------------------------------------------------------------
  function renderAiView() {
    const ai = state.data.workersAi;
    return `
      <div class="breadcrumbs">
        <span>Build</span>
        <span class="divider">/</span>
        <span>AI</span>
        <span class="divider">/</span>
        <span class="current">Workers AI & Gateway</span>
      </div>

      <div class="page-title-row">
        <div>
          <div style="display: flex; align-items: center; gap: 0.75rem;">
            <h1 class="page-title">🤖 Workers AI & AI Gateway</h1>
            <span class="cf-pill active"><span class="status-dot"></span><span>GPU Edge Active</span></span>
            <span class="badge-tag-beta">Tensor Core</span>
          </div>
          <p class="page-subtitle">Serverless GPU inference running Meta Llama 3.3, Mistral, Whisper, and SDXL across Cloudflare's Anycast network.</p>
        </div>
        <button class="btn-cf btn-cf-primary" onclick="runAiInference()">▶ Run Prompt Inference</button>
      </div>

      <!-- AI Gateway Telemetry Cards -->
      <div class="stats-overview-grid">
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">24h AI Inferences</span></div>
          <div class="metric-value">${ai.gateway.totalRequests}</div>
          <div class="metric-footer">Tokens: ${ai.gateway.tokensGenerated}</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Gateway Cache Hit</span></div>
          <div class="metric-value" style="color: #059669;">${ai.gateway.cacheHitRate}</div>
          <div class="metric-footer">Saved: ${ai.gateway.costSaved}</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Median Latency</span></div>
          <div class="metric-value">${ai.gateway.avgLatency}</div>
          <div class="metric-footer">Streaming first-token</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Active GPU Cluster</span></div>
          <div class="metric-value" style="font-size: 1.25rem;">NVIDIA H100 / L40S</div>
          <div class="metric-footer">Serverless edge execution</div>
        </div>
      </div>

      <!-- Interactive AI Prompt Playground -->
      <div class="ui-kit-section" style="margin-top: 1.5rem;">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
          <div>
            <h3 class="ui-kit-title" style="margin: 0;">Interactive AI Prompt Playground</h3>
            <p class="ui-kit-desc" style="margin: 0.25rem 0 0 0;">Test real-time model inference and response streaming directly from Cloudflare V8 isolates.</p>
          </div>
          <div style="display: flex; gap: 0.5rem; align-items: center;">
            <select id="cf-ai-model-select" class="cf-select" style="width: 240px; font-size: 0.82rem;">
              <option value="llama-3.3-70b">Meta Llama 3.3 70B Instruct</option>
              <option value="mistral-7b-v0.3">Mistral 7B Instruct v0.3</option>
              <option value="whisper-large-v3">OpenAI Whisper Large v3</option>
              <option value="bge-large-en">BAAI BGE-Large-EN v1.5</option>
              <option value="sdxl-turbo">Stability AI SDXL Turbo</option>
            </select>
            <button class="btn-cf btn-cf-primary btn-cf-sm" onclick="runAiInference()">▶ Run (Ctrl+Enter)</button>
          </div>
        </div>

        <div class="cf-form-group">
          <label>Prompt Input</label>
          <textarea id="cf-ai-prompt-input" class="cf-input" style="height: 70px; resize: vertical; font-family: var(--font-mono); font-size: 0.84rem;">${ai.samplePrompt}</textarea>
        </div>

        <div>
          <label style="font-size: 0.82rem; font-weight: 600; color: var(--text-color-kumo-strong); margin-bottom: 0.35rem; display: block;">Streaming Inference Output</label>
          <div id="cf-ai-response-box" style="background: #0f172a; border: 1px solid #1e293b; border-radius: var(--radius-md); padding: 1.15rem; color: #38bdf8; font-family: var(--font-mono); font-size: 0.84rem; line-height: 1.6; min-height: 110px;">
            ${ai.sampleResponse}
          </div>
          <div style="display: flex; justify-content: space-between; margin-top: 0.5rem; font-size: 0.76rem; color: var(--text-color-kumo-subtle);">
            <span>Generation speed: <strong>142 tokens/sec</strong></span>
            <span>Edge latency: <strong>114 ms</strong> • Tokens used: <strong>89 tokens</strong></span>
          </div>
        </div>
      </div>

      <!-- Models Catalog Cards -->
      <div class="kumo-table-container">
        <div class="table-header-title">
          <span>Supported Edge Models Catalog (${ai.models.length})</span>
          <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">Available via <code>@cloudflare/ai</code> SDK</span>
        </div>
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Model Name</th>
              <th>Author</th>
              <th>Task / Category</th>
              <th>Inference Speed</th>
              <th>Context Window</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            ${ai.models.map(m => `
              <tr>
                <td><strong style="color: var(--text-color-kumo-strong);">🧠 ${m.name}</strong></td>
                <td>${m.author}</td>
                <td><span class="badge-tag-beta">${m.category}</span></td>
                <td>${m.speed}</td>
                <td><code>${m.context}</code></td>
                <td><span class="cf-pill active"><span class="status-dot"></span><span>${m.status}</span></span></td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 23: Workers KV Store
  // --------------------------------------------------------------------------
  function renderKvView() {
    const kv = state.data.workersKv;
    return `
      <div class="breadcrumbs">
        <span>Storage & Databases</span>
        <span class="divider">/</span>
        <span class="current">Workers KV Namespaces</span>
      </div>

      <div class="page-title-row">
        <div>
          <h1 class="page-title">Workers KV Key-Value Store</h1>
          <p class="page-subtitle">Low-latency distributed global key-value store with read replication to 330+ Cloudflare edge locations.</p>
        </div>
        <button class="btn-cf btn-cf-primary" onclick="alert('Create KV Namespace wizard opened.')">+ Create Namespace</button>
      </div>

      <!-- Namespaces Table -->
      <div class="kumo-table-container" style="margin-bottom: 2rem;">
        <div class="table-header-title">
          <span>Active KV Namespaces (${kv.namespaces.length})</span>
        </div>
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Namespace Name</th>
              <th>Namespace ID</th>
              <th>Total Keys</th>
              <th>Storage Size</th>
              <th>24h Reads / Writes</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            ${kv.namespaces.map(ns => `
              <tr>
                <td><strong style="color: var(--text-color-kumo-strong);">🔑 ${ns.name}</strong></td>
                <td><code>${ns.id}</code></td>
                <td>${ns.keys_count.toLocaleString()} keys</td>
                <td>${ns.size}</td>
                <td>${ns.read_ops_24h} / ${ns.write_ops_24h}</td>
                <td>
                  <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="alert('Viewing keys in ${ns.name}')">Browse Keys</button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>

      <!-- Sample Key-Value Browser -->
      <div class="kumo-table-container">
        <div class="table-header-title">
          <span>Keys in EDGE_CONFIG_STORE</span>
          <div style="display: flex; gap: 0.5rem;">
            <input type="text" class="cf-input" placeholder="Search keys by prefix..." style="width: 220px; font-size: 0.8rem; padding: 0.3rem 0.6rem;" />
            <button class="btn-cf btn-cf-primary btn-cf-sm" onclick="alert('Add Key modal.')">+ Add Key</button>
          </div>
        </div>
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Key Name</th>
              <th>Value (JSON / String)</th>
              <th>Expiration TTL</th>
              <th>Metadata</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            ${kv.sampleKeys.map(k => `
              <tr>
                <td><strong>${k.key}</strong></td>
                <td><code>${k.value}</code></td>
                <td>${k.ttl}</td>
                <td><span class="badge-tag-beta">${k.metadata}</span></td>
                <td>
                  <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="copyToClipboard('${k.value}', 'kv_${k.key}')">
                    ${state.copiedKey === 'kv_' + k.key ? '✓ Copied' : 'Copy'}
                  </button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // View 24: Turnstile CAPTCHA & Bot Protection
  // --------------------------------------------------------------------------
  function renderTurnstileView() {
    const t = state.data.turnstile;
    return `
      <div class="breadcrumbs">
        <span>Protect & Connect</span>
        <span class="divider">/</span>
        <span class="current">Turnstile CAPTCHA</span>
      </div>

      <div class="page-title-row">
        <div>
          <div style="display: flex; align-items: center; gap: 0.75rem;">
            <h1 class="page-title">🛡️ Cloudflare Turnstile</h1>
            <span class="cf-pill active"><span class="status-dot"></span><span>Active</span></span>
            <span class="badge-tag-beta">Privacy First</span>
          </div>
          <p class="page-subtitle">Smart CAPTCHA alternative protecting websites from bad bots with seamless zero-friction verification.</p>
        </div>
        <button class="btn-cf btn-cf-primary" onclick="alert('Add Turnstile Widget wizard.')">+ Add Widget</button>
      </div>

      <!-- Stats Strip -->
      <div class="stats-overview-grid">
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Total Solves (24h)</span></div>
          <div class="metric-value">145,430</div>
          <div class="metric-footer">Across 3 active domains</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Visitor Pass Rate</span></div>
          <div class="metric-value" style="color: #059669;">99.8%</div>
          <div class="metric-footer">Zero friction human pass</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Average Solve Time</span></div>
          <div class="metric-value">0.4 sec</div>
          <div class="metric-footer">Non-interactive challenge</div>
        </div>
        <div class="metric-card">
          <div class="metric-header"><span class="metric-label">Bots Mitigated</span></div>
          <div class="metric-value" style="color: #dc2626;">290</div>
          <div class="metric-footer">Automated scrapers blocked</div>
        </div>
      </div>

      <!-- Interactive Turnstile Widget Simulator -->
      <div class="ui-kit-section" style="margin-top: 1.5rem;">
        <h3 class="ui-kit-title">Interactive Turnstile Widget Simulator</h3>
        <p class="ui-kit-desc">Experience the exact visitor verification challenge embedded in authentication forms.</p>

        <div style="background: var(--color-kumo-canvas); border: 1px solid var(--color-kumo-line); border-radius: var(--radius-lg); padding: 1.5rem; max-width: 380px;">
          <div id="cf-turnstile-box" style="display: flex; align-items: center; justify-content: space-between; background: var(--color-kumo-base); border: 1px solid var(--color-kumo-line); border-radius: var(--radius-md); padding: 0.85rem 1.15rem; cursor: pointer;" onclick="solveTurnstile()">
            <div style="display: flex; align-items: center; gap: 0.85rem;">
              <div id="cf-turnstile-checkbox" style="width: 24px; height: 24px; border: 2px solid var(--color-kumo-line); border-radius: 4px; display: flex; align-items: center; justify-content: center; background: var(--color-kumo-base); transition: var(--kumo-transition);"></div>
              <span id="cf-turnstile-label" style="font-size: 0.88rem; font-weight: 500; color: var(--text-color-kumo-strong);">Verify you are human</span>
            </div>
            <div style="display: flex; flex-direction: column; align-items: flex-end;">
              <img src="assets/cloudflare-logo.svg" alt="Cloudflare" style="height: 18px;" />
              <span style="font-size: 0.65rem; color: var(--text-color-kumo-subtle); margin-top: 2px;">Cloudflare Turnstile</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Active Widgets Table -->
      <div class="kumo-table-container">
        <div class="table-header-title">
          <span>Configured Turnstile Widgets (${t.widgets.length})</span>
        </div>
        <table class="kumo-table">
          <thead>
            <tr>
              <th>Widget Name</th>
              <th>Site Key</th>
              <th>Domain Scope</th>
              <th>Challenge Mode</th>
              <th>24h Solves</th>
              <th>Pass Rate</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            ${t.widgets.map(w => `
              <tr>
                <td><strong>${w.name}</strong></td>
                <td><code>${w.id}</code></td>
                <td>${w.domain}</td>
                <td><span class="badge-tag-beta">${w.mode}</span></td>
                <td>${w.solves_24h}</td>
                <td><strong style="color: #059669;">${w.passRate}</strong></td>
                <td><span class="cf-pill active"><span class="status-dot"></span><span>${w.status}</span></span></td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `;
  }

  // --------------------------------------------------------------------------
  // Modals Manager
  // --------------------------------------------------------------------------
  function renderModals() {
    if (!modalContainer) return;
    if (!state.activeModal) {
      modalContainer.innerHTML = '';
      return;
    }

    if (state.activeModal === 'add_dns_record') {
      modalContainer.innerHTML = `
        <div class="cf-modal-backdrop" onclick="if(event.target === this) closeModal()">
          <div class="cf-modal-dialog">
            <div class="cf-modal-head">
              <h3>Add DNS Record to ${state.data.activeZone.name}</h3>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="closeModal()">✕</button>
            </div>
            <div class="cf-modal-body">
              <div class="cf-form-group">
                <label>Record Type</label>
                <select class="cf-select">
                  <option>A (IPv4 Address)</option>
                  <option>AAAA (IPv6 Address)</option>
                  <option>CNAME (Domain Alias)</option>
                  <option>TXT (Text Record)</option>
                  <option>MX (Mail Exchanger)</option>
                </select>
              </div>
              <div class="cf-form-group">
                <label>Name (Use @ for root)</label>
                <input type="text" class="cf-input" placeholder="subdomain or @" />
              </div>
              <div class="cf-form-group">
                <label>IPv4 Address / Target</label>
                <input type="text" class="cf-input" placeholder="198.51.100.24" />
              </div>
              <div class="cf-form-group">
                <label>Proxy Status</label>
                <div style="display: flex; align-items: center; gap: 0.75rem;">
                  <span class="dns-status-cloud proxied"><span>☁️</span><span>Proxied (Orange Cloud)</span></span>
                  <span style="font-size: 0.8rem; color: var(--text-color-kumo-subtle);">Protects origin IP and accelerates traffic</span>
                </div>
              </div>
            </div>
            <div class="cf-modal-foot">
              <button class="btn-cf btn-cf-secondary" onclick="closeModal()">Cancel</button>
              <button class="btn-cf btn-cf-primary" onclick="alert('DNS Record added!'); closeModal();">Save Record</button>
            </div>
          </div>
        </div>
      `;
    } else if (state.activeModal === 'purge_cache') {
      modalContainer.innerHTML = `
        <div class="cf-modal-backdrop" onclick="if(event.target === this) closeModal()">
          <div class="cf-modal-dialog">
            <div class="cf-modal-head">
              <h3>Purge Edge Cache</h3>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="closeModal()">✕</button>
            </div>
            <div class="cf-modal-body">
              <p style="font-size: 0.88rem; color: var(--text-color-kumo-default); margin-bottom: 1rem;">
                Purging clears cached assets across all 330+ Cloudflare data centers immediately.
              </p>
              <div style="display: flex; gap: 0.75rem;">
                <button class="btn-cf btn-cf-primary" style="flex: 1;" onclick="alert('Cache purged across all global edge servers.'); closeModal();">Purge Everything</button>
                <button class="btn-cf btn-cf-secondary" style="flex: 1;" onclick="alert('Custom purge by URL.'); closeModal();">Custom Purge by URL</button>
              </div>
            </div>
          </div>
        </div>
      `;
    } else if (state.activeModal === 'invite_member') {
      modalContainer.innerHTML = `
        <div class="cf-modal-backdrop" onclick="if(event.target === this) closeModal()">
          <div class="cf-modal-dialog">
            <div class="cf-modal-head">
              <h3>Invite Member to Account</h3>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="closeModal()">✕</button>
            </div>
            <div class="cf-modal-body">
              <div class="cf-form-group">
                <label>Email Address</label>
                <input type="email" class="cf-input" placeholder="colleague@company.com" />
              </div>
              <div class="cf-form-group">
                <label>Assigned Role</label>
                <select class="cf-select">
                  <option>Administrator</option>
                  <option>Developer</option>
                  <option>DNS Manager</option>
                  <option>Billing Only</option>
                  <option>Read Only</option>
                </select>
              </div>
            </div>
            <div class="cf-modal-foot">
              <button class="btn-cf btn-cf-secondary" onclick="closeModal()">Cancel</button>
              <button class="btn-cf btn-cf-primary" onclick="alert('Invitation sent successfully.'); closeModal();">Send Invite</button>
            </div>
          </div>
        </div>
      `;
    } else if (state.activeModal === 'create_token') {
      modalContainer.innerHTML = `
        <div class="cf-modal-backdrop" onclick="if(event.target === this) closeModal()">
          <div class="cf-modal-dialog">
            <div class="cf-modal-head">
              <h3>Create Scoped API Token</h3>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" onclick="closeModal()">✕</button>
            </div>
            <div class="cf-modal-body">
              <div class="cf-form-group">
                <label>Token Name</label>
                <input type="text" class="cf-input" placeholder="e.g. GitHub Actions CI/CD" />
              </div>
              <div class="cf-form-group">
                <label>Permissions Template</label>
                <select class="cf-select">
                  <option>Edit Cloudflare Workers</option>
                  <option>Edit DNS Records (Zone DNS)</option>
                  <option>Read All Resources</option>
                  <option>Custom Scoped Permissions</option>
                </select>
              </div>
            </div>
            <div class="cf-modal-foot">
              <button class="btn-cf btn-cf-secondary" onclick="closeModal()">Cancel</button>
              <button class="btn-cf btn-cf-primary" onclick="alert('API Token created: cf_pat_981a... (Copied to clipboard)'); closeModal();">Generate Token</button>
            </div>
          </div>
        </div>
      `;
    }
  }

  // SQL Runner Helper
  window.executeSqlQuery = function () {
    const el = document.getElementById('cf-sql-input');
    if (el) state.sqlQuery = el.value;
    alert('Query executed: ' + state.sqlQuery + '\\n5 rows returned in 0.8ms');
    render();
  };

  // AI Inference Runner Helper with Streaming Typing Effect
  window.runAiInference = function () {
    const box = document.getElementById('cf-ai-response-box');
    const input = document.getElementById('cf-ai-prompt-input');
    const promptText = input ? input.value : state.data.workersAi.samplePrompt;
    if (!box) return;

    box.innerHTML = '<span style="color: #94a3b8;">Streaming from Cloudflare Anycast GPU (IAD Edge)...</span>';
    const fullText = state.data.workersAi.sampleResponse;
    let i = 0;
    box.innerHTML = '';
    const timer = setInterval(() => {
      if (i < fullText.length) {
        box.innerHTML += fullText.charAt(i);
        i += 3;
      } else {
        clearInterval(timer);
        box.innerHTML = fullText;
      }
    }, 15);
  };

  // Turnstile Interactive Widget Simulator Helper
  window.solveTurnstile = function () {
    const checkbox = document.getElementById('cf-turnstile-checkbox');
    const label = document.getElementById('cf-turnstile-label');
    if (!checkbox || !label) return;

    checkbox.innerHTML = '<div style="width: 14px; height: 14px; border: 2px solid var(--cf-brand-orange); border-top-color: transparent; border-radius: 50%; animation: spinSlow 0.6s linear infinite;"></div>';
    label.textContent = 'Verifying...';

    setTimeout(() => {
      checkbox.style.background = '#059669';
      checkbox.style.borderColor = '#059669';
      checkbox.innerHTML = '<span style="color: #ffffff; font-size: 14px; font-weight: 700;">✓</span>';
      label.textContent = 'Verification successful (Human visitor verified)';
      label.style.color = '#059669';
    }, 450);
  };

  window.queryTable = function (tableName) {
    state.sqlQuery = 'SELECT * FROM ' + tableName + ' LIMIT 10;';
    state.subTab = 'overview';
    render();
  };

  window.setSslMode = function (mode) {
    state.data.activeZone.ssl.mode = mode;
    render();
  };

  // --------------------------------------------------------------------------
  // Quick Search ⌘K Spotlight Palette
  // --------------------------------------------------------------------------
  window.openQuickSearch = function () {
    state.searchOpen = true;
    if (quickSearchBackdrop) {
      quickSearchBackdrop.style.display = 'flex';
      setTimeout(() => {
        if (quickSearchInput) {
          quickSearchInput.value = '';
          quickSearchInput.focus();
          handleSearchInput();
        }
      }, 50);
    }
  };

  window.closeQuickSearch = function () {
    state.searchOpen = false;
    if (quickSearchBackdrop) {
      quickSearchBackdrop.style.display = 'none';
    }
  };

  window.handleSearchInput = function () {
    const q = (quickSearchInput.value || '').toLowerCase().trim();
    if (!quickSearchResults) return;

    let items = [
      { type: 'Page', title: 'Account Home', tab: 'home', badge: 'Overview' },
      { type: 'Page', title: 'Zone Overview (northlab.dev)', tab: 'zone-overview', badge: 'Zone' },
      { type: 'Page', title: 'DNS Records & Anycast', tab: 'dns', badge: 'DNS' },
      { type: 'Page', title: 'Workers & Pages', tab: 'workers', badge: `${state.data.workers.length} workers` },
      { type: 'Page', title: 'Worker: code-collection', tab: 'worker-detail', badge: 'Production' },
      { type: 'Page', title: 'Workers AI & AI Gateway', tab: 'ai', badge: 'Llama 3.3' },
      { type: 'Page', title: 'D1 SQL Databases', tab: 'd1', badge: 'SQLite' },
      { type: 'Page', title: 'D1: code_collection Console', tab: 'd1-detail', badge: 'SQL' },
      { type: 'Page', title: 'R2 Object Storage', tab: 'r2', badge: 'Zero Egress' },
      { type: 'Page', title: 'R2: code-collection Bucket', tab: 'r2-detail', badge: 'Files' },
      { type: 'Page', title: 'Workers KV Namespaces', tab: 'kv', badge: 'Key-Value' },
      { type: 'Page', title: 'SSL/TLS Encryption', tab: 'ssl', badge: 'Strict' },
      { type: 'Page', title: 'WAF & Security Rules', tab: 'waf', badge: 'Firewall' },
      { type: 'Page', title: 'Turnstile CAPTCHA', tab: 'turnstile', badge: 'Privacy' },
      { type: 'Page', title: 'Caching & Optimization', tab: 'caching', badge: 'CDN' },
      { type: 'Page', title: 'Edge Network Protocols', tab: 'network', badge: 'HTTP/3' },
      { type: 'Page', title: 'Members & RBAC', tab: 'team', badge: '5 Members' },
      { type: 'Page', title: 'Billing & Subscriptions', tab: 'billing', badge: 'Invoices' },
      { type: 'Page', title: 'API Tokens & Keys', tab: 'tokens', badge: 'Credentials' },
      { type: 'Page', title: 'Audit Logs Stream', tab: 'audit-logs', badge: 'Security' },
      { type: 'Page', title: 'Profile & Preferences', tab: 'settings', badge: 'Account' },
      { type: 'Page', title: 'UI Components Kit', tab: 'components', badge: 'Showcase' },
    ];

    const filtered = q ? items.filter(i => i.title.toLowerCase().includes(q) || i.type.toLowerCase().includes(q)) : items;

    quickSearchResults.innerHTML = `
      <div class="cf-command-section-title">Navigation & Quick Switcher (${filtered.length} results)</div>
      ${filtered.map(item => `
        <div class="cf-command-item" onclick="switchTab('${item.tab}'); closeQuickSearch();">
          <div class="item-left">
            <span>📄</span>
            <strong>${item.title}</strong>
          </div>
          <span class="cf-pill stopped">${item.badge}</span>
        </div>
      `).join('')}
    `;
  };

  // Keyboard Shortcuts Listener
  window.addEventListener('keydown', function (e) {
    if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
      e.preventDefault();
      if (state.searchOpen) {
        closeQuickSearch();
      } else {
        openQuickSearch();
      }
    } else if (e.key === 'Escape') {
      if (state.searchOpen) closeQuickSearch();
      if (state.activeModal) closeModal();
      if (state.customizerOpen) toggleCustomizer(false);
    }
  });

  // URL Hash Listener
  window.addEventListener('hashchange', () => {
    const hash = window.location.hash.replace('#', '');
    if (hash && hash !== state.currentTab) {
      state.currentTab = hash;
      render();
    }
  });

  // Init
  document.addEventListener('DOMContentLoaded', () => {
    const hash = window.location.hash.replace('#', '');
    if (hash) state.currentTab = hash;
    applyTheme();
    setAccentColor(state.accent);
    if (state.sidebarMini) document.body.classList.add('sidebar-mini');
    render();
  });

})();
