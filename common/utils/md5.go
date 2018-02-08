// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func Md5(raw string) string {
	h := md5.New()
	io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
}

// for file
func CreateMd5(f *os.File) ([]byte, error) {
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	md5_sum := h.Sum(nil)
	return md5_sum, nil
}

// check
func CheckSum(f *os.File, origin_sum []byte) (bool, error) {
	new_sum, err := CreateMd5(f)
	if err != nil {
		return false, err
	}

	if bytes.Equal(new_sum, origin_sum) {
		return true, nil

	} else {
		return false, nil
	}

}
