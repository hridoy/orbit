// Small fetch wrapper for the agent API. Sends the saved dashboard token and turns
// error responses into thrown Errors with the agent's message.
export const getToken = () => {
  try {
    return localStorage.getItem('orbit_token') || ''
  } catch (e) {
    return ''
  }
}

export async function api(path, { method = 'GET', body } = {}) {
  const headers = {}
  const token = getToken()
  if (token) headers.Authorization = 'Bearer ' + token
  if (body !== undefined) headers['Content-Type'] = 'application/json'

  const res = await fetch(path, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })
  if (res.status === 401) {
    try { localStorage.removeItem('orbit_token') } catch (e) {}
    window.location.reload()
    throw new Error('Session expired. Sign in again.')
  }
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || `Request failed (HTTP ${res.status})`)
  return data
}
