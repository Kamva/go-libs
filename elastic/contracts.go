package elastic

import (
	"encoding/json"
)

type Settings map[string]interface{}

type Map map[string]TypeMap

type TypeMap struct {
	Properties TypeMapProperties `json:"properties"`
}

type TypeMapProperties map[string]FieldSetting

type FieldSetting map[string]interface{}

type Index interface {
	GetMapping() Mapping
}

type Type interface {
	GetIndex() Index
	GetBody() interface{}
}

type Mapping struct {
	Settings Settings `json:"settings"`
	Mappings Map      `json:"mappings"`
}

func (m Mapping) Serialize() string {
	result, err := json.Marshal(m)

	if err != nil {
		throwException(err, "LIB_ERR")
	}

	return string(result)
}
