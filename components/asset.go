package components

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/minio/minio-go"
	"github.com/parnurzeal/gorequest"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type State struct {
	Status    string
	AssetID   string
	Id        float64
	Filename  string
	Priority  float64
	TotalSize float64
	Md5sum    string
}

type Progress struct {
	Uploaded int64
	Total    int64
}

func (progress *Progress) UpdateProgress(url, params string, wg *sync.WaitGroup) {
	sleepDur, _ := time.ParseDuration("10s")
	for {
		cur := atomic.LoadInt64(&progress.Uploaded)
		glog.V(2).Infof("Progress is %d", int(cur))
		postBody := fmt.Sprintf(`{"uploaded_size": %d, "total_size": %d}`, int(cur), int(progress.Total))
		callEndPoint(url, params, http.MethodPut, postBody)
		if cur >= progress.Total {
			wg.Done()
			return
		}
		time.Sleep(sleepDur)
	}
}

func (progress *Progress) Read(p []byte) (int, error) {
	n := len(p)
	atomic.StoreInt64(&progress.Uploaded, atomic.AddInt64(&progress.Uploaded, int64(len(p))))
	return n, nil
}

func CalculateMd5sum(filepath string) string {
	hash := md5.New()
	file, _ := os.Open(filepath)
	defer file.Close()
	io.Copy(hash, file)
	md5sumByte := hash.Sum(nil)
	return string(md5sumByte)

}

func UpdateAsset(state State) State {
	return State{}
}

func CreateAsset(url string, token string, feed string, acc_domain string, state State) State {

	var newState []map[string]interface{}

	assetId := strings.Replace(state.Filename, path.Ext(state.Filename), "", -1)
	params := fmt.Sprintf(`{"auth_token":"%s","feed_id":"%s","account_id": "%s","filename":"%s","md5sum":"%s","size":"%f","asset_id":"%s"}`, token, feed, acc_domain, state.Filename, state.Md5sum, state.TotalSize, path.Base(assetId))
	resp, err := callEndPoint(url, params, http.MethodPost, `{}`)

	if err != nil {
		return State{}
	}

	respBytes := []byte(resp)

	json.Unmarshal(respBytes, &newState)

	if len(newState) == 0 {
		return State{
			"not_present",
			state.Filename,
			-1,
			state.Filename,
			-1,
			-1,
			"n/a",
		}
	}

	curAssetState := newState[0]

	priority := 0.00
	if curAssetState["Priority"] != nil {
		priority = curAssetState["priority"].(float64)
	}

	curState := State{
		curAssetState["state"].(string),
		curAssetState["asset_id"].(string),
		curAssetState["id"].(float64),
		curAssetState["filename"].(string),
		priority,
		curAssetState["size"].(float64),
		curAssetState["md5sum"].(string),
	}

	return curState
}

func ValidateAsset(filepath string) State {
	return State{}
}

func GetAssetState(filepath, stateUrl, stateParams string) (State, error) {
	var state []map[string]interface{}
	resp, err := callEndPoint(stateUrl, stateParams, http.MethodGet, `{}`)

	if err != nil {
		return State{}, err
	}

	respBytes := []byte(resp)

	json.Unmarshal(respBytes, &state)

	if len(state) == 0 {
		return State{
			"not_present",
			filepath,
			-1,
			filepath,
			-1,
			-1,
			"n/a",
		}, err
	}

	curAssetState := state[0]

	priority := 0.00
	if curAssetState["Priority"] != nil {
		priority = curAssetState["priority"].(float64)
	}

	curState := State{
		curAssetState["state"].(string),
		curAssetState["asset_id"].(string),
		curAssetState["id"].(float64),
		filepath,
		priority,
		curAssetState["size"].(float64),
		curAssetState["md5sum"].(string),
	}

	return curState, nil

}

func getObjectKeyFromFilepath(filepath string) string {
	return strings.Replace(filepath, path.Dir(path.Dir(filepath))+"/", "", -1)
}

func callEndPoint(url string, queryParams string, method string, body string) (string, error) {

	request := gorequest.New()

	glog.V(2).Infof("%s %s %s %s", url, queryParams, method, body)
	resp, respBody, errs := request.CustomMethod(method, url).
		Query(queryParams).
		Send(body).
		End()

	if len(errs) != 0 {
		return "", errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Status is " + string(resp.Status))
	}

	return string(respBody), nil
}

func UploadFile(filepath, updateUrl, updateParams, key, secret, bucket string) (bool, error) {
	ssl := true
	var wg sync.WaitGroup

	objectKey := getObjectKeyFromFilepath(filepath)
	s3Client, err := minio.New("s3.amazonaws.com", key, secret, ssl)
	if err != nil {
		glog.V(2).Infof("Encoutered error %s", err.Error())
		return false, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		glog.V(2).Infof("Encoutered error %s", err.Error())
		return false, err
	}

	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		glog.V(2).Infof("Encoutered error %s", err.Error())
		return false, err
	}

	progress := Progress{Total: fileStat.Size(), Uploaded: 0}

	wg.Add(1)

	go progress.UpdateProgress(updateUrl, updateParams, &wg)

	bytesUploaded, err := s3Client.PutObject(bucket, objectKey, file, fileStat.Size(), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
		Progress:    &progress,
	})

	wg.Wait()

	if err != nil {
		glog.V(2).Infof("Encoutered error %s", err.Error())
		return false, err
	}
	glog.V(2).Infof("Successfully uploaded bytes: %d", bytesUploaded)
	return true, nil
}
