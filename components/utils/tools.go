package utils

import (
	"bytes"
	"crypto/md5"
	"github.com/golang/glog"
	"github.com/minio/minio-go"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type State struct {
	Status   string
	AssetID  string
	Filename string
	Priority int
	Md5sum   string
}

type Progress struct {
	Uploaded int64
	Total    int64
}

func (progress *Progress) PrintValue(wg *sync.WaitGroup) {
	sleepDur, _ := time.ParseDuration("10s")
	for {
		cur := atomic.LoadInt64(&progress.Uploaded)
		glog.V(2).Infof("Progress is %d", int(cur))
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
	byteN := bytes.IndexByte(md5sumByte, 0)
	return string(md5sumByte[:byteN])

}

func UpdateAsset(state State) State {
	return State{}
}

func CreateAsset(filepath string) State {
	return State{}
}

func ValidateAsset(filepath string) State {
	return State{}
}

func GetAssetState(filepath string) (State, error) {
	curState := State{
		"new",
		"ASDF",
		filepath,
		100,
		"ad43f4567cbd321456bdaeaefcdefdefd3421567896",
	}
	return curState, nil

}

func getObjectKeyFromFilepath(filepath string) string {
	return strings.Replace(filepath, path.Dir(path.Dir(filepath))+"/", "", -1)
}

func UploadFile(filepath, key, secret, bucket string) (bool, error) {
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

	go progress.PrintValue(&wg)

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
