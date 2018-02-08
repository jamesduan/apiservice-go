// Download

package model

import (
	"fmt"
)

// the status of download file
const (
	DLInvalid      = iota
	DLStart        // 1
	DLEnd          // 2
	DLGetBasicErr  // 3 get basic info error
	DLGetConcrtErr //   get concrete info error
	DLWrtFileErr   //   write file error
	DLSizeErr      //   file size error
	DLChkSumErr    //   checksum error
	DLDecmpErr     //   decompress error
)

type DLStatus int

type BasicInfoRequest struct {
	IP          string
	FileName    string
	FileVersion string
}

func (this *BasicInfoRequest) String() string {
	return fmt.Sprintf(
		"<IP:%s, FileName:%s, FileVersion:%s>",
		this.IP,
		this.FileName,
		this.FileVersion,
	)
}

//file basic information
type BasicInfoResponse struct {
	Size     int64  `json:"size"` //total size
	CheckSum []byte `json:"checksum"`
	Step     int    `json:"step"`
	Version  string `json:"version"`
}

type ConcreteInfoRequest struct {
	IP          string
	FileName    string
	FileVersion string
	Offset      int64
	Length      int
}

// file content
type ConcreteInfoResponse struct {
	Length  int64 // the length of the content
	Content []byte
}

type DownloadResultRequest struct {
	IP          string
	FileName    string
	FileVersion string
	Status      DLStatus
}
