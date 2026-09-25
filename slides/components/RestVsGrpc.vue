<script setup lang="ts">
import { useSlideContext } from "@slidev/client";

const { $clicks } = useSlideContext();
</script>

<template>
  <div class="split">
    <section class="side" :class="{ dim: $clicks >= 2 }">
      <header>With REST <span>pick a workaround</span></header>

      <h4>Browser → BFF</h4>
      <ul>
        <li><b>Polling</b>: noisy, and stale between polls</li>
        <li><b>Long polling</b>: less chatter, messy timeouts</li>
        <li><b>SSE / WebSockets</b>: fine, but outside your API contract</li>
      </ul>

      <template v-if="$clicks >= 1">
        <h4>Service → service</h4>
        <ul>
          <li><b>Polling</b>: the same noise, now between servers</li>
          <li><b>Webhooks</b>: retries, signatures, dead letters</li>
          <li><b>A broker</b>: topics, schemas, one more thing to run</li>
        </ul>
      </template>
    </section>

    <section v-show="$clicks >= 2" class="side primary">
      <header>With gRPC <span>one keyword</span></header>

      <pre class="rpc">
rpc WatchAuction(WatchAuctionRequest)
  returns (<b>stream</b> WatchAuctionResponse);</pre>

      <ul>
        <li>Server and client <strong>generated</strong> from that line</li>
        <li>Deadlines, status codes, keepalives: <strong>built in</strong></li>
        <li>At the edge we still used SSE: browsers don’t speak gRPC</li>
      </ul>
    </section>
  </div>
</template>

<style scoped>
.split {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.4rem;
}
.side {
  border: 2px solid color-mix(in srgb, currentColor 18%, transparent);
  border-radius: 0.75rem;
  padding: 1rem 1.2rem;
  transition: opacity 0.3s ease;
}
.side.dim {
  opacity: 0.6;
}
.side.primary {
  border-color: var(--se-color-primary);
  background: color-mix(in srgb, var(--se-color-primary) 6%, transparent);
}
header {
  font-size: 1.15rem;
  font-weight: 800;
}
header span {
  margin-left: 0.4rem;
  font-size: 0.85rem;
  font-weight: 600;
  opacity: 0.6;
}
h4 {
  margin: 0.8rem 0 0.3rem;
  font-size: 0.75rem;
  font-weight: 800;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  opacity: 0.6;
}
ul {
  padding-left: 1.1rem;
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  font-size: 0.92rem;
}
.rpc {
  font-family: "Fira Code", monospace;
  font-size: 0.8rem;
  line-height: 1.3;
  margin: 0.9rem 0;
  padding: 0.5rem 0.7rem;
  border-radius: 0.4rem;
  background: color-mix(in srgb, currentColor 5%, transparent);
  white-space: pre;
}
.rpc b {
  color: var(--se-color-primary);
}
</style>
