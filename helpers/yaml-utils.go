package helpers

import (
	"gopkg.in/yaml.v2"
)

func YamlConfigParser(in map[string] interface{}, output interface{}) error {
	out, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(out, output)
	if err != nil {
		return err
	}
	return nil
}

func StructToMap(in interface{}) (map[string]interface{}, error) {
	output_yaml, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}
	var output map[string]interface{}
	err = yaml.Unmarshal(output_yaml, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}