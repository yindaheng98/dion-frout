package player

import (
	"github.com/cloudwebrtc/nats-discovery/pkg/discovery"
	log "github.com/pion/ion-log"
	"github.com/pion/ion/pkg/proto"
	"github.com/yindaheng98/dion/config"
	"github.com/yindaheng98/dion/pkg/islb"
	"github.com/yindaheng98/dion/pkg/sfu"
	"github.com/yindaheng98/dion/util"
)

type Player struct {
	*sfu.Subscriber
	node      *islb.Node
	aliveExec *util.SingleExec
}

func NewPlayerWithID(nid string) Player {
	node := islb.NewNode(nid)
	return Player{
		Subscriber: sfu.NewSubscriber(&node),
		node:       &node,
		aliveExec:  util.NewSingleExec(),
	}
}

func NewPlayer() Player {
	return NewPlayerWithID("player-" + util.RandomString(16))
}

func (p Player) Name() string {
	return "player.Player"
}

func (p Player) Start(conf config.Common) {
	err := p.node.Start(conf.Nats.URL)
	if err != nil {
		panic(err)
	}
	err = p.node.Watch(proto.ServiceALL)
	if err != nil {
		log.Errorf("Node.Watch(proto.ServiceALL) error %v", err)
	}
	p.aliveExec.Do(func() {
		err = p.node.KeepAlive(discovery.Node{
			DC:      conf.Global.Dc,
			Service: config.ServiceClient,
			NID:     p.node.NID,
			RPC: discovery.RPC{
				Protocol: discovery.NGRPC,
				Addr:     conf.Nats.URL,
				//Params:   map[string]string{"username": "foo", "password": "bar"},
			},
		})
		if err != nil {
			log.Errorf("Node.KeepAlive error %v", err)
		}
	})
}
