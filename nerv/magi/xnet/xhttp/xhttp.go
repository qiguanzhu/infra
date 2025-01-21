package xhttp

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"
)

func GetDefaultClient() HttpClientWrapper {
	cfg := &httpConfig{
		Prefix:   []ReqPrefixFunc{},
		Suffix:   []ResSuffixFunc{},
		bodySize: 10,
	}
	cfg.Use(defaultLoader())
	cfg.Use(defaultReader())
	return &HcWrapper{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 128,
				MaxConnsPerHost:     1024,
				Proxy:               http.ProxyFromEnvironment,
				DialContext: defaultTransportDialContext(&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}),
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     true,
				ForceAttemptHTTP2:     true,
			},
			Timeout: 60 * time.Second,
		},
		cfg: cfg,
	}
}

func HttpReqGetOk(url string, timeout time.Duration) ([]byte, error) {
	_, bs, err := GetDefaultClient().CallOpOk(context.Background(), map[string]any{},
		Get(url),
		Timeout(timeout),
	)
	return bs, err
}

func HttpReqPostOk(url string, data map[string]any, timeout time.Duration) ([]byte, error) {
	_, bs, err := GetDefaultClient().CallOpOk(context.Background(), data,
		Post(url),
		Timeout(timeout),
	)
	return bs, err
}

func HttpReqOk(url, method string, data map[string]any, timeout time.Duration) ([]byte, error) {
	if _, ok := _MethodMap[method]; !ok {
		return nil, errors.New("not supported http method")
	}
	_, bs, err := GetDefaultClient().CallOpOk(context.Background(), data,
		_MethodMap[method](url),
		Timeout(timeout),
	)
	return bs, err
}

func HttpReqPost(url string, data map[string]any, timeout time.Duration) ([]byte, int, error) {
	res, bs, err := GetDefaultClient().CallOp(context.Background(), data,
		Post(url),
		Timeout(timeout),
	)
	if err != nil {
		return nil, 0, err
	}
	return bs, res.StatusCode, err
}

func HttpReq(url, method string, data map[string]any, timeout time.Duration) ([]byte, int, error) {
	if _, ok := _MethodMap[method]; !ok {
		return nil, 0, errors.New("not supported http method")
	}
	res, bs, err := GetDefaultClient().CallOpOk(context.Background(), data,
		_MethodMap[method](url),
		Timeout(timeout),
	)
	if err != nil {
		return nil, 0, err
	}
	return bs, res.StatusCode, err
}

func HttpReqWithHeadOk(url, method string, heads map[string]string, data map[string]any, timeout time.Duration) ([]byte, error) {
	if _, ok := _MethodMap[method]; !ok {
		return nil, errors.New("not supported http method")
	}
	_, bs, err := GetDefaultClient().CallOpOk(context.Background(), data,
		_MethodMap[method](url),
		Timeout(timeout),
		Header(heads),
	)
	return bs, err
}

func HttpReqWithHead(url, method string, heads map[string]string, data map[string]any, timeout time.Duration) ([]byte, int, error) {
	if _, ok := _MethodMap[method]; !ok {
		return nil, 0, errors.New("not supported http method")
	}
	res, bs, err := GetDefaultClient().CallOp(context.Background(), data,
		_MethodMap[method](url),
		Timeout(timeout),
		Header(heads),
	)
	if err != nil {
		return nil, 0, err
	}
	return bs, res.StatusCode, err
}

func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
