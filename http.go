package nano

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	base14 "github.com/fumiama/go-base16384"
)

// HTTPRequsetConstructer ...
type HTTPRequsetConstructer func(ep string, contenttype string, auth string, body io.Reader) (*http.Request, error)

func newHTTPEndpointRequestWithAuth(method, contenttype, ep string, auth string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, OpenAPI+ep, body)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", auth)
	if contenttype == "" {
		contenttype = "application/json"
	}
	req.Header.Add("Content-Type", contenttype)
	return
}

// NewHTTPEndpointGetRequestWithAuth 新建带鉴权头的 HTTP GET 请求
func NewHTTPEndpointGetRequestWithAuth(ep string, contenttype string, auth string, body io.Reader) (*http.Request, error) {
	return newHTTPEndpointRequestWithAuth("GET", contenttype, ep, auth, body)
}

// NewHTTPEndpointPutRequestWithAuth 新建带鉴权头的 HTTP PUT 请求
func NewHTTPEndpointPutRequestWithAuth(ep string, contenttype string, auth string, body io.Reader) (*http.Request, error) {
	return newHTTPEndpointRequestWithAuth("PUT", contenttype, ep, auth, body)
}

// NewHTTPEndpointDeleteRequestWithAuth 新建带鉴权头的 HTTP DELETE 请求
func NewHTTPEndpointDeleteRequestWithAuth(ep string, contenttype string, auth string, body io.Reader) (*http.Request, error) {
	return newHTTPEndpointRequestWithAuth("DELETE", contenttype, ep, auth, body)
}

// NewHTTPEndpointPostRequestWithAuth 新建带鉴权头的 HTTP POST 请求
func NewHTTPEndpointPostRequestWithAuth(ep string, contenttype string, auth string, body io.Reader) (*http.Request, error) {
	return newHTTPEndpointRequestWithAuth("POST", contenttype, ep, auth, body)
}

// NewHTTPEndpointPatchRequestWithAuth 新建带鉴权头的 HTTP PATCH 请求
func NewHTTPEndpointPatchRequestWithAuth(ep string, contenttype string, auth string, body io.Reader) (*http.Request, error) {
	return newHTTPEndpointRequestWithAuth("PATCH", contenttype, ep, auth, body)
}

// WriteHTTPQueryIfNotNil 如果非空则将请求添加到 baseurl 后
//
// ex. WriteHTTPQueryIfNotNil("http://a.com/api", "a", 0, "b", 1, "c", 2) is http://a.com/api?b=1&c=2
func WriteHTTPQueryIfNotNil(baseurl string, queries ...any) string {
	if len(queries) < 2 {
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

// WriteBodyByMultipartFormData 使用 multipart/form-data 上传
func WriteBodyByMultipartFormData(params ...any) (*bytes.Buffer, string, error) {
	if len(params)%2 != 0 {
		panic("invalid params to " + getThisFuncName())
	}
	fieldname := ""
	buf := bytes.NewBuffer(make([]byte, 0, 65536))
	w := multipart.NewWriter(buf)
	defer w.Close()
	for i, x := range params {
		if i%2 == 0 { // 参数
			fieldname = x.(string)
			continue
		}
		rx := reflect.ValueOf(x)
		if rx.IsZero() {
			continue
		}
		r, err := w.CreateFormField(fieldname)
		if err != nil {
			return nil, "", err
		}
		if rx.Kind() == reflect.Pointer && rx.Elem().Kind() == reflect.Struct { // 使用 json 编码
			err = json.NewEncoder(r).Encode(x)
			if err != nil {
				return nil, "", err
			}
			continue
		}
		switch o := x.(type) {
		case string:
			if strings.HasPrefix(o, "file:///") { // 是文件路径
				f, err := os.Open(o[8:])
				if err != nil {
					return nil, "", err
				}
				defer f.Close()
				_, err = io.Copy(r, f)
				if err != nil {
					return nil, "", err
				}
				continue
			}
			if strings.HasPrefix(o, "base64://") { // 是 base64
				_, err = io.Copy(r, base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(o[9:])))
				if err != nil {
					return nil, "", err
				}
				continue
			}
			if strings.HasPrefix(o, "base16384://") { // 是 base16384
				_, err = io.Copy(r, base14.NewDecoder(bytes.NewBufferString(o[12:])))
				if err != nil {
					return nil, "", err
				}
				continue
			}
			_, err = io.WriteString(r, o)
			if err != nil {
				return nil, "", err
			}
			continue
		case []byte:
			_, err = r.Write(o)
			if err != nil {
				return nil, "", err
			}
			continue
		default:
			_, err = io.WriteString(r, fmt.Sprint(o))
			if err != nil {
				return nil, "", err
			}
			continue
		}
	}
	return buf, w.FormDataContentType(), nil
}
