package yaml

import (
	"encoding/xml"
	"fproxy/config"
)

func init() {
	config.Register(&Decoder{}, "xml")
}

type Decoder struct {
}

func (d *Decoder) Decode(v []byte, out interface{}) error {
	return xml.Unmarshal(v, out)
}
