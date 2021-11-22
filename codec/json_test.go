package codec

import (
	"encoding/json"
	"fmt"
	jsonIter "github.com/json-iterator/go"
	"github.com/lensesio/tableprinter"
	"os"
	"testing"
)

type Student struct {
	ID       uint     `json:"id"`
	Age      uint8    `json:"age"`
	Gender   uint8    `json:"gender"`
	Name     []byte   `json:"name"`
	Location Location `json:"location"`
}

type Location struct {
	Country  string
	Province string
	City     string
	District string
}

func TestEncodeJson(t *testing.T) {
	var s = Student{
		ID:     1,
		Age:    27,
		Gender: 1,
		Name:   []byte("yuchanns"),
		Location: Location{
			Country:  "China",
			Province: "Guangdong",
			City:     "Shenzhen",
			District: "Nanshan",
		},
	}
	rst, err := EncodeJson(s)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Printf("%s", rst)

	v := Student{}
	DecodeJson(rst, &v)
	fmt.Printf("\n%#v\n", v)
	fmt.Printf("\n%#v\n", v.Name)
	fmt.Printf("\n%s\n", v.Name)
}

func TestDecodeJson(t *testing.T) {
	s := `{
 "id": 1,
 "age": 27,
 "gender": 1,
 "name": "eXVjaGFubnM=",
 "location": {
  "Country": "China",
  "Province": "Guangdong",
  "City": "Shenzhen",
  "District": "Nanshan"
 }
}`
	v := Student{}
	DecodeJson([]byte(s), &v)
	fmt.Printf("\n%#v\n", v)

	b, _ := json.MarshalIndent(v, "", "    ")
	tableprinter.PrintJSON(os.Stdout, b)

}

func TestGetJson(t *testing.T) {
	s := `{
 "id": 1,
 "age": 27.345,
 "gender": 1,
 "name": "eXVjaGFubnM=",
 "location": {
  "Country": ["China", "US"],
  "Province": "Guangdong",
  "City": "Shenzhen",
  "District": "Nanshan"
 }
}`
	v1 := jsonIter.Get([]byte(s), "id").ValueType()
	fmt.Printf("\n%#v\n", v1)
	v2 := jsonIter.Get([]byte(s), "age")
	fmt.Printf("\n%#v\n", v2)
	v3 := jsonIter.Get([]byte(s), "location")
	fmt.Printf("\n%#v\n", v3)
	v4 := jsonIter.Get([]byte(s), "location", "Country")
	fmt.Printf("\n%#v\n", v4)
	v5 := jsonIter.Get([]byte(s), "location", "Country", 0)
	fmt.Printf("\n%#v\n", v5)

	v6 := jsonIter.Get([]byte(s), "name")
	fmt.Printf("\n%s\n", v6.ToString())
}
