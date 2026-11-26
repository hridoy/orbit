<template>
  <div class="home-page">
    <div class="home-hero">
      <div class="home-pill">
        <span class="home-dot"></span>
        <span>{{ server.hostname || 'This server' }}<template v-if="server.os"> · {{ server.os }} {{ server.arch }}</template></span>
      </div>
      <h1 class="home-title">VPS Controller</h1>
      <p class="home-sub">Apps, databases and domains running on this server.</p>
    </div>

    <div class="home-grid">
      <section class="home-card">
        <div class="home-card-head">
          <router-link to="/ssl" class="home-link">Domains</router-link>
          <router-link to="/settings" class="home-link-muted">Settings</router-link>
        </div>
        <p v-if="!server.domains.length && !server.panelDomain" class="home-empty">No domains yet. Add one in SSL/TLS.</p>
        <ul class="home-list">
          <li v-if="server.panelDomain">
            <router-link to="/settings">{{ server.panelDomain }}</router-link>
            <span class="home-tag">panel</span>
          </li>
          <li v-for="d in server.domains" :key="d.domain">
            <router-link to="/ssl">{{ d.domain }}</router-link>
            <span :class="['home-tag', d.status]">{{ d.status }}</span>
          </li>
        </ul>
      </section>

      <section class="home-card">
        <div class="home-card-head">
          <router-link to="/workers" class="home-link">Apps</router-link>
          <span class="home-count">{{ server.apps.length }}</span>
        </div>
        <p v-if="!server.apps.length" class="home-empty">No apps yet.</p>
        <ul class="home-list">
          <li v-for="a in server.apps" :key="a.app_id">
            <router-link to="/workers">{{ a.name || a.app_id }}</router-link>
            <span :class="['home-tag', a.status]">{{ a.status }}</span>
          </li>
        </ul>
      </section>

      <section class="home-card">
        <div class="home-card-head">
          <router-link to="/d1" class="home-link">Databases</router-link>
          <span class="home-count">{{ server.databases.length }}</span>
        </div>
        <p v-if="!server.databases.length" class="home-empty">No databases yet.</p>
        <ul class="home-list">
          <li v-for="d in server.databases" :key="d.db_id">
            <router-link :to="`/d1/${d.db_id}`">{{ d.name }}</router-link>
            <span :class="['home-tag', d.status]">{{ d.engine }}</span>
          </li>
        </ul>
      </section>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { server, loadServer } from '@/composables/serverState.js'

let timer
onMounted(() => {
  loadServer()
  timer = setInterval(loadServer, 10000)
})
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.home-page { display: flex; flex-direction: column; gap: 1.25rem; min-width: 0; }
.home-hero { text-align: center; padding: 1rem 0; }
.home-pill { display: inline-flex; align-items: center; gap: 0.5rem; font-size: 13px; padding: 4px 12px; border-radius: 999px; background: #f3f4f6; }
.home-dot { width: 8px; height: 8px; border-radius: 50%; background: #16a34a; }
.home-title { font-size: 28px; margin: 0.75rem 0 0.25rem; }
.home-sub { font-size: 14px; color: #6b7280; margin: 0; }
.home-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 1rem; }
.home-card { background: var(--kumo-base, #fff); border: 1px solid var(--kumo-border, #e2e4e8); border-radius: 10px; padding: 1rem; min-width: 0; box-sizing: border-box; }
.home-card-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.home-link { font-weight: 600; color: inherit; text-decoration: none; }
.home-link-muted { font-size: 13px; color: #6b7280; text-decoration: none; }
.home-count { font-size: 13px; color: #6b7280; }
.home-empty { font-size: 13px; color: #6b7280; margin: 0; }
.home-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 0.5rem; }
.home-list li { display: flex; justify-content: space-between; gap: 0.5rem; align-items: center; font-size: 14px; min-width: 0; }
.home-list a { overflow-wrap: anywhere; min-width: 0; }
.home-tag { font-size: 11px; padding: 2px 8px; border-radius: 999px; background: #e5e7eb; color: #374151; white-space: nowrap; }
.home-tag.running, .home-tag.active { background: #dcfce7; color: #166534; }
.home-tag.stopped, .home-tag.failed { background: #fee2e2; color: #991b1b; }
@media (max-width: 900px) {
  .home-grid { grid-template-columns: minmax(0, 1fr); }
}
</style>
