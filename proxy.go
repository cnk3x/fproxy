package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	HeaderXForwardedFor   = "X-Forwarded-For"
	HeaderXForwardedProto = "X-Forwarded-Proto"
	HeaderXRealIP         = "X-Real-IP"
)

func Ps(target *url.URL) http.Handler{
	return &proxy{target:target}
}

type proxy struct {
	target *url.URL
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	var realIp string
	if ip := r.Header.Get(HeaderXForwardedFor); ip != "" {
		realIp = strings.Split(ip, ", ")[0]
	} else if ip := r.Header.Get(HeaderXRealIP); ip != "" {
		realIp = ip
	} else {
		realIp, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	// Fix header
	if w.Header().Get(HeaderXRealIP) == "" {
		w.Header().Set(HeaderXRealIP, realIp)
	}
	if w.Header().Get(HeaderXForwardedProto) == "" {
		w.Header().Set(HeaderXForwardedProto, r.URL.Scheme)
	}
	upgrade := r.Header.Get("upgrade")
	isWs := upgrade == "websocket" || upgrade == "Websocket"
	if isWs && w.Header().Get(HeaderXForwardedFor) == "" { // For HTTP, it is automatically set by Go HTTP reverse proxy.
		w.Header().Set(HeaderXForwardedFor, realIp)
	}

	if isWs {
		proxyRaw(w, r, p.target)
	} else {
		httputil.NewSingleHostReverseProxy(p.target).ServeHTTP(w, r)
	}
}

func proxyRaw(w http.ResponseWriter, r *http.Request, target *url.URL) {
	in, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, fmt.Sprintf("proxy raw, hijack error=%v, url=%s", target, err), http.StatusInternalServerError)
		return
	}
	defer in.Close()

	out, err := net.Dial("tcp", target.Host)
	if err != nil {
		http.Error(w, fmt.Sprintf("proxy raw, dial error=%v, url=%s", target, err), http.StatusBadGateway)
		return
	}
	defer out.Close()

	// Write header
	err = r.Write(out)
	if err != nil {
		http.Error(w, fmt.Sprintf("proxy raw, request header copy error=%v, url=%s", target, err), http.StatusBadGateway)
		return
	}

	errCh := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader) {
		_, err = io.Copy(dst, src)
		errCh <- err
	}

	go cp(out, in)
	go cp(in, out)
	err = <-errCh
	if err != nil && err != io.EOF {
		http.Error(w, fmt.Sprintf("proxy raw, copy body error=%v, url=%s", target, err), http.StatusBadGateway)
	}
}
