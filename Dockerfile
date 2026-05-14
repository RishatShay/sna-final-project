FROM golang:1.24-bookworm AS build


RUN apt-get update && apt-get install -y gcc && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /out/raftnode ./cmd/raftnode
RUN CGO_ENABLED=1 go build -o /out/raftctl ./cmd/raftctl

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/raftnode /usr/local/bin/raftnode
COPY --from=build /out/raftctl /usr/local/bin/raftctl

EXPOSE 9001 8001
ENTRYPOINT ["/usr/local/bin/raftnode"]
