<script setup lang="ts">
  import type { Sale } from '../api'

  withDefaults(defineProps<{ winners: Sale[], layout?: 'list' | 'strip' }>(), { layout: 'list' })
</script>

<template>
  <section :class="layout">
    <h3 class="heading font-bold">Owners of the goose</h3>
    <p v-if="winners.length === 0" class="text-muted">
      Nobody has taken the goose home yet. Win this round and your name goes up here first.
    </p>
    <ol v-else class="entries">
      <li v-for="sale in winners" :key="sale.round" class="entry">
        <svg aria-hidden="true" class="egg shrink-0" viewBox="0 0 30 38">
          <path d="M15 1C7 1 1 13 1 23c0 8 6 14 14 14s14-6 14-14C29 13 23 1 15 1z" fill="var(--v0-gold-300)" stroke="var(--v0-gold-700)" stroke-width="2" />
        </svg>
        <span class="text-muted tabular-nums">Round {{ sale.round }}</span>
        <strong class="name truncate">{{ sale.bidder }}</strong>
        <span class="font-display font-semibold tabular-nums">€{{ sale.amount }}</span>
      </li>
    </ol>
  </section>
</template>

<style scoped>
  .list .heading { margin-bottom: 0.75rem; }
  .list .entries { display: flex; flex-direction: column; gap: 0.25rem; }
  .list .entry {
    display: grid;
    grid-template-columns: auto 4.5rem minmax(0, 1fr) auto;
    align-items: center;
    gap: 0.75rem;
    padding: 0.4rem 0;
    border-bottom: 1px solid var(--v0-border);
  }
  .list .entry:first-child .name { color: var(--v0-bronze); }
  .list .egg { width: 0.9rem; }

  .strip { display: flex; align-items: center; gap: 2rem; font-size: 1.5rem; white-space: nowrap; }
  .strip .heading { flex-shrink: 0; }
  .strip .entries {
    display: flex;
    flex: 1;
    min-width: 0;
    gap: 2.5rem;
    overflow: hidden;
    /* Fills the row, so only an entry that runs off the screen edge fades, instead of being cut off mid-name. */
    mask-image: linear-gradient(to right, #000 90%, transparent);
  }
  .strip .entry { display: flex; align-items: center; gap: 0.6rem; }
  .strip .egg { width: 1.2rem; }
</style>
