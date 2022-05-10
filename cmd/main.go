package main

import (
	"go-currency/http"
	"go-currency/service"

	"github.com/didip/tollbooth"

	"github.com/gin-gonic/gin"
)

func main() {
	//Default返回一个默认的路由引擎
	r := gin.Default()
	limiter := tollbooth.NewLimiter(1000, nil)
	r.Use(http.Limiter(limiter))
	service.Init()
	http.Init(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
