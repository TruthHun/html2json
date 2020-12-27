package html2json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var rt = NewDefault()

func toJSON(v interface{}) (js string) {
	b, _ := json.Marshal(v)
	return string(b)
}

func TestRichText_ParseMarkdownByByte(t *testing.T) {
	b, _ := ioutil.ReadFile("examples/bookstack-readme.md")
	now := time.Now()
	nodes, err := rt.ParseMarkdownByByte(b, "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(time.Since(now))
	t.Log(toJSON(nodes))
}

func TestParse_BigHTML(t *testing.T) {
	b, _ := ioutil.ReadFile("examples/gin.html")
	now := time.Now()
	nodes, err := rt.ParseByByte(b, "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(time.Since(now))
	t.Log(toJSON(nodes))
}

func TestParse_Media(t *testing.T) {
	b, _ := ioutil.ReadFile("examples/media.html")
	nodes, err := rt.ParseByByteV2(b, "")
	if err != nil {
		t.Error(err)
	}
	t.Log(toJSON(nodes))
	ioutil.WriteFile("examples/media.json", []byte(toJSON(nodes)), os.ModePerm)
}

func TestParseByURL(t *testing.T) {
	nodes, err := rt.ParseByURL("https://my.oschina.net/huanghaibin/blog/3106432", "")
	if err != nil {
		t.Error(err)
	}
	t.Log(toJSON(nodes))
}

func BenchmarkParse(b *testing.B) {
	h, _ := ioutil.ReadFile("examples/gin.html")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt.ParseByByteV2(h, "")
	}
}

func BenchmarkParse_uniapp(b *testing.B) {
	h, _ := ioutil.ReadFile("examples/uniapp.html")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt.ParseByByte(h, "")
	}
}

func BenchmarkParse2(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt.Parse("<div>hello world</div>", "")
	}
}
