<script setup lang="ts">
  import type { Lot } from '../api'
  import GoldenDuck from './GoldenDuck.vue'
  import VeiledLot from './VeiledLot.vue'

  /** `veiled`: PlaceBid isn't implemented yet, so the lot stays under its cloth. */
  defineProps<{ lot: Lot, shimmer?: boolean, veiled?: boolean }>()
</script>

<template>
  <!-- Cloth and artwork share one cell: on reveal the cloth lifts off while the lot fades in underneath. -->
  <div class="grid size-full *:col-start-1 *:row-start-1">
    <Transition name="unveil">
      <!-- Artwork we drew for a lot, looked up by id. Any other lot falls back to its emoji. -->
      <GoldenDuck v-if="!veiled && lot.id === 'golden-goose'" class="size-full" :shimmer />
      <div v-else-if="!veiled" aria-hidden="true" class="grid size-full place-items-center [container-type:size]">
        <span class="text-[70cqmin] leading-none">{{ lot.emoji }}</span>
      </div>
    </Transition>
    <Transition name="lift">
      <VeiledLot v-if="veiled" class="size-full" />
    </Transition>
  </div>
</template>

<style scoped>
  .lift-leave-active {
    animation: lift 1.1s cubic-bezier(0.5, 0, 0.8, 0.4) both;
    transform-origin: 50% 20%;
  }
  @keyframes lift {
    0% { transform: none; }
    25% { transform: translateY(4%) scaleY(0.96); }
    100% { transform: translate(18%, -130%) rotate(16deg); opacity: 0; }
  }
  .unveil-enter-active {
    animation: unveil 0.9s 0.35s ease-out both;
  }
  @keyframes unveil {
    0% { transform: scale(0.82); opacity: 0; filter: brightness(1.8) drop-shadow(0 0 0 transparent); }
    60% { filter: brightness(1.25) drop-shadow(0 0 18px var(--v0-gold-300)); }
    100% { transform: none; opacity: 1; filter: none; }
  }
</style>
