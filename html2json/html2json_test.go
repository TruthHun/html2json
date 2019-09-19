package html2json

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

var rt = NewDefault()

func toJSON(v interface{}) (js string) {
	b, _ := json.Marshal(v)
	return string(b)
}

func TestParse_BigHTML(t *testing.T) {
	b, _ := ioutil.ReadFile("examples/gin.html")
	nodes, err := rt.ParseByByte(b)
	if err != nil {
		t.Error(err)
	}
	t.Log(toJSON(nodes))
}

func TestParse_Media(t *testing.T) {
	b, _ := ioutil.ReadFile("examples/media.html")
	nodes, err := rt.ParseByByte(b)
	if err != nil {
		t.Error(err)
	}
	t.Log(toJSON(nodes))
}

func TestParseByURL(t *testing.T) {
	nodes, err := rt.ParseByURL("https://my.oschina.net/huanghaibin/blog/3106432")
	if err != nil {
		t.Error(err)
	}
	t.Log(toJSON(nodes))
}

func BenchmarkParse(b *testing.B) {
	h, _ := ioutil.ReadFile("examples/gin.html")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt.ParseByByte(h)
	}
}

func BenchmarkParse_uniapp(b *testing.B) {
	h, _ := ioutil.ReadFile("examples/uniapp.html")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt.ParseByByte(h)
	}
}

func BenchmarkParse2(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rt.Parse("<div>hello world</div>")
	}
}
