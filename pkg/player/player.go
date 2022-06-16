package player

import (
	log "github.com/pion/ion-log"
	"github.com/pion/ion/pkg/proto"
	"github.com/pion/webrtc/v3"
	"github.com/yindaheng98/dion/pkg/islb"
	"github.com/yindaheng98/dion/pkg/sfu"
	"github.com/yindaheng98/dion/pkg/sxu/room"
	pb "github.com/yindaheng98/dion/proto"
	"github.com/yindaheng98/dion/util"
)

type Player struct {
	node      *islb.Node
	nats      string
	watchExec *util.SingleExec
	sfu       *sfu.Subscriber
	room      *room.Client
	switcher  *ClientSwitcher
}

func NewPlayerWithID(nats, nid string) Player {
	node := islb.NewNode(nid)
	switcher := NewClientSwitcher(&node)
	return Player{
		node:      &node,
		nats:      nats,
		watchExec: util.NewSingleExec(),
		sfu:       sfu.NewSubscriber(&node),
		room:      room.NewClient(switcher),
		switcher:  switcher,
	}
}

func NewPlayer(nats string) Player {
	return NewPlayerWithID(nats, "player-"+util.RandomString(16))
}

func (p Player) Name() string {
	return "player.Player"
}

func (p Player) Connect() {
	err := p.node.Start(p.nats)
	if err != nil {
		panic(err)
	}
	p.watchExec.Do(func() {
		err := p.node.Watch(proto.ServiceALL)
		if err != nil {
			log.Errorf("Node.Watch(proto.ServiceALL) error %v", err)
		}
	})
	p.room.Connect()
	p.sfu.Connect()
}

func (p Player) Close() {
	p.sfu.Close()
	p.room.Close()
	p.node.Close()
}

func (p Player) Connected() bool {
	return p.sfu.Connected() && p.room.Connected()
}

func (p Player) SwitchSession(session *pb.ClientNeededSession) {
	p.room.UpdateSession(session)
	p.sfu.SwitchSession(session)
}

func (p Player) SwitchNode(peerNID string, parameters map[string]interface{}) {
	p.switcher.Switch(peerNID, parameters)
	p.room.RefreshConn()
	p.sfu.SwitchNode(peerNID, parameters)
}

func (p Player) Switch(peerNID string, parameters map[string]interface{}, session *pb.ClientNeededSession) {
	p.room.UpdateSession(session)
	p.switcher.Switch(peerNID, parameters)
	p.room.RefreshConn()
	p.sfu.Switch(peerNID, parameters, session)
}

func (p Player) OnTrack(f func(remote *webrtc.TrackRemote, receiver *webrtc.RTPReceiver)) {
	p.sfu.OnTrack = f
}
