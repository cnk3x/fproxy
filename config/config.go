package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	decoders   map[string]Decoder
	defaultExt = "json"
)

type Decoder interface {
	Decode([]byte, interface{}) error
}

func Register(decoder Decoder, ext ...string) {
	if decoders == nil {
		decoders = make(map[string]Decoder)
	}
	for _, i := range ext {
		decoders[i] = decoder
	}
}

func Unmarshal(in string, out interface{}) error {
	var (
		ext    = strings.TrimPrefix(filepath.Ext(in), ".")
		v, err = ioutil.ReadFile(in)
	)

	if ext == "" {
		ext = defaultExt
	}

	if err != nil {
		return err
	}

	decoder, find := decoders[ext]
	if find {
		return decoder.Decode(v, out)
	}
	return fmt.Errorf("config decoder not found for %s", ext)
}
