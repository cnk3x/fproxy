package yaml

import (
	"encoding/json"
	"fproxy/config"
)

func init() {
	config.Register(&Decoder{}, "yaml", "yml")
}

type Decoder struct {
}

func (d *Decoder) Decode(v []byte, out interface{}) error {
	txt := []rune(string(v))
	r := make([]rune, len(txt))
	j := 0
	closed := true
	comment := false
	for i, c := range txt {
		switch c {
		case '\n':
			closed = true
			comment = false
			r[i] = c
			i++
		case '"':
			if comment {
				continue
			}
			if txt[i-1] == '\\' {
				
			}
			closed = !closed
		}

		if closed {

		}
	}
	return json.Unmarshal(v, out)
}
