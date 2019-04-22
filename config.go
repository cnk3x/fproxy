package main
/*
[app]
listen  =":3000"

[[proxy]]
name    ="backend"
prefix  ="/api"
target  ="http://119.23.215.60:48080/"

[[proxy]]
name    ="front"
prefix  ="/"
target  ="www"
*/
type Config struct {
	App struct {
		ListenAt string `toml:"listen" json:"listen"`
	} `json:"app" toml:"app"`

	Proxy []struct {
		Name   string `json:"name" toml:"name"`
		Prefix string `json:"prefix" toml:"prefix"`
		Target string `json:"target" toml:"target"`
	} `json:"proxy" toml:"proxy"`
}
