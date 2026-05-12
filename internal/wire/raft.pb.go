package wire

import proto "github.com/golang/protobuf/proto"

const _ = proto.ProtoPackageIsVersion3

type LogEntry struct {
	Index   uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Term    uint64 `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	Command []byte `protobuf:"bytes,3,opt,name=command,proto3" json:"command,omitempty"`
}

func (m *LogEntry) Reset()         { *m = LogEntry{} }
func (m *LogEntry) String() string { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()    {}
func (m *LogEntry) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}
func (m *LogEntry) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}
func (m *LogEntry) GetCommand() []byte {
	if m != nil {
		return m.Command
	}
	return nil
}

type RequestVoteRequest struct {
	Term         uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	CandidateId  string `protobuf:"bytes,2,opt,name=candidate_id,json=candidateId,proto3" json:"candidate_id,omitempty"`
	LastLogIndex uint64 `protobuf:"varint,3,opt,name=last_log_index,json=lastLogIndex,proto3" json:"last_log_index,omitempty"`
	LastLogTerm  uint64 `protobuf:"varint,4,opt,name=last_log_term,json=lastLogTerm,proto3" json:"last_log_term,omitempty"`
}

func (m *RequestVoteRequest) Reset()         { *m = RequestVoteRequest{} }
func (m *RequestVoteRequest) String() string { return proto.CompactTextString(m) }
func (*RequestVoteRequest) ProtoMessage()    {}
func (m *RequestVoteRequest) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}
func (m *RequestVoteRequest) GetCandidateId() string {
	if m != nil {
		return m.CandidateId
	}
	return ""
}
func (m *RequestVoteRequest) GetLastLogIndex() uint64 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}
func (m *RequestVoteRequest) GetLastLogTerm() uint64 {
	if m != nil {
		return m.LastLogTerm
	}
	return 0
}

type RequestVoteResponse struct {
	Term        uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	VoteGranted bool   `protobuf:"varint,2,opt,name=vote_granted,json=voteGranted,proto3" json:"vote_granted,omitempty"`
}

func (m *RequestVoteResponse) Reset()         { *m = RequestVoteResponse{} }
func (m *RequestVoteResponse) String() string { return proto.CompactTextString(m) }
func (*RequestVoteResponse) ProtoMessage()    {}
func (m *RequestVoteResponse) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}
func (m *RequestVoteResponse) GetVoteGranted() bool {
	if m != nil {
		return m.VoteGranted
	}
	return false
}

type AppendEntriesRequest struct {
	Term         uint64      `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId     string      `protobuf:"bytes,2,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	PrevLogIndex uint64      `protobuf:"varint,3,opt,name=prev_log_index,json=prevLogIndex,proto3" json:"prev_log_index,omitempty"`
	PrevLogTerm  uint64      `protobuf:"varint,4,opt,name=prev_log_term,json=prevLogTerm,proto3" json:"prev_log_term,omitempty"`
	Entries      []*LogEntry `protobuf:"bytes,5,rep,name=entries,proto3" json:"entries,omitempty"`
	LeaderCommit uint64      `protobuf:"varint,6,opt,name=leader_commit,json=leaderCommit,proto3" json:"leader_commit,omitempty"`
}

func (m *AppendEntriesRequest) Reset()         { *m = AppendEntriesRequest{} }
func (m *AppendEntriesRequest) String() string { return proto.CompactTextString(m) }
func (*AppendEntriesRequest) ProtoMessage()    {}
func (m *AppendEntriesRequest) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}
func (m *AppendEntriesRequest) GetLeaderId() string {
	if m != nil {
		return m.LeaderId
	}
	return ""
}
func (m *AppendEntriesRequest) GetPrevLogIndex() uint64 {
	if m != nil {
		return m.PrevLogIndex
	}
	return 0
}
func (m *AppendEntriesRequest) GetPrevLogTerm() uint64 {
	if m != nil {
		return m.PrevLogTerm
	}
	return 0
}
func (m *AppendEntriesRequest) GetEntries() []*LogEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}
func (m *AppendEntriesRequest) GetLeaderCommit() uint64 {
	if m != nil {
		return m.LeaderCommit
	}
	return 0
}

type AppendEntriesResponse struct {
	Term       uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Success    bool   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	MatchIndex uint64 `protobuf:"varint,3,opt,name=match_index,json=matchIndex,proto3" json:"match_index,omitempty"`
}

func (m *AppendEntriesResponse) Reset()         { *m = AppendEntriesResponse{} }
func (m *AppendEntriesResponse) String() string { return proto.CompactTextString(m) }
func (*AppendEntriesResponse) ProtoMessage()    {}
func (m *AppendEntriesResponse) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}
func (m *AppendEntriesResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}
func (m *AppendEntriesResponse) GetMatchIndex() uint64 {
	if m != nil {
		return m.MatchIndex
	}
	return 0
}

type InstallSnapshotRequest struct {
	Term              uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId          string `protobuf:"bytes,2,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	LastIncludedIndex uint64 `protobuf:"varint,3,opt,name=last_included_index,json=lastIncludedIndex,proto3" json:"last_included_index,omitempty"`
	LastIncludedTerm  uint64 `protobuf:"varint,4,opt,name=last_included_term,json=lastIncludedTerm,proto3" json:"last_included_term,omitempty"`
	Data              []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *InstallSnapshotRequest) Reset()         { *m = InstallSnapshotRequest{} }
func (m *InstallSnapshotRequest) String() string { return proto.CompactTextString(m) }
func (*InstallSnapshotRequest) ProtoMessage()    {}
func (m *InstallSnapshotRequest) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}
func (m *InstallSnapshotRequest) GetLeaderId() string {
	if m != nil {
		return m.LeaderId
	}
	return ""
}
func (m *InstallSnapshotRequest) GetLastIncludedIndex() uint64 {
	if m != nil {
		return m.LastIncludedIndex
	}
	return 0
}
func (m *InstallSnapshotRequest) GetLastIncludedTerm() uint64 {
	if m != nil {
		return m.LastIncludedTerm
	}
	return 0
}
func (m *InstallSnapshotRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type InstallSnapshotResponse struct {
	Term uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
}

func (m *InstallSnapshotResponse) Reset()         { *m = InstallSnapshotResponse{} }
func (m *InstallSnapshotResponse) String() string { return proto.CompactTextString(m) }
func (*InstallSnapshotResponse) ProtoMessage()    {}
func (m *InstallSnapshotResponse) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

type ClientWriteRequest struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *ClientWriteRequest) Reset()         { *m = ClientWriteRequest{} }
func (m *ClientWriteRequest) String() string { return proto.CompactTextString(m) }
func (*ClientWriteRequest) ProtoMessage()    {}
func (m *ClientWriteRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}
func (m *ClientWriteRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type ClientWriteResponse struct {
	Success  bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	LeaderId string `protobuf:"bytes,2,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	Error    string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	Index    uint64 `protobuf:"varint,4,opt,name=index,proto3" json:"index,omitempty"`
	Term     uint64 `protobuf:"varint,5,opt,name=term,proto3" json:"term,omitempty"`
}

func (m *ClientWriteResponse) Reset()         { *m = ClientWriteResponse{} }
func (m *ClientWriteResponse) String() string { return proto.CompactTextString(m) }
func (*ClientWriteResponse) ProtoMessage()    {}
func (m *ClientWriteResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}
func (m *ClientWriteResponse) GetLeaderId() string {
	if m != nil {
		return m.LeaderId
	}
	return ""
}
func (m *ClientWriteResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}
func (m *ClientWriteResponse) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}
func (m *ClientWriteResponse) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}

type ClientReadRequest struct {
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *ClientReadRequest) Reset()         { *m = ClientReadRequest{} }
func (m *ClientReadRequest) String() string { return proto.CompactTextString(m) }
func (*ClientReadRequest) ProtoMessage()    {}
func (m *ClientReadRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type ClientReadResponse struct {
	Success  bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	LeaderId string `protobuf:"bytes,2,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	Error    string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	Value    string `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *ClientReadResponse) Reset()         { *m = ClientReadResponse{} }
func (m *ClientReadResponse) String() string { return proto.CompactTextString(m) }
func (*ClientReadResponse) ProtoMessage()    {}
func (m *ClientReadResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}
func (m *ClientReadResponse) GetLeaderId() string {
	if m != nil {
		return m.LeaderId
	}
	return ""
}
func (m *ClientReadResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}
func (m *ClientReadResponse) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type StatusRequest struct{}

func (m *StatusRequest) Reset()         { *m = StatusRequest{} }
func (m *StatusRequest) String() string { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()    {}

type PeerStatus struct {
	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Address    string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	MatchIndex uint64 `protobuf:"varint,3,opt,name=match_index,json=matchIndex,proto3" json:"match_index,omitempty"`
	NextIndex  uint64 `protobuf:"varint,4,opt,name=next_index,json=nextIndex,proto3" json:"next_index,omitempty"`
}

func (m *PeerStatus) Reset()         { *m = PeerStatus{} }
func (m *PeerStatus) String() string { return proto.CompactTextString(m) }
func (*PeerStatus) ProtoMessage()    {}
func (m *PeerStatus) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}
func (m *PeerStatus) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}
func (m *PeerStatus) GetMatchIndex() uint64 {
	if m != nil {
		return m.MatchIndex
	}
	return 0
}
func (m *PeerStatus) GetNextIndex() uint64 {
	if m != nil {
		return m.NextIndex
	}
	return 0
}

type StatusResponse struct {
	Id           string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Role         string        `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
	Term         uint64        `protobuf:"varint,3,opt,name=term,proto3" json:"term,omitempty"`
	LeaderId     string        `protobuf:"bytes,4,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	CommitIndex  uint64        `protobuf:"varint,5,opt,name=commit_index,json=commitIndex,proto3" json:"commit_index,omitempty"`
	LastApplied  uint64        `protobuf:"varint,6,opt,name=last_applied,json=lastApplied,proto3" json:"last_applied,omitempty"`
	LastLogIndex uint64        `protobuf:"varint,7,opt,name=last_log_index,json=lastLogIndex,proto3" json:"last_log_index,omitempty"`
	Peers        []*PeerStatus `protobuf:"bytes,8,rep,name=peers,proto3" json:"peers,omitempty"`
}

func (m *StatusResponse) Reset()         { *m = StatusResponse{} }
func (m *StatusResponse) String() string { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()    {}
func (m *StatusResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}
func (m *StatusResponse) GetRole() string {
	if m != nil {
		return m.Role
	}
	return ""
}
func (m *StatusResponse) GetTerm() uint64 {
	if m != nil {
		return m.Term
	}
	return 0
}
func (m *StatusResponse) GetLeaderId() string {
	if m != nil {
		return m.LeaderId
	}
	return ""
}
func (m *StatusResponse) GetCommitIndex() uint64 {
	if m != nil {
		return m.CommitIndex
	}
	return 0
}
func (m *StatusResponse) GetLastApplied() uint64 {
	if m != nil {
		return m.LastApplied
	}
	return 0
}
func (m *StatusResponse) GetLastLogIndex() uint64 {
	if m != nil {
		return m.LastLogIndex
	}
	return 0
}
func (m *StatusResponse) GetPeers() []*PeerStatus {
	if m != nil {
		return m.Peers
	}
	return nil
}

func init() {
	proto.RegisterType((*LogEntry)(nil), "raft.v1.LogEntry")
	proto.RegisterType((*RequestVoteRequest)(nil), "raft.v1.RequestVoteRequest")
	proto.RegisterType((*RequestVoteResponse)(nil), "raft.v1.RequestVoteResponse")
	proto.RegisterType((*AppendEntriesRequest)(nil), "raft.v1.AppendEntriesRequest")
	proto.RegisterType((*AppendEntriesResponse)(nil), "raft.v1.AppendEntriesResponse")
	proto.RegisterType((*InstallSnapshotRequest)(nil), "raft.v1.InstallSnapshotRequest")
	proto.RegisterType((*InstallSnapshotResponse)(nil), "raft.v1.InstallSnapshotResponse")
	proto.RegisterType((*ClientWriteRequest)(nil), "raft.v1.ClientWriteRequest")
	proto.RegisterType((*ClientWriteResponse)(nil), "raft.v1.ClientWriteResponse")
	proto.RegisterType((*ClientReadRequest)(nil), "raft.v1.ClientReadRequest")
	proto.RegisterType((*ClientReadResponse)(nil), "raft.v1.ClientReadResponse")
	proto.RegisterType((*StatusRequest)(nil), "raft.v1.StatusRequest")
	proto.RegisterType((*PeerStatus)(nil), "raft.v1.PeerStatus")
	proto.RegisterType((*StatusResponse)(nil), "raft.v1.StatusResponse")
}
