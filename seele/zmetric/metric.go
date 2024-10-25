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
 @Time    : 2024/10/12 -- 11:34
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: metric.go
*/

package zmetric

// Counter describes a metric that accumulates values monotonically.
// An example of a counter is the number of received rpc requests.
type Counter interface {
	With(labelValues ...string) Counter
	Add(delta float64)
	Inc()
}

// Gauge describes a metric that takes specific values over time.
// An example of a gauge is the current number of connections.
type Gauge interface {
	With(labelValues ...string) Gauge
	Set(value float64)
	Add(delta float64)
}

// Histogram describes a metric that takes repeated observations of the same
// kind of thing, and produces a statistical summary of those observations,
// typically expressed as quantiles or buckets. An example of a histogram is
// rpc request latencies.
type Histogram interface {
	With(labelValues ...string) Histogram
	Observe(value float64)
}

// VectorOpts contains the common arguments for creating vec Metric..
type VectorOpts struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	LabelNames []string
}
