package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"sort"
	"strings"
)

func NewApp() *App {
	return &App{
		proxies:  make(map[string]*Proxy),
		prefixes: make([]string, 0),
	}
}

type App struct {
	proxies  map[string]*Proxy
	prefixes []string
}

type Proxy struct {
	Prefix  string
	Target  string
	handler http.Handler
}

func (sa *App) Handle(prefix, targetUrl string) {
	prefix = path.Clean("/" + prefix + "/")

	var handler http.Handler
	strip := strings.HasSuffix(targetUrl, "/")
	targetUrl = strings.TrimSuffix(targetUrl, "/")
	target, err := url.Parse(targetUrl)
	if err != nil {
		log.Fatal(err)
	}

	if target.IsAbs() && target.Scheme != "file" {
		handler = Ps(target)
	} else {
		handler = Fs(targetUrl)
	}

	if strip {
		handler = http.StripPrefix(prefix, handler)
	}

	p := &Proxy{Prefix: prefix, Target: targetUrl, handler: handler}

	if _, find := sa.proxies[prefix]; !find {
		sa.prefixes = append(sa.prefixes, prefix)
		sort.Sort(PathSorter(sa.prefixes))
	}

	sa.proxies[prefix] = p
}

func (sa *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.Proto + " " + r.Method + " " + r.Host + r.URL.String() + " " + r.RemoteAddr + " " + r.UserAgent())
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			stack := make([]byte, 4<<10)
			length := runtime.Stack(stack, false)
			log.Printf("[PANIC RECOVER] %v %s\n", err, stack[:length])
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	urlPath := r.URL.Path
	for _, prefix := range sa.prefixes {
		if strings.HasPrefix(strings.ToLower(urlPath), prefix) {
			proxy := sa.proxies[prefix]
			proxy.handler.ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, urlPath+" not found", http.StatusNotFound)
}

type PathSorter []string

func (sa PathSorter) Len() int {
	return len(sa)
}

func (sa PathSorter) Less(i, j int) bool {
	var (
		l1 = len(strings.Split(sa[i], "/"))
		l2 = len(strings.Split(sa[j], "/"))
	)
	if l1 == l2 {
		l1 = len(sa[j])
		l2 = len(sa[i])
	}
	return l1 < l2
}

func (sa PathSorter) Swap(i, j int) {
	(sa)[i], (sa)[j] = (sa)[j], (sa)[i]
}
