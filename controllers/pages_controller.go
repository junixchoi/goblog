package controllers

import (
	"fmt"
	"net/http"
)

type PagesController struct {
}

func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, Goblog</h1>")
}

func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "about page")
}

func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>잘못된 경로입니다</h1>")
}
