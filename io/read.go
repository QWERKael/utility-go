package io

import (
	"errors"
	"github.com/QWERKael/utility-go/path"
	"io/ioutil"
)

func ReadFile(filePath string) ([]byte, error) {
	if fileType, err := path.CheckPath(filePath); err == nil {
		if fileType == path.File {
			b, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, errors.New("打开文件失败: " + err.Error())
			}
			return b, nil
		} else {
			return nil, errors.New("指定路径不是文件: " + filePath)
		}
	} else {
		return nil, errors.New("路径检测失败: " + err.Error())
	}
}
