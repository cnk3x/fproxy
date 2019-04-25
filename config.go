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

	Proxy map[string]string `json:"proxy" toml:"proxy" xml:"proxy" yaml:"proxy"`
}
