<script setup lang="ts">
  import type { ServerLoad } from '../api'

  defineProps<{ load: ServerLoad, watchers: number }>()

  const rate = (perSecond: number) => perSecond.toFixed(perSecond < 10 ? 1 : 0)
</script>

<template>
  <section class="rounded-2xl border border-border bg-surface p-6 shadow-sm">
    <h3 class="mb-3 text-xs font-bold tracking-widest text-muted uppercase">The whole room, as the auction server sees it</h3>

    <div class="grid grid-cols-2 gap-3">
      <div class="rounded-lg bg-background p-4" :class="load.pollsPerSecond > 0 ? '' : 'opacity-50'">
        <p class="text-xs font-bold tracking-widest text-muted uppercase">Polling</p>
        <p class="text-4xl font-black tabular-nums">
          {{ rate(load.pollsPerSecond) }}<span class="text-lg font-bold text-muted"> calls/s</span>
        </p>
        <p class="mt-1 text-sm text-muted">
          <template v-if="load.pollsPerSecond > 0">
            <strong class="text-on-surface">{{ Math.round(load.emptyPollRatio * 100) }}%</strong> found nothing new
          </template>
          <template v-else>nobody is polling</template>
        </p>
      </div>

      <div class="rounded-lg bg-background p-4" :class="watchers > 0 ? 'ring-2 ring-primary/40' : 'opacity-50'">
        <p class="text-xs font-bold tracking-widest text-primary uppercase">Streaming</p>
        <p class="text-4xl font-black tabular-nums">
          {{ watchers }}<span class="text-lg font-bold text-muted"> open</span>
        </p>
        <p class="mt-1 text-sm text-muted">
          <strong class="text-on-surface">{{ rate(load.pushesPerSecond) }}</strong> pushes/s, only when something changed
        </p>
      </div>
    </div>
  </section>
</template>
