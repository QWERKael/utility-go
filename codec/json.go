package codec

import (
	"bufio"
	path2 "github.com/QWERKael/utility-go/path"
	jsonIter "github.com/json-iterator/go"
	"io/ioutil"
)

func EncodeJson(v interface{}) ([]byte, error) {
	//var json = jsonIter.ConfigCompatibleWithStandardLibrary
	//return jsonIter.Marshal(&v)
	return jsonIter.MarshalIndent(&v, "", " ")
}

func DecodeJson(b []byte, v interface{}) error {
	//var json = jsonIter.ConfigCompatibleWithStandardLibrary
	//return jsonIter.Marshal(&v)
	return jsonIter.Unmarshal(b, v)
}

func DecodeJsonFromFile(path string, v interface{}) error {
	file, err := path2.OpenFileIfExist(path)
	if err != nil {
		return err
	}
	var b []byte
	b, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	return DecodeJson(b, v)
}

func EncodeJsonToFile(path string, v interface{}) error {
	file, err := path2.CreateOrOpenFileForOverWrite(path)
	if err != nil {
		return err
	}
	bufWriter := bufio.NewWriter(file)
	defer func() {
		bufWriter.Flush()
		file.Close()
	}()
	var b []byte
	b, err = EncodeJson(v)
	if err != nil {
		return err
	}
	_, err = bufWriter.Write(b)
	if err != nil {
		return err
	}
	return nil
}
