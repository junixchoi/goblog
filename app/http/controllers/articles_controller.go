package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"

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