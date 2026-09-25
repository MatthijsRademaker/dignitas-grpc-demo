import { onScopeDispose, shallowRef, watch } from 'vue'

const STORAGE_KEY = 'auction:revealed'
const PROBE_EVERY_MS = 5000

/**
 * Whether this browser's BFF implements PlaceBid: the Golden Goose stays under a cloth until it does.
 *
 * The probe can never place a bid: an empty bidder. The workshop stub answers 501. A working BFF passes
 * the call on, the auction server rejects it with INVALID_ARGUMENT, and ToProblem() makes that a 400
 * titled "InvalidArgument": proof the call reached the server. While veiled it keeps asking, because
 * `--watch` restarts the BFF on every save, so the reveal needs no reload.
 */
export function usePlaceBidReady () {
  // Remembered, so a participant who already finished doesn't see the cloth flash on every reload.
  const revealed = shallowRef(localStorage.getItem(STORAGE_KEY) === 'true')
  watch(revealed, value => localStorage.setItem(STORAGE_KEY, String(value)))

  async function probe () {
    try {
      const response = await fetch('/api/bids', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ lotId: '', bidder: '', amount: 0 }),
      })
      if (response.status === 501) {
        revealed.value = false
      } else if (response.status === 400) {
        const problem = await response.json().catch(() => null)
        if (problem?.title === 'InvalidArgument') revealed.value = true
      }
      // Anything else (BFF restarting, server unreachable) says nothing either way: ask again later.
    } catch {
      // BFF not reachable right now: ask again later.
    }
  }

  // Once on load, which also corrects a stale remembered value, then every few seconds while veiled.
  probe()
  const timer = setInterval(() => {
    if (!revealed.value) probe()
  }, PROBE_EVERY_MS)
  onScopeDispose(() => clearInterval(timer))

  /** A bid went through: PlaceBid clearly works, even if its error handling doesn't yet. */
  function markReady () {
    revealed.value = true
  }

  return { revealed, markReady }
}
