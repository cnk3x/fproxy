package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strings"
)

type App []*Proxy

type Proxy struct {
	Name    string
	Prefix  string
	Target  string
	handler http.Handler
}

func (sa *App) Handle(name, prefix, targetUrl string) {
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

	*sa = append(*sa, &Proxy{Name: name, Prefix: prefix, Target: targetUrl, handler: handler})
	sort.Sort(sa)
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

	path := r.URL.Path
	for _, sr := range *sa {
		if strings.HasPrefix(strings.ToLower(path), sr.Prefix) {
			sr.handler.ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, path+" not found", http.StatusNotFound)
}

func (sa *App) Len() int {
	return len(*sa)
}

func (sa *App) Less(i, j int) bool {
	var (
		l1 = len(strings.Split((*sa)[i].Prefix, "/"))
		l2 = len(strings.Split((*sa)[j].Prefix, "/"))
	)
	if l1 == l2 {
		l1 = len((*sa)[j].Prefix)
		l2 = len((*sa)[i].Prefix)
	}
	return l1 > l2
}

func (sa *App) Swap(i, j int) {
	(*sa)[i], (*sa)[j] = (*sa)[j], (*sa)[i]
}
