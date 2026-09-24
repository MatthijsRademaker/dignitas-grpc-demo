import { shallowRef, watch, type Ref } from 'vue'
import type { Auction } from './api'

export const ROOM_HISTORY_MS = 90_000

/** The room load as the server reported it, at the moment this page received it. */
export interface RoomSample {
  at: number
  pollsPerSecond: number
  watchers: number
}

/**
 * The last 90 seconds of room load, sampled every time an auction arrives, whichever way it came.
 * Unlike the per-browser stats in useAuctionFeed, it never resets when this page switches
 * transport or poll interval: the whole point is to watch the room flip.
 */
export function useRoomHistory (auction: Ref<Auction | null>) {
  const samples = shallowRef<RoomSample[]>([])

  watch(auction, next => {
    if (!next) return
    const at = Date.now()
    const cutoff = at - ROOM_HISTORY_MS
    const kept = samples.value
    // Keep the last sample from before the window too: its value is still in force at the left edge.
    let first = 0
    while (first + 1 < kept.length && kept[first + 1]!.at <= cutoff) first++
    samples.value = [...kept.slice(first), { at, pollsPerSecond: next.load.pollsPerSecond, watchers: next.watchers }]
  })

  return samples
}
