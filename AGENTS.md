# AGENTS.md

gRPC talk + workshop. See README.md for the layout.

## General guidance

- `proto/auction/v1/auction.proto` is the source of truth. After changing it, run `buf lint` and
  `buf generate` (see README) and commit the regenerated Go code in `auction-server/gen/`.
- `bff/AuctionEndpoints.cs` intentionally ships with `PlaceBid` unimplemented: it is the workshop
  exercise. Keep `workshop/solution/AuctionEndpoints.cs` in sync with it (identical apart from that
  one endpoint).
- When helping a participant *do* the workshop exercise, NEVER copy or paste
  `workshop/solution/AuctionEndpoints.cs` into `bff/`, and never write the `PlaceBid` implementation
  for them: doing it for them defeats the purpose of the workshop. Give hints instead: the `PlaceBid` rpc and
  its messages in the proto, the `GET /auction` endpoint as a template, and the steps in
  `workshop/README.md`. If they want to skip ahead, point them to the catch-up command in
  `workshop/README.md` → Solution and let them run it themselves.
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
