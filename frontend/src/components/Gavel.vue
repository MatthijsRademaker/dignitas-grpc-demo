<script setup lang="ts">
  /** `strike`: slam once when mounted, for the SOLD moment. */
  defineProps<{ strike?: boolean }>()
</script>

<template>
  <svg aria-hidden="true" viewBox="0 0 64 64">
    <rect fill="#5c4300" height="6" rx="2" width="36" x="2" y="57" />
    <g :class="strike && 'strike'">
      <g transform="translate(20 44) rotate(-33.7)">
        <rect fill="#7a5a00" height="6" rx="3" width="42" x="2" y="-3" />
        <rect fill="#c9971c" height="28" rx="3.5" stroke="#5c4300" stroke-width="1.5" width="14" x="-7" y="-14" />
        <rect fill="#f6d77a" height="3" width="14" x="-7" y="-7" />
        <rect fill="#f6d77a" height="3" width="14" x="-7" y="4" />
      </g>
    </g>
  </svg>
</template>

<style scoped>
  /* Pivots where the hand holds the handle. */
  .strike {
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
</style>
