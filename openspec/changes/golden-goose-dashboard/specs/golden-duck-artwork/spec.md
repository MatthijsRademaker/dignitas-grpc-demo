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
The page favicon SHALL be a static golden duck SVG, replacing the 🔨 emoji favicon.

#### Scenario: Browser tab
- **WHEN** the app is open in a browser tab
- **THEN** the tab icon is the golden duck

### Requirement: Gold is decoration, not text
Gold tones SHALL be used only for artwork, borders, glows and other decoration. Any text in a gold family SHALL
use the bronze token, with a contrast of at least 4.5:1 against its background.

#### Scenario: Contrast check
- **WHEN** a gold-family colour is used for text
- **THEN** it is the bronze token, and passes 4.5:1 on white and on the background colour
