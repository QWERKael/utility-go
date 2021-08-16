package io

import (
	"fmt"
	"path"
	"testing"
)

func TestReadFile(t *testing.T) {
	filePath := path.Join("..", "testDir", "testFile")
	if b, err := ReadFile(filePath); err != nil {
		t.Errorf("read file [%s] failed: %s", filePath, err.Error())
	} else {
		fmt.Printf("get content from [%s]:\n%s\n", filePath, b)
		if fmt.Sprintf("%s", b) != "123456" {
			t.Errorf("content is wrong")
		}
	}
}
