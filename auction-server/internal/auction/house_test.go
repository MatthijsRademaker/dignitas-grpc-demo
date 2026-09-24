package auction

import (
	"errors"
	"testing"
	"time"

	auctionv1 "github.com/dignitas/dignitas-grpc-demo/auction-server/gen/auction/v1"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time          { return c.t }
func (c *clock) advance(d time.Duration) { c.t = c.t.Add(d) }

func newTestHouse(c *clock) *House {
	return NewHouse(Config{
		Lots: []*auctionv1.Lot{
			{Id: "a", StartingPrice: 10, MinIncrement: 5},
			{Id: "b", StartingPrice: 1, MinIncrement: 1},
		},
		LotDuration:  60 * time.Second,
		SnipeWindow:  10 * time.Second,
		Intermission: 5 * time.Second,
		Now:          c.now,
	})
}

func TestPlaceBidEnforcesMinimum(t *testing.T) {
	h := newTestHouse(&clock{t: time.Unix(0, 0)})

	if _, err := h.PlaceBid("a", "ada", 9); !errors.Is(err, ErrBidRejected) {
		t.Fatalf("below starting price: got %v, want ErrBidRejected", err)
	}
	if _, err := h.PlaceBid("a", "ada", 10); err != nil {
		t.Fatalf("starting price: %v", err)
	}
	if _, err := h.PlaceBid("a", "grace", 14); !errors.Is(err, ErrBidRejected) {
		t.Fatalf("below increment: got %v, want ErrBidRejected", err)
	}
	a, err := h.PlaceBid("a", "grace", 15)
	if err != nil {
		t.Fatalf("valid raise: %v", err)
	}
	if a.GetHighestBid().GetBidder() != "grace" || a.GetBidCount() != 2 || len(a.GetRecentBids()) != 2 {
		t.Fatalf("unexpected state: %v", a)
	}
}

func TestPlaceBidValidatesInput(t *testing.T) {
	h := newTestHouse(&clock{t: time.Unix(0, 0)})

	if _, err := h.PlaceBid("a", "", 10); !errors.Is(err, ErrInvalidBid) {
		t.Fatalf("empty bidder: got %v, want ErrInvalidBid", err)
	}
	if _, err := h.PlaceBid("a", "ada", 0); !errors.Is(err, ErrInvalidBid) {
		t.Fatalf("zero amount: got %v, want ErrInvalidBid", err)
	}
	if _, err := h.PlaceBid("b", "ada", 10); !errors.Is(err, ErrBidRejected) {
		t.Fatalf("wrong lot: got %v, want ErrBidRejected", err)
	}
}

func TestLateBidExtendsLot(t *testing.T) {
	c := &clock{t: time.Unix(0, 0)}
	h := newTestHouse(c)

	c.advance(55 * time.Second)
	a, err := h.PlaceBid("a", "ada", 10)
	if err != nil {
		t.Fatal(err)
	}
	if a.GetRemainingMs() != 10_000 {
		t.Fatalf("remaining after late bid: got %dms, want 10000ms", a.GetRemainingMs())
	}
}

func TestLotLifecycle(t *testing.T) {
	c := &clock{t: time.Unix(0, 0)}
	h := newTestHouse(c)
	events, stop := h.Watch()
	defer stop()

	if e := <-events; e.GetKind() != auctionv1.WatchAuctionResponse_KIND_SNAPSHOT || e.GetAuction().GetWatchers() != 1 {
		t.Fatalf("first event: got %v", e)
	}
	if _, err := h.PlaceBid("a", "ada", 10); err != nil {
		t.Fatal(err)
	}
	if e := <-events; e.GetKind() != auctionv1.WatchAuctionResponse_KIND_BID_PLACED {
		t.Fatalf("after bid: got %v", e.GetKind())
	}

	c.advance(60 * time.Second)
	h.Tick()
	if e := <-events; e.GetKind() != auctionv1.WatchAuctionResponse_KIND_LOT_CLOSED || e.GetAuction().GetStatus() != auctionv1.LotStatus_LOT_STATUS_SOLD {
		t.Fatalf("after close: got %v", e)
	}
	if _, err := h.PlaceBid("a", "grace", 100); !errors.Is(err, ErrBidRejected) {
		t.Fatalf("bid on closed lot: got %v, want ErrBidRejected", err)
	}

	c.advance(5 * time.Second)
	h.Tick()
	e := <-events
	if e.GetKind() != auctionv1.WatchAuctionResponse_KIND_LOT_OPENED || e.GetAuction().GetLot().GetId() != "b" || e.GetAuction().GetHighestBid() != nil {
		t.Fatalf("after reopen: got %v", e)
	}
}

func TestPollTellsNewsFromNoNews(t *testing.T) {
	c := &clock{t: time.Unix(0, 0)}
	h := newTestHouse(c)

	h.Poll("bff-1") // first poll from a client: news
	h.Poll("bff-1") // nothing changed: empty
	if _, err := h.PlaceBid("a", "ada", 10); err != nil {
		t.Fatal(err)
	}
	h.Poll("bff-1") // saw the bid: news
	a := h.Poll("bff-2")

	if got := a.GetLoad().GetPollsPerSecond(); got != 4/loadWindow.Seconds() {
		t.Fatalf("polls per second: got %v, want %v", got, 4/loadWindow.Seconds())
	}
	if got := a.GetLoad().GetEmptyPollRatio(); got != 0.25 {
		t.Fatalf("empty poll ratio: got %v, want 0.25", got)
	}

	c.advance(loadWindow + time.Second)
	if got := h.Snapshot().GetLoad().GetPollsPerSecond(); got != 0 {
		t.Fatalf("polls per second after the window: got %v, want 0", got)
	}
}

func TestBidAgeIsPerSnapshot(t *testing.T) {
	c := &clock{t: time.Unix(0, 0)}
	h := newTestHouse(c)
	if _, err := h.PlaceBid("a", "ada", 10); err != nil {
		t.Fatal(err)
	}

	c.advance(1500 * time.Millisecond)
	first := h.Snapshot()
	c.advance(time.Second)
	second := h.Snapshot()

	if got := first.GetRecentBids()[0].GetAgeMs(); got != 1500 {
		t.Fatalf("first snapshot: got %dms, want 1500ms", got)
	}
	if got := second.GetHighestBid().GetAgeMs(); got != 2500 {
		t.Fatalf("second snapshot: got %dms, want 2500ms", got)
	}
}

func TestLoadChangedIsNotNews(t *testing.T) {
	c := &clock{t: time.Unix(0, 0)}
	h := newTestHouse(c)
	events, stop := h.Watch()
	defer stop()
	<-events // snapshot

	h.Poll("bff-1")
	c.advance(loadEvery)
	h.Tick()
	if e := <-events; e.GetKind() != auctionv1.WatchAuctionResponse_KIND_LOAD_CHANGED || e.GetAuction().GetLoad().GetPollsPerSecond() == 0 {
		t.Fatalf("after polling started: got %v", e)
	}

	a := h.Poll("bff-1") // only the load figure changed since the last poll
	if got := a.GetLoad().GetEmptyPollRatio(); got != 0.5 {
		t.Fatalf("empty poll ratio: got %v, want 0.5", got)
	}

	c.advance(loadEvery / 2)
	h.Tick()
	select {
	case e := <-events:
		t.Fatalf("load announced again within loadEvery: got %v", e.GetKind())
	default:
	}
}
