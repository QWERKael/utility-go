package codec

import (
	"fmt"
	"gopkg.in/ini.v1"
)

func DecodeIniAllSection(path string) (map[string]map[string]string, error) {
	f, err := ini.LoadSources(ini.LoadOptions{
		AllowBooleanKeys: true,
	}, path)
	if err != nil {
		return nil, err
	}

	f.Sections()

	m := make(map[string]map[string]string)

	sections := f.Sections()
	for _, section := range sections {
		m[section.Name()] = section.KeysHash()
	}
	return m, nil
}

func DecodeIniSection(b []byte, sectionName string) (map[string]string, error) {
	f, err := ini.LoadSources(ini.LoadOptions{
		AllowBooleanKeys: true,
	}, b)

	fmt.Printf("\nf值：%#v\n", f)
	if err != nil {
		return nil, err
	}

	var section *ini.Section
	section, err = f.GetSection(sectionName)
	if err != nil {
		return nil, err
	}
	return section.KeysHash(), nil
}
