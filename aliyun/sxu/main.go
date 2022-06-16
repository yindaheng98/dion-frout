package main

import (
	"flag"
	"fmt"
	sfu2 "github.com/pion/ion-sfu/pkg/sfu"
	"github.com/yindaheng98/dion-frout/algorithms/processor"
	"github.com/yindaheng98/dion/pkg/islb"
	"github.com/yindaheng98/dion/pkg/sfu"
	"github.com/yindaheng98/dion/pkg/sxu"
	"github.com/yindaheng98/dion/pkg/sxu/syncer"
	"github.com/yindaheng98/dion/util"

	"os"
	"os/signal"
	"syscall"

	log "github.com/pion/ion-log"
)

var (
	conf = sfu.Config{}
	file string
)

func showHelp() {
	fmt.Printf("Usage:%s {params}\n", os.Args[0])
	fmt.Println("      -c {config file}")
	fmt.Println("      -h (show help info)")
	os.Exit(-1)
}

func main() {
	var id, ffmpeg, bandwidth, filter string
	flag.StringVar(&id, "id", "sxu-"+util.RandomString(8), "id of sxu")
	flag.StringVar(&ffmpeg, "ffmpeg", "ffmpeg", "path to ffmpeg executable")
	flag.StringVar(&ffmpeg, "bandwith", "bandwith", "encode bandwidth")
	flag.StringVar(&filter, "filter", "drawtext=text='%{localtime\\:%Y-%m-%d %H.%M.%S}':fontsize=60:x=(w-text_w)/2:y=0", "ffmpeg -vf ???")
	flag.StringVar(&file, "c", "cmd/sxu/sfu.toml", "config file")
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

	log.Infof("--- starting isglb node ---")
	node := sxu.NewWithID(id, sxu.NewDefaultToolBoxBuilder(func(box *syncer.ToolBox, node *islb.Node, i *sfu2.SFU) {
		f := processor.NewFFmpegIVFProcessorFactory(ffmpeg)
		f.Filter = filter
		f.Bandwidth = bandwidth
		box.TrackProcessor = sxu.NewProceedRouter(i, f)
	}))
	if err := node.Start(conf); err != nil {
		log.Errorf("isglb start error: %v", err)
		os.Exit(-1)
	}
	defer node.Close()

	// Press Ctrl+C to exit the process
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}
