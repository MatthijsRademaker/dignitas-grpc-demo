<script setup lang="ts">
  import { computed } from 'vue'
  import type { Auction } from '../api'
  import type { Egg } from '../useEggs'
  import Countdown from './Countdown.vue'
  import EggNest from './EggNest.vue'
  import LotArt from './LotArt.vue'
  import PriceTicker from './PriceTicker.vue'
  import SoldMoment from './SoldMoment.vue'

  const props = defineProps<{ auction: Auction, remainingMs: number, eggs: Egg[], overflow: number }>()

  const price = computed(() => props.auction.highestBid?.amount ?? props.auction.lot.startingPrice)
</script>

<template>
  <!-- The lot plate: the one framed thing on the page, like a plate in an auction catalogue. -->
  <section class="relative overflow-hidden rounded-[1.75rem] bg-surface p-2 shadow-[0_1px_2px_rgb(31_27_46/0.06),0_12px_32px_-12px_rgb(115_63_203/0.25)]">
    <div class="grid rounded-[1.35rem] border border-gold-300 sm:grid-cols-[minmax(0,5fr)_minmax(0,7fr)]">
      <div class="flex flex-col items-center justify-end rounded-t-[1.3rem] bg-[radial-gradient(circle_at_50%_40%,var(--v0-gold-50),var(--v0-surface)_70%)] px-4 pt-6 sm:rounded-l-[1.3rem] sm:rounded-tr-none">
        <div class="aspect-square w-[70%] max-w-64">
          <LotArt :lot="auction.lot" :shimmer="auction.status === 'open'" />
        </div>
        <EggNest class="-mt-2" :eggs :overflow />
      </div>

      <div class="flex flex-col gap-5 p-6 sm:py-8 sm:pr-8">
        <div>
          <p v-if="auction.round" class="mb-2 inline-block rounded-full border border-gold-300 bg-gold-50 px-3 py-0.5 text-sm font-bold text-bronze">
            Round {{ auction.round }}
          </p>
          <h2 class="font-display text-4xl leading-tight font-semibold tracking-tight">{{ auction.lot.title }}</h2>
          <p class="mt-2 max-w-prose text-muted">{{ auction.lot.description }}</p>
        </div>

        <div class="mt-auto flex items-end justify-between gap-6">
          <div class="min-w-0">
            <p class="text-sm font-bold text-muted">{{ auction.highestBid ? 'Highest bid' : 'Starting at' }}</p>
            <PriceTicker class="text-6xl font-semibold sm:text-7xl" :value="price" />
            <p class="mt-2 text-muted">
              <template v-if="auction.highestBid">by <strong class="text-on-surface">{{ auction.highestBid.bidder }}</strong></template>
              <template v-else>no bids yet</template>
            </p>
          </div>
          <Countdown
            v-if="auction.status === 'open'"
            class="w-32 shrink-0 text-sm sm:w-36"
            :lot-duration-ms="auction.lotDurationMs"
            :remaining-ms
          />
        </div>

        <p class="flex flex-wrap gap-x-5 gap-y-1 border-t border-border pt-4 text-sm text-muted">
          <span>{{ auction.bidCount }} {{ auction.bidCount === 1 ? 'bid' : 'bids' }} this round</span>
          <span>{{ auction.watchers }} live {{ auction.watchers === 1 ? 'stream' : 'streams' }} on the server</span>
        </p>
      </div>
    </div>

    <Transition name="fade">
      <SoldMoment
        v-if="auction.status !== 'open'"
        class="absolute inset-0"
        :remaining-ms
        :round="auction.round"
        :winner="auction.status === 'sold' ? auction.highestBid : null"
      />
    </Transition>
  </section>
</template>

<style scoped>
  .fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
  .fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
