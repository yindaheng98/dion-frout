package main

import (
	"bufio"
	"flag"
	"fmt"
	log "github.com/pion/ion-log"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/ivfwriter"
	"github.com/yindaheng98/dion-frout/algorithms"
	"github.com/yindaheng98/dion-frout/pkg/player"
	"github.com/yindaheng98/dion/config"
	pb "github.com/yindaheng98/dion/proto"
	"github.com/yindaheng98/dion/util"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	conf = config.Common{}
	file string
)

func showHelp() {
	fmt.Printf("Usage:%s {params}\n", os.Args[0])
	fmt.Println("      -c {config file}")
	fmt.Println("      -h (show help info)")
	os.Exit(-1)
}

func main() {
	var ffplay, nid, sid, uid string
	flag.StringVar(&ffplay, "ffplay", "ffplay", "path to ffplay executable")
	flag.StringVar(&nid, "nid", algorithms.ServiceNameQingdao, "target node id")
	flag.StringVar(&sid, "sid", config.ServiceSessionStupid, "target session id")
	flag.StringVar(&uid, "uid", algorithms.UserPath, "your user id")
	flag.StringVar(&file, "c", "aliyun/conf/islb.toml", "config file")
	help := flag.Bool("h", false, "help info")
	flag.Parse()
	if *help {
		showHelp()
	}

	err := conf.Load(file)
	if err != nil {
		fmt.Printf("config file %s read failed. %v\n", file, err)
		showHelp()
	}

	fmt.Printf("config %s load ok!\n", file)

	log.Init(conf.Log.Level)

	sub := player.NewPlayer()
	sub.OnTrack = func(remote *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Infof("OnTrack started: %+v\n", remote)
		ffplay := exec.Command(ffplay, "-f", "ivf", "-i", "pipe:0")
		stdin, stdout, err := util.GetStdPipes(ffplay)
		if err != nil {
			panic(err)
		}
		defer ffplay.Process.Kill()
		go func(stdout io.ReadCloser) {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}(stdout)
		ivfWriter, err := ivfwriter.NewWith(stdin)
		if err != nil {
			panic(err)
		}

		for {
			// Read RTP packets being sent to Pion
			rtp, _, readErr := remote.ReadRTP()
			log.Infof("Subscriber get a RTP Packet")
			if readErr != nil {
				log.Errorf("Subscriber RTP Packet read error %+v", readErr)
				return
			}

			if ivfWriterErr := ivfWriter.WriteRTP(rtp); ivfWriterErr != nil {
				log.Errorf("RTP Packet write error: %+v", ivfWriterErr)
				return
			}
		}
	}
	sub.Start(conf)
	sub.Switch(nid, map[string]interface{}{}, &pb.ClientNeededSession{
		Session: sid,
		User:    uid,
	})

	// Press Ctrl+C to exit the process
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}
