package wire

import (
	"testing"

	newproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
)

func TestMessagesRoundTripThroughNewProtobufRuntime(t *testing.T) {
	in := &AppendEntriesRequest{
		Term:         4,
		LeaderId:     "node1",
		PrevLogIndex: 7,
		PrevLogTerm:  3,
		Entries: []*LogEntry{
			{Index: 8, Term: 4, Command: []byte(`{"op":"set","key":"k","value":"v"}`)},
		},
		LeaderCommit: 8,
	}

	raw, err := newproto.Marshal(protoadapt.MessageV2Of(in))
	if err != nil {
		t.Fatal(err)
	}
	var out AppendEntriesRequest
	if err := newproto.Unmarshal(raw, protoadapt.MessageV2Of(&out)); err != nil {
		t.Fatal(err)
	}
	if out.GetTerm() != in.GetTerm() || out.GetLeaderId() != in.GetLeaderId() || len(out.GetEntries()) != 1 {
		t.Fatalf("round trip mismatch: %#v", out)
	}
	if string(out.GetEntries()[0].GetCommand()) != string(in.GetEntries()[0].GetCommand()) {
		t.Fatalf("command = %q, want %q", out.GetEntries()[0].GetCommand(), in.GetEntries()[0].GetCommand())
	}
}
