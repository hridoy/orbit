<template>
  <footer class="cf-kumo-footer">
    <div class="cf-footer-inner">
      <nav class="cf-footer-nav" aria-label="Footer">
        <a href="https://support.cloudflare.com" target="_blank" rel="noopener noreferrer" class="cf-footer-link">Support</a>
        <a href="https://www.cloudflarestatus.com" target="_blank" rel="noopener noreferrer" class="cf-footer-link">System Status</a>
        <a href="https://www.cloudflare.com/careers/" target="_blank" rel="noopener noreferrer" class="cf-footer-link">Careers</a>
        <a href="https://www.cloudflare.com/terms/" target="_blank" rel="noopener noreferrer" class="cf-footer-link">Terms of Use</a>
        <a href="https://www.cloudflare.com/disclosure/" target="_blank" rel="noopener noreferrer" class="cf-footer-link">Report Security Issues</a>
        <a href="https://www.cloudflare.com/privacypolicy/" target="_blank" rel="noopener noreferrer" class="cf-footer-link">Privacy Policy</a>
        <button type="button" class="cf-footer-btn" @click="openCookiePreferences">Cookie Preferences</button>
      </nav>
      <div class="cf-footer-copyright">
        © {{ currentYear }} Cloudflare, Inc.
      </div>
    </div>
  </footer>
</template>

<script setup>
const currentYear = new Date().getFullYear()

function openCookiePreferences() {
  alert('Cloudflare Cookie Preferences dialog.')
}
</script>

<style scoped>
.cf-kumo-footer {
  width: 100%;
  height: 48px;
  min-height: 48px;
  background-color: var(--kumo-base, #ffffff);
  border-top: 1px solid var(--color-kumo-line, #e2e8f0);
  display: flex;
  align-items: center;
  box-sizing: border-box;
}

[data-theme="dark"] .cf-kumo-footer {
  background-color: #0b0f19;
  border-top-color: #1e293b;
}

.cf-footer-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0 2rem;
  font-size: 12px;
  color: #737373;
}

.cf-footer-nav {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  flex-wrap: wrap;
}

.cf-footer-link {
  color: #737373;
  text-decoration: none;
  transition: color 0.15s ease;
}

.cf-footer-link:hover {
  color: #18181b;
  text-decoration: underline;
}

[data-theme="dark"] .cf-footer-link:hover {
  color: #f8fafc;
}

.cf-footer-btn {
  background: transparent;
  border: none;
  padding: 0;
  font-size: 12px;
  color: #737373;
  cursor: pointer;
  transition: color 0.15s ease;
}

.cf-footer-btn:hover {
  color: #18181b;
  text-decoration: underline;
}

[data-theme="dark"] .cf-footer-btn:hover {
  color: #f8fafc;
}

.cf-footer-copyright {
  white-space: nowrap;
}
</style>
