# Bookkeepify

Turns a Spotify Liked Songs pile into a small set of coherent playlists, then keeps it sorted by routing every new like to the right playlist automatically.

[Design Doc](docs/design_v1.md)

## Layout

```
cmd/bookkeepify/   one binary: api | worker | scheduler subcommands
internal/config/   environment config
internal/server/   HTTP API (reads Postgres only; never calls upstreams on the request path)
```

## Running locally

Requires Go 1.27+ and Docker.

```sh
cp .env.example .env          # fill in Spotify and Last.fm credentials
docker compose up -d postgres
set -a; source .env; set +a
go run ./cmd/bookkeepify api  # http://localhost:8080/healthz
```

```sh
go test ./...
```

## Milestones

| | Phase | Exit criteria |
| --- | --- | --- |
| M0 | API spike | Go/no-go per data source (Spotify endpoints, artist `genres`, Last.fm tag quality) |
| M1 | Ingest | Full library in Postgres; re-running sync changes nothing |
| M2 | Enrich | ≥ 90% of liked songs have features |
| M3 | Cluster | Most clusters acceptable as playlists |
| M4 | Apply + UI | Real playlists created, no duplicates on retry |
| M5 | Auto-sort | A week of new likes sorted with ≤ 15% corrections |
| M6 | Scale lab | Simulate 10,000 users and find the first bottleneck |
