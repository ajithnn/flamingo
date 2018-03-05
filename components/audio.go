package components

import (
	"github.com/golang/glog"
	//  "os"
	//  "os/exec"
	//  "path"
	//  "strings"
)

type Audio struct {
}

func (a Audio) Process(filepath string, config interface{}, postProcess func()) {
	defer postProcess()
	glog.V(2).Infof("File path ", filepath, " Media file is being processed.")
	return
}
