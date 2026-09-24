package auction

import auctionv1 "github.com/dignitas/dignitas-grpc-demo/auction-server/gen/auction/v1"

// Catalog is what the house auctions off: one prize, sold again and again in rounds.
var Catalog = []*auctionv1.Lot{
	{Id: "golden-goose", Emoji: "🦆", Title: "The Golden Goose", Description: "Lays one golden egg for every bid it hears about. Keep it streaming.", StartingPrice: 5, MinIncrement: 5},
}
