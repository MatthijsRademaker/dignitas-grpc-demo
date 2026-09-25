<script setup lang="ts">
  import { computed, onMounted, shallowRef, watch } from 'vue'
  import { getInfo, type Auction } from './api'
  import { DEFAULT_POLL_INTERVAL_MS, useAuctionFeed, type Transport } from './useAuctionFeed'
  import { useEggs } from './useEggs'
  import { useNow } from './useNow'
  import { usePlaceBidReady } from './usePlaceBidReady'
  import { useRoomHistory } from './useRoomHistory'
  import ParticipantView from './views/ParticipantView.vue'
  import StageView from './views/StageView.vue'

  // Two layouts, one app: `?view=stage` is the projector. `?transport=stream` works in both.
  const params = new URLSearchParams(location.search)
  const view = params.get('view') === 'stage' ? 'stage' : 'participant'
  const transport = shallowRef<Transport>(params.get('transport') === 'stream' ? 'stream' : 'poll')
  const pollIntervalMs = shallowRef<number>(DEFAULT_POLL_INTERVAL_MS)
  const bidder = shallowRef(localStorage.getItem('auction:bidder') ?? '')
  const auctionHost = shallowRef('')

  // Shared by both layouts. The room history and the eggs outlive transport switches.
  const { auction, receivedAt, error, stats, seenAfter, onNewBids, apply } = useAuctionFeed(transport, pollIntervalMs)
  const samples = useRoomHistory(auction)
  const { eggs, overflow } = useEggs(auction, onNewBids)
  // The goose is an inside joke: it stays under a cloth until this browser's BFF implements PlaceBid.
  const { revealed, markReady } = usePlaceBidReady()
  const now = useNow()
  const remainingMs = computed(() =>
    auction.value ? Math.max(0, auction.value.remainingMs - (now.value - receivedAt.value)) : 0,
  )

  watch(bidder, name => localStorage.setItem('auction:bidder', name))
  watch(revealed, value => {
    document.querySelector<HTMLLinkElement>('link[rel="icon"]')!.href = value ? '/favicon.svg' : '/favicon-gavel.svg'
  }, { immediate: true })

  function placed (next: Auction) {
    markReady()
    apply(next)
  }

  onMounted(async () => {
    if (view === 'participant') auctionHost.value = (await getInfo()).auctionHost
  })
</script>

<template>
  <StageView
    v-if="view === 'stage'"
    v-model:poll-interval-ms="pollIntervalMs"
    v-model:transport="transport"
    :auction
    :eggs
    :error
    :now
    :overflow
    :remaining-ms
    :revealed
    :samples
    :seen-after
    :stats
  />
  <ParticipantView
    v-else
    v-model:bidder="bidder"
    v-model:poll-interval-ms="pollIntervalMs"
    v-model:transport="transport"
    :auction
    :auction-host
    :eggs
    :error
    :now
    :overflow
    :remaining-ms
    :revealed
    :samples
    :seen-after
    :stats
    @placed="placed"
  />
</template>
