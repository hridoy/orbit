<template>
  <div class="cf-env-card">
    <div class="cf-env-head">
      <div>
        <strong class="cf-env-title">Variables and Secrets</strong>
        <div class="cf-env-sub">Encrypted runtime bindings per environment</div>
      </div>
      <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
        <CfTabs v-model="envScope" :tabs="['Production', 'Previews']" />
        <CfButton variant="primary" size="sm" icon="plus" @click="$emit('add')">Add</CfButton>
      </div>
    </div>

    <div class="cf-env-rows">
      <div v-for="v in activeVars" :key="v.key" class="cf-env-row">
        <div class="cf-env-key-col">
          <CfBadge :variant="v.secret ? 'warning' : 'neutral'" size="sm">
            {{ v.secret ? 'Secret' : 'Text' }}
          </CfBadge>
          <code class="cf-env-key">{{ v.key }}</code>
        </div>
        <div class="cf-env-val-col">
          <code>{{ v.secret && !revealed[v.key] ? '••••••••••••' : v.value }}</code>
          <button
            v-if="v.secret"
            type="button"
            class="cf-env-toggle"
            @click="revealed[v.key] = !revealed[v.key]"
          >
            {{ revealed[v.key] ? 'Hide' : 'Show' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'

const props = defineProps({
  variables: {
    type: Array,
    default: () => [
      { key: 'NODE_ENV', value: 'production', secret: false, env: 'Production' },
      { key: 'JWT_SIGNING_KEY', value: 'sk_live_99aabbcc1122', secret: true, env: 'Production' },
      { key: 'PREVIEW_MODE', value: 'true', secret: false, env: 'Previews' }
    ]
  }
})

defineEmits(['add'])
const envScope = ref('Production')
const revealed = reactive({})

const activeVars = computed(() =>
  props.variables.filter(v => !v.env || v.env === envScope.value)
)
</script>

<style scoped>
.cf-env-card {
  border-radius: var(--radius-lg, 8px);
  border: 1px solid var(--color-kumo-line);
  background-color: var(--color-kumo-base);
  overflow: hidden;
}

.cf-env-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  flex-wrap: wrap;
}

.cf-env-title {
  font-size: 14px;
  color: var(--text-color-kumo-strong);
}

.cf-env-sub {
  font-size: 12px;
  color: var(--text-color-kumo-subtle);
}

.cf-env-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  border-bottom: 1px solid var(--color-kumo-hairline);
  flex-wrap: wrap;
}

.cf-env-row:last-child {
  border-bottom: none;
}

.cf-env-key-col,
.cf-env-val-col {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cf-env-key {
  font-weight: 600;
}

.cf-env-toggle {
  border: none;
  background: transparent;
  color: var(--text-color-kumo-link);
  font-size: 12px;
  cursor: pointer;
}
</style>
