<script setup lang="ts">
  import { computed } from 'vue'

  const props = defineProps<{ remainingMs: number, lotDurationMs?: number }>()

  const RADIUS = 45
  const CIRCUMFERENCE = 2 * Math.PI * RADIUS

  const seconds = computed(() => Math.ceil(props.remainingMs / 1000))
  const text = computed(() => `${Math.floor(seconds.value / 60)}:${String(seconds.value % 60).padStart(2, '0')}`)
  const closing = computed(() => seconds.value <= 10)
  // A late bid can extend the round past its configured duration: then the ring is simply full.
  const fill = computed(() => {
    const full = Math.max(props.lotDurationMs ?? 0, props.remainingMs)
    return full > 0 ? props.remainingMs / full : 0
  })
</script>

<template>
  <div class="flex flex-col items-center gap-[0.35em] [container-type:inline-size]">
    <div class="relative w-full" :class="closing && 'animate-pulse-closing'">
      <svg aria-hidden="true" class="block w-full -rotate-90" viewBox="0 0 100 100">
        <circle class="stroke-border" cx="50" cy="50" fill="none" :r="RADIUS" stroke-width="5" />
        <circle
          class="ring"
          :class="closing ? 'stroke-error' : 'stroke-primary'"
          cx="50"
          cy="50"
          fill="none"
          :r="RADIUS"
          :stroke-dasharray="CIRCUMFERENCE"
          :stroke-dashoffset="CIRCUMFERENCE * (1 - fill)"
          stroke-linecap="round"
          stroke-width="5"
        />
      </svg>
      <span
        class="absolute inset-0 grid place-items-center font-display text-[26cqi] leading-none font-semibold tabular-nums"
        :class="closing && 'text-error'"
        role="timer"
      >{{ text }}</span>
    </div>
    <span class="text-center font-bold" :class="closing ? 'text-error' : 'text-muted'">
      {{ closing ? 'Going once, going twice…' : 'left in this round' }}
    </span>
  </div>
</template>

<style scoped>
  .ring {
    transition: stroke-dashoffset 0.25s linear, stroke 0.3s;
  }
</style>
