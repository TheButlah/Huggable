package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

//StaticContent is a `http.Handler` that serves static files to the client
type StaticContent struct {
}

//NewStaticContent constructs a `StaticContent` handler
func NewStaticContent(root string) StaticContent {
	return StaticContent{}
}

func (sc StaticContent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	log.Println("[StaticContent] serving", r.URL.String())
	fmt.Fprintf(w, "Hello world! Its currently %s", t)
}

//Compile-time check that `StaticContent` implements `http.Handler`
var _ http.Handler = (*StaticContent)(nil)
