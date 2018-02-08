// Copyright 2018 JingDong, Inc.
// Written by Zhangyunyang 2018/01/10
//
// Download File Set

package download

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

type DLFileSet struct {
	sync.RWMutex
	Files map[string]*DownloadFile
}

var GFileSet *DLFileSet

func NewFileSet() error {
	GFileSet = &DLFileSet{Files: make(map[string]*DownloadFile)}
	return nil
}

//download file set
func (fs *DLFileSet) Get(fileKey string) (*DownloadFile, error) {
	fs.RLock()
	defer fs.RUnlock()

	f, ok := fs.Files[fileKey]
	if ok {
		return f, nil
	} else {
		err := fmt.Errorf("The download file not exist, fileKey: %d", fileKey)
		return nil, err
	}

}

// relative path = path + filename + "_" + fileversion
func (fs *DLFileSet) Add(file *BasicFileInfo) error {
	fs.Lock()
	defer fs.Unlock()

	fileName := strings.TrimSpace(file.Name)
	fileVersion := strings.TrimSpace(file.Version)

	// check if exsit
	fileKey := fileName + "_" + fileVersion
	if _, ok := fs.Files[fileKey]; ok {
		return nil
	}

	path := strings.TrimSpace(file.Path)
	if len(path) != 0 && path[len(path)-1] != '/' {
		path += "/"
	}

	relativePath := path + fileName + "_" + fileVersion
	f, err := NewDownloadFile(relativePath)
	if err != nil {
		log.Printf("New File failed, file:%v, error:%v\n", file, err)
		return err
	}

	// insert new one
	fs.Files[fileKey] = f
	return nil
}
