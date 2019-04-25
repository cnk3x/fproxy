package main

import (
	"fmt"
	"fproxy/config"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"

	_ "fproxy/config/decoders/json"
	_ "fproxy/config/decoders/toml"
	_ "fproxy/config/decoders/yaml"
)

const (
	Name        = "fproxy"
	Description = "前端代理工具"
	Version     = "1.0.0"
	BuildTime   = ""
)

func main() {
	fmt.Printf("%s - %s - %s\n\n", Name, Description, Version)

	var (
		listenAt string
		proxy    []string
		help     bool
		dbg      bool
	)

	var cfg string
	flag := pflag.NewFlagSet(Name, pflag.ContinueOnError)
	flag.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}

	flag.StringVarP(&cfg, "config", "c", "", "配置文件路径")
	flag.StringVarP(&listenAt, "listen", "l", ":80", "发布端口")
	flag.StringSliceVarP(&proxy, "proxy", "p", []string{}, "代理链")
	flag.BoolVarP(&dbg, "debug", "d", false, "调试模式")
	flag.BoolVarP(&help, "help", "h", false, "帮助信息")

	flag.Usage = func() {
		if len(BuildTime) > 0 {
			fmt.Printf("编译时间:\n  %s\n", BuildTime)
		}
		fmt.Printf("参数:\n")
		flag.PrintDefaults()
		fmt.Println()
	}

	_ = flag.Parse(os.Args[1:])

	log.SetFlags(log.Ldate | log.Ltime)

	if help {
		flag.Usage()
		os.Exit(0)
	}

	var c AppConfig

	if cfg != "" {
		err := config.Unmarshal(cfg, &c)
		if err != nil {
			log.Fatal(err)
		}
	}

	if c.App.ListenAt == "" {
		c.App.ListenAt = listenAt
	}

	if c.App.ListenAt == "" {
		c.App.ListenAt = ":80"
	}

	if c.Proxy == nil {
		c.Proxy = make(map[string]string)
	}

	for _, link := range proxy {
		i := strings.Index(link, "=")
		l := len(link) - 1

		switch i {
		case -1:
			c.Proxy["/"] = link
		case 0:
			c.Proxy["/"] = link[1:]
		case l:
			c.Proxy[link] = "./"
		default:
			c.Proxy[link[0:i]] = link[i+1:]
		}
	}

	if len(c.Proxy) == 0 {
		www := os.Getenv("www")
		api := os.Getenv("api")

		if api != "" {
			c.Proxy["api"] = api

			if www == "" {
				www = "www"
			}

			c.Proxy["/"] = www
		}
	}

	if len(c.Proxy) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if dbg {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	Run(c)
}

func Run(c AppConfig) {
	sApp := NewApp()

	for prefix, value := range c.Proxy {
		sApp.Handle(prefix, value)
	}

	for _, prefix := range sApp.prefixes {
		app := sApp.proxies[prefix]
		log.Printf("%s -> %s\n", app.Prefix, app.Target)
	}

	log.Printf("listen at %s", c.App.ListenAt)
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
