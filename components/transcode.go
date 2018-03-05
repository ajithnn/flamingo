package components

import (
	"github.com/golang/glog"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Transcode struct {
}

func (t Transcode) Process(filepath string, config interface{}, postProcess func()) {
	defer postProcess()
	glog.V(2).Infof("File path ", filepath, " Media file is being processed.")
	cmd := exec.Command("ffmpeg", "-y", "-i", filepath, path.Join("Inbox", "Track", strings.Split(path.Base(filepath), ".")[0]+".ts"))
	_, err := cmd.CombinedOutput()
	if err != nil {
		glog.V(2).Infof("Processing failed for ", filepath, "Moving file to error folder.")
		glog.V(2).Infof(err.Error())
		err = os.Rename(filepath, path.Join("outbox", "errors", path.Base(filepath)))
		if err != nil {
			glog.V(2).Infof("Error Movement failed ", err)
		}
	} else {
		glog.V(2).Infof("Successfully complete processing for ", filepath)
		err = os.Rename(filepath, path.Join("outbox", "transcode", path.Base(filepath)))
		if err != nil {
			glog.V(2).Infof("Error Movement failed ", err)
		}
	}
	return
}
