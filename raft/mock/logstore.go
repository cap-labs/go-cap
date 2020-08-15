package mock

import (
	"github.com/cap-labs/go-cap/raft"
	"github.com/libs4go/scf4go"
	"github.com/libs4go/smf4go"
)

// LogStore .
type LogStore interface {
	raft.LogStore
}

// mockLogStore .
type mockLogStore struct {
}

// NewLogStore .
func NewLogStore(config scf4go.Config) (smf4go.Service, error) {
	return &mockLogStore{}, nil
}

func (store *mockLogStore) GetByTerm(term int64) ([]*raft.Entry, error) {
	return nil, nil
}

func (store *mockLogStore) Get(index int64) (*raft.Entry, error) {
	return nil, nil
}

func (store *mockLogStore) Save(entries []*raft.Entry) error {
	return nil
}

func (store *mockLogStore) Commit(index int64) error {
	return nil
}

func (store *mockLogStore) LastCommitted() (*raft.Entry, error) {
	return nil, nil
}

func (store *mockLogStore) LastEntry() (*raft.Entry, error) {
	return nil, nil
}

func (store *mockLogStore) CreateSnapshot() error {
	return nil
}

func (store *mockLogStore) Snapshot() (chan *raft.RequestInstallSnapshot, error) {
	return nil, nil
}

func (store *mockLogStore) InstallSnapshot(chan *raft.RequestInstallSnapshot) error {
	return nil
}
