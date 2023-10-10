package nano

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// NewHTTPEndpointGetRequestWithAuth 新建带鉴权头的 HTTP GET 请求
func NewHTTPEndpointGetRequestWithAuth(ep string, auth string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", StandardAPI+ep, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", auth)
	return
}

// NewHTTPEndpointDeleteRequestWithAuth 新建带鉴权头的 HTTP DELETE 请求
func NewHTTPEndpointDeleteRequestWithAuth(ep string, auth string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest("DELETE", StandardAPI+ep, body)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", auth)
	return
}

// NewHTTPEndpointPostRequestWithAuth 新建带鉴权头的 HTTP POST 请求
func NewHTTPEndpointPostRequestWithAuth(ep string, auth string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest("POST", StandardAPI+ep, body)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", auth)
	return
}

// NewHTTPEndpointPatchRequestWithAuth 新建带鉴权头的 HTTP PATCH 请求
func NewHTTPEndpointPatchRequestWithAuth(ep string, auth string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest("PATCH", StandardAPI+ep, body)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", auth)
	return
}

// WriteHTTPQueryIfNotNil 如果非空则将请求添加到 baseurl 后
//
// ex. WriteHTTPQueryIfNotNil("http://a.com/api", "a", 0, "b", 1, "c", 2) is http://a.com/api?b=1&c=2
func WriteHTTPQueryIfNotNil(baseurl string, queries ...any) string {
	if len(queries) == 0 {
		return baseurl
	}
	hasstart := false
	queryname := ""
	sb := strings.Builder{}
	for i, q := range queries {
		if i%2 == 0 {
			queryname = q.(string)
			continue
		}
		if reflect.ValueOf(q).IsZero() {
			continue
		}
		if !hasstart {
			sb.WriteString(baseurl)
			sb.WriteByte('?')
			hasstart = true
		}
		sb.WriteString(queryname)
		sb.WriteByte('=')
		sb.WriteString(url.QueryEscape(fmt.Sprint(q)))
		sb.WriteByte('&')
	}
	if sb.Len() <= 4 {
		return baseurl
	}
	return sb.String()[:sb.Len()-1]
}

// WriteBodyFromJSON 从 json 结构体 ptr 写入 bytes.Buffer, 忽略 error (内部使用不会出错)
func WriteBodyFromJSON(ptr any) *bytes.Buffer {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	_ = json.NewEncoder(buf).Encode(ptr)
	return buf
}
