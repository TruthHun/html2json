# html2json

使用Go语言开发的HTML转JSON工具，将HTML内容转换为符合各种小程序`rich-text`组件内容渲染所需格式的`JSON`

## 介绍

在开发`BookStack`的配套微信小程序`BookChat`以及使用 `uni-app` 开发配套的手机APP应用`BookChatApp`的过程中，
我尝试了很多种开源的小程序HTML解析渲染工具，但是都不是很满意，主要体现在以下几点：

1. 性能不好，影响体验。
    
    书栈网的HTML内容长度大小不一，大一点(100kb左右)的HTML内容渲染，要10多秒，这种对于用户体验来说，是难以忍受的。

1. 稳定性不高，容错不够。
    
    如果把HTML文本字符串传递给渲染组件，渲染组件在将HTML解析成元素节点的过程中，使用到了正则表达式，内容中
    出现不可预料的内容的时候，解析就会出错，内容变成一片空白。
    
1. 渲染效果不理想。
    
    表格、代码块渲染效果不好，样式也不好控制

基于以上，所以用Go语言开发实现了这么个转换工具，将HTML转为JSON。

对我来说，在后端将HTML转为JSON，并配合小程序`rich-text`组件对内容进行渲染，性能和稳定性以及渲染效果，
都比较符合预期，尽管并没有第三方HTML渲染工具那样提供了图片预览的功能。

目前已经在`BookStack` v2.1 版本中使用了。


## 使用方式

```
./html2json --help
```

### resetful 方式使用

#### 启动服务
```
./html2json serve --port 8888 --tags weixin-html-tags.json
```

- `--port` - 指定服务端口，默认为 8888
- `--tags` - 指定信任的HTML元素，json数组文件，里面存放各个支持的HTML标签。默认使用 uni-app 信任的HTML标签

各小程序支持的HTML标签

> - 微信小程序：https://developers.weixin.qq.com/miniprogram/dev/component/rich-text.html
> - 支付宝小程序：https://docs.alipay.com/mini/component/rich-text
> - 百度小程序：https://smartprogram.baidu.com/docs/develop/component/base/#rich-text-%E5%AF%8C%E6%96%87%E6%9C%AC/
> - 头条小程序：https://developer.toutiao.com/dev/miniapp/uEDMy4SMwIjLxAjM
> - QQ小程序：https://q.qq.com/wiki/develop/miniprogram/component/basic-content/rich-text.html
> - uni-app: https://uniapp.dcloud.io/component/rich-text?id=rich-text

#### API接口

##### 解析来自url链接的HTML

**请求方法**

GET

**请求接口**
```
/html2json
```

**使用示例**

> http://localhost:8888/html2json?timeout=5&url=https://gitee.com/truthhun/BookStack


**参数说明**

- `timeout` - 超时时间，单位为秒，默认为10秒
- `url` - 需要解析的内容链接

### 以包的形式引用

## 演示地址
