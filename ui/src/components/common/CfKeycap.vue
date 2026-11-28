<template>
  <kbd class="cf-keycap">
    <slot />
  </kbd>
</template>

<style scoped>
.cf-keycap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 22px;
  height: 22px;
  padding: 0 6px;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 600;
  line-height: 1;
  border-radius: var(--radius-sm, 4px);
  background-color: var(--color-kumo-recessed);
  color: var(--text-color-kumo-default);
  border: 1px solid var(--color-kumo-line);
  box-shadow: 0 1px 0 var(--color-kumo-line);
  user-select: none;
}
</style>
