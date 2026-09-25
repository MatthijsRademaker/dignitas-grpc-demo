## ADDED Requirements

### Requirement: Stage view by query parameter
With `?view=stage`, the frontend SHALL render the stage view instead of the participant view. It SHALL combine
with `?transport=stream`.

#### Scenario: Presenter opens the stage
- **WHEN** the presenter opens `http://localhost:8080/?view=stage&transport=stream`
- **THEN** the stage view renders and uses streaming

### Requirement: Projector layout
The stage view SHALL fill a 1920×1080 viewport with no scrolling. It SHALL show the golden duck artwork large, with
its egg nest, the round number, the price with its leading bidder, the countdown ring, the latest three bids, the
room-load chart as the largest data element, this screen's "bids seen after" figure, and the winners wall as a
strip. It SHALL NOT show the bid form or the name field. The price and countdown text SHALL be at least 96px
tall.

#### Scenario: Readable from the back
- **WHEN** the stage view runs at 1920×1080
- **THEN** there is no scrollbar and no bid form, and the price is ≥ 96px

#### Scenario: Smaller projector
- **WHEN** the viewport is 1280×720
- **THEN** the same layout scales down and still does not scroll

### Requirement: Presenter keeps the transport toggle
The stage view SHALL show the transport toggle and the poll interval toggle, visually quieter than the auction content,
so the presenter can flip the projector between polling and streaming.

#### Scenario: Flip on stage
- **WHEN** the presenter switches the stage view to streaming
- **THEN** the "bids seen after" figure resets for this screen, and the room-load chart keeps its history

### Requirement: Full-screen SOLD moment
When a round closes sold, the stage view SHALL show the SOLD moment across the whole screen for the intermission
(gavel, winner, amount, next-round countdown), then return to the auction layout when the next round opens.

#### Scenario: Round closes on stage
- **WHEN** a `LOT_CLOSED` event with a winner arrives
- **THEN** the full-screen SOLD moment shows until `LOT_OPENED`

### Requirement: Stage reveals like any other screen
The stage view SHALL follow the same veil rule as the participant view, based on the presenter's own BFF.

#### Scenario: Presenter applies the solution live
- **WHEN** the presenter live-codes PlaceBid during the workshop
- **THEN** the projector lifts the cloth and shows the Golden Goose to the room
