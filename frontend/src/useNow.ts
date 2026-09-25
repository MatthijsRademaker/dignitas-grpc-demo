import { onScopeDispose, shallowRef } from 'vue'

export function useNow (intervalMs = 200) {
  const now = shallowRef(Date.now())
  const timer = setInterval(() => (now.value = Date.now()), intervalMs)
  onScopeDispose(() => clearInterval(timer))
  return now
}
