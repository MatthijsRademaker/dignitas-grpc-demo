using Dignitas.Auction.V1;
using Grpc.Core;

namespace Bff;

public static partial class AuctionEndpoints
{
    private static void MapPlaceBid(RouteGroupBuilder api)
    {
        // WORKSHOP solution: unary gRPC -> JSON.
        api.MapPost(
            "/bids",
            async (
                PlaceBidBody body,
                AuctionService.AuctionServiceClient auction,
                CancellationToken ct
            ) =>
            {
                try
                {
                    var response = await auction.PlaceBidAsync(
                        new PlaceBidRequest
                        {
                            LotId = body.LotId,
                            Bidder = body.Bidder,
                            Amount = body.Amount,
                        },
                        deadline: DateTime.UtcNow.AddSeconds(2),
                        cancellationToken: ct
                    );
                    return Results.Ok(AuctionDto.From(response.Auction));
                }
                catch (RpcException ex)
                {
                    return ex.ToProblem();
                }
            }
        );
    }
}
