package stream

import (
	"context"
	"io"
	"sync"

	"github.com/cap-labs/go-cap"
	"github.com/libs4go/errors"
	"github.com/libs4go/slf4go"
	"google.golang.org/grpc"
)

const errVendor = "cap.stream"

// errors
var (
	ErrClosed = errors.New("the channel closed", errors.WithVendor(errVendor))
)

// Builder the stream builde facade interface
type Builder interface {
	// Create create new stream channel by remote peer object
	Create(cap.Peer) (Channel, error)
}

// Channel peer stream channel for consensus engine send
type Channel interface {
	// Stream channel is closable
	io.Closer
	// Remote return the channel remote peer information
	Remote() cap.Peer
	// Recv recv message from remote peer
	Recv() ([]byte, error)
	// Send send message to remote peer
	Send([]byte) error
}

type capNetworkStreamBuilder struct {
	sync.RWMutex
	slf4go.Logger
	client       cap.NetworkStreamClient
	handleClient cap.NetworkStream_HandleClient
	context      context.Context
	recvBuffSize int
	callOps      []grpc.CallOption
	channels     map[string]*capNetworkStreamChannel
}

// Option bulder create option
type Option func(builder *capNetworkStreamBuilder)

// Ctx option context.Context
func Ctx(ctx context.Context) Option {
	return func(builder *capNetworkStreamBuilder) {
		builder.context = ctx
	}
}

// RecvBuffSize each channel recv message queue size
func RecvBuffSize(size int) Option {
	return func(builder *capNetworkStreamBuilder) {
		builder.recvBuffSize = size
	}
}

// CallOps grpc call ops for call cap.NetworkStreamClient Handle method
func CallOps(callOps ...grpc.CallOption) Option {
	return func(builder *capNetworkStreamBuilder) {
		builder.callOps = callOps
	}
}

// New create stream builder with cap.NetworkStreamClient
func New(client cap.NetworkStreamClient, options ...Option) (Builder, error) {
	builder := &capNetworkStreamBuilder{
		Logger:       slf4go.Get("cap-stream"),
		client:       client,
		context:      context.Background(),
		recvBuffSize: 10,
	}

	for _, option := range options {
		option(builder)
	}

	return builder, nil
}

func (builder *capNetworkStreamBuilder) recvLoop() error {
	client, err := builder.client.Handle(builder.context, builder.callOps...)

	if err != nil {
		return errors.Wrap(err, "call cap.NetworkStreamClient Handle error")
	}

	builder.handleClient = client

	go func() {

		for {
			packet, err := client.Recv()

			if err != nil {
				builder.E("call cap.NetworkStreamClient#Handle recv stream {@error}", err)
				continue
			}

			builder.dispatch(packet)
		}
	}()

	return nil
}

func (builder *capNetworkStreamBuilder) dispatch(packet *cap.StreamPacket) {
	builder.RLock()
	channel, ok := builder.channels[packet.PeerId]
	builder.RUnlock()

	if !ok {
		builder.W("peer({@peer}) channel not found", packet.PeerId)
		return
	}

	channel.in(packet)
}

func (builder *capNetworkStreamBuilder) Create(peer cap.Peer) (Channel, error) {
	builder.Lock()
	defer builder.Unlock()

	channel, ok := builder.channels[peer.Id]

	if ok {
		return channel, nil
	}

	channel = &capNetworkStreamChannel{
		recvChan: make(chan *cap.StreamPacket, builder.recvBuffSize),
		peer:     peer,
		builder:  builder,
	}

	builder.channels[peer.Id] = channel

	return channel, nil
}

type capNetworkStreamChannel struct {
	recvChan chan *cap.StreamPacket
	peer     cap.Peer
	builder  *capNetworkStreamBuilder
}

func (channel *capNetworkStreamChannel) Close() error {

	channel.builder.Lock()
	defer channel.builder.Unlock()

	delete(channel.builder.channels, channel.peer.Id)

	close(channel.recvChan)

	return nil
}

func (channel *capNetworkStreamChannel) Remote() cap.Peer {
	return channel.peer
}

func (channel *capNetworkStreamChannel) Recv() ([]byte, error) {
	packet, ok := <-channel.recvChan

	if !ok {
		return nil, io.EOF
	}

	return packet.Content, nil
}

func (channel *capNetworkStreamChannel) Send(content []byte) error {
	return channel.builder.handleClient.Send(&cap.StreamPacket{
		PeerId:  channel.peer.Id,
		Content: content,
	})
}

func (channel *capNetworkStreamChannel) in(packet *cap.StreamPacket) {
	defer func() {
		if err := recover(); err != nil {
			channel.builder.E("send to closed channel {@peer}", channel.peer.Id)
		}
	}()

	channel.recvChan <- packet
}
