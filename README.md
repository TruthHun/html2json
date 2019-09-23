# html2json

使用Go语言开发的HTML和Markdown转JSON工具，将HTML/Markdown内容转换为符合各种小程序`rich-text`组件内容渲染所需格式的`JSON`

## 介绍

在开发 [BookStack](https://gitee.com/truthhun/BookStack) 的配套微信小程序 [BookChat](https://gitee.com/truthhun/BookChat) 
以及使用 `uni-app` 开发配套的手机APP应用 [BookChatApp](https://gitee.com/truthhun/BookChatApp)的过程中，
我尝试了很多种开源的小程序HTML解析渲染组件，但都不是很满意。

后来将HTML内容在后端解析，并在小程序段使用自带的`rich-text`组件进行渲染，性能、稳定性以及渲染效果都比较符合预期，尽管并没有第三方HTML渲染工具那样提供图片预览的功能。

目前已经在`BookStack` v2.1 版本中使用了。

## 特点

1. 支持markdown和HTML转JSON
1. 可作为 resetful 服务提供API接口进行访问
1. Go语言开发者，可作为包的形式进行引用

## 使用方式

更多使用帮助，使用如下方式查看

```
./html2json --help
```

### resetful 方式使用

#### 启动服务
```
./html2json serve --port 8888 --tags weixin-html-tags.json
```

- `--port` - [非必需参数]指定服务端口，默认为 8888
- `--tags` - [非必须参数]指定信任的HTML元素。json数组文件，里面存放各个支持的HTML标签。默认使用 uni-app 信任的HTML标签

各小程序支持的HTML标签

> - 微信小程序：https://developers.weixin.qq.com/miniprogram/dev/component/rich-text.html
> - 支付宝小程序：https://docs.alipay.com/mini/component/rich-text
> - 百度小程序：https://smartprogram.baidu.com/docs/develop/component/base/#rich-text-%E5%AF%8C%E6%96%87%E6%9C%AC/
> - 头条小程序：https://developer.toutiao.com/dev/miniapp/uEDMy4SMwIjLxAjM
> - QQ小程序：https://q.qq.com/wiki/develop/miniprogram/component/basic-content/rich-text.html
> - uni-app: https://uniapp.dcloud.io/component/rich-text?id=rich-text

`weixin-html-tags.json`文件示例：
```
["a", "abbr", "address", "article", "aside", "b", "bdi", "bdo", "big", "blockquote", "br", "caption", "center", "cite", "code", "col", "colgroup", "dd", "del", "div", "dl", "dt", "em", "fieldset", "font", "footer", "h1", "h2", "h3", "h4", "h5", "h6", "header", "hr", "i", "img", "ins", "label", "legend", "li", "mark", "nav", "ol", "p", "pre", "q", "rt", "ruby", "s", "section", "small", "span", "strong", "sub", "sup", "table", "tbody", "td", "tfoot", "th", "thead", "tr", "tt", "u", "ul"]
```

#### API接口

##### 解析来自url链接的HTML

**请求方法**

GET

**请求接口**
```
/html2json
```

**请求参数**

- `url` - [必需]需要解析的内容链接。
- `timeout` - 超时时间，单位为秒，默认为10秒
- `domain` - 图片等静态资源域名，用于拼装图片等链接。需带 `http` 或 `https`，如 `https://static.bookstack.cn`

> 注意：程序只解析 HTML 中的 Body 内容

**使用示例**

> http://localhost:8888/html2json?timeout=5&url=https://gitee.com/truthhun/BookStack


##### 解析Form表单提交HTML的内容

**请求方法**

POST

**请求接口**
```
/html2json
```

**请求参数**

- `html` - HTML内容字符串
- `domain` - 图片等静态资源域名，用于拼装图片等链接。需带 `http` 或 `https`，如 `https://static.bookstack.cn`


##### 解析form表单提交的markdown内容

**请求方法**

POST

**请求接口**
```
/md2json
```

**请求参数**

- `markdown` - [必需] markdown内容字符串
- `domain` - 图片等静态资源域名，用于拼装图片等链接。需带 `http` 或 `https`，如 `https://static.bookstack.cn`


### 以包的形式引用(针对Go语言)


#### 安装
```
go get -v github.com/TruthHun/html2json
```

#### 使用示例

```
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TruthHun/html2json/html2json"
)

func main()  {
	//rt:=html2json.NewDefault()
	appTags:=html2json.GetTags(html2json.TagUniAPP)
	rt:=html2json.New(appTags)
	htmlStr:=`
<div>
	hello world!
	<span>this is a span</span>
	a b c d e
	<img src="https://www.bookstack.cn/static/images/logo.png"/>
	<audio src="helloworld.mp3"></audio>
	<video src="../bookstack.mp4"></video>
<a href="https://www.bookstack.cn">书栈网 - 分享知识，共享智慧</a>
</div>
<iframe src="https://www.baidu.com" frameborder="0"></iframe>
<pre>
	this is pre code
</pre>
`
	now:=time.Now()
	nodes,err:=rt.Parse(htmlStr,"https://www.bookstack.cn/static/")
	if err!=nil{
		panic(err)
	}
	fmt.Println("spend time",time.Since(now))
	fmt.Println(toJSON(nodes))
}

func toJSON(v interface{}) (js string) {
	b,_:=json.Marshal(v)
	return string(b)
}
```

**示例代码输出结果**
```
[{
	"name": "div",
	"attrs": {
		"class": "tag-div"
	},
	"children": [{
		"type": "text",
		"text": "\n\thello world!\n\t"
	}, {
		"name": "span",
		"attrs": {
			"class": "tag-span"
		},
		"children": [{
			"type": "text",
			"text": "this is a span"
		}]
	}, {
		"type": "text",
		"text": "\n\ta b c d e\n\t"
	}, {
		"name": "img",
		"attrs": {
			"class": "tag-img",
			"src": "https://www.bookstack.cn/static/images/logo.png"
		}
	}, {
		"type": "text",
		"text": "\n\t"
	}, {
		"name": "a",
		"attrs": {
			"class": "tag-audio",
			"href": "https://www.bookstack.cn/static/helloworld.mp3",
		},
		"children": [{
			"type": "text",
			"text": " [audio] https://www.bookstack.cn/static/helloworld.mp3 "
		}]
	}, {
		"type": "text",
		"text": "\n\t"
	}, {
		"name": "a",
		"attrs": {
			"class": "tag-video",
			"href": "https://www.bookstack.cn/bookstack.mp4",
		},
		"children": [{
			"type": "text",
			"text": " [video] https://www.bookstack.cn/bookstack.mp4 "
		}]
	}, {
		"type": "text",
		"text": "\n"
	}, {
		"name": "a",
		"attrs": {
			"class": "tag-a",
			"href": "https://www.bookstack.cn"
		},
		"children": [{
			"type": "text",
			"text": "书栈网 - 分享知识，共享智慧"
		}]
	}, {
		"type": "text",
		"text": "\n"
	}]
}, {
	"type": "text",
	"text": "\n"
}, {
	"name": "a",
	"attrs": {
		"class": "tag-iframe",
		"frameborder": "0",
		"href": "https://www.baidu.com",
	},
	"children": [{
		"type": "text",
		"text": " [iframe] https://www.baidu.com "
	}]
}, {
	"type": "text",
	"text": "\n"
}, {
	"name": "div",
	"attrs": {
		"class": "tag-pre",
		"style": "display: block;font-family: monospace;white-space: pre;margin: 1em 0;"
	},
	"children": [{
		"type": "text",
		"text": "\tthis is pre code\n"
	}]
}, {
	"type": "text",
	"text": "\n"
}]
```

## 说明

所有标签都会生成一个 `"tag-"+标签名`的`class`，以便于对标签进行样式控制。

比如 `a`标签，会添加上`tag-a`的class，`div`标签会添加一个`tag-div`，`code`标签会添加一个 `tag-code`的class，以此类推。

**特别注释事项**

由于部分小程序`rich-text`组件并不支持`pre`标签，所以`pre`标签会被转为`div`标签，并且多出一个`tag-pre`的class，同时会在增加一个
`pre`标签本身默认的css样式：

```
display: block;
font-family: monospace;
white-space: pre;
margin: 1em 0;
```

同时，如果`video`、`iframe`、`audio`标签，如果不在信任的标签里面，则作为`a`标签处理


## 程序体验

编译好了的程序，只有一个可执行文件，部署和使用都和简单。

如果是部署到线上生产环境，需要设置守护进程。

> 注意：如果服务不对外提供API请求服务，建议把端口不对外网开放