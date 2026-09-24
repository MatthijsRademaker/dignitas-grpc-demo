<script setup lang="ts">
  import type { Lot } from '../api'
  import GoldenDuck from './GoldenDuck.vue'

  defineProps<{ lot: Lot, shimmer?: boolean }>()
</script>

<template>
  <!-- Artwork we drew for a lot, looked up by id. Any other lot falls back to its emoji. -->
  <GoldenDuck v-if="lot.id === 'golden-goose'" class="size-full" :shimmer />
  <div v-else aria-hidden="true" class="grid size-full place-items-center [container-type:size]">
    <span class="text-[70cqmin] leading-none">{{ lot.emoji }}</span>
  </div>
</template>
