package supervisor

import (
	"context"

	cap "github.com/cap-labs/go-cap"
	"github.com/libs4go/errors"
	"github.com/libs4go/scf4go"
	"github.com/libs4go/slf4go"
	"github.com/libs4go/smf4go"
)

type supervisorImpl struct {
	slf4go.Logger
	Networking             cap.NetworkingClient `inject:"supervisor.Networking"`
	Consensus              cap.ConsensusClient  `inject:"supervisor.Consensus"`
	Storage                cap.StorageClient    `inject:"supervisor.Storage"`
	networkingHandleClient cap.Networking_HandleClient
}

// New create supervisor with config
func New(config scf4go.Config) (smf4go.Service, error) {
	return &supervisorImpl{
		Logger: slf4go.Get("cap-supervisor"),
	}, nil
}

func (sr *supervisorImpl) Start() error {
	networkingHandleClient, err := sr.Networking.Handle(context.Background())

	if err != nil {
		return errors.Wrap(err, "call Networking#handle error")
	}

	sr.networkingHandleClient = networkingHandleClient

	go sr.networkingHandle()

	return nil
}

func (sr *supervisorImpl) networkingHandle() {

}
