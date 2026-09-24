<script setup lang="ts">
  import { computed, shallowRef, watch } from 'vue'

  const props = defineProps<{ value: number }>()

  // Keyed by place from the right, so the units column stays the units column when the price
  // gains a digit, and only the digits that changed roll.
  const columns = computed(() => {
    const digits = String(Math.max(0, Math.round(props.value))).split('').map(Number)
    return digits.map((digit, i) => ({ digit, place: digits.length - 1 - i }))
  })

  // Bumped on every change: a new key restarts the highlight.
  const changes = shallowRef(0)
  watch(() => props.value, () => changes.value++)
</script>

<template>
  <span class="relative isolate inline-flex font-display leading-none tabular-nums">
    <span class="sr-only">€{{ value }}</span>
    <span aria-hidden="true" class="flex">
      <span class="pr-[0.06em]">€</span>
      <span v-for="column in columns" :key="column.place" class="digit">
        <span class="strip" :style="{ transform: `translateY(${-column.digit * 10}%)` }">
          <span v-for="n in 10" :key="n">{{ n - 1 }}</span>
        </span>
      </span>
    </span>
    <span v-if="changes > 0" :key="changes" aria-hidden="true" class="glow" />
  </span>
</template>

<style scoped>
  .digit {
    display: inline-block;
    height: 1.1em;
    overflow: hidden;
  }
  .strip {
    display: flex;
    flex-direction: column;
    transition: transform 0.7s cubic-bezier(0.2, 0.8, 0.2, 1);
  }
  .strip > span {
    height: 1.1em;
    line-height: 1.1em;
  }
  .glow {
    position: absolute;
    inset: -0.08em -0.25em;
    z-index: -1;
    border-radius: 0.2em;
    background: radial-gradient(closest-side, var(--v0-gold-200), transparent);
    opacity: 0;
    animation: glow 1.1s ease-out;
  }
  @keyframes glow {
    0% { opacity: 0.9; transform: scale(1.08); }
    100% { opacity: 0; transform: scale(1); }
  }
</style>
