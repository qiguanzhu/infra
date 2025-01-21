/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/8/27 -- 17:10
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: operate.go
*/

package xhttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var _MethodMap = map[string]func(url string) CfgOp{
	http.MethodGet:     Get,
	http.MethodPost:    Post,
	http.MethodPut:     Put,
	http.MethodPatch:   Patch,
	http.MethodDelete:  Delete,
	http.MethodOptions: Option,
}

type HttpConfig httpConfig

type httpConfig struct {
	url         string
	method      string
	timeout     time.Duration
	headers     http.Header
	cookies     *http.Cookie
	tlsCfg      *tls.Config
	requestType string
	reader      reader
	loader      loadFunc
	bodySize    int

	Prefix []ReqPrefixFunc
	Suffix []ResSuffixFunc
}

func (h *httpConfig) Use(opts ...CfgOp) {
	for _, opt := range opts {
		opt(h)
	}
}

func (h *httpConfig) MakeReq(ctx context.Context, v any) (req *http.Request, err error) {
	body, err := h.reader(v)
	if err != nil {
		return nil, err
	}
	req, err = http.NewRequestWithContext(ctx, h.method, h.url, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (h *httpConfig) MakeRes(ctx context.Context, res *http.Response, bs []byte, tar any) (*http.Response, error) {
	err := h.loader(bs, tar)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *httpConfig) MakeResOk(ctx context.Context, res *http.Response, bs []byte, tar any) (*http.Response, error) {
	if res.StatusCode != http.StatusOK {
		return res, fmt.Errorf("StatusCode(%d) != 200", res.StatusCode)
	}
	return h.MakeRes(ctx, res, bs, tar)
}

type CfgOp func(*httpConfig)

func Post(url string) CfgOp {
	return func(c *httpConfig) {
		c.method = http.MethodPost
		c.url = url
	}
}

func Get(url string) CfgOp {
	return func(c *httpConfig) {
		c.method = http.MethodGet
		c.url = url
	}
}

func Put(url string) CfgOp {
	return func(c *httpConfig) {
		c.method = http.MethodPut
		c.url = url
	}
}

func Delete(url string) CfgOp {
	return func(c *httpConfig) {
		c.method = http.MethodDelete
		c.url = url
	}
}

func Patch(url string) CfgOp {
	return func(c *httpConfig) {
		c.method = http.MethodPatch
		c.url = url
	}
}

func Option(url string) CfgOp {
	return func(c *httpConfig) {
		c.method = http.MethodOptions
		c.url = url
	}
}

func Header(headers map[string]string) CfgOp {
	return func(c *httpConfig) {
		if c.headers == nil {
			c.headers = http.Header{}
		}
		if headers == nil {
			return
		}
		for k, v := range headers {
			c.headers.Set(k, v)
		}
	}
}

func Cookies(cookies *http.Cookie) CfgOp {
	return func(c *httpConfig) {
		c.cookies = cookies
	}
}

func Prefix(prefix ...ReqPrefixFunc) CfgOp {
	return func(c *httpConfig) {
		c.Prefix = append([]ReqPrefixFunc{}, prefix...)
	}
}

func Suffix(suffix ...ResSuffixFunc) CfgOp {
	return func(c *httpConfig) {
		c.Suffix = append([]ResSuffixFunc{}, suffix...)
	}
}

func AppendPrefix(prefix ...ReqPrefixFunc) CfgOp {
	return func(c *httpConfig) {
		c.Prefix = append(c.Prefix, prefix...)
	}
}

func AppendSuffix(suffix ...ResSuffixFunc) CfgOp {
	return func(c *httpConfig) {
		c.Suffix = append(c.Suffix, suffix...)
	}
}

func TLS(tlsC *tls.Config) CfgOp {
	return func(c *httpConfig) {
		c.tlsCfg = tlsC
	}
}

func Req(formType string) CfgOp {
	return func(c *httpConfig) {
		c.reader = _ReqContentTypeReader[formType](c)
		Header(map[string]string{
			"Content-Type": _ReqContentTypeMap[formType],
		})(c)
	}
}

func Res(formType string) CfgOp {
	return func(c *httpConfig) {
		c.loader = _ResContentTypeLoader[formType]
	}
}

// BodySize set body size (MB), default is 10MB
func BodySize(sizeMB int) CfgOp {
	return func(c *httpConfig) {
		c.bodySize = sizeMB
	}
}

func Timeout(timeout time.Duration) CfgOp {
	return func(c *httpConfig) {
		c.timeout = timeout
	}
}

// ReqPrefixFunc a hook Before the request started
type ReqPrefixFunc func(context context.Context, req *http.Request) context.Context

// ResSuffixFunc a hook After the request finished
type ResSuffixFunc func(context context.Context, res *http.Response) context.Context

func DefaultReqPrefix(ctx context.Context, req *http.Request) context.Context {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("[Read Req body] error: %s", err.Error())
		return ctx
	}
	enEscapeUrl, err := url.QueryUnescape(string(body))
	if err == nil {
		fmt.Printf("[Req] %s", enEscapeUrl)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	return ctx
}

func DefaultResSuffix(ctx context.Context, res *http.Response) context.Context {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("[Read Res body] error: %s", err.Error())
		return ctx
	}
	res.Body = io.NopCloser(bytes.NewBuffer(body))
	fmt.Printf("[Res] %s", string(body))
	return ctx
}
