package auction

import auctionv1 "github.com/dignitas/dignitas-grpc-demo/auction-server/gen/auction/v1"

// Catalog is the rotation of lots the house auctions off, in order, forever.
var Catalog = []*auctionv1.Lot{
	{Id: "rubber-duck", Emoji: "🦆", Title: "Senior Rubber Duck", Description: "Has reviewed ten thousand pull requests. Approved none of them.", StartingPrice: 5, MinIncrement: 5},
	{Id: "http11-connection", Emoji: "🔌", Title: "A slightly used HTTP/1.1 connection", Description: "One request at a time. Lovingly kept alive since 1997.", StartingPrice: 1, MinIncrement: 1},
	{Id: "keyboard", Emoji: "⌨️", Title: "Mechanical keyboard, Cherry MX Blue", Description: "Your open-plan office neighbours will adore it.", StartingPrice: 40, MinIncrement: 10},
	{Id: "friday-deploy", Emoji: "🚀", Title: "Friday afternoon deploy permit", Description: "Valid for one production release after 16:00 on a Friday. No refunds.", StartingPrice: 50, MinIncrement: 10},
	{Id: "coffee-machine", Emoji: "☕", Title: "The good coffee machine from the 3rd floor", Description: "The one that actually does oat milk properly.", StartingPrice: 100, MinIncrement: 25},
	{Id: "rfc-9113", Emoji: "📜", Title: "Printed copy of RFC 9113 (HTTP/2)", Description: "Bedtime reading about streams, frames and flow control.", StartingPrice: 10, MinIncrement: 5},
}
