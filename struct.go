package gisk

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type Value struct {
	ValueType ValueType
	Value     interface{}
}

type RawMessage []byte

func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}
func (m RawMessage) MarshalYAML() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for YAMLRawMessage.
func (m *RawMessage) UnmarshalYAML(value *yaml.Node) error {
	if m == nil {
		return errors.New("yaml.RawMessage: UnmarshalYAML on nil pointer")
	}

	data, err := yaml.Marshal(value)
	if err != nil {
		return err
	}

	*m = append((*m)[0:0], data...)
	return nil
}
