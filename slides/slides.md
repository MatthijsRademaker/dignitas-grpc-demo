---
theme: '@dignitas/slidev-theme'
title: 'gRPC Auction House'
info: |
  ## gRPC Auction House
  Contracts, codegen and streams: a one-hour talk and workshop where the room becomes the demo.

  A GenTech session by Matthijs Rademaker.
duration: 60min
---

::title::
gRPC Auction House

::subtitle::
Contracts, streams and a room full of bidders · Matthijs Rademaker · GenTech

<!--
One hour. The whole talk hangs on one app: a live auction everyone in this room
takes part in. You will build one piece of it yourself, then watch the whole room
light up when we switch it to streaming.
-->

---

# The hour

<div class="agenda">
  <div class="slot"><span class="t">7′</span> Why gRPC: the mental model</div>
  <div class="slot"><span class="t">7′</span> Reading a <code>.proto</code> file</div>
  <div class="slot"><span class="t">5′</span> Under the hood: why it is fast</div>
  <div class="slot hi"><span class="t">15′</span> Workshop: implement <code>PlaceBid</code></div>
  <div class="slot hi"><span class="t">10′</span> Live auction: the whole room</div>
  <div class="slot"><span class="t">7′</span> Streaming shapes and what they are for</div>
  <div class="slot"><span class="t">4′</span> Trade-offs</div>
  <div class="slot"><span class="t">5′</span> Questions</div>
</div>

<!--
Point at the two highlighted rows: that is the part people remember, and the
part I protect if we run late. If time slips, "under the hood" gets compressed.
-->

<style>
.agenda { margin-top: 1.2rem; display: grid; grid-template-columns: 1fr 1fr; gap: 0.6rem 1.2rem; font-size: 1.15rem; }
.slot { display: flex; align-items: center; gap: 0.9rem; padding: 0.55rem 0.9rem; border-radius: 0.6rem;
  border: 1px solid color-mix(in srgb, currentColor 14%, transparent); }
.slot.hi { border-color: var(--se-color-primary); background: color-mix(in srgb, var(--se-color-primary) 8%, transparent); font-weight: 700; }
.slot .t { flex: none; width: 2.8rem; font-weight: 800; color: var(--se-color-primary); }
</style>

---

# Before we start: kick off the build

<div class="setup">

```bash
git clone https://github.com/MatthijsRademaker/dignitas-grpc-demo.git
cd dignitas-grpc-demo && cp .env.example .env
docker compose build    # a few minutes: start it now
```

</div>

<p class="mt-6 text-lg">Only Docker needed. No Go, no .NET, no Node on your machine. Got the .NET SDK? Even better: <code>cd bff && dotnet build</code>.</p>

<!--
Get the slow part out of the way now so the workshop is 15 minutes of code,
not 15 minutes of downloading base images. Copying .env now saves a step
later; the real IP comes on the "Get connected" slide. Walk around briefly if people get stuck; don't block on it.
-->

<style>
.setup pre { font-size: 1.25rem !important; }
</style>

---
layout: section
---

# Why gRPC

---

# Two ways to ask the same question

<div class="grid grid-cols-2 gap-8 mt-6">
<div class="ask">
  <h3>REST</h3>
  <p class="sub">“Here are resources, go fetch them.”</p>

```http
GET  /auctions/current
POST /auctions/current/bids
     { "bidder": "Ada", "amount": 120 }
```

  <p class="note">You design URLs, verbs, JSON shapes and status codes, then write a client by hand.</p>
</div>
<div v-click class="ask primary">
  <h3>gRPC</h3>
  <p class="sub">“Here’s a service, call its methods.”</p>

```csharp
await auction.GetAuctionAsync(new());
await auction.PlaceBidAsync(new() {
    Bidder = "Ada", Amount = 120 });
```

  <p class="note">You design a <strong>contract</strong>. The client is generated from it.</p>
</div>
</div>

<!--
The mental model shift, before any detail. REST thinks in nouns you fetch.
gRPC thinks in verbs you call, a remote procedure call, hence the name.
Both are fine. The rest of the hour is about when the second one pays off.
-->

<style>
.ask { border: 2px solid color-mix(in srgb, currentColor 15%, transparent); border-radius: 0.75rem; padding: 0.9rem 1.2rem; }
.ask.primary { border-color: var(--se-color-primary); background: color-mix(in srgb, var(--se-color-primary) 6%, transparent); }
.ask h3 { font-weight: 800; font-size: 1.3rem; }
.ask .sub { opacity: 0.7; margin: 0.1rem 0 0.6rem; }
.ask .note { font-size: 0.95rem; margin-top: 0.6rem; }
</style>

---

# gRPC in one sentence

<p class="text-2xl mt-6">Remote procedure calls over <span class="text-primary font-bold">HTTP/2</span>, described by a <span class="text-primary font-bold">Protobuf contract</span>.</p>

<div class="pillars mt-8">
  <div v-click class="pillar">
    <div class="h">📜 Contract first</div>
    <p>A <code>.proto</code> file defines services, methods and messages. It is the single source of truth.</p>
  </div>
  <div v-click class="pillar">
    <div class="h">⚙️ Generated code</div>
    <p>Typed clients and server stubs for Go, C#, Java, TypeScript… from the same file.</p>
  </div>
  <div v-click class="pillar">
    <div class="h">🚄 HTTP/2 + binary</div>
    <p>Compact messages, many calls over one connection, and <strong>streaming</strong> built in.</p>
  </div>
</div>

<!--
Three pillars, one click each. Originally built at Google, now a CNCF project.
Don't go deep yet: each pillar gets its own moment in the next 15 minutes.
-->

<style>
.pillars { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; }
.pillar { border-radius: 0.75rem; padding: 1rem 1.1rem; border: 1px solid color-mix(in srgb, var(--se-color-primary) 30%, transparent);
  background: color-mix(in srgb, var(--se-color-primary) 6%, transparent); }
.pillar .h { font-weight: 800; font-size: 1.1rem; margin-bottom: 0.4rem; }
.pillar p { font-size: 0.95rem; }
</style>

---
clicks: 3
---

# “But my browser can’t speak gRPC…”

<ArchBoundary class="mt-2" />

<!--
Say the objection before they do. Browsers can't do native gRPC: they don't
expose HTTP/2 framing and trailers. gRPC-Web exists (with a proxy), but it isn't
identical. The usual answer: the browser speaks HTTP/JSON to a BFF or gateway.
[click] Behind that boundary, services speak gRPC. [click] Service to service,
polyglot, that is where it shines. [click] The punchline: choose per boundary.
This is exactly the architecture of tonight's app.
-->

---
layout: section
---

# Reading a `.proto` file

---
clicks: 5
---

# `auction.proto`, top to bottom

<div class="proto-grid">
<div>

```proto {1-2|4-7|6|9-13|10-12|15-20}
syntax = "proto3";
package auction.v1;

service AuctionService {
  rpc GetAuction(GetAuctionRequest) returns (GetAuctionResponse);
  rpc PlaceBid(PlaceBidRequest) returns (PlaceBidResponse);
}

message PlaceBidRequest {
  string lot_id = 1;
  string bidder = 2;
  int64 amount = 3;  // whole euros
}

message Auction {
  Lot lot = 1;
  LotStatus status = 2;          // an enum
  Bid highest_bid = 3;           // unset until the first bid
  repeated Bid recent_bids = 4;  // a list
}
```

</div>
<div class="terms">
  <div class="term"><b>syntax · package</b> Protobuf version, and a namespace so contracts don’t collide.</div>
  <div v-click="1" class="term"><b>service · rpc</b> A service is a set of methods. Each <code>rpc</code> is one method you can call.</div>
  <div v-click="2" class="term"><b>unary</b> One request in, one response out. A function call over the network.</div>
  <div v-click="3" class="term"><b>message</b> A typed struct. Scalars, other messages, enums.</div>
  <div v-click="4" class="term"><b>field number</b> The field’s identity on the wire. Rename freely; never change or reuse the number.</div>
  <div v-click="5" class="term"><b>nesting · repeated</b> Messages contain messages; <code>repeated</code> makes a list.</div>
</div>
</div>

<!--
Teach the terms on contact, not as a glossary dump. Walk the file top to bottom,
one click per concept. This is the real contract you'll use in the workshop in
ten minutes; the full file in the repo has comments and one more method. We
get to that third method later.
-->

<style>
.proto-grid { display: grid; grid-template-columns: 1.25fr 1fr; gap: 1.4rem; margin-top: 0.8rem; align-items: start; }
.proto-grid pre { font-size: 0.8rem !important; }
.terms { display: flex; flex-direction: column; gap: 0.35rem; }
.term { font-size: 0.82rem; line-height: 1.35; padding: 0.35rem 0.65rem; border-left: 3px solid var(--se-color-primary);
  background: color-mix(in srgb, var(--se-color-primary) 6%, transparent); border-radius: 0 0.4rem 0.4rem 0; }
.term b { color: var(--se-color-primary); font-weight: 800; margin-right: 0.35rem; }
</style>

---

# One contract, every language

<div class="grid grid-cols-2 gap-6 mt-4">
<div>
  <div class="lang-label">Go · the auction server implements</div>

```go
func (s *Service) PlaceBid(
    ctx context.Context,
    req *auctionv1.PlaceBidRequest,
) (*auctionv1.PlaceBidResponse, error)
```

</div>
<div>
  <div class="lang-label">C# · your BFF calls</div>

```csharp
var response = await auction.PlaceBidAsync(
    new PlaceBidRequest {
        LotId = "rubber-duck",
        Bidder = "Ada",
        Amount = 120 });
```

</div>
</div>

<div class="gen mt-6">
  <span class="chip">auction.proto</span>
  <span class="arr">→</span>
  <span class="chip tool">buf generate</span><span class="arr">→</span><span class="chip">Go server interface</span>
  <span class="sep">·</span>
  <span class="chip tool">dotnet build</span><span class="arr">→</span><span class="chip">C# client</span>
</div>

<p v-click class="mt-5 text-xl text-center">No hand-written HTTP client. No URL building. No JSON mapping. <span class="text-primary font-bold">Change the contract, and the compiler tells you what broke.</span></p>

<!--
Code generation is part of the build. In the BFF, Grpc.Tools regenerates the
client on every dotnet build. Go code is generated once with buf and committed.
Contract first: the proto is the source of truth, not the implementation.
-->

<style>
.lang-label { font-size: 0.8rem; font-weight: 800; letter-spacing: 0.05em; text-transform: uppercase; opacity: 0.6; margin-bottom: 0.3rem; }
.gen { display: flex; align-items: center; justify-content: center; gap: 0.5rem; flex-wrap: wrap; }
.gen .chip { padding: 0.3rem 0.7rem; border-radius: 0.5rem; font-weight: 700; font-size: 0.9rem;
  border: 2px solid color-mix(in srgb, currentColor 20%, transparent); }
.gen .chip.tool { border-color: var(--se-color-primary); color: var(--se-color-primary); font-family: 'Fira Code', monospace; }
.gen .arr { opacity: 0.5; font-weight: 800; }
.gen .sep { opacity: 0.3; margin: 0 0.4rem; }
</style>

---
layout: section
---

# Under the hood

---
clicks: 2
---

# Why it is small: the same bid, on the wire

<WireBytes class="mt-6" />

<!--
Bid { bidder: "Ada", amount: 120 }. JSON repeats every field name in every
message. [click] Protobuf sends a field number plus a type tag, then the value.
[click] That's why field numbers are sacred: they ARE the field on the wire.
Smaller is nice; not having to parse text is the bigger CPU win. These are real
bytes, taken from the generated Go code.
-->

---
clicks: 1
---

# Why it is fast: one connection, many streams

<Multiplexing class="mt-4" />

<div v-click="1" class="grid grid-cols-3 gap-4 mt-2 text-base">
  <div class="fact"><b>Multiplexing</b> calls don’t wait for each other or need extra connections</div>
  <div class="fact"><b>Header compression</b> repeated headers become a few bytes</div>
  <div class="fact"><b>Long-lived streams</b> just another lane on a connection you already have</div>
</div>

<!--
HTTP/1.1: one request at a time per connection, so clients open several and
poll. [click] HTTP/2: binary frames from many streams interleaved on ONE
connection. That last point is what makes streaming cheap, and it sets up the
demo. If time is slipping, this is the slide to compress.
-->

<style>
.fact { padding: 0.5rem 0.8rem; border-radius: 0.5rem; background: color-mix(in srgb, var(--se-color-primary) 7%, transparent); }
.fact b { display: block; color: var(--se-color-primary); }
</style>

---
layout: section
---

# Workshop: `PlaceBid`

---

# The room is the system

<div class="flex justify-center mt-2">
  <RoomHub />
</div>

<p class="text-center text-lg mt-2">My laptop runs the auction. <span class="text-primary font-bold">Your BFF talks gRPC to it</span> across the room’s Wi-Fi.</p>

<!--
Hub and spokes. The auction server runs on my machine, listening on all
interfaces. Every one of you runs a frontend and a BFF locally, pointed at my IP.
When the workshop works, every bid from every laptop lands in the same auction.
-->

---

# Get connected

<div class="steps">

```bash
cp .env.example .env
# AUCTION_HOST=http://192.168.1.20:50051   ← my IP, see below
docker compose up --build --watch
```

</div>

<div class="grid grid-cols-2 gap-6 mt-5 text-lg">
  <div class="check">Open <code>http://localhost:8080</code>: you should see the lot and live bids.</div>
  <div class="check">Try to bid: you get <span class="text-primary font-bold">“Workshop time”</span>. That’s your job.</div>
</div>

<p class="ip">Presenter: <code>AUCTION_HOST=http://192.168.1.20:50051</code></p>

<!--
Update the IP on this slide right before the session (ip -4 addr / ipconfig).
--watch rebuilds the BFF container on every save, so the edit loop is: save,
wait a few seconds, refresh. Seeing the lot proves GetAuction already works
through your BFF: that call is prebuilt, and it's your template.
-->

<style>
.steps pre { font-size: 1.15rem !important; }
.check { padding: 0.6rem 0.9rem; border-radius: 0.6rem; border: 1px solid color-mix(in srgb, currentColor 15%, transparent); }
.ip { margin-top: 1.4rem; text-align: center; font-size: 1.6rem; font-weight: 800; }
.ip code { color: var(--se-color-primary); }
</style>

---

# Your task: `bff/AuctionEndpoints.cs`

<div class="task-grid">
<div>

```csharp
api.MapPost("/bids", (PlaceBidBody body,
    AuctionService.AuctionServiceClient auction,
    CancellationToken ct) =>
{
    // 1. Build a PlaceBidRequest from `body`
    // 2. await auction.PlaceBidAsync(request,
    //        deadline: ..., cancellationToken: ct)
    // 3. return Results.Ok(AuctionDto.From(...))
    // 4. catch RpcException → ex.ToProblem()
});
```

</div>
<div class="hints">
  <div class="hint"><b>Template</b> <code>GET /auction</code>, right above it, does the same for <code>GetAuction</code>.</div>
  <div class="hint"><b>Discover</b> Type <code>auction.</code> and let the IDE show you the generated methods.</div>
  <div class="hint"><b>Stretch 1</b> Set the deadline to 1 ms. What happens?</div>
  <div class="hint"><b>Stretch 2</b> Bid too low. Which gRPC status comes back, and what does the browser see?</div>
</div>
</div>

<p class="text-center mt-4 text-lg opacity-70">⏱ 15 minutes · pairing encouraged · stuck? <code>workshop/README.md</code></p>

<!--
About ten lines of real code. Walk the room. Most common snags: forgetting
async on the lambda, AUCTION_HOST typos, firewall on the presenter machine.
Stretch 1 shows DeadlineExceeded mapped to 504. Stretch 2 shows
FailedPrecondition mapped to 409: gRPC status codes are not HTTP codes, and the
boundary translates them (GrpcErrors.cs).
-->

<style>
.task-grid { display: grid; grid-template-columns: 1.3fr 1fr; gap: 1.2rem; margin-top: 0.6rem; align-items: start; }
.task-grid pre { font-size: 0.85rem !important; }
.hints { display: flex; flex-direction: column; gap: 0.5rem; }
.hint { font-size: 0.92rem; padding: 0.45rem 0.7rem; border-radius: 0.5rem; background: color-mix(in srgb, var(--se-color-primary) 7%, transparent); }
.hint b { color: var(--se-color-primary); margin-right: 0.3rem; }
</style>

---

# One possible solution

```csharp {3-7|8|9-12|all}
api.MapPost("/bids", async (PlaceBidBody body, AuctionService.AuctionServiceClient auction, CancellationToken ct) =>
{
    try
    {
        var response = await auction.PlaceBidAsync(
            new PlaceBidRequest { LotId = body.LotId, Bidder = body.Bidder, Amount = body.Amount },
            deadline: DateTime.UtcNow.AddSeconds(2),
            cancellationToken: ct);
        return Results.Ok(AuctionDto.From(response.Auction));
    }
    catch (RpcException ex)
    {
        return ex.ToProblem();   // FailedPrecondition → 409, InvalidArgument → 400, …
    }
});
```

<p v-click class="mt-5 text-xl text-center">A typed call across the network. <span class="text-primary font-bold">The contract did the rest.</span></p>

<p v-click class="mt-4 text-lg text-center">Not done? <code>cp workshop/solution/AuctionEndpoints.cs bff/</code>, wait a few seconds, then bid.</p>

<!--
Walk the highlights: the call itself, the deadline (gRPC has no default
timeout; always set one), the mapping to JSON for the browser, the status-code
translation. Ask who got it working: hands up. Their bids are in MY auction now.
[click] Leave the catch-up line on screen for a moment before "Everyone: bid!":
--watch picks up the copied file, so nobody sits out the live auction.
-->

---
layout: section
---

# Live auction

---

# Everyone: bid!

<div class="cue">
  <div class="bar">you → the room · polling mode</div>
  <div class="msg">Pick a name. Outbid your neighbour. Snipe the rubber duck.</div>
  <div class="watch">
    <strong>Watch for:</strong> other people’s bids show up <span class="text-primary font-bold">about a second late</span> (“bids seen after”),
    the “on the wire” strip fills with <span class="text-primary font-bold">grey dots</span>, and on my screen the room’s
    <span class="text-primary font-bold">calls per second</span>, mostly finding nothing new.
  </div>
</div>

<!--
Let this run for a minute or two. Keep the presenter screen on the auction:
the "whole room" panel shows every laptop's polls hitting my server, and how
many found nothing new. The delay should feel slightly annoying: that's the
point. Late bids extend the clock, so it gets frantic at the end of each lot.
Then the polling dilemma: switch the projector to 0.5s. Bids arrive sooner,
and the room's calls per second jump. Faster polling buys freshness with load.
-->

<style>
.cue { margin-top: 1.5rem; border: 2px solid color-mix(in srgb, var(--se-color-primary) 35%, transparent); border-radius: 0.75rem; overflow: hidden; max-width: 800px; }
.cue .bar { background: color-mix(in srgb, var(--se-color-primary) 14%, transparent); padding: 0.45rem 1rem; font-size: 0.85rem; font-weight: 700; letter-spacing: 0.04em; }
.cue .msg { padding: 1rem; font-size: 1.4rem; font-weight: 800; }
.cue .watch { padding: 0.9rem 1rem; font-size: 1.05rem; border-top: 1px solid color-mix(in srgb, var(--se-color-primary) 25%, transparent);
  background: color-mix(in srgb, var(--se-color-primary) 5%, transparent); }
</style>

---
clicks: 1
---

# Polling asks. Streaming is told.

<PollVsStream class="mt-4" />

<!--
Top row: what actually happened. Middle: polling every two seconds, mostly
grey (nothing new), and every bid arrives late. [click] Bottom: one request
that stays open; each bid arrives the moment it happens. "Polling an order
status every 2 seconds? You're actually asking to subscribe to changes."
-->

---

# The one-line change

```proto {2|3}
service AuctionService {
  rpc GetAuction(GetAuctionRequest) returns (GetAuctionResponse);
  rpc WatchAuction(WatchAuctionRequest) returns (stream WatchAuctionResponse);
}
```

<p class="mt-6 text-2xl text-center">Same service, same contract. One keyword: <code class="text-primary font-bold">stream</code>.</p>

<p v-click class="mt-4 text-lg text-center opacity-80">The first message is a snapshot, then one message per change. The stream stays open until someone hangs up.</p>

<!--
Here is the third method I skipped earlier. This is the whole delta between
"ask me" and "tell me". No broker, no websocket server, no second API.
-->

---

# Stream translation at the boundary

<div class="grid grid-cols-2 gap-5 mt-3">
<div>
  <div class="lang-label">Go · auction server sends</div>

```go
events, stop := s.house.Watch()
defer stop()
for {
    select {
    case <-stream.Context().Done():
        return nil // client hung up
    case event := <-events:
        if err := stream.Send(event); err != nil {
            return err
        }
    }
}
```

</div>
<div>
  <div class="lang-label">C# · BFF relays as Server-Sent Events</div>

```csharp
api.MapGet("/auction/stream", (client, ct) =>
    TypedResults.ServerSentEvents(Watch(client, ct)));

async IAsyncEnumerable<AuctionEventDto> Watch(
    AuctionServiceClient client, CancellationToken ct)
{
    using var call = client.WatchAuction(new(),
        cancellationToken: ct);
    await foreach (var m in call.ResponseStream
                                .ReadAllAsync(ct))
        yield return AuctionEventDto.From(m);
}
```

</div>
</div>

<p class="mt-3 text-center text-lg">Browser closes the tab → <code>ct</code> fires → gRPC stream cancelled → server stops sending. <span class="text-primary font-bold">Cancellation flows end to end.</span></p>

<!--
The browser can't speak gRPC, so the BFF translates the stream into SSE,
which every browser supports. That's the architectural point: each boundary
uses the right protocol. The C# is lightly abbreviated; the real thing is in
bff/AuctionEndpoints.cs, and it's already in your BFF.
-->

<style>
.lang-label { font-size: 0.8rem; font-weight: 800; letter-spacing: 0.05em; text-transform: uppercase; opacity: 0.6; margin-bottom: 0.3rem; }
</style>

---

# Flip the switch

<div class="flex justify-center">
  <RoomHub streaming />
</div>

<div class="cue">
  <div class="msg">Everyone: toggle <span class="text-primary">Streaming · SSE ← gRPC</span>, and keep bidding.</div>
  <div class="watch">
    <strong>Watch for:</strong> on my screen, <span class="text-primary font-bold">calls per second dropping to zero</span> while open streams climb.
    On yours, “bids seen after” going to <span class="text-primary font-bold">0.00s</span>, and your HTTP request counter <span class="text-primary font-bold">standing still</span>.
  </div>
</div>

<!--
The payoff. Switch the projector to the auction app, in streaming mode. The
"whole room" panel is the server's view: polls per second fall away over ten
seconds as the room flips, and open WatchAuction streams jump. Then run a lot or two in streaming mode. Show the server
terminal briefly: "stream opened" lines from every IP in the room.
-->

<style>
.cue { margin: 0.4rem auto 0; border: 2px solid color-mix(in srgb, var(--se-color-primary) 35%, transparent); border-radius: 0.75rem; overflow: hidden; max-width: 820px; }
.cue .msg { padding: 0.7rem 1rem; font-size: 1.2rem; font-weight: 800; }
.cue .watch { padding: 0.7rem 1rem; font-size: 0.95rem; border-top: 1px solid color-mix(in srgb, var(--se-color-primary) 25%, transparent);
  background: color-mix(in srgb, var(--se-color-primary) 5%, transparent); }
</style>

---
layout: section
---

# Streaming shapes

---
clicks: 3
---

# Four shapes, one auction

<StreamShapes class="mt-2" />

<!--
Name the shapes using the auction as the example. Unary: you built it.
[click] Server streaming: you just watched it. [click] Client streaming: many
messages up, one answer back, like uploading a photo in chunks.
[click] Bidirectional: both sides send whenever they like on one stream.
-->

---

# It was never about the auction

<p class="text-xl mt-2">What we built was a server saying: <span class="text-primary font-bold">“I’ll tell you when something changes.”</span></p>

<table class="shapes">
  <thead>
    <tr><th>What we want</th><th>How we fake it with HTTP</th><th>What it actually is</th><th>The gRPC shape</th></tr>
  </thead>
  <tbody>
    <tr v-click>
      <td>Order status on screen</td><td>Poll <code>GET /orders/42</code> every 2s</td><td>A subscription to changes</td>
      <td><code>WatchOrder(…) returns (stream OrderUpdate)</code></td>
    </tr>
    <tr v-click>
      <td>Deployment progress</td><td>Poll the job, guess when it’s done</td><td>A progress feed</td>
      <td><code>Deploy(…) returns (stream Stage)</code></td>
    </tr>
    <tr v-click>
      <td>Device telemetry and commands</td><td>POST every reading, poll for commands</td><td>A conversation</td>
      <td><code>Connect(stream Reading) returns (stream Command)</code></td>
    </tr>
  </tbody>
</table>

<p v-click class="mt-5 text-center text-lg opacity-80">Not the only way to do it. gRPC makes it natural, inside the same typed contract.</p>

<!--
The generalisation moment, while it's fresh, for anyone who can't yet map an
auction to their day job. Same pattern three times: what we want, how we fake
it with HTTP, what the interaction really is, and the gRPC shape.
-->

<style>
.shapes { margin-top: 1.2rem; width: 100%; font-size: 0.95rem; border-collapse: separate; border-spacing: 0 0.35rem; }
.shapes th { text-align: left; font-size: 0.75rem; font-weight: 800; letter-spacing: 0.05em; text-transform: uppercase; opacity: 0.6; padding: 0 0.7rem; }
.shapes td { padding: 0.55rem 0.7rem; background: color-mix(in srgb, currentColor 4%, transparent); }
.shapes td:first-child { font-weight: 800; border-radius: 0.5rem 0 0 0.5rem; }
.shapes td:last-child { border-radius: 0 0.5rem 0.5rem 0; background: color-mix(in srgb, var(--se-color-primary) 9%, transparent); }
.shapes code { font-size: 0.82rem; }
</style>

---

# Streaming isn’t only for live updates

<div class="story">
  <div class="bar">a war story · client streaming</div>
  <div class="body">
    <p>We needed large file uploads. Our gateway didn’t do chunked uploads well.</p>
    <p class="mt-3">So we sent the blob as a <strong>client stream of chunks</strong> to a gRPC service, and it was <span class="text-primary font-bold">ridiculously fast</span> compared to the plain HTTP upload.</p>

```proto
rpc Upload(stream Chunk) returns (UploadResult);
```

  </div>
</div>

<p class="mt-5 text-lg text-center">The win wasn’t gRPC magic. It was the <span class="text-primary font-bold">streaming model</span>, and fewer gateway constraints in our stack.</p>

<!--
Thirty seconds, your own story. A completely different reason to stream:
not "live", just "I don't want one huge request". Don't dive into transport
internals; keep it at "here's a class of problem and here's how streaming helped".
-->

<style>
.story { margin-top: 1.2rem; border: 2px solid color-mix(in srgb, var(--se-color-primary) 35%, transparent); border-radius: 0.75rem; overflow: hidden; max-width: 820px; }
.story .bar { background: color-mix(in srgb, var(--se-color-primary) 14%, transparent); padding: 0.45rem 1rem; font-size: 0.85rem; font-weight: 700; letter-spacing: 0.04em; }
.story .body { padding: 1rem 1.2rem; font-size: 1.1rem; }
.story pre { margin-top: 0.8rem; font-size: 1rem !important; }
</style>

---

# “Isn’t this just eventing?”

<StreamVsBroker class="mt-4" />

<p v-click class="mt-5 text-lg text-center">
Live bidders watching <em>right now</em> → <span class="text-primary font-bold">stream</span>.
Invoice the winner tomorrow, replay for audit → <span class="font-bold">broker</span>.
</p>

<!--
Ninety seconds. Streaming is event-driven in spirit, but it's a direct pipe:
producer and consumer connected at the same time. A broker adds decoupling in
time and responsibility. Streaming is a communication pattern, not an event
backbone. You don't need a Service Bus or MQTT broker for a one-to-one live feed.
-->

---
layout: section
---

# Trade-offs

---

# Honest trade-offs

<div class="grid grid-cols-2 gap-6 mt-4">
<div class="side good">
  <h3>You gain</h3>
  <ul>
    <li><b>One contract</b>, typed clients in every language</li>
    <li><b>Performance</b>: small binary messages, multiplexed connections</li>
    <li><b>Streaming</b> as a first-class shape, not a bolt-on</li>
    <li><b>Evolvable APIs</b>: add fields without breaking callers</li>
  </ul>
</div>
<div class="side">
  <h3>You pay</h3>
  <ul>
    <li><b>Browsers</b> need a boundary: a BFF, gateway or gRPC-Web proxy</li>
    <li><b>Binary</b> is harder to eyeball: use <code>grpcurl</code>, Postman, reflection</li>
    <li><b>Infrastructure</b> must understand HTTP/2: proxies, load balancers</li>
    <li><b>New concepts</b>: deadlines, status codes, stream lifecycles</li>
  </ul>
</div>
</div>

<p class="mt-6 text-2xl text-center font-bold">Choose per boundary, <span class="text-primary">not one protocol everywhere.</span></p>

<!--
One slide, both columns, no hedging. If someone asks "so should we do this
everywhere?", calmly say no. That balance builds trust. The auction server has
reflection enabled: grpcurl -plaintext <ip>:50051 list works right now.
-->

<style>
.side { border: 2px solid color-mix(in srgb, currentColor 15%, transparent); border-radius: 0.75rem; padding: 0.9rem 1.2rem; }
.side.good { border-color: var(--se-color-primary); background: color-mix(in srgb, var(--se-color-primary) 6%, transparent); }
.side h3 { font-weight: 800; font-size: 1.2rem; margin-bottom: 0.5rem; }
.side ul { padding-left: 1.1rem; display: flex; flex-direction: column; gap: 0.35rem; font-size: 0.98rem; }
</style>

---
clicks: 5
---

# Five things to take home

<div class="mt-6 grid grid-cols-1 gap-3 text-lg">
  <div v-click class="take"><span class="n">1</span> gRPC: call methods on a service over HTTP/2, described by a <code>.proto</code> contract.</div>
  <div v-click class="take"><span class="n">2</span> The contract is the source of truth; clients and servers are generated from it.</div>
  <div v-click class="take"><span class="n">3</span> Field numbers, binary frames and multiplexing are why it’s small and fast.</div>
  <div v-click class="take"><span class="n">4</span> “Tell me when something changes” is one keyword: <code>stream</code>.</div>
  <div v-click class="take"><span class="n">5</span> Choose per boundary: JSON at the edge, gRPC between services, a broker when time matters.</div>
</div>

<style>
.take { display: flex; align-items: center; gap: 0.9rem; padding: 0.7rem 1rem; border-radius: 0.6rem;
  border: 1px solid color-mix(in srgb, var(--se-color-primary) 30%, transparent);
  background: color-mix(in srgb, var(--se-color-primary) 6%, transparent); }
.take .n { flex: none; width: 2rem; height: 2rem; border-radius: 50%; display: flex; align-items: center; justify-content: center;
  font-weight: 800; color: #fff; background: var(--se-color-primary); }
</style>

<!--
Recap, don't dwell. Then the closing line and questions.
-->

---
layout: section
---

# REST: here are resources, go fetch them.<br>gRPC: here’s a service. Call it, <span class="text-primary">or stay on the line.</span>

<style>
.slidev-layout.section h1 { font-size: 2.8rem; line-height: 1.25; }
</style>

<!--
Closing line. Thank you. Questions? Leave the auction running on the second
screen during Q&A: people will keep bidding.
-->

---
layout: section
---

# Backup

<!--
Only if asked during Q&A.
-->

---

# “Why bother, if my frontend can’t use it?”

<p class="text-lg mt-2">Sometimes you shouldn’t. Simple CRUD from browser to one backend? REST is fine. It pays off <strong>behind</strong> the BFF: fan-out, polyglot services, real-time pipelines.</p>

<div class="path mt-6">
  <div v-click class="step"><span class="n">1</span><b>Start simple</b><p>Unary gRPC behind the BFF. The browser keeps its JSON.</p></div>
  <div v-click class="step"><span class="n">2</span><b>Need real-time?</b><p>The backend streams, the BFF translates to SSE or WebSockets. Tonight’s demo.</p></div>
  <div v-click class="step"><span class="n">3</span><b>Advanced</b><p>The BFF keeps a live snapshot from streams; the frontend polls the BFF cheaply. Costs: a stateful BFF.</p></div>
</div>

<!--
Backup: jump here (press g) when someone asks. The challenge you'll get: "the frontend still polls, so why invest?" Answer:
because it cleans up service-to-service contracts, and each boundary gets the
right protocol. Present option 3 as an option, not a default: statefulness
in the BFF is a real cost.
-->

<style>
.path { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; }
.step { position: relative; padding: 1rem 1.1rem 1rem; border-radius: 0.75rem; border: 1px solid color-mix(in srgb, var(--se-color-primary) 30%, transparent);
  background: color-mix(in srgb, var(--se-color-primary) 5%, transparent); }
.step .n { display: inline-flex; width: 1.8rem; height: 1.8rem; border-radius: 50%; align-items: center; justify-content: center;
  font-weight: 800; color: #fff; background: var(--se-color-primary); margin-right: 0.5rem; }
.step b { font-size: 1.1rem; }
.step p { font-size: 0.92rem; margin-top: 0.5rem; }
</style>

---

# Sharp edges, for next time

<div class="edges mt-4">
  <div class="edge"><b>Deadlines</b> gRPC has no default timeout. Every BFF call sets one: <code>deadline: 2s</code>.</div>
  <div class="edge"><b>Cancellation</b> Pass the token on, so hung-up clients stop costing you. <code>ct</code> → stream.</div>
  <div class="edge"><b>Slow consumers</b> A watcher that falls behind skips events. Safe here, because every event carries full state.</div>
  <div class="edge"><b>Keepalives</b> Long-lived streams die silently on Wi-Fi and NAT. The BFF pings every 20s.</div>
  <div class="edge"><b>Load balancing</b> One long-lived HTTP/2 connection defeats per-connection balancing. Balance per call (L7).</div>
  <div class="edge"><b>Status codes</b> gRPC codes aren’t HTTP codes. The boundary maps them: <code>GrpcErrors.cs</code>.</div>
</div>

<p class="mt-5 text-center text-lg opacity-75">Every one of these is handled somewhere in the demo repo. Go read it.</p>

<!--
Backup: jump here (press g) if someone asks about production concerns. Name
them, point at where the repo handles each, promise a follow-up. No deep dive.
-->

<style>
.edges { display: grid; grid-template-columns: 1fr 1fr; gap: 0.6rem 1rem; }
.edge { font-size: 0.95rem; padding: 0.55rem 0.8rem; border-left: 3px solid var(--se-color-primary);
  background: color-mix(in srgb, var(--se-color-primary) 6%, transparent); border-radius: 0 0.5rem 0.5rem 0; }
.edge b { display: block; color: var(--se-color-primary); font-weight: 800; }
</style>
