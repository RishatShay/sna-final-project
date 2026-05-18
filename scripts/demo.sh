#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

NODES=(node1 node2 node3 node4 node5)
PROMETHEUS_URL="${PROMETHEUS_URL:-http://localhost:9090}"
GRAFANA_URL="${GRAFANA_URL:-http://localhost:3000}"
LOKI_URL="${LOKI_URL:-http://localhost:3100}"

COMPOSE_CMD=()
if docker compose version >/dev/null 2>&1; then
  COMPOSE_CMD=(docker compose)
elif command -v docker-compose >/dev/null 2>&1; then
  COMPOSE_CMD=(docker-compose)
else
  echo "ERROR: Docker Compose is not installed. Install docker compose or docker-compose." >&2
  exit 1
fi

compose() {
  "${COMPOSE_CMD[@]}" "$@"
}

usage() {
  cat <<'USAGE'
Raft SNA demo helper

Usage:
  ./scripts/demo.sh preflight      Validate local demo prerequisites and config shape.
  ./scripts/demo.sh start          Build and start the full cluster/monitoring stack
  ./scripts/demo.sh status         Print raftctl status from every reachable node
  ./scripts/demo.sh workload       Write/read demo keys and show status
  ./scripts/demo.sh metrics        Query the main Prometheus metrics
  ./scripts/demo.sh logs           Query recent Raft logs from Loki.
  ./scripts/demo.sh failover       Stop the current leader, wait for a new leader, then write again
  ./scripts/demo.sh degrade        Stop one non-leader while preserving quorum, then write again.
  ./scripts/demo.sh quorum-loss    Stop enough nodes to lose quorum and show alert-facing metrics
  ./scripts/demo.sh restore        Start all Raft nodes again and show recovered status
  ./scripts/demo.sh stop-node NODE Stop one specific node
  ./scripts/demo.sh down           Stop the whole compose stack.

Typical presentation flow:
  preflight -> start -> status -> workload -> metrics -> logs
  -> failover -> degrade -> quorum-loss -> restore -> down
USAGE
}

section() {
  printf '\n== %s ==\n' "$1"
}

info() {
  printf '%s\n' "$*"
}

warn() {
  printf 'WARNING: %s\n' "$*" >&2
}

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "ERROR: required command is missing: $1" >&2
    exit 1
  fi
}

wait_http() {
  local name="$1"
  local url="$2"
  local attempts="${3:-60}"

  info "Waiting for ${name}: ${url}"
  for _ in $(seq 1 "$attempts"); do
    if curl -fsS "$url" >/dev/null 2>&1; then
      info "${name} is ready"
      return 0
    fi
    sleep 2
  done

  echo "ERROR: ${name} did not become ready at ${url}" >&2
  return 1
}

raftctl() {
  local node="$1"
  shift
  compose exec -T "$node" raftctl -addr "${node}:9001" "$@"
}

status_json_from() {
  local node="$1"
  raftctl "$node" status 2>/dev/null
}

json_field() {
  local field="$1"
  sed -n "s/.*\"${field}\": \"\\([^\"]*\\)\".*/\\1/p" | head -n 1
}

find_leader() {
  local node status role id leader
  for node in "${NODES[@]}"; do
    if status="$(status_json_from "$node")"; then
      role="$(printf '%s\n' "$status" | json_field role)"
      id="$(printf '%s\n' "$status" | json_field id)"
      if [[ "$role" == "leader" && -n "$id" ]]; then
        printf '%s\n' "$id"
        return 0
      fi
      leader="$(printf '%s\n' "$status" | json_field leader_id)"
      if [[ -n "$leader" ]]; then
        printf '%s\n' "$leader"
        return 0
      fi
    fi
  done
  return 1
}

running_nodes() {
  local service
  if compose ps --services --filter status=running >/dev/null 2>&1; then
    compose ps --services --filter status=running | grep -E '^node[1-5]$' || true
    return 0
  fi

  for service in "${NODES[@]}"; do
    if compose ps -q "$service" >/dev/null 2>&1; then
      printf '%s\n' "$service"
    fi
  done
}

first_reachable_node() {
  local node
  for node in "${NODES[@]}"; do
    if status_json_from "$node" >/dev/null 2>&1; then
      printf '%s\n' "$node"
      return 0
    fi
  done
  return 1
}

prom_query() {
  local query="$1"
  local raw
  info "Prometheus query: ${query}"
  raw="$(curl -fsS -G "${PROMETHEUS_URL}/api/v1/query" --data-urlencode "query=${query}")"
  if command -v jq >/dev/null 2>&1; then
    printf '%s\n' "$raw" | jq -r '
      if (.data.result | length) == 0 then
        "  no series returned"
      else
        .data.result[] |
        (.metric | del(.__name__) | to_entries | sort_by(.key) | map("\(.key)=\(.value)") | join(", ")) as $labels |
        "  " + (if $labels == "" then "value" else $labels end) + " => " + .value[1]
      end
    '
  else
    printf '%s\n' "$raw"
  fi
}

preflight() {
  section "Compose command"
  info "Using: ${COMPOSE_CMD[*]}"

  section "Docker daemon"
  if ! docker ps >/dev/null 2>&1; then
    warn "Docker daemon is not reachable by this user. Start Docker or add permissions before the live demo."
  else
    info "Docker daemon is reachable"
  fi

  section "Config validation"
  compose config >/dev/null
  info "docker compose config is valid"

  section "Secrets check"
  if grep -q 'GF_SMTP_PASSWORD: example' docker-compose.yml; then
    warn "Grafana SMTP password is still the placeholder 'example'. Set it locally before email-alert demo."
  fi
  if grep -q 'addresses: "example.com"' deployments/grafana/provisioning/alerting/contactpoints.yaml; then
    warn "Grafana alert email address is still example.com. Set a real recipient locally before email-alert demo."
  fi
  info "Do not commit real SMTP passwords or mailbox addresses."

  section "Important URLs"
  info "Grafana:    ${GRAFANA_URL}  admin/admin"
  info "Prometheus: ${PROMETHEUS_URL}"
  info "Loki:       ${LOKI_URL}"
  info "Targets:    ${PROMETHEUS_URL}/targets"
}

start_stack() {
  require_command curl
  section "Starting stack"
  compose up --build -d

  section "Waiting for endpoints"
  wait_http "Prometheus" "${PROMETHEUS_URL}/-/ready" 60
  wait_http "Grafana" "${GRAFANA_URL}/api/health" 60
  wait_http "node1 healthz" "http://localhost:8001/healthz" 60

  section "Initial status"
  sleep 5
  status_all
}

status_all() {
  local node
  section "Raft status"
  for node in "${NODES[@]}"; do
    printf '\n--- %s ---\n' "$node"
    if ! raftctl "$node" status; then
      info "${node} is not reachable"
    fi
  done
}

workload() {
  local target now key
  target="$(first_reachable_node)"
  now="$(date +%H%M%S)"

  section "Client writes"
  for key in key1 key2 key3; do
    raftctl "$target" write "$key" "value-${key}-${now}"
  done

  section "Client reads"
  raftctl "$target" read key1
  raftctl "$target" read key2

  section "Status after workload"
  raftctl "$target" status
}

metrics() {
  require_command curl
  section "Prometheus metrics"
  prom_query 'up{job="raft"}'
  prom_query 'raft_is_leader'
  prom_query 'raft_commit_index'
  prom_query 'sum by (method,result) (increase(raft_rpc_total[5m]))'
  prom_query 'raft_replication_lag_entries'
}

logs() {
  require_command curl
  local raw
  section "Loki logs"
  info 'Loki query: {node_id=~"node[1-5]"} |= "client"'
  raw="$(curl -fsS -G "${LOKI_URL}/loki/api/v1/query_range" \
    --data-urlencode 'query={node_id=~"node[1-5]"} |= "client"' \
    --data-urlencode 'limit=30')"
  if command -v jq >/dev/null 2>&1; then
    printf '%s\n' "$raw" | jq -r '
      if (.data.result | length) == 0 then
        "  no logs returned; run ./scripts/demo.sh workload once and retry"
      else
        .data.result[] |
        .stream as $stream |
        .values[] |
        "  " + ((.[0] | tonumber / 1000000000) | strftime("%H:%M:%S")) + " " + ($stream.node_id // $stream.container // "-") + " " + .[1]
      end
    '
  else
    printf '%s\n' "$raw"
  fi
}

stop_node() {
  local node="${1:-}"
  if [[ -z "$node" ]]; then
    echo "ERROR: stop-node requires a node name, for example node3" >&2
    exit 2
  fi
  section "Stopping ${node}"
  compose stop "$node"
}

failover() {
  local old_leader new_leader target
  old_leader="$(find_leader)"
  section "Stopping current leader"
  info "Current leader: ${old_leader}"
  compose stop "$old_leader"

  section "Waiting for a new leader"
  sleep 8
  new_leader="$(find_leader)"
  info "New leader: ${new_leader}"

  section "Write after leader failover"
  target="$(first_reachable_node)"
  raftctl "$target" write after_failover ok
  raftctl "$target" read after_failover
  status_all
}

degrade() {
  local leader node target=""
  leader="$(find_leader || true)"
  section "Stopping one non-leader while keeping quorum"
  while IFS= read -r node; do
    if [[ "$node" != "$leader" ]]; then
      target="$node"
      break
    fi
  done < <(running_nodes)

  if [[ -z "$target" ]]; then
    echo "ERROR: no non-leader node is available to stop" >&2
    exit 1
  fi

  info "Leader remains: ${leader:-unknown}"
  info "Stopping: ${target}"
  compose stop "$target"
  sleep 5

  section "Write with reduced cluster"
  node="$(first_reachable_node)"
  raftctl "$node" write degraded_cluster still_has_quorum
  raftctl "$node" read degraded_cluster
  status_all
}

quorum_loss() {
  local running_count down_count need leader node targets=()
  running_count="$(running_nodes | wc -l | tr -d ' ')"
  down_count="$((5 - running_count))"
  need="$((3 - down_count))"

  section "Quorum loss"
  info "Running nodes: ${running_count}; already down: ${down_count}; need to stop: ${need}"
  if (( need <= 0 )); then
    info "Cluster has already lost quorum."
  else
    leader="$(find_leader || true)"
    if [[ -n "$leader" ]]; then
      targets+=("$leader")
    fi
    while IFS= read -r node; do
      if (( ${#targets[@]} >= need )); then
        break
      fi
      if [[ "$node" != "$leader" ]]; then
        targets+=("$node")
      fi
    done < <(running_nodes)

    info "Stopping nodes: ${targets[*]}"
    compose stop "${targets[@]}"
  fi

  section "Alert-facing metrics"
  sleep 10
  prom_query 'up{job="raft"}'
  prom_query 'count(up{job="raft"} == 0)'
  info "Grafana email alerts are evaluated every 30s; the 3-node quorum-loss rule has for: 30s."
}

restore() {
  section "Restoring all Raft nodes"
  compose up -d node1 node2 node3 node4 node5
  sleep 10
  status_all
  metrics
}

down() {
  section "Stopping compose stack"
  compose down
}

main() {
  local command="${1:-}"
  case "$command" in
    preflight) preflight ;;
    tests) tests ;;
    start) start_stack ;;
    status) status_all ;;
    workload) workload ;;
    metrics) metrics ;;
    logs) logs ;;
    failover) failover ;;
    degrade) degrade ;;
    quorum-loss) quorum_loss ;;
    restore) restore ;;
    stop-node) shift; stop_node "${1:-}" ;;
    down) down ;;
    ""|-h|--help|help) usage ;;
    *)
      echo "ERROR: unknown command: ${command}" >&2
      usage >&2
      exit 2
      ;;
  esac
}

main "$@"
