using Dignitas.Auction.V1;

namespace Bff;

// The browser gets plain JSON shaped for the UI, not the protobuf messages themselves.

public record PlaceBidBody(string LotId, string Bidder, long Amount);

public record LotDto(string Id, string Title, string Description, string Emoji, long StartingPrice, long MinIncrement)
{
    public static LotDto From(Lot lot) =>
        new(lot.Id, lot.Title, lot.Description, lot.Emoji, lot.StartingPrice, lot.MinIncrement);
}

public record BidDto(string Bidder, long Amount, DateTimeOffset PlacedAt)
{
    public static BidDto From(Bid bid) => new(bid.Bidder, bid.Amount, bid.PlacedAt.ToDateTimeOffset());
}

public record AuctionDto(
    LotDto Lot,
    string Status,
    BidDto? HighestBid,
    IReadOnlyList<BidDto> RecentBids,
    int BidCount,
    long RemainingMs,
    int Watchers)
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
        auction.Watchers);
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
            _ => throw new ArgumentOutOfRangeException(nameof(message), message.Kind, "Unknown event kind"),
        },
        AuctionDto.From(message.Auction));
}
