package yaml

import (
	"fproxy/config"

	"github.com/BurntSushi/toml"
)

func init() {
	config.Register(&Decoder{}, "toml", "tml")
}

type Decoder struct {
}

func (d *Decoder) Decode(v []byte, out interface{}) error {
	return toml.Unmarshal(v, out)
}
