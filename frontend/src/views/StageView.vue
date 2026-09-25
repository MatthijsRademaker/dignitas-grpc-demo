<script setup lang="ts">
  import { computed, onScopeDispose, shallowRef } from 'vue'
  import type { Auction } from '../api'
  import BidFeed from '../components/BidFeed.vue'
  import Countdown from '../components/Countdown.vue'
  import EggNest from '../components/EggNest.vue'
  import Gavel from '../components/Gavel.vue'
  import GoldenDuck from '../components/GoldenDuck.vue'
  import LotArt from '../components/LotArt.vue'
  import PriceTicker from '../components/PriceTicker.vue'
  import RoomChart from '../components/RoomChart.vue'
  import SoldMoment from '../components/SoldMoment.vue'
  import TransportPanel from '../components/TransportPanel.vue'
  import TransportToggles from '../components/TransportToggles.vue'
  import WinnersWall from '../components/WinnersWall.vue'
  import type { FeedStats, Transport } from '../useAuctionFeed'
  import type { Egg } from '../useEggs'
  import type { RoomSample } from '../useRoomHistory'

  const props = defineProps<{
    auction: Auction | null
    remainingMs: number
    now: number
    error: string | null
    stats: FeedStats
    seenAfter: Map<string, number | null>
    eggs: Egg[]
    overflow: number
    samples: RoomSample[]
    revealed: boolean
  }>()
  const transport = defineModel<Transport>('transport', { required: true })
  const pollIntervalMs = defineModel<number>('pollIntervalMs', { required: true })

  // Laid out once, at 1920×1080, and scaled to whatever the projector is: a 1280×720 beamer gets
  // exactly the same layout, smaller, and never a scrollbar.
  const WIDTH = 1920
  const HEIGHT = 1080
  const fit = () => Math.min(innerWidth / WIDTH, innerHeight / HEIGHT)
  const scale = shallowRef(fit())
  const onResize = () => (scale.value = fit())
  addEventListener('resize', onResize)
  onScopeDispose(() => removeEventListener('resize', onResize))

  const price = computed(() => props.auction ? props.auction.highestBid?.amount ?? props.auction.lot.startingPrice : 0)
</script>

<template>
  <div class="fixed inset-0 overflow-hidden bg-background">
    <div
      class="absolute top-1/2 left-1/2 grid h-[1080px] w-[1920px] grid-cols-1 grid-rows-[56px_470px_minmax(0,1fr)_56px] gap-7 px-12 py-10"
      :style="{ transform: `translate(-50%, -50%) scale(${scale})` }"
    >
      <header class="flex items-center justify-between">
        <div class="flex items-center gap-4">
          <GoldenDuck v-if="revealed" class="size-14" :shimmer="false" />
          <Gavel v-else class="size-14" />
          <h1 class="font-display text-[2.4rem] font-semibold tracking-tight">gRPC Auction House</h1>
        </div>
        <TransportToggles v-model:poll-interval-ms="pollIntervalMs" v-model:transport="transport" quiet />
      </header>

      <template v-if="auction">
        <section class="rounded-[2.5rem] bg-surface p-3 shadow-[0_2px_4px_rgb(31_27_46/0.05),0_24px_48px_-20px_rgb(115_63_203/0.3)]">
          <div class="grid h-full grid-cols-[520px_minmax(0,1fr)_400px] items-center gap-10 rounded-[2rem] border-2 border-gold-300 pr-12">
            <div class="flex h-full flex-col items-center justify-end rounded-l-[1.9rem] bg-[radial-gradient(circle_at_50%_40%,var(--v0-gold-50),var(--v0-surface)_70%)] px-8 pt-6">
              <div class="size-[290px]">
                <LotArt :lot="auction.lot" :shimmer="auction.status === 'open'" :veiled="!revealed" />
              </div>
              <Transition mode="out-in" name="nest-in">
                <EggNest v-if="revealed" class="-mt-3" :eggs :overflow />
                <p v-else class="pt-2 pb-6 text-center text-[1.5rem] text-muted">Implement PlaceBid to find out what's under the cloth.</p>
              </Transition>
            </div>

            <div class="flex min-w-0 flex-col">
              <p v-if="auction.round" class="mb-3 self-start rounded-full border-2 border-gold-300 bg-gold-50 px-5 py-1 text-[1.75rem] font-bold text-bronze">
                Round {{ auction.round }}
              </p>
              <h2 class="font-display text-[4rem] leading-none font-semibold tracking-tight">{{ revealed ? auction.lot.title : 'A mystery lot' }}</h2>
              <p class="mt-5 text-[1.75rem] font-bold text-muted">{{ auction.highestBid ? 'Highest bid' : 'Starting at' }}</p>
              <PriceTicker class="-my-1 text-[9.5rem] font-semibold" :value="price" />
              <p class="mt-1 truncate text-[2.25rem] text-muted">
                <template v-if="auction.highestBid">by <strong class="text-on-surface">{{ auction.highestBid.bidder }}</strong></template>
                <template v-else>no bids yet</template>
              </p>
            </div>

            <Countdown
              class="w-[380px] text-[1.75rem]"
              :class="auction.status !== 'open' && 'invisible'"
              :lot-duration-ms="auction.lotDurationMs"
              :remaining-ms
            />
          </div>
        </section>

        <div class="grid min-h-0 grid-cols-[minmax(0,1fr)_520px_330px] gap-10">
          <RoomChart :load="auction.load" :now :samples size="xl" />
          <BidFeed :bids="auction.recentBids.slice(0, 3)" :seen-after size="xl" />
          <TransportPanel :error :now :poll-interval-ms :stats :transport variant="compact" />
        </div>

        <WinnersWall v-if="auction.winners" layout="strip" :revealed :winners="auction.winners" />
      </template>

      <p v-else class="row-span-3 grid place-items-center text-[2.5rem]" :class="error ? 'font-bold text-error' : 'text-muted'">
        {{ error ?? 'Connecting to the auction…' }}
      </p>

      <Transition name="fade">
        <SoldMoment
          v-if="auction && auction.status !== 'open'"
          class="absolute inset-0 z-20"
          :remaining-ms
          :revealed
          :round="auction.round"
          size="fullscreen"
          :winner="auction.status === 'sold' ? auction.highestBid : null"
        />
      </Transition>
    </div>
  </div>
</template>

<style scoped>
  .fade-enter-active, .fade-leave-active { transition: opacity 0.4s ease; }
  .fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
