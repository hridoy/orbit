<template>
  <div>
    <div class="breadcrumbs">
      <router-link to="/">Account</router-link>
      <span class="divider">/</span>
      <span class="current">SSL / TLS Encryption</span>
    </div>

    <div class="page-title-row">
      <div>
        <h1 class="page-title">SSL / TLS Encryption</h1>
        <p class="page-subtitle">Manage end-to-end edge cryptography, universal certificates, and minimum TLS protocols.</p>
      </div>
    </div>

    <!-- 4-Card Encryption Modes Grid -->
    <div class="ssl-modes-grid">
      <div class="ssl-mode-card" :class="{ active: ssl.mode === 'Off' }" @click="ssl.mode = 'Off'">
        <div class="ssl-mode-head">🔓 Off</div>
        <div class="ssl-mode-desc">No encryption applied. Visitors connect over plain HTTP.</div>
      </div>
      <div class="ssl-mode-card" :class="{ active: ssl.mode === 'Flexible' }" @click="ssl.mode = 'Flexible'">
        <div class="ssl-mode-head">🔐 Flexible</div>
        <div class="ssl-mode-desc">Encrypted between visitor and Cloudflare; unencrypted to origin.</div>
      </div>
      <div class="ssl-mode-card" :class="{ active: ssl.mode === 'Full' }" @click="ssl.mode = 'Full'">
        <div class="ssl-mode-head">🔒 Full</div>
        <div class="ssl-mode-desc">End-to-end encrypted to origin with self-signed certificate.</div>
      </div>
      <div class="ssl-mode-card" :class="{ active: ssl.mode === 'Full (Strict)' }" @click="ssl.mode = 'Full (Strict)'">
        <div class="ssl-mode-head">🛡️ Full (Strict)</div>
        <div class="ssl-mode-desc">Enforces trusted CA certificates on origin server. Highly recommended.</div>
      </div>
    </div>

    <!-- Edge Certificates Notice -->
    <div class="ui-kit-section">
      <h3 class="ui-kit-title">Universal SSL & Edge Certificates</h3>
      <p class="ui-kit-desc">Automatically provisions and renews ECC and RSA SSL certificates from Cloudflare CA and Let's Encrypt.</p>
      <div style="display: flex; gap: 1rem; align-items: center; margin-top: 1rem;">
        <span class="cf-pill active"><span class="status-dot"></span><span>Universal SSL Active</span></span>
        <span class="badge-tag-beta">TLS 1.3 Enabled</span>
        <span class="badge-tag-beta">HSTS Preload Enabled</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { mockData } from '@/data/mockData.js'

const ssl = mockData.activeZone.ssl
</script>
