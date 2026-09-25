## Why

The dashboard is the stage for the talk's key moment: the whole room flipping from polling to streaming.
Right now it reads like an admin panel. Every card has the same weight, the prize rotates through six
joke lots, and the room-load panel shows only the current number, so you can't see the flip happen. A single
memorable prize, the Golden Goose, plus a polished participant view and a dedicated projector view
turns the live auction into something the room remembers. The load chart makes the teaching point visible.

## What Changes

- **BREAKING (catalog):** the six-lot rotation is replaced by one lot, `golden-goose` ("The Golden Goose"),
  auctioned again and again in **rounds**. After each round closes and the intermission ends, the goose reopens at its
  starting price.
- Proto: `Auction` gains `round` (which round is on the block), `lot_duration_ms` (how long a round
  opens for, for the countdown ring) and `winners` (previous sales, newest first, capped). These are new
  field numbers only, so older clients keep working. Run `buf lint` / `buf generate` and commit `gen/`.
- Auction server: counts rounds, records a winner when a round closes sold, and exposes both in every
  snapshot. The rehearsal bots reset their budget per round instead of per lot id, because the lot id never changes now.
- BFF: `AuctionDto` carries `Round`, `LotDurationMs` and `Winners`. `bff/AuctionEndpoints.cs` and
  `workshop/solution/AuctionEndpoints.cs` stay in sync, and `PlaceBid` stays unimplemented in `bff/`.
- Frontend, participant view (default): a redesigned, light-themed dashboard. It has a hero lot card with a golden
  duck artwork, rolling price digits, a countdown ring, a SOLD moment, golden eggs that drop for every bid
  received, a bid form with an outbid toast, a bid feed with initials avatars, a winners wall, the
  per-browser wire panel, and a room-load chart over time.
- Frontend, stage view (`?view=stage`): a 1920×1080 projector layout with no bid form, large type, a room-load chart
  as the main element, a full-screen SOLD moment, and the transport toggle kept so the presenter can flip it.
- **The reveal:** the goose is an inside joke, and finishing the workshop exercise unlocks it. Until this
  browser's BFF implements `PlaceBid`, the lot is a veiled "mystery lot": a cloth instead of the artwork, no
  name or description, no golden eggs, and no goose wording. Once `PlaceBid` works, the cloth lifts and the
  Golden Goose appears, with no page reload. The raw JSON still says `golden-goose`, so a curious dev can find it early.
- Golden duck SVG artwork, recoloured from an open-licensed duck emoji SVG, with gold gradients, a highlight
  and a shimmer. Also used as the favicon. Gold is used for artwork and decoration only, never for text.
- Fonts bundled offline via `@fontsource` (display serif, Nunito Sans, Fira Code), so nothing is fetched from
  the network at the venue.
- Skeleton loading state, and `prefers-reduced-motion` support for every animation.
- Slides and workshop text updated from the rubber duck to the golden goose.

## Capabilities

### New Capabilities
- `golden-goose-auction`: single-prize auction run in rounds, including round numbering, the winners
  record, and how both reach the browser through the proto and BFF.
- `auction-dashboard`: the participant view: hero lot, countdown, SOLD moment, golden eggs, bidding
  with outbid feedback, bid feed, winners wall, wire panel, room-load chart, loading and motion rules.
- `stage-view`: the projector view selected with `?view=stage`.
- `golden-duck-artwork`: the golden duck SVG, its lookup by lot id with emoji fallback, favicon, and
  the gold-palette rules.

### Modified Capabilities
<!-- none: openspec/specs/ is empty -->

## Impact

- `proto/auction/v1/auction.proto`, `auction-server/gen/` (regenerated)
- `auction-server/internal/auction/{lots,house,bots,house_test}.go`
- `bff/Dtos.cs` (`bff/AuctionEndpoints.cs` and the workshop solution only if their shared shape changes)
- `frontend/`: nearly every component rewritten or added, `index.html`, `package.json`
  (`@fontsource` packages; no chart library), `src/styles/main.css`, the theme plugin
- `slides/slides.md` ("Snipe the rubber duck", `LotId = "rubber-duck"`), `workshop/README.md` (grpcurl
  example), `workshop/PRESENTER.md` (the stage view URL, "lot" → "round" wording)
- Third-party asset licence notice for the duck SVG
- Participants must pull and rebuild before the day. Nothing changes in the workshop exercise itself.
