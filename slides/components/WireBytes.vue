<script setup lang="ts">
import { useSlideContext } from '@slidev/client'

// Verified with proto.Marshal(&auctionv1.Bid{Bidder: "Ada", Amount: 120}).
const json = '{"bidder":"Ada","amount":120}'
const groups = [
  { bytes: ['0a'], note: 'field 1 · length-delimited' },
  { bytes: ['03'], note: 'length: 3' },
  { bytes: ['41', '64', '61'], note: '"Ada"' },
  { bytes: ['10'], note: 'field 2 · varint' },
  { bytes: ['78'], note: '120' },
]

const { $clicks } = useSlideContext()
</script>

<template>
  <div class="wire">
    <div class="row">
      <div class="label">JSON <span class="count">{{ json.length }} bytes</span></div>
      <code class="json">{{ json }}</code>
    </div>

    <div v-show="$clicks >= 1" class="row">
      <div class="label primary">Protobuf <span class="count">7 bytes</span></div>
      <div class="bytes">
        <div v-for="(group, i) in groups" :key="i" class="group">
          <div class="hex">
            <span v-for="b in group.bytes" :key="b" class="byte">{{ b }}</span>
          </div>
          <div v-show="$clicks >= 2" class="note">{{ group.note }}</div>
        </div>
      </div>
    </div>

    <p v-show="$clicks >= 2" class="punch">
      The field <em>names</em> never travel. Only the <span class="text-primary">field numbers</span> from the .proto do.
    </p>
  </div>
</template>

<style scoped>
.wire {
  display: flex;
  flex-direction: column;
  gap: 1.6rem;
}
.row {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.label {
  font-size: 0.85rem;
  font-weight: 800;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  opacity: 0.7;
}
.label.primary {
  color: var(--se-color-primary);
  opacity: 1;
}
.count {
  margin-left: 0.6rem;
  font-size: 1.1rem;
  letter-spacing: 0;
  text-transform: none;
}
.json {
  font-family: 'Fira Code', monospace;
  font-size: 1.6rem;
  padding: 0.7rem 1rem;
  border-radius: 0.6rem;
  background: color-mix(in srgb, currentColor 6%, transparent);
  width: fit-content;
}
.bytes {
  display: flex;
  gap: 0.9rem;
  align-items: flex-start;
}
.group {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.45rem;
}
.hex {
  display: flex;
  gap: 0.3rem;
  padding: 0.35rem;
  border-radius: 0.6rem;
  background: color-mix(in srgb, var(--se-color-primary) 12%, transparent);
}
.byte {
  font-family: 'Fira Code', monospace;
  font-size: 1.6rem;
  font-weight: 700;
  padding: 0.2rem 0.5rem;
  border-radius: 0.35rem;
  color: #fff;
  background: var(--se-color-primary);
}
.note {
  font-size: 0.8rem;
  font-weight: 600;
  opacity: 0.7;
  max-width: 9rem;
  text-align: center;
}
.punch {
  font-size: 1.2rem;
  font-weight: 700;
}
</style>
