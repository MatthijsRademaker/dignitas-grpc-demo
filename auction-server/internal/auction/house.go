// Package auction runs the live auction: one open lot at a time, bids, and a
// fan-out of every change to all watchers.
package auction

import (
	"context"
	"errors"
	"fmt"
	"math"
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

	// Load is averaged over loadWindow. Streams hear about a change in polling traffic at
	// most every loadEvery, so the projector shows the room flipping without flooding anyone.
	loadWindow = 10 * time.Second
	loadEvery  = 2 * time.Second
	// A poller that has not polled for this long is forgotten.
	pollerTTL = time.Minute
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

	// version counts changes to the auction, so a poll can tell whether it brought news.
	version      uint64
	pollers      map[string]poller
	polls        meter
	emptyPolls   meter
	pushes       meter
	lastLoad     *auctionv1.ServerLoad // as last announced to streams
	lastLoadSent time.Time
}

type poller struct {
	version uint64 // the version this poller saw last
	at      time.Time
}

func NewHouse(cfg Config) *House {
	if len(cfg.Lots) == 0 {
		panic("auction: at least one lot is required")
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	h := &House{
		cfg:        cfg,
		watchers:   map[uint64]chan *auctionv1.WatchAuctionResponse{},
		pollers:    map[string]poller{},
		polls:      meter{window: loadWindow},
		emptyPolls: meter{window: loadWindow},
		pushes:     meter{window: loadWindow},
		lastLoad:   &auctionv1.ServerLoad{},
	}
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

// Tick closes the open lot or opens the next one when its deadline has passed, and
// tells streams when polling traffic has changed.
func (h *House) Tick() {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := h.cfg.Now()
	h.announceLoad(now)
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

// Poll is Snapshot for a polling client, identified by its address. It records the
// poll, and whether it found anything the same client had not seen yet.
func (h *House) Poll(client string) *auctionv1.Auction {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := h.cfg.Now()
	last, known := h.pollers[client]
	h.polls.mark(now)
	if known && last.version == h.version {
		h.emptyPolls.mark(now)
	}
	h.pollers[client] = poller{version: h.version, at: now}
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
	h.pushes.mark(h.cfg.Now())
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
	now := h.cfg.Now()
	// Copies, not the stored bids: age_ms is different in every snapshot, and earlier
	// snapshots may still be marshalled by other goroutines.
	withAge := func(bid *auctionv1.Bid) *auctionv1.Bid {
		return &auctionv1.Bid{Bidder: bid.Bidder, Amount: bid.Amount, PlacedAt: bid.PlacedAt,
			AgeMs: max(0, now.Sub(bid.PlacedAt.AsTime()).Milliseconds())}
	}
	recent := make([]*auctionv1.Bid, len(h.recent))
	for i, bid := range h.recent {
		recent[i] = withAge(bid)
	}
	var highest *auctionv1.Bid
	if h.highest != nil {
		highest = withAge(h.highest)
	}
	return &auctionv1.Auction{
		Lot:         h.cfg.Lots[h.lotIndex],
		Status:      h.status,
		HighestBid:  highest,
		RecentBids:  recent,
		BidCount:    h.bidCount,
		RemainingMs: max(0, h.deadline.Sub(now).Milliseconds()),
		Watchers:    int32(len(h.watchers)),
		Load:        h.load(now),
	}
}

func (h *House) load(now time.Time) *auctionv1.ServerLoad {
	load := &auctionv1.ServerLoad{
		PollsPerSecond:  h.polls.perSecond(now),
		PushesPerSecond: h.pushes.perSecond(now),
	}
	if polls := h.polls.count(now); polls > 0 {
		load.EmptyPollRatio = float64(h.emptyPolls.count(now)) / float64(polls)
	}
	return load
}

// announceLoad pushes a LOAD_CHANGED event when the poll rate has clearly changed since
// streams last heard about it: someone started or stopped polling. The other figures ride
// along, but never trigger one on their own: they jitter with every poll, and pushing
// about pushes would feed itself.
func (h *House) announceLoad(now time.Time) {
	for client, p := range h.pollers {
		if now.Sub(p.at) > pollerTTL {
			delete(h.pollers, client)
		}
	}
	if len(h.watchers) == 0 || now.Sub(h.lastLoadSent) < loadEvery {
		return
	}
	if from, to := h.lastLoad.PollsPerSecond, h.polls.perSecond(now); !clearlyDifferent(from, to) {
		return
	}
	h.broadcast(auctionv1.WatchAuctionResponse_KIND_LOAD_CHANGED)
}

func clearlyDifferent(from, to float64) bool {
	if from == 0 || to == 0 {
		return from != to
	}
	return math.Abs(to-from)/from > 0.25
}

func (h *House) broadcast(kind auctionv1.WatchAuctionResponse_Kind) {
	h.broadcastExcept(^uint64(0), kind)
}

func (h *House) broadcastExcept(skip uint64, kind auctionv1.WatchAuctionResponse_Kind) {
	now := h.cfg.Now()
	// Load changes are not news about the auction itself: a poll that only finds a new
	// load figure still counts as empty.
	if kind != auctionv1.WatchAuctionResponse_KIND_LOAD_CHANGED {
		h.version++
	}
	event := &auctionv1.WatchAuctionResponse{Kind: kind, Auction: h.snapshot()}
	h.lastLoad, h.lastLoadSent = event.Auction.Load, now
	for id, events := range h.watchers {
		if id == skip {
			continue
		}
		select {
		case events <- event:
			h.pushes.mark(now)
		default: // watcher is behind; it catches up on the next event
		}
	}
}
