package auction

import (
	"context"
	"math/rand/v2"
	"time"

	auctionv1 "github.com/dignitas/dignitas-grpc-demo/auction-server/gen/auction/v1"
)

var botNames = []string{"bot-ada", "bot-grace", "bot-linus", "bot-barbara", "bot-ken", "bot-margaret"}

// frenzyBelow is when bots stop browsing and start fighting. Above the snipe window, so
// the last-second extensions keep them in it until their budgets run out.
const frenzyBelow = 15 * time.Second

// RunBots starts n simulated bidders so the auction feels alive when you are
// rehearsing alone. They bid in-process, not over gRPC.
//
// Like people, they are quiet while the clock is long and pile in when it runs out.
// That shows both costs of polling in one round: requests that find nothing early on,
// and news that arrives late during the frenzy at the end.
func RunBots(ctx context.Context, house *House, n int) {
	for i := range n {
		go func() {
			name := botNames[i%len(botNames)]
			var round int32
			var budget int64
			for {
				auction := house.Snapshot()
				lot := auction.GetLot()
				if auction.GetRound() != round {
					// A new round: decide how far this bot goes before it drops out.
					round = auction.GetRound()
					budget = lot.GetStartingPrice() + lot.GetMinIncrement()*(15+rand.Int64N(30))
				}

				frenzy := auction.GetStatus() == auctionv1.LotStatus_LOT_STATUS_OPEN &&
					time.Duration(auction.GetRemainingMs())*time.Millisecond < frenzyBelow
				wait := time.Duration(8000+rand.IntN(10000)) * time.Millisecond
				if frenzy {
					wait = time.Duration(300+rand.IntN(1200)) * time.Millisecond
				}
				select {
				case <-ctx.Done():
					return
				case <-time.After(wait):
				}

				if house.Snapshot().GetHighestBid().GetBidder() == name {
					continue // already winning
				}
				amount := house.MinimumBid()
				if frenzy {
					amount += rand.Int64N(2) * lot.GetMinIncrement()
				}
				if amount <= budget {
					// Rejections (lot closed, outbid in the meantime) are expected here.
					_, _ = house.PlaceBid(lot.GetId(), name, amount)
				}
			}
		}()
	}
}
