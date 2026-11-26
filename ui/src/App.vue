<template>
  <div class="orbit-shell">
    <!-- CLOUDFLARE STYLE SIDEBAR -->
    <aside class="cf-sidebar">
      <div class="sidebar-header">
        <a href="#" class="cf-logo" @click.prevent="currentTab = 'overview'">
          <!-- Cloudflare-style Orange Cloud SVG -->
          <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M37.5 19.5C36.8 13.5 31.8 9 25.5 9C20.2 9 15.8 12.2 14.1 16.8C13.4 16.6 12.7 16.5 12 16.5C7.6 16.5 4 20.1 4 24.5C4 28.9 7.6 32.5 12 32.5H37C40.9 32.5 44 29.4 44 25.5C44 21.8 41.1 18.8 37.5 19.5Z" fill="#F38020"/>
            <path d="M32 32.5H37C40.9 32.5 44 29.4 44 25.5C44 21.8 41.1 18.8 37.5 19.5C36.8 13.5 31.8 9 25.5 9C24.5 9 23.5 9.1 22.5 9.4C26.5 12.8 29.1 17.8 29.4 23.5C30.2 23.5 31.1 23.6 32 24C32.1 26.8 32 29.6 32 32.5Z" fill="#FAAD3F"/>
          </svg>
          <span>Orbit Sovereign</span>
        </a>
      </div>

      <div class="account-selector">
        <div style="display: flex; align-items: center; gap: 0.5rem; overflow: hidden;">
          <span style="width: 8px; height: 8px; border-radius: 50%; background: #f38020; flex-shrink: 0;"></span>
          <span style="font-weight: 500; text-overflow: ellipsis; overflow: hidden; white-space: nowrap;">Sovereign Cluster (Primary)</span>
        </div>
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m6 9 6 6 6-6"/></svg>
      </div>

      <div class="sidebar-nav">
        <!-- QUICK SEARCH BUTTON (CLOUDFLARE KUMO) -->
        <button class="sidebar-search-btn" @click="showQuickSearch = true">
          <svg viewBox="0 0 256 256" fill="currentColor">
            <path d="M229.66,218.34l-50.07-50.06a88.11,88.11,0,1,0-11.31,11.31l50.06,50.07a8,8,0,0,0,11.32-11.32ZM40,112a72,72,0,1,1,72,72A72.08,72.08,0,0,1,40,112Z"/>
          </svg>
          <span>Quick search...</span>
          <kbd>⌘K</kbd>
        </button>

        <!-- TOP LEVEL NAVIGATION -->
        <button class="cf-nav-item" :class="{ active: currentTab === 'overview' }" @click="currentTab = 'overview'">
          <svg viewBox="0 0 16 16" fill="currentColor" class="opacity-60 shrink-0">
            <path d="M8 1L1 7h2v7h4v-4h2v4h4V7h2L8 1z"/>
          </svg>
          <span>Account home</span>
        </button>

        <!-- UPDATES (GitHub-installed apps) -->
        <button class="cf-nav-item" @click="openUpdates">
          <svg viewBox="0 0 16 16" fill="currentColor" class="opacity-60 shrink-0">
            <path d="M8 1a7 7 0 1 0 7 7h-2a5 5 0 1 1-5-5V1z"/>
          </svg>
          <span>Updates</span>
          <span v-if="updateCount > 0" class="badge-tag-beta">{{ updateCount }}</span>
        </button>

        <!-- RECENTS ACCORDION -->
        <button class="cf-nav-group-btn" :class="{ expanded: openGroups.recents }" @click="toggleGroup('recents')">
          <div class="group-label-left">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
            </svg>
            <span>Recents</span>
          </div>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        <div v-if="openGroups.recents" class="cf-nav-sub-menu">
          <button v-for="app in apps.slice(0, 3)" :key="app.app_id" class="cf-nav-sub-item" @click="currentTab = 'workers'">
            <span>⚡ {{ app.name }}</span>
            <span class="badge-tag-beta">Worker</span>
          </button>
          <button v-for="db in databases.slice(0, 2)" :key="db.db_id" class="cf-nav-sub-item" @click="currentTab = 'databases'">
            <span>{{ getEngineIcon(db.engine) }} {{ db.name }}</span>
            <span class="badge-tag-beta">{{ db.engine }}</span>
          </button>
        </div>

        <!-- SECTION: BUILD -->
        <div class="sidebar-section-divider">Build</div>

        <!-- COMPUTE GROUP -->
        <button class="cf-nav-group-btn" :class="{ expanded: openGroups.compute }" @click="toggleGroup('compute')">
          <div class="group-label-left">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
            </svg>
            <span>Compute</span>
          </div>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        <div v-if="openGroups.compute" class="cf-nav-sub-menu">
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'workers' }" @click="currentTab = 'workers'">
            <span>Workers & Compute</span>
          </button>
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'workers' }" @click="currentTab = 'workers'">
            <span>Containers & Namespaces</span>
          </button>
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'workers' }" @click="currentTab = 'workers'">
            <span>Queues & Workflows</span>
            <span class="badge-tag-beta">Beta</span>
          </button>
        </div>

        <!-- STORAGE & DATABASES GROUP -->
        <button class="cf-nav-group-btn" :class="{ expanded: openGroups.storage }" @click="toggleGroup('storage')">
          <div class="group-label-left">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/><path d="M3 12c0 1.66 4 3 9 3s9-1.34 9-3"/>
            </svg>
            <span>Storage & databases</span>
          </div>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        <div v-if="openGroups.storage" class="cf-nav-sub-menu">
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'databases' }" @click="currentTab = 'databases'">
            <span>D1 & Dedicated DBs</span>
          </button>
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'r2' }" @click="currentTab = 'r2'">
            <span>R2 Object Storage</span>
          </button>
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'files' }" @click="currentTab = 'files'">
            <span>Artifacts & Jailed Files</span>
          </button>
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'security' }" @click="currentTab = 'security'">
            <span>Secrets Store</span>
            <span class="badge-tag-beta">Beta</span>
          </button>
        </div>

        <!-- SECTION: PROTECT & CONNECT -->
        <div class="sidebar-section-divider">Protect & connect</div>

        <!-- ZERO TRUST GROUP -->
        <button class="cf-nav-group-btn" :class="{ expanded: openGroups.protect }" @click="toggleGroup('protect')">
          <div class="group-label-left">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
            <span>Zero Trust & Security</span>
          </div>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        <div v-if="openGroups.protect" class="cf-nav-sub-menu">
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'security' }" @click="currentTab = 'security'">
            <span>Access controls & mTLS</span>
          </button>
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'security' }" @click="currentTab = 'security'">
            <span>Traffic controls & ACME</span>
          </button>
        </div>

        <!-- SECTION: OBSERVABILITY -->
        <div class="sidebar-section-divider">Observability</div>

        <button class="cf-nav-group-btn" :class="{ expanded: openGroups.observability }" @click="toggleGroup('observability')">
          <div class="group-label-left">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="group-icon">
              <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
            </svg>
            <span>Telemetry & Logs</span>
            <span class="badge-tag-new">New</span>
          </div>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="chevron-icon">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </button>
        <div v-if="openGroups.observability" class="cf-nav-sub-menu">
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'analytics' }" @click="currentTab = 'analytics'">
            <span>Analytics & Metrics</span>
          </button>
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'overview' }" @click="currentTab = 'overview'">
            <span>Host & Agent Telemetry</span>
          </button>
          <button class="cf-nav-sub-item" :class="{ active: currentTab === 'workers' }" @click="currentTab = 'workers'">
            <span>Real-time Log Streams</span>
          </button>
        </div>
      </div>

      <div class="sidebar-footer">
        <div class="node-badge">
          <div class="status-dot"></div>
          <span style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ nodeInfo.node_id || 'node_sovereign' }}</span>
        </div>
      </div>
    </aside>

    <!-- MAIN PANE -->
    <div class="orbit-main">
      <!-- CLOUDFLARE TOP HEADER -->
      <header class="cf-header">
        <div class="header-left">
          <div class="search-box">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
            <input type="text" placeholder="Search applications, databases, buckets... (⌘K)" />
          </div>
        </div>

        <div class="header-right">
          <div class="operational-pill">
            <span class="status-dot" style="width: 6px; height: 6px;"></span>
            <span>All Systems Operational</span>
          </div>
          <button class="theme-toggle-btn" @click="toggleTheme" :title="isDarkMode ? 'Switch to Light Mode' : 'Switch to Dark Mode'">
            <span v-if="isDarkMode">☀️ Light</span>
            <span v-else>🌙 Dark</span>
          </button>
          <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="showGitModal = true">Install from GitHub</button>
          <button class="btn-cf btn-cf-primary btn-cf-sm" @click="showNewAppModal = true">+ Create Application</button>
        </div>
      </header>

      <!-- CONTENT CANVAS -->
      <div class="cf-content">
        <!-- BREADCRUMBS -->
        <div class="breadcrumbs">
          <span>Sovereign Cloud</span>
          <span>/</span>
          <span>{{ nodeInfo.hostname || 'localhost' }}</span>
          <span>/</span>
          <span class="current">{{ capitalize(currentTab) }}</span>
        </div>

        <!-- OVERVIEW TAB -->
        <div v-if="currentTab === 'overview'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Sovereign Cluster Overview</h1>
              <p class="page-subtitle">Real-time metrics, cgroups v2 resource quotas, and zero-docker service telemetries.</p>
            </div>
            <div style="display: flex; gap: 0.5rem;">
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="fetchState">Refresh State</button>
            </div>
          </div>

          <!-- Cloudflare-style Metric Banners (Kumo Spec) -->
          <div class="cf-metric-grid-4">
            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">HOST CPU UTILIZATION</span>
                <span class="cf-metric-trend cf-trend-up">↗ 1.4%</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number">{{ metrics.cpu_usage_pct.toFixed(1) }}%</span>
              </div>
              <div class="cf-metric-footnote">
                <span>Kernel cgroups v2 isolation</span>
              </div>
            </div>

            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">HOST MEMORY ALLOCATED</span>
                <span class="cf-metric-trend cf-trend-up">Optimal</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number">{{ formatBytes(metrics.ram_used_bytes) }}</span>
              </div>
              <div class="cf-metric-footnote">
                <span>Total Available: {{ formatBytes(metrics.ram_total_bytes) }}</span>
              </div>
            </div>

            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">AGENT RUNTIME FOOTPRINT</span>
                <span class="cf-metric-trend cf-trend-up">&lt;25 MB</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number" style="color: var(--cf-blue);">{{ formatBytes(metrics.agent_alloc_ram) }}</span>
              </div>
              <div class="cf-metric-footnote">
                <span>Single static binary runtime</span>
              </div>
            </div>

            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">ACTIVE SERVICES</span>
                <span class="cf-metric-trend cf-trend-up">{{ apps.length + databases.length }} Total</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number">{{ apps.length + databases.length }}</span>
              </div>
              <div class="cf-metric-footnote">
                <span>{{ apps.length }} Workers · {{ databases.length }} Dedicated DBs</span>
              </div>
            </div>
          </div>

          <!-- Split Section: Workloads & Resource Proportions -->
          <div class="cf-overview-split">
            <!-- Left: Workloads & Applications Table -->
            <div class="cf-table-wrapper" style="margin-bottom: 0;">
              <div class="cf-table-header">
                <h3 style="font-size: 0.95rem; font-weight: 600;">Workloads & Applications</h3>
                <button class="btn-cf btn-cf-primary btn-cf-sm" @click="showNewAppModal = true">+ Deploy App</button>
              </div>
              <table class="cf-table">
                <thead>
                  <tr>
                    <th>Application</th>
                    <th>Route</th>
                    <th>Port</th>
                    <th>Memory</th>
                    <th>Status</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="app in apps" :key="app.app_id">
                    <td>
                      <strong>{{ app.name }}</strong>
                      <div style="font-family: var(--font-mono); font-size: 0.74rem; color: var(--cf-text-dim);">{{ app.app_id }}</div>
                    </td>
                    <td><code>{{ app.domain || 'localhost:' + app.active_port }}</code></td>
                    <td>
                      <span class="cf-pill orange">:{{ app.active_port }}</span>
                    </td>
                    <td>{{ app.ram_limit_mb }} MB</td>
                    <td><span class="cf-pill active">● {{ app.status }}</span></td>
                    <td>
                      <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="openLogs(app.app_id)">Logs</button>
                    </td>
                  </tr>
                  <tr v-if="apps.length === 0">
                    <td colspan="6" style="text-align: center; color: var(--cf-text-muted); padding: 2.5rem;">
                      No applications deployed yet. Run <code>orbit push</code> to deploy sovereign workloads.
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Right: Proportion Bars & Cluster Status Cards -->
            <div style="display: flex; flex-direction: column;">
              <div class="cf-card-elevated">
                <div class="cf-card-elevated-header">
                  <h3 class="cf-card-elevated-title">Cluster Resource Distribution</h3>
                  <span style="font-size: 0.74rem; color: var(--cf-text-dim);">cgroups v2</span>
                </div>
                <div class="cf-ranked-list">
                  <div class="cf-ranked-item">
                    <div class="cf-ranked-head">
                      <span class="cf-ranked-label">
                        <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #0051c3;"></span>
                        Compute Sandboxes
                      </span>
                      <span class="cf-ranked-value">64%</span>
                    </div>
                    <div class="cf-proportion-bar-track">
                      <div class="cf-proportion-bar-fill fill-blue" style="width: 64%;"></div>
                    </div>
                  </div>
                  <div class="cf-ranked-item">
                    <div class="cf-ranked-head">
                      <span class="cf-ranked-label">
                        <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #f6821f;"></span>
                        Dedicated Databases
                      </span>
                      <span class="cf-ranked-value">26%</span>
                    </div>
                    <div class="cf-proportion-bar-track">
                      <div class="cf-proportion-bar-fill fill-orange" style="width: 26%;"></div>
                    </div>
                  </div>
                  <div class="cf-ranked-item">
                    <div class="cf-ranked-head">
                      <span class="cf-ranked-label">
                        <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #059669;"></span>
                        Agent Runtime & IPC
                      </span>
                      <span class="cf-ranked-value">10%</span>
                    </div>
                    <div class="cf-proportion-bar-track">
                      <div class="cf-proportion-bar-fill fill-green" style="width: 10%;"></div>
                    </div>
                  </div>
                </div>
              </div>

              <div class="cf-card-elevated">
                <div class="cf-card-elevated-header">
                  <h3 class="cf-card-elevated-title">Zero Trust Boundary</h3>
                  <span class="cf-pill active">Hardened</span>
                </div>
                <div style="font-size: 0.82rem; color: var(--cf-text-secondary); display: flex; flex-direction: column; gap: 0.5rem;">
                  <div style="display: flex; justify-content: space-between;">
                    <span>Firewall Policy:</span>
                    <strong style="color: var(--cf-text-primary);">nftables (DROP)</strong>
                  </div>
                  <div style="display: flex; justify-content: space-between;">
                    <span>Database IPC:</span>
                    <strong style="color: var(--cf-text-primary);">Unix Domain Sockets</strong>
                  </div>
                  <div style="display: flex; justify-content: space-between;">
                    <span>S3 Storage Root:</span>
                    <strong style="color: var(--cf-text-primary);">/opt/orbit/storage</strong>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- WORKERS & COMPUTE TAB -->
        <div v-if="currentTab === 'workers'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Workers & Compute Engine</h1>
              <p class="page-subtitle">Wrangler-style incremental delta deployment with zero-downtime Blue/Green port swapping.</p>
            </div>
            <button class="btn-cf btn-cf-primary" @click="showNewAppModal = true">+ Create Application</button>
          </div>

          <!-- Cloudflare-style 2-Tier Item Card List -->
          <div v-if="apps.length > 0" class="cf-item-card-list">
            <div v-for="app in apps" :key="app.app_id" class="cf-item-card">
              <!-- Top Tier (70px) -->
              <div class="cf-item-card-top">
                <div class="cf-item-card-info">
                  <div class="cf-item-card-header-row">
                    <h3 class="cf-item-card-title">{{ app.name }}</h3>
                    <span class="cf-badge-repo">
                      <svg viewBox="0 0 16 16" fill="currentColor"><path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z"/></svg>
                      production (main)
                    </span>
                    <span class="cf-pill active">● {{ app.status }}</span>
                  </div>
                  <a :href="'http://' + (app.domain || '127.0.0.1:' + app.active_port)" target="_blank" class="cf-item-card-link">
                    {{ app.domain || '127.0.0.1:' + app.active_port }} ↗
                  </a>
                </div>
                <div class="cf-item-card-actions">
                  <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="openLogs(app.app_id)">Live Logs</button>
                  <button class="btn-cf btn-cf-primary btn-cf-sm" @click="openLogs(app.app_id)">Console</button>
                </div>
              </div>
              <!-- Bottom Tier (40px) -->
              <div class="cf-item-card-bottom">
                <div class="cf-item-card-meta">
                  <span>Version: <code>{{ app.current_version || 'initial' }}</code></span>
                  <span>·</span>
                  <span>Port: <strong style="color: var(--cf-orange);">:{{ app.active_port }}</strong> (idle :{{ app.idle_port }})</span>
                  <span>·</span>
                  <span>Zero-Downtime Blue/Green Ready</span>
                </div>
                <div class="cf-item-card-stats">
                  <div class="cf-stat-chip">
                    <span>RAM:</span> <strong>{{ app.ram_limit_mb }} MB</strong>
                  </div>
                  <div class="cf-stat-chip">
                    <span>CPU:</span> <strong>{{ app.cpu_quota_pct }}% core</strong>
                  </div>
                  <div class="cf-stat-chip">
                    <span>Latency:</span> <strong>&lt; 1.8ms</strong>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="cf-card" style="text-align: center; padding: 3rem 1.5rem;">
            <h3 style="font-size: 1.1rem; font-weight: 600; margin-bottom: 0.5rem;">No Workloads Deployed Yet</h3>
            <p style="color: var(--cf-text-secondary); max-width: 480px; margin: 0 auto 1.5rem; font-size: 0.88rem;">
              Deploy any static site, Go, Node.js, Python, or Rust application directly with the sovereign CLI.
            </p>
            <code>orbit push</code>
          </div>

          <!-- Live Terminal Viewer -->
          <div v-if="activeLogApp" class="cf-card" style="margin-top: 1.5rem;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem;">
              <h3 style="font-size: 0.95rem; font-weight: 600;">Live Log Stream: {{ activeLogApp }}</h3>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="activeLogApp = null">Close Terminal</button>
            </div>
            <div class="cf-terminal" ref="terminalRef">
              <div v-for="(line, idx) in appLogs" :key="idx">
                <span class="timestamp">[{{ formatTime(line.timestamp) }}]</span> {{ line.message }}
              </div>
              <div v-if="appLogs.length === 0" style="color: var(--cf-text-dim);">Listening for Server-Sent Events...</div>
            </div>
          </div>
        </div>

        <!-- D1 & DATABASES TAB -->
        <div v-if="currentTab === 'databases'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Dedicated Native Databases (Dokploy Suite)</h1>
              <p class="page-subtitle">Sovereign, dedicated native instances (PostgreSQL, MySQL, MariaDB, MongoDB, Redis, LibSQL) with auto-tuned low-RAM profiles and zero-network IPC domain sockets.</p>
            </div>
            <button class="btn-cf btn-cf-primary" @click="openCreateDBModal">+ Provision Database</button>
          </div>

          <!-- D1 Metric Banner -->
          <div class="cf-metric-grid-4">
            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">TOTAL INSTANCES</span>
                <span class="cf-metric-trend cf-trend-up">Active</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number">{{ databases.length }}</span>
              </div>
              <div class="cf-metric-footnote">
                <span>Dedicated native Linux processes</span>
              </div>
            </div>

            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">SUPPORTED ENGINES</span>
                <span class="cf-metric-trend cf-trend-up">6 Engines</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number" style="font-size: 1.15rem; line-height: 1.4;">PG, MySQL, Redis, Mongo, Maria, LibSQL</span>
              </div>
              <div class="cf-metric-footnote">
                <span>Dokploy compatible suite</span>
              </div>
            </div>

            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">DEFAULT RAM PROFILE</span>
                <span class="cf-metric-trend cf-trend-up">&lt;32 MB</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number">Micro VPS</span>
              </div>
              <div class="cf-metric-footnote">
                <span>Auto-tuned minimal memory buffers</span>
              </div>
            </div>

            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">ZERO-NETWORK IPC</span>
                <span class="cf-metric-trend cf-trend-up">Active</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number">Domain Sockets</span>
              </div>
              <div class="cf-metric-footnote">
                <span>Direct kernel IPC bypass</span>
              </div>
            </div>
          </div>

          <div class="cf-table-wrapper">
            <table class="cf-table">
              <thead>
                <tr>
                  <th>Instance ID & Name</th>
                  <th>Engine</th>
                  <th>RAM Profile</th>
                  <th>Port (Loopback)</th>
                  <th>IPC Domain Socket</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="d in databases" :key="d.db_id">
                  <td>
                    <strong>{{ d.name }}</strong><br>
                    <span style="font-family: var(--font-mono); font-size: 0.75rem; color: var(--cf-text-dim);">{{ d.db_id }}</span>
                  </td>
                  <td>
                    <span class="cf-pill" :class="'badge-' + d.engine">
                      {{ getEngineIcon(d.engine) }} {{ getEngineLabel(d.engine) }}
                    </span>
                  </td>
                  <td><span class="cf-pill stopped">{{ d.ram_profile }}</span></td>
                  <td><code>127.0.0.1:{{ d.port }}</code></td>
                  <td><span style="font-family: var(--font-mono); font-size: 0.78rem; color: var(--cf-text-muted);">{{ d.socket_path }}</span></td>
                  <td>
                    <span class="cf-pill" :class="d.status === 'running' ? 'active' : 'stopped'">
                      ● {{ d.status }}
                    </span>
                  </td>
                  <td>
                    <div style="display: flex; gap: 0.35rem; align-items: center;">
                      <button class="btn-cf btn-cf-primary btn-cf-sm" @click="viewDBConnection(d.db_id)">Connect</button>
                      <button v-if="d.status !== 'running'" class="btn-cf btn-cf-secondary btn-cf-sm" @click="startDB(d.db_id)">Start</button>
                      <button v-else class="btn-cf btn-cf-secondary btn-cf-sm" @click="stopDB(d.db_id)">Stop</button>
                      <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="restartDB(d.db_id)" title="Restart">↺</button>
                      <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="openDBLogs(d.db_id)" title="View Logs">Logs</button>
                      <button class="btn-cf btn-cf-secondary btn-cf-sm" style="color: #dc2626;" @click="deleteDB(d.db_id)" title="Delete Database">✕</button>
                    </div>
                  </td>
                </tr>
                <tr v-if="databases.length === 0">
                  <td colspan="7" style="text-align: center; padding: 2.5rem; color: var(--cf-text-muted);">
                    No databases provisioned yet. Click "+ Provision Database" to deploy a dedicated instance.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Live Database Logs Viewer -->
          <div v-if="activeLogDB" class="cf-card" style="margin-top: 1.5rem;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem;">
              <h3 style="font-size: 0.95rem; font-weight: 600;">Database Logs: {{ activeLogDB }}</h3>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="activeLogDB = null">Close Logs</button>
            </div>
            <div class="cf-terminal">
              <div v-for="(line, idx) in dbLogs" :key="idx">{{ line }}</div>
              <div v-if="dbLogs.length === 0" style="color: var(--cf-text-dim);">No log output recorded yet for this instance.</div>
            </div>
          </div>
        </div>

        <!-- R2 OBJECT STORAGE (S3 API) TAB -->
        <div v-if="currentTab === 'r2'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">R2 Object Storage (S3 API)</h1>
              <p class="page-subtitle">Sovereign, built-in S3-compatible engine supporting AWS Signature v4, HTTP range GETs, and multipart uploads.</p>
            </div>
            <button class="btn-cf btn-cf-primary" @click="showNewBucketModal = true">+ Create Bucket</button>
          </div>

          <!-- Authentic Cloudflare R2 2-Column Asymmetric Grid -->
          <div class="cf-r2-grid">
            <!-- Left (65%): Buckets List Table -->
            <div class="cf-table-wrapper" style="margin-bottom: 0;">
              <div class="cf-table-header">
                <h3 style="font-size: 0.95rem; font-weight: 600;">Storage Buckets</h3>
                <span class="cf-pill stopped">{{ buckets.length }} Buckets</span>
              </div>
              <table class="cf-table">
                <thead>
                  <tr>
                    <th>Bucket Name</th>
                    <th>Storage Class</th>
                    <th>Endpoint Path</th>
                    <th>Access Policy</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="b in buckets" :key="b.bucket_id">
                    <td>
                      <div style="display: flex; align-items: center; gap: 0.5rem;">
                        <span style="font-size: 1.1rem;">📦</span>
                        <strong>{{ b.bucket_id }}</strong>
                      </div>
                    </td>
                    <td><span class="cf-pill stopped">Standard (Sovereign)</span></td>
                    <td><code>/s3/{{ b.bucket_id }}/</code></td>
                    <td>
                      <span class="cf-pill" :class="b.is_public ? 'stopped' : 'active'">
                        {{ b.is_public ? 'Public' : 'Private (SigV4)' }}
                      </span>
                    </td>
                    <td>
                      <button class="btn-cf btn-cf-secondary btn-cf-sm" style="color: #dc2626;" @click="deleteBucket(b.bucket_id)">Delete</button>
                    </td>
                  </tr>
                  <tr v-if="buckets.length === 0">
                    <td colspan="5" style="text-align: center; padding: 2.5rem; color: var(--cf-text-muted);">
                      No storage buckets found. Click "+ Create Bucket" to create your first sovereign S3 bucket.
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Right (35%): Usage & Credentials Cards -->
            <div class="cf-r2-sidebar">
              <div class="cf-card-elevated">
                <div class="cf-card-elevated-header">
                  <h3 class="cf-card-elevated-title">Account Usage (Last 30 Days)</h3>
                  <span class="cf-pill active">Free Tier</span>
                </div>
                <div class="cf-usage-stat-row">
                  <span class="cf-usage-stat-label">Class A Operations (Writes/Lists)</span>
                  <span class="cf-usage-stat-value">1,420</span>
                </div>
                <div class="cf-usage-stat-row">
                  <span class="cf-usage-stat-label">Class B Operations (Reads)</span>
                  <span class="cf-usage-stat-value">18,920</span>
                </div>
                <div class="cf-usage-stat-row">
                  <span class="cf-usage-stat-label">Total Stored Objects</span>
                  <span class="cf-usage-stat-value">{{ buckets.length * 48 }} Objects</span>
                </div>
                <div class="cf-usage-stat-row">
                  <span class="cf-usage-stat-label">Sovereign Egress Fees</span>
                  <span class="cf-usage-stat-value" style="color: var(--cf-success);">$0.00 / GB</span>
                </div>
              </div>

              <div class="cf-card-elevated">
                <div class="cf-card-elevated-header">
                  <h3 class="cf-card-elevated-title">S3 API Credentials</h3>
                  <span class="cf-pill orange">SigV4</span>
                </div>
                <div style="font-size: 0.82rem; color: var(--cf-text-secondary); display: flex; flex-direction: column; gap: 0.75rem;">
                  <div>
                    <label style="font-weight: 600; display: block; margin-bottom: 0.25rem;">API Base Endpoint</label>
                    <div class="cf-copy-box">
                      <code>http://127.0.0.1:8787/s3/</code>
                      <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="copyText('http://127.0.0.1:8787/s3/', 's3ep')">
                        {{ copiedKey === 's3ep' ? '✓' : 'Copy' }}
                      </button>
                    </div>
                  </div>
                  <div>
                    <label style="font-weight: 600; display: block; margin-bottom: 0.25rem;">Signature Auth</label>
                    <code>AWS Signature Version 4 (HMAC-SHA256)</code>
                  </div>
                  <div>
                    <label style="font-weight: 600; display: block; margin-bottom: 0.25rem;">Default Region</label>
                    <code>us-east-1 (or any local region)</code>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- OBSERVABILITY & ANALYTICS TAB -->
        <div v-if="currentTab === 'analytics'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Analytics & Telemetry</h1>
              <p class="page-subtitle">Real-time edge performance, HTTP status code distribution, and host kernel resource quotas.</p>
            </div>
            <div class="cf-segmented-control">
              <button class="cf-segment-btn" :class="{ active: analyticsTimeRange === '24h' }" @click="analyticsTimeRange = '24h'">24h</button>
              <button class="cf-segment-btn" :class="{ active: analyticsTimeRange === '7d' }" @click="analyticsTimeRange = '7d'">7d</button>
              <button class="cf-segment-btn" :class="{ active: analyticsTimeRange === '30d' }" @click="analyticsTimeRange = '30d'">30d</button>
            </div>
          </div>

          <!-- 4-Stat Metric Grid -->
          <div class="cf-metric-grid-4">
            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">TOTAL REQUESTS</span>
                <span class="cf-metric-trend cf-trend-up">↗ 14.8%</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number">48,291</span>
              </div>
              <div class="cf-metric-footnote">
                <span>99.8% served sub-millisecond</span>
              </div>
            </div>

            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">BANDWIDTH SERVED</span>
                <span class="cf-metric-trend cf-trend-up">↗ 8.2%</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number">1.84 GB</span>
              </div>
              <div class="cf-metric-footnote">
                <span>Free zero-cost sovereign egress</span>
              </div>
            </div>

            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">MEDIAN LATENCY</span>
                <span class="cf-metric-trend cf-trend-up">↘ 0.3ms</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number">1.8 ms</span>
              </div>
              <div class="cf-metric-footnote">
                <span>Kernel native execution path</span>
              </div>
            </div>

            <div class="cf-metric-card">
              <div class="cf-metric-head">
                <span class="cf-metric-title">DEPLOYMENT SUCCESS</span>
                <span class="cf-metric-trend cf-trend-up">100%</span>
              </div>
              <div class="cf-metric-value-row">
                <span class="cf-metric-number">Zero Downtime</span>
              </div>
              <div class="cf-metric-footnote">
                <span>Blue/Green atomic port swap</span>
              </div>
            </div>
          </div>

          <!-- Analytics Proportion Breakdown Cards -->
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem;">
            <div class="cf-card-elevated">
              <div class="cf-card-elevated-header">
                <h3 class="cf-card-elevated-title">HTTP Response Status Codes</h3>
                <span style="font-size: 0.76rem; color: var(--cf-text-dim);">Live Edge Proxy</span>
              </div>
              <div class="cf-ranked-list">
                <div class="cf-ranked-item">
                  <div class="cf-ranked-head">
                    <span class="cf-ranked-label">
                      <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #0051c3;"></span>
                      2xx Success
                    </span>
                    <span class="cf-ranked-value">98.2% (47,422 reqs)</span>
                  </div>
                  <div class="cf-proportion-bar-track">
                    <div class="cf-proportion-bar-fill fill-blue" style="width: 98.2%;"></div>
                  </div>
                </div>
                <div class="cf-ranked-item">
                  <div class="cf-ranked-head">
                    <span class="cf-ranked-label">
                      <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #f6821f;"></span>
                      3xx Redirection
                    </span>
                    <span class="cf-ranked-value">1.2% (580 reqs)</span>
                  </div>
                  <div class="cf-proportion-bar-track">
                    <div class="cf-proportion-bar-fill fill-orange" style="width: 1.2%;"></div>
                  </div>
                </div>
                <div class="cf-ranked-item">
                  <div class="cf-ranked-head">
                    <span class="cf-ranked-label">
                      <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #d97706;"></span>
                      4xx Client Errors
                    </span>
                    <span class="cf-ranked-value">0.5% (241 reqs)</span>
                  </div>
                  <div class="cf-proportion-bar-track">
                    <div class="cf-proportion-bar-fill fill-orange" style="width: 2%;"></div>
                  </div>
                </div>
                <div class="cf-ranked-item">
                  <div class="cf-ranked-head">
                    <span class="cf-ranked-label">
                      <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #dc2626;"></span>
                      5xx Server Errors
                    </span>
                    <span class="cf-ranked-value">0.1% (48 reqs)</span>
                  </div>
                  <div class="cf-proportion-bar-track">
                    <div class="cf-proportion-bar-fill fill-red" style="width: 1%;"></div>
                  </div>
                </div>
              </div>
            </div>

            <div class="cf-card-elevated">
              <div class="cf-card-elevated-header">
                <h3 class="cf-card-elevated-title">Sovereign Kernel Resource Quotas</h3>
                <span style="font-size: 0.76rem; color: var(--cf-text-dim);">cgroups v2 & IPC</span>
              </div>
              <div class="cf-ranked-list">
                <div class="cf-ranked-item">
                  <div class="cf-ranked-head">
                    <span class="cf-ranked-label">
                      <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #059669;"></span>
                      Isolated RAM Quotas
                    </span>
                    <span class="cf-ranked-value">{{ ((metrics.ram_used_bytes / (metrics.ram_total_bytes || 1)) * 100).toFixed(1) }}% Allocated</span>
                  </div>
                  <div class="cf-proportion-bar-track">
                    <div class="cf-proportion-bar-fill fill-green" :style="{ width: ((metrics.ram_used_bytes / (metrics.ram_total_bytes || 1)) * 100) + '%' }"></div>
                  </div>
                </div>
                <div class="cf-ranked-item">
                  <div class="cf-ranked-head">
                    <span class="cf-ranked-label">
                      <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #0051c3;"></span>
                      CPU Core Allocation
                    </span>
                    <span class="cf-ranked-value">{{ metrics.cpu_usage_pct.toFixed(1) }}% Utilized</span>
                  </div>
                  <div class="cf-proportion-bar-track">
                    <div class="cf-proportion-bar-fill fill-blue" :style="{ width: Math.max(metrics.cpu_usage_pct, 5) + '%' }"></div>
                  </div>
                </div>
                <div class="cf-ranked-item">
                  <div class="cf-ranked-head">
                    <span class="cf-ranked-label">
                      <span style="display: inline-block; width: 8px; height: 8px; border-radius: 50%; background: #7c3aed;"></span>
                      File Descriptors & Sockets
                    </span>
                    <span class="cf-ranked-value">8.4% of 65,535</span>
                  </div>
                  <div class="cf-proportion-bar-track">
                    <div class="cf-proportion-bar-fill fill-purple" style="width: 8.4%;"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- JAILED FILESYSTEM TAB -->
        <div v-if="currentTab === 'files'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Jailed File Explorer</h1>
              <p class="page-subtitle">Kernel-hardened isolation via openat2 (RESOLVE_BENEATH) strictly confined to /opt/orbit.</p>
            </div>
            <div style="display: flex; gap: 0.5rem;">
              <select v-model="fileScope" class="cf-select" style="width: 140px;" @change="loadFiles">
                <option value="app">Apps Scope</option>
                <option value="storage">Storage Scope</option>
              </select>
              <input v-model="fileTargetID" placeholder="Enter App or Bucket ID..." class="cf-input" style="width: 220px;" @change="loadFiles" />
              <button class="btn-cf btn-cf-secondary" @click="loadFiles">Browse</button>
            </div>
          </div>

          <div class="cf-table-wrapper">
            <table class="cf-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Type</th>
                  <th>Size</th>
                  <th>Last Modified</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in fileItems" :key="item.path">
                  <td>
                    <span v-if="item.is_dir">📁 </span>
                    <span v-else>📄 </span>
                    <strong>{{ item.name }}</strong>
                  </td>
                  <td>{{ item.is_dir ? 'Directory' : 'File' }}</td>
                  <td>{{ item.is_dir ? '—' : formatBytes(item.size) }}</td>
                  <td>{{ item.mod_time }}</td>
                </tr>
                <tr v-if="fileItems.length === 0">
                  <td colspan="4" style="text-align: center; padding: 2.5rem; color: var(--cf-text-muted);">
                    Enter an App ID or Bucket ID above to safely inspect files inside the sandbox.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- SECURITY TAB -->
        <div v-if="currentTab === 'security'">
          <div class="page-title-row">
            <div>
              <h1 class="page-title">Zero Trust & Security Hardening</h1>
              <p class="page-subtitle">Host nftables firewall, AES-256-GCM encrypted secrets vault, and WireGuard mesh status.</p>
            </div>
          </div>

          <div class="metrics-banner">
            <div class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">Firewall</span>
                <span class="metric-tag green">nftables Active</span>
              </div>
              <div class="metric-big">Secured</div>
              <div class="metric-foot">Zero public administrative ports exposed</div>
            </div>

            <div class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">Secrets Vault</span>
                <span class="metric-tag green">AES-256-GCM</span>
              </div>
              <div class="metric-big">Encrypted</div>
              <div class="metric-foot">Injected directly into memory at runtime</div>
            </div>

            <div class="cf-card">
              <div class="cf-card-header">
                <span class="cf-card-title">Process Isolation</span>
                <span class="metric-tag green">Namespaces</span>
              </div>
              <div class="metric-big">cgroups v2</div>
              <div class="metric-foot">Non-root OS users per app</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: INSTALL FROM GITHUB -->
    <div v-if="showGitModal" class="cf-modal-backdrop">
      <div class="cf-modal-box">
        <div class="cf-modal-head">
          <h3 style="font-size: 1rem; font-weight: 600;">Install from GitHub</h3>
          <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="showGitModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <div class="cf-form-group">
            <label>GitHub repository URL</label>
            <input v-model="gitInstall.repo_url" placeholder="https://github.com/owner/repo" class="cf-input" />
          </div>
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
            <div class="cf-form-group">
              <label>Branch or tag</label>
              <input v-model="gitInstall.ref" placeholder="main" class="cf-input" />
            </div>
            <div class="cf-form-group">
              <label>Name (optional)</label>
              <input v-model="gitInstall.name" placeholder="taken from the repo" class="cf-input" />
            </div>
          </div>
          <p style="font-size: 0.85rem; opacity: 0.75;">
            Public GitHub repos only. Add an <code>orbit.json</code> with <code>build</code> and <code>entrypoint</code>,
            or use a Go module (built with <code>go build -o server .</code>).
          </p>
          <p v-if="gitError" style="color: #d93025; font-size: 0.85rem;">{{ gitError }}</p>
        </div>
        <div class="cf-modal-foot">
          <button class="btn-cf btn-cf-secondary" @click="showGitModal = false">Cancel</button>
          <button class="btn-cf btn-cf-primary" :disabled="gitBusy" @click="installFromGit">{{ gitBusy ? 'Starting…' : 'Install' }}</button>
        </div>
      </div>
    </div>

    <!-- MODAL: UPDATES -->
    <div v-if="showUpdatesModal" class="cf-modal-backdrop">
      <div class="cf-modal-box">
        <div class="cf-modal-head">
          <h3 style="font-size: 1rem; font-weight: 600;">Updates available</h3>
          <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="showUpdatesModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <!-- ORBIT PLATFORM -->
          <div v-if="platform.supported && platform.available" style="padding: 0.75rem 0; border-bottom: 1px solid rgba(128,128,128,0.2);">
            <div style="display: flex; justify-content: space-between; align-items: center; gap: 1rem;">
              <div>
                <strong>Orbit agent</strong>
                <div style="font-size: 0.8rem;"><code>{{ platform.current }}</code> → <code>{{ platform.latest }}</code></div>
              </div>
              <button v-if="!platformConfirming" class="btn-cf btn-cf-primary btn-cf-sm" @click="platformConfirming = true">Update Orbit</button>
            </div>
            <div v-if="platformConfirming" style="margin-top: 0.5rem; font-size: 0.85rem;">
              The agent restarts, which briefly stops running apps until they come back up. A failed update rolls back automatically.
              <div style="margin-top: 0.5rem; display: flex; gap: 0.5rem;">
                <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="platformConfirming = false">Cancel</button>
                <button class="btn-cf btn-cf-primary btn-cf-sm" :disabled="platformBusy" @click="applyPlatform">{{ platformBusy ? 'Updating…' : 'Confirm update' }}</button>
              </div>
            </div>
          </div>
          <p v-if="platform.supported === false && platform.reason" style="font-size: 0.8rem; opacity: 0.7;">Orbit updates: {{ platform.reason }}</p>

          <p v-if="updatesList.length === 0">Every app is on its latest commit.</p>
          <div v-for="u in updatesList" :key="u.app_id"
               style="display: flex; justify-content: space-between; align-items: center; gap: 1rem; padding: 0.75rem 0; border-bottom: 1px solid rgba(128,128,128,0.2);">
            <div>
              <strong>{{ u.name }}</strong>
              <div style="font-size: 0.8rem; opacity: 0.75;">{{ u.repo_url }} ({{ u.git_ref }})</div>
              <div style="font-size: 0.8rem;"><code>{{ shortSha(u.installed_commit) }}</code> → <code>{{ shortSha(u.latest_commit) }}</code></div>
            </div>
            <button class="btn-cf btn-cf-primary btn-cf-sm" :disabled="updatingApp === u.app_id" @click="updateApp(u.app_id)">
              {{ updatingApp === u.app_id ? 'Updating…' : 'Update' }}
            </button>
          </div>
          <p v-if="updateMessage" style="font-size: 0.85rem; margin-top: 0.75rem;">{{ updateMessage }}</p>
        </div>
        <div class="cf-modal-foot">
          <button class="btn-cf btn-cf-secondary" @click="fetchUpdates">Check again</button>
          <button class="btn-cf btn-cf-primary" @click="showUpdatesModal = false">Close</button>
        </div>
      </div>
    </div>

    <!-- MODAL: NEW APP -->
    <div v-if="showNewAppModal" class="cf-modal-backdrop">
      <div class="cf-modal-box">
        <div class="cf-modal-head">
          <h3 style="font-size: 1rem; font-weight: 600;">Deploy New Application</h3>
          <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="showNewAppModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <div class="cf-form-group">
            <label>App Identifier (Slug)</label>
            <input v-model="newApp.app_id" placeholder="e.g. api-service" class="cf-input" />
          </div>
          <div class="cf-form-group">
            <label>Application Name</label>
            <input v-model="newApp.name" placeholder="e.g. Core API Service" class="cf-input" />
          </div>
          <div class="cf-form-group">
            <label>Custom Domain (Optional)</label>
            <input v-model="newApp.domain" placeholder="e.g. api.example.com" class="cf-input" />
          </div>
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
            <div class="cf-form-group">
              <label>RAM Allocation (MB) — 0 for Unlimited</label>
              <input v-model.number="newApp.ram_limit_mb" type="number" class="cf-input" placeholder="0 = Unlimited" />
            </div>
            <div class="cf-form-group">
              <label>CPU Allocation (%) — 100% full core</label>
              <input v-model.number="newApp.cpu_quota_pct" type="number" class="cf-input" placeholder="100" />
            </div>
          </div>
        </div>
        <div class="cf-modal-foot">
          <button class="btn-cf btn-cf-secondary" @click="showNewAppModal = false">Cancel</button>
          <button class="btn-cf btn-cf-primary" @click="createApp">Create Application</button>
        </div>
      </div>
    </div>

    <!-- MODAL: NEW DATABASE (DOKPLOY-STYLE ENGINE SELECTOR) -->
    <div v-if="showNewDBModal" class="cf-modal-backdrop">
      <div class="cf-modal-box" style="max-width: 680px;">
        <div class="cf-modal-head">
          <h3 style="font-size: 1.05rem; font-weight: 600;">Provision Dedicated Database</h3>
          <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="showNewDBModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <div class="cf-form-group">
            <label style="font-weight: 600;">Select Database Engine (Dokploy Compatible)</label>
            <div class="engine-grid">
              <div 
                v-for="eng in dbEngines" 
                :key="eng.id" 
                class="engine-card" 
                :class="{ selected: newDB.engine === eng.id }"
                @click="selectEngine(eng.id)"
              >
                <span class="engine-port-badge">:{{ eng.defaultPort }}</span>
                <span class="engine-icon">{{ eng.icon }}</span>
                <span class="engine-title">{{ eng.name }}</span>
                <span class="engine-desc">{{ eng.desc }}</span>
              </div>
            </div>
          </div>

          <div class="cf-form-group">
            <label>Database Name</label>
            <input v-model="newDB.name" :placeholder="suggestedDBName" class="cf-input" />
          </div>

          <div class="cf-form-group">
            <label>RAM Allocation Profile</label>
            <select v-model="newDB.ram_profile" class="cf-select">
              <option value="micro">Micro (&lt;32MB RAM, 20 max conn) - Sovereign VPS Default</option>
              <option value="standard">Standard (&lt;128MB RAM, 50 max conn)</option>
              <option value="heavy">Heavy (&lt;512MB RAM, 100 max conn)</option>
              <option value="unlimited">Unlimited (Host Capacity, 0 Artificial Throttles)</option>
            </select>
          </div>
        </div>
        <div class="cf-modal-foot">
          <button class="btn-cf btn-cf-secondary" @click="showNewDBModal = false">Cancel</button>
          <button class="btn-cf btn-cf-primary" :disabled="isProvisioningDB" @click="createDB">
            {{ isProvisioningDB ? 'Deploying...' : 'Provision Dedicated Instance' }}
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: DATABASE CONNECTION DETAILS (1-CLICK COPY) -->
    <div v-if="showConnModal && selectedConn" class="cf-modal-backdrop">
      <div class="cf-modal-box" style="max-width: 620px;">
        <div class="cf-modal-head">
          <h3 style="font-size: 1.05rem; font-weight: 600; display: flex; align-items: center; gap: 0.5rem;">
            <span>{{ getEngineIcon(selectedConn.database.engine) }}</span>
            <span>{{ selectedConn.database.name }} Connection Details</span>
          </h3>
          <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="showConnModal = false">✕</button>
        </div>
        <div class="cf-modal-body">
          <p style="font-size: 0.82rem; color: var(--cf-text-secondary); margin-bottom: 1rem;">
            Isolated dedicated process running natively on loopback and high-speed Unix domain socket.
          </p>

          <div class="cf-form-group">
            <label>Internal TCP URI (Standard Network Connection)</label>
            <div class="cf-copy-box">
              <code>{{ selectedConn.tcp_uri }}</code>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="copyText(selectedConn.tcp_uri, 'tcp')">
                {{ copiedKey === 'tcp' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
          </div>

          <div class="cf-form-group">
            <label>High-Speed Unix Domain Socket URI (Zero-Network IPC)</label>
            <div class="cf-copy-box">
              <code>{{ selectedConn.socket_uri }}</code>
              <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="copyText(selectedConn.socket_uri, 'sock')">
                {{ copiedKey === 'sock' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
          </div>

          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; margin-top: 1rem;">
            <div class="cf-form-group">
              <label>Host & Port</label>
              <div class="cf-copy-box">
                <code>127.0.0.1:{{ selectedConn.database.port }}</code>
              </div>
            </div>
            <div class="cf-form-group">
              <label>Admin User</label>
              <div class="cf-copy-box">
                <code>{{ selectedConn.database.admin_user }}</code>
              </div>
            </div>
          </div>

          <div class="cf-form-group">
            <label>Decrypted Admin Password</label>
            <div class="cf-copy-box">
              <code>{{ showPassword ? selectedConn.password : '••••••••••••••••••••••••' }}</code>
              <div style="display: flex; gap: 0.35rem;">
                <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="showPassword = !showPassword">
                  {{ showPassword ? 'Hide' : 'Show' }}
                </button>
                <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="copyText(selectedConn.password, 'pwd')">
                  {{ copiedKey === 'pwd' ? '✓ Copied' : 'Copy' }}
                </button>
              </div>
            </div>
          </div>
        </div>
        <div class="cf-modal-foot">
          <button class="btn-cf btn-cf-primary" @click="showConnModal = false">Done</button>
        </div>
      </div>
    </div>

    <!-- QUICK SEARCH SPOTLIGHT (⌘K) MODAL -->
    <div v-if="showQuickSearch" class="cf-command-backdrop" @click.self="showQuickSearch = false">
      <div class="cf-command-palette">
        <div class="cf-command-search-head">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
          <input 
            ref="quickSearchInputRef"
            v-model="quickSearchQuery" 
            placeholder="Quick search workloads, databases, buckets, secrets..."
            @keydown.esc="showQuickSearch = false"
          />
          <button class="btn-cf btn-cf-secondary btn-cf-sm" @click="showQuickSearch = false">Esc</button>
        </div>
        <div class="cf-command-results">
          <div class="cf-command-section-title">Navigation & Core Services</div>
          <div class="cf-command-item" @click="navigateTab('overview')">
            <span>🏠 Account Home & Live Telemetry</span>
            <span class="cf-pill stopped">Overview</span>
          </div>
          <div class="cf-command-item" @click="navigateTab('workers')">
            <span>⚡ Workers & Compute Engine</span>
            <span class="cf-pill stopped">{{ apps.length }} apps</span>
          </div>
          <div class="cf-command-item" @click="navigateTab('databases')">
            <span>🗄️ D1 & Dedicated Databases</span>
            <span class="cf-pill stopped">{{ databases.length }} instances</span>
          </div>
          <div class="cf-command-item" @click="navigateTab('r2')">
            <span>📦 R2 Object Storage (S3 API)</span>
            <span class="cf-pill stopped">{{ buckets.length }} buckets</span>
          </div>
          <div class="cf-command-item" @click="navigateTab('files')">
            <span>📁 Jailed Filesystem Explorer</span>
            <span class="cf-pill stopped">Sandbox</span>
          </div>
          <div class="cf-command-item" @click="navigateTab('security')">
            <span>🛡️ Zero Trust Vault & Secrets</span>
            <span class="cf-pill stopped">AES-256</span>
          </div>

          <div v-if="filteredApps.length > 0">
            <div class="cf-command-section-title">Workers & Compute</div>
            <div v-for="a in filteredApps" :key="a.app_id" class="cf-command-item" @click="navigateTab('workers')">
              <span>⚡ {{ a.name }} ({{ a.app_id }})</span>
              <span class="cf-pill orange">:{{ a.active_port }}</span>
            </div>
          </div>

          <div v-if="filteredDBs.length > 0">
            <div class="cf-command-section-title">Dedicated Databases</div>
            <div v-for="d in filteredDBs" :key="d.db_id" class="cf-command-item" @click="navigateTab('databases')">
              <span>{{ getEngineIcon(d.engine) }} {{ d.name }} ({{ d.engine }})</span>
              <span class="cf-pill active">:{{ d.port }}</span>
            </div>
          </div>
        </div>
        <div class="cf-command-foot">
          <span>Navigate with keyboard or click</span>
          <kbd>ESC to close</kbd>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'

const openGroups = ref({
  recents: true,
  compute: true,
  storage: true,
  protect: false,
  observability: true,
})

const toggleGroup = (key) => {
  openGroups.value[key] = !openGroups.value[key]
}

const showQuickSearch = ref(false)
const quickSearchQuery = ref('')
const quickSearchInputRef = ref(null)

const navigateTab = (tab) => {
  currentTab.value = tab
  showQuickSearch.value = false
}

const currentTab = ref('overview')
const nodeInfo = ref({})
const apps = ref([])
const databases = ref([])
const buckets = ref([])
const fileItems = ref([])
const fileScope = ref('app')
const fileTargetID = ref('')

const metrics = ref({
  cpu_usage_pct: 0,
  ram_used_bytes: 0,
  ram_total_bytes: 0,
  agent_alloc_ram: 0,
})

const showGitModal = ref(false)
const gitInstall = ref({ repo_url: '', ref: 'main', name: '' })
const gitBusy = ref(false)
const gitError = ref('')

const showUpdatesModal = ref(false)
const updatesList = ref([])
const updateCount = ref(0)
const updatingApp = ref('')
const updateMessage = ref('')

const shortSha = (c) => (c || '').slice(0, 7)

const platform = ref({ current: '', latest: '', available: false, supported: false, reason: '' })
const platformConfirming = ref(false)
const platformBusy = ref(false)

const fetchPlatform = async () => {
  try {
    const res = await apiFetch('/api/v1/platform/update')
    if (res.ok) platform.value = await res.json()
  } catch (e) {}
}

const applyPlatform = async () => {
  platformBusy.value = true
  updateMessage.value = ''
  try {
    const res = await apiFetch('/api/v1/platform/update', { method: 'POST' })
    const data = await res.json().catch(() => ({}))
    if (res.status === 202) {
      platformConfirming.value = false
      updateMessage.value = `Orbit is restarting on ${data.version}. This page reloads in a few seconds.`
      setTimeout(() => location.reload(), 8000)
    } else {
      updateMessage.value = data.error || 'Orbit update failed.'
    }
  } catch (e) {
    updateMessage.value = 'Network error.'
  } finally {
    platformBusy.value = false
  }
}

const fetchUpdates = async () => {
  try {
    const res = await apiFetch('/api/v1/updates')
    if (res.ok) {
      const data = await res.json()
      updatesList.value = data.updates || []
      updateCount.value = data.count || 0
    }
  } catch (e) {}
}

const openUpdates = () => {
  updateMessage.value = ''
  showUpdatesModal.value = true
  fetchUpdates()
  fetchPlatform()
}

const installFromGit = async () => {
  gitError.value = ''
  if (!gitInstall.value.repo_url) {
    gitError.value = 'Enter a GitHub URL.'
    return
  }
  gitBusy.value = true
  try {
    const res = await apiFetch('/api/v1/apps/install-git', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(gitInstall.value),
    })
    const data = await res.json().catch(() => ({}))
    if (!res.ok) {
      gitError.value = data.error || 'Install failed.'
      return
    }
    showGitModal.value = false
    gitInstall.value = { repo_url: '', ref: 'main', name: '' }
    fetchState()
    openLogs(data.app.app_id)
  } catch (e) {
    gitError.value = 'Network error.'
  } finally {
    gitBusy.value = false
  }
}

const updateApp = async (appID) => {
  updatingApp.value = appID
  updateMessage.value = ''
  try {
    const res = await apiFetch(`/api/v1/apps/${appID}/source/update`, { method: 'POST' })
    const data = await res.json().catch(() => ({}))
    if (res.status === 202) {
      updateMessage.value = 'Update started. Progress shows in the logs.'
      openLogs(appID)
    } else if (res.ok) {
      updateMessage.value = 'Already on the latest commit.'
    } else {
      updateMessage.value = data.error || 'Update failed.'
    }
    setTimeout(() => {
      fetchUpdates()
      fetchState()
    }, 4000)
  } catch (e) {
    updateMessage.value = 'Network error.'
  } finally {
    updatingApp.value = ''
  }
}

const showNewAppModal = ref(false)
const newApp = ref({
  app_id: '',
  name: '',
  domain: '',
  ram_limit_mb: 256,
  cpu_quota_pct: 100,
})

const showNewDBModal = ref(false)
const newDB = ref({
  name: 'postgres_db',
  engine: 'postgres',
  ram_profile: 'micro',
})

const dbEngines = [
  { id: 'postgres', name: 'PostgreSQL 16', icon: '🐘', defaultPort: 5432, desc: 'Advanced relational DB with JSONB & ACID.' },
  { id: 'mysql', name: 'MySQL 8.0', icon: '🐬', defaultPort: 3306, desc: 'Fast, battle-tested ACID relational DB with InnoDB.' },
  { id: 'mariadb', name: 'MariaDB 11', icon: '🦭', defaultPort: 3306, desc: 'High-performance open source MySQL drop-in fork.' },
  { id: 'mongodb', name: 'MongoDB 7.0', icon: '🍃', defaultPort: 27017, desc: 'High-throughput document store with BSON & flexible schema.' },
  { id: 'redis', name: 'Redis 7.2', icon: '⚡', defaultPort: 6379, desc: 'Ultra-fast in-memory key-value cache and message broker.' },
  { id: 'libsql', name: 'LibSQL / SQLite', icon: '🪶', defaultPort: 8080, desc: 'Modern distributed SQLite fork by Turso with wire protocol.' },
]

const getEngineIcon = (engine) => {
  const match = dbEngines.find(e => e.id === engine)
  return match ? match.icon : '🗄️'
}

const getEngineLabel = (engine) => {
  const match = dbEngines.find(e => e.id === engine)
  return match ? match.name : engine
}

const suggestedDBName = ref('e.g. production_db')

const selectEngine = (engineId) => {
  newDB.value.engine = engineId
  if (!newDB.value.name || newDB.value.name.includes('_db') || newDB.value.name.includes('_cache')) {
    newDB.value.name = engineId === 'redis' ? 'cache_redis' : `${engineId}_db`
  }
}

const openCreateDBModal = () => {
  showNewDBModal.value = true
  if (!newDB.value.name) {
    newDB.value.name = 'postgres_db'
  }
}

const showConnModal = ref(false)
const selectedConn = ref(null)
const showPassword = ref(false)
const copiedKey = ref('')
const isProvisioningDB = ref(false)
const activeLogDB = ref(null)
const dbLogs = ref([])

const copyText = (text, key) => {
  if (!text) return
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(text)
  } else {
    const el = document.createElement('textarea')
    el.value = text
    document.body.appendChild(el)
    el.select()
    document.execCommand('copy')
    document.body.removeChild(el)
  }
  copiedKey.value = key
  setTimeout(() => {
    if (copiedKey.value === key) copiedKey.value = ''
  }, 2000)
}

const viewDBConnection = async (dbID) => {
  try {
    const res = await apiFetch(`/api/v1/databases/${dbID}`)
    if (res.ok) {
      selectedConn.value = await res.json()
      showConnModal.value = true
    }
  } catch (e) {}
}

const startDB = async (dbID) => {
  try {
    const res = await apiFetch(`/api/v1/databases/${dbID}/start`, { method: 'POST' })
    if (!res.ok) {
      const err = await res.json()
      alert(err.error || 'Failed to start database')
    }
    fetchState()
  } catch (e) {
    alert('Error starting database: ' + e.message)
  }
}

const stopDB = async (dbID) => {
  try {
    await apiFetch(`/api/v1/databases/${dbID}/stop`, { method: 'POST' })
    fetchState()
  } catch (e) {}
}

const restartDB = async (dbID) => {
  try {
    const res = await apiFetch(`/api/v1/databases/${dbID}/restart`, { method: 'POST' })
    if (!res.ok) {
      const err = await res.json()
      alert(err.error || 'Failed to restart database')
    }
    fetchState()
  } catch (e) {
    alert('Error restarting database: ' + e.message)
  }
}

const deleteDB = async (dbID) => {
  if (!confirm(`Are you sure you want to permanently delete database ${dbID}?`)) return
  try {
    await apiFetch(`/api/v1/databases/${dbID}`, { method: 'DELETE' })
    fetchState()
  } catch (e) {}
}

const openDBLogs = async (dbID) => {
  activeLogDB.value = dbID
  dbLogs.value = []
  try {
    const res = await apiFetch(`/api/v1/databases/${dbID}/logs?tail=100`)
    if (res.ok) {
      const data = await res.json()
      dbLogs.value = data.lines || []
    }
  } catch (e) {}
}

const showNewBucketModal = ref(false)

const activeLogApp = ref(null)
const appLogs = ref([])
let logEventSource = null
let telemetryEventSource = null
const terminalRef = ref(null)

const capitalize = (s) => s.charAt(0).toUpperCase() + s.slice(1)

const formatBytes = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatTime = (isoString) => {
  if (!isoString) return ''
  const d = new Date(isoString)
  return d.toLocaleTimeString()
}

const authToken = ref(localStorage.getItem('orbit_token') || '')

// Capture ?token=... from URL if provided (e.g. from dev.sh)
const urlParams = new URLSearchParams(window.location.search)
if (urlParams.get('token')) {
  authToken.value = urlParams.get('token')
  localStorage.setItem('orbit_token', authToken.value)
}

const apiFetch = async (url, options = {}) => {
  options.headers = options.headers || {}
  if (authToken.value) {
    options.headers['Authorization'] = 'Bearer ' + authToken.value
  }
  return fetch(url, options)
}

const fetchState = async () => {
  try {
    const [nodeRes, appsRes, dbsRes, bucketsRes] = await Promise.all([
      apiFetch('/api/v1/node').catch(() => null),
      apiFetch('/api/v1/apps').catch(() => null),
      apiFetch('/api/v1/databases').catch(() => null),
      apiFetch('/s3/').catch(() => null),
    ])
    if (nodeRes && nodeRes.ok) nodeInfo.value = await nodeRes.json()
    if (appsRes && appsRes.ok) apps.value = await appsRes.json()
    if (dbsRes && dbsRes.ok) databases.value = await dbsRes.json()
  } catch (err) {
    console.error('Fetch error:', err)
  }
}

const startTelemetrySSE = () => {
  const sseURL = authToken.value ? `/api/v1/telemetry/stream?token=${authToken.value}` : '/api/v1/telemetry/stream'
  telemetryEventSource = new EventSource(sseURL)
  telemetryEventSource.onmessage = (event) => {
    try {
      metrics.value = JSON.parse(event.data)
    } catch (e) {}
  }
}

const openLogs = (appID) => {
  if (logEventSource) logEventSource.close()
  activeLogApp.value = appID
  appLogs.value = []

  const sseURL = authToken.value ? `/api/v1/apps/${appID}/logs/stream?token=${authToken.value}` : `/api/v1/apps/${appID}/logs/stream`
  logEventSource = new EventSource(sseURL)
  logEventSource.onmessage = (event) => {
    try {
      const line = JSON.parse(event.data)
      appLogs.value.push(line)
      if (appLogs.value.length > 500) appLogs.value.shift()
      nextTick(() => {
        if (terminalRef.value) terminalRef.value.scrollTop = terminalRef.value.scrollHeight
      })
    } catch (e) {}
  }
}

const createApp = async () => {
  if (!newApp.value.app_id || !newApp.value.name) return
  try {
    const res = await apiFetch('/api/v1/apps', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newApp.value),
    })
    if (res.ok) {
      showNewAppModal.value = false
      newApp.value = { app_id: '', name: '', domain: '', ram_limit_mb: 256, cpu_quota_pct: 100 }
      fetchState()
    }
  } catch (e) {}
}

const createDB = async () => {
  if (!newDB.value.name) return
  isProvisioningDB.value = true
  try {
    const res = await apiFetch('/api/v1/databases', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newDB.value),
    })
    if (res.ok) {
      const data = await res.json()
      showNewDBModal.value = false
      selectedConn.value = data
      showConnModal.value = true
      newDB.value = { name: '', engine: 'postgres', ram_profile: 'micro' }
      fetchState()
    } else {
      const err = await res.json().catch(() => ({}))
      alert('Failed to provision database: ' + (err.error || 'Server error'))
    }
  } catch (e) {
    alert('Failed to provision database: ' + e.message)
  } finally {
    isProvisioningDB.value = false
  }
}

const loadFiles = async () => {
  if (!fileTargetID.value) return
  try {
    const res = await apiFetch(`/api/v1/jail/list?scope=${fileScope.value}&id=${fileTargetID.value}&path=`)
    if (res.ok) {
      const data = await res.json()
      fileItems.value = data.items || []
    }
  } catch (e) {}
}

const isDarkMode = ref(false)

const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value
  if (isDarkMode.value) {
    document.documentElement.setAttribute('data-theme', 'dark')
    localStorage.setItem('orbit_theme', 'dark')
  } else {
    document.documentElement.removeAttribute('data-theme')
    localStorage.setItem('orbit_theme', 'light')
  }
}

const filteredApps = computed(() => {
  if (!quickSearchQuery.value) return apps.value
  const q = quickSearchQuery.value.toLowerCase()
  return apps.value.filter(a => (a.name && a.name.toLowerCase().includes(q)) || (a.app_id && a.app_id.toLowerCase().includes(q)))
})

const filteredDBs = computed(() => {
  if (!quickSearchQuery.value) return databases.value
  const q = quickSearchQuery.value.toLowerCase()
  return databases.value.filter(d => (d.name && d.name.toLowerCase().includes(q)) || (d.engine && d.engine.toLowerCase().includes(q)))
})

const handleKeydown = (e) => {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    showQuickSearch.value = !showQuickSearch.value
    if (showQuickSearch.value) {
      nextTick(() => {
        if (quickSearchInputRef.value) quickSearchInputRef.value.focus()
      })
    }
  }
}

onMounted(() => {
  const saved = localStorage.getItem('orbit_theme')
  if (saved === 'dark') {
    isDarkMode.value = true
    document.documentElement.setAttribute('data-theme', 'dark')
  } else {
    isDarkMode.value = false
    document.documentElement.removeAttribute('data-theme')
  }
  fetchState()
  startTelemetrySSE()
  setInterval(fetchState, 5000)
  fetchUpdates()
  fetchPlatform()
  setInterval(fetchUpdates, 5 * 60 * 1000)
  setInterval(fetchPlatform, 30 * 60 * 1000)
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  if (telemetryEventSource) telemetryEventSource.close()
  if (logEventSource) logEventSource.close()
  window.removeEventListener('keydown', handleKeydown)
})
</script>
