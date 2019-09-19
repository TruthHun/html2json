package html2json

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
)

// html to json struct
type h2j struct {
	Name     string            `json:"name,omitempty"` // 对应 HTML 标签
	Type     string            `json:"type,omitempty"` // element 或者 text
	Text     string            `json:"text,omitempty"`
	Attrs    map[string]string `json:"attrs,omitempty"`
	Children []h2j             `json:"children,omitempty"`
}

type RichText struct {
	tagsMap sync.Map
}

func NewDefault() *RichText {
	return New(defaultTags)
}

func New(customTags []string) *RichText {
	if len(customTags) == 0 {
		customTags = defaultTags
	}
	tagsMap := sync.Map{}
	for _, tag := range customTags {
		tagsMap.Store(strings.ToLower(tag), true)
	}
	return &RichText{
		tagsMap: tagsMap,
	}
}

func (r *RichText) Parse(htmlStr string) (data []h2j, err error) {
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return
	}
	doc.Find("body").Each(func(i int, selection *goquery.Selection) {
		data = r.parse(selection)
	})
	return
}

func (r *RichText) ParseByByte(htmlByte []byte) (data []h2j, err error) {
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(htmlByte))
	if err != nil {
		return
	}
	doc.Find("body").Each(func(i int, selection *goquery.Selection) {
		data = r.parse(selection)
	})
	return
}

func (r *RichText) ParseByURL(urlStr string, expire ...int) (data []h2j, err error) {
	var (
		resp *http.Response
		b    []byte
	)
	req := httplib.Get(urlStr)
	if strings.HasPrefix(strings.ToLower(urlStr), "https://") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	req.Header("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.87 Safari/537.36")
	ex := 10 * time.Second
	if len(expire) > 0 && expire[0] > 0 {
		ex = time.Duration(expire[0]) * time.Second
	}
	resp, err = req.SetTimeout(ex, ex).Response()
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return r.ParseByByte(b)
}

func (r *RichText) parse(sel *goquery.Selection) (data []h2j) {
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
				h.Name = strings.ToLower(item.Data)

				// 忽略script
				if h.Name == "script" {
					continue
				}

				// attrs
				attr := make(map[string]string)
				for _, a := range item.Attr {
					attr[a.Key] = a.Val
				}

				if class, ok := attr["class"]; ok {
					attr["class"] = fmt.Sprintf("tag-%v %v", h.Name, class)
				} else {
					attr["class"] = "tag-" + h.Name
				}

				// 小程序不支持的HTML标签，全部转为div标签
				if _, ok := r.tagsMap.Load(h.Name); !ok {
					switch h.Name {
					case "pre":
						h.Name = "div"
						defaultStyle := "display: block;font-family: monospace;white-space: pre;margin: 1em 0;" // set default <pre> css
						if style, ok := attr["style"]; ok {
							attr["style"] = defaultStyle + style
						} else {
							attr["style"] = defaultStyle
						}
					case "audio", "video", "iframe":
						h.Name = "a"
						if src, ok := attr["src"]; ok {
							attr["href"] = src
							h.Children = []h2j{{Type: "text", Text: fmt.Sprintf(" [audio]%v ", src)}}
						}
					default:
						h.Name = "div"
					}
				}
				h.Attrs = attr
				h.Children = r.parse(goquery.NewDocumentFromNode(item).Selection)
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
