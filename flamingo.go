package main

import (
  "github.com/golang/glog"
  "github.com/ajithnn/go-flow/flow"
  "flag"
)

func Init(){
  flag.Parse()
}

func main(){

  inputArgs := flag.Args()[0:]
  if len(inputArgs) != 5 {

    glog.V(2).Infof("Usage:")
    glog.V(2).Infof("go run scan_folder.go -logtostderr=true -v=2 <Inbox Path> <Comma separated whitelist of folders>")
    os.Exit(1)
  }

  typeMap :=  map[string]flow.Asset{
    "Media": Media{},
    "Meta": Meta{},
    "Transcode": Transcode{},
    "Graphics": Graphic{},
    "Subtitles": Subtitle{},
    "Audio": Audio{},
    "Track": Track{},
  }

  config := flow.Flow{
    inputArgs[0],
    inputArgs[1],
    inputArgs[2],
    float64(inputArgs[3]),
    int(inputArgs[4]),
    typeMap,
    getPrioritizedList
  }

  glog.V(2).Info("Starting Flow Service")
  flow.Trigger(config)

}

