package player

import (
	"github.com/cloudwebrtc/nats-grpc/pkg/rpc"
	log "github.com/pion/ion-log"
	"github.com/yindaheng98/dion/config"
	"github.com/yindaheng98/dion/pkg/islb"
	"sync/atomic"
)

type ClientSwitcher struct {
	*islb.Node
	nid   atomic.Value
	param atomic.Value
}

func NewClientSwitcher(node *islb.Node) *ClientSwitcher {
	nid, param := atomic.Value{}, atomic.Value{}
	nid.Store("*")
	param.Store(map[string]interface{}{})
	return &ClientSwitcher{
		Node:  node,
		nid:   nid,
		param: param,
	}
}

func (f ClientSwitcher) NewClient() *rpc.Client {
	nid, param := f.nid.Load().(string), f.param.Load().(map[string]interface{})
	for _, node := range f.GetNeighborNodes() {
		if node.Service == config.ServiceSXU && (node.NID == nid || nid == "*") {
			client, err := f.NewNatsRPCClient(config.ServiceSXU, node.NID, param)
			if err != nil {
				log.Errorf("cannot NewNatsRPCClient: %v, try next", err)
			}
			return client
		}
	}
	return nil
}

func (f ClientSwitcher) Switch(peerNID string, parameters map[string]interface{}) {
	f.nid.Store(peerNID)
	f.param.Store(parameters)
}
