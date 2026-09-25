<script setup lang="ts">
import { useSlideContext } from '@slidev/client'

const { $clicks } = useSlideContext()

// HTTP/1.1: each connection carries one request/response at a time.
const lanes = [
  [{ x: 0, w: 110, label: 'GET auction' }, { x: 190, w: 110, label: 'GET auction' }],
  [{ x: 40, w: 110, label: 'POST bid' }, { x: 230, w: 110, label: 'GET auction' }],
  [{ x: 90, w: 110, label: 'GET auction' }],
]

// HTTP/2: frames of three streams interleaved on a single connection.
// w = WatchAuction (long-lived), g = GetAuction, p = PlaceBid
const frames = 'wgwpwgwwpwgw'.split('')
</script>

<template>
  <div class="w-full flex flex-col items-center">
    <svg viewBox="0 0 900 270" class="w-full max-w-[880px]" role="img"
      aria-label="HTTP/1.1 needs a connection per concurrent request; HTTP/2 interleaves many streams over one connection.">
      <!-- HTTP/1.1 -->
      <text x="10" y="24" class="title">HTTP/1.1</text>
      <g v-for="(lane, i) in lanes" :key="i" :transform="`translate(10 ${50 + i * 60})`">
        <rect x="0" y="0" width="400" height="40" rx="8" class="track" />
        <text x="-2" y="-6" class="lane-label">connection {{ i + 1 }}</text>
        <g v-for="(req, j) in lane" :key="j">
          <rect :x="req.x + 6" y="6" :width="req.w" height="28" rx="6" class="req" />
          <text :x="req.x + 6 + req.w / 2" y="25" text-anchor="middle" class="req-text">{{ req.label }}</text>
        </g>
      </g>
      <text x="10" y="252" class="caption">one request at a time per connection</text>

      <!-- HTTP/2 -->
      <g v-show="$clicks >= 1">
        <text x="480" y="24" class="title primary">HTTP/2</text>
        <text x="478" y="104" class="lane-label">one connection</text>
        <rect x="480" y="110" width="410" height="52" rx="10" class="pipe" />
        <g v-for="(f, i) in frames" :key="i">
          <rect :x="488 + i * 33.5" y="120" width="28" height="32" rx="5" :class="['frame', f]" />
        </g>
        <g transform="translate(480 190)" class="legend">
          <rect x="0" y="0" width="14" height="14" rx="3" class="frame w" />
          <text x="20" y="12">WatchAuction (stays open)</text>
          <rect x="210" y="0" width="14" height="14" rx="3" class="frame g" />
          <text x="230" y="12">GetAuction</text>
          <rect x="320" y="0" width="14" height="14" rx="3" class="frame p" />
          <text x="340" y="12">PlaceBid</text>
        </g>
        <text x="480" y="252" class="caption">many streams, interleaved as binary frames</text>
      </g>
    </svg>
  </div>
</template>

<style scoped>
.title {
  fill: currentColor;
  font-size: 20px;
  font-weight: 800;
}
.title.primary {
  fill: var(--se-color-primary);
}
.track {
  fill: color-mix(in srgb, currentColor 5%, transparent);
  stroke: currentColor;
  stroke-opacity: 0.2;
}
.lane-label {
  fill: currentColor;
  opacity: 0.5;
  font-size: 11px;
  font-weight: 700;
}
.req {
  fill: color-mix(in srgb, currentColor 18%, transparent);
}
.req-text {
  fill: currentColor;
  font-size: 12px;
  font-weight: 700;
}
.pipe {
  fill: color-mix(in srgb, var(--se-color-primary) 8%, transparent);
  stroke: var(--se-color-primary);
  stroke-width: 2;
}
.frame.w {
  fill: var(--se-color-primary);
}
.frame.g {
  fill: color-mix(in srgb, var(--se-color-primary) 45%, transparent);
}
.frame.p {
  fill: color-mix(in srgb, currentColor 55%, transparent);
}
.legend text {
  fill: currentColor;
  font-size: 12px;
  font-weight: 600;
}
.caption {
  fill: currentColor;
  opacity: 0.65;
  font-size: 14px;
  font-weight: 700;
}
</style>
