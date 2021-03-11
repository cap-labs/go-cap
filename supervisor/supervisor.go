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
	Application            cap.ApplicationClient `inject:"caplabs.supervisor.Application"`
	Networking             cap.NetworkingClient  `inject:"caplabs.supervisor.Networking"`
	Consensus              cap.ConsensusClient   `inject:"caplabs.supervisor.Consensus"`
	Storage                cap.StorageClient     `inject:"caplabs.supervisor.Storage"`
	networkingHandleClient cap.Networking_HandleClient
}

// New create supervisor with config
func New(config scf4go.Config) (smf4go.Service, error) {
	return &supervisorImpl{
		Logger: slf4go.Get("cap-supervisor"),
	}, nil
}

func (sr *supervisorImpl) Start() error {
	// networkingHandleClient, err := sr.Networking.Handle(context.Background())

	// if err != nil {
	// 	return errors.Wrap(err, "call Networking#handle error")
	// }

	// sr.networkingHandleClient = networkingHandleClient

	// go sr.networkingHandle()

	return nil
}

func (sr *supervisorImpl) networkingHandle() {
	for {
		msg, err := sr.networkingHandleClient.Recv()

		if err != nil {
			sr.E("networking handle error: {@error}", err)
			continue
		}

		sr.D("handle msg from {@from}", msg.From)

		consensusResponse, err := sr.Consensus.Handle(context.Background(), &cap.ConsensusRequest{
			Content: msg.Content,
		})

		if err != nil {
			sr.E("consensus handle error: {@error}", err)
			continue
		}

		switch v := consensusResponse.Content.(type) {
		case *cap.ConsensusResponse_AppRequest:

			err := sr.handleAppRequest(v.AppRequest)

			if err != nil {
				sr.E("handleAppRequest error: {@error}", err)
				continue
			}

		case *cap.ConsensusResponse_SendMessage:
			err := sr.handleSendMessage(v.SendMessage)

			if err != nil {
				sr.E("handleAppRequest error: {@error}", err)
				continue
			}
		}
	}

}

func (sr *supervisorImpl) handleSendMessage(request *cap.Message) error {
	err := sr.networkingHandleClient.Send(request)

	if err != nil {
		return errors.Wrap(err, "networking#send error")
	}

	return nil
}

func (sr *supervisorImpl) handleAppRequest(request *cap.ApplicationRequest) error {
	resp, err := sr.Application.Handle(context.Background(), request)

	if err != nil {
		return errors.Wrap(err, "application#handle error")
	}

	if resp.StateSet == nil {
		return nil
	}

	_, err = sr.Storage.StateSet(context.Background(), resp.StateSet)

	return err
}
