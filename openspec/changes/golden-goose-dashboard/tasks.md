## 1. Contract

- [x] 1.1 Add `Sale` message and `Auction.round` (9), `lot_duration_ms` (10), `winners` (11) to `proto/auction/v1/auction.proto`; update the service comment ("lots rotate forever" → rounds)
- [x] 1.2 Run `buf lint` and `buf generate`; commit the regenerated `auction-server/gen/`

## 2. Auction server

- [x] 2.1 Replace `Catalog` in `lots.go` with the single `golden-goose` lot (emoji `🦆`, title "The Golden Goose", description, starting price, increment)
- [x] 2.2 Add a `round` counter to `House`, incremented in `openLot`, and include `round` and `lot_duration_ms` in `snapshot()`
- [x] 2.3 Record a `Sale` in `Tick` when a round closes sold (newest first, capped at 10), and include `winners` in `snapshot()`
- [x] 2.4 Make the bots in `bots.go` reset their budget when `auction.GetRound()` changes, instead of when the lot id does
- [x] 2.5 Tests in `house_test.go`: round increments on reopen and holds during the intermission; a sold round is recorded and an unsold one is not; the cap drops the oldest; `lot_duration_ms` matches the config
- [x] 2.6 `go test ./...` passes

## 3. BFF

- [x] 3.1 Add `SaleDto` and extend `AuctionDto` in `bff/Dtos.cs` with `Round`, `LotDurationMs`, `Winners`
- [x] 3.2 Confirm `bff/AuctionEndpoints.cs` and `workshop/solution/AuctionEndpoints.cs` need no change, and are still identical apart from `PlaceBid`; `dotnet build` passes
- [x] 3.3 Check the JSON with `curl localhost:8080/api/auction` against a running server (round, lotDurationMs, winners)

## 4. Frontend foundation

- [x] 4.1 Add `@fontsource-variable/fraunces`, `@fontsource-variable/nunito-sans`, `@fontsource/fira-code` (Latin subsets), import them in `main.ts`, set `--font-display` in `main.css`
- [x] 4.2 Extend the v0 theme and `@theme inline` with the gold ramp (decoration) and the `bronze` text token; verify bronze reaches ≥ 4.5:1 on `surface` and `background`
- [x] 4.3 Add `round`, `lotDurationMs`, `Sale` and `winners` to `src/api.ts` (optional, for older servers)
- [x] 4.4 Add an `onNewBids` hook to `useAuctionFeed`, fired with the bids first seen in each `apply`, with a flag for "already there at page load"
- [x] 4.5 Add `useRoomHistory`: a 90s ring buffer of `{at, pollsPerSecond, watchers}`, sampled on every auction and never reset by transport changes
- [x] 4.6 Add `useEggs`: turns `onNewBids` batches into egg drops, clears them on a new round, caps the nest at 24 plus "+N"
- [x] 4.7 Add a `prefers-reduced-motion` base layer in `main.css` that disables the shimmer, pulse, burst, digit roll and egg fall

## 5. Golden duck artwork

- [x] 5.1 Take the Noto Emoji duck (`emoji_u1f986.svg`, Apache 2.0) and build `GoldenDuck.vue` with a gold gradient ramp, a specular highlight, a darker outline and a masked shimmer sweep
- [x] 5.2 `LotArt.vue`: `golden-goose` → `GoldenDuck`, otherwise `lot.emoji` as text
- [x] 5.3 Static `public/favicon.svg` and update `index.html`
- [x] 5.4 `frontend/THIRD_PARTY_NOTICES.md`, plus a source and licence comment in `GoldenDuck.vue`

## 6. Shared dashboard components

- [x] 6.1 `PriceTicker.vue`: rolling digits per column, with a highlight on change
- [x] 6.2 `Countdown.vue`: SVG ring filled to `remaining / max(lotDurationMs, remaining)`, `m:ss`, a closing state in the last 10s
- [x] 6.3 `EggNest.vue`: falling and bouncing golden eggs from `useEggs`, with a static nest for bids that were already there at page load
- [x] 6.4 `SoldMoment.vue`: gavel, "SOLD to X · €Y", next-round countdown; an unsold variant; a `card` / `fullscreen` size prop
- [x] 6.5 `LotHero.vue`: artwork, round badge, title in the display serif, description, price, leader, countdown, nest, and the SOLD overlay
- [x] 6.6 `RoomChart.vue`: two small multiples (calls/s, open streams) on a shared 90s axis, step lines to "now", direct value labels, `aria-label`s; follow the `dataviz` skill
- [x] 6.7 `WinnersWall.vue`: newest first, an empty invitation, hidden when `winners` is missing
- [x] 6.8 `BidFeed.vue`: initials avatars with a stable name-hash colour, and keep the delay and "you" markers
- [x] 6.9 `BidPanel.vue`: restyle, keep the 501 workshop message, add the outbid toast with a one-click "Bid €min"
- [x] 6.10 `TransportPanel.vue`: restyle to match; add a compact variant exposing "bids seen after" for the stage
- [x] 6.11 Skeleton placeholders for the hero, feed and chart, replacing "Connecting to the auction…", with the error text kept

## 7. Views

- [x] 7.1 Move the shared state (feed, transport, history, eggs) into `App.vue`, and pick the view from `?view=stage`
- [x] 7.2 `ParticipantView.vue`: header (brand, wire path, toggles), hero with a sticky bid panel, feed with the wire panel, chart with the winners wall
- [x] 7.3 `StageView.vue`: a 1920×1080 grid with no scroll; large duck and nest, price and countdown ≥ 96px, top 3 bids, a large chart, a "bids seen after" figure, a winners strip, quiet toggles, a full-screen SOLD moment; scales down to 1280×720
- [x] 7.4 `npm run build` (vue-tsc) passes

## 8. Docs and slides

- [x] 8.1 `slides/slides.md`: "Snipe the rubber duck" → the golden goose; `LotId = "rubber-duck"` → `"golden-goose"`; "lot" → "round" wording in the live-auction speaker notes
- [x] 8.2 `workshop/README.md`: grpcurl example `lot_id` → `golden-goose`; "the lot currently on the block" wording
- [x] 8.3 `workshop/PRESENTER.md`: projector on `?view=stage`, rounds instead of lots, a tip to switch the stage to streaming after the reveal
- [x] 8.4 Update the proto comment and the README table ("lots rotate forever") where needed

## 9. Verification

- [x] 9.1 Rehearse with `AUCTION_BOTS=3 docker compose --profile presenter up --build`: bots bid across several rounds, the winners wall fills, the round number increments
- [x] 9.2 Participant view in the browser: bid, get outbid (toast), watch eggs in batches (poll 2s) against one at a time (stream), then SOLD and the next round
- [x] 9.3 Stage view at 1920×1080 and 1280×720: no scroll, readable, full-screen SOLD; flipping transport keeps the chart history
- [ ] 9.4 With the network disconnected (no internet): fonts still render; DevTools shows no external requests
- [x] 9.5 With reduced motion emulated in DevTools: no animations
- [ ] 9.6 Screenshot both views on the real projector during rehearsal, and tune the gold gradient if it looks muddy
