import { onScopeDispose, reactive, shallowRef, watch, type Ref } from 'vue'
import { getAuction, type Auction, type AuctionEvent } from './api'

export type Transport = 'poll' | 'stream'

export const POLL_INTERVAL_MS = 2000

/** One dot on the transport timeline: a poll response or a stream event. */
export interface Tick {
  at: number
  changed: boolean
}

/**
 * Keeps `auction` up to date, either by polling GET /api/auction (unary gRPC behind the BFF)
 * or by listening to GET /api/auction/stream (server-streaming gRPC, translated to SSE).
 */
export function useAuctionFeed (transport: Ref<Transport>) {
  const auction = shallowRef<Auction | null>(null)
  const receivedAt = shallowRef(0)
  const error = shallowRef<string | null>(null)
  const stats = reactive({ requests: 0, messages: 0, lastKind: '' as string, ticks: [] as Tick[] })

  let stop = () => {}

  function apply (next: Auction) {
    const previous = auction.value
    const changed = !previous
      || previous.lot.id !== next.lot.id
      || previous.status !== next.status
      || previous.bidCount !== next.bidCount
      || previous.watchers !== next.watchers
    auction.value = next
    receivedAt.value = Date.now()
    return changed
  }

  function record (changed: boolean) {
    const now = Date.now()
    stats.ticks = [...stats.ticks.filter(t => now - t.at < 30_000), { at: now, changed }]
  }

  function startPolling () {
    let cancelled = false
    async function poll () {
      stats.requests++
      try {
        const result = await getAuction()
        if (cancelled) return
        if (result.ok) {
          error.value = null
          record(apply(result.value))
        } else {
          error.value = result.problem.detail ?? result.problem.title
        }
      } catch {
        if (!cancelled) error.value = 'Cannot reach the BFF'
      }
    }
    poll()
    const timer = setInterval(poll, POLL_INTERVAL_MS)
    return () => {
      cancelled = true
      clearInterval(timer)
    }
  }

  function startStreaming () {
    stats.requests++
    const source = new EventSource('/api/auction/stream')
    source.addEventListener('auction', message => {
      const event: AuctionEvent = JSON.parse(message.data)
      error.value = null
      stats.messages++
      stats.lastKind = event.kind
      apply(event.auction)
      record(true)
    })
    source.addEventListener('open', () => (error.value = null))
    // EventSource reconnects by itself; every reconnect is a new request.
    source.addEventListener('error', () => {
      error.value = 'Stream interrupted, reconnecting…'
      stats.requests++
    })
    return () => source.close()
  }

  watch(transport, mode => {
    stop()
    stats.requests = 0
    stats.messages = 0
    stats.lastKind = ''
    stats.ticks = []
    stop = mode === 'poll' ? startPolling() : startStreaming()
  }, { immediate: true })

  onScopeDispose(() => stop())

  return { auction, receivedAt, error, stats, apply }
}
