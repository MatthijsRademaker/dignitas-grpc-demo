## ADDED Requirements

### Requirement: Participant view is the default
Without a `view` query parameter, the frontend SHALL render the participant view: a hero lot card, a bid form,
the bid feed, the per-browser wire panel, the room-load chart and the winners wall, in a light theme. The
existing `?transport=stream` parameter SHALL keep working.

#### Scenario: Default URL
- **WHEN** a participant opens `http://localhost:8080`
- **THEN** the participant view renders, with the bid form visible

### Requirement: Hero lot card
The hero card SHALL show the lot artwork (see `golden-duck-artwork`), the title in the display serif, the
description, the round number, the current price (the highest bid, or the starting price), the leading bidder and a
countdown.

#### Scenario: No bids yet
- **WHEN** round 4 has just opened
- **THEN** the card shows "Round 4", "Starting at €5" and "no bids yet"

### Requirement: Rolling price digits
When the price changes, each changed digit SHALL roll to its new value instead of being replaced.

#### Scenario: New highest bid
- **WHEN** the price goes from €85 to €90
- **THEN** the tens and units digits roll to 9 and 0, and the hero briefly highlights

### Requirement: Countdown ring
While a round is open, the countdown SHALL show `m:ss` inside a ring whose fill is
`remaining / max(lotDurationMs, remaining)`. In the last 10 seconds, the ring and the text SHALL switch to the error colour
and pulse, and the label SHALL read "Going once, going twice…".

#### Scenario: Snipe extension
- **WHEN** a bid lands with 3s left and the snipe window extends the round to 10s
- **THEN** the ring refills to 10s' worth, and stays in the closing state

#### Scenario: After a refresh
- **WHEN** the page is reloaded with 30s of a 60s round left
- **THEN** the ring is half full

### Requirement: SOLD moment
When a round closes sold, the hero card SHALL show a SOLD state for the intermission: a gavel animation, the
winner's name and the amount in the display serif, and "next round in Ns". An unsold round SHALL show a quiet
"The goose flew off, nobody bid" state instead.

#### Scenario: Sold
- **WHEN** a `LOT_CLOSED` event arrives with Alice's €85 as the highest bid
- **THEN** the hero shows "SOLD to Alice · €85" and counts down to the next round

### Requirement: Golden eggs
Every bid this page sees for the first time SHALL drop one golden egg into a nest on the hero card. Bids that arrive in the
same response SHALL drop together. Bids that were already there when the page loaded SHALL appear in the nest without
animation. The nest SHALL empty when a new round opens, and SHALL show at most 24 eggs plus a "+N" count.

#### Scenario: Polling brings a batch
- **WHEN** a single poll response contains three bids this page had not seen
- **THEN** three eggs drop at the same moment

#### Scenario: Streaming brings them one by one
- **WHEN** three bids are placed a second apart while streaming
- **THEN** three eggs drop a second apart

#### Scenario: New round
- **WHEN** the next round opens
- **THEN** the nest is empty

### Requirement: Bid form with outbid feedback
The bid form SHALL keep its current behaviour (name, amount at or above the minimum, a bid button and a quick bid, and
the workshop message on 501). When the leading bidder changes from the participant's name to someone else within the
same round, a toast SHALL announce the new amount and offer a one-click bid at the new minimum.

#### Scenario: Outbid
- **WHEN** Bob outbids Alice while Alice has this page open under the name "Alice"
- **THEN** Alice sees a toast "Bob outbid you: €90" with a "Bid €95" action

#### Scenario: Workshop not done yet
- **WHEN** the BFF answers `POST /api/bids` with 501
- **THEN** the form shows the workshop message, as it does today

### Requirement: Bid feed with avatars
The bid feed SHALL list the latest bids, newest first. Each row SHALL show an initials avatar (its colour a stable
function of the name), the name, the time, how late this page saw the bid, and the amount. The participant's
own bids SHALL be marked.

#### Scenario: Same colour everywhere
- **WHEN** "Alice" bids twice
- **THEN** both rows show the same avatar colour, on every screen

### Requirement: Room-load chart over time
The dashboard SHALL chart the last 90 seconds of room load as two small multiples on one shared time axis:
polling calls per second, and open streams, each with its current value as a label. The history SHALL NOT
reset when this browser changes transport or poll interval.

#### Scenario: Room flips to streaming
- **WHEN** the room switches from polling to streaming over 20 seconds
- **THEN** the calls/s panel slopes down to 0, and the open streams panel climbs, in one continuous 90s window

#### Scenario: This browser switches transport
- **WHEN** this browser switches from polling to streaming
- **THEN** the chart keeps the samples collected before the switch

### Requirement: Winners wall
The dashboard SHALL list earlier winners from `winners` (round, name, amount), newest first. When there are no winners,
it SHALL show an invitation to become the first owner of the goose. When the field is missing (an older server),
the wall SHALL be hidden.

#### Scenario: After two sold rounds
- **WHEN** Alice won round 1 for €85 and Bob won round 2 for €120
- **THEN** the wall lists "Round 2 · Bob · €120", then "Round 1 · Alice · €85"

### Requirement: Loading and error states
Until the first auction arrives, the dashboard SHALL show skeleton placeholders shaped like the hero, the feed and the
chart. When the connection fails, the error text SHALL be visible.

#### Scenario: BFF down
- **WHEN** the BFF is unreachable at startup
- **THEN** the skeletons show together with "Cannot reach the BFF"

### Requirement: Offline fonts
The frontend SHALL bundle every font it uses (a display serif, Nunito Sans, Fira Code) and SHALL NOT request
fonts or other assets from the network at runtime.

#### Scenario: No internet
- **WHEN** the app runs on a machine without internet access
- **THEN** the lot title renders in the display serif, and the UI renders in Nunito Sans

### Requirement: Reduced motion
With `prefers-reduced-motion: reduce`, the shimmer, pulse and SOLD burst SHALL be disabled, price digits SHALL change
instantly, and eggs SHALL appear without falling.

#### Scenario: Reduced motion
- **WHEN** the OS requests reduced motion and a bid arrives
- **THEN** the price and the egg update with no animation
