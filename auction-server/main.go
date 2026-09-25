// Command auction-server is the central gRPC auction house everyone in the room connects to.
package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	auctionv1 "github.com/dignitas/dignitas-grpc-demo/auction-server/gen/auction/v1"
	"github.com/dignitas/dignitas-grpc-demo/auction-server/internal/auction"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func main() {
	addr := flag.String("addr", ":50051", "listen address (all interfaces, so the whole room can reach it)")
	lotDuration := flag.Duration("lot-duration", 60*time.Second, "how long each round is open")
	snipeWindow := flag.Duration("snipe-window", 10*time.Second, "a bid this close to closing extends the round by this much")
	intermission := flag.Duration("intermission", 8*time.Second, "pause between rounds")
	bots := flag.Int("bots", 0, "number of simulated bidders, for rehearsing alone")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	house := auction.NewHouse(auction.Config{
		Lots:         auction.Catalog,
		LotDuration:  *lotDuration,
		SnipeWindow:  *snipeWindow,
		Intermission: *intermission,
	})
	go house.Run(ctx)
	auction.RunBots(ctx, house, *bots)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(logUnary),
		grpc.StreamInterceptor(logStream),
		// Streams can sit idle for a while on conference Wi-Fi; let clients ping to keep them alive.
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{MinTime: 10 * time.Second, PermitWithoutStream: true}),
	)
	auctionv1.RegisterAuctionServiceServer(server, auction.NewService(house))
	// Reflection lets tools like grpcurl and Postman discover the API without the .proto file.
	reflection.Register(server)

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		slog.Error("listen failed", "addr", *addr, "err", err)
		os.Exit(1)
	}
	go func() {
		<-ctx.Done()
		server.Stop()
	}()

	slog.Info("auction house open", "addr", lis.Addr().String(), "bots", *bots)
	if err := server.Serve(lis); err != nil {
		slog.Error("serve failed", "err", err)
		os.Exit(1)
	}
}

func logUnary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	slog.Info("unary", "method", info.FullMethod, "from", peerAddr(ctx), "code", status.Code(err).String(), "took", time.Since(start))
	return resp, err
}

func logStream(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	slog.Info("stream opened", "method", info.FullMethod, "from", peerAddr(ss.Context()))
	err := handler(srv, ss)
	slog.Info("stream closed", "method", info.FullMethod, "from", peerAddr(ss.Context()), "code", status.Code(err).String(), "lasted", time.Since(start).Round(time.Second))
	return err
}

func peerAddr(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}
	return "unknown"
}
