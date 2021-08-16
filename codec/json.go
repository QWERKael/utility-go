package codec

import (
	jsonIter "github.com/json-iterator/go"
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