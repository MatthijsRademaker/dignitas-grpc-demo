using Bff;
using Dignitas.Auction.V1;

var builder = WebApplication.CreateBuilder(args);

// The presenter's auction server, e.g. http://192.168.1.20:50051 (plain HTTP/2, no TLS, room network only).
var auctionHost = builder.Configuration["AUCTION_HOST"]
    ?? throw new InvalidOperationException("AUCTION_HOST is not set. Point it at the auction server, e.g. http://192.168.1.20:50051");

// One typed client, generated from auction.proto. The factory reuses a single HTTP/2 connection
// for every call and stream: that is the multiplexing from the slides.
builder.Services
    .AddGrpcClient<AuctionService.AuctionServiceClient>(options => options.Address = new Uri(auctionHost))
    .ConfigureChannel(channel => channel.ThrowOperationCanceledOnCancellation = true)
    .ConfigurePrimaryHttpMessageHandler(() => new SocketsHttpHandler
    {
        // Keep long-lived streams alive through Wi-Fi hiccups and idle NAT timeouts.
        KeepAlivePingDelay = TimeSpan.FromSeconds(20),
        KeepAlivePingTimeout = TimeSpan.FromSeconds(10),
        KeepAlivePingPolicy = HttpKeepAlivePingPolicy.Always,
        EnableMultipleHttp2Connections = true,
    });

builder.Services.AddProblemDetails();

var app = builder.Build();

app.MapGet("/api/info", () => new { AuctionHost = auctionHost });
app.MapAuctionEndpoints();

app.Run();
