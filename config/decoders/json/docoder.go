package yaml

import (
	"encoding/json"
	"fproxy/config"
)

func init() {
	config.Register(&Decoder{}, "json")
}

type Decoder struct {
}

func (d *Decoder) Decode(v []byte, out interface{}) error {
	txt := []rune(string(v))
	l := len(txt)
	r := make([]rune, l)
	j := 0
	closed := true
	comment := false
	for i, c := range txt {
		switch c {
		case '\n':
			closed = true
			comment = false
		case '"':
			if comment {
				continue
			}
			if i > 0 && txt[i-1] != '\\' {
				closed = !closed
			}
		case '/':
			if comment {
				continue
			}
			if closed && i == 0 && i < l-1 && txt[i+1] == '/' {
				comment = true
				continue
			}
		}

		r[j] = c
		j++
	}
	return json.Unmarshal([]byte(string(r)), out)
}
