package primitives

import (
	"bytes"
	"encoding/json"

	yamlv3 "gopkg.in/yaml.v3"

	"github.com/palantir/stacktrace"
)

const (
	empty = ""
	tab   = "\t"
)

func IndentedJSON(data interface{}) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
func IndentedYAML(data interface{}) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	enc := yamlv3.NewEncoder(buffer)
	defer enc.Close()
	enc.SetIndent(2)
	err := enc.Encode(data)
	if err != nil {
		err = stacktrace.Propagate(err, "could not encode data")
		return nil, err
	}
	return buffer, nil
}
