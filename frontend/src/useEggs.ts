import { computed, shallowRef, watch, type Ref } from 'vue'
import { bidKey, roundKey, type Auction } from './api'
import type { NewBidsListener } from './useAuctionFeed'

export const NEST_CAPACITY = 24

export interface Egg {
  key: string
  /** Fell in while this page was watching. `false`: the bid was already there at page load. */
  dropped: boolean
}

/**
 * One golden egg per bid this page sees for the first time. Every call of the listener is one
 * batch, and a batch mounts at once: a poll that brings three bids drops three eggs together,
 * a stream drops them one at a time. The nest empties when a new round opens.
 */
export function useEggs (auction: Ref<Auction | null>, onNewBids: (listener: NewBidsListener) => void) {
  const eggs = shallowRef<Egg[]>([])

  // Sync: useAuctionFeed sets `auction` before it reports new bids, so a new round's first
  // bids land in an empty nest instead of being cleared right after they fall.
  watch(() => auction.value && roundKey(auction.value), (round, previous) => {
    if (previous !== undefined && round !== previous) eggs.value = []
  }, { flush: 'sync' })

  onNewBids((bids, initial) => {
    // Oldest first, so the newest egg lands last.
    const laid = [...bids].reverse().map(bid => ({ key: bidKey(bid), dropped: !initial }))
    if (initial) {
      // Only the latest bids travel with the auction: fill in the rest of the count.
      const missing = Math.max(0, (auction.value?.bidCount ?? 0) - laid.length)
      laid.unshift(...Array.from({ length: missing }, (_, i) => ({ key: `earlier-${i}`, dropped: false })))
    }
    eggs.value = [...eggs.value, ...laid]
  })

  const visible = computed(() => eggs.value.slice(-NEST_CAPACITY))
  const overflow = computed(() => Math.max(0, eggs.value.length - NEST_CAPACITY))

  return { eggs: visible, overflow }
}
