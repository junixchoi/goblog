package article

import (
	"goblog/app/models"
	"goblog/pkg/route"
	"strconv"
)

type Article struct {
	models.BaseModel
	Title string
	Body  string
}

func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(a.ID, 10))
}
