package html2json

import (
	"encoding/json"
	"testing"
)

func toJSON(v interface{}) (js string) {
	b, _ := json.Marshal(v)
	return string(b)
}

func TestParse(t *testing.T) {

}
