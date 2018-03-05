package main

import (
	"flag"
	"github.com/ajithnn/flamingo/components"
	"github.com/ajithnn/go-flow/flow"
	"github.com/golang/glog"
	"os"
	"strconv"
	"time"
)

func init() {
	flag.Parse()
}

func main() {

	inputArgs := flag.Args()[0:]
	if len(inputArgs) != 5 {

		glog.V(2).Infof("Usage:")
		glog.V(2).Infof("go run flamingo.go -logtostderr=true -v=2 <Inbox Path> <pipe json path> <Comma separated whitelist of folders> <Stable file Timeout Seconds> <Scan folder repeat Timeout seconds>")
		os.Exit(1)
	}

	typeMap := map[string]flow.Stage{
		"Media":     components.Media{},
		"Meta":      components.Meta{},
		"Transcode": components.Transcode{},
		"Graphics":  components.Graphic{},
		"Subtitles": components.Subtitle{},
		"Audio":     components.Audio{},
		"Track":     components.Track{},
	}

	timeOut, err := strconv.ParseFloat(inputArgs[3], 64)

	if err != nil {
		glog.V(2).Infof(" Timeout Invalid Parameters")
		os.Exit(0)
	}

	scanTimeOut, err := time.ParseDuration(inputArgs[4])

	if err != nil {
		glog.V(2).Infof("Scan Timeout Invalid Parameters")
		os.Exit(0)
	}

	config := flow.Flow{
		inputArgs[0],
		inputArgs[1],
		inputArgs[2],
		timeOut,
		scanTimeOut,
		typeMap,
		GetPrioritizedList,
	}

	glog.V(2).Infof("Starting Flow Service")
	flow.Trigger(config)

}

func GetPrioritizedList(fileList []string) []string {
	return fileList
}
