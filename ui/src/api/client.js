// Centralized Orbit API client with Bearer token auth & 401 handling

export function getToken() {
  try {
    return localStorage.getItem('orbit_token') || ''
  } catch {
    return ''
  }
}

export function setToken(token) {
  try {
    if (token) localStorage.setItem('orbit_token', token)
    else localStorage.removeItem('orbit_token')
  } catch {}
}

export function authHeaders(extra = {}) {
  const token = getToken()
  return token ? { Authorization: `Bearer ${token}`, ...extra } : { ...extra }
}

export async function apiFetch(url, options = {}) {
  const headers = authHeaders(options.headers || {})
  const res = await fetch(url, { ...options, headers })
  if (res.status === 401) {
    setToken('')
  }
  return res
}

export async function apiJSON(url, options = {}) {
  const res = await apiFetch(url, options)
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(data.error || `HTTP ${res.status}`)
  }
  return data
}
