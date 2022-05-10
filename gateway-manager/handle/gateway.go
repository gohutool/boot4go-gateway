package handle

import (
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
	"reflect"

	. "gateway-manager/model"
	. "gateway-manager/util"
	. "github.com/gohutool/boot4go-etcd/client"
	. "github.com/gohutool/boot4go-util/http"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : gateway.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/30 23:19
* 修改历史 : 1. [2022/4/30 23:19] 创建文件 by LongYong
*/

const (
	GATEWAY_ACTIVE_PREFIX = DB_PREFIX + "/gateway-active/"

	GATEWAY_ACTIVE_FORMAT        = GATEWAY_ACTIVE_PREFIX + "%s/"
	GATEWAY_ACTIVE_SEARCH_FORMAT = GATEWAY_ACTIVE_PREFIX + "%s"
)

type gatewayHandler struct {
}

var GatewayHandler = &gatewayHandler{}

func (u *gatewayHandler) InitRouter(router *routing.Router, routerGroup *routing.RouteGroup) {
	routerGroup.Get("/gateway/list", TokenInterceptorHandler(u.ListGateway))
	//routerGroup.Put("/cert/<cert-id>", TokenInterceptorHandler(u.SaveCert))
	//routerGroup.Get("/cert/<cert-id>", TokenInterceptorHandler(u.GetCert))
	//routerGroup.Delete("/cert/<cert-id>", TokenInterceptorHandler(u.DeleteCert))
}

func (u *gatewayHandler) ListGateway(context *routing.Context) error {
	Logger.Debug("%v", "ListGateway")

	prefix, _ := Param(context.RequestCtx, "prefix")

	gateways := EtcdClient.GetKeyObjectsWithPrefix(SearchPrefix(prefix), reflect.TypeOf((*Gateway)(nil)), nil,
		0, 0, ReadTimeout)

	Result.Success(PageResultDataBuilder.New().Data(gateways).Build(), "OK").Response(context.RequestCtx)

	return nil
}

func SearchPrefix(prefix string) string {
	return fmt.Sprintf(GATEWAY_ACTIVE_SEARCH_FORMAT, prefix)
}

func (u *gatewayHandler) GetGatewayCount() int {
	return EtcdClient.CountWithPrefix(GATEWAY_ACTIVE_PREFIX, ReadTimeout)
}

func (u *gatewayHandler) GetGatewayNodeCount() int {
	return EtcdClient.CountWithPrefix(GATEWAY_ACTIVE_PREFIX, ReadTimeout)
}
