package codec

import (
	"github.com/BurntSushi/toml"
	path2 "github.com/QWERKael/utility-go/path"
	"io/ioutil"
)

func DecodeTomlFromFile(path string, v interface{}) error {
	if file, err := path2.OpenFileIfExist(path); err != nil {
		return err
	} else if b, err := ioutil.ReadAll(file); err != nil {
		return err
	} else {
		return toml.Unmarshal(b, v)
	}
}
