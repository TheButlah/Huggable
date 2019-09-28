package handlers

import (
	"log"
	"net/http"
	"path/filepath"
)

//StaticContent is a `http.Handler` that serves static files to the client
type StaticContent struct {
	fs http.Handler
}

//Compile-time check that `StaticContent` implements `http.Handler`
var _ http.Handler = (*StaticContent)(nil)

//NewStaticContent constructs a `StaticContent` handler
func NewStaticContent(root string) StaticContent {
	root = filepath.FromSlash(root)
	return StaticContent{
		fs: http.FileServer(http.Dir(root)),
	}
}

func (sc StaticContent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[StaticContent] serving", r.URL.String())
	// Don't cache service worker
	if r.URL.Path == "/sw.js" {
		w.Header().Set("Cache-Control", "No-Cache")
		sc.fs.ServeHTTP(w, r)
		return
	}
	sc.fs.ServeHTTP(w, r)
}
