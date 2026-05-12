package metrics

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
)

type Metrics struct {
	mu sync.Mutex

	nodeID       string
	role         string
	term         uint64
	commitIndex  uint64
	lastApplied  uint64
	lastLogIndex uint64
	elections    uint64
	rpcs         map[string]uint64
	lag          map[string]uint64
}

func New(nodeID string) *Metrics {
	return &Metrics{
		nodeID: nodeID,
		role:   "follower",
		rpcs:   map[string]uint64{},
		lag:    map[string]uint64{},
	}
}

func (m *Metrics) SetState(role string, term, commitIndex, lastApplied, lastLogIndex uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.role = role
	m.term = term
	m.commitIndex = commitIndex
	m.lastApplied = lastApplied
	m.lastLogIndex = lastLogIndex
}

func (m *Metrics) IncElection() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.elections++
}

func (m *Metrics) IncRPC(method, result string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rpcs[method+"|"+result]++
}

func (m *Metrics) SetReplicationLag(peerID string, lag uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lag[peerID] = lag
}

func (m *Metrics) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	labels := fmt.Sprintf(`node_id="%s"`, escapeLabel(m.nodeID))
	leader := 0
	if m.role == "leader" {
		leader = 1
	}

	fmt.Fprintf(w, "# HELP raft_node_up Node liveness marker.\n")
	fmt.Fprintf(w, "# TYPE raft_node_up gauge\n")
	fmt.Fprintf(w, "raft_node_up{%s} 1\n", labels)
	fmt.Fprintf(w, "# HELP raft_is_leader Whether this node currently believes it is leader.\n")
	fmt.Fprintf(w, "# TYPE raft_is_leader gauge\n")
	fmt.Fprintf(w, "raft_is_leader{%s} %d\n", labels, leader)
	fmt.Fprintf(w, "# HELP raft_current_term Current Raft term.\n")
	fmt.Fprintf(w, "# TYPE raft_current_term gauge\n")
	fmt.Fprintf(w, "raft_current_term{%s} %d\n", labels, m.term)
	fmt.Fprintf(w, "# HELP raft_commit_index Highest committed log index.\n")
	fmt.Fprintf(w, "# TYPE raft_commit_index gauge\n")
	fmt.Fprintf(w, "raft_commit_index{%s} %d\n", labels, m.commitIndex)
	fmt.Fprintf(w, "# HELP raft_last_applied Highest state-machine-applied log index.\n")
	fmt.Fprintf(w, "# TYPE raft_last_applied gauge\n")
	fmt.Fprintf(w, "raft_last_applied{%s} %d\n", labels, m.lastApplied)
	fmt.Fprintf(w, "# HELP raft_last_log_index Highest local log index.\n")
	fmt.Fprintf(w, "# TYPE raft_last_log_index gauge\n")
	fmt.Fprintf(w, "raft_last_log_index{%s} %d\n", labels, m.lastLogIndex)
	fmt.Fprintf(w, "# HELP raft_elections_total Elections started by this node.\n")
	fmt.Fprintf(w, "# TYPE raft_elections_total counter\n")
	fmt.Fprintf(w, "raft_elections_total{%s} %d\n", labels, m.elections)

	rpcKeys := make([]string, 0, len(m.rpcs))
	for key := range m.rpcs {
		rpcKeys = append(rpcKeys, key)
	}
	sort.Strings(rpcKeys)
	fmt.Fprintf(w, "# HELP raft_rpc_total Consensus RPCs handled by method and result.\n")
	fmt.Fprintf(w, "# TYPE raft_rpc_total counter\n")
	for _, key := range rpcKeys {
		parts := strings.SplitN(key, "|", 2)
		fmt.Fprintf(w, "raft_rpc_total{%s,method=\"%s\",result=\"%s\"} %d\n", labels, escapeLabel(parts[0]), escapeLabel(parts[1]), m.rpcs[key])
	}

	peers := make([]string, 0, len(m.lag))
	for peer := range m.lag {
		peers = append(peers, peer)
	}
	sort.Strings(peers)
	fmt.Fprintf(w, "# HELP raft_replication_lag_entries Leader-side lag by peer.\n")
	fmt.Fprintf(w, "# TYPE raft_replication_lag_entries gauge\n")
	for _, peer := range peers {
		fmt.Fprintf(w, "raft_replication_lag_entries{%s,peer_id=\"%s\"} %d\n", labels, escapeLabel(peer), m.lag[peer])
	}
}

func escapeLabel(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}
