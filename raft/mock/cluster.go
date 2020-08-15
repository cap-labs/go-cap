package mock

import (
	"github.com/cap-labs/go-cap/raft"
	"github.com/libs4go/scf4go"
	"github.com/libs4go/smf4go"
)

// ClusterManager .
type ClusterManager interface {
	raft.ClusterManager
}

type mockClusterManager struct {
	cluster raft.Cluster
}

// NewClusterManager .
func NewClusterManager(config scf4go.Config) (smf4go.Service, error) {

	var cluster raft.Cluster

	err := config.Scan(&cluster)

	if err != nil {
		return nil, err
	}

	return &mockClusterManager{
		cluster: cluster,
	}, nil
}

func (cm *mockClusterManager) Get() (*raft.Cluster, error) {
	return &cm.cluster, nil
}

func (cm *mockClusterManager) Set(cluster *raft.Cluster) error {
	cm.cluster = *cluster

	return nil
}
