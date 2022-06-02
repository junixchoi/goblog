package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var db *sql.DB

type Article struct {
	Title, Body string
	ID          int64
}

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// generate article link
func (a Article) Link() string {
	showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
	if err != nil {
		checkError(err)
		return ""
	}
	return showURL.String()
}

func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(a.ID, 10))

	if err != nil {
		return 0, err
	}

	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

// 라우트 명칭으로 URL 획득
func RouteName2URL(routeName string, pairs ...string) string {
	url, err := router.Get(routeName).URL(pairs...)
	if err != nil {
		checkError(err)
		return ""
	}

	return url.String()
}

// int64 를 스트링으로 전환
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
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
	// URL 파라미터
	id := getRouteVariable("id", r)

	// 문장 데이터 획득
	article, err := getArticleByID(id)

	// 에러 발생 시
	if err != nil {
		if err == sql.ErrNoRows {
			// 404 에러
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 error")
		} else {
			// 서버 에러
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 Internal Server Error")
		}
	} else {
		// 데이터 조회 성공
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL": RouteName2URL,
				"Int64ToString": Int64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)

		err = tmpl.Execute(w, article)
		checkError(err)
	}
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	// 리스트 조회
	rows, err := db.Query("SELECT * FROM articles")
	checkError(err)
	defer rows.Close()

	var articles []Article
	// loop
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		checkError(err)
		articles = append(articles, article)
	}

	err = rows.Err()
	checkError(err)

	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	checkError(err)

	err = tmpl.Execute(w, articles)
	checkError(err)
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
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "insert success, ID is : "+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 error")
		}
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

func saveArticleToDB(title string, body string) (int64, error) {
	// 변수 초기화
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	// 1. prepare
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
	if err != nil {
		return 0, err
	}

	// 2. defer
	defer stmt.Close()

	// 3. execute
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	// 4. return ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
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

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	// URL 파라미터
	id := getRouteVariable("id", r)

	// 문장 데이터 획득
	article, err := getArticleByID(id)

	// 에러 발생 시
	if err != nil {
		if err == sql.ErrNoRows {
			// 404 에러
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 error")
		} else {
			// 서버 에러
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 Internal Server Error")
		}
	} else {
		// 데이터 조회 성공, 폼 표시
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		checkError(err)

		err = tmpl.Execute(w, data)
		checkError(err)
	}
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// URL 파라미터
	id := getRouteVariable("id", r)

	// 문장 데이터 획득
	_, err := getArticleByID(id)

	// 에러 발생 시
	if err != nil {
		if err == sql.ErrNoRows {
			// 404 에러
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 error")
		} else {
			// 서버 에러
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 Internal Server Error")
		}
	} else {
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := make(map[string]string)

		// 타이틀 검증
		if title == "" {
			errors["title"] = "title is empty."
		} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
			errors["title"] = "title length must be 3 ~40"
		}

		// 내용 검증
		if body == "" {
			errors["body"] = "body is empty."
		} else if utf8.RuneCountInString(body) < 10 {
			errors["body"] = "body length must be > 10"
		}

		if len(errors) == 0 {
			// 데이터 업데이트
			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
			rs, err := db.Exec(query, title, body, id)

			if err != nil {
				checkError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 Internal Server Error")
			}

			// 업데이트 성공 이후 상세 페이지로 이동
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "수정한 내용이 없습니다.")
			}
		} else {
			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			checkError(err)

			err = tmpl.Execute(w, data)
			checkError(err)
		}
	}
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// URL 파라미터
	id := getRouteVariable("id", r)

	// 문장 데이터 획득
	article, err := getArticleByID(id)

	// 에러 발생 시
	if err != nil {
		if err == sql.ErrNoRows {
			// 404 에러
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 error")
		} else {
			// 서버 에러
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 Internal Server Error")
		}
	} else {
		// 삭제
		rowsAffected, err := article.Delete()

		if err != nil {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 Internal Server Error")
		} else {
			if rowsAffected > 0 {
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 error")
			}
		}
	}
}

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `

	_, err := db.Exec(createArticlesSQL)
	checkError(err)
}

func main() {
	initDB()
	createTables()
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

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
