package path

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
)

type FileType int32

const (
	Unknown  FileType = 0
	NotExist FileType = 1
	File     FileType = 2
	Dir      FileType = 3
)

func CheckPath(path string) (FileType, error) {
	if path == "" {
		return Unknown, errors.New("path is nil")
	}
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NotExist, nil
		} else {
			return Unknown, err
		}
	}
	if fi.IsDir() {
		return Dir, nil
	} else {
		return File, nil
	}
}

func CheckAndCreateDir(path string) error {
	ft, err := CheckPath(path)
	if err != nil {
		return err
	}
	switch ft {
	case Dir:
		return nil
	case NotExist:
		err = os.MkdirAll(path, os.ModePerm)
		return err
	case File:
		return errors.New("目标路径是一个文件")
	}
	return errors.New("确认目录时遇到了一个未知的错误")
}

func CheckAndOpenFile(path string) (*os.File, error) {
	ft, err := CheckPath(path)
	var file *os.File
	switch ft {
	case NotExist:
		file, err = os.Create(path)
	case File:
		file, err = os.Open(path)
	case Dir:
		err = errors.New("目标路径是一个目录")
	}
	if err != nil {
		return nil, err
	}
	return file, nil
}

func SumMd5FromFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	m := md5.New()
	_, err = io.Copy(m, file)
	if err != nil {
		return "", err
	}
	Md5Str := fmt.Sprintf("%x", m.Sum(nil))
	return Md5Str, nil
}


