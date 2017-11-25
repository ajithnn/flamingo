package components

import (
  "os"
  "crypto/md5"
  "bytes"
  "io"
)


type State struct{
  Status string
  AssetID string
  Filename string
  Priority int
  Md5sum string
}

func CalculateMd5sum(filepath string) string{
  hash := md5.New()
  file,_ := os.Open(filepath)
  defer file.Close()
  io.Copy(hash,file)
  md5sumByte := hash.Sum(nil)
  byteN := bytes.IndexByte(md5sumByte,0)
  return string(md5sumByte[:byteN])

}

func UpdateAsset(state State) State{
  return State{}
}

func CreateAsset(filepath string) State{
  return State{}
}

func ValidateAsset(filepath string) State{
  return State{}
}

func GetAssetState(filepath string) (State,error) {
 curState := State{
    "uploaded",
    "ASDF",
    filepath,
    100,
    "ad43f4567cbd321456bdaeaefcdefdefd3421567896",
  }
  return curState,nil

}


