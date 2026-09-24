package auction

import (
	"context"
	"math/rand/v2"
	"time"
)

var botNames = []string{"bot-ada", "bot-grace", "bot-linus", "bot-barbara", "bot-ken", "bot-margaret"}

// RunBots starts n simulated bidders so the auction feels alive when you are
// rehearsing alone. They bid in-process, not over gRPC.
func RunBots(ctx context.Context, house *House, n int) {
	for i := range n {
		go func() {
			name := botNames[i%len(botNames)]
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Duration(3000+rand.IntN(7000)) * time.Millisecond):
					lot := house.Snapshot().GetLot()
					// Rejections (lot closed, outbid in the meantime) are expected here.
					_, _ = house.PlaceBid(lot.GetId(), name, house.MinimumBid()+rand.Int64N(3)*lot.GetMinIncrement())
				}
			}
		}()
	}
}
