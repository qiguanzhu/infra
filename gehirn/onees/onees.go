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
 @Time    : 2024/11/4 -- 18:26
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: onees.go
*/

package onees

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func Connect(dsn string, user, pwd string) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL(dsn),
		elastic.SetBasicAuth(user, pwd),
		elastic.SetHealthcheck(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetGzip(true),
		elastic.SetSniff(false),
		elastic.SetHttpClient(&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConnsPerHost:   1024,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 10 * time.Second,
		}),
	)
}

func MustConnect(dsn string, user, pwd string) *elastic.Client {
	cli, err := Connect(dsn, user, pwd)
	if err != nil {
		panic(fmt.Sprintf("connect to es %s error %+v", dsn, err))
	}
	return cli
}
