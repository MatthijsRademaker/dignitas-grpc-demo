---
theme: '@dignitas/slidev-theme'
title: "gRPC Auction House"
info: |
  ## gRPC Auction House
  Contracts, codegen and streams
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
  <div class="slot">Why gRPC</div>
  <div class="slot">Reading a <code>.proto</code> file</div>
  <div class="slot">Under the hood: why it is fast</div>
  <div class="slot hi">Workshop: implement <code>PlaceBid</code></div>
  <div class="slot hi">Live auction</div>
  <div class="slot">Four kinds of call and what they are for</div>
  <div class="slot">Trade-offs</div>
  <div class="slot">Questions</div>
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

  <p class="note">You design URLs, verbs and JSON. Then you write the client. By hand.</p>
</div>
<div v-click class="ask primary">
  <h3>gRPC</h3>
  <p class="sub">“Here’s a service, call its methods.”</p>

```csharp
await auction.GetAuctionAsync(new());
await auction.PlaceBidAsync(new() {
    Bidder = "Ada", Amount = 120
  });
```

  <p class="note">You design a <strong>contract</strong>. The client is generated from it.</p>
</div>
</div>

<!--
REST thinks in nouns you fetch. gRPC thinks in verbs you call: remote procedure
call, hence the name. Both are fine. The rest of the hour: when the second pays off.
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
    <p>Services, methods and messages, in one <code>.proto</code> file.</p>
  </div>
  <div v-click class="pillar">
    <div class="h">⚙️ Generated code</div>
    <p>Typed clients and servers for Go, C#, Java, TypeScript… from that one file.</p>
  </div>
  <div v-click class="pillar">
    <div class="h">🚄 HTTP/2 + binary</div>
    <p>Compact messages, many calls over one connection, and <strong>streaming</strong> built in.</p>
  </div>
</div>

<!--
One click per pillar. Born at Google, now CNCF. Don't go deep: each one gets its
own moment in the next 15 minutes.
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
Say it before they do. Browsers don't expose HTTP/2 framing and trailers, so no
native gRPC. gRPC-Web exists, with a proxy and some caveats. The usual answer: the
browser speaks HTTP/JSON to a BFF. [click] Behind it, gRPC. [click] Service to
service, polyglot: where it shines. [click] Choose per boundary. This is tonight's app.
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
</div>

<div class="gen mt-6">
  <span class="chip">auction.proto</span>
  <span class="arr">→</span>
  <span class="chip tool">buf generate</span><span class="arr">→</span><span class="chip">Go server interface</span>
  <span class="sep">·</span>
  <span class="chip tool">dotnet build</span><span class="arr">→</span><span class="chip">C# client</span>
</div>

<p v-click class="mt-5 text-xl text-center">Change the contract, and <span class="text-primary font-bold">the compiler tells you what broke.</span></p>

<!--
Codegen is part of the build: Grpc.Tools regenerates the C# client on every
dotnet build. The Go side is generated with buf and committed.
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
Smaller is nice; not parsing text is the bigger CPU win. Real bytes, from the Go code.
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

<div class="flex justify-center mt-2">
  <RoomHub />
</div>

<!--
Hub and spokes: one auction server on my machine, a frontend and BFF on yours,
pointed at my IP. Every bid from every laptop lands in the same auction.
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
  <div class="check">Open <code>http://localhost:8080</code>: a mystery lot, with live bids.</div>
  <div class="check">Try to bid: <span class="text-primary font-bold">“Workshop time”</span>. Fix that, and you find out what’s under the cloth.</div>
</div>

<p class="ip">Presenter: <code>AUCTION_HOST=http://192.168.1.20:50051</code></p>

<!--
Update the IP right before the session (ip -4 addr / ipconfig). --watch rebuilds
the BFF on every save: save, wait a few seconds. No refresh needed: the cloth
lifts by itself once PlaceBid works. Seeing the lot at all proves GetAuction works
through your BFF: that call is prebuilt, and it's your template. Don't spoil
what's under the cloth.
-->

<style>
.steps pre { font-size: 1.15rem !important; }
.check { padding: 0.6rem 0.9rem; border-radius: 0.6rem; border: 1px solid color-mix(in srgb, currentColor 15%, transparent); }
.ip { margin-top: 1.4rem; text-align: center; font-size: 1.6rem; font-weight: 800; }
.ip code { color: var(--se-color-primary); }
</style>

---

# Your task: `bff/AuctionEndpoints.Workshop.cs`

<div class="task-grid">
<div>

```csharp
api.MapPost("/bids", (PlaceBidBody body,
    // TODO inject RPC Service
    // Generated from Protobuf contract
    CancellationToken ct) =>
{
    // make this place a real bid
});
```

</div>
<div class="hints">
  <div class="hint done"><b>Done when</b>
    <ul>
      <li>your bid lands on everyone’s screen</li>
      <li>a bid that’s too low shows the reason, not a crash</li>
      <li>the call can’t hang forever</li>
    </ul>
  </div>
  <div class="hint"><b>Stuck?</b> <code>workshop/hints/1.md</code>, then 2, then 3. Each gives away a bit more.</div>
</div>
</div>

<p class="text-center mt-4 text-lg opacity-70">⏱ 15 minutes · pairing welcome · the full guide: <code>workshop/README.md</code></p>

<!--
About ten lines, and no recipe on screen on purpose: this is hard mode. Walk the
room and send stuck people to the next hint, not the answer. Usual snags: getting the
client injected (hint 2 names the type), no async on the lambda, AUCTION_HOST typos,
my firewall. "Can't hang forever" means a deadline: gRPC has none by default. A cheer from a corner means someone just met the goose.
-->

<style>
.task-grid { display: grid; grid-template-columns: 1.3fr 1fr; gap: 1.2rem; margin-top: 0.6rem; align-items: start; }
.task-grid pre { font-size: 0.85rem !important; }
.hints { display: flex; flex-direction: column; gap: 0.5rem; }
.hint { font-size: 0.92rem; padding: 0.45rem 0.7rem; border-radius: 0.5rem; background: color-mix(in srgb, var(--se-color-primary) 7%, transparent); }
.hint b { color: var(--se-color-primary); margin-right: 0.3rem; }
.hint.done ul { margin: 0.3rem 0 0; padding-left: 1.1rem; list-style: disc; }
.hint.done li { margin: 0.1rem 0; }
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

<p v-click class="mt-4 text-lg text-center">Didn’t finish in time? <code>cp workshop/solution/AuctionEndpoints.Workshop.cs bff/</code> and you’re in.</p>

<!--
Walk the highlights: the call, the deadline (gRPC has no default timeout), the
JSON for the browser, the status-code translation. Hands up if you've seen what's
under the cloth. [click] Leave the catch-up line up for a moment: --watch picks up
the copied file, so nobody sits out the auction, or misses the goose.
-->

---
layout: section
---

# Live auction

---

# Everyone: bid!

<div class="cue">
  <div class="bar">polling mode</div>
  <div class="watch">
    <strong>Watch for:</strong> bids landing <span class="text-primary font-bold">a second late</span>,
    and golden eggs dropping <span class="text-primary font-bold">in batches</span>: one batch per poll.
  </div>
</div>

<!--
A round or two, projector on the stage view (localhost:8080/?view=stage). Point
at the room chart: every laptop's polls, most finding nothing new. The delay should
feel slightly annoying: that's the point. Late bids extend the clock, so the end of
each round gets frantic, and the winner goes up on the wall.
Then the polling dilemma: switch the projector to 0.5s. Bids arrive sooner, and
calls per second jump. Freshness, paid for in load.
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
Top: what happened. Middle: polling every 2s, mostly grey, every bid late.
[click] Bottom: one request that stays open; each bid arrives as it happens.
"Polling an order status every 2 seconds? You wanted a subscription."
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
The third method I skipped earlier. That's the whole delta between "ask me"
and "tell me".
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

<!--
The BFF turns the gRPC stream into SSE, which every browser speaks. Lightly
abbreviated; the real thing is already in your bff/AuctionEndpoints.cs.
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
    <strong>Watch for:</strong> calls per second <span class="text-primary font-bold">sliding to zero</span>,
    “bids seen after” hitting <span class="text-primary font-bold">0.00s</span>, and the eggs dropping <span class="text-primary font-bold">one by one</span>.
  </div>
</div>

<!--
The payoff. Flip the stage view to streaming too. Calls per second slope down
over ten seconds as the room flips, open streams climb, all in one 90s window.
Run a round or two like this. Flash the server terminal: "stream opened" from
every IP in the room.
-->

<style>
.cue { margin: 0.4rem auto 0; border: 2px solid color-mix(in srgb, var(--se-color-primary) 35%, transparent); border-radius: 0.75rem; overflow: hidden; max-width: 820px; }
.cue .msg { padding: 0.7rem 1rem; font-size: 1.2rem; font-weight: 800; }
.cue .watch { padding: 0.7rem 1rem; font-size: 0.95rem; border-top: 1px solid color-mix(in srgb, var(--se-color-primary) 25%, transparent);
  background: color-mix(in srgb, var(--se-color-primary) 5%, transparent); }
</style>

---

# What we just built

<table class="built">
  <thead>
    <tr><th></th><th>What you did</th><th>What gRPC did</th></tr>
  </thead>
  <tbody>
    <tr v-click>
      <td>Contract</td><td>Read <code>auction.proto</code></td>
      <td>Generated the Go server and your C# client from it</td>
    </tr>
    <tr v-click>
      <td>Unary</td><td>Wrote <code>POST /bids</code></td>
      <td>One typed <code>PlaceBid</code> call, a deadline, status codes mapped to HTTP</td>
    </tr>
    <tr v-click>
      <td>Polling</td><td>Asked <code>GetAuction</code> every 2s</td>
      <td>Calls from every browser, most of them finding nothing new</td>
    </tr>
    <tr v-click>
      <td>Streaming</td><td>Flipped one switch</td>
      <td>One <code>WatchAuction</code> stream per browser, bids pushed as they land</td>
    </tr>
  </tbody>
</table>

<p v-click class="mt-5 text-center text-lg">Two languages, one contract, <span class="text-primary font-bold">no hand-written client</span>.</p>

<!--
Quick recap while the auction keeps running. Each row is something they did
themselves in the last half hour. Don't dwell: the next section names the shapes.
-->

<style>
.built { margin-top: 1.2rem; width: 100%; font-size: 0.95rem; border-collapse: separate; border-spacing: 0 0.35rem; }
.built th { text-align: left; font-size: 0.75rem; font-weight: 800; letter-spacing: 0.05em; text-transform: uppercase; opacity: 0.6; padding: 0 0.7rem; }
.built td { padding: 0.55rem 0.7rem; background: color-mix(in srgb, currentColor 4%, transparent); }
.built td:first-child { font-weight: 800; border-radius: 0.5rem 0 0 0.5rem; }
.built td:last-child { border-radius: 0 0.5rem 0.5rem 0; background: color-mix(in srgb, var(--se-color-primary) 9%, transparent); }
.built code { font-size: 0.82rem; }
</style>

---
layout: section
---

# Four kinds of call

---
clicks: 3
---

# One unary, three streams

<StreamShapes class="mt-2" />

<!--
Unary: you built it, and it is not a stream: one message each way. [click] Server
streaming: you just watched it.
[click] Client streaming: many up, one back, like a chunked upload.
[click] Bidirectional: both sides talk whenever they like.
-->

---
clicks: 3
---

# How would you do this with REST?

<RestVsGrpc class="mt-4" />

<p v-click="3" class="mt-5 text-lg text-center">
Need it delivered tomorrow, or replayed for audit? <span class="font-bold">That’s still a broker.</span>
</p>

<!--
The obvious question: how would you build this without gRPC? At the edge
the answer is the usual toolbox, and we used it too: the BFF turns the stream
into SSE. [click] Behind the BFF it gets ad hoc: polling, webhooks, or a broker.
[click] With gRPC it's one keyword in the contract, and the generated code on
both sides handles cancellation and deadlines. [click] A stream is a direct pipe:
both ends connected at the same time. Durability and replay are a broker's job.
-->

---
layout: section
---

# Trade-offs

---

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

<!--
Both columns, no hedging. "Should we do this everywhere?" No. Evolvable APIs: we
did that today. round, lot_duration_ms and winners are fields 9 to 11, and no BFF
broke. Reflection is on: grpcurl -plaintext <ip>:50051 list works right now.
-->

<style>
.side { border: 2px solid color-mix(in srgb, currentColor 15%, transparent); border-radius: 0.75rem; padding: 0.9rem 1.2rem; }
.side.good { border-color: var(--se-color-primary); background: color-mix(in srgb, var(--se-color-primary) 6%, transparent); }
.side h3 { font-weight: 800; font-size: 1.2rem; margin-bottom: 0.5rem; }
.side ul { padding-left: 1.1rem; display: flex; flex-direction: column; gap: 0.35rem; font-size: 0.98rem; }
</style>

---
layout: section
---

# REST: here are resources, go fetch them.<br>gRPC: here’s a service. Call it, <span class="text-primary">or stay on the line.</span>

<style>
.slidev-layout.section h1 { font-size: 2.8rem; line-height: 1.25; }
</style>

<!--
Thank you. Questions? Leave the stage view up during Q&A: people will keep
bidding for the goose.
-->

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
Backup (press g) for production questions. Name them, point at the repo,
promise a follow-up. No deep dive.
-->

<style>
.edges { display: grid; grid-template-columns: 1fr 1fr; gap: 0.6rem 1rem; }
.edge { font-size: 0.95rem; padding: 0.55rem 0.8rem; border-left: 3px solid var(--se-color-primary);
  background: color-mix(in srgb, var(--se-color-primary) 6%, transparent); border-radius: 0 0.5rem 0.5rem 0; }
.edge b { display: block; color: var(--se-color-primary); font-weight: 800; }
</style>
