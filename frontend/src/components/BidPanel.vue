<script setup lang="ts">
  import { Button, Input, NumberField } from '@vuetify/v0'
  import { computed, shallowRef, watch } from 'vue'
  import { minimumBid, placeBid, type Auction } from '../api'

  const props = defineProps<{ auction: Auction }>()
  const emit = defineEmits<{ placed: [auction: Auction] }>()
  const bidder = defineModel<string>('bidder', { required: true })

  type Feedback = { kind: 'success' | 'error' | 'workshop', text: string }

  const minimum = computed(() => minimumBid(props.auction))
  const amount = shallowRef(minimum.value)
  const sending = shallowRef(false)
  const feedback = shallowRef<Feedback | null>(null)
  const canBid = computed(() => props.auction.status === 'open' && bidder.value.trim() !== '' && !sending.value)

  // Someone outbid you, or a new lot opened: never leave a stale amount in the field.
  watch(minimum, min => {
    if (amount.value < min) amount.value = min
  })
  watch(() => props.auction.lot.id, () => {
    amount.value = minimum.value
    feedback.value = null
  })

  async function bid (value: number) {
    sending.value = true
    try {
      const result = await placeBid(props.auction.lot.id, bidder.value.trim(), value)
      if (result.ok) {
        emit('placed', result.value)
        feedback.value = { kind: 'success', text: `You are the highest bidder at €${value}.` }
      } else if (result.problem.status === 501) {
        feedback.value = { kind: 'workshop', text: result.problem.detail ?? result.problem.title }
      } else {
        feedback.value = { kind: 'error', text: result.problem.detail ?? result.problem.title }
      }
    } catch {
      feedback.value = { kind: 'error', text: 'Cannot reach the BFF.' }
    } finally {
      sending.value = false
    }
  }
</script>

<template>
  <section class="flex flex-col gap-4 rounded-2xl border border-border bg-surface p-6 shadow-sm">
    <h3 class="text-xs font-bold tracking-widest text-muted uppercase">Place a bid</h3>

    <Input.Root v-model="bidder" class="flex flex-col gap-1" label="Your name">
      <label class="text-sm font-bold">Your name</label>
      <Input.Control
        class="rounded-lg border border-border bg-background px-3 py-2 outline-none focus:border-primary"
        placeholder="Shown to the whole room"
      />
    </Input.Root>

    <NumberField.Root
      v-model="amount"
      class="flex flex-col gap-1"
      label="Amount"
      :min="minimum"
      :step="auction.lot.minIncrement"
    >
      <span class="text-sm font-bold">Amount (minimum €{{ minimum }})</span>
      <div class="flex overflow-hidden rounded-lg border border-border bg-background focus-within:border-primary">
        <NumberField.Decrement class="w-11 cursor-pointer text-xl font-bold hover:bg-primary/10">−</NumberField.Decrement>
        <NumberField.Control class="w-full bg-transparent px-3 py-2 text-center text-lg font-bold tabular-nums outline-none" />
        <NumberField.Increment class="w-11 cursor-pointer text-xl font-bold hover:bg-primary/10">+</NumberField.Increment>
      </div>
    </NumberField.Root>

    <div class="grid grid-cols-2 gap-2">
      <Button.Root
        class="cursor-pointer rounded-lg bg-primary px-4 py-3 font-extrabold text-on-primary transition hover:brightness-110 data-[disabled]:cursor-not-allowed data-[disabled]:opacity-40"
        :disabled="!canBid"
        @click="bid(amount)"
      >
        Bid €{{ amount }}
      </Button.Root>
      <Button.Root
        class="cursor-pointer rounded-lg border-2 border-primary px-4 py-3 font-extrabold text-primary transition hover:bg-primary/10 data-[disabled]:cursor-not-allowed data-[disabled]:opacity-40"
        :disabled="!canBid"
        @click="bid(minimum)"
      >
        Quick bid €{{ minimum }}
      </Button.Root>
    </div>

    <p
      v-if="feedback"
      class="rounded-lg px-3 py-2 text-sm"
      :class="{
        'bg-success/10 text-success': feedback.kind === 'success',
        'bg-error/10 text-error': feedback.kind === 'error',
        'bg-warning/10 text-warning': feedback.kind === 'workshop',
      }"
    >
      <template v-if="feedback.kind === 'workshop'">🛠 <strong>Workshop time:</strong> </template>{{ feedback.text }}
    </p>
  </section>
</template>
