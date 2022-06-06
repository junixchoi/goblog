package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

// 라우트 명칭으로 URL 획득
func Name2URL(routeName string, pairs ...string) string {
	var route *mux.Router
	url, err := route.Get(routeName).URL(pairs...)
	if err != nil {
		// checkError(err)
		return ""
	}

	return url.String()
}

// URI 라우트 파라미터 획득
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
