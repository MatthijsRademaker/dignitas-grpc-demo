## Context

The app is `browser ─HTTP/JSON, SSE→ BFF (.NET) ─gRPC→ auction-server (Go)`. It runs on two kinds of
screen at once: every participant's laptop, used for bidding, and the presenter's projector, used for the
polling-to-streaming reveal described in `workshop/PRESENTER.md`. The frontend is Vue 3 with
Vuetify0 (headless) and Tailwind v4. Theme colours come from v0 variables scoped to `#app`, hence
`@theme inline` in `src/styles/main.css`.

Relevant current behaviour:

- `House.Tick` opens `(lotIndex + 1) % len(Lots)` after the intermission, so a one-lot catalog
  already reopens the same lot. There is no notion of a round, and past results are dropped in `openLot`.
- The browser only gets `remaining_ms`, not the lot's total duration.
- `useAuctionFeed` resets all its stats when the transport or poll interval changes.
- `ServerLoad` is averaged over a 10s window (`loadWindow`) and pushed to streams at most every 2s.
- `bots.go` resets a bot's budget when `lot.id` changes.
- Venue Wi-Fi is unreliable. Today `Nunito Sans` and `Fira Code` only render if they happen to be installed.

## Goals / Non-Goals

**Goals:**
- One prize, the Golden Goose, auctioned in rounds, with the round number and a winners wall shared by
  every screen.
- A participant view and a stage view that feel like one polished product in a light theme.
- Make the flip visible: a room-load chart whose history survives the transport switch, and golden eggs
  that show whether bids arrive in batches (polling) or one at a time (streaming).
- Everything works offline once the Docker images are built.

**Non-Goals:**
- A dark theme, or a theme switcher.
- Persisting winners across server restarts.
- A chart or animation library: the charts are small hand-written SVG and CSS.
- Changing the workshop exercise, or any existing RPC or its semantics.
- Showing the presenter's LAN IP on the stage view (it stays on the "Get connected" slide).

## Decisions

### D1. Rounds live on the server, as new proto fields

```proto
message Auction {
  ...
  int32 round = 9;             // 1 for the first lot opened, +1 every time a lot opens
  int64 lot_duration_ms = 10;  // how long a lot opens for; late bids can extend it past this
  repeated Sale winners = 11;  // earlier sold rounds, newest first, capped
}

message Sale {
  int32 round = 1;
  string bidder = 2;
  int64 amount = 3;
  google.protobuf.Timestamp sold_at = 4;
}
```

`House` gets a `round` counter, incremented in `openLot`, and a `winners` slice capped at 10. A winner is
added when `Tick` moves a round from open to sold, so it appears in the `LOT_CLOSED` event. Unsold rounds
are not recorded. The winners stay in memory and are lost on restart.

*Alternatives:* let the browser count rounds and remember winners. Rejected because a laptop that joins late,
or a projector that refreshes, would show a different round number and an empty wall than everyone else.
The server owning the story is also a good proto-evolution example: new fields with new numbers, and
nothing breaks.

*Why `lot_duration_ms`:* the countdown ring needs to know what a full ring means. Guessing it from the largest
`remaining_ms` seen shows a full ring after a mid-round refresh, which looks wrong on the projector. Ring
fill = `remaining / max(lot_duration_ms, remaining)`, so a snipe extension that goes past the
configured duration just shows a full ring.

### D2. The catalog becomes one lot; bots key on the round

`Catalog` holds just `{Id: "golden-goose", Emoji: "🦆", Title: "The Golden Goose", ...}`. The Emoji stays a
real duck, so grpcurl output and the text fallback still make sense. Bots re-randomise their budget when
`auction.GetRound()` changes, not when the lot id does.

### D3. BFF: extend `AuctionDto`, leave the endpoints alone

`AuctionDto` gains `int Round`, `long LotDurationMs` and `IReadOnlyList<SaleDto> Winners`.
`SaleDto.From(Sale)` follows the `BidDto` pattern. The endpoints already return `AuctionDto.From(...)`, so
`bff/AuctionEndpoints.cs` and `workshop/solution/AuctionEndpoints.cs` need no change and stay in sync
without effort.

### D4. One app, two layouts, chosen by the URL

`?view=stage` selects `StageView.vue`. Otherwise `ParticipantView.vue` renders. `?transport=stream` keeps
working in both. `App.vue` keeps the shared state (the feed, the transport, the room history) and passes
it down, so both layouts use the same components (`LotHero`, `Countdown`, `PriceTicker`, `RoomChart`,
`WinnersWall`, `EggNest`) in different sizes.

```
 App.vue ── useAuctionFeed ─┬─ useRoomHistory (survives transport switches)
                            ├─ useEggs (newly seen bids → egg drops)
                            └─ view = participant | stage
      ParticipantView                     StageView (1920×1080)
 ┌───────────────────────┬────────┐   ┌───────────────┬───────────────┐
 │ LotHero + EggNest     │BidPanel│   │ GoldenDuck    │ PriceTicker   │
 │ PriceTicker Countdown │(sticky)│   │ + EggNest     │ Countdown     │
 ├───────────────────────┼────────┤   │               │ top 3 bids    │
 │ BidFeed               │ Wire   │   ├───────────────┴───────────────┤
 ├───────────────────────┴────────┤   │ RoomChart (large) · Wire stat │
 │ RoomChart · WinnersWall        │   │ WinnersWall strip             │
 └────────────────────────────────┘   └───────────────────────────────┘
```

A route library is not needed for two layouts.

### D5. The room-load history is separate from per-browser stats

A new `useRoomHistory(auction)` keeps a ring buffer of `{at, pollsPerSecond, watchers}` for the last 90s.
It samples every time an auction arrives, whatever the transport, and it is never cleared by a
transport or interval change. `RoomChart` draws two small multiples, **calls/s** and **open
streams**, on one shared 90s time axis, as step lines that carry the last value forward to "now". It does
not use a dual-axis chart, because the two series have different units. Each panel labels its current
value directly.

The 10s server average turns the flip into a slope of about 10s rather than a sudden drop. That reads fine and needs
no change.

### D6. Golden eggs come from bids seen for the first time

`useAuctionFeed.apply` already knows which bids are new to this page (`seenAfter`). It gets an
`onNewBids(bids)` hook, and `useEggs` turns each call into one batch of egg drops. A poll that brings three
bids drops three eggs together, while a stream drops them one at a time. Bids that were already there when the page loaded appear
in the nest without animation. The nest empties when a new round opens, and shows at most about 24 eggs, with
"+N" beyond that. Your own bids also drop an egg, as soon as the `PlaceBid` response comes back.

### D7. Golden duck artwork is a recoloured open-licensed SVG

The source is the Noto Emoji duck (`emoji_u1f986.svg`, Apache 2.0). Its flat fills are replaced by an SVG
`linearGradient` gold ramp (roughly `#7a5a00 → #c9971c → #f6d77a`), with a specular highlight path and a
shimmer: an animated gradient band under a mask of the duck shape. It lives in a `GoldenDuck.vue`
component so the gradients and animation can be themed and paused. `LotArt.vue` maps
`lot.id → component` and otherwise renders `lot.emoji`. The licence goes in
`frontend/THIRD_PARTY_NOTICES.md`, with a comment in the component. `public/favicon.svg` is a static,
unanimated copy.

*Alternative:* CSS `filter` on the 🦆 glyph. Rejected because every OS draws a different duck, so the gold
would look different on each laptop. Drawing a duck from scratch was also rejected: a professionally drawn
shape looks better at projector size.

### D8. Palette and type

The light theme keeps the purple `primary` (`#733FCB`). New tokens: `gold-50…900` for artwork and
decoration only, and `bronze` (a gold-700 around `#7a5a00`) as the one gold that may appear as text, at ≥ 4.5:1 on
white. Type: **Fraunces** (variable, display serif) for the lot title, price and SOLD moment, **Nunito
Sans** for the UI, and **Fira Code** for wire and code details. All three come from `@fontsource`
packages and are bundled by Vite, so no network is needed at runtime.

### D9. Motion is CSS, and everything obeys `prefers-reduced-motion`

The rolling price digits use one vertically translated digit strip per column. The countdown ring is an SVG circle
animated through `stroke-dashoffset`. It turns red and pulses in the last 10s. The SOLD moment is a gavel slam
and a gold burst over the hero card (participant view) or the whole screen (stage view), lasting for the
intermission. Eggs drop with a keyframe fall and a bounce. Under `prefers-reduced-motion: reduce`, the shimmer,
pulse and burst are off, digits change instantly, and eggs appear without falling.

### D10. Small polish items

- **Outbid toast:** shown when `highestBid` changes from you (by name) to someone else in the same
  round. It shows the new amount and a one-click "Bid €min" action, and disappears after a few seconds.
- **Initials avatars** in the bid feed. The colour is a stable hash of the name over a fixed, accessible set.
- **Skeleton** versions of the hero, feed and chart replace "Connecting to the auction…", and the
  error text still shows when the connection fails.

### D11. The goose is unlocked by PlaceBid

Each browser asks its own BFF whether `PlaceBid` works, with a probe that can never place a bid:
`POST /api/bids` with an empty bidder. The workshop stub answers 501. A finished BFF passes the call on,
and the auction server rejects the empty bidder with `INVALID_ARGUMENT`, which `ToProblem()` turns into
400 with the title `InvalidArgument`. Only that answer reveals the goose, because it proves the call reached the server. A
successful real bid reveals it too, so an implementation without error handling (500 on the probe)
still gets there with its first bid.

The probe runs on page load and every 5s while the lot is veiled, and stops once it is revealed. `--watch` restarts the BFF on
save, so the reveal happens a few seconds after the participant's fix, with no reload. The last result is
kept in `localStorage`, so a finished participant doesn't see the cloth flash on every reload. The probe on
load still corrects it (for example after `git checkout bff/AuctionEndpoints.cs`).

While veiled: `LotArt` shows a draped cloth (`VeiledLot.vue`), the title is "A mystery lot", and the
description tells you to implement `PlaceBid`. The egg nest is hidden. Goose wording (the winners heading, the unsold copy) is
neutral, and the header and favicon show a gavel. The reveal lifts the cloth off, and the duck scales in
and starts to shimmer. With reduced motion, it simply swaps. The stage view uses the same rule, so the presenter's projector reveals once the solution is
applied. The auction server and the JSON are unchanged. The secret is in the frontend only, and a dev who looks at
DevTools or grpcurl finds it early. That's fine for an inside joke.

*Alternative:* a new RPC or flag that reports whether the BFF is implemented. Rejected because it would change the contract and the
exercise. The probe needs nothing new.

## Risks / Trade-offs

- [Recoloured Noto duck looks muddy on a projector] → Tune the gradient on the real projector during
  rehearsal. Keep a darker outline stroke for definition. The favicon uses a simplified copy.
- [Participants on an old image miss the new fields] → The frontend treats a missing `round` or `winners` as
  "hide that element", and the BFF maps an unset `Sale` list to empty. Participants must pull and rebuild anyway
  for the goose.
- [Stage browser's own polling adds to the room load it displays] → It's one client among many.
  PRESENTER.md suggests opening the stage in streaming mode once the reveal is done.
- [Fraunces plus two other font families make the bundle bigger] → Use Latin subsets only and variable fonts. It's
  still well under 300 KB.
- [Hand-rolled SVG charts need extra care with accessibility] → Each panel has an `aria-label` with the current value and a
  visible text value, and colour is never the only signal.
- [Veiled browsers probe every 5s, which adds `PlaceBid … code=InvalidArgument` lines to the server log] → They stop
  at the reveal, and the first of those lines is a nice signal in the presenter's log that a participant just finished.
- [Winners are lost on server restart] → Acceptable for a talk. The wall simply starts empty again.

## Migration Plan

1. Change the proto, run `buf lint` and `buf generate`, and commit `auction-server/gen/`.
2. Update the server and tests, then the BFF DTOs, then the frontend.
3. Update the slides, `workshop/README.md` and `PRESENTER.md`.
4. Rehearse with `AUCTION_BOTS=3 docker compose --profile presenter up --build` on the actual projector.

Rollback is `git revert`. Nothing is stored, so there's no data migration.

## Open Questions

- The exact goose description copy (e.g. "Lays one golden egg for every server push.").
- Should the stage view show a "Round N" badge larger than the lot title, or next to it?
