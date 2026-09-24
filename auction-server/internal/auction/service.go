package auction

import (
	"context"
	"errors"
	"log/slog"

	auctionv1 "github.com/dignitas/dignitas-grpc-demo/auction-server/gen/auction/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service exposes the House over gRPC. It implements the generated
// auctionv1.AuctionServiceServer interface.
type Service struct {
	auctionv1.UnimplementedAuctionServiceServer
	house *House
}

func NewService(house *House) *Service {
	return &Service{house: house}
}

func (s *Service) GetAuction(context.Context, *auctionv1.GetAuctionRequest) (*auctionv1.GetAuctionResponse, error) {
	return &auctionv1.GetAuctionResponse{Auction: s.house.Snapshot()}, nil
}

func (s *Service) PlaceBid(_ context.Context, req *auctionv1.PlaceBidRequest) (*auctionv1.PlaceBidResponse, error) {
	auction, err := s.house.PlaceBid(req.GetLotId(), req.GetBidder(), req.GetAmount())
	switch {
	case errors.Is(err, ErrInvalidBid):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrBidRejected):
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	case err != nil:
		return nil, status.Error(codes.Internal, err.Error())
	}
	slog.Info("bid accepted", "bidder", req.GetBidder(), "amount", req.GetAmount(), "lot", req.GetLotId())
	return &auctionv1.PlaceBidResponse{Auction: auction}, nil
}

func (s *Service) WatchAuction(_ *auctionv1.WatchAuctionRequest, stream grpc.ServerStreamingServer[auctionv1.WatchAuctionResponse]) error {
	events, stop := s.house.Watch()
	defer stop()

	for {
		select {
		case <-stream.Context().Done():
			// The client went away or cancelled: that is how a stream normally ends.
			return nil
		case event := <-events:
			if err := stream.Send(event); err != nil {
				return err
			}
		}
	}
}
