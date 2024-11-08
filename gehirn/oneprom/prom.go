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
 @Time    : 2024/11/1 -- 18:40
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: prom.go
*/

package oneprom

import (
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qiguanzhu/infra/nerv/magi/xenv"
	"github.com/spf13/cast"
)

type prom struct {
	app    string
	env    string
	uptime *prometheus.CounterVec
	statsd *prometheus.CounterVec
	timing *prometheus.HistogramVec

	serverStartedCounter   *prometheus.CounterVec
	serverHandledCounter   *prometheus.CounterVec
	serverHandledHistogram *prometheus.HistogramVec
}

func (p *prom) Statsd(target, key string, add int) {
	if add <= 0 {
		return
	}
	p.statsd.WithLabelValues(p.app, p.env, target, key).Add(cast.ToFloat64(add))
}

func (p *prom) Timing(target, key string, cost float64) {
	p.timing.WithLabelValues(p.app, p.env, target, key).Observe(cost)
}

var (
	_once       sync.Once
	defaultProm prom
)

func Init(cfg xenv.Enver) {
	_once.Do(func() {
		uptime := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "server_uptime_total",
				Help: "service uptime.",
			}, []string{"app", "env"},
		)
		statsd := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "server_call_total",
		}, []string{"app", "env", "target", "key"})
		timing := prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "server_call_seconds",
			Buckets: []float64{1, 10, 20, 30, 50, 80, 100, 200, 300, 500, 800, 1000, 3000, 5000, 10000},
		}, []string{"app", "env", "target", "key"})

		prometheus.MustRegister(uptime)
		prometheus.MustRegister(statsd)
		prometheus.MustRegister(timing)

		serverStartedCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "server_started_total",
		}, []string{"app", "env", "handler"})
		serverHandledCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "server_handled_total",
		}, []string{"app", "env", "handler", "code"})
		serverHandledHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "server_handling_seconds",
			Buckets: []float64{1, 10, 20, 30, 50, 80, 100, 200, 300, 500, 800, 1000, 3000, 5000, 10000},
		}, []string{"app", "env", "handler"})

		prometheus.MustRegister(serverStartedCounter)
		prometheus.MustRegister(serverHandledCounter)
		prometheus.MustRegister(serverHandledHistogram)

		go func() {
			for range time.Tick(time.Second) {
				uptime.WithLabelValues(cfg.GetApp(), cfg.GetEnv()).Inc()
			}
		}()

		defaultProm = prom{
			app:                    cfg.GetApp(),
			env:                    cfg.GetApp(),
			uptime:                 uptime,
			statsd:                 statsd,
			timing:                 timing,
			serverStartedCounter:   serverStartedCounter,
			serverHandledCounter:   serverHandledCounter,
			serverHandledHistogram: serverHandledHistogram,
		}

		http.Handle("/metrics", promhttp.Handler())
	})
}

func Statsd(target, key string, add int) {
	if add <= 0 {
		return
	}
	defaultProm.statsd.WithLabelValues(defaultProm.app, defaultProm.env, target, key).Add(cast.ToFloat64(add))
}

func Timing(target, key string, cost float64) {
	defaultProm.timing.WithLabelValues(defaultProm.app, defaultProm.env, target, key).Observe(cost)
}

func ServerStartedCounter(handlerName string) {
	defaultProm.serverStartedCounter.WithLabelValues(defaultProm.app, defaultProm.env, handlerName).Inc()
}

func ServerHandledCounter(handlerName string, code int) {
	defaultProm.serverHandledCounter.WithLabelValues(defaultProm.app, defaultProm.env, handlerName, cast.ToString(code)).Inc()
}

func ServerHandledHistogram(handlerName string, cost time.Duration) {
	defaultProm.serverHandledHistogram.WithLabelValues(defaultProm.app, defaultProm.env, handlerName).Observe(cast.ToFloat64(cost.Milliseconds()))
}
