# Implement `PlaceBid`, step by step

The presenter runs one gRPC auction server for the whole room. You run a frontend and a BFF (backend for
frontend) on your laptop. Your BFF can already *read* the auction over gRPC. You'll teach it to *place bids*.

```
browser ──HTTP/JSON──▶ your BFF (.NET) ──gRPC──▶ presenter's auction-server (Go)
         localhost:8080                           http://<presenter-ip>:50051
```

About 15 minutes. Pairing welcome.

## 1. Connect

Did the [quickstart](../QUICKSTART.md)? Then you only need the `.env` edit and the last line.

```bash
git clone https://github.com/MatthijsRademaker/dignitas-grpc-demo.git && cd dignitas-grpc-demo
cp .env.example .env
# edit .env: AUCTION_HOST=http://<presenter-ip>:50051   (the IP is on the slide)
docker compose up --build --watch
```

Open <http://localhost:8080>: a mystery lot, under a cloth, with live bids. Reading already works through the
prebuilt `GET /api/auction` endpoint, which makes a unary `GetAuction` call.

Try to bid: **"Workshop time"**. That's the endpoint you're about to write.

`--watch` rebuilds the BFF every time you save a file in `bff/`. Give it a few seconds after each save.

## 2. Read the contract

Open [`proto/auction/v1/auction.proto`](../proto/auction/v1/auction.proto) and find:

- the `PlaceBid` rpc on `AuctionService`
- `PlaceBidRequest`: which fields, which types?
- `PlaceBidResponse`
- the comment above `PlaceBid`: which gRPC status codes a rejected bid returns

## 3. Find the generated client

You don't write the gRPC client. `Grpc.Tools` generates it from the proto on every build (see the `<Protobuf>` item
in [`bff/Bff.csproj`](../bff/Bff.csproj)).

- **With the .NET SDK:** run `dotnet build` in `bff/` and type `auction.` in your IDE, or look for `PlaceBidAsync`
  in `bff/obj/Debug/net10.0/auction/v1/AuctionGrpc.cs`.
- **Docker only:** trust the proto. Every `rpc Foo` becomes `FooAsync(...)`, and every message becomes a C# class
  with PascalCase properties (`lot_id` becomes `LotId`).

## 4. Write the endpoint

This is the challenge. Make `POST /bids` in [`bff/AuctionEndpoints.cs`](../bff/AuctionEndpoints.cs) place a real
bid. `GET /auction`, right above it, is the same pattern for a different rpc.

You're done when:

- your bid shows up in the feed on **everyone's** screen
- a bid that's too low shows the server's reason ("minimum bid is €…"), not a crash
- the call can't hang forever if the server disappears

Stuck? Each hint gives away a bit more, so only open the next one:

1. [Where to look](hints/1.md): no code, just directions
2. [The pieces](hints/2.md): each fragment on its own
3. [The shape](hints/3.md): the whole endpoint, with blanks

## 5. Bid

Save, wait a few seconds, and bid. Keep an eye on the cloth.

## Going further

1. **Deadlines.** Set the deadline to 1 millisecond. Which gRPC status comes back, and what does `GrpcErrors.cs`
   turn it into? (gRPC has no default timeout. Always set one.)
2. **Status codes.** Bid too low, or with an empty name. Which gRPC codes come back, and which HTTP codes does the
   browser see?
3. **Talk gRPC directly**, no BFF and no proto file: the server has reflection enabled.

   ```bash
   docker run --rm fullstorydev/grpcurl -plaintext <presenter-ip>:50051 list
   docker run --rm fullstorydev/grpcurl -plaintext <presenter-ip>:50051 auction.v1.AuctionService/GetAuction
   # take the lot id from that answer; €1 is too low on purpose
   docker run --rm fullstorydev/grpcurl -plaintext \
     -d '{"lot_id":"<lot-id>","bidder":"grpcurl","amount":1}' \
     <presenter-ip>:50051 auction.v1.AuctionService/PlaceBid
   ```

4. **Read the streaming side.** `GET /auction/stream` in the same file turns the `WatchAuction` server stream into
   Server-Sent Events. How does cancellation reach the server when you close the tab?

## Out of time?

Try it yourself first. The full file is in [`solution/AuctionEndpoints.cs`](solution/AuctionEndpoints.cs). To catch
up for the live auction:

```bash
cp workshop/solution/AuctionEndpoints.cs bff/
```

`--watch` rebuilds the BFF by itself. Running natively? Restart `dotnet run`.

<details>
<summary>Show the solution</summary>

```csharp
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
        return ex.ToProblem();
    }
});
```

</details>

## No Docker?

Needs the .NET 10 SDK and Node 22+.

```bash
cd bff && AUCTION_HOST=http://<presenter-ip>:50051 dotnet run   # http://localhost:5080
cd frontend && npm ci && npm run dev                            # http://localhost:5173
```

## Troubleshooting

| Symptom | Likely cause |
| --- | --- |
| BFF exits with `AUCTION_HOST is not set…` | `.env` is missing or has no `AUCTION_HOST` |
| Page says "Connecting…" or shows `Unavailable` | Wrong IP or port in `AUCTION_HOST`, or the Wi-Fi blocks laptop-to-laptop traffic. Test with the `grpcurl … list` command above |
| `DeadlineExceeded` on every call | Same as above: the server can't be reached in time |
| Port 8080 already in use | Change `"8080:80"` in `compose.yaml` to e.g. `"8081:80"` |
| `--watch` not recognised | Update Docker Compose, or rerun `docker compose up --build` after each save |
