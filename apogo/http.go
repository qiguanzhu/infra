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
 @Time    : 2024/10/12 -- 16:52
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: request.go
*/

package apogo

import (
	"context"
	"github.com/Bishoptylaor/go-toolkit/xnet/xhttp"
	"net/http"
	"time"
)

func NewHclient(duration time.Duration) xhttp.HttpClientWrapper {
	return xhttp.NewHttpClientWrapper(&http.Client{
		Timeout: duration,
	})
}

func Request(hclient xhttp.HttpClientWrapper, url string) ([]byte, error) {
	_, bs, err := hclient.CallOpOk(context.TODO(), map[string]any{}, xhttp.Get(url))
	return bs, err
}
