package raft

import (
	"github.com/libs4go/errors"
)

//go:generate protoc --proto_path=../proto-cap --go_out=plugins=grpc,paths=source_relative:. raft.proto

// ScopeOfAPIError .
const errVendor = "raft"

// errors
var (
	ErrOutOfCluster = errors.New("the peer is not in cluster config", errors.WithVendor(errVendor), errors.WithCode(-1))
)

// Network .
type Network interface {
	// Connect to remote raft peer
	Connect(remote *Peer) (RaftStateMachineClient, error)

	Serve(local *Peer, server RaftStateMachineServer) error
}

// ClusterManager .
type ClusterManager interface {
	Get() (*Cluster, error)
	Set(*Cluster) error
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
	// Last entry
	LastEntry() (*Entry, error)
	// Create snapshot
	CreateSnapshot() error
	// Get snapshot stream
	Snapshot() (chan *RequestInstallSnapshot, error)
	// Install snapshot
	InstallSnapshot(chan *RequestInstallSnapshot) error
}
