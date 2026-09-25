<template>
  <div class="split">
    <section class="side primary">
      <header>gRPC stream <span>a direct pipe</span></header>
      <div class="flow">
        <span class="box">Producer</span>
        <span class="pipe" />
        <span class="box">Consumer</span>
      </div>
      <ul>
        <li>Both sides are connected <strong>at the same time</strong></li>
        <li>Nothing in between to run or operate</li>
        <li>Not connected? You missed it</li>
        <li>Great for <strong>live conversations</strong></li>
      </ul>
    </section>

    <section class="side">
      <header>Message broker <span>a post office</span></header>
      <div class="flow">
        <span class="box">Producer</span>
        <span class="arrow">→</span>
        <span class="box broker">📬 Broker</span>
        <span class="arrow">→</span>
        <span class="box">Consumers</span>
      </div>
      <ul>
        <li>Decoupled <strong>in time</strong>: consumer can be offline</li>
        <li>Durable, retries, replay</li>
        <li>Fan-out to many independent consumers</li>
        <li>Extra infrastructure: Service Bus, Kafka, MQTT…</li>
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
.flow {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0.9rem 0;
}
.box {
  padding: 0.35rem 0.7rem;
  border-radius: 0.5rem;
  font-size: 0.85rem;
  font-weight: 700;
  border: 2px solid color-mix(in srgb, currentColor 25%, transparent);
}
.box.broker {
  border-color: currentColor;
}
.pipe {
  flex: 1;
  height: 0.55rem;
  border-radius: 1rem;
  background: var(--se-color-primary);
}
.arrow {
  font-weight: 800;
  opacity: 0.55;
}
ul {
  padding-left: 1.1rem;
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  font-size: 0.95rem;
}
</style>
