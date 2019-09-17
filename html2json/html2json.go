package html2json

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
)

type h2j struct {
	Name     string            `json:"name,omitempty"` // 对应 HTML 标签
	Type     string            `json:"type,omitempty"` // element 或者 text
	Text     string            `json:"text,omitempty"`
	Attrs    map[string]string `json:"attrs,omitempty"`
	Children []h2j             `json:"children,omitempty"`
}

// replace
type rp struct {
	Tag          string
	DefaultStyle string
}

// 微信小程序支持的 HTML 标签：https://developers.weixin.qq.com/miniprogram/dev/component/rich-text.html
var richTextTags = map[string]bool{"a": true, "abbr": true, "address": true, "article": true, "aside": true, "b": true,
	"bdi": true, "bdo": true, "big": true, "blockquote": true, "br": true, "caption": true, "center": true,
	"cite": true, "code": true, "col": true, "colgroup": true, "dd": true, "del": true, "div": true, "dl": true, "dt": true, "em": true,
	"fieldset": true, "font": true, "footer": true, "h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"header": true, "hr": true, "i": true, "img": true, "ins": true, "label": true, "legend": true, "li": true, "mark": true,
	"nav": true, "ol": true, "p": true, "q": true, "rt": true, "ruby": true,
	"s": true, "section": true, "small": true, "span": true, "strong": true, "sub": true, "sup": true,
	"table": true, "tbody": true, "td": true, "tfoot": true, "th": true, " thead": true, "tr": true, "tt": true, "u": true, "ul": true}

func Parse(htmlStr string) (data []h2j, err error) {
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return
	}
	doc.Find("body").Each(func(i int, selection *goquery.Selection) {
		data = parse(selection)
	})
	return
}

func parse(sel *goquery.Selection) (data []h2j) {
	nodes := sel.Children().Nodes
	if len(nodes) == 0 {
		if txt := sel.Text(); txt != "" {
			data = []h2j{{Text: txt, Type: "text"}}
		}
		return
	}
	sel.Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
		ns := s.Nodes
		for _, item := range ns {
			var h h2j
			if item.Type != html.TextNode {
				h.Name = item.Data
				attr := make(map[string]string)
				for _, a := range item.Attr {
					attr[a.Key] = a.Val
				}

				if class, ok := attr["class"]; ok {
					attr["class"] = fmt.Sprintf("tag-%v %v", h.Name, class)
				} else {
					attr["class"] = "tag-" + h.Name
				}
				if h.Name == "pre" {
					h.Name = "div"
					// set default <pre> css
					defualtStyle := "display: block;font-family: monospace;white-space: pre;margin: 1em 0;"
					if style, ok := attr["style"]; ok {
						attr["style"] = defualtStyle + style
					} else {
						attr["style"] = defualtStyle
					}
				}
				h.Attrs = attr
				h.Children = parse(goquery.NewDocumentFromNode(item).Selection)
			} else {
				h.Type = "text"
				h.Text = goquery.NewDocumentFromNode(item).Selection.Text()
			}
			data = append(data, h)
		}
		return true
	})
	return
}

func videoToNode() {

}

func iframeToNode() {

}

func audioToNode() {

}
