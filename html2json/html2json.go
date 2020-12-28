package html2json

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/russross/blackfriday"

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

type inode struct {
	Type string `json:"type"`
	Data []h2j  `json:"data"`
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

func (r *RichText) ParseMarkdown(md, domain string) (data []h2j, err error) {
	return r.ParseMarkdownByByte([]byte(md), domain)
}

func (r *RichText) ParseMarkdownByByte(mdByte []byte, domain string) (data []h2j, err error) {
	return r.ParseByByte(blackfriday.Run(mdByte), domain)
}

func (r *RichText) Parse(htmlStr string, domain string) (data []h2j, err error) {
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return
	}
	doc.Find("body").Each(func(i int, selection *goquery.Selection) {
		data = r.parse(selection, domain)
	})
	return
}

func (r *RichText) ParseByByte(htmlByte []byte, domain string) (data []h2j, err error) {
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(htmlByte))
	if err != nil {
		return
	}
	doc.Find("body").Each(func(i int, selection *goquery.Selection) {
		data = r.parse(selection, domain)
	})
	return
}

func (r *RichText) ParseByByteV2(htmlByte []byte, domain string) (inodes []inode, err error) {
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(htmlByte))
	if err != nil {
		return
	}

	splitMark := "$@$@$@$"
	mediaTags := map[string]bool{"audio": true, "video": true, "iframe": true, "img": true}
	blockTags := map[string]bool{"article": true, "aside": true, "base": true, "body": true, "center": true, "figure": true, "nav": true, "title": true, "h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true, "p": true, "div": true}
	doc.Find("body").Each(func(i int, selection *goquery.Selection) {
		for tag, _ := range mediaTags {
			doc.Find(tag).Each(func(idx int, sel *goquery.Selection) {
				if tag != "img" {
					sel.BeforeHtml(splitMark)
					sel.AfterHtml(splitMark)
				} else {
					tagName := "body"
					if len(sel.Parent().Nodes) > 0 {
						tagName = strings.ToLower(sel.Parent().Nodes[0].DataAtom.String())
					}
					// 父节点为块节点且除了图片之外没有其他内容
					if _, ok := blockTags[tagName]; ok && strings.TrimSpace(sel.Parent().Text()) == "" {
						sel.BeforeHtml(splitMark)
						sel.AfterHtml(splitMark)
					}
				}
			})
		}
	})

	ret, _ := doc.Find("body").Html()
	slice := strings.Split(ret, splitMark)

	var data []h2j
	for _, item := range slice {
		if strings.TrimSpace(item) != "" {
			doc2, _ := goquery.NewDocumentFromReader(strings.NewReader(item))
			data = append(data, r.parseV2(doc2.Find("body"), domain)...)
		}
	}

	var (
		idata []h2j
		l     = len(data)
	)

	for idx, item := range data {
		if _, ok := mediaTags[item.Name]; ok {
			if len(idata) > 0 {
				inodes = append(inodes, inode{"richtext", idata})
			}
			inodes = append(inodes, inode{item.Name, []h2j{item}})
			idata = make([]h2j, 0)
		} else {
			idata = append(idata, item)
			if idx == l-1 {
				inodes = append(inodes, inode{"richtext", idata})
			}
		}
	}

	return
}

func (r *RichText) ParseByURL(urlStr string, domain string, timeout ...int) (data []h2j, err error) {
	var (
		resp *http.Response
		b    []byte
	)
	req := httplib.Get(urlStr)
	if strings.HasPrefix(strings.ToLower(urlStr), "https://") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	req.Header("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.87 Safari/537.36")
	to := 10 * time.Second
	if len(timeout) > 0 && timeout[0] > 0 {
		to = time.Duration(timeout[0]) * time.Second
	}
	resp, err = req.SetTimeout(to, to).Response()
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return r.ParseByByte(b, domain)
}

func (r *RichText) parseV2(sel *goquery.Selection, domain string) (data []h2j) {
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
				if h.Name == "script" || h.Name == "link" {
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

				switch h.Name {
				case "img", "audio", "video", "iframe":
					if src, ok := attr["src"]; ok {
						attr["src"] = r.fixSourceLink(domain, src)
					}
				case "a":
					if href, ok := attr["href"]; ok {
						attr["href"] = r.fixSourceLink(domain, href)
					}
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
						if src, ok := attr["src"]; ok {
							src = r.fixSourceLink(domain, src)
						}
					default:
						h.Name = "div"
					}
				}
				h.Attrs = attr
				if len(h.Children) == 0 {
					h.Children = r.parseV2(goquery.NewDocumentFromNode(item).Selection, domain)
				}
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

func (r *RichText) parse(sel *goquery.Selection, domain string) (data []h2j) {
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
				if h.Name == "script" || h.Name == "link" {
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

				switch h.Name {
				case "img", "audio", "video":
					if src, ok := attr["src"]; ok {
						attr["src"] = r.fixSourceLink(domain, src)
					}
				case "a":
					if href, ok := attr["href"]; ok {
						attr["href"] = r.fixSourceLink(domain, href)
					}
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
						if src, ok := attr["src"]; ok {
							src = r.fixSourceLink(domain, src)
							attr["href"] = src
							delete(attr, "src")
							h.Children = []h2j{{Type: "text", Text: fmt.Sprintf(" [%v] %v ", h.Name, src)}}
						}
						h.Name = "a"
					default:
						h.Name = "div"
					}
				}
				h.Attrs = attr
				if len(h.Children) == 0 {
					h.Children = r.parse(goquery.NewDocumentFromNode(item).Selection, domain)
				}
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

func (r *RichText) fixSourceLink(domain, link string) string {
	if domain == "" {
		return link
	}

	link = strings.ReplaceAll(link, "\\", "/")

	if strings.HasPrefix(link, "//") {
		return "http:" + link
	}

	linkLower := strings.ToLower(link)
	if strings.HasPrefix(linkLower, "https://") || strings.HasPrefix(linkLower, "http://") {
		return link
	}

	u, err := url.Parse(domain)

	if err != nil {
		return link
	}

	if strings.HasPrefix(link, "/") {
		return u.Scheme + "://" + u.Host + "/" + strings.TrimLeft(link, "/")
	}
	u.Path = path.Join(strings.TrimRight(u.Path, "/")+"/", link)

	// return u.String() // 会对中文进行编码
	return u.Scheme + "://" + u.Host + "/" + strings.Trim(u.Path, "/")
}
