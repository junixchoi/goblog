package view

import (
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

func Render(w io.Writer, name string, data interface{}) {
	// set template directory
	viewDir := "resources/views/"

	// "articles.show" => "articles/show"
	name = strings.Replace(name, ".", "/", -1)

	// all template files slice
	files, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	// slice append target files
	newFiles := append(files, viewDir+name+".gohtml")

	// parse template files
	tmpl, err := template.New(name + ".gohtml").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(newFiles...)
	logger.LogError(err)

	// rendering template
	err = tmpl.ExecuteTemplate(w, "app", data)
	logger.LogError(err)
}
