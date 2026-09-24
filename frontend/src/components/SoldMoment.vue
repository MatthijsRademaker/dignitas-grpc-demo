<script setup lang="ts">
  import { computed } from 'vue'
  import type { Bid } from '../api'

  const props = withDefaults(defineProps<{
    winner: Bid | null
    /** Until the next round opens. */
    remainingMs: number
    round?: number
    size?: 'card' | 'fullscreen'
  }>(), { size: 'card' })

  const seconds = computed(() => Math.ceil(props.remainingMs / 1000))
  const rays = Array.from({ length: 16 }, (_, i) => i * (360 / 16))
</script>

<template>
  <div
    aria-live="polite"
    class="sold grid place-items-center overflow-hidden text-center"
    :class="[size, winner ? 'bg-surface/95' : 'bg-surface/90']"
    role="status"
  >
    <div class="relative flex flex-col items-center gap-[0.4em] px-[1em]">
      <template v-if="winner">
        <svg aria-hidden="true" class="burst motion-decoration" viewBox="-100 -100 200 200">
          <path
            v-for="angle in rays"
            :key="angle"
            d="M-4 -30 L0 -98 L4 -30 Z"
            :fill="angle % 45 === 0 ? 'var(--v0-gold-300)' : 'var(--v0-gold-100)'"
            :transform="`rotate(${angle})`"
          />
        </svg>

        <svg aria-hidden="true" class="gavel relative" viewBox="0 0 64 64">
          <rect fill="#5c4300" height="6" rx="2" width="36" x="2" y="57" />
          <g class="swing">
            <g transform="translate(20 44) rotate(-33.7)">
              <rect fill="#7a5a00" height="6" rx="3" width="42" x="2" y="-3" />
              <rect fill="#c9971c" height="28" rx="3.5" stroke="#5c4300" stroke-width="1.5" width="14" x="-7" y="-14" />
              <rect fill="#f6d77a" height="3" width="14" x="-7" y="-7" />
              <rect fill="#f6d77a" height="3" width="14" x="-7" y="4" />
            </g>
          </g>
        </svg>

        <p v-if="round" class="relative font-bold text-bronze">Round {{ round }}</p>
        <p class="relative font-display text-[2.4em] leading-tight font-semibold">
          SOLD to {{ winner.bidder }} · €{{ winner.amount }}
        </p>
      </template>

      <template v-else>
        <p class="font-display text-[1.8em] leading-tight font-medium text-muted">The goose flew off, nobody bid</p>
      </template>

      <p class="relative text-[0.9em] text-muted tabular-nums">next round in {{ seconds }}s</p>
    </div>
  </div>
</template>

<style scoped>
  .card {
    font-size: 1rem;
  }
  .fullscreen {
    font-size: 2.6rem;
  }
  .gavel {
    width: 4.5em;
  }
  .swing {
    transform-origin: 53px 22px;
    animation: slam 0.7s cubic-bezier(0.3, 0, 0.4, 1) both;
  }
  @keyframes slam {
    0% { transform: rotate(-38deg); }
    55% { transform: rotate(5deg); }
    70% { transform: rotate(-4deg); }
    85% { transform: rotate(1deg); }
    100% { transform: rotate(0); }
  }
  .burst {
    position: absolute;
    top: 50%;
    left: 50%;
    width: 22em;
    translate: -50% -50%;
    animation: burst 1.2s 0.35s ease-out both, turn 60s linear 1.5s infinite;
  }
  @keyframes burst {
    0% { transform: scale(0.2); opacity: 0; }
    40% { opacity: 1; }
    100% { transform: scale(1); opacity: 0.55; }
  }
  @keyframes turn {
    from { transform: rotate(0); opacity: 0.55; }
    to { transform: rotate(360deg); opacity: 0.55; }
  }
</style>
