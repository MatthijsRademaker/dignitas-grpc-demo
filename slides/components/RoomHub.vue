<script setup lang="ts">
// Hub and spokes: the presenter runs the auction server, every laptop in the room runs its own
// BFF + frontend and connects to it. `streaming` animates the spokes once WatchAuction is live.
defineProps<{ streaming?: boolean }>()

const cx = 450
const cy = 165
const spokes = Array.from({ length: 8 }, (_, i) => {
  const angle = (i / 8) * Math.PI * 2 - Math.PI / 2 + Math.PI / 8
  return { x: cx + Math.cos(angle) * 330, y: cy + Math.sin(angle) * 135 }
})
</script>

<template>
  <svg viewBox="0 0 900 330" class="w-full max-w-[840px]" role="img"
    aria-label="One auction server on the presenter's laptop, with every participant's BFF connected to it over gRPC.">
    <line v-for="(s, i) in spokes" :key="`l${i}`" :x1="cx" :y1="cy" :x2="s.x" :y2="s.y"
      class="spoke" :class="{ streaming }" />

    <g v-for="(s, i) in spokes" :key="`n${i}`" :transform="`translate(${s.x} ${s.y})`">
      <rect x="-70" y="-24" width="140" height="48" rx="10" class="laptop" />
      <text y="-3" text-anchor="middle" class="laptop-title">💻 your laptop</text>
      <text y="14" text-anchor="middle" class="laptop-sub">frontend + BFF</text>
    </g>

    <g :transform="`translate(${cx} ${cy})`">
      <rect x="-105" y="-40" width="210" height="80" rx="14" class="hub" />
      <text y="-8" text-anchor="middle" class="hub-title">🔨 auction-server</text>
      <text y="16" text-anchor="middle" class="hub-sub">presenter · :50051 · Go</text>
    </g>
  </svg>
</template>

<style scoped>
.spoke {
  stroke: var(--se-color-primary);
  stroke-opacity: 0.45;
  stroke-width: 2.5;
}
.spoke.streaming {
  stroke-opacity: 0.9;
  stroke-dasharray: 8 8;
  animation: flow 0.8s linear infinite;
}
@keyframes flow {
  to {
    stroke-dashoffset: 16;
  }
}
.laptop {
  fill: var(--slidev-code-background, #fff);
  stroke: currentColor;
  stroke-opacity: 0.3;
  stroke-width: 1.5;
}
.laptop-title {
  fill: currentColor;
  font-size: 13px;
  font-weight: 800;
}
.laptop-sub {
  fill: currentColor;
  opacity: 0.55;
  font-size: 11px;
  font-weight: 600;
}
.hub {
  fill: var(--se-color-primary);
}
.hub-title {
  fill: #fff;
  font-size: 18px;
  font-weight: 800;
}
.hub-sub {
  fill: #fff;
  opacity: 0.85;
  font-size: 12px;
  font-weight: 600;
}
</style>
