package yaml

import (
	"fproxy/config"
	"gopkg.in/yaml.v2"
)

func init() {
	config.Register(&Decoder{}, "yaml", "yml")
}

type Decoder struct {
}

func (d *Decoder) Decode(v []byte, out interface{}) error {
	return yaml.Unmarshal(v, out)
}
