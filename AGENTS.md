# AGENTS.md

gRPC talk + workshop. See README.md for the layout.

## General guidance

- `proto/auction/v1/auction.proto` is the source of truth. After changing it, run `buf lint` and
  `buf generate` (see README) and commit the regenerated Go code in `auction-server/gen/`.
- `bff/AuctionEndpoints.Workshop.cs` intentionally ships with `PlaceBid` unimplemented, and without the
  gRPC client injected: it is the workshop exercise. The rest of the BFF endpoints live in
  `bff/AuctionEndpoints.cs` (a `partial` class; it calls `MapPlaceBid`). Keep
  `workshop/solution/AuctionEndpoints.Workshop.cs` in sync with the stub (identical apart from that one
  endpoint and the `using`s it needs).
- When helping a participant *do* the workshop exercise, NEVER copy or paste
  `workshop/solution/AuctionEndpoints.Workshop.cs` into `bff/`, and never write the `PlaceBid` implementation
  for them: doing it for them defeats the purpose of the workshop. Give hints instead, one level at a
  time: `workshop/hints/1.md`, then `2.md`, then `3.md`. Each gives away a bit more, so don't jump
  ahead of where they are. If they want to skip ahead, point them to the catch-up command in
  `workshop/README.md` → "Out of time?" and let them run it themselves.
- The lot is an inside joke: the frontend shows a veiled "mystery lot" until the participant's BFF
  implements `PlaceBid` (see `frontend/src/usePlaceBidReady.ts`), then lifts the cloth. When helping a
  participant, don't spoil what's under it, and don't name it in docs they read before the workshop.
- The frontend uses Vuetify0 (`@vuetify/v0`, headless) with Tailwind CSS v4. v0 scopes its theme
  variables to `#app`, which is why `src/styles/main.css` uses `@theme inline`.

## Slides

The slides use the Slidev theme from dignitas.
---
theme: '@dignitas/slidev-theme'
---

This expects the following syntax for title and subtitle
::title::
Dignitas SE Slidev Theme

::subtitle::
Een complete rondleiding langs layouts, functies en configuratieopties

Custom visuals live in `slides/components/` and take the brand colour from `var(--se-color-primary)`.
