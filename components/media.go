package components

import (
  "github.com/golang/glog"
//  "os"
//  "os/exec"
//  "path"
//  "strings"
)

type Media struct {
}

func (m Media) Process(filepath string, postProcess func()) {
  defer postProcess()
  glog.V(2).Info("File path ", filepath, " Media file is being processed.")
    // TODO:
    // Check state on Cloud , Is uploadable ?
    // If not uploadable, Create/Update on blip as required and release lock. 
    // Next run will be used for upload.
    // If not known what to do , Move file to Track.
  return
}
