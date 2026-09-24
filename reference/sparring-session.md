# Sparring session: designing the gRPC talk

The voice conversation the talk, demo app and workshop were designed from. Condensed to the
decisions it produced; the slides and app implement these.

## Audience and framing

- Mixed audience, mostly backend, varying experience. They think in CRUD over HTTP, frontend to
  backend, not backend to backend.
- Start with the mental model: RPC over HTTP/2 with Protobuf contracts, instead of ad hoc JSON over
  REST. "REST says: here are resources, go fetch them. gRPC says: here's a service, call its
  methods, maybe stream."
- Mention the browser early, then frame it: the browser speaks HTTP/JSON to a BFF or gateway;
  behind that boundary, services use gRPC. gRPC-Web exists but needs a proxy and isn't identical.
- Terminology just in time, not as a glossary dump: a "reading a proto file" pass top to bottom
  (syntax, messages, services, RPCs, field types and numbers, unary). Later, when streaming comes up,
  show only the one-line delta: `returns (stream …)`.

## Demo and workshop

- One app, a live auction. Unary `PlaceBid` for the workshop, server-streaming `WatchAuction` for
  the demo.
- Hub and spokes: the presenter runs the central auction service on all interfaces (port 50051);
  everyone runs Docker Compose with a frontend and BFF, pointed at the presenter's IP via
  `AUCTION_HOST`. Backup: a small cloud VM with the same compose file. Test connectivity beforehand.
- The workshop is unary-only: implement `PlaceBid` in the BFF, 10–20 lines, 15 minutes, small and
  guided. Streaming has lifecycle, backpressure and bridging gotchas, so it is not a first operation.
- Then flip on streaming and the room becomes the demo: everyone sees bids almost instantly.
- The BFF translates the stream to SSE: make "stream translation" explicit.
- Show polling first for contrast, then flip to streaming so the delay and flicker vanish.

## Landing the use cases

- After the auction, a translation slide: "What we just did wasn't about an auction. It was a
  server saying: tell me when something changes." Map it to order status, deployment progress, and
  telemetry, in the format: what we want / how we fake it with HTTP / what the interaction actually
  is / the gRPC shape.
- "Why invest if my frontend can't use it?" Sometimes you shouldn't. It pays off behind the BFF:
  fan-out, polyglot services, real-time pipelines. Choose per boundary, not one protocol everywhere.
  Start simple (unary behind the BFF), stream when you need real-time. An advanced option: the BFF
  keeps a live snapshot from streams and the frontend polls it; present as an option, since it makes
  the BFF stateful.
- Streaming vs eventing (90 seconds): a gRPC stream is a direct pipe, both sides connected at the
  same time. A broker adds decoupling in time and responsibility: durability, fan-out, replay. Live
  bidders watching now fit a stream; processing tomorrow or replay needs a broker.
- War story (30 seconds): a file upload where the gateway didn't chunk well; a client stream of
  chunks to a gRPC service was far faster. The win was the streaming model and fewer gateway
  constraints, not gRPC magic.
- Sharp edges to name briefly and defer: cancellation, deadlines, flow control.

## Timing (60 minutes)

| Minutes | Segment |
| --- | --- |
| 7 | Why gRPC, one diagram |
| 7 | Reading a proto (the auction proto) |
| 5 | How it works under the hood (compress first if time slips) |
| 15 | Unary workshop: `PlaceBid` |
| 10 | Live auction payoff |
| 7 | Streaming shapes, using the auction, plus the translation slide |
| 4 | Trade-offs, one slide |
| 4 | Q&A |

Protect the workshop and the live auction: that is what people remember.
