<script setup lang="ts">
import { useSlideContext } from '@slidev/client'

const { $clicks } = useSlideContext()

// requests/responses: how many messages travel each way.
const shapes = [
  { name: 'Unary', gist: 'one request, one response', requests: 1, responses: 1, rpc: 'PlaceBid(PlaceBidRequest)\n  returns (PlaceBidResponse)', example: 'Place a bid. You built this.' },
  { name: 'Server streaming', gist: 'one request, many responses', requests: 1, responses: 4, rpc: 'WatchAuction(WatchAuctionRequest)\n  returns (stream WatchAuctionResponse)', example: 'Watch the auction. You just saw it.' },
  { name: 'Client streaming', gist: 'many requests, one response', requests: 4, responses: 1, rpc: 'UploadLotPhoto(stream Chunk)\n  returns (UploadResult)', example: 'Upload a lot photo in chunks.' },
  { name: 'Bidirectional', gist: 'both sides talk, whenever', requests: 3, responses: 3, rpc: 'Auctioneer(stream Command)\n  returns (stream Event)', example: 'Auctioneer console: bids in, “going once…” out.' },
]
</script>

<template>
  <div class="grid">
    <section v-for="(shape, i) in shapes" :key="shape.name" v-show="$clicks >= i" class="card"
      :class="{ current: $clicks === i }">
      <header>
        <span class="name">{{ shape.name }}</span>
        <span class="gist">{{ shape.gist }}</span>
      </header>
      <div class="wire">
        <span class="end">client</span>
        <div class="lanes">
          <div class="lane"><span v-for="n in shape.requests" :key="n" class="msg">→</span></div>
          <div class="lane back"><span v-for="n in shape.responses" :key="n" class="msg">←</span></div>
        </div>
        <span class="end">server</span>
      </div>
      <pre class="rpc">rpc {{ shape.rpc }}</pre>
      <p class="example">{{ shape.example }}</p>
    </section>
  </div>
</template>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.9rem;
}
.card {
  border: 2px solid color-mix(in srgb, currentColor 15%, transparent);
  border-radius: 0.75rem;
  padding: 0.7rem 1rem;
  transition: all 0.3s ease;
}
.card.current {
  border-color: var(--se-color-primary);
  background: color-mix(in srgb, var(--se-color-primary) 7%, transparent);
}
header {
  display: flex;
  align-items: baseline;
  gap: 0.6rem;
}
.name {
  font-weight: 800;
  font-size: 1.05rem;
}
.gist {
  font-size: 0.8rem;
  opacity: 0.6;
}
.wire {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0.45rem 0;
}
.end {
  font-size: 0.7rem;
  font-weight: 700;
  opacity: 0.55;
  text-transform: uppercase;
}
.lanes {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}
.lane {
  display: flex;
  gap: 0.25rem;
  padding: 0.1rem 0.3rem;
  border-radius: 0.4rem;
  background: color-mix(in srgb, currentColor 5%, transparent);
}
.lane.back {
  justify-content: flex-end;
}
.msg {
  font-size: 0.75rem;
  font-weight: 800;
  line-height: 1.1rem;
  padding: 0 0.45rem;
  border-radius: 0.3rem;
  color: #fff;
  background: var(--se-color-primary);
}
.lane.back .msg {
  background: color-mix(in srgb, currentColor 55%, transparent);
}
.rpc {
  font-family: 'Fira Code', monospace;
  font-size: 0.72rem;
  line-height: 1.25;
  margin: 0;
  padding: 0.35rem 0.5rem;
  border-radius: 0.4rem;
  background: color-mix(in srgb, currentColor 5%, transparent);
  white-space: pre;
}
.example {
  margin: 0.4rem 0 0;
  font-size: 0.85rem;
}
</style>
