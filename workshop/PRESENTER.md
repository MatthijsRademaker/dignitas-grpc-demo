# Presenter guide

Everything you need to run the live auction on the day. The talk itself is in
[`slides/`](../slides); speaker notes (press `P` in Slidev) carry the timings and cues.

## The day before

1. **Rehearse alone** with simulated bidders:

   ```bash
   AUCTION_BOTS=3 docker compose --profile presenter up --build
   ```

   Leave `AUCTION_HOST` unset (or empty in `.env`): your BFF then finds the local auction-server
   through Docker DNS. Open <http://localhost:8080>. Your BFF still has the workshop TODO. Apply the solution so you can
   bid too: `cp workshop/solution/AuctionEndpoints.cs bff/` (or live-code it during the workshop and
   `git checkout bff/AuctionEndpoints.cs` afterwards).

2. **Prove the room can reach you.** From a *second* machine on the same network:

   ```bash
   docker run --rm fullstorydev/grpcurl -plaintext <your-ip>:50051 list
   ```

   You should see `auction.v1.AuctionService`. If not, see [Networking](#networking).

3. **Pre-pull images** on the venue Wi-Fi if you can, or at least on your own machine:
   `docker compose --profile presenter build`.

4. Put your IP on the "Get connected" slide.

## On the day

| When | What you run / do |
| --- | --- |
| Before people arrive | `docker compose --profile presenter up --build` (no bots). Check the IP again: venue DHCP may have changed it |
| "Before we start" slide | Everyone runs `cp .env.example .env` and starts `docker compose build` |
| Workshop | Walk the room. Watch the server log: every participant's first `PlaceBid` shows up as `unary … code=OK` from their IP |
| End of the workshop | Not everyone finishes. Leave the catch-up line on the "One possible solution" slide up: `cp workshop/solution/AuctionEndpoints.cs bff/` |
| Live auction | Projector on the app in **polling** mode. Let the room bid for a lot, and point at "the whole room" panel: calls per second, and how many found nothing new. Switch the projector to 0.5s to show the polling dilemma. Then everyone flips to **streaming**: calls per second fall to zero, open streams climb, "bids seen after" drops to 0.00s |
| Q&A | Leave the auction running |

Useful server flags (append to `command` in `compose.yaml`): `-lot-duration=60s`,
`-snipe-window=10s`, `-intermission=8s`, `-bots=N`.

## Networking

The server listens on `:50051` on all interfaces. Participants need to reach **your LAN IP** on that
port over plain TCP (no TLS).

- **Firewall:** allow inbound TCP 50051. On Windows (PowerShell as admin):
  `New-NetFirewallRule -DisplayName "gRPC auction" -Direction Inbound -Protocol TCP -LocalPort 50051 -Action Allow`
- **WSL2 with Docker Engine inside WSL:** WSL's default NAT networking hides the port from the LAN.
  Either switch to mirrored networking (`%UserProfile%\.wslconfig`: `[wsl2]` then
  `networkingMode=mirrored`, then `wsl --shutdown`) and allow inbound traffic through the Hyper-V
  firewall, or run the server with Docker Desktop, which publishes ports on the Windows host.
- **Guest Wi-Fi with client isolation** blocks laptop-to-laptop traffic entirely. Fallbacks:
  - Run the server on a small cloud VM (`docker compose --profile presenter up auction-server`) and
    hand out its IP instead. Same compose file, same flow.
  - Worst case, everyone runs solo: `docker compose --profile presenter up`, with `AUCTION_HOST` unset.
    The workshop still works; the shared auction doesn't.

## What to point at during the talk

| Concept | Where it lives |
| --- | --- |
| The contract | `proto/auction/v1/auction.proto` |
| Codegen at build time | `<Protobuf>` item in `bff/Bff.csproj`; generated Go in `auction-server/gen/` (`buf generate`) |
| Unary call + deadline | `GET /auction` in `bff/AuctionEndpoints.cs` |
| gRPC status → HTTP | `bff/GrpcErrors.cs` |
| Server streaming, Go side | `WatchAuction` in `auction-server/internal/auction/service.go` |
| Stream translation to SSE | `GET /auction/stream` in `bff/AuctionEndpoints.cs` |
| Slow consumers / flow control | `broadcastExcept` in `auction-server/internal/auction/house.go` |
| Room load, empty polls, bid delay | `Poll`, `announceLoad` and `snapshot` in `auction-server/internal/auction/house.go` |
| Keepalives | `SocketsHttpHandler` in `bff/Program.cs`; `KeepaliveEnforcementPolicy` in `auction-server/main.go` |
| Reflection (for grpcurl) | `reflection.Register` in `auction-server/main.go` |
