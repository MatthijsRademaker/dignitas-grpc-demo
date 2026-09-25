<script setup lang="ts">
  import { useId } from 'vue'

  const id = useId()
</script>

<template>
  <!-- A lot under an auction cloth: purple velvet, gold trim, and a tag that gives nothing away. -->
  <svg aria-hidden="true" class="overflow-visible" viewBox="0 0 128 128">
    <defs>
      <linearGradient :id="`${id}-velvet`" x1="0" x2="1" y1="0" y2="0.3">
        <stop offset="0" stop-color="#3b1f6e" />
        <stop offset="0.35" stop-color="#6a3bbd" />
        <stop offset="0.55" stop-color="#5a31a6" />
        <stop offset="1" stop-color="#2e1757" />
      </linearGradient>
      <filter :id="`${id}-soft`" height="200%" width="200%" x="-50%" y="-50%">
        <feGaussianBlur stdDeviation="1.6" />
      </filter>
    </defs>

    <ellipse cx="64" cy="116" fill="#1f1b2e" opacity="0.12" rx="52" ry="5" />

    <!-- The cloth, and the folds where it hangs over whatever is underneath. -->
    <path
      d="M64 16C46 16 38 38 34 60C30 82 20 98 10 110Q24 117 37 111Q50 117 64 111Q78 117 91 111Q104 117 118 110C108 98 98 82 94 60C90 38 82 16 64 16Z"
      :fill="`url(#${id}-velvet)`"
      stroke="#241246"
      stroke-width="1.5"
    />
    <g fill="none" stroke-linecap="round">
      <path d="M57 24C51 48 45 82 37 111" stroke="#241246" stroke-opacity="0.55" stroke-width="2" />
      <path d="M71 24C78 50 84 82 91 111" stroke="#241246" stroke-opacity="0.55" stroke-width="2" />
      <path d="M64 30C64 60 64 88 64 111" stroke="#9c7ae0" stroke-opacity="0.35" stroke-width="1.5" />
    </g>
    <path :filter="`url(#${id}-soft)`" d="M50 26C44 44 40 66 36 88" fill="none" stroke="#c9b3f5" stroke-linecap="round" stroke-opacity="0.55" stroke-width="3" />

    <!-- Gold trim along the hem, with tassels. -->
    <path d="M10 110Q24 117 37 111Q50 117 64 111Q78 117 91 111Q104 117 118 110" fill="none" stroke="#e0ad32" stroke-width="3" />
    <g fill="#c9971c" stroke="#7a5a00" stroke-width="0.8">
      <path v-for="x in [10, 37, 64, 91, 118]" :key="x" :d="`M${x - 2.5} 112h5l-1 7h-3z`" />
    </g>

    <!-- The lot tag. -->
    <path d="M94 50C100 46 104 44 108 42" fill="none" stroke="#7a5a00" stroke-width="1" />
    <g transform="rotate(14 112 46)">
      <rect fill="#fdf8e6" height="15" rx="2" stroke="#a67a0f" stroke-width="1.2" width="17" x="105" y="40" />
      <text fill="#7a5a00" font-family="Fraunces Variable, serif" font-size="12" font-weight="700" text-anchor="middle" x="113.5" y="51.5">?</text>
    </g>
  </svg>
</template>
