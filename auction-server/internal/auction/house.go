// Package auction runs the live auction: one open lot at a time, bids, and a
// fan-out of every change to all watchers.
package auction

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	auctionv1 "github.com/dignitas/dignitas-grpc-demo/auction-server/gen/auction/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	// ErrInvalidBid means the request itself is malformed.
	ErrInvalidBid = errors.New("invalid bid")
	// ErrBidRejected means the request is fine, but the auction state does not allow it.
	ErrBidRejected = errors.New("bid rejected")
)

const (
	recentBidsCap = 12
	// Each watcher gets a small buffer. Every event carries the full state, so a
	// watcher that falls behind can safely skip events: the next one heals it.
	watcherBuffer = 32
)

type Config struct {
	Lots         []*auctionv1.Lot
	LotDuration  time.Duration
	SnipeWindow  time.Duration // a bid inside this window extends the lot to now+SnipeWindow
	Intermission time.Duration // pause between a lot closing and the next one opening
	Now          func() time.Time
}

type House struct {
	cfg Config

	mu        sync.Mutex
	lotIndex  int
	status    auctionv1.LotStatus
	deadline  time.Time // open: closing time; sold/unsold: next lot opens
	highest   *auctionv1.Bid
	recent    []*auctionv1.Bid
	bidCount  int32
	watchers  map[uint64]chan *auctionv1.WatchAuctionResponse
	nextWatch uint64
}

func NewHouse(cfg Config) *House {
	if len(cfg.Lots) == 0 {
		panic("auction: at least one lot is required")
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	h := &House{cfg: cfg, watchers: map[uint64]chan *auctionv1.WatchAuctionResponse{}}
	h.openLot(0)
	return h
}

// Run advances the auction clock until ctx is cancelled.
func (h *House) Run(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.Tick()
		}
	}
}

// Tick closes the open lot or opens the next one when its deadline has passed.
func (h *House) Tick() {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := h.cfg.Now()
	if now.Before(h.deadline) {
		return
	}
	if h.status == auctionv1.LotStatus_LOT_STATUS_OPEN {
		h.status = auctionv1.LotStatus_LOT_STATUS_UNSOLD
		if h.highest != nil {
			h.status = auctionv1.LotStatus_LOT_STATUS_SOLD
		}
		h.deadline = now.Add(h.cfg.Intermission)
		h.broadcast(auctionv1.WatchAuctionResponse_KIND_LOT_CLOSED)
		return
	}
	h.openLot((h.lotIndex + 1) % len(h.cfg.Lots))
	h.broadcast(auctionv1.WatchAuctionResponse_KIND_LOT_OPENED)
}

func (h *House) Snapshot() *auctionv1.Auction {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.snapshot()
}

func (h *House) PlaceBid(lotID, bidder string, amount int64) (*auctionv1.Auction, error) {
	if bidder == "" {
		return nil, fmt.Errorf("%w: bidder is required", ErrInvalidBid)
	}
	if amount <= 0 {
		return nil, fmt.Errorf("%w: amount must be positive", ErrInvalidBid)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	lot := h.cfg.Lots[h.lotIndex]
	if lotID != lot.Id || h.status != auctionv1.LotStatus_LOT_STATUS_OPEN {
		return nil, fmt.Errorf("%w: lot %q is not open for bidding", ErrBidRejected, lotID)
	}
	if minimum := h.minimumBid(); amount < minimum {
		return nil, fmt.Errorf("%w: minimum bid is €%d", ErrBidRejected, minimum)
	}

	now := h.cfg.Now()
	bid := &auctionv1.Bid{Bidder: bidder, Amount: amount, PlacedAt: timestamppb.New(now)}
	h.highest = bid
	h.recent = append([]*auctionv1.Bid{bid}, h.recent...)
	if len(h.recent) > recentBidsCap {
		h.recent = h.recent[:recentBidsCap]
	}
	h.bidCount++
	if snipeDeadline := now.Add(h.cfg.SnipeWindow); h.deadline.Before(snipeDeadline) {
		h.deadline = snipeDeadline
	}

	h.broadcast(auctionv1.WatchAuctionResponse_KIND_BID_PLACED)
	return h.snapshot(), nil
}

// MinimumBid is the lowest amount PlaceBid currently accepts.
func (h *House) MinimumBid() int64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.minimumBid()
}

// Watch registers a watcher. The first event on the channel is always a
// snapshot. Call the returned function to stop watching.
func (h *House) Watch() (<-chan *auctionv1.WatchAuctionResponse, func()) {
	h.mu.Lock()
	defer h.mu.Unlock()

	id := h.nextWatch
	h.nextWatch++
	events := make(chan *auctionv1.WatchAuctionResponse, watcherBuffer)
	h.watchers[id] = events
	events <- &auctionv1.WatchAuctionResponse{Kind: auctionv1.WatchAuctionResponse_KIND_SNAPSHOT, Auction: h.snapshot()}
	h.broadcastExcept(id, auctionv1.WatchAuctionResponse_KIND_WATCHERS_CHANGED)

	return events, func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		delete(h.watchers, id)
		h.broadcast(auctionv1.WatchAuctionResponse_KIND_WATCHERS_CHANGED)
	}
}

func (h *House) openLot(index int) {
	h.lotIndex = index
	h.status = auctionv1.LotStatus_LOT_STATUS_OPEN
	h.deadline = h.cfg.Now().Add(h.cfg.LotDuration)
	h.highest = nil
	h.recent = nil
	h.bidCount = 0
}

func (h *House) minimumBid() int64 {
	if h.highest == nil {
		return h.cfg.Lots[h.lotIndex].StartingPrice
	}
	return h.highest.Amount + h.cfg.Lots[h.lotIndex].MinIncrement
}

func (h *House) snapshot() *auctionv1.Auction {
	return &auctionv1.Auction{
		Lot:         h.cfg.Lots[h.lotIndex],
		Status:      h.status,
		HighestBid:  h.highest,
		RecentBids:  h.recent,
		BidCount:    h.bidCount,
		RemainingMs: max(0, h.deadline.Sub(h.cfg.Now()).Milliseconds()),
		Watchers:    int32(len(h.watchers)),
	}
}

func (h *House) broadcast(kind auctionv1.WatchAuctionResponse_Kind) {
	h.broadcastExcept(^uint64(0), kind)
}

func (h *House) broadcastExcept(skip uint64, kind auctionv1.WatchAuctionResponse_Kind) {
	event := &auctionv1.WatchAuctionResponse{Kind: kind, Auction: h.snapshot()}
	for id, events := range h.watchers {
		if id == skip {
			continue
		}
		select {
		case events <- event:
		default: // watcher is behind; it catches up on the next event
		}
	}
}
