package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	NAME    = "前端代理"
	VERSION = "1.0.0"
)

func main() {
	log.SetFlags(log.Ltime)

	fmt.Printf("%s - %s", NAME, VERSION)
	fmt.Println()

	var cfg string
	flag.StringVar(&cfg, "c", "fproxy.toml", "配置文件路径")
	flag.Parse()

	var c Config
	_, err := toml.DecodeFile(cfg, &c)
	if err != nil {
		if cfg == "fproxy.toml" {
			flag.Usage()
			fmt.Println()
		}
		log.Fatal(err)
	}

	Run(c)
}

func Run(c Config) {
	sApp := new(App)

	for _, s := range c.Proxy {
		sApp.Handle(s.Name, s.Prefix, s.Target)
	}

	for _, app := range *sApp {
		log.Printf("name=%s, prefix=%s, target=%s\n", app.Name, app.Prefix, app.Target)
	}

	log.Printf("%-6s%s", "发布地址 ", c.App.ListenAt)
	s := &http.Server{
		Handler:      sApp,
		Addr:         c.App.ListenAt,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
