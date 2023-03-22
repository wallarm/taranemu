package tarantella

import (
	_ "embed"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

var (
	dummyInstanceID = uuid.New().String()
	schemaVersion   = 0x56 // has no special meaning for the stub
)

//go:embed dummy-281_vspace.yaml
var dummySpacesYaml []byte
var dummySpaces = parseYamlArray(dummySpacesYaml)

//go:embed dummy-289_vindex.yaml
var dummyIndexesYaml []byte
var dummyIndexes = parseYamlArray(dummyIndexesYaml)

//go:embed dummy-execute.1.yaml
var dummyExecute1Yaml []byte
var dummyExecute1 = parseYamlMap(dummyExecute1Yaml)

// parseYamlArray just parses YAML content into array
func parseYamlArray(content []byte) []any {
	var m []any
	err := yaml.Unmarshal(content, &m)
	if err != nil {
		panic(err)
	}
	return m
}

// parseYamlArray just parses YAML content into map
func parseYamlMap(content []byte) map[any]any {
	m := make(map[any]any)
	err := yaml.Unmarshal(content, &m)
	if err != nil {
		panic(err)
	}
	return m
}
