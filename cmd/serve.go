// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	ginzip "github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/truthhun/html2json/html2json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动HTTP服务",
	Long:  `以HTTP接口的形式提供HTML转JSON的服务`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			tags []string
			port int
			err  error
			b    []byte
		)
		if port, err = strconv.Atoi(cmd.Flag("port").Value.String()); err != nil {
			fmt.Println(err.Error())
			fmt.Println("使用默认端口: ", 8888)
			port = 8888
		}

		tagsFile := cmd.Flag("tags").Value.String()
		if tagsFile != "" {
			if b, err = ioutil.ReadFile(tagsFile); err != nil {
				fmt.Println(err.Error())
				fmt.Println("使用默认HTML标签")
			} else {
				if err = json.Unmarshal(b, &tags); err != nil {
					fmt.Println(err.Error())
					fmt.Println("使用默认HTML标签")
				}
			}
		}
		serve(port, tags...)
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	serveCmd.PersistentFlags().Int("port", 8888, "服务监听端口")
	serveCmd.PersistentFlags().String("tags", "", "自定义的可信任的HTML标签所在的json文件路径")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Response struct {
	Error string      `json:"error,omitempty"`
	IsOK  bool        `json:"is_ok"`
	Data  interface{} `json:"data,omitempty"`
}

var rt = html2json.NewDefault()

func serve(port int, tag ...string) {
	app := gin.New()

	if len(tag) > 0 {
		rt = html2json.New(tag)
	}

	// 设置跨域和gzip
	app.Use(ginzip.Gzip(gzip.BestCompression), cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}), gin.Recovery())

	app.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"pong": "hello html2json!"}) })
	app.GET("/html2json", html2JSON)  // params: url, expire
	app.POST("/html2json", html2JSON) // params: html

	fmt.Println("serve on port:", port)
	err := app.Run(fmt.Sprintf(":%v", port))
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
			resp.Data, err = rt.Parse(htmlStr)
		}
	case http.MethodGet:
		urlStr := ctx.DefaultQuery("url", "")
		exp, _ := strconv.Atoi(ctx.DefaultQuery("expire", "10"))
		if urlStr == "" {
			err = errors.New("url is empty")
		} else {
			resp.Data, err = rt.ParseByURL(urlStr, exp)
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
