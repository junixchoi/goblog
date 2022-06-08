package article

import (
	"goblog/pkg/route"
	"goblog/pkg/types"
)

type Article struct {
	ID    uint64
	Title string
	Body  string
}

func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", types.Uint64ToString(a.ID))
}
