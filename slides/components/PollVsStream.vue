<script setup lang="ts">
import { useSlideContext } from '@slidev/client'

const { $clicks } = useSlideContext()

const SECONDS = 20
const INTERVAL = 2
const x = (t: number) => 90 + t * 39

// Six bids in twenty seconds, bursty like a real room.
const events = [1.3, 1.9, 7.2, 12.5, 13.1, 17.8]
const polls = Array.from({ length: SECONDS / INTERVAL + 1 }, (_, i) => {
  const t = i * INTERVAL
  return { t, news: events.some(e => e > t - INTERVAL && e <= t) }
})
const delays = events.map(e => Math.ceil(e / INTERVAL) * INTERVAL - e)
const withNews = polls.filter(p => p.news).length
const maxDelay = Math.max(...delays).toFixed(1)
</script>

<template>
  <div class="w-full flex flex-col items-center">
    <svg viewBox="0 0 900 290" class="w-full max-w-[880px]" role="img"
      aria-label="Polling every two seconds mostly returns nothing new and delivers bids late; streaming delivers each bid the moment it happens.">
      <!-- time axis -->
      <line :x1="x(0)" y1="270" :x2="x(SECONDS)" y2="270" class="axis" />
      <text v-for="t in [0, 5, 10, 15, 20]" :key="t" :x="x(t)" y="286" text-anchor="middle" class="tick">{{ t }}s</text>

      <!-- what actually happened -->
      <text x="0" y="54" class="row-label">bids</text>
      <line :x1="x(0)" y1="50" :x2="x(SECONDS)" y2="50" class="guide" />
      <g v-for="e in events" :key="e" :transform="`translate(${x(e)} 50)`">
        <rect x="-7" y="-7" width="14" height="14" transform="rotate(45)" class="event" />
      </g>

      <!-- polling -->
      <text x="0" y="134" class="row-label">polling</text>
      <line :x1="x(0)" y1="130" :x2="x(SECONDS)" y2="130" class="guide" />
      <line v-for="(e, i) in events" :key="`d${i}`" :x1="x(e)" y1="58" :x2="x(e + delays[i])" y2="122" class="delay" />
      <circle v-for="p in polls" :key="p.t" :cx="x(p.t)" cy="130" r="8" :class="p.news ? 'hit' : 'miss'" />

      <!-- streaming -->
      <g v-show="$clicks >= 1">
        <text x="0" y="214" class="row-label primary">stream</text>
        <line :x1="x(0)" y1="210" :x2="x(SECONDS)" y2="210" class="guide" />
        <line :x1="x(0)" y1="210" :x2="x(SECONDS)" y2="210" class="open" />
        <line v-for="e in events" :key="`s${e}`" :x1="x(e)" y1="58" :x2="x(e)" y2="202" class="instant" />
        <circle v-for="e in events" :key="`c${e}`" :cx="x(e)" cy="210" r="8" class="hit" />
      </g>
    </svg>

    <div class="stats">
      <span class="stat">Polling: {{ polls.length }} requests · {{ withNews }} with news · up to {{ maxDelay }}s late</span>
      <span v-show="$clicks >= 1" class="stat primary">Streaming: 1 request · {{ events.length }} events · on time</span>
    </div>
  </div>
</template>

<style scoped>
.axis {
  stroke: currentColor;
  stroke-opacity: 0.3;
  stroke-width: 1.5;
}
.tick {
  fill: currentColor;
  opacity: 0.5;
  font-size: 12px;
}
.row-label {
  fill: currentColor;
  font-size: 15px;
  font-weight: 800;
}
.row-label.primary {
  fill: var(--se-color-primary);
}
.guide {
  stroke: currentColor;
  stroke-opacity: 0.12;
  stroke-width: 1;
}
.open {
  stroke: var(--se-color-primary);
  stroke-opacity: 0.35;
  stroke-width: 6;
  stroke-linecap: round;
}
.event {
  fill: currentColor;
}
.delay {
  stroke: currentColor;
  stroke-opacity: 0.35;
  stroke-width: 1.5;
  stroke-dasharray: 4 4;
}
.instant {
  stroke: var(--se-color-primary);
  stroke-opacity: 0.5;
  stroke-width: 1.5;
}
.hit {
  fill: var(--se-color-primary);
}
.miss {
  fill: color-mix(in srgb, currentColor 18%, transparent);
  stroke: currentColor;
  stroke-opacity: 0.25;
}
.stats {
  display: flex;
  gap: 1rem;
  margin-top: 0.4rem;
}
.stat {
  padding: 0.35rem 0.8rem;
  border-radius: 0.5rem;
  font-size: 0.95rem;
  font-weight: 700;
  background: color-mix(in srgb, currentColor 6%, transparent);
}
.stat.primary {
  color: var(--se-color-primary);
  background: color-mix(in srgb, var(--se-color-primary) 10%, transparent);
}
</style>
