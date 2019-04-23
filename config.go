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

/*
app:
  listen_at: :3000

proxy:
  - name: front
    prefix: /
    target: ./www

  - name: backend
    prefix: /api
    target: http://119.23.215.60:48080/
*/
type AppConfig struct {
	App struct {
		ListenAt string `json:"listen" xml:"listen" toml:"listen" yaml:"listen"`
	} `json:"app" xml:"listen" toml:"app" xml:"app" yaml:"app"`

	Proxy [] *ProxyConfig `json:"proxy" toml:"proxy" xml:"proxy" yaml:"proxy"`
}

type ProxyConfig struct {
	Name   string `json:"name" xml:"listen" toml:"name" xml:"name" yaml:"name"`
	Prefix string `json:"prefix" xml:"listen" toml:"prefix" xml:"prefix" yaml:"prefix"`
	Target string `json:"target" xml:"listen" toml:"target" xml:"target" yaml:"target"`
}
