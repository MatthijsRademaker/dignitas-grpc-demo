<script setup lang="ts">
  import { Button, Input, NumberField } from '@vuetify/v0'
  import { computed, onScopeDispose, shallowRef, watch } from 'vue'
  import { minimumBid, placeBid, roundKey, type Auction } from '../api'

  const props = defineProps<{ auction: Auction }>()
  const emit = defineEmits<{ placed: [auction: Auction] }>()
  const bidder = defineModel<string>('bidder', { required: true })

  type Feedback = { kind: 'success' | 'error' | 'workshop', text: string }
  type Outbid = { bidder: string, amount: number }

  const minimum = computed(() => minimumBid(props.auction))
  const amount = shallowRef(minimum.value)
  const sending = shallowRef(false)
  const feedback = shallowRef<Feedback | null>(null)
  const outbid = shallowRef<Outbid | null>(null)
  const canBid = computed(() => props.auction.status === 'open' && bidder.value.trim() !== '' && !sending.value)

  // Someone outbid you, or a new round opened: never leave a stale amount in the field.
  watch(minimum, min => {
    if (amount.value < min) amount.value = min
  })
  watch(() => roundKey(props.auction), () => {
    amount.value = minimum.value
    feedback.value = null
    outbid.value = null
  })

  let hideOutbid: ReturnType<typeof setTimeout> | undefined
  onScopeDispose(() => clearTimeout(hideOutbid))

  // You were leading, and now someone else is, at a higher price, in the same round. The price
  // check skips a poll that was already in flight when your own bid came back.
  watch(() => props.auction, (next, previous) => {
    const me = bidder.value.trim()
    const was = previous?.highestBid
    const now = next.highestBid
    if (!me || !previous || !was || !now || roundKey(previous) !== roundKey(next)) return
    if (was.bidder === me && now.bidder !== me && now.amount > was.amount) {
      outbid.value = { bidder: now.bidder, amount: now.amount }
      feedback.value = null
      clearTimeout(hideOutbid)
      hideOutbid = setTimeout(() => (outbid.value = null), 6000)
    }
  })

  async function bid (value: number) {
    sending.value = true
    outbid.value = null
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
  <section class="flex flex-col gap-4 rounded-3xl border border-border bg-surface p-6">
    <h3 class="font-display text-2xl font-semibold">Place a bid</h3>

    <Input.Root v-model="bidder" class="flex flex-col gap-1" label="Your name">
      <label class="text-sm font-bold">Your name</label>
      <Input.Control
        class="rounded-xl border border-border bg-background px-3 py-2.5 outline-none focus:border-primary"
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
      <span class="text-sm font-bold">Amount <span class="font-normal text-muted">(at least €{{ minimum }})</span></span>
      <div class="flex overflow-hidden rounded-xl border border-border bg-background focus-within:border-primary">
        <NumberField.Decrement class="w-12 cursor-pointer text-xl font-bold hover:bg-primary/10">−</NumberField.Decrement>
        <NumberField.Control class="w-full bg-transparent px-3 py-2.5 text-center font-display text-2xl font-semibold tabular-nums outline-none" />
        <NumberField.Increment class="w-12 cursor-pointer text-xl font-bold hover:bg-primary/10">+</NumberField.Increment>
      </div>
    </NumberField.Root>

    <div class="flex flex-col gap-2">
      <Button.Root
        class="cursor-pointer rounded-xl bg-primary px-4 py-3.5 text-lg font-extrabold text-on-primary transition hover:brightness-110 data-[disabled]:cursor-not-allowed data-[disabled]:opacity-40"
        :disabled="!canBid"
        @click="bid(amount)"
      >
        Bid €{{ amount }}
      </Button.Root>
      <Button.Root
        class="cursor-pointer rounded-xl border-2 border-primary/30 px-4 py-2.5 font-bold text-primary transition hover:border-primary hover:bg-primary/5 data-[disabled]:cursor-not-allowed data-[disabled]:opacity-40"
        :disabled="!canBid"
        @click="bid(minimum)"
      >
        Quick bid €{{ minimum }}
      </Button.Root>
    </div>

    <p v-if="!bidder.trim()" class="text-sm text-muted">Enter your name to bid.</p>

    <p
      v-if="feedback"
      class="rounded-xl px-3 py-2 text-sm"
      :class="{
        'bg-success/10 text-success': feedback.kind === 'success',
        'bg-error/10 text-error': feedback.kind === 'error',
        'bg-warning/10 text-warning': feedback.kind === 'workshop',
      }"
      role="status"
    >
      <template v-if="feedback.kind === 'workshop'">🛠 <strong>Workshop time:</strong> </template>{{ feedback.text }}
    </p>

    <!-- The outbid toast. Fixed, but still inside #app, where the theme variables live. -->
    <Transition name="toast">
      <div
        v-if="outbid"
        aria-live="assertive"
        class="fixed inset-x-4 bottom-4 z-40 flex items-center gap-4 rounded-2xl bg-on-surface py-3 pr-3 pl-5 text-surface shadow-2xl sm:inset-x-auto sm:right-6 sm:bottom-6"
        role="alert"
      >
        <p class="flex-1"><strong>{{ outbid.bidder }}</strong> outbid you: €{{ outbid.amount }}</p>
        <Button.Root
          class="cursor-pointer rounded-xl bg-primary px-4 py-2 font-extrabold text-on-primary hover:brightness-110 data-[disabled]:opacity-40"
          :disabled="!canBid"
          @click="bid(minimum)"
        >
          Bid €{{ minimum }}
        </Button.Root>
        <button aria-label="Dismiss" class="cursor-pointer rounded-lg px-2 py-1 text-surface/70 hover:text-surface" type="button" @click="outbid = null">✕</button>
      </div>
    </Transition>
  </section>
</template>

<style scoped>
  .toast-enter-active, .toast-leave-active { transition: opacity 0.25s ease, transform 0.25s ease; }
  .toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(1rem); }
</style>
