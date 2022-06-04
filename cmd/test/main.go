package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"github.com/cloudwebrtc/nats-discovery/pkg/discovery"
	"github.com/yindaheng98/dion-frout/pkg/controller"
	"github.com/yindaheng98/dion/algorithms/impl/random"
	pb "github.com/yindaheng98/dion/proto"
	"github.com/yindaheng98/dion/util"
)

type TestClient struct {
	nodes map[string]discovery.Node
}

func NewTestClient(n uint) TestClient {
	c := TestClient{nodes: map[string]discovery.Node{}}
	for i := uint(0); i < n; i++ {
		id := util.RandomString(8)
		c.nodes[id] = discovery.Node{
			DC:        util.RandomString(8),
			Service:   "test-service",
			NID:       util.RandomString(8),
			RPC:       discovery.RPC{},
			ExtraInfo: nil,
		}
	}
	return c
}

func (t TestClient) Connect(nats_addr string) error {
	fmt.Println(nats_addr)
	return nil
}

func (t TestClient) GetNodes() map[string]discovery.Node {
	for id, node := range t.nodes {
		if random.RandBool() {
			node.DC = util.RandomString(8)
			t.nodes[id] = node
		}
	}
	return t.nodes
}

func (t TestClient) SwitchNode(id string) {
	fmt.Printf("Switch to: %s\n", id)
}

func (t TestClient) SwitchSession(session *pb.ClientNeededSession) {
	//TODO implement me
	panic("implement me")
}

func main() {
	a := app.New()
	w := a.NewWindow("dion system")
	go controller.Control(w, NewTestClient(4))
	w.ShowAndRun()
}
