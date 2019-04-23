package yaml

import (
	"encoding/xml"
	"fproxy/config"
)

func init() {
	config.Register(&Decoder{}, "yaml", "yml")
}

type Decoder struct {
}

func (d *Decoder) Decode(v []byte, out interface{}) error {
	return xml.Unmarshal(v, out)
}
