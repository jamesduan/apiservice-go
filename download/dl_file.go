// Copyright 2018 JingDong, Inc.
// Written by Zhangyunyang 2018/01/10
//
// DownloadFile
package download

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"apiservice/g"
)

//file basic information
type BasicFileInfo struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Path     string `json:"path"`
	Size     int64  `json:"size"` //total size
	CheckSum []byte `json:"checksum"`
	Step     int    `json:"step"`
}

func (this *BasicFileInfo) String() string {
	return fmt.Sprintf("<Id:%d, Name:%s, Version:%s, Path:%s, Size:%d, CheckSum:%x, Step:%d>",
		this.Id,
		this.Name,
		this.Version,
		this.Path,
		this.Size,
		this.CheckSum,
		this.Step,
	)
}

type DownloadFile struct {
	sync.RWMutex
	fileName string // fullname
	fileId   int64  //DB table pk
	fileDesc *os.File
}

func NewDownloadFile(fileName string) (*DownloadFile, error) {
	var (
		fileDesc *os.File
	)

	download_dir := g.Config().Download
	download_dir = strings.TrimSpace(download_dir)
	if download_dir[len(download_dir)-1] != '/' {
		download_dir += "/"
	}

	fullName := download_dir + fileName + ".tar.gz"

	if fd, err := os.Open(fullName); err != nil {
		log.Printf("Open failed, file:%s, error:%v\n", fullName, err)
		return nil, err
	} else {
		fileDesc = fd
	}

	dlf := &DownloadFile{fileName: fullName, fileDesc: fileDesc}
	return dlf, nil
}

func (dlf *DownloadFile) GetFD() (*os.File, error) {
	dlf.RLock()
	defer dlf.RUnlock()

	if dlf.fileDesc == nil {
		log.Println("FD is nil, file name:", dlf.fileName)
		return nil, fmt.Errorf("FD is nil")
	}

	return dlf.fileDesc, nil
}

func (dlf *DownloadFile) Destroy() error {
	dlf.Lock()
	defer dlf.Unlock()

	if dlf.fileDesc != nil {
		dlf.fileDesc.Close()
		dlf.fileDesc = nil
	}

	return nil
}
