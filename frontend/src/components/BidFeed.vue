<script setup lang="ts">
  import type { Bid } from '../api'

  defineProps<{ bids: Bid[], me: string }>()

  const time = (iso: string) => new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
</script>

<template>
  <section class="rounded-2xl border border-border bg-surface p-6 shadow-sm">
    <h3 class="mb-3 text-xs font-bold tracking-widest text-muted uppercase">Latest bids</h3>
    <p v-if="bids.length === 0" class="text-muted">Be the first. Nobody likes an empty room.</p>
    <TransitionGroup v-else class="flex flex-col gap-1" name="bid" tag="ol">
      <li
        v-for="(bid, index) in bids"
        :key="bid.placedAt + bid.bidder"
        class="flex items-center justify-between rounded-lg px-3 py-2"
        :class="[index === 0 ? 'bg-primary/10 font-bold' : '', bid.bidder === me ? 'ring-2 ring-primary/40' : '']"
      >
        <span>{{ bid.bidder }}<span v-if="bid.bidder === me" class="ml-2 text-xs text-primary">(you)</span></span>
        <span class="flex items-baseline gap-3">
          <span class="font-mono text-xs text-muted">{{ time(bid.placedAt) }}</span>
          <span class="tabular-nums">€{{ bid.amount }}</span>
        </span>
      </li>
    </TransitionGroup>
  </section>
</template>

<style scoped>
  .bid-enter-active { transition: all 0.4s ease-out; }
  .bid-enter-from { opacity: 0; transform: translateY(-12px); }
  .bid-move { transition: transform 0.4s ease; }
</style>
