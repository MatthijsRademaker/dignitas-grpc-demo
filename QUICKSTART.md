# gRPC Auction House: quickstart

Five minutes now saves fifteen at the venue, where the Wi-Fi is shared by everyone.

## You need

- **Git**
- **Docker**: Docker Desktop on Windows or macOS, or Docker Engine with the Compose plugin on Linux. Compose
  2.22 or newer (`docker compose version`), for `--watch`.
- Optional: the **.NET 10 SDK** and your C# IDE of choice, for autocomplete while you code. Nothing else: no Go,
  no Node.

## Before the day

```bash
git clone https://github.com/MatthijsRademaker/dignitas-grpc-demo.git
cd dignitas-grpc-demo
cp .env.example .env
```

Open `.env` and set it up for a dry run on your own machine:

```bash
AUCTION_HOST=
AUCTION_BOTS=3
```

An empty `AUCTION_HOST` means "use a local auction server". The three bots bid so you're not alone. Then:

```bash
docker compose --profile presenter up --build
```

The first build takes a few minutes. Open <http://localhost:8080>: you should see a mystery lot on the block and
bots bidding on it. Try to bid yourself: "Workshop time" means everything works. That's the part we build together.

Stop with `Ctrl+C`. Have the .NET SDK? Run `cd bff && dotnet build` once, so your IDE knows the generated gRPC client.

## On the day

1. Same Wi-Fi as the presenter, and no VPN: your laptop connects straight to the presenter's machine.
2. Put the presenter's IP in `.env` (it's on the slide):

   ```bash
   AUCTION_HOST=http://<presenter-ip>:50051
   ```

3. `docker compose up --build --watch`, then follow [`workshop/README.md`](workshop/README.md).

Stuck? The troubleshooting table is at the bottom of [`workshop/README.md`](workshop/README.md).
