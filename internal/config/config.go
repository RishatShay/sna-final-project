package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Peer struct {
	ID      string
	Address string
}

type Config struct {
	NodeID            string
	GRPCAddr          string
	HTTPAddr          string
	DataDir           string
	Peers             []Peer
	ElectionMin       time.Duration
	ElectionMax       time.Duration
	HeartbeatInterval time.Duration
	SnapshotThreshold uint64
}

func FromEnv() (Config, error) {
	cfg := Config{
		NodeID:            env("NODE_ID", "node1"),
		GRPCAddr:          env("RAFT_GRPC_ADDR", ":9001"),
		HTTPAddr:          env("RAFT_HTTP_ADDR", ":8001"),
		DataDir:           env("RAFT_DATA_DIR", "data/node1"),
		ElectionMin:       envDurationMS("RAFT_ELECTION_MIN_MS", 500),
		ElectionMax:       envDurationMS("RAFT_ELECTION_MAX_MS", 900),
		HeartbeatInterval: envDurationMS("RAFT_HEARTBEAT_MS", 100),
		SnapshotThreshold: envUint64("RAFT_SNAPSHOT_THRESHOLD", 10000),
	}

	if cfg.ElectionMin <= cfg.HeartbeatInterval {
		return Config{}, fmt.Errorf("RAFT_ELECTION_MIN_MS must be greater than RAFT_HEARTBEAT_MS")
	}
	if cfg.ElectionMax < cfg.ElectionMin {
		return Config{}, fmt.Errorf("RAFT_ELECTION_MAX_MS must be >= RAFT_ELECTION_MIN_MS")
	}

	peers, err := parsePeers(os.Getenv("RAFT_PEERS"), cfg.NodeID)
	if err != nil {
		return Config{}, err
	}
	cfg.Peers = peers
	return cfg, nil
}

func parsePeers(raw, selfID string) ([]Peer, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	peers := make([]Peer, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		pair := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(pair) != 2 || pair[0] == "" || pair[1] == "" {
			return nil, fmt.Errorf("invalid RAFT_PEERS entry %q, expected id=host:port", part)
		}
		id := pair[0]
		if _, ok := seen[id]; ok {
			return nil, fmt.Errorf("duplicate peer id %q", id)
		}
		seen[id] = struct{}{}
		if id == selfID {
			continue
		}
		peers = append(peers, Peer{ID: id, Address: pair[1]})
	}
	return peers, nil
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envDurationMS(key string, fallback int) time.Duration {
	value := env(key, strconv.Itoa(fallback))
	ms, err := strconv.Atoi(value)
	if err != nil || ms <= 0 {
		return time.Duration(fallback) * time.Millisecond
	}
	return time.Duration(ms) * time.Millisecond
}

func envUint64(key string, fallback uint64) uint64 {
	value := env(key, strconv.FormatUint(fallback, 10))
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
