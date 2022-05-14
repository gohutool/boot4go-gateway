package main

import (
	"fmt"
	"gateway-server/util"
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

	LM15 := util.GetSummaryLastMetrics("192.168.56.1:9000", "15", "2022/05/10 10:1/4")

	fmt.Printf("%+v\n", LM15)

	time.Sleep(1 * time.Second)
}

func TestBuildSummaryLastMetrics(t *testing.T) {

	overall := util.GetMetrics("")
	t1 := time.Now()
	dimension := util.Time2Dimension(t1)
	OM1, OM15, OM30, OM60, OM1D := util.BuildSummaryMetrics(dimension, "", overall)

	fmt.Printf("%+v\n", OM1)
	fmt.Printf("%+v\n", OM15)
	fmt.Printf("%+v\n", OM30)
	fmt.Printf("%+v\n", OM60)
	fmt.Printf("%+v\n", OM1D)

	time.Sleep(1 * time.Second)
}
