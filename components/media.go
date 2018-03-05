package components

import (
	. "github.com/ajithnn/flamingo/components/utils"
	"github.com/golang/glog"
	"os"
	"path"
)

type Media struct {
	mediaPath string
}

func (m Media) Process(filepath string, config interface{}, postProcess func()) {
	defer postProcess()
	glog.V(2).Infof("File path %s Media file is being processed.", filepath)

	parsedConfig := config.(map[string]interface{})
	inbox := parsedConfig["inbox_path"].(string)
	outbox := parsedConfig["outbox_path"].(string)

	glog.V(2).Infof("Parsed Config Inbox is %s Outbox is %s.", inbox, outbox)

	state, err := GetAssetState(filepath)

	switch state.Status {
	case "new":
		// Start Uploading
		status, err := UploadFile(filepath, parsedConfig["access_key"].(string), parsedConfig["secret"].(string), parsedConfig["bucket"].(string))
		if err != nil {
			glog.V(2).Infof("Error Uploading file, marking failed.")
			state.Status = "failed"
			UpdateAsset(state)
			err = os.Rename(filepath, path.Join(outbox, path.Base(filepath)))
		}
		if status {
			newPath := path.Dir(path.Dir(filepath))
			glog.V(2).Infof("Completed Upload moving to Track ", newPath)
			err = os.Rename(filepath, path.Join(outbox, path.Base(filepath)))
		}
	case "uploading":
		// Start Move to Track Folder
		glog.V(2).Infof("File path ", filepath, " Media file is being uploaded.Nothing to do.")
	case "transcoding":
		// Start Move to Track Folder
		glog.V(2).Infof("File path ", filepath, " Media file is being transcoded.Nothing to do.")
	case "transcoded":
		// Start Move to Track Folder
		err = os.Rename(filepath, path.Join("./Inbox", "Track", path.Base(filepath)))
	case "processing":
		// Start Move to Track Folder
		err = os.Rename(filepath, path.Join("./Inbox", "Track", path.Base(filepath)))
	case "queued":
		// Start Move to Track Folder
		err = os.Rename(filepath, path.Join("./Inbox", "Track", path.Base(filepath)))
	case "uploaded":
		// Calculate Md5Sum and Set State to New if different, Else Move to Outbox
		md5 := CalculateMd5sum(filepath)
		if md5 == state.Md5sum {
			err = os.Rename(filepath, path.Join("outbox", "media", path.Base(filepath)))
			return
		}
		state.Md5sum = md5
		newState := UpdateAsset(state)
		if newState.Status == "new" {
			glog.V(2).Infof("Successfully Updated New Md5Sum, Will upload in next Cycle.")
			return
		}
	case "initial":
		// Create Entry on Cloud and update Md5sum , move to new state.
	default:
		// Error in State , Retry In the next run. No Op
		if err != nil {
			glog.V(2).Infof("Processing failed for ", filepath, "Moving file to error folder.")
			glog.V(2).Infof(err.Error())
			err = os.Rename(filepath, path.Join("outbox", "errors", path.Base(filepath)))
			if err != nil {
				glog.V(2).Infof("Error Movement failed ", err, " deleting file.")
				os.Remove(filepath)
			}
		} else {
			glog.V(2).Infof("Successfully complete processing for ", filepath)
			err = os.Rename(filepath, path.Join("outbox", "media", path.Base(filepath)))
			if err != nil {
				glog.V(2).Infof("Error Movement failed ", err, " deleting file.")
				os.Remove(filepath)
			}
		}
	}
	return
}
