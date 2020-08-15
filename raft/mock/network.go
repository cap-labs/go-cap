package mock

import (
	"github.com/cap-labs/go-cap/raft"
	"github.com/libs4go/scf4go"
	"github.com/libs4go/smf4go"
)

// Network .
type Network interface {
	raft.Network
	CheckServer() bool
}

type mockNetwork struct {
	local  raft.Peer
	server raft.RaftStateMachineServer
}

// NewNetwork .
func NewNetwork(config scf4go.Config) (smf4go.Service, error) {
	return &mockNetwork{}, nil
}

func (mock *mockNetwork) Connect(remote *raft.Peer) (raft.RaftStateMachineClient, error) {
	return newStateMachine(remote), nil
}

func (mock *mockNetwork) Serve(local *raft.Peer, server raft.RaftStateMachineServer) error {
	mock.server = server
	mock.local = *local
	return nil
}
