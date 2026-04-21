FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o /out/gateway ./cmd/main.go

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update \
	&& apt-get install -y --no-install-recommends bash ca-certificates tzdata \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/gateway /app/gateway
COPY secret /app/secret

EXPOSE 8080

ENTRYPOINT ["/bin/bash", "-lc", "set -a; . /app/secret/.env; set +a; exec /app/gateway"]
