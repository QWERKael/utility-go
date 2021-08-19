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

//
// CheckPath
//  @Description: 检查文件的路径
//  @param path 传入路径
//  @return FileType 返回文件的类型
//  @return error 当路径为空，或者获取路径信息失败时，报错
//
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

func CheckPathAs(path string, checkFt FileType) (bool, error) {
	if ft, err := CheckPath(path); err != nil {
		return false, err
	} else {
		if ft == checkFt {
			return true, nil
		} else {
			return false, nil
		}
	}
}

//
// CheckAndCreateDir
//  @Description: 检查并创建目录
//  @param path 传入路径
//  @return error 无法识别path或者path为文件时，报错
//
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

//
// OpenFileIfExist
//  @Description: 如果文件存在，则打开，否则报错
//  @param path
//  @return *os.File
//  @return error
//
func OpenFileIfExist(path string) (*os.File, error) {
	if ok, err := CheckPathAs(path, File); err != nil {
		return nil, err
	} else if ok {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return file, nil
	} else {
		return nil, errors.New("传入路径不是文件类型")
	}
}

//
// CreateFileIfNotExist
//  @Description: 如果文件不存在，则创建，否则报错
//  @param path
//  @return *os.File
//  @return error
//
func CreateFileIfNotExist(path string) (*os.File, error) {
	if ok, err := CheckPathAs(path, NotExist); err != nil {
		return nil, err
	} else if ok {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		return file, nil
	} else {
		return nil, errors.New("传入路径已存在")
	}
}

//
// CheckAndOpenFile
//  @Description: 打开一个文件，如果文件存在就直接打开，如果不存在就创建然后打开
//  @param path
//  @return *os.File
//  @return error
//
func CheckAndOpenFile(path string) (*os.File, error) {
	ft, err := CheckPath(path)
	if err != nil {
		return nil, err
	}
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

//
// SumMd5FromFile
//  @Description: 校验文件的md5值，返回字符串
//  @param fileName
//  @return string
//  @return error
//
func SumMd5FromFile(fileName string) (string, error) {
	file, err := OpenFileIfExist(fileName)
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
