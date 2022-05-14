package util

import (
	prometheusfasthttp "github.com/gohutool/boot4go-prometheus/fasthttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : prometheus.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/28 20:52
* 修改历史 : 1. [2022/4/28 20:52] 创建文件 by LongYong
*/

var (
	/**/
	TotalCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "go",
			Name:      "forward_total",
			Help:      "Total number of host forward",
		},
		[]string{"host"},
	)

	//m = metrics.NewMeter()

	// 监控实时并发量（处理中的请求）
	ConcurrentRequestsGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Subsystem: "go",
			Name:      "request_handle_concurrent",
			Help:      "Number of incoming HTTP Requests handling concurrently now.",
		},
	)

	// 监控请求量，请求耗时等
	LatencyRequestsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: "go",
			Name:      "request_handle_latency",
			Help:      "Histogram statistics of http(s) requests latency second",
			Buckets:   []float64{0.01, 0.1, 0.2, 0.5, 1, 2, 5, 10, 30},
		},
		[]string{"code"},
	)
)

func init() {
	//prometheus.MustRegister(totalCounterVec)
	prometheus.MustRegister(ConcurrentRequestsGauge)
	prometheus.MustRegister(LatencyRequestsHistogram)
	prometheus.MustRegister(TotalCounterVec)
}

func PrometheusRequestHandler(next fasthttp.RequestHandler) fasthttp.RequestHandler {

	return prometheusfasthttp.RequestCounterHandler(func(ctx *fasthttp.RequestCtx) {

		host := string(ctx.Host())

		if string(ctx.Path()) == "/metrics" {
			prometheusfasthttp.PrometheusHandler(prometheusfasthttp.HandlerOpts{})(ctx)

			return
		} else {
			if next != nil {
				TotalCounterVec.WithLabelValues(host).Inc()
				startTime := time.Now()

				ConcurrentRequestsGauge.Inc()
				defer ConcurrentRequestsGauge.Dec()

				next(ctx)

				finishTime := time.Now()
				LatencyRequestsHistogram.WithLabelValues(
					strconv.Itoa(ctx.Response.StatusCode())).Observe(float64(finishTime.Sub(startTime) / time.Second))
			}
		}

	})

}
