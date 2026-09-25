<script setup lang="ts">
  // Artwork: the duck from Noto Emoji (github.com/googlefonts/noto-emoji, 2D/svg/emoji_u1f986.svg),
  // Copyright 2013 Google Inc., licensed under the Apache License 2.0. Modified: the mallard's fills
  // are replaced by gold gradients, with an outline, a highlight and a shimmer added.
  // See frontend/THIRD_PARTY_NOTICES.md.
  import { useId } from 'vue'
  import * as duck from '../art/duck'

  withDefaults(defineProps<{ shimmer?: boolean }>(), { shimmer: true })

  // Several ducks can be on one page: every gradient, filter and mask needs its own id.
  const id = useId()
  const url = (name: string) => `url(#${id}-${name})`
</script>

<template>
  <svg aria-hidden="true" class="overflow-visible" viewBox="0 0 128 128">
    <defs>
      <linearGradient :id="`${id}-head`" gradientUnits="userSpaceOnUse" x1="46" x2="52" y1="6" y2="52">
        <stop offset="0" stop-color="#fbe9a8" />
        <stop offset="0.45" stop-color="#eec455" />
        <stop offset="1" stop-color="#b8870f" />
      </linearGradient>
      <linearGradient :id="`${id}-body`" gradientUnits="userSpaceOnUse" x1="52" x2="72" y1="50" y2="102">
        <stop offset="0" stop-color="#f6d77a" />
        <stop offset="0.45" stop-color="#e0ad32" />
        <stop offset="1" stop-color="#9a700a" />
      </linearGradient>
      <linearGradient :id="`${id}-wing`" gradientUnits="userSpaceOnUse" x1="64" x2="80" y1="50" y2="86">
        <stop offset="0" stop-color="#fff4cf" />
        <stop offset="0.5" stop-color="#f6d77a" />
        <stop offset="1" stop-color="#d9a52a" />
      </linearGradient>
      <linearGradient :id="`${id}-tail`" gradientUnits="userSpaceOnUse" x1="100" x2="90" y1="44" y2="98">
        <stop offset="0" stop-color="#c9971c" />
        <stop offset="1" stop-color="#7a5a00" />
      </linearGradient>
      <linearGradient :id="`${id}-tail-inner`" gradientUnits="userSpaceOnUse" x1="96" x2="90" y1="52" y2="94">
        <stop offset="0" stop-color="#f6d77a" />
        <stop offset="1" stop-color="#c9971c" />
      </linearGradient>
      <linearGradient :id="`${id}-beak`" gradientUnits="userSpaceOnUse" x1="20" x2="20" y1="15" y2="31">
        <stop offset="0" stop-color="#f3c14a" />
        <stop offset="1" stop-color="#b8870f" />
      </linearGradient>
      <linearGradient :id="`${id}-feet`" gradientUnits="userSpaceOnUse" x1="60" x2="60" y1="96" y2="126">
        <stop offset="0" stop-color="#c9971c" />
        <stop offset="1" stop-color="#8a6400" />
      </linearGradient>
      <linearGradient :id="`${id}-shine`" x1="0" x2="1" y1="0" y2="0">
        <stop offset="0" stop-color="#fffdf5" stop-opacity="0" />
        <stop offset="0.5" stop-color="#fffdf5" stop-opacity="0.75" />
        <stop offset="1" stop-color="#fffdf5" stop-opacity="0" />
      </linearGradient>
      <filter :id="`${id}-soft`" height="200%" width="200%" x="-50%" y="-50%">
        <feGaussianBlur stdDeviation="1.4" />
      </filter>
      <mask :id="`${id}-shape`">
        <path v-for="d in duck.silhouette" :key="d" :d fill="#fff" />
      </mask>
    </defs>

    <!-- The outline: every shape stroked behind the fills, so only the outer edge shows. -->
    <g fill="#5c4300" stroke="#5c4300" stroke-linejoin="round" stroke-width="3.5">
      <path v-for="d in duck.silhouette" :key="d" :d />
    </g>

    <path :d="duck.footBack" :fill="url('feet')" />
    <path :d="duck.footFront" :fill="url('feet')" />
    <path :d="duck.tail" :fill="url('tail')" />
    <path :d="duck.tailInner" :fill="url('tail-inner')" />
    <path :d="duck.body" :fill="url('body')" />
    <path :d="duck.head" :fill="url('head')" />
    <path :d="duck.beak" :fill="url('beak')" />
    <path :d="duck.collar" fill="#fff6d6" />
    <path :d="duck.wing" :fill="url('wing')" />
    <path :d="duck.eye" fill="#3d2c00" />
    <path :d="duck.eyeGlint" fill="#fffaf0" opacity="0.8" />

    <!-- Specular highlights: where the light hits the crown and the back. -->
    <g :filter="url('soft')" fill="#fffdf5">
      <ellipse cx="45" cy="13" opacity="0.7" rx="9" ry="3.5" transform="rotate(-18 45 13)" />
      <ellipse cx="46" cy="64" opacity="0.45" rx="10" ry="3" transform="rotate(-8 46 64)" />
      <ellipse cx="72" cy="58" opacity="0.5" rx="8" ry="2.2" />
    </g>

    <g v-if="shimmer" class="motion-decoration" :mask="url('shape')">
      <g transform="skewX(-20)">
        <rect class="shimmer-band" :fill="url('shine')" height="140" width="36" x="0" y="-6" />
      </g>
    </g>
  </svg>
</template>

<style scoped>
  .shimmer-band {
    animation: shimmer 6s ease-in-out infinite;
    transform: translateX(-60px);
  }
  @keyframes shimmer {
    0%, 60% { transform: translateX(-60px); }
    100% { transform: translateX(190px); }
  }
</style>
