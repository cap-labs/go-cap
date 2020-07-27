package raft

import "github.com/cap-labs/go-cap"

//go:generate protoc --proto_path=../proto-cap --go_out=plugins=grpc,paths=source_relative:. raft.proto

// Network .
type Network interface {
	// Connect to remote raft peer
	Connect(remote *cap.Peer) (RaftStateMachineClient, error)

	Serve(local *cap.Peer, server RaftStateMachineServer) error
}

// ClusterManager .
type ClusterManager interface {
	Get(*cap.Cluster) error
	Set(*cap.Cluster) error
}

// LogStore .
type LogStore interface {
	// Get log entries by term id
	GetByTerm(term int64) ([]*Entry, error)
	// Get log entry by offset index
	Get(index int64) (*Entry, error)
	// Store entries
	Save(entries []*Entry) error
	// Commit logs before index (include index self)
	Commit(index int64) error
	// Last commited entry
	LastCommitted() (*Entry, error)
	// Create snapshot
	CreateSnapshot() error
	// Get snapshot stream
	Snapshot() (chan *RequestInstallSnapshot, error)
	// Install snapshot
	InstallSnapshot(chan *RequestInstallSnapshot) error
}
