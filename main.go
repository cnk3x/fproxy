package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

const (
	NAME        = "fproxy"
	DESCRIPTION = "前端代理工具"
	VERSION     = "1.0.0"
	BUILD       = "master"
)

func main() {
	fmt.Fprintf(os.Stdin, "%s - %s\n\n版本:\n  %s\n编译时间:\n  %s\n", NAME, DESCRIPTION, VERSION, BUILD)

	var (
		listenAt string
		proxy    []string
	)

	var cfg string
	flag := pflag.NewFlagSet(NAME, pflag.ContinueOnError)
	flag.StringVarP(&cfg, "config", "c", "fproxy.toml", "配置文件路径")
	flag.StringVarP(&listenAt, "listen", "l", ":3000", "发布端口")
	flag.StringSliceVarP(&proxy, "proxy", "p", []string{}, "代理链")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "参数:\n")
		flag.PrintDefaults()
		fmt.Println()
		os.Exit(1)
	}
	fmt.Println()
	flag.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
	flag.Parse(os.Args[1:])

	log.SetFlags(log.LUTC | log.Lshortfile)

	var c AppConfig

	if cfg != "" {
		v, err := ioutil.ReadFile(cfg)
		if err == nil {
			err = yaml.Unmarshal(v, &c)
		}
		if err != nil {
			if cfg != "fproxy.yml" {
				log.Fatal(err)
			}
		}
	}

	if c.App.ListenAt == "" {
		c.App.ListenAt = listenAt
	}

	if c.App.ListenAt == "" {
		c.App.ListenAt = ":3000"
	}

	for _, link := range proxy {
		if pc := parseProxyLink(link); pc != nil {
			c.Proxy = append(c.Proxy, pc)
		}
	}

	if len(c.Proxy) == 0 {
		flag.Usage()
	}

	fmt.Println(proxy)

	Run(c)
}

func Run(c AppConfig) {
	sApp := new(App)

	for _, s := range c.Proxy {
		sApp.Handle(s.Name, s.Prefix, s.Target)
	}

	for _, app := range *sApp {
		fmt.Printf("name=%s, prefix=%s, target=%s\n", app.Name, app.Prefix, app.Target)
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

func parseProxyLink(link string) *ProxyConfig {
	if link != "" {
		var (
			name   string
			prefix string
			target string
		)
		s := strings.Split(link, ",")

		for i := range s {
			if strings.Contains(name, ":") || strings.HasPrefix(name, ".") {
				target = s[i]
			} else if strings.Contains(name, "/") {
				prefix = s[i]
			} else {
				name = s[i]
			}
		}

		if prefix == "" {
			prefix = "/"
		}

		if name == "" {
			name = strings.Replace(prefix, "/", "_", -1)
		}

		return &ProxyConfig{Name: name, Prefix: prefix, Target: target}
	}
	return nil
}
