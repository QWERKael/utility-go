package path

import (
	"fmt"
	"path"
	"testing"
)

func TestCheckPath(t *testing.T) {
	dir := path.Join("..", "testDir")
	file := path.Join(dir, "testFile")
	notExist := "notExist"
	var fileType FileType
	if fileType, _ = CheckPath(dir); fileType != Dir {
		t.Errorf("[%s] is not a dir", dir)
	}
	if fileType, _ = CheckPath(file); fileType != File {
		t.Errorf("[%s] is not a file", file)
	}
	if fileType, _ = CheckPath(notExist); fileType != NotExist {
		t.Errorf("[%s] is exist", notExist)
	}
	if fileType, err := CheckPath(""); err != nil {
		t.Errorf("path [%s] is wrong: %s", "", err.Error())
	} else {
		fmt.Printf("get file type: [%v]", fileType)
	}
}
