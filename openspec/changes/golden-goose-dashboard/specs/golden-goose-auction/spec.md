## ADDED Requirements

### Requirement: Single prize auctioned in rounds
The auction server SHALL auction exactly one lot, `golden-goose` ("The Golden Goose", emoji `🦆`), and
SHALL reopen it at its starting price, with no bids, after every round's intermission.

#### Scenario: Goose reopens after a sold round
- **WHEN** a round of `golden-goose` closes sold and the intermission elapses
- **THEN** a `LOT_OPENED` event carries lot `golden-goose` with status open, no highest bid, no recent bids and a bid count of 0

#### Scenario: Goose reopens after an unsold round
- **WHEN** a round closes with no bids and the intermission elapses
- **THEN** `golden-goose` opens again as a new round

### Requirement: Round number
Every `Auction` message SHALL carry `round`, which is 1 for the first round after the server starts and
increases by exactly 1 each time a lot opens.

#### Scenario: Round increments on reopen
- **WHEN** round 3 closes and the next round opens
- **THEN** every snapshot and event after that carries `round = 4`

#### Scenario: Round survives the intermission
- **WHEN** round 2 has closed and the intermission is running
- **THEN** snapshots still carry `round = 2`

### Requirement: Lot duration is exposed
Every `Auction` message SHALL carry `lot_duration_ms`, the configured time a round opens for, whether or not
late bids have extended it.

#### Scenario: Configured duration
- **WHEN** the server runs with `-lot-duration=60s`
- **THEN** every snapshot carries `lot_duration_ms = 60000`

### Requirement: Winners record
The server SHALL record a `Sale` (round, bidder, amount, sold_at) when a round closes sold, and every
`Auction` message SHALL carry the recorded sales in `winners`, newest first, capped at 10. Unsold rounds
SHALL NOT be recorded. The record is held in memory only.

#### Scenario: Sold round is recorded
- **WHEN** round 1 closes with Alice's highest bid of €85
- **THEN** the `LOT_CLOSED` event's `winners` starts with `{round: 1, bidder: "Alice", amount: 85}`

#### Scenario: Unsold round is not recorded
- **WHEN** a round closes without bids
- **THEN** `winners` is unchanged

#### Scenario: Cap
- **WHEN** 11 rounds have been sold
- **THEN** `winners` holds the 10 most recent sales, and the sale from round 1 has been dropped

### Requirement: Additive proto change
The new fields and the `Sale` message SHALL only use new field numbers, and no existing field, RPC or enum value
SHALL change. `buf lint` SHALL pass and the regenerated `auction-server/gen/` SHALL be committed.

#### Scenario: Old client keeps working
- **WHEN** a BFF built from the previous proto calls the new server
- **THEN** `GetAuction`, `PlaceBid` and `WatchAuction` behave as before, and the new fields are ignored

### Requirement: BFF exposes rounds and winners
The BFF's auction JSON SHALL include `round`, `lotDurationMs` and `winners` (an array of
`{round, bidder, amount, soldAt}`, and an empty array when there are none) on `GET /api/auction`, `POST /api/bids`
and every SSE `auction` event. `bff/AuctionEndpoints.cs` SHALL keep `PlaceBid` unimplemented, and
`workshop/solution/AuctionEndpoints.cs` SHALL stay identical to it apart from that endpoint.

#### Scenario: JSON shape
- **WHEN** the browser calls `GET /api/auction` during round 2 after Alice won round 1
- **THEN** the response contains `"round": 2` and `"winners": [{"round": 1, "bidder": "Alice", "amount": 85, "soldAt": "…"}]`

### Requirement: Rehearsal bots bid every round
The simulated bidders SHALL pick a fresh budget whenever the round changes, so they keep bidding in every
round of the single lot.

#### Scenario: Bots in round 2
- **WHEN** the server runs with `-bots=3` and round 2 opens
- **THEN** bots place bids in round 2, as they did in round 1
