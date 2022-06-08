package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"
	"strconv"
	"unicode/utf8"

	"gorm.io/gorm"
)

type ArticlesController struct {
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// URL 파라미터
	id := route.GetRouteVariable("id", r)

	// 문장 데이터 획득
	article, err := article.Get(id)

	// 에러 발생 시
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 404 에러
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 error")
		} else {
			// 서버 에러
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 Internal Server Error")
		}
	} else {
		// 데이터 조회 성공
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL": route.Name2URL,
				"Int64ToString": types.Int64ToString,
			}).
			ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		tmpl.Execute(w, article)
	}
}

func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 Error")
	} else {
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, articles)
	}
}

type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
}

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {

	storeURL := route.Name2URL("articles.store")
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

	tmpl.Execute(w, data)
}

func validateArticleFormData(title string, body string) map[string]string {
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

	return errors
}

func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "Insert Success, ID is: "+strconv.FormatUint(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Create Fail")
		}
	} else {

		storeURL := route.Name2URL("articles.store")

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

		logger.LogError(err)

		tmpl.Execute(w, data)
	}
}
