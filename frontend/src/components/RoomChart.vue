<script setup lang="ts">
  import { computed, shallowRef, useId } from 'vue'
  import type { ServerLoad } from '../api'
  import { ROOM_HISTORY_MS, type RoomSample } from '../useRoomHistory'

  const props = withDefaults(defineProps<{
    samples: RoomSample[]
    /** The latest figures, for the secondary line under each panel title. */
    load: ServerLoad
    now: number
    size?: 'md' | 'xl'
  }>(), { size: 'md' })

  type Key = 'pollsPerSecond' | 'watchers'

  const helpId = `${useId()}-help`

  const rate = (perSecond: number) => perSecond.toFixed(perSecond < 10 ? 1 : 0)
  const count = (n: number) => String(Math.round(n))

  // Two small multiples on one time axis, never one chart with two y-scales: the units differ.
  const PANELS = [
    { key: 'pollsPerSecond' as Key, title: 'Polling calls per second', unit: 'calls/s', color: 'var(--series-poll)', format: rate },
    { key: 'watchers' as Key, title: 'Open streams', unit: 'open', color: 'var(--series-stream)', format: count },
  ]

  /** 0 at the left edge (90s ago), 1 at now. */
  const position = (at: number) => Math.min(1, Math.max(0, 1 - (props.now - at) / ROOM_HISTORY_MS))

  /** The value in force at `at`: the samples are steps, each one holds until the next. */
  function valueAt (key: Key, at: number): number | null {
    let value: number | null = null
    for (const sample of props.samples) {
      if (sample.at > at) break
      value = sample[key]
    }
    return value
  }

  /** A round number above the peak, so the line never touches the top of its panel. */
  function niceMax (peak: number) {
    const target = Math.max(peak, 1) * 1.15
    const magnitude = 10 ** Math.floor(Math.log10(target))
    return [1, 2, 2.5, 5, 10].map(step => step * magnitude).find(max => max >= target)!
  }

  // Paths are drawn in a 1000×100 box stretched to the panel; strokes stay 2px regardless.
  const panels = computed(() => PANELS.map(panel => {
    const values = props.samples.map(sample => sample[panel.key])
    const peak = Math.max(0, ...values)
    const max = niceMax(peak)
    const y = (value: number) => 100 - (value / max) * 100
    const last = values.at(-1) ?? null

    let line = ''
    if (props.samples.length > 0) {
      const [first, ...rest] = props.samples
      line = `M${position(first!.at) * 1000},${y(first![panel.key])}`
      for (const sample of rest) line += `H${position(sample.at) * 1000}V${y(sample[panel.key])}`
      line += 'H1000' // carry the last value forward to now
    }
    const area = line && `${line}V100H${position(props.samples[0]!.at) * 1000}Z`

    return {
      ...panel,
      max,
      line,
      area,
      last,
      /** The end of the line, as a fraction from the top. */
      end: last === null ? null : y(last) / 100,
      label: last === null
        ? `${panel.title}: no samples yet`
        : `${panel.title}, last 90 seconds: now ${panel.format(last)}, peak ${panel.format(peak)}`,
    }
  }))

  // The crosshair, as a fraction of the plot width, so it stays under the pointer as time moves.
  const hover = shallowRef<number | null>(null)
  const hoverAt = computed(() => hover.value === null ? null : props.now - (1 - hover.value) * ROOM_HISTORY_MS)
  const ago = (at: number) => {
    const seconds = Math.round((props.now - at) / 1000)
    return seconds <= 0 ? 'now' : `${seconds}s ago`
  }

  function onPointer (event: PointerEvent) {
    const box = (event.currentTarget as HTMLElement).getBoundingClientRect()
    hover.value = Math.min(1, Math.max(0, (event.clientX - box.left) / box.width))
  }

  function onKey (event: KeyboardEvent) {
    const step = 5000 / ROOM_HISTORY_MS
    const current = hover.value ?? 1
    const next = {
      ArrowLeft: current - step,
      ArrowRight: current + step,
      Home: 0,
      End: 1,
    }[event.key]
    if (event.key === 'Escape') hover.value = null
    if (next === undefined) return
    event.preventDefault()
    hover.value = Math.min(1, Math.max(0, next))
  }

  const tableRows = computed(() => Array.from({ length: 7 }, (_, i) => {
    const at = props.now - (6 - i) * 15_000
    return { at, polls: valueAt('pollsPerSecond', at), watchers: valueAt('watchers', at) }
  }))
</script>

<template>
  <figure class="room-chart flex flex-col" :class="size">
    <figcaption class="caption font-bold">The whole room, as the auction server sees it</figcaption>

    <div
      class="plots relative flex min-h-0 flex-1 flex-col outline-none"
      tabindex="0"
      :aria-describedby="helpId"
      @blur="hover = null"
      @focus="hover ??= 1"
      @keydown="onKey"
    >
      <span :id="helpId" class="sr-only">Use the arrow keys to read earlier values.</span>

      <div v-for="panel in panels" :key="panel.key" class="panel flex min-h-0 flex-1 flex-col">
        <div class="panel-head flex items-baseline gap-2">
          <svg aria-hidden="true" class="key shrink-0" viewBox="0 0 16 4"><line :stroke="panel.color" stroke-linecap="round" stroke-width="3" x1="2" x2="14" y1="2" y2="2" /></svg>
          <span class="font-bold">{{ panel.title }}</span>
          <span class="text-muted">
            <template v-if="panel.key === 'pollsPerSecond'">
              <template v-if="load.pollsPerSecond > 0">{{ Math.round(load.emptyPollRatio * 100) }}% found nothing new</template>
              <template v-else>nobody is polling</template>
            </template>
            <template v-else>{{ rate(load.pushesPerSecond) }} pushes/s, only when something changed</template>
          </span>
        </div>

        <div class="row flex min-h-0 flex-1">
          <div aria-hidden="true" class="ticks flex flex-col justify-between text-right text-muted tabular-nums">
            <span>{{ panel.format(panel.max) }}</span>
            <span>0</span>
          </div>

          <div
            :aria-label="panel.label"
            class="plot relative min-h-0 flex-1"
            role="img"
            @pointerleave="hover = null"
            @pointermove="onPointer"
          >
            <svg aria-hidden="true" class="absolute inset-0 size-full overflow-visible" preserveAspectRatio="none" viewBox="0 0 1000 100">
              <line class="grid-line" vector-effect="non-scaling-stroke" x1="0" x2="1000" y1="0" y2="0" />
              <line class="grid-line" vector-effect="non-scaling-stroke" x1="0" x2="1000" y1="100" y2="100" />
              <path v-if="panel.area" :d="panel.area" :fill="panel.color" fill-opacity="0.1" />
              <path
                v-if="panel.line"
                :d="panel.line"
                fill="none"
                :stroke="panel.color"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                vector-effect="non-scaling-stroke"
              />
            </svg>
            <span v-if="hover !== null" class="crosshair" :style="{ left: `${hover * 100}%` }" />
            <span v-if="panel.end !== null" class="end-dot" :style="{ top: `${panel.end * 100}%`, background: panel.color }" />
            <p v-if="!panel.line" class="absolute inset-0 grid place-items-center text-muted">Waiting for the first figures</p>
          </div>

          <div aria-hidden="true" class="gutter relative">
            <span
              v-if="panel.last !== null"
              class="value absolute -translate-y-1/2 leading-none font-extrabold whitespace-nowrap"
              :style="{ top: `${Math.min(88, Math.max(12, panel.end! * 100))}%` }"
            >{{ panel.format(panel.last) }}<span class="unit font-bold text-muted"> {{ panel.unit }}</span></span>
          </div>
        </div>
      </div>

      <div aria-hidden="true" class="axis flex text-muted tabular-nums">
        <span class="ticks" />
        <span class="relative flex-1">
          <span class="absolute left-0">90s ago</span>
          <span class="absolute left-1/3 -translate-x-1/2">60s</span>
          <span class="absolute left-2/3 -translate-x-1/2">30s</span>
          <span class="absolute right-0">now</span>
        </span>
        <span class="gutter" />
      </div>

      <div
        v-if="hoverAt !== null"
        aria-live="polite"
        class="tooltip pointer-events-none absolute top-0 z-10 rounded-lg bg-on-surface px-3 py-2 text-surface shadow-lg"
        :style="{ left: `calc(var(--ticks) + (100% - var(--ticks) - var(--gutter)) * ${hover})` }"
      >
        <p class="text-surface/70">{{ ago(hoverAt) }}</p>
        <p v-for="panel in panels" :key="panel.key" class="flex items-center gap-2 whitespace-nowrap">
          <svg aria-hidden="true" class="key" viewBox="0 0 16 4"><line :stroke="panel.color" stroke-linecap="round" stroke-width="3" x1="2" x2="14" y1="2" y2="2" /></svg>
          <strong class="tabular-nums">{{ valueAt(panel.key, hoverAt) === null ? '–' : panel.format(valueAt(panel.key, hoverAt)!) }}</strong>
          <span class="text-surface/70">{{ panel.title.toLowerCase() }}</span>
        </p>
      </div>
    </div>

    <table class="sr-only">
      <caption>Room load over the last 90 seconds</caption>
      <thead>
        <tr><th scope="col">When</th><th scope="col">Polling calls per second</th><th scope="col">Open streams</th></tr>
      </thead>
      <tbody>
        <tr v-for="row in tableRows" :key="row.at">
          <th scope="row">{{ ago(row.at) }}</th>
          <td>{{ row.polls === null ? 'no data' : rate(row.polls) }}</td>
          <td>{{ row.watchers === null ? 'no data' : count(row.watchers) }}</td>
        </tr>
      </tbody>
    </table>
  </figure>
</template>

<style scoped>
  .room-chart {
    /* Validated as a pair (dataviz validate_palette.js, light): CVD ΔE 28.8, both ≥ 3:1 on white. */
    --series-poll: #d9480f;
    --series-stream: var(--v0-primary);
  }

  .md {
    --ticks: 2.25rem;
    --gutter: 6.5rem;
    height: 19rem;
    gap: 0.75rem;
    font-size: 0.875rem;
  }
  .md .value { left: 0.75rem; font-size: 1.5rem; }
  .md .unit { font-size: 0.8rem; }
  .md .panel { gap: 0.4rem; }
  .md .panel + .panel { margin-top: 1rem; }
  .md .axis { height: 1.5rem; padding-top: 0.35rem; font-size: 0.75rem; }
  .md .ticks { font-size: 0.75rem; }
  .md .key { width: 1rem; }

  .xl {
    --ticks: 3.5rem;
    --gutter: 13rem;
    height: 100%;
    gap: 1rem;
    font-size: 1.5rem;
  }
  .xl .caption { font-size: 1.6rem; }
  .xl .value { left: 1.25rem; font-size: 3.5rem; }
  .xl .unit { font-size: 1.4rem; }
  .xl .panel { gap: 0.5rem; }
  .xl .panel + .panel { margin-top: 1.5rem; }
  .xl .axis { height: 2.25rem; padding-top: 0.5rem; font-size: 1.1rem; }
  .xl .ticks { font-size: 1.1rem; }
  .xl .key { width: 1.6rem; }

  .ticks { width: var(--ticks); padding-right: 0.5rem; }
  .gutter { width: var(--gutter); }
  .ticks > span { translate: 0 -0.5em; }
  .ticks > span:last-child { translate: 0 0.5em; }

  .grid-line {
    stroke: var(--v0-border);
    stroke-width: 1;
  }
  .end-dot {
    position: absolute;
    right: 0;
    width: 0.625rem;
    height: 0.625rem;
    border-radius: 9999px;
    box-shadow: 0 0 0 2px var(--v0-surface);
    translate: 50% -50%;
  }
  .xl .end-dot { width: 1rem; height: 1rem; }
  .crosshair {
    position: absolute;
    top: 0;
    bottom: 0;
    width: 1px;
    background: var(--v0-muted);
  }
  .tooltip {
    translate: -50% calc(-100% - 0.5rem);
    font-size: 0.8125rem;
  }
  .plots:focus-visible {
    outline: 2px solid var(--v0-primary);
    outline-offset: 6px;
    border-radius: 0.5rem;
  }
</style>
