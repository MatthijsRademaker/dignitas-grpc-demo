<script setup lang="ts">
  import { bidKey, type Bid } from '../api'

  defineProps<{ bids: Bid[], me: string, seenAfter: Map<string, number | null> }>()

  const time = (iso: string) => new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  const delay = (ms: number) => `+${(ms / 1000).toFixed(ms < 1000 ? 2 : 1)}s`
</script>

<template>
  <section class="rounded-2xl border border-border bg-surface p-6 shadow-sm">
    <h3 class="mb-3 text-xs font-bold tracking-widest text-muted uppercase">Latest bids</h3>
    <p v-if="bids.length === 0" class="text-muted">Be the first. Nobody likes an empty room.</p>
    <TransitionGroup v-else class="flex flex-col gap-1" name="bid" tag="ol">
      <li
        v-for="(bid, index) in bids"
        :key="bidKey(bid)"
        class="flex items-center justify-between rounded-lg px-3 py-2"
        :class="[index === 0 ? 'bg-primary/10 font-bold' : '', bid.bidder === me ? 'ring-2 ring-primary/40' : '']"
      >
        <span>{{ bid.bidder }}<span v-if="bid.bidder === me" class="ml-2 text-xs text-primary">(you)</span></span>
        <span class="flex items-baseline gap-3">
          <span class="font-mono text-xs text-muted">{{ time(bid.placedAt) }}</span>
          <span
            v-if="seenAfter.get(bidKey(bid)) != null"
            class="w-14 text-right font-mono text-xs"
            :class="seenAfter.get(bidKey(bid))! < 300 ? 'text-primary' : 'text-warning'"
            title="How long after it was placed this page saw the bid"
          >{{ delay(seenAfter.get(bidKey(bid))!) }}</span>
          <span v-else class="w-14" />
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
