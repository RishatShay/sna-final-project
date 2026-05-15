.PHONY: test race build compose-up compose-down

test:
	go test ./...

race:
	go test -race ./...

build:
	go build ./cmd/raftnode
	go build ./cmd/raftctl

compose-up:
	docker compose up --build

compose-down:
	docker compose down
