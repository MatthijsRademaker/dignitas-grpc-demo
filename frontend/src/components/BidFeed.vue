<script setup lang="ts">
  import { avatarColour, initials } from '../avatar'
  import { bidKey, type Bid } from '../api'

  withDefaults(defineProps<{
    bids: Bid[]
    me?: string
    seenAfter: Map<string, number | null>
    size?: 'md' | 'xl'
  }>(), { me: '', size: 'md' })

  const time = (iso: string) => new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  const delay = (ms: number) => `+${(ms / 1000).toFixed(ms < 1000 ? 2 : 1)}s`
</script>

<template>
  <section class="feed" :class="size">
    <h3 class="heading font-bold">Latest bids</h3>
    <p v-if="bids.length === 0" class="text-muted">Be the first. Nobody likes an empty room.</p>
    <TransitionGroup v-else class="flex flex-col" name="bid" tag="ol">
      <li
        v-for="(bid, index) in bids"
        :key="bidKey(bid)"
        class="row flex items-center gap-3 rounded-xl"
        :class="[index === 0 && 'bg-primary/8', me && bid.bidder === me && 'ring-2 ring-primary/40 ring-inset']"
      >
        <span
          aria-hidden="true"
          class="avatar grid shrink-0 place-items-center rounded-full font-extrabold text-white"
          :style="{ background: avatarColour(bid.bidder) }"
        >{{ initials(bid.bidder) }}</span>
        <span class="min-w-0 flex-1">
          <span class="block truncate font-bold">
            {{ bid.bidder }}<span v-if="me && bid.bidder === me" class="ml-2 text-sm font-bold text-primary">(you)</span>
          </span>
          <span class="meta flex gap-3 text-muted tabular-nums">
            <span>{{ time(bid.placedAt) }}</span>
            <span
              v-if="seenAfter.get(bidKey(bid)) != null"
              :class="seenAfter.get(bidKey(bid))! < 300 ? 'text-primary' : 'text-warning'"
              title="How long after it was placed this page saw the bid"
            >seen {{ delay(seenAfter.get(bidKey(bid))!) }} later</span>
          </span>
        </span>
        <span class="amount font-display font-semibold tabular-nums">€{{ bid.amount }}</span>
      </li>
    </TransitionGroup>
  </section>
</template>

<style scoped>
  .md .heading { margin-bottom: 0.75rem; }
  .md .row { padding: 0.5rem 0.75rem; }
  .md .avatar { width: 2.25rem; height: 2.25rem; font-size: 0.8rem; }
  .md .meta { font-size: 0.8rem; }
  .md .amount { font-size: 1.35rem; }

  .xl { font-size: 1.75rem; }
  .xl .heading { margin-bottom: 0.5rem; font-size: 1.6rem; }
  .xl .row { padding: 0.6rem 1rem; }
  .xl .avatar { width: 3.75rem; height: 3.75rem; font-size: 1.35rem; }
  .xl .meta { font-size: 1.1rem; }
  .xl .amount { font-size: 2.6rem; }

  .bid-enter-active { transition: all 0.4s ease-out; }
  .bid-enter-from { opacity: 0; transform: translateY(-12px); }
  .bid-move { transition: transform 0.4s ease; }
</style>
