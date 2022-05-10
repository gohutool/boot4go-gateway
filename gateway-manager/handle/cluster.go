package handle

import (
	"fmt"
	. "gateway-manager/model"
	. "gateway-manager/util"
	. "github.com/gohutool/boot4go-etcd/client"
	util4go "github.com/gohutool/boot4go-util"
	. "github.com/gohutool/boot4go-util/http"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/satori/go.uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"reflect"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : cluster.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/25 14:48
* 修改历史 : 1. [2022/4/25 14:48] 创建文件 by LongYong
*/

const (
	CLUSTER_DATA_PREFIX       = DB_PREFIX + "/cluster-data/"
	CLUSTER_DATA_CLUSTER_PATH = CLUSTER_DATA_PREFIX + "%s/"

	CLUSTER_BAK_DATA_PREFIX       = DB_PREFIX + "/cluster-data-bak/"
	CLUSTER_BAK_DATA_CLUSTER_PATH = CLUSTER_BAK_DATA_PREFIX + "%s/"
)

type clusterHandler struct {
}

var ClusterHandler = &clusterHandler{}

func (u *clusterHandler) InitRouter(router *routing.Router, routerGroup *routing.RouteGroup) {
	routerGroup.Get("/cluster/list", TokenInterceptorHandler(u.ListCluster))
	routerGroup.Put("/cluster/<cluster-id>", TokenInterceptorHandler(u.SaveCluster))
	routerGroup.Get("/cluster/<cluster-id>", TokenInterceptorHandler(u.GetCluster))
	routerGroup.Delete("/cluster/<cluster-id>", TokenInterceptorHandler(u.DeleteCluster))
}

func (u *clusterHandler) ListCluster(context *routing.Context) error {
	Logger.Debug("%v", "ListCluster")

	clusters := EtcdClient.GetKeyObjectsWithPrefix(CLUSTER_DATA_PREFIX, reflect.TypeOf((*Cluster)(nil)), nil,
		0, 0, ReadTimeout)

	Result.Success(PageResultDataBuilder.New().Data(clusters).Build(), "OK").Response(context.RequestCtx)
	return nil
}

func (u *clusterHandler) GetCluster(context *routing.Context) error {
	Logger.Debug("%v", "GetCluster")

	id := context.Param("cluster-id")

	if len(id) == 0 {
		panic("没有clusterId参数")
	}

	cluster := EtcdClient.KeyObject(ClusterKey(id), reflect.TypeOf((*Cluster)(nil)), ReadTimeout)

	if cluster == nil {
		Result.Fail("没有找到" + id + "对应的集群数据").Response(context.RequestCtx)
		// panic("没有找到" + id + "对应的集群数据")
	} else {
		Result.Success(cluster, "OK").Response(context.RequestCtx)
	}

	return nil
}

func (u *clusterHandler) SaveCluster(context *routing.Context) error {
	Logger.Debug("%v", "SaveCluster")
	id := context.Param("cluster-id")

	if len(id) == 0 {
		panic("没有clusterId参数")
	}

	cluster := &Cluster{}

	//Unmarshal
	if err := JsonBodyUnmarshalObject(context.RequestCtx, cluster); err != nil {
		panic("参数解析出错 " + err.Error())
	}

	if util4go.IsEmpty(string(cluster.ClusterType)) {
		panic("集群类型不能为空")
	}

	if util4go.IsEmpty(cluster.ClusterName) {
		panic("集群名称不能为空")
	}

	//if len(strings.TrimSpace(cluster.LbType)) == 0 {
	//	panic("负载均衡不能为空")
	//}

	//urlParse, err := url.ParseRequestURI(cluster.ClusterUrl)
	//if err != nil {
	//	panic("集群地址格式不正确")
	//}
	//
	//cluster.ClusterUrl = urlParse.Host

	nt := time.Now().Format("2006/1/2 15:04:05")

	if id == "add" || len(id) == 0 {
		id = uuid.Must(uuid.NewV4(), nil).String()
		cluster.Id = id
		cluster.SetTime = nt
		cluster.UpdateTime = ""
	} else {
		o := EtcdClient.KeyObject(ClusterKey(id), reflect.TypeOf((*Cluster)(nil)), ReadTimeout)

		if o == nil {
			panic("没有找到" + id + "对应的集群数据")
		}

		oldCluster := o.(*Cluster)
		cluster.Id = id

		cluster.SetTime = oldCluster.SetTime
		cluster.UpdateTime = nt
	}

	if _, err := EtcdClient.PutValue(ClusterKey(id), cluster, WriteTimeout); err == nil {
		Result.Success(cluster, "保存集群成功").Response(context.RequestCtx)
	} else {
		Logger.Error("保存集群失败 %+v", cluster)
		panic("保存集群失败")
	}

	return nil
}

func (u *clusterHandler) DeleteCluster(context *routing.Context) error {
	Logger.Debug("%v", "DeleteCluster")
	id := context.Param("cluster-id")

	if len(id) == 0 {
		panic("没有clusterId参数")
	}

	o, err := EtcdClient.KeyValue(ClusterKey(id), ReadTimeout)

	if err != nil {
		panic("没有找到" + id + "对应的集群数据")
	}

	if _, err := EtcdClient.BulkOps(func(leaseID clientv3.LeaseID) ([]clientv3.Op, error) {
		return []clientv3.Op{
			clientv3.OpDelete(ClusterKey(id)),
			clientv3.OpPut(ClusterBakKey(id), o, clientv3.WithLease(leaseID)),
		}, nil
	}, BakDataTTL, WriteTimeout); err != nil {
		Logger.Error("删除集群名称失败 %+v", err)
		panic("删除集群失败")
	}
	Result.Success("", "删除集群成功").Response(context.RequestCtx)
	return nil
}

func ClusterKey(clusterId string) string {
	return fmt.Sprintf(CLUSTER_DATA_CLUSTER_PATH, clusterId)
}

func ClusterBakKey(clusterId string) string {
	return fmt.Sprintf(CLUSTER_BAK_DATA_CLUSTER_PATH, clusterId)
}

func (u *clusterHandler) GetClusterCount() int {
	return EtcdClient.CountWithPrefix(CLUSTER_DATA_PREFIX, ReadTimeout)
}

func (u *clusterHandler) GetEndpointCount() int {

	clusters := EtcdClient.GetKeyObjectsWithPrefix(CLUSTER_DATA_PREFIX, reflect.TypeOf((*Cluster)(nil)), nil,
		0, 0, ReadTimeout)

	if clusters == nil || len(clusters) == 0 {
		return 0
	}

	var count int = 0

	for _, cluster := range clusters {
		count += cluster.(*Cluster).NodeCount
	}

	return count
}
