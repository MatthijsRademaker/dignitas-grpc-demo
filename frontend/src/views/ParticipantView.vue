<script setup lang="ts">
  import type { Auction } from '../api'
  import BidFeed from '../components/BidFeed.vue'
  import BidPanel from '../components/BidPanel.vue'
  import DashboardSkeleton from '../components/DashboardSkeleton.vue'
  import Gavel from '../components/Gavel.vue'
  import GoldenDuck from '../components/GoldenDuck.vue'
  import LotHero from '../components/LotHero.vue'
  import RoomChart from '../components/RoomChart.vue'
  import TransportPanel from '../components/TransportPanel.vue'
  import TransportToggles from '../components/TransportToggles.vue'
  import WinnersWall from '../components/WinnersWall.vue'
  import type { FeedStats, Transport } from '../useAuctionFeed'
  import type { Egg } from '../useEggs'
  import type { RoomSample } from '../useRoomHistory'

  defineProps<{
    auction: Auction | null
    remainingMs: number
    now: number
    error: string | null
    stats: FeedStats
    seenAfter: Map<string, number | null>
    eggs: Egg[]
    overflow: number
    samples: RoomSample[]
    auctionHost: string
    revealed: boolean
  }>()
  const emit = defineEmits<{ placed: [auction: Auction] }>()
  const transport = defineModel<Transport>('transport', { required: true })
  const pollIntervalMs = defineModel<number>('pollIntervalMs', { required: true })
  const bidder = defineModel<string>('bidder', { required: true })
</script>

<template>
  <div class="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-6 sm:px-6 lg:py-8">
    <header class="flex flex-wrap items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <GoldenDuck v-if="revealed" class="size-11" :shimmer="false" />
        <Gavel v-else class="size-11" />
        <div>
          <h1 class="font-display text-2xl leading-tight font-semibold tracking-tight">gRPC Auction House</h1>
          <p class="text-sm text-muted">
            browser <span class="font-mono">─HTTP/JSON→</span> your BFF
            <span class="font-mono">─gRPC→</span>
            <span class="font-mono">{{ auctionHost || '…' }}</span>
          </p>
        </div>
      </div>
      <TransportToggles v-model:poll-interval-ms="pollIntervalMs" v-model:transport="transport" />
    </header>

    <DashboardSkeleton v-if="!auction" :error="error" />

    <!-- On small screens the left column dissolves, so the bid form can sit right under the hero. -->
    <main v-else class="grid gap-8 lg:grid-cols-[minmax(0,1fr)_22rem]">
      <div class="contents lg:col-start-1 lg:row-start-1 lg:flex lg:min-w-0 lg:flex-col lg:gap-10">
        <LotHero class="order-1" :auction :eggs :overflow :remaining-ms :revealed />

        <div class="order-3 grid gap-8 md:grid-cols-2">
          <BidFeed :bids="auction.recentBids.slice(0, 8)" :me="bidder.trim()" :seen-after />
          <TransportPanel :error="error" :now :poll-interval-ms :stats="stats" :transport />
        </div>

        <div class="order-3 grid gap-8 md:grid-cols-[minmax(0,1fr)_16rem]">
          <RoomChart :load="auction.load" :now :samples />
          <WinnersWall v-if="auction.winners" :revealed :winners="auction.winners" />
        </div>
      </div>

      <!-- On wide screens the bid form stays in reach while you watch the chart. -->
      <aside class="order-2 lg:col-start-2 lg:row-start-1">
        <div class="lg:sticky lg:top-6">
          <BidPanel v-model:bidder="bidder" :auction @placed="emit('placed', $event)" />
        </div>
      </aside>
    </main>
  </div>
</template>
