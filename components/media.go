package components

import (
  "github.com/golang/glog"
  "os"
  "os/exec"
  "path"
  "strings"
)

type Media struct {
  mediaPath string
}

func (m Media) Process(filepath string, postProcess func()) {
  defer postProcess()
  glog.V(2).Info("File path ", filepath, " Media file is being processed.")

  state,err := getState(filepath)
  switch state.state {
  case "new":
    // Start Uploading
  case "uploading":
    // Start Move to Track Folder
      err = os.Rename(filepath, path.Join("Track", path.Base(filepath)))
  case "transcoding":
    // Start Move to Track Folder
      err = os.Rename(filepath, path.Join("Track", path.Base(filepath)))
  case "transcoded":
    // Start Move to Track Folder
      err = os.Rename(filepath, path.Join("Track", path.Base(filepath)))
  case "processing":
    // Start Move to Track Folder
      err = os.Rename(filepath, path.Join("Track", path.Base(filepath)))
  case "queued":
    // Start Move to Track Folder
      err = os.Rename(filepath, path.Join("Track", path.Base(filepath)))
  case "uploaded":
    // Calculate Md5Sum and Set State to New if different, Else Move to Outbox
    md5 := calculateMd5Sum(filepath)
    if md5 == state.md5{
      err = os.Rename(filepath, path.Join("outbox","media", path.Base(filepath)))
      return
    }
    newState := updateMd5Sum(filepath)
    if(newState.state == "new"){
      glog.V(2).Info("Successfully Updated New Md5Sum, Will upload in next Cycle.")
      return
    }
  case "initial":
    // Create Entry on Cloud and update Md5sum , move to new state.
  default:
    // Error in State , Retry In the next run. No Op
    if err != nil {
      glog.V(2).Info("Processing failed for ", filepath, "Moving file to error folder.")
      glog.V(2).Info(err)
      err = os.Rename(filepath, path.Join("outbox", "errors", path.Base(filepath)))
      if err != nil {
        glog.V(2).Info("Error Movement failed ", err, " deleting file.")
        os.Unlink(filepath)
      }
    } else {
      glog.V(2).Info("Successfully complete processing for ", filepath)
      err = os.Rename(filepath, path.Join("outbox","media", path.Base(filepath)))
      if err != nil {
        glog.V(2).Info("Error Movement failed ", err, " deleting file.")
        os.Unlink(filepath)
      }
    }
  }
  return
}
