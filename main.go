package main

import (
	"fmt"
	"net/http"
	"strings"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, Goblog</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>잘못된 경로입니다</h1>")
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "about page")
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/about", aboutHandler)

	router.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.SplitN(r.URL.Path, "/", 3)[2]
		fmt.Fprint(w, "Article ID: "+id)
	})

	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprint(w, "articles list")
		case "POST":
			fmt.Fprint(w, "new article")
		}
	})

	http.ListenAndServe(":3000", router)
}
