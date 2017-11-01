package components

import (
  "github.com/golang/glog"
//  "os"
//  "os/exec"
//  "path"
//  "strings"
)

type Audio struct {
  mediaPath string
}

func (a Audio) Process(filepath string, postProcess func()) {
  defer postProcess()
  glog.V(2).Info("File path ", filepath, " Media file is being processed.")
  return
}
