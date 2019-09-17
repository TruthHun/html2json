package main

import (
	"compress/gzip"
	"errors"
	"github.com/TruthHun/html2json/html2json"
	"github.com/gin-contrib/cors"
	ginzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Error string      `json:"error,omitempty"`
	IsOK  bool        `json:"is_ok"`
	Data  interface{} `json:"data,omitempty"`
}

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
	app.GET("/html2json", html2JSON)  // params: url, expire
	app.POST("/html2json", html2JSON) // params: html
	err := app.Run(":8888")
	if err != nil {
		panic(err)
	}
}

func html2JSON(ctx *gin.Context) {
	var err error
	resp := Response{IsOK: true}
	switch ctx.Request.Method {
	case http.MethodPost:
		htmlStr := ctx.DefaultPostForm("html", "")
		if htmlStr == "" {
			err = errors.New("html is empty")
		} else {
			resp.Data, err = html2json.Parse(htmlStr)
		}
	case http.MethodGet:
		urlStr := ctx.DefaultQuery("url", "")
		exp, _ := strconv.Atoi(ctx.DefaultQuery("expire", "10"))
		if urlStr == "" {
			err = errors.New("url is empty")
		} else {
			resp.Data, err = html2json.ParseByURL(urlStr, exp)
		}
	default:
		err = errors.New("request method is not allow")
	}
	resp.IsOK = err == nil
	if err != nil {
		resp.Error = err.Error()
	}
	ctx.JSON(http.StatusOK, resp)
}
