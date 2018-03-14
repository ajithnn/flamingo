package components

import (
	"fmt"
	"github.com/golang/glog"
	"os"
	"path"
	"strings"
)

type Video struct {
	mediaPath string
}

func (m Video) Process(filepath string, config interface{}, postProcess func()) {
	defer postProcess()
	glog.V(2).Infof("File path %s Video file is being processed.", filepath)

	parsedConfig := config.(map[string]interface{})

	inPath := parsedConfig["in_path"].(string)
	outPath := parsedConfig["out_path"].(string)
	errPath := parsedConfig["err_path"].(string)
	baseUrl := parsedConfig["api_base"].(string)
	token := parsedConfig["auth"].(string)

	s3Access := parsedConfig["access_key"].(string)
	s3Secret := parsedConfig["secret"].(string)
	bucket := parsedConfig["bucket"].(string)

	feed := parsedConfig["id"].(string)
	accDomain := parsedConfig["domain"].(string)

	stateEndpoint := parsedConfig["state_endpoint"].(string)
	updateEndpoint := parsedConfig["update_endpoint"].(string)
	createEndpoint := parsedConfig["create_endpoint"].(string)

	assetId := strings.Replace(path.Base(filepath), path.Ext(filepath), "", -1)
	stateUrl := fmt.Sprintf("%s%s", baseUrl, stateEndpoint)
	stateParams := fmt.Sprintf("auth_token=%s&assets=%s&feed_id=%s&account_id=%s", token, assetId, feed, accDomain)
	createUrl := fmt.Sprintf("%s%s", baseUrl, createEndpoint)

	glog.V(2).Infof("Retrieving State for %s Calling %s.", assetId, stateUrl)
	state, err := GetAssetState(filepath, stateUrl, stateParams)
	state.Filename = strings.Replace(state.Filename, inPath, "", -1)

	updateUrl := fmt.Sprintf("%s%s%d.json", baseUrl, updateEndpoint, int(state.Id))
	updateParams := fmt.Sprintf("auth_token=%s&feed_id=%s&account_id=%s", token, feed, accDomain)

	switch state.Status {
	case "new":
		status, err := UploadFile(filepath, updateUrl, updateParams, s3Access, s3Secret, bucket)
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

		fileObject, err := os.Stat(filepath)

		if err != nil {
			glog.V(2).Infof("Error file not found to find size Error: %s", err.Error())
			return
		}

		state.Md5sum = fmt.Sprintf("%x", md5)
		state.Status = "new"
		state.TotalSize = float64(fileObject.Size())
		glog.V(2).Infof("Md5: %s Status: %s Size: %d", state.Md5sum, state.Status, state.TotalSize)

		newState := CreateAsset(createUrl, token, feed, accDomain, state)
		if newState.Status == "new" {
			glog.V(2).Infof("Successfully Updated New Md5Sum, Will upload in next Cycle.")
			return
		}

	default:
		// Error in State , Retry In the next run. No Op
		if err != nil {
			glog.V(2).Infof("Processing failed for ", filepath, "Moving file to error folder.")
			glog.V(2).Infof(err.Error())
			err = os.Rename(filepath, path.Join(outPath, path.Base(filepath)))
			if err != nil {
				glog.V(2).Infof("Error Movement failed ", err, " deleting file.")
				os.Remove(filepath)
			}
		} else {
			glog.V(2).Infof("Successfully complete processing for ", filepath)
			err = os.Rename(filepath, path.Join(outPath, path.Base(filepath)))
			if err != nil {
				glog.V(2).Infof("Error Movement failed ", err, " deleting file.")
				os.Remove(filepath)
			}
		}
	}
	return
}
