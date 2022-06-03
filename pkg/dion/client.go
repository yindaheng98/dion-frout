package client

import (
	nrpc "github.com/cloudwebrtc/nats-grpc/pkg/rpc"
	log "github.com/pion/ion-log"
	"github.com/yindaheng98/dion/config"
	"github.com/yindaheng98/dion/pkg/islb"
	"github.com/yindaheng98/dion/pkg/sfu"
	"github.com/yindaheng98/dion/pkg/sxu/room"
	pb2 "github.com/yindaheng98/dion/proto"
	"github.com/yindaheng98/dion/util"
	"sync/atomic"
)

type conf struct {
	session    *pb2.ClientNeededSession
	peerNID    string
	parameters map[string]interface{}
	version    uint32
}

type HealthWithSubscriber struct {
	*islb.Node
	sub     *sfu.Subscriber
	conf    atomic.Value
	current uint32
}

func (h HealthWithSubscriber) NewClient() *nrpc.Client {
	c := h.conf.Load()
	if c != nil {
		return nil
	}
	peerNID, parameters, session := c.(conf).peerNID, c.(conf).parameters, c.(conf).session
	client, err := h.NewNatsRPCClient(config.ServiceSXU, peerNID, parameters)
	if err != nil {
		log.Errorf("cannot NewNatsRPCClient: %v, try next", err)
	}
	current := atomic.LoadUint32(&h.current)
	if current != c.(conf).version {
		h.sub.SwitchNode(session, peerNID, parameters)
		atomic.StoreUint32(&h.current, c.(conf).version)
	}
	return client
}

func (h HealthWithSubscriber) Switch(session *pb2.ClientNeededSession, peerNID string, parameters map[string]interface{}) {
	current := atomic.LoadUint32(&h.current)
	h.conf.Store(conf{
		session:    session,
		peerNID:    peerNID,
		parameters: parameters,
		version:    current + 1,
	})
}

type Client struct {
	*islb.Node
	HealthFactory HealthWithSubscriber
	SFU           *sfu.Subscriber
	Room          *room.Client
}

func NewClient(uid string) *Client {
	node := islb.NewNode(uid)
	sub := sfu.NewSubscriber(&node)
	health := HealthWithSubscriber{
		Node:    &node,
		sub:     sub,
		current: 0,
	}
	health.conf.Store(conf{
		session: &pb2.ClientNeededSession{
			Session: "stupid",
			User:    util.RandomString(8),
		},
		peerNID:    "*",
		parameters: map[string]interface{}{},
		version:    0,
	})
	return &Client{
		Node:          &node,
		HealthFactory: health,
		SFU:           sfu.NewSubscriber(&node),
		Room:          room.NewClient(health),
	}
}
