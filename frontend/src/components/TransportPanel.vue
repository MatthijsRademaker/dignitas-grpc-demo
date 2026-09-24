<script setup lang="ts">
  import { computed } from 'vue'
  import type { Tick, Transport } from '../useAuctionFeed'

  const props = defineProps<{
    transport: Transport
    pollIntervalMs: number
    stats: { requests: number, messages: number, lastKind: string, ticks: Tick[], delays: number[] }
    error: string | null
    now: number
  }>()

  const WINDOW_MS = 30_000

  const averageDelay = computed(() => {
    const { delays } = props.stats
    return delays.length === 0 ? null : delays.reduce((sum, d) => sum + d, 0) / delays.length
  })
</script>

<template>
  <section class="flex flex-col gap-4 rounded-2xl border border-border bg-surface p-6 shadow-sm">
    <h3 class="text-xs font-bold tracking-widest text-muted uppercase">On the wire · this browser · last 30s</h3>

    <div class="relative h-10 overflow-hidden rounded-lg bg-background">
      <div class="absolute inset-y-0 right-0 w-px bg-primary/40" />
      <span
        v-for="tick in stats.ticks"
        :key="tick.at"
        class="absolute top-1/2 size-3 -translate-1/2 rounded-full"
        :class="tick.changed ? 'bg-primary' : 'bg-muted/30'"
        :style="{ left: `${100 - ((now - tick.at) / WINDOW_MS) * 100}%` }"
      />
    </div>

    <dl class="grid grid-cols-3 gap-3 text-sm">
      <div class="rounded-lg bg-background p-3">
        <dt class="text-muted">HTTP requests</dt>
        <dd class="text-2xl font-black tabular-nums">{{ stats.requests }}</dd>
      </div>
      <div class="rounded-lg bg-background p-3">
        <dt class="text-muted">{{ transport === 'poll' ? 'With news' : 'Events pushed' }}</dt>
        <dd class="text-2xl font-black tabular-nums">
          {{ transport === 'poll' ? stats.ticks.filter(t => t.changed).length : stats.messages }}
        </dd>
      </div>
      <div class="rounded-lg bg-background p-3">
        <dt class="text-muted">Bids seen after</dt>
        <dd class="text-2xl font-black tabular-nums" :class="averageDelay !== null && averageDelay < 300 ? 'text-primary' : ''">
          {{ averageDelay === null ? '–' : `${(averageDelay / 1000).toFixed(averageDelay < 1000 ? 2 : 1)}s` }}
        </dd>
      </div>
    </dl>

    <p class="text-sm text-muted">
      <template v-if="transport === 'poll'">
        Asking the BFF every {{ pollIntervalMs / 1000 }}s: “anything new?” Grey dots mean nothing was.
        The BFF makes a unary <span class="font-mono">GetAuction</span> call each time.
        Poll faster to see bids sooner, and the room load goes up.
      </template>
      <template v-else>
        One open connection. The server pushes only when something changes
        <span v-if="stats.lastKind">(last: <span class="font-mono">{{ stats.lastKind }}</span>)</span>.
        The BFF relays a <span class="font-mono">WatchAuction</span> gRPC stream as Server-Sent Events.
      </template>
    </p>

    <p v-if="error" class="rounded-lg bg-error/10 px-3 py-2 text-sm text-error">{{ error }}</p>
  </section>
</template>
