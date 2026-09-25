# gRPC Auction House

A one-hour talk and hands-on workshop about gRPC, where the room becomes the demo. The presenter runs
one gRPC auction server; every participant runs a frontend and BFF, implements the `PlaceBid` call,
and joins a live auction. Then everyone flips from polling to streaming and watches bids land on
every screen at once.

```
browser ──HTTP/JSON, SSE──▶ BFF (.NET 10) ──gRPC──▶ auction-server (Go)
 Vue + Vuetify0 + Tailwind    translates           one per room, :50051
```

| Path | What |
| --- | --- |
| [`proto/`](proto) | The contract: `auction.proto` |
| [`auction-server/`](auction-server) | Go gRPC server: unary `GetAuction` / `PlaceBid`, server-streaming `WatchAuction` |
| [`bff/`](bff) | ASP.NET Core minimal API; typed client generated from the proto at build time |
| [`frontend/`](frontend) | Vue 3, [Vuetify0](https://0.vuetifyjs.com) headless components, Tailwind CSS v4 |
| [`slides/`](slides) | Slidev deck on the Dignitas theme |
| [`workshop/`](workshop) | Participant handout, presenter guide, solution |
| [`reference/`](reference) | The conversation the talk was designed from |

## Quick start

Participants: [`QUICKSTART.md`](QUICKSTART.md) is the pre-workshop setup, ready to post in a chat.

```bash
# Participant: connect to the presenter's auction server
cp .env.example .env          # set AUCTION_HOST=http://<presenter-ip>:50051
docker compose up --build --watch

# Presenter, or anyone solo: also run the auction server (add AUCTION_BOTS=3 to rehearse alone).
# Leave AUCTION_HOST unset or empty: the BFF finds the local auction-server through Docker DNS.
docker compose --profile presenter up --build
```

Open <http://localhost:8080>. Then follow [`workshop/README.md`](workshop/README.md). The projector layout for the
presenter is <http://localhost:8080/?view=stage>.

## Development

| Task | Command |
| --- | --- |
| Regenerate Go code after changing the proto | `docker run --rm -v "$PWD:/workspace" -w /workspace bufbuild/buf generate` |
| Lint the proto | `docker run --rm -v "$PWD:/workspace" -w /workspace bufbuild/buf lint` |
| Test the auction server | `cd auction-server && go test ./...` |
| Run natively | `cd auction-server && go run . -bots 3` · `cd bff && dotnet run` · `cd frontend && npm run dev` |
| Slides | `cd slides && npm install && npm run dev` (needs the Dignitas npm feed) |

The C# client needs no generation step: `Grpc.Tools` regenerates it on every `dotnet build`.
