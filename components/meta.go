package components

import (
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"path"
)

type Meta struct {
	metaPath string
}

func (m Meta) Process(filepath string, config interface{}, postProcess func()) {
	defer postProcess()
	glog.V(2).Infof("Processing Meta file ", filepath)
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		glog.V(2).Infof("Error reading meta file ", filepath)
	} else {
		glog.V(2).Infof("Length of meta file ", filepath, " is ", len(dat))
		err = os.Rename(filepath, path.Join("outbox", "meta", path.Base(filepath)))
		if err != nil {
			glog.V(2).Infof("Error moving meta ", filepath, " error is ", err)
		}
	}
	return
}
