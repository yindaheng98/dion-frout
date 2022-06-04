package client

import (
	"bufio"
	"fmt"
	"github.com/cloudwebrtc/nats-discovery/pkg/discovery"
	"github.com/pion/ion/pkg/proto"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/ivfwriter"
	"github.com/yindaheng98/dion/pkg/islb"
	"github.com/yindaheng98/dion/pkg/sfu"
	pb "github.com/yindaheng98/dion/proto"
	"github.com/yindaheng98/dion/util"
	"io"
	"log"
	"os/exec"
)

type Client struct {
	islb.Node
	watchExec util.SingleExec
	sub       *sfu.Subscriber
}

func (h *Client) Connect(addr, cmd string) error {
	log.Println("Connecting......")
	err := h.Node.Start(addr)
	if err != nil {
		return err
	}
	h.watchExec.Do(func() {
		log.Println("Start watching......")
		err := h.Node.Watch(proto.ServiceALL)
		if err != nil {
			log.Fatalf("Node.Watch(proto.ServiceALL) error %v\n", err)
		}
	})
	h.sub = sfu.NewSubscriber(&h.Node)
	h.sub.OnTrack = func(remote *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("OnTrack started: %+v\n", remote)
		ffplay := exec.Command(cmd, "-f", "ivf", "-i", "pipe:0")
		defer ffplay.Process.Kill()
		stdin, stdout, err := util.GetStdPipes(ffplay)
		if err != nil {
			log.Fatalf("Cannot start ffplay: %+v\n", err)
		}
		go func(stdout io.ReadCloser) {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}(stdout)
		ivfWriter, err := ivfwriter.NewWith(stdin)
		if err != nil {
			log.Fatalf("Cannot create ivfwriter: %+v\n", err)
		}

		for {
			// Read RTP packets being sent to Pion
			rtp, _, readErr := remote.ReadRTP()
			log.Println("Subscriber get a RTP Packet")
			if readErr != nil {
				log.Printf("Subscriber RTP Packet read error %+v\n", readErr)
				return
			}

			if ivfWriterErr := ivfWriter.WriteRTP(rtp); ivfWriterErr != nil {
				log.Printf("RTP Packet write error: %+v\n", ivfWriterErr)
				return
			}
		}
	}
	return nil
}

func (h *Client) GetNodes() map[string]discovery.Node {
	return h.Node.GetNeighborNodes()
}

func (h *Client) SwitchNode(id string) {
	fmt.Printf("Switch node to: %s\n", id)
	h.sub.SwitchNode(id, map[string]interface{}{})
}

func (h *Client) SwitchSession(session *pb.ClientNeededSession) {
	fmt.Printf("Switch session to: %+v\n", session)
	h.sub.SwitchSession(session)
}
