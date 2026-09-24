<script setup lang="ts">
  import { computed } from 'vue'
  import type { Auction } from '../api'

  const props = defineProps<{ auction: Auction, remainingMs: number }>()

  const price = computed(() => props.auction.highestBid?.amount ?? props.auction.lot.startingPrice)
  const seconds = computed(() => Math.ceil(props.remainingMs / 1000))
  const countdown = computed(() => `${Math.floor(seconds.value / 60)}:${String(seconds.value % 60).padStart(2, '0')}`)
  const closing = computed(() => props.auction.status === 'open' && seconds.value <= 10)
</script>

<template>
  <section class="overflow-hidden rounded-2xl border border-border bg-surface shadow-sm">
    <div class="flex items-start gap-5 p-6">
      <div class="grid size-24 shrink-0 place-items-center rounded-2xl bg-background text-6xl">
        {{ auction.lot.emoji }}
      </div>
      <div class="min-w-0">
        <p class="text-xs font-bold tracking-widest text-primary uppercase">Now on the block</p>
        <h2 class="text-2xl font-extrabold">{{ auction.lot.title }}</h2>
        <p class="mt-1 text-muted">{{ auction.lot.description }}</p>
      </div>
    </div>

    <div class="grid grid-cols-2 border-t border-border">
      <div class="p-6">
        <p class="text-xs font-bold tracking-widest text-muted uppercase">
          {{ auction.highestBid ? 'Highest bid' : 'Starting at' }}
        </p>
        <p :key="price" class="animate-flash -mx-2 inline-block rounded-lg px-2 text-5xl font-black tabular-nums">
          €{{ price }}
        </p>
        <p class="mt-1 h-6 text-sm text-muted">
          <template v-if="auction.highestBid">by <strong class="text-on-surface">{{ auction.highestBid.bidder }}</strong></template>
          <template v-else>no bids yet</template>
        </p>
      </div>

      <div class="border-l border-border p-6">
        <template v-if="auction.status === 'open'">
          <p class="text-xs font-bold tracking-widest text-muted uppercase">
            {{ closing ? 'Going once, going twice…' : 'Closes in' }}
          </p>
          <p class="text-5xl font-black tabular-nums" :class="closing ? 'text-error' : ''">{{ countdown }}</p>
          <p class="mt-1 text-sm text-muted">late bids extend the clock</p>
        </template>
        <template v-else>
          <p class="text-xs font-bold tracking-widest uppercase" :class="auction.status === 'sold' ? 'text-success' : 'text-muted'">
            {{ auction.status === 'sold' ? 'Sold! 🎉' : 'No bids, unsold' }}
          </p>
          <p class="text-2xl font-extrabold">
            <template v-if="auction.highestBid">{{ auction.highestBid.bidder }} · €{{ auction.highestBid.amount }}</template>
            <template v-else>Better luck next lot</template>
          </p>
          <p class="mt-1 text-sm text-muted">next lot in {{ seconds }}s</p>
        </template>
      </div>
    </div>

    <div class="flex justify-between border-t border-border bg-background/60 px-6 py-3 text-sm text-muted">
      <span>{{ auction.bidCount }} {{ auction.bidCount === 1 ? 'bid' : 'bids' }} on this lot</span>
      <span>👀 {{ auction.watchers }} live {{ auction.watchers === 1 ? 'stream' : 'streams' }} on the server</span>
    </div>
  </section>
</template>
