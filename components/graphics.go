package components

import (
  "github.com/golang/glog"
 // "os"
 // "os/exec"
 // "path"
 // "strings"
)

type Graphic struct {
}

func (g Graphic) Process(filepath string, postProcess func()) {
  defer postProcess()
  glog.V(2).Info("File path ", filepath, " Media file is being processed.")
  return
}
