## ADDED Requirements

### Requirement: Golden duck SVG
The frontend SHALL include a golden duck SVG artwork, derived from an open-licensed duck emoji SVG and recoloured
with gold gradients, a specular highlight and a slow shimmer sweep. It SHALL look the same on every OS, because it
does not depend on the system emoji font.

#### Scenario: Windows 10 laptop
- **WHEN** the dashboard runs on a machine whose emoji font lacks newer glyphs
- **THEN** the golden duck renders identically to other machines

#### Scenario: Reduced motion
- **WHEN** `prefers-reduced-motion: reduce` is set
- **THEN** the duck renders gold without the shimmer

### Requirement: Artwork lookup with emoji fallback
The lot artwork SHALL be chosen by `lot.id`: `golden-goose` renders the golden duck SVG, and any other lot renders
`lot.emoji` as text.

#### Scenario: Unknown lot
- **WHEN** a server sends a lot with id `keyboard` and emoji `⌨️`
- **THEN** the hero shows `⌨️`

### Requirement: Licence notice
The artwork's source and licence SHALL be recorded in `frontend/THIRD_PARTY_NOTICES.md`, and in a comment in the
artwork component.

#### Scenario: Attribution present
- **WHEN** someone inspects the repository
- **THEN** the notice names the source project, the file and its licence

### Requirement: Favicon
Once the lot is revealed, the page favicon SHALL be a static golden duck SVG. While it is veiled, the favicon SHALL be a
gavel.

#### Scenario: Browser tab
- **WHEN** the app is open in a browser tab and PlaceBid works
- **THEN** the tab icon is the golden duck

#### Scenario: Not done yet
- **WHEN** the BFF still answers PlaceBid with 501
- **THEN** the tab icon is a gavel

### Requirement: Gold is decoration, not text
Gold tones SHALL be used only for artwork, borders, glows and other decoration. Any text in a gold family SHALL
use the bronze token, with a contrast of at least 4.5:1 against its background.

#### Scenario: Contrast check
- **WHEN** a gold-family colour is used for text
- **THEN** it is the bronze token, and passes 4.5:1 on white and on the background colour

### Requirement: The goose is unlocked by PlaceBid
Until this browser's BFF implements PlaceBid, the lot SHALL be veiled: a cloth instead of the artwork, the title
"A mystery lot", a description asking to implement PlaceBid, no egg nest, and no goose wording anywhere on the page. The
frontend SHALL find out with a probe that cannot place a bid (`POST /api/bids` with an empty bidder). It SHALL reveal the
lot when the probe comes back as 400 `InvalidArgument`, or when a bid succeeds. While veiled, it SHALL probe again every
few seconds, and it SHALL reveal without a page reload. The auction server and the BFF JSON SHALL NOT change.

#### Scenario: Workshop not done yet
- **WHEN** a participant opens the dashboard and the BFF answers PlaceBid with 501
- **THEN** the hero shows the cloth and "A mystery lot", and the words "goose" and "Golden Goose" appear nowhere on the page

#### Scenario: Participant finishes the exercise
- **WHEN** the participant saves a working PlaceBid and `--watch` restarts the BFF
- **THEN** within a few seconds, without a reload, the cloth lifts and the Golden Goose appears

#### Scenario: Probe places no bid
- **WHEN** the probe runs
- **THEN** the auction server rejects it with INVALID_ARGUMENT, and no bid is recorded

#### Scenario: Reload after finishing
- **WHEN** a participant who already revealed the goose reloads the page
- **THEN** the goose shows right away, without the cloth flashing first

#### Scenario: Reduced motion
- **WHEN** `prefers-reduced-motion: reduce` is set and the lot is revealed
- **THEN** the cloth is replaced by the duck with no animation
