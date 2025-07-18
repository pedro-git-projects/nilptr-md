package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pedro-git-projects/nilptr-md/httpext"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed assets/hello.md
var helloMD []byte

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	httpext.ContentType.Add(w.Header(), httpext.TextHTML)
	w.WriteHeader(http.StatusOK)

	md := goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)

	if err := md.Convert(helloMD, w); err != nil {
		http.Error(w, "Failed to render markdown", http.StatusInternalServerError)
		return
	}
}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", HelloWorldHandler)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        m,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Starting server on localhost:8080...")
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
