package testsuite

import (
	"github.com/libs4go/scf4go"
	"github.com/libs4go/smf4go"
)

type raftTestSuite struct {
}

// New create raft standard testsuite
func New(config scf4go.Config) (smf4go.Service, error) {
	return &raftTestSuite{}, nil
}
