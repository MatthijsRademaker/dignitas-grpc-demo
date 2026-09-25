using System.Runtime.CompilerServices;
using Dignitas.Auction.V1;
using Grpc.Core;

namespace Bff;

public static class AuctionEndpoints
{
    public static void MapAuctionEndpoints(this IEndpointRouteBuilder app)
    {
        var api = app.MapGroup("/api");

        // Unary gRPC -> JSON. In polling mode the frontend calls this every 2 seconds.
        api.MapGet("/auction", async (AuctionService.AuctionServiceClient auction, CancellationToken ct) =>
        {
            try
            {
                var response = await auction.GetAuctionAsync(
                    new GetAuctionRequest(),
                    deadline: DateTime.UtcNow.AddSeconds(2),
                    cancellationToken: ct);
                return Results.Ok(AuctionDto.From(response.Auction));
            }
            catch (RpcException ex)
            {
                return ex.ToProblem();
            }
        });

        // ------------------------------------------------------------------------------------------
        // WORKSHOP: make this place a real bid. The frontend POSTs { lotId, bidder, amount } here.
        // GET /auction above is your template. When you're done, and hints if you're stuck:
        // workshop/README.md, step 4.
        // ------------------------------------------------------------------------------------------
        api.MapPost("/bids", (PlaceBidBody body, AuctionService.AuctionServiceClient auction, CancellationToken ct) =>
        {
            return Results.Problem(
                title: "Not implemented yet",
                detail: "PlaceBid is waiting for you in bff/AuctionEndpoints.cs",
                statusCode: StatusCodes.Status501NotImplemented);
        });

        // Server-streaming gRPC -> Server-Sent Events: stream translation at the boundary.
        api.MapGet("/auction/stream", (AuctionService.AuctionServiceClient auction, CancellationToken ct) =>
            TypedResults.ServerSentEvents(Watch(auction, ct), eventType: "auction"));
    }

    private static async IAsyncEnumerable<AuctionEventDto> Watch(
        AuctionService.AuctionServiceClient auction,
        [EnumeratorCancellation] CancellationToken ct)
    {
        // `ct` fires when the browser disconnects. Passing it on cancels the gRPC stream as well,
        // so the auction server stops sending to someone who is no longer listening.
        using var call = auction.WatchAuction(new WatchAuctionRequest(), cancellationToken: ct);
        await foreach (var message in call.ResponseStream.ReadAllAsync(ct))
        {
            yield return AuctionEventDto.From(message);
        }
    }
}
