package main

import (
	"fmt"
	. "gateway-manager/model"
	"gateway-manager/util"
	. "github.com/gohutool/boot4go-etcd/client"
	"testing"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : server_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/10 10:32
* 修改历史 : 1. [2022/5/10 10:32] 创建文件 by LongYong
*/

func init() {
	err := EtcdClient.Init([]string{"192.168.56.101:32379"}, "", "", 0)

	if err == nil {
		fmt.Println("Etcd is connect")
	} else {
		panic("Etcd can not connect")
	}
}

func TestGetSummaryLastMetrics(t *testing.T) {

	datas := util.GetVisitSummaryMetrics("192.168.56.1:9000", "1", 12)

	fmt.Printf("%+v", datas)

	time.Sleep(1 * time.Second)
}

func TestGenerateXAxis(t *testing.T) {
	datas := GenerateXAxis(time.Now(), 1, 10)
	fmt.Printf("%+v\n", datas)

	datas = GenerateXAxis(time.Now(), 15, 10)
	fmt.Printf("%+v\n", datas)

	datas = GenerateXAxis(time.Now(), 30, 10)
	fmt.Printf("%+v\n", datas)

	datas = GenerateXAxis(time.Now(), 60, 10)
	fmt.Printf("%+v\n", datas)

	datas = GenerateXAxis(time.Now(), 100, 10)
	fmt.Printf("%+v\n", datas)

	time.Sleep(1 * time.Second)
}

func TestGetSummaryMetrics(t *testing.T) {

	datas := util.GetSummaryMetrics("", 60, 10)

	for _, data := range datas {
		fmt.Printf("%+v %+v\n", data.Point, *data.Value)
	}

	time.Sleep(1 * time.Second)
}

func TestMakeLineChartData(t *testing.T) {

	datas := util.MakeLineChartData("192.168.56.1:9000", 15, 10)

	fmt.Printf("%+v \n", datas)

	time.Sleep(1 * time.Second)
}
