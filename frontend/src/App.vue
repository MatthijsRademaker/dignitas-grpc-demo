<script setup lang="ts">
  import { Toggle } from '@vuetify/v0'
  import { computed, onMounted, shallowRef, watch } from 'vue'
  import { getInfo } from './api'
  import BidFeed from './components/BidFeed.vue'
  import BidPanel from './components/BidPanel.vue'
  import LotCard from './components/LotCard.vue'
  import RoomLoad from './components/RoomLoad.vue'
  import TransportPanel from './components/TransportPanel.vue'
  import { DEFAULT_POLL_INTERVAL_MS, POLL_INTERVALS_MS, useAuctionFeed, type Transport } from './useAuctionFeed'
  import { useNow } from './useNow'

  const initial = new URLSearchParams(location.search).get('transport')
  const transport = shallowRef<Transport>(initial === 'stream' ? 'stream' : 'poll')
  const pollIntervalMs = shallowRef<number>(DEFAULT_POLL_INTERVAL_MS)
  const bidder = shallowRef(localStorage.getItem('auction:bidder') ?? '')
  const auctionHost = shallowRef('')

  const { auction, receivedAt, error, stats, seenAfter, apply } = useAuctionFeed(transport, pollIntervalMs)
  const now = useNow()
  const remainingMs = computed(() =>
    auction.value ? Math.max(0, auction.value.remainingMs - (now.value - receivedAt.value)) : 0,
  )

  watch(bidder, name => localStorage.setItem('auction:bidder', name))
  onMounted(async () => (auctionHost.value = (await getInfo()).auctionHost))
</script>

<template>
  <div class="mx-auto flex max-w-6xl flex-col gap-6 p-6">
    <header class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-extrabold tracking-tight">🔨 gRPC Auction House</h1>
        <p class="text-sm text-muted">
          browser <span class="font-mono">─HTTP/JSON→</span> your BFF
          <span class="font-mono">─gRPC→</span>
          <span class="font-mono">{{ auctionHost || '…' }}</span>
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <Toggle.Group
          v-if="transport === 'poll'"
          v-model="pollIntervalMs"
          class="flex rounded-xl border border-border bg-surface p-1 shadow-sm"
          label="How often to poll"
          mandatory
        >
          <Toggle.Root
            v-for="ms in POLL_INTERVALS_MS"
            :key="ms"
            class="cursor-pointer rounded-lg px-3 py-2 text-sm font-bold text-muted tabular-nums transition data-[state=on]:bg-on-surface data-[state=on]:text-surface"
            :value="ms"
          >
            {{ ms / 1000 }}s
          </Toggle.Root>
        </Toggle.Group>

        <Toggle.Group
          v-model="transport"
          class="flex rounded-xl border border-border bg-surface p-1 shadow-sm"
          label="How this page gets updates"
          mandatory
        >
          <Toggle.Root
            class="cursor-pointer rounded-lg px-4 py-2 text-sm font-bold text-muted transition data-[state=on]:bg-primary data-[state=on]:text-on-primary"
            value="poll"
          >
            Polling · every {{ pollIntervalMs / 1000 }}s
          </Toggle.Root>
          <Toggle.Root
            class="cursor-pointer rounded-lg px-4 py-2 text-sm font-bold text-muted transition data-[state=on]:bg-primary data-[state=on]:text-on-primary"
            value="stream"
          >
            Streaming · SSE ← gRPC
          </Toggle.Root>
        </Toggle.Group>
      </div>
    </header>

    <div v-if="!auction" class="rounded-2xl border border-border bg-surface p-10 text-center text-muted">
      {{ error ?? 'Connecting to the auction…' }}
    </div>

    <main v-else class="grid gap-6 lg:grid-cols-[1.5fr_1fr]">
      <div class="flex flex-col gap-6">
        <LotCard :auction :remaining-ms />
        <RoomLoad :load="auction.load" :watchers="auction.watchers" />
        <BidFeed :bids="auction.recentBids.slice(0, 8)" :me="bidder" :seen-after />
      </div>
      <div class="flex flex-col gap-6">
        <BidPanel v-model:bidder="bidder" :auction @placed="apply" />
        <TransportPanel :error :now :poll-interval-ms :stats :transport />
      </div>
    </main>
  </div>
</template>
