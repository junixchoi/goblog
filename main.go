package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var db *sql.DB

func initDB() {
	var err error
	config := mysql.Config{
		User: "root",
		Passwd: "",
		Addr: "127.0.0.1:3306"
		Net: "tcp"
		DBName: "goblog"
		AllowNativePasswords: true,
	}

	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, Goblog</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "about page")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>잘못된 경로입니다</h1>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "Article ID: "+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "articles list")
}

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errors := make(map[string]string)

	// title validator
	if title == "" {
		errors["title"] = "title is empty"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "title's length is > 3, < 40"
	}

	// body validator
	if body == "" {
		errors["body"] = "body is empty"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "body's length must >= 10"
	}

	if len(errors) == 0 {
		fmt.Fprintf(w, "validated! <br>")
		fmt.Fprintf(w, "title is: %v <br>", title)
		fmt.Fprintf(w, "title's length is: %v <br>", utf8.RuneCountInString(title))
		fmt.Fprintf(w, "body is: %v <br>", body)
		fmt.Fprintf(w, "body's length is: %v <br>", utf8.RuneCountInString(body))
	} else {
		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 메인 페이지 이외 모든 경로에 대한 요청에서 슬래시를 제거한다.
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()

	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func main() {
	initDB()
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// 404 페이지
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 미들웨어 사용
	router.Use(forceHTMLMiddleware)

	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)
	articleURL, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL: ", articleURL)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
