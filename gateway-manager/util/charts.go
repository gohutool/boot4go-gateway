package util

import (
	"gateway-manager/model"
	util4go "github.com/gohutool/boot4go-util"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : charts.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/10 16:10
* 修改历史 : 1. [2022/5/10 16:10] 创建文件 by LongYong
*/

func Reverse[T any](source []T) []T {
	var rtn []T

	for idx := len(source) - 1; idx >= 0; idx-- {
		rtn = append(rtn, source[idx])
	}

	return rtn
}

func MakeLineChartData(sername string, dimension, count int) map[string]any {
	rtn := make(map[string]any)
	datas := Reverse(GetSummaryMetrics(sername, dimension, count))

	rtn["xAxis"] = util4go.Collect(datas, func(one model.MetricsPoint) string {
		return one.Point.Text
	})

	rtn["current_blackip"] = util4go.Collect(datas, func(one model.MetricsPoint) int64 {
		return one.Value.Current.BlackIpLose
	})

	rtn["current_total"] = util4go.Collect(datas, func(one model.MetricsPoint) int64 {
		return one.Value.Current.Total
	})

	rtn["current_ok"] = util4go.Collect(datas, func(one model.MetricsPoint) int64 {
		return one.Value.Current.Ok
	})

	rtn["current_lose"] = util4go.Collect(datas, func(one model.MetricsPoint) int64 {
		return one.Value.Current.Lose
	})

	rtn["increase_total"] = util4go.Collect(datas, func(one model.MetricsPoint) int64 {
		return one.Value.Increase.Total
	})

	rtn["increase_ok"] = util4go.Collect(datas, func(one model.MetricsPoint) int64 {
		return one.Value.Increase.Ok
	})

	rtn["increase_lose"] = util4go.Collect(datas, func(one model.MetricsPoint) int64 {
		return one.Value.Increase.Lose
	})

	rtn["increase_black_ip"] = util4go.Collect(datas, func(one model.MetricsPoint) int64 {
		return one.Value.Increase.BlackIpLose
	})

	rtn["data"] = datas

	return rtn
}
