<script setup lang="ts">
  import { computed } from 'vue'
  import type { FeedStats, Transport } from '../useAuctionFeed'

  const props = withDefaults(defineProps<{
    transport: Transport
    pollIntervalMs: number
    stats: FeedStats
    error: string | null
    now: number
    /** `compact`: only the "bids seen after" figure, for the stage. */
    variant?: 'full' | 'compact'
  }>(), { variant: 'full' })

  const WINDOW_MS = 30_000

  const averageDelay = computed(() => {
    const { delays } = props.stats
    return delays.length === 0 ? null : delays.reduce((sum, d) => sum + d, 0) / delays.length
  })
  const delayText = computed(() => averageDelay.value === null
    ? '–'
    : `${(averageDelay.value / 1000).toFixed(averageDelay.value < 1000 ? 2 : 1)}s`)
  const fast = computed(() => averageDelay.value !== null && averageDelay.value < 300)
</script>

<template>
  <section v-if="variant === 'compact'" class="flex flex-col gap-2">
    <h3 class="text-[1.6rem] font-bold">Bids reach this screen after</h3>
    <p class="text-[7rem] leading-none font-extrabold" :class="fast && 'text-primary'">{{ delayText }}</p>
    <p class="text-[1.3rem] text-muted">
      <template v-if="transport === 'poll'">Polling every {{ pollIntervalMs / 1000 }}s: news waits for the next ask</template>
      <template v-else>Streaming: the server pushes as it happens</template>
    </p>
    <p v-if="error" class="text-[1.3rem] font-bold text-error">{{ error }}</p>
  </section>

  <section v-else class="flex flex-col gap-4">
    <div>
      <h3 class="font-bold">On the wire</h3>
      <p class="text-sm text-muted">This browser, last 30 seconds</p>
    </div>

    <div class="relative h-10 overflow-hidden rounded-xl bg-surface">
      <div class="absolute inset-y-0 right-0 w-px bg-primary/40" />
      <span
        v-for="tick in stats.ticks"
        :key="tick.at"
        class="absolute top-1/2 size-3 -translate-1/2 rounded-full"
        :class="tick.changed ? 'bg-primary' : 'bg-muted/30'"
        :style="{ left: `${100 - ((now - tick.at) / WINDOW_MS) * 100}%` }"
      />
    </div>

    <dl class="grid grid-cols-3 gap-2 text-sm">
      <div class="rounded-xl bg-surface p-3">
        <dt class="text-muted">HTTP requests</dt>
        <dd class="text-2xl font-extrabold tabular-nums">{{ stats.requests }}</dd>
      </div>
      <div class="rounded-xl bg-surface p-3">
        <dt class="text-muted">{{ transport === 'poll' ? 'With news' : 'Events pushed' }}</dt>
        <dd class="text-2xl font-extrabold tabular-nums">
          {{ transport === 'poll' ? stats.ticks.filter(t => t.changed).length : stats.messages }}
        </dd>
      </div>
      <div class="rounded-xl bg-surface p-3">
        <dt class="text-muted">Bids seen after</dt>
        <dd class="text-2xl font-extrabold tabular-nums" :class="fast && 'text-primary'">{{ delayText }}</dd>
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

    <p v-if="error" class="rounded-xl bg-error/10 px-3 py-2 text-sm text-error">{{ error }}</p>
  </section>
</template>
