<script setup lang="ts">
  import { Toggle } from '@vuetify/v0'
  import { POLL_INTERVALS_MS, type Transport } from '../useAuctionFeed'

  /** `quiet`: for the stage, where the presenter needs the toggles but the room shouldn't look at them. */
  defineProps<{ quiet?: boolean }>()
  const transport = defineModel<Transport>('transport', { required: true })
  const pollIntervalMs = defineModel<number>('pollIntervalMs', { required: true })
</script>

<template>
  <div class="flex flex-wrap items-center gap-2" :class="quiet && 'opacity-60 transition-opacity focus-within:opacity-100 hover:opacity-100'">
    <Toggle.Group
      v-if="transport === 'poll'"
      v-model="pollIntervalMs"
      class="flex rounded-xl border border-border bg-surface p-1"
      label="How often to poll"
      mandatory
    >
      <Toggle.Root
        v-for="ms in POLL_INTERVALS_MS"
        :key="ms"
        class="cursor-pointer rounded-lg px-3 py-1.5 text-sm font-bold text-muted tabular-nums transition data-[state=on]:bg-on-surface data-[state=on]:text-surface"
        :value="ms"
      >
        {{ ms / 1000 }}s
      </Toggle.Root>
    </Toggle.Group>

    <Toggle.Group
      v-model="transport"
      class="flex rounded-xl border border-border bg-surface p-1"
      label="How this page gets updates"
      mandatory
    >
      <Toggle.Root
        class="cursor-pointer rounded-lg px-4 py-1.5 text-sm font-bold text-muted transition data-[state=on]:bg-primary data-[state=on]:text-on-primary"
        value="poll"
      >
        Polling · every {{ pollIntervalMs / 1000 }}s
      </Toggle.Root>
      <Toggle.Root
        class="cursor-pointer rounded-lg px-4 py-1.5 text-sm font-bold text-muted transition data-[state=on]:bg-primary data-[state=on]:text-on-primary"
        value="stream"
      >
        Streaming · SSE ← gRPC
      </Toggle.Root>
    </Toggle.Group>
  </div>
</template>
