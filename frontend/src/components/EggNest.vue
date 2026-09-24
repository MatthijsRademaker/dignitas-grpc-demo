<script setup lang="ts">
  import { computed, useId } from 'vue'
  import type { Egg } from '../useEggs'

  const props = defineProps<{ eggs: Egg[], overflow: number }>()

  const id = useId()
  const total = computed(() => props.eggs.length + props.overflow)
  const label = computed(() => total.value === 0
    ? 'An empty nest: no bids yet this round'
    : `${total.value} golden ${total.value === 1 ? 'egg' : 'eggs'}, one for each bid this round`)
</script>

<template>
  <!-- Sized in container units, so the same nest works on a laptop and on the projector. -->
  <div :aria-label="label" class="relative w-full [container-type:inline-size]" role="img">
    <svg aria-hidden="true" class="absolute size-0">
      <defs>
        <radialGradient :id="`${id}-egg`" cx="0.38" cy="0.3" r="0.8">
          <stop offset="0" stop-color="#fff4cf" />
          <stop offset="0.35" stop-color="#f6d77a" />
          <stop offset="0.8" stop-color="#d9a52a" />
          <stop offset="1" stop-color="#a67a0f" />
        </radialGradient>
      </defs>
    </svg>

    <ol class="flex h-[19cqi] flex-wrap-reverse content-start items-end justify-center gap-x-[0.6cqi] px-[14cqi]">
      <li v-for="egg in eggs" :key="egg.key" class="egg" :class="egg.dropped && 'dropped'">
        <svg class="block size-full" viewBox="0 0 30 38">
          <path d="M15 1C7 1 1 13 1 23c0 8 6 14 14 14s14-6 14-14C29 13 23 1 15 1z" :fill="`url(#${id}-egg)`" stroke="#7a5a00" stroke-width="1.2" />
          <ellipse cx="10" cy="12" fill="#fffdf5" opacity="0.7" rx="3" ry="5" transform="rotate(20 10 12)" />
        </svg>
      </li>
    </ol>

    <!-- The front of the nest, woven from bronze twigs, drawn over the bottom row of eggs. -->
    <svg aria-hidden="true" class="relative mx-auto -mt-[8cqi] block w-[84%]" viewBox="0 0 200 40">
      <path d="M4 8 C30 16 170 16 196 8 C194 26 160 38 100 38 C40 38 6 26 4 8 Z" fill="#8a6400" />
      <g fill="none" stroke-linecap="round">
        <path d="M8 12 C40 22 160 22 192 12" stroke="#e0ad32" stroke-width="2.2" />
        <path d="M12 19 C44 29 156 29 188 19" stroke="#5c4300" stroke-width="2" />
        <path d="M22 27 C56 35 144 35 178 27" stroke="#c9971c" stroke-width="2" />
        <path d="M18 10 C50 26 90 30 130 22 M70 30 C110 34 150 26 184 12 M36 24 C70 12 120 14 164 30" stroke="#f6d77a" stroke-width="1.2" />
        <path d="M4 8 C30 16 170 16 196 8" stroke="#5c4300" stroke-width="2.5" />
      </g>
    </svg>

    <span
      v-if="overflow > 0"
      class="absolute right-[4cqi] bottom-[5cqi] rounded-full bg-surface px-[1.2cqi] font-bold text-bronze text-[3cqi] shadow-sm"
    >+{{ overflow }}</span>
  </div>
</template>

<style scoped>
  .egg {
    width: 5.2cqi;
    height: 6.6cqi;
    margin-bottom: 0.2cqi;
    transform-origin: 50% 100%;
  }
  /* Laid by the goose above: a fall, a squash, and two little bounces. */
  .dropped {
    animation: drop 0.95s cubic-bezier(0.5, 0, 0.75, 0) both;
  }
  @keyframes drop {
    0% { transform: translateY(-45cqi); opacity: 0; }
    12% { opacity: 1; }
    55% { transform: translateY(0) scale(1.08, 0.88); animation-timing-function: ease-out; }
    72% { transform: translateY(-2.5cqi) scale(0.97, 1.04); animation-timing-function: ease-in; }
    86% { transform: translateY(0) scale(1.03, 0.96); animation-timing-function: ease-out; }
    93% { transform: translateY(-0.8cqi); animation-timing-function: ease-in; }
    100% { transform: translateY(0); opacity: 1; }
  }
</style>
