# Workshop: implement `PlaceBid`

The presenter runs one central gRPC auction server. You run a frontend and a BFF (backend for
frontend) on your own laptop. Your BFF already *reads* the auction over gRPC. Your job: let it
*place bids*, so you can join the live auction.

```
browser ──HTTP/JSON──▶ your BFF (.NET) ──gRPC──▶ presenter's auction-server (Go)
         localhost:8080                           http://<presenter-ip>:50051
```

**Timebox: 15 minutes.** Pairing is encouraged.

## 0. Setup

You only need Docker (Compose 2.22 or newer, for `--watch`).

```bash
git clone <repo-url> && cd dignitas-grpc-demo
cp .env.example .env
# edit .env: AUCTION_HOST=http://<presenter-ip>:50051   (the IP is on the slide)
docker compose up --build --watch
```

Open <http://localhost:8080>. You should see the lot currently on the block, with live bids.
That already works: the prebuilt `GET /api/auction` endpoint makes a unary `GetAuction` gRPC call.

Now try to place a bid. You'll get **"Workshop time"**. That's the part you build.

`--watch` rebuilds and restarts the BFF container every time you save a file in `bff/`. Give it
a few seconds, then bid again.

## 1. Read the contract

Open [`proto/auction/v1/auction.proto`](../proto/auction/v1/auction.proto) and find:

- the `PlaceBid` rpc on `AuctionService`
- its request message, `PlaceBidRequest`: which fields, which types?
- its response message, `PlaceBidResponse`
- the comment above `PlaceBid` that says which gRPC status codes a rejected bid returns

## 2. Find the generated client

You never write the gRPC client yourself. `Grpc.Tools` generates it from the proto on every build
(see the `<Protobuf>` item in [`bff/Bff.csproj`](../bff/Bff.csproj)).

- **With the .NET SDK:** run `dotnet build` in `bff/`, then open
  `bff/obj/Debug/net10.0/auction/v1/AuctionGrpc.cs` and look for `PlaceBidAsync`. Or type
  `auction.` in your IDE and let autocomplete show you.
- **Docker only:** trust the proto. Every `rpc Foo` becomes `FooAsync(...)` on the client, and
  every message becomes a C# class with PascalCase properties (`lot_id` becomes `LotId`).

## 3. Implement the endpoint

Open [`bff/AuctionEndpoints.cs`](../bff/AuctionEndpoints.cs) and find the `WORKSHOP` block.
`GET /auction`, right above it, is your template. It does the same thing for `GetAuction`.

1. Build a `PlaceBidRequest` from `body`.
2. `await auction.PlaceBidAsync(request, deadline: ..., cancellationToken: ct)`. Remember to make
   the lambda `async`.
3. Return `Results.Ok(AuctionDto.From(response.Auction))`.
4. Catch `RpcException` and return `ex.ToProblem()`.

## 4. Try it

- Bid on the lot. Your name should appear in the bid feed on **everyone's** screen.
- Bid below the minimum. What does the frontend show?

## Stretch goals

1. **Deadlines.** Set the deadline to 1 millisecond. Which gRPC status do you get, and what does
   `GrpcErrors.cs` turn it into? (gRPC has no default timeout. Always set one.)
2. **Status codes.** Bid too low, or with an empty name. Which gRPC codes come back, and which HTTP
   codes does the browser see?
3. **Talk to gRPC directly**, no BFF, no proto file. The server has reflection enabled:

   ```bash
   docker run --rm fullstorydev/grpcurl -plaintext <presenter-ip>:50051 list
   docker run --rm fullstorydev/grpcurl -plaintext <presenter-ip>:50051 describe auction.v1.PlaceBidRequest
   docker run --rm fullstorydev/grpcurl -plaintext \
     -d '{"lot_id":"rubber-duck","bidder":"grpcurl","amount":1}' \
     <presenter-ip>:50051 auction.v1.AuctionService/PlaceBid
   ```

4. **Read the streaming side.** `GET /auction/stream` in the same file turns the `WatchAuction`
   server stream into Server-Sent Events. How does cancellation reach the server when you close
   the tab?

## Troubleshooting

| Symptom | Likely cause |
| --- | --- |
| `docker compose` says `Set AUCTION_HOST…` | `.env` is missing or has no `AUCTION_HOST` |
| Page says "Connecting…" or shows `Unavailable` | Wrong IP or port in `AUCTION_HOST`, or the Wi-Fi blocks laptop-to-laptop traffic. Test with the `grpcurl … list` command above |
| `DeadlineExceeded` on every call | Same as above: the server can't be reached in time |
| Port 8080 already in use | Change `"8080:80"` in `compose.yaml` to e.g. `"8081:80"` |
| `--watch` not recognised | Update Docker Compose, or rerun `docker compose up --build` after each save |

## No Docker? Run it natively

Needs the .NET 10 SDK and Node 22+.

```bash
cd bff && AUCTION_HOST=http://<presenter-ip>:50051 dotnet run   # http://localhost:5080
cd frontend && npm ci && npm run dev                            # http://localhost:5173
```

## Solution

Try it yourself first. The full file is in [`solution/AuctionEndpoints.cs`](solution/AuctionEndpoints.cs).

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
