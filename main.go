package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, Goblog</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "about page")
	} else {
		fmt.Fprint(w, "<h1>잘못된 경로입니다</h1>")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
