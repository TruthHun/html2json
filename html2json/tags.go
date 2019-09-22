package html2json

// 各小程序支持的HTML标签
//		微信小程序：https://developers.weixin.qq.com/miniprogram/dev/component/rich-text.html
//		支付宝小程序：https://docs.alipay.com/mini/component/rich-text
//		百度小程序：https://smartprogram.baidu.com/docs/develop/component/base/#rich-text-%E5%AF%8C%E6%96%87%E6%9C%AC/
//		头条小程序：https://developer.toutiao.com/dev/miniapp/uEDMy4SMwIjLxAjM
//		QQ小程序：https://q.qq.com/wiki/develop/miniprogram/component/basic-content/rich-text.html
//		uni-app: https://uniapp.dcloud.io/component/rich-text?id=rich-text
//  我们这里以 uni-app 支持的标签为默认支持的标签
var (
	defaultTags = []string{
		"a", "abbr", "b", "blockquote", "br", "code", "col", "colgroup", "dd", "del", "div", "dl", "dt", "em", "fieldset", "h1", "h2", "h3", "h4", "h5", "h6", "header", "hr", "i", "img", "ins", "label", "legend", "li", "ol", "p", "q", "span", "strong", "sub", "sup", "table", "tbody", "td", "tfoot", "th", "thead", "tr", "tt", "ul",
	}
	alipayTags = []string{
		"a", "abbr", "b", "blockquote", "br", "code", "col", "colgroup", "dd", "del", "div", "dl", "dt", "em", "fieldset", "h1", "h2", "h3", "h4", "h5", "h6", "hr", "i", "img", "ins", "label", "legend", "li", "ol", "p", "q", "span", "strong", "sub", "sup", "table", "tbody", "td", "tfoot", "th", "thead", "tr", "ul",
	}
	baiduTags = []string{
		"a", "abbr", "b", "blockquote", "br", "code", "col", "colgroup", "dd", "del", "div", "dl", "dt", "em", "fieldset", "h1", "h2", "h3", "h4", "h5", "h6", "hr", "i", "img", "ins", "label", "legend", "li", "ol", "p", "q", "span", "strong", "sub", "sup", "table", "tbody", "td", "tfoot", "th", "thead", "tr", "ul",
	}
	qqTags = []string{
		"a", "abbr", "b", "blockquote", "br", "code", "col", "colgroup", "dd", "del", "div", "dl", "dt", "em", "fieldset", "h1", "h2", "h3", "h4", "h5", "h6", "hr", "i", "img", "ins", "label", "legend", "li", "ol", "p", "q", "span", "strong", "sub", "sup", "table", "tbody", "td", "tfoot", "th", "thead", "tr", "ul",
	}
	mpTags = []string{
		"a", "abbr", "address", "article", "aside", "b", "bdi", "bdo", "big", "blockquote", "br", "caption", "center", "cite", "code", "col", "colgroup", "dd", "del", "div", "dl", "dt", "em", "fieldset", "font", "footer", "h1", "h2", "h3", "h4", "h5", "h6", "header", "hr", "i", "img", "ins", "label", "legend", "li", "mark", "nav", "ol", "p", "pre", "q", "rt", "ruby", "s", "section", "small", "span", "strong", "sub", "sup", "table", "tbody", "td", "tfoot", "th", "thead", "tr", "tt", "u", "ul",
	}
	toutiaoTags = defaultTags // 没看到有限定的信任标签
)

type Tag string

const (
	TagBaidu   = "baidu"
	TagAplipay = "alipay"
	TagQQ      = "qq"
	TagWeixin  = "weixin"
	TagToutiao = "toutiao"
	TagUniAPP  = "uni-app"
)

func GetTags(cate Tag) []string {
	switch cate {
	case TagBaidu:
		return baiduTags
	case TagAplipay:
		return alipayTags
	case TagQQ:
		return qqTags
	case "mp", TagWeixin:
		return mpTags
	case "tt", TagToutiao:
		return toutiaoTags
	case TagUniAPP:
		return defaultTags
	default:
		return defaultTags
	}
}
