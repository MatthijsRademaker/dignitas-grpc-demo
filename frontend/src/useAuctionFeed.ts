import { onScopeDispose, reactive, shallowRef, watch, type Ref } from 'vue'
import { bidKey, getAuction, roundKey, type Auction, type AuctionEvent, type Bid } from './api'

export type Transport = 'poll' | 'stream'

export const POLL_INTERVALS_MS = [500, 2000, 5000] as const
export const DEFAULT_POLL_INTERVAL_MS = 2000

/** One dot on the transport timeline: a poll response or a stream event. */
export interface Tick {
  at: number
  changed: boolean
}

/** What this browser sent and received since it last changed transport or poll interval. */
export interface FeedStats {
  requests: number
  messages: number
  lastKind: string
  ticks: Tick[]
  delays: number[]
}

/**
 * Called with the bids this page sees for the first time, once per response or event, newest
 * first. `initial`: they were already there when the page loaded, rather than new.
 */
export type NewBidsListener = (bids: Bid[], initial: boolean) => void

/**
 * Keeps `auction` up to date, either by polling GET /api/auction (unary gRPC behind the BFF)
 * or by listening to GET /api/auction/stream (server-streaming gRPC, translated to SSE).
 *
 * It also measures how late every bid arrived here: the server stamps each bid with its age
 * when it sends it, so no clock has to agree with another.
 */
export function useAuctionFeed (transport: Ref<Transport>, pollIntervalMs: Ref<number>) {
  const auction = shallowRef<Auction | null>(null)
  const receivedAt = shallowRef(0)
  const error = shallowRef<string | null>(null)
  const stats = reactive<FeedStats>({ requests: 0, messages: 0, lastKind: '', ticks: [], delays: [] })
  /** Delay per bid, the first time this page saw it. `null`: it was already there when the page loaded. */
  const seenAfter = reactive(new Map<string, number | null>())

  const newBidsListeners: NewBidsListener[] = []

  let stop = () => {}

  function apply (next: Auction, fromFeed = true) {
    const previous = auction.value
    const changed = !previous
      || roundKey(previous) !== roundKey(next)
      || previous.status !== next.status
      || previous.bidCount !== next.bidCount
      || previous.watchers !== next.watchers
    if (previous && roundKey(previous) !== roundKey(next)) seenAfter.clear()
    const newBids: Bid[] = []
    for (const bid of next.recentBids) {
      if (seenAfter.has(bidKey(bid))) continue
      seenAfter.set(bidKey(bid), previous ? bid.ageMs : null)
      newBids.push(bid)
      // Your own bids come back in the PlaceBid response: only count what the feed delivered.
      if (previous && fromFeed) stats.delays = [...stats.delays.slice(-19), bid.ageMs]
    }
    auction.value = next
    receivedAt.value = Date.now()
    // After `auction` is set, so a listener that watches it already sees the new round.
    if (newBids.length > 0) {
      for (const listener of newBidsListeners) listener(newBids, !previous)
    }
    return changed
  }

  function record (changed: boolean) {
    const now = Date.now()
    stats.ticks = [...stats.ticks.filter(t => now - t.at < 30_000), { at: now, changed }]
  }

  function startPolling (intervalMs: number) {
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
    const timer = setInterval(poll, intervalMs)
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

  watch([transport, pollIntervalMs], ([mode, intervalMs]) => {
    stop()
    stats.requests = 0
    stats.messages = 0
    stats.lastKind = ''
    stats.ticks = []
    stats.delays = []
    stop = mode === 'poll' ? startPolling(intervalMs) : startStreaming()
  }, { immediate: true })

  onScopeDispose(() => stop())

  function onNewBids (listener: NewBidsListener) {
    newBidsListeners.push(listener)
  }

  return { auction, receivedAt, error, stats, seenAfter, onNewBids, apply: (next: Auction) => apply(next, false) }
}
