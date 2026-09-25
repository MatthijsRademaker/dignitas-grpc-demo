<script setup lang="ts">
import { useSlideContext } from "@slidev/client";

const { $clicks } = useSlideContext();
</script>

<template>
  <div class="w-full flex flex-col items-center">
    <svg
      viewBox="0 0 900 330"
      class="w-full max-w-[860px]"
      role="img"
      aria-label="The browser speaks HTTP/JSON to a BFF; behind the BFF, services speak gRPC to each other."
    >
      <defs>
        <marker
          id="ab-arrow"
          viewBox="0 0 10 10"
          refX="9"
          refY="5"
          markerWidth="7"
          markerHeight="7"
          orient="auto-start-reverse"
        >
          <path d="M 0 0 L 10 5 L 0 10 z" class="head" />
        </marker>
      </defs>

      <!-- zones -->
      <rect
        x="10"
        y="16"
        width="250"
        height="298"
        rx="14"
        class="zone"
        :class="{ lit: $clicks >= 3 }"
      />
      <text x="30" y="44" class="zone-label">THE EDGE</text>
      <g v-show="$clicks >= 1">
        <rect
          x="530"
          y="16"
          width="360"
          height="298"
          rx="14"
          class="zone"
          :class="{ lit: $clicks >= 3 }"
        />
        <text x="550" y="44" class="zone-label">SERVICE TO SERVICE</text>
      </g>

      <!-- browser -> BFF: always -->
      <rect x="45" y="135" width="170" height="60" rx="10" class="node" />
      <text x="130" y="171" text-anchor="middle" class="node-text">
        🖥️ Browser
      </text>
      <line
        x1="215"
        y1="165"
        x2="335"
        y2="165"
        class="link http"
        marker-end="url(#ab-arrow)"
      />
      <text x="275" y="150" text-anchor="middle" class="link-label">
        HTTP / JSON
      </text>
      <text x="275" y="190" text-anchor="middle" class="link-sub">
        SSE for live data
      </text>

      <rect x="340" y="125" width="150" height="80" rx="10" class="node bff" />
      <text x="415" y="162" text-anchor="middle" class="node-text on-primary">
        BFF
      </text>
      <text x="415" y="184" text-anchor="middle" class="node-sub on-primary">
        translates
      </text>

      <!-- BFF -> services -->
      <g v-show="$clicks >= 1">
        <line
          x1="490"
          y1="150"
          x2="598"
          y2="80"
          class="link grpc"
          marker-end="url(#ab-arrow)"
        />
        <line
          x1="490"
          y1="165"
          x2="598"
          y2="165"
          class="link grpc"
          marker-end="url(#ab-arrow)"
        />
        <line
          x1="490"
          y1="180"
          x2="598"
          y2="250"
          class="link grpc"
          marker-end="url(#ab-arrow)"
        />
        <text x="545" y="120" text-anchor="middle" class="link-label primary">
          gRPC
        </text>

        <rect x="600" y="52" width="200" height="56" rx="10" class="node svc" />
        <text x="700" y="86" text-anchor="middle" class="node-text">
          Auction
          <tspan class="lang">Go</tspan>
        </text>
        <rect
          x="600"
          y="137"
          width="200"
          height="56"
          rx="10"
          class="node svc"
        />
        <text x="700" y="171" text-anchor="middle" class="node-text">
          Payments
          <tspan class="lang">Java</tspan>
        </text>
        <rect
          x="600"
          y="222"
          width="200"
          height="56"
          rx="10"
          class="node svc"
        />
        <text x="700" y="256" text-anchor="middle" class="node-text">
          Inventory
          <tspan class="lang">.NET</tspan>
        </text>
      </g>
    </svg>

    <p v-show="$clicks >= 3" class="punch">
      Not one protocol everywhere.
      <span class="text-primary">Choose per boundary.</span>
    </p>
  </div>
</template>

<style scoped>
.zone {
  fill: color-mix(in srgb, currentColor 3%, transparent);
  stroke: currentColor;
  stroke-opacity: 0.25;
  stroke-width: 1.5;
  stroke-dasharray: 6 6;
  transition: all 0.3s ease;
}
.zone.lit {
  stroke: var(--se-color-primary);
  stroke-opacity: 0.8;
  fill: color-mix(in srgb, var(--se-color-primary) 6%, transparent);
}
.zone-label {
  fill: currentColor;
  opacity: 0.55;
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.08em;
}
.node {
  fill: var(--slidev-code-background, #fff);
  stroke: currentColor;
  stroke-opacity: 0.3;
  stroke-width: 2;
}
.node.svc {
  stroke: var(--se-color-primary);
  stroke-opacity: 0.7;
}
.node.bff {
  fill: var(--se-color-primary);
  stroke: none;
}
.node-text {
  fill: currentColor;
  font-size: 18px;
  font-weight: 800;
}
.node-sub {
  font-size: 13px;
  font-weight: 600;
}
.on-primary {
  fill: #fff;
}
.lang {
  fill: currentColor;
  opacity: 0.55;
  font-size: 13px;
  font-weight: 700;
}
.link {
  fill: none;
  stroke-width: 2.5;
}
.link.http {
  stroke: currentColor;
  stroke-opacity: 0.55;
}
.link.grpc {
  stroke: var(--se-color-primary);
}
.head {
  fill: var(--se-color-primary);
}
.link-label {
  fill: currentColor;
  font-size: 14px;
  font-weight: 800;
}
.link-label.primary {
  fill: var(--se-color-primary);
}
.link-sub {
  fill: currentColor;
  opacity: 0.55;
  font-size: 12px;
}
.punch {
  margin-top: 0.5rem;
  font-size: 1.3rem;
  font-weight: 800;
  text-align: center;
}
</style>
