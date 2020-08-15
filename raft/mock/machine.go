package mock

import (
	"context"

	"github.com/cap-labs/go-cap/raft"
	"google.golang.org/grpc"
)

// StateMachine .
type StateMachine interface {
}

type mockStateMachine struct {
	local raft.Peer
}

func newStateMachine(peer *raft.Peer) *mockStateMachine {
	return &mockStateMachine{
		local: *peer,
	}
}

func (m *mockStateMachine) Vote(ctx context.Context, in *raft.RequestVote, opts ...grpc.CallOption) (*raft.ResponseVote, error) {
	return nil, nil
}

func (m *mockStateMachine) AppendEntries(ctx context.Context, in *raft.RequestAppendEntries, opts ...grpc.CallOption) (*raft.ResponseAppendEntries, error) {
	return nil, nil
}

func (m *mockStateMachine) InstallSnapshot(ctx context.Context, in *raft.RequestInstallSnapshot, opts ...grpc.CallOption) (*raft.ResponseInstallSnapshot, error) {
	return nil, nil
}
