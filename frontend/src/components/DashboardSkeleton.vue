<script setup lang="ts">
  defineProps<{ error: string | null }>()
</script>

<template>
  <!-- Shaped like what is coming, so the page doesn't jump when the first auction arrives. -->
  <div aria-busy="true" class="flex flex-col gap-6">
    <p
      class="rounded-2xl px-4 py-3 font-bold"
      :class="error ? 'bg-error/10 text-error' : 'text-muted'"
      role="status"
    >
      {{ error ?? 'Connecting to the auction…' }}
    </p>
    <div class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_22rem]">
      <div class="flex flex-col gap-6">
        <div class="bone grid h-80 grid-cols-[5fr_7fr] gap-6 rounded-[1.75rem] p-8">
          <div class="rounded-full bg-border/70" />
          <div class="flex flex-col gap-3 pt-4">
            <div class="h-5 w-24 rounded-full bg-border/70" />
            <div class="h-9 w-3/4 rounded-lg bg-border/70" />
            <div class="h-4 w-full rounded bg-border/70" />
            <div class="mt-auto h-16 w-1/2 rounded-lg bg-border/70" />
          </div>
        </div>
        <div class="grid gap-6 md:grid-cols-2">
          <div class="flex flex-col gap-2">
            <div v-for="n in 5" :key="n" class="bone h-12 rounded-xl" />
          </div>
          <div class="bone h-64 rounded-xl" />
        </div>
        <div class="bone h-72 rounded-xl" />
      </div>
      <div class="bone h-96 rounded-3xl" />
    </div>
  </div>
</template>

<style scoped>
  .bone {
    background: var(--v0-surface);
    animation: breathe 1.6s ease-in-out infinite;
  }
  @keyframes breathe {
    50% { opacity: 0.55; }
  }
</style>
