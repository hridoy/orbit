<template>
  <div class="cf-create-overlay" role="dialog" aria-label="Create application">
    <!-- Top 56px White Header Bar -->
    <header class="cf-create-topbar">
      <div class="cf-create-topbar-left">
        <svg width="28" height="18" viewBox="0 0 65 30" fill="none" aria-label="Orbit">
          <path d="M44.2 14.2c-.6 0-1.1.1-1.6.3-.3-4.6-4.1-8.2-8.8-8.2-3.8 0-7.1 2.4-8.3 5.8-.7-.3-1.4-.5-2.2-.5-2.9 0-5.3 2.4-5.3 5.3 0 .2 0 .4.1.6-2.4.4-4.2 2.5-4.2 5 0 2.8 2.3 5.1 5.1 5.1h25.2c3.5 0 6.3-2.8 6.3-6.3s-2.8-6.1-6.3-6.1z" fill="#F38020"/>
        </svg>
        <div class="cf-account-pill">
          <span>Default account</span>
          <span style="color: #9ca3af; font-size: 11px;">⇅</span>
        </div>
      </div>
      <div class="cf-create-topbar-right">
        <button type="button" class="cf-topbar-close" title="Close" @click="closeWizard">✕</button>
      </div>
    </header>

    <!-- Main Centered Neutral Canvas -->
    <main class="cf-create-canvas">
      <!-- STEP: LEGACY PAGES WORKFLOW -->
      <div v-if="step === 'pages'" class="cf-pages-shell">
        <button type="button" class="cf-back-link" @click="step = 'landing'">‹ Create a Worker</button>
        <h1 class="cf-pages-h1">Get started</h1>
        <p class="cf-pages-sub">Get started with Pages. How would you like to begin?</p>
        <div class="cf-pages-cards">
          <div class="cf-pages-card">
            <div class="cf-pages-card-left">
              <div class="cf-opt-icon-box">
                <CfIcon name="git-branch" size="16" style="color: #111827;" />
              </div>
              <div>
                <div class="cf-pages-card-title">Import an existing Git repository</div>
                <div class="cf-pages-card-desc">Start by importing an existing Git repository.</div>
              </div>
            </div>
            <button type="button" class="cf-btn-outline-sm" @click="step = 'git-select'; gitExpanded = true">Get started</button>
          </div>
          <div class="cf-pages-card">
            <div class="cf-pages-card-left">
              <div class="cf-opt-icon-box">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
              </div>
              <div>
                <div class="cf-pages-card-title">Drag and drop your files</div>
                <div class="cf-pages-card-desc">Upload your site's assets including HTML, CSS, and JS files directly from your computer.</div>
              </div>
            </div>
            <button type="button" class="cf-btn-outline-sm" @click="step = 'upload'">Get started</button>
          </div>
        </div>
      </div>

      <!-- CENTERED CARD SHELL (Landing, Git, Hello World, Templates, Setup, Upload) -->
      <div v-else class="cf-create-card-wrap">
        <div class="cf-create-card">
          <!-- 1. LANDING: Make something new -->
          <template v-if="step === 'landing'">
            <div class="cf-card-body">
              <h2 class="cf-card-h2">Make something new</h2>
              <p class="cf-card-sub">Start from a template, connect a repository, or upload your code.</p>

              <div class="cf-option-list">
                <button type="button" class="cf-option-row" @click="step = 'git-select'">
                  <span class="cf-opt-icon-box">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="#111827">
                      <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
                    </svg>
                    <span class="cf-dot-green"></span>
                  </span>
                  <span class="cf-opt-label">Continue with GitHub</span>
                </button>

                <button type="button" class="cf-option-row" @click="step = 'git-select'; gitExpanded = true">
                  <span class="cf-opt-icon-box">
                    <CfIcon name="git-branch" size="15" style="color: #fc6d26;" />
                    <span class="cf-dot-gray"></span>
                  </span>
                  <span class="cf-opt-label">Connect with GitLab</span>
                </button>

                <button type="button" class="cf-option-row" @click="step = 'hello-world'">
                  <span class="cf-opt-icon-box">
                    <CfIcon name="globe" size="15" style="color: #16a34a;" />
                  </span>
                  <span class="cf-opt-label">Start with Hello World!</span>
                </button>

                <button type="button" class="cf-option-row" @click="step = 'templates'">
                  <span class="cf-opt-icon-box">
                    <CfIcon name="layers" size="15" style="color: #2563eb;" />
                  </span>
                  <span class="cf-opt-label">Select a template</span>
                </button>

                <button type="button" class="cf-option-row" @click="step = 'upload'">
                  <span class="cf-opt-icon-box">
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
                  </span>
                  <span class="cf-opt-label">Upload your static files</span>
                </button>
              </div>
            </div>
          </template>

          <!-- 2. GIT REPOSITORY SELECTOR -->
          <template v-else-if="step === 'git-select'">
            <div class="cf-card-body">
              <h2 class="cf-card-h2">Select a repository</h2>
              <p class="cf-card-sub">Select a repository from your GitHub account or clone any public Git URL to deploy.</p>

              <div class="cf-two-col">
                <select v-model="gitConnection" class="cf-input-ctrl">
                  <option value="github">GitHub (Public / Token)</option>
                  <option value="gitlab">GitLab</option>
                </select>
                <input v-model="repoSearch" type="text" class="cf-input-ctrl" placeholder="Search repositories..." />
              </div>

              <div class="cf-repo-listbox">
                <button
                  v-for="r in filteredRepos"
                  :key="r.url"
                  type="button"
                  class="cf-repo-item"
                  :class="{ selected: selectedRepoUrl === r.url }"
                  @click="selectRepo(r)"
                >
                  <span>{{ r.name }}</span>
                  <span v-if="selectedRepoUrl === r.url" style="color: #0051c3;">✓</span>
                </button>
              </div>

              <div class="cf-accordion-box">
                <button type="button" class="cf-accordion-trigger" @click="gitExpanded = !gitExpanded">
                  <span>Clone a public repository via Git URL</span>
                  <span>{{ gitExpanded ? '⌄' : '›' }}</span>
                </button>
                <div v-if="gitExpanded" class="cf-accordion-body">
                  <label class="cf-field-label">Git repository URL</label>
                  <input
                    v-model="customGitUrl"
                    type="url"
                    class="cf-input-ctrl"
                    placeholder="https://github.com/owner/repo"
                    @input="selectedRepoUrl = customGitUrl"
                  />
                  <div class="cf-field-hint">Must be a GitHub or GitLab repository URL</div>
                </div>
              </div>
            </div>

            <div class="cf-card-footer">
              <button type="button" class="cf-btn-ghost" @click="step = 'landing'">Back</button>
              <button
                type="button"
                class="cf-btn-primary"
                :disabled="!canProceedGit"
                @click="proceedFromGit"
              >
                Next →
              </button>
            </div>
          </template>

          <!-- 3. SET UP YOUR APPLICATION (Git & Template Config) -->
          <template v-else-if="step === 'setup-app'">
            <div class="cf-card-body cf-card-scroll">
              <h2 class="cf-card-h2">Set up your application</h2>
              <p class="cf-card-sub">Configure your application project and deploy it to Orbit.</p>

              <div v-if="setupForm.repoUrl" class="cf-repo-pill">
                <span>{{ setupForm.repoUrl }}</span>
                <span style="color: #6b7280;">↗</span>
              </div>

              <div class="cf-form-stack">
                <div>
                  <label class="cf-field-label">Project name</label>
                  <input v-model="setupForm.name" type="text" class="cf-input-ctrl" placeholder="my-app" />
                </div>

                <div>
                  <div class="cf-label-row">
                    <label class="cf-field-label">Build command</label>
                    <span class="cf-optional-tag">Optional</span>
                  </div>
                  <input v-model="setupForm.buildCmd" type="text" class="cf-input-ctrl cf-mono" placeholder="npm run build" />
                </div>

                <div>
                  <label class="cf-field-label">Deploy / Start command</label>
                  <input v-model="setupForm.entrypoint" type="text" class="cf-input-ctrl cf-mono" placeholder="./server" />
                </div>

                <div>
                  <label class="cf-field-label">Git branch</label>
                  <input v-model="setupForm.ref" type="text" class="cf-input-ctrl cf-mono" placeholder="main" />
                </div>

                <label class="cf-toggle-pill">
                  <span style="display: flex; align-items: center; gap: 0.5rem;">
                    <CfIcon name="git-branch" size="14" style="color: #6b7280;" />
                    <span>Enable Preview builds</span>
                  </span>
                  <input type="checkbox" v-model="setupForm.previewBuilds" />
                </label>

                <label class="cf-toggle-pill">
                  <span style="display: flex; align-items: center; gap: 0.5rem;">
                    <CfIcon name="lock" size="14" style="color: #6b7280;" />
                    <span>Protect with Orbit Access</span>
                  </span>
                  <input type="checkbox" v-model="setupForm.orbitAccess" />
                </label>

                <div class="cf-accordion-box">
                  <button type="button" class="cf-accordion-trigger" @click="advancedOpen = !advancedOpen">
                    <span>Advanced settings</span>
                    <span>{{ advancedOpen ? '⌃' : '⌄' }}</span>
                  </button>
                  <div v-if="advancedOpen" class="cf-accordion-body">
                    <div>
                      <label class="cf-field-label">Domain (optional)</label>
                      <input v-model="setupForm.domain" type="text" class="cf-input-ctrl" placeholder="app.example.com" />
                    </div>
                    <div class="cf-two-col" style="margin-top: 0.65rem;">
                      <div>
                        <label class="cf-field-label">Variable name</label>
                        <input v-model="setupForm.envKey" type="text" class="cf-input-ctrl cf-mono" placeholder="NODE_ENV" />
                      </div>
                      <div>
                        <label class="cf-field-label">Variable value</label>
                        <input v-model="setupForm.envVal" type="text" class="cf-input-ctrl cf-mono" placeholder="production" />
                      </div>
                    </div>
                  </div>
                </div>

                <CfAlert v-if="deployError" variant="danger" title="Deployment Error">{{ deployError }}</CfAlert>
              </div>
            </div>

            <div class="cf-card-footer">
              <button type="button" class="cf-btn-ghost" @click="step = setupOrigin">Back</button>
              <button type="button" class="cf-btn-primary" :disabled="deploying || !setupForm.name" @click="deployConfiguredApp">
                {{ deploying ? 'Deploying...' : 'Deploy' }}
              </button>
            </div>
          </template>

          <!-- 4. START WITH HELLO WORLD -->
          <template v-else-if="step === 'hello-world'">
            <div class="cf-card-body">
              <h2 class="cf-card-h2">Deploy Hello World</h2>
              <p class="cf-card-sub">A simple HTTP application that returns 'Hello World!'. Perfect for getting started.</p>

              <div class="cf-form-stack">
                <div>
                  <label class="cf-field-label">Worker name</label>
                  <div class="cf-composite-input">
                    <input v-model="helloName" type="text" class="cf-composite-field" />
                    <span class="cf-composite-suffix">.orbit.local</span>
                    <span style="color: #0051c3; font-size: 14px; padding-right: 0.65rem;">✓</span>
                  </div>
                </div>

                <div>
                  <label class="cf-field-label">Worker preview</label>
                  <pre class="cf-code-preview">export default {
  async fetch(request, env, ctx) {
    return new Response('Hello World!');
  },
};</pre>
                </div>

                <label class="cf-toggle-pill">
                  <span style="display: flex; align-items: center; gap: 0.5rem;">
                    <CfIcon name="lock" size="14" style="color: #6b7280;" />
                    <span>Protect with Orbit Access</span>
                  </span>
                  <input type="checkbox" v-model="helloAccess" />
                </label>

                <CfAlert v-if="deployError" variant="danger" title="Deployment Error">{{ deployError }}</CfAlert>
              </div>
            </div>

            <div class="cf-card-footer">
              <button type="button" class="cf-btn-ghost" @click="step = 'landing'">Back</button>
              <button type="button" class="cf-btn-primary" :disabled="deploying || !helloName" @click="deployHelloWorld">
                {{ deploying ? 'Deploying...' : 'Deploy' }}
              </button>
            </div>
          </template>

          <!-- 5. SELECT A TEMPLATE -->
          <template v-else-if="step === 'templates'">
            <div class="cf-card-body cf-card-scroll">
              <h2 class="cf-card-h2">Select a template</h2>
              <p class="cf-card-sub">Get up and running with one of our full-stack application examples.</p>

              <div class="cf-template-list">
                <button
                  v-for="tpl in templates"
                  :key="tpl.id"
                  type="button"
                  class="cf-template-card"
                  @click="pickTemplate(tpl)"
                >
                  <div class="cf-tpl-icon">{{ tpl.emoji }}</div>
                  <div class="cf-tpl-info">
                    <div class="cf-tpl-title-row">
                      <strong>{{ tpl.title }}</strong>
                      <span v-if="tpl.featured" class="cf-featured-pill">✦ Featured</span>
                    </div>
                    <div class="cf-tpl-desc">{{ tpl.desc }}</div>
                  </div>
                </button>
              </div>
            </div>

            <div class="cf-card-footer">
              <button type="button" class="cf-btn-ghost" @click="step = 'landing'">Back</button>
              <a href="https://workers.new/templates/" target="_blank" rel="noopener" class="cf-ext-link">Browse all templates ↗</a>
            </div>
          </template>

          <!-- 6. UPLOAD STATIC FILES -->
          <template v-else-if="step === 'upload'">
            <div class="cf-card-body">
              <h2 class="cf-card-h2">Upload and deploy</h2>
              <p class="cf-card-sub">Drag and drop your static files and configure deployment settings.</p>

              <div class="cf-form-stack">
                <div>
                  <label class="cf-field-label">Project name</label>
                  <input v-model="uploadAppName" type="text" class="cf-input-ctrl" placeholder="my-static-site" />
                </div>

                <label class="cf-dropzone">
                  <input type="file" multiple style="display: none;" @change="onStaticFilesPicked" />
                  <div class="cf-opt-icon-box" style="margin: 0 auto 0.65rem;">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
                  </div>
                  <div style="font-size: 0.875rem; color: #111827;">
                    Drag in or click to upload a <strong style="text-decoration: underline;">file</strong> or <strong style="text-decoration: underline;">folder</strong>.
                  </div>
                  <div style="font-size: 0.78rem; color: #6b7280; margin-top: 0.25rem;">
                    {{ uploadFiles.length ? `${uploadFiles.length} file(s) ready to upload` : 'Contents you drag here will be uploaded to your account' }}
                  </div>
                </label>

                <CfAlert v-if="deployError" variant="danger" title="Upload Error">{{ deployError }}</CfAlert>
              </div>
            </div>

            <div class="cf-card-footer">
              <button type="button" class="cf-btn-ghost" @click="step = 'landing'">Back</button>
              <button
                type="button"
                class="cf-btn-primary"
                :disabled="deploying || !uploadAppName || uploadFiles.length === 0"
                @click="deployStaticFiles"
              >
                {{ deploying ? 'Deploying...' : 'Deploy' }}
              </button>
            </div>
          </template>
        </div>

        <!-- Below-card legacy Pages link (shown on landing) -->
        <p v-if="step === 'landing'" class="cf-legacy-note">
          Need to use the legacy Pages workflow?
          <button type="button" class="cf-legacy-link" @click="step = 'pages'">Continue to Pages</button>
        </p>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { apiFetch } from '@/api/client.js'

const router = useRouter()
const step = ref('landing')
const setupOrigin = ref('landing')

// Git step state
const gitConnection = ref('github')
const repoSearch = ref('')
const gitExpanded = ref(false)
const selectedRepoUrl = ref('')
const customGitUrl = ref('')
const sampleRepos = [
  { name: 'orbit-starter-go', url: 'https://github.com/orbit-paas/orbit-starter-go' },
  { name: 'vite-react-app', url: 'https://github.com/vitejs/vite-react-starter' }
]
const filteredRepos = computed(() => {
  const q = repoSearch.value.trim().toLowerCase()
  if (!q) return sampleRepos
  return sampleRepos.filter(r => r.name.toLowerCase().includes(q))
})
const canProceedGit = computed(() => Boolean((selectedRepoUrl.value || customGitUrl.value).trim()))

function selectRepo(r) {
  selectedRepoUrl.value = r.url
  customGitUrl.value = r.url
}

function proceedFromGit() {
  const url = (selectedRepoUrl.value || customGitUrl.value).trim()
  const baseName = url.split('/').filter(Boolean).pop()?.replace(/\.git$/, '') || 'git-app'
  setupForm.value.repoUrl = url
  setupForm.value.name = baseName.toLowerCase().replace(/[^a-z0-9-]/g, '-')
  setupForm.value.buildCmd = 'npm run build'
  setupForm.value.entrypoint = './server'
  setupOrigin.value = 'git-select'
  step.value = 'setup-app'
}

// Setup App form
const advancedOpen = ref(false)
const deploying = ref(false)
const deployError = ref('')
const setupForm = ref({
  name: '',
  repoUrl: '',
  ref: 'main',
  buildCmd: 'npm run build',
  entrypoint: './server',
  domain: '',
  previewBuilds: true,
  orbitAccess: false,
  envKey: '',
  envVal: ''
})

// Hello World state
const helloName = ref(`billowing-art-${Math.floor(1000 + Math.random() * 9000)}`)
const helloAccess = ref(false)

// Templates state
const templates = [
  { id: 'astro', emoji: '🚀', title: 'Astro Framework Starter', desc: 'Build fast, content-focused websites and SSR apps.', featured: true, build: 'npm run build', entry: 'node ./dist/server/entry.mjs' },
  { id: 'react-router', emoji: '⚛️', title: 'React Router Framework Starter', desc: 'Full-stack React Router template with server rendering.', featured: true, build: 'npm run build', entry: 'npm run start' },
  { id: 'vite-react', emoji: '⚡', title: 'Vite React Template', desc: 'Single-page React application powered by Vite.', featured: true, build: 'npm run build', entry: './server' },
  { id: 'go-http', emoji: '🐹', title: 'Go HTTP Microservice', desc: 'Zero-dependency Go binary API server for Orbit VPS.', featured: true, build: 'go build -o server .', entry: './server' },
  { id: 'openapi', emoji: '📘', title: 'Backend OpenAPI Template', desc: 'Typed REST API service with OpenAPI documentation.', featured: false, build: 'npm run build', entry: './server' }
]

function pickTemplate(tpl) {
  setupForm.value.name = tpl.id + '-app'
  setupForm.value.repoUrl = ''
  setupForm.value.buildCmd = tpl.build
  setupForm.value.entrypoint = tpl.entry
  setupOrigin.value = 'templates'
  step.value = 'setup-app'
}

// Static upload state
const uploadAppName = ref('static-site-' + Math.floor(100 + Math.random() * 900))
const uploadFiles = ref([])

function onStaticFilesPicked(e) {
  uploadFiles.value = Array.from(e.target.files || [])
}

function slugifyId(name) {
  return name.toLowerCase().trim().replace(/[^a-z0-9-]/g, '-').replace(/^-+|-+$/g, '') || `app-${Date.now()}`
}

async function deployConfiguredApp() {
  deployError.value = ''
  deploying.value = true
  try {
    const appId = slugifyId(setupForm.value.name)
    if (setupForm.value.repoUrl) {
      const res = await apiFetch('/api/v1/apps/install-git', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          app_id: appId,
          name: setupForm.value.name,
          repo_url: setupForm.value.repoUrl,
          ref: setupForm.value.ref || 'main'
        })
      })
      if (!res.ok) throw new Error(await res.text())
    } else {
      const res = await apiFetch('/api/v1/apps', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          app_id: appId,
          name: setupForm.value.name,
          domain: setupForm.value.domain,
          binary_entrypoint: setupForm.value.entrypoint || './server'
        })
      })
      if (!res.ok) throw new Error(await res.text())
    }
    router.push(`/apps/${appId}`)
  } catch (e) {
    deployError.value = e.message || 'Deployment failed'
  } finally {
    deploying.value = false
  }
}

async function deployHelloWorld() {
  deployError.value = ''
  deploying.value = true
  try {
    const appId = slugifyId(helloName.value)
    const res = await apiFetch('/api/v1/apps', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        app_id: appId,
        name: helloName.value,
        binary_entrypoint: './server'
      })
    })
    if (!res.ok) throw new Error(await res.text())
    router.push(`/apps/${appId}`)
  } catch (e) {
    deployError.value = e.message || 'Failed to create Hello World app'
  } finally {
    deploying.value = false
  }
}

async function deployStaticFiles() {
  deployError.value = ''
  deploying.value = true
  try {
    const appId = slugifyId(uploadAppName.value)
    // Upload files into S3 assets bucket under app prefix
    for (const file of uploadFiles.value) {
      await apiFetch(`/s3/assets/${encodeURIComponent(appId)}/${encodeURIComponent(file.name)}`, {
        method: 'PUT',
        headers: { 'Content-Type': file.type || 'application/octet-stream' },
        body: file
      })
    }
    const res = await apiFetch('/api/v1/apps', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        app_id: appId,
        name: uploadAppName.value,
        binary_entrypoint: './static-serve'
      })
    })
    if (!res.ok) throw new Error(await res.text())
    router.push(`/apps/${appId}`)
  } catch (e) {
    deployError.value = e.message || 'Static upload failed'
  } finally {
    deploying.value = false
  }
}

function closeWizard() {
  router.push('/apps')
}
</script>

<style scoped>
.cf-create-overlay {
  position: fixed;
  inset: 0;
  z-index: 60;
  background: #fafafa;
  display: flex;
  flex-direction: column;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.cf-create-topbar {
  height: 56px;
  background: #ffffff;
  border-bottom: 1px solid #e5e7eb;
  padding: 0 1.25rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.cf-create-topbar-left {
  display: flex;
  align-items: center;
  gap: 0.9rem;
}

.cf-account-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.35rem 0.7rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 0.84rem;
  font-weight: 500;
  color: #111827;
  background: #ffffff;
}

.cf-topbar-close {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: #6b7280;
  font-size: 16px;
  cursor: pointer;
}

.cf-topbar-close:hover {
  background: #f3f4f6;
  color: #111827;
}

.cf-create-canvas {
  flex: 1;
  overflow-y: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem 1rem;
}

.cf-create-card-wrap {
  width: 100%;
  max-width: 560px;
}

.cf-create-card {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.cf-card-body {
  padding: 1.75rem;
}

.cf-card-scroll {
  max-height: 70vh;
  overflow-y: auto;
}

.cf-card-h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #111827;
  margin: 0 0 0.3rem;
}

.cf-card-sub {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0 0 1.25rem;
}

.cf-option-list {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.cf-option-row {
  display: flex;
  align-items: center;
  gap: 0.85rem;
  width: 100%;
  padding: 0.75rem 0.9rem;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  text-align: left;
  transition: border-color 0.15s, background 0.15s;
}

.cf-option-row:hover {
  border-color: #cbd5e1;
  background: #f9fafb;
}

.cf-opt-icon-box {
  width: 32px;
  height: 32px;
  border: 1px solid #e5e7eb;
  border-radius: 7px;
  background: #ffffff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  position: relative;
  flex-shrink: 0;
}

.cf-dot-green {
  position: absolute;
  top: -3px;
  right: -3px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #16a34a;
  border: 1.5px solid #ffffff;
}

.cf-dot-gray {
  position: absolute;
  top: -3px;
  right: -3px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #9ca3af;
  border: 1.5px solid #ffffff;
}

.cf-opt-label {
  font-size: 0.9rem;
  font-weight: 500;
  color: #111827;
}

.cf-card-footer {
  background: #f9fafb;
  border-top: 1px solid #e5e7eb;
  padding: 0.85rem 1.75rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cf-btn-ghost {
  border: none;
  background: transparent;
  color: #4b5563;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  padding: 0.45rem 0.75rem;
  border-radius: 6px;
}

.cf-btn-ghost:hover {
  background: #e5e7eb;
}

.cf-btn-primary {
  background: #0051c3;
  color: #ffffff;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
}

.cf-btn-primary:disabled {
  background: #93c5fd;
  cursor: not-allowed;
}

.cf-btn-outline-sm {
  border: 1px solid #d1d5db;
  background: #ffffff;
  color: #111827;
  border-radius: 6px;
  padding: 0.4rem 0.85rem;
  font-size: 0.84rem;
  font-weight: 500;
  cursor: pointer;
}

.cf-legacy-note {
  text-align: center;
  font-size: 0.82rem;
  color: #6b7280;
  margin-top: 1rem;
}

.cf-legacy-link {
  background: none;
  border: none;
  color: #111827;
  text-decoration: underline;
  cursor: pointer;
  font-size: 0.82rem;
}

.cf-two-col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.65rem;
  margin-bottom: 0.75rem;
}

.cf-input-ctrl {
  width: 100%;
  padding: 0.5rem 0.7rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 0.875rem;
  color: #111827;
  background: #ffffff;
  box-sizing: border-box;
}

.cf-mono {
  font-family: monospace;
}

.cf-repo-listbox {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  max-height: 150px;
  overflow-y: auto;
  margin-bottom: 0.85rem;
}

.cf-repo-item {
  width: 100%;
  padding: 0.6rem 0.85rem;
  border: none;
  border-bottom: 1px solid #f3f4f6;
  background: #ffffff;
  display: flex;
  justify-content: space-between;
  font-size: 0.875rem;
  cursor: pointer;
  text-align: left;
}

.cf-repo-item:last-child {
  border-bottom: none;
}

.cf-repo-item.selected {
  background: #eff6ff;
  color: #0051c3;
  font-weight: 500;
}

.cf-accordion-box {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

.cf-accordion-trigger {
  width: 100%;
  padding: 0.65rem 0.85rem;
  background: #f9fafb;
  border: none;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  font-weight: 500;
  color: #111827;
  cursor: pointer;
}

.cf-accordion-body {
  padding: 0.85rem;
  border-top: 1px solid #e5e7eb;
}

.cf-field-label {
  display: block;
  font-size: 0.82rem;
  font-weight: 500;
  color: #111827;
  margin-bottom: 0.3rem;
}

.cf-field-hint {
  font-size: 0.75rem;
  color: #6b7280;
  margin-top: 0.3rem;
}

.cf-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cf-optional-tag {
  font-size: 0.75rem;
  color: #6b7280;
}

.cf-form-stack {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}

.cf-repo-pill {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.55rem 0.85rem;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #f9fafb;
  font-size: 0.85rem;
  margin-bottom: 1rem;
}

.cf-toggle-pill {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.65rem 0.85rem;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 0.86rem;
  cursor: pointer;
}

.cf-composite-input {
  display: flex;
  align-items: center;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  overflow: hidden;
  background: #ffffff;
}

.cf-composite-field {
  flex: 1;
  border: none;
  padding: 0.5rem 0.7rem;
  font-size: 0.875rem;
  outline: none;
}

.cf-composite-suffix {
  font-size: 0.84rem;
  color: #6b7280;
  padding-right: 0.45rem;
}

.cf-code-preview {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 0.85rem;
  font-family: monospace;
  font-size: 0.8rem;
  color: #111827;
  margin: 0;
  overflow-x: auto;
}

.cf-template-list {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.cf-template-card {
  display: flex;
  align-items: center;
  gap: 0.85rem;
  padding: 0.75rem 0.9rem;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #ffffff;
  cursor: pointer;
  text-align: left;
}

.cf-template-card:hover {
  background: #f9fafb;
  border-color: #cbd5e1;
}

.cf-tpl-icon {
  width: 36px;
  height: 36px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}

.cf-tpl-info {
  flex: 1;
  min-width: 0;
}

.cf-tpl-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  font-size: 0.88rem;
  color: #111827;
}

.cf-featured-pill {
  background: #fef3c7;
  color: #b45309;
  font-size: 0.72rem;
  font-weight: 600;
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
  white-space: nowrap;
}

.cf-tpl-desc {
  font-size: 0.8rem;
  color: #6b7280;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 0.15rem;
}

.cf-ext-link {
  font-size: 0.84rem;
  color: #0051c3;
  text-decoration: none;
  font-weight: 500;
}

.cf-dropzone {
  display: block;
  border: 1px dashed #d1d5db;
  background: #f9fafb;
  border-radius: 10px;
  padding: 2rem 1rem;
  text-align: center;
  cursor: pointer;
}

.cf-pages-shell {
  width: 100%;
  max-width: 680px;
}

.cf-back-link {
  background: none;
  border: none;
  color: #0051c3;
  font-size: 0.85rem;
  cursor: pointer;
  padding: 0;
  margin-bottom: 0.75rem;
}

.cf-pages-h1 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #111827;
  margin: 0 0 0.25rem;
}

.cf-pages-sub {
  font-size: 0.9rem;
  color: #6b7280;
  margin: 0 0 1.25rem;
}

.cf-pages-cards {
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
}

.cf-pages-card {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 1.15rem 1.25rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.cf-pages-card-left {
  display: flex;
  align-items: center;
  gap: 0.9rem;
}

.cf-pages-card-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: #111827;
}

.cf-pages-card-desc {
  font-size: 0.82rem;
  color: #6b7280;
  margin-top: 0.15rem;
}
</style>
