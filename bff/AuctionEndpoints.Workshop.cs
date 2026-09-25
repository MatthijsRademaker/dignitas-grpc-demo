namespace Bff;

using Dignitas.Auction.V1;
using Grpc.Core;

public static partial class AuctionEndpoints
{
    private static void MapPlaceBid(RouteGroupBuilder api)
    {
        // ------------------------------------------------------------------------------------------
        // WORKSHOP: make this place a real bid. The frontend POSTs { lotId, bidder, amount } here.
        // GET /auction in AuctionEndpoints.cs is your template. When you're done, and hints if
        // you're stuck: workshop/README.md, step 4.
        // ------------------------------------------------------------------------------------------
        api.MapPost(
            "/bids",
            (
                PlaceBidBody body,
                // TODO inject RPC Service
                // Generated from Protobuf contract
                CancellationToken ct
            ) =>
            {
                return Results.Problem(
                    title: "Not implemented yet",
                    detail: "PlaceBid is waiting for you in bff/AuctionEndpoints.Workshop.cs",
                    statusCode: StatusCodes.Status501NotImplemented
                );
            }
        );
    }
}
