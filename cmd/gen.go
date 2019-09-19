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
	"encoding/json"
	"fmt"
	"github.com/TruthHun/html2json/html2json"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成小程序 rich-text 组件支持的HTML标签到json文件中",
	Long: `
html2json gen --cate app		生成 uni-app 支持的HTML标签
html2json gen --cate alipay	生成支付宝小程序支持的HTML标签
html2json gen --cate weixin	生成微信小程序支持的HTML标签
html2json gen --cate baidu		生成百度小程序支持的HTML标签
html2json gen --cate qq		生成QQ小程序支持的HTML标签
html2json gen --cate toutiao	生成头条小程序支持的HTML标签
`,
	Run: func(cmd *cobra.Command, args []string) {
		cate := cmd.Flag("cate").Value.String()
		tags := html2json.GetTags(cate)
		file := fmt.Sprintf("%v.json", cate)
		b, err := json.Marshal(tags)
		if err != nil {
			panic(err)
		}
		if err = ioutil.WriteFile(file, b, os.ModePerm); err != nil {
			panic(err)
		}
		fmt.Printf("write to file : %v\n", file)
	},
}

func init() {
	RootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	genCmd.PersistentFlags().String("cate", "app", "小程序分类")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
