import { ref, reactive } from 'vue'

const toasts = ref([])
let idCounter = 0

export function useToast() {
  function show({ title = '', message = '', type = 'info', duration = 4000 }) {
    const id = ++idCounter
    const toast = { id, title, message, type, duration }
    toasts.value.push(toast)

    if (duration > 0) {
      setTimeout(() => {
        remove(id)
      }, duration)
    }
    return id
  }

  function remove(id) {
    const idx = toasts.value.findIndex(t => t.id === id)
    if (idx !== -1) {
      toasts.value.splice(idx, 1)
    }
  }

  function success(message, title = 'Success') {
    return show({ title, message, type: 'success' })
  }

  function error(message, title = 'Error') {
    return show({ title, message, type: 'error' })
  }

  function info(message, title = 'Information') {
    return show({ title, message, type: 'info' })
  }

  function warning(message, title = 'Warning') {
    return show({ title, message, type: 'warning' })
  }

  return {
    toasts,
    show,
    remove,
    success,
    error,
    info,
    warning
  }
}
