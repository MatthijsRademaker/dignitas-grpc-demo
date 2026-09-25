using Dignitas.Auction.V1;

namespace Bff;

// The browser gets plain JSON shaped for the UI, not the protobuf messages themselves.

public record PlaceBidBody(string LotId, string Bidder, long Amount);

public record LotDto(string Id, string Title, string Description, string Emoji, long StartingPrice, long MinIncrement)
{
    public static LotDto From(Lot lot) =>
        new(lot.Id, lot.Title, lot.Description, lot.Emoji, lot.StartingPrice, lot.MinIncrement);
}

public record BidDto(string Bidder, long Amount, DateTimeOffset PlacedAt, long AgeMs)
{
    public static BidDto From(Bid bid) => new(bid.Bidder, bid.Amount, bid.PlacedAt.ToDateTimeOffset(), bid.AgeMs);
}

public record SaleDto(int Round, string Bidder, long Amount, DateTimeOffset SoldAt)
{
    public static SaleDto From(Sale sale) => new(sale.Round, sale.Bidder, sale.Amount, sale.SoldAt.ToDateTimeOffset());
}

public record ServerLoadDto(double PollsPerSecond, double EmptyPollRatio, double PushesPerSecond)
{
    // An older server that doesn't send load yet: report nothing rather than fail.
    public static ServerLoadDto From(ServerLoad? load) =>
        load is null ? new(0, 0, 0) : new(load.PollsPerSecond, load.EmptyPollRatio, load.PushesPerSecond);
}

public record AuctionDto(
    LotDto Lot,
    string Status,
    BidDto? HighestBid,
    IReadOnlyList<BidDto> RecentBids,
    int BidCount,
    long RemainingMs,
    int Watchers,
    ServerLoadDto Load,
    int Round,
    long LotDurationMs,
    IReadOnlyList<SaleDto> Winners)
{
    public static AuctionDto From(Auction auction) => new(
        LotDto.From(auction.Lot),
        auction.Status switch
        {
            LotStatus.Open => "open",
            LotStatus.Sold => "sold",
            LotStatus.Unsold => "unsold",
            _ => throw new ArgumentOutOfRangeException(nameof(auction), auction.Status, "Unknown lot status"),
        },
        // Message fields are null when unset: that is how proto3 says "no bid yet".
        auction.HighestBid is null ? null : BidDto.From(auction.HighestBid),
        auction.RecentBids.Select(BidDto.From).ToList(),
        auction.BidCount,
        auction.RemainingMs,
        auction.Watchers,
        ServerLoadDto.From(auction.Load),
        auction.Round,
        auction.LotDurationMs,
        // Repeated fields are never null: an older server without winners just sends an empty list.
        auction.Winners.Select(SaleDto.From).ToList());
}

public record AuctionEventDto(string Kind, AuctionDto Auction)
{
    public static AuctionEventDto From(WatchAuctionResponse message) => new(
        message.Kind switch
        {
            WatchAuctionResponse.Types.Kind.Snapshot => "snapshot",
            WatchAuctionResponse.Types.Kind.BidPlaced => "bid_placed",
            WatchAuctionResponse.Types.Kind.LotOpened => "lot_opened",
            WatchAuctionResponse.Types.Kind.LotClosed => "lot_closed",
            WatchAuctionResponse.Types.Kind.WatchersChanged => "watchers_changed",
            WatchAuctionResponse.Types.Kind.LoadChanged => "load_changed",
            _ => throw new ArgumentOutOfRangeException(nameof(message), message.Kind, "Unknown event kind"),
        },
        AuctionDto.From(message.Auction));
}
