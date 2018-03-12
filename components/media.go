package components

import (
	"fmt"
	"github.com/golang/glog"
	"os"
	"path"
	"strings"
)

type Media struct {
	mediaPath string
}

func (m Media) Process(filepath string, config interface{}, postProcess func()) {
	defer postProcess()
	glog.V(2).Infof("File path %s Media file is being processed.", filepath)

	parsedConfig := config.(map[string]interface{})

	inPath := parsedConfig["in_path"].(string)
	outPath := parsedConfig["out_path"].(string)
	errPath := parsedConfig["err_path"].(string)
	baseUrl := parsedConfig["api_base"].(string)
	token := parsedConfig["token"].(string)

	s3_access := parsedConfig["access_key"].(string)
	s3_secret := parsedConfig["secret"].(string)
	bucket := parsedConfig["bucket"].(string)

	feed := parsedConfig["feed_id"].(string)
	acc_domain := parsedConfig["account_domain"].(string)
	//acc_id := parsedConfig["account_id"].(string)

	assetId := strings.Replace(path.Base(filepath), path.Ext(filepath), "", -1)
	stateUrl := fmt.Sprintf("%sts/assets/state.json", baseUrl)
	stateParams := fmt.Sprintf("auth_token=%s&assets=%s&feed_id=%s&account_id=%s", token, assetId, feed, acc_domain)
	createUrl := fmt.Sprintf("%sts/assets.json", baseUrl)

	glog.V(2).Infof("Retrieving State for %s Calling %s.", assetId, stateUrl)
	state, err := GetAssetState(filepath, stateUrl, stateParams)
	state.Filename = strings.Replace(state.Filename, inPath, "", -1)
	updateUrl := fmt.Sprintf("%sts/assets/%d.json", baseUrl, int(state.Id))
	updateParams := fmt.Sprintf("auth_token=%s&feed_id=%s&account_id=%s", token, feed, acc_domain)

	switch state.Status {
	case "new":
		status, err := UploadFile(filepath, updateUrl, updateParams, s3_access, s3_secret, bucket)
		if err != nil {
			glog.V(2).Infof("Error Uploading file, marking failed.")
			state.Status = "failed"
			UpdateAsset(state)
			err = os.Rename(filepath, path.Join(errPath, path.Base(filepath)))
		}
		if status {
			newPath := path.Dir(path.Dir(filepath))
			glog.V(2).Infof("Completed Upload moving to Track ", newPath)
			err = os.Rename(filepath, path.Join(errPath, path.Base(filepath)))
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
		state.Md5sum = fmt.Sprintf("%x", md5)
		newState := UpdateAsset(state)
		if newState.Status == "new" {
			glog.V(2).Infof("Successfully Updated New Md5Sum, Will upload in next Cycle.")
			return
		}
	case "not_present":

		// Create Entry on Cloud and update Md5sum , move to new state.
		md5 := CalculateMd5sum(filepath)
		if md5 == state.Md5sum {
			err = os.Rename(filepath, path.Join(outPath, path.Base(filepath)))
			return
		}

		state.Md5sum = fmt.Sprintf("%x", md5)
		state.Status = "new"
		state.TotalSize = 12345678
		glog.V(2).Infof("Md5: %s Status: %s Size: %d", state.Md5sum, state.Status, state.TotalSize)

		newState := CreateAsset(createUrl, token, feed, acc_domain, state)
		if newState.Status == "new" {
			glog.V(2).Infof("Successfully Updated New Md5Sum, Will upload in next Cycle.")
			return
		}

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
