package main

import (
	"net/http"

	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/pkg/logger"
)

func main() {
	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	err := http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)
}
