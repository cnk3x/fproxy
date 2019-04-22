package main

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Fs string

// ServeHTTP serves a request, attempting to reply with an embedded file.
func (fs Fs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/")

	p, err := url.PathUnescape(p)
	if err != nil {
		return
	}

	fs.File(w, r, p)
}

func (fs Fs) File(w http.ResponseWriter, r *http.Request, p string) {
	p = path.Clean("/" + p)

	name := filepath.Join(string(fs), p)
	fi, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			if p != "/index.html" {
				fs.File(w, r, "index.html")
			} else {
				http.NotFound(w, r)
			}
			return
		}
		return
	}

	if fi.IsDir() {
		fs.File(w, r, filepath.Join(p, "index.html"))
		return
	}

	f, err := os.Open(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
}
