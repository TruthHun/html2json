package main

import (
	"compress/gzip"
	"github.com/gin-contrib/cors"
	ginzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	app := gin.New()

	// 设置跨域和gzip
	app.Use(ginzip.Gzip(gzip.BestCompression), cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}), gin.Recovery())
	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"pong": "hello html2json!"})
	})
	app.GET("/html2json")  // params: url, expire
	app.POST("/html2json") // params: html
	err := app.Run(":8888")
	if err != nil {
		panic(err)
	}
}

func Html2JSON(ctx *gin.Context) {

}
