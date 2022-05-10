package handle

import (
	"fmt"
	. "gateway-manager/model"
	. "gateway-manager/util"
	. "github.com/gohutool/boot4go-etcd/client"
	. "github.com/gohutool/boot4go-util"
	. "github.com/gohutool/boot4go-util/http"
	. "github.com/gohutool/boot4go-util/jwt"
	routing "github.com/qiangxue/fasthttp-routing"
	"reflect"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : user.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/24 10:06
* 修改历史 : 1. [2022/4/24 10:06] 创建文件 by LongYong
*/

type userHandler struct {
}

var UserHandler = &userHandler{}

func (u *userHandler) InitRouter(router *routing.Router, routerGroup *routing.RouteGroup) {
	router.Post("/login", u.Login)
	router.Post("/init", u.Init)
	router.Get("/logout", u.Logout)

	routerGroup.Get("/user/menu", TokenInterceptorHandler(u.MenuTree))
	routerGroup.Put("/user/pwd", TokenInterceptorHandler(u.MenuTree))
	routerGroup.Put("/user/profile", TokenInterceptorHandler(u.MenuTree))
	routerGroup.Get("/user/portal", TokenInterceptorHandler(u.Portal))
	routerGroup.Get("/portal/linereport", TokenInterceptorHandler(u.LineMetric))
	routerGroup.Get("/portal/barreport", TokenInterceptorHandler(u.BarMetric))

}

func (u *userHandler) Login(context *routing.Context) error {
	username := GetParams(context.RequestCtx, "username", "")
	password := GetParams(context.RequestCtx, "password", "")

	if username == "" || password == "" {
		Result.Fail("请填写登录用户名和用户密码").Response(context.RequestCtx)
		return nil
	}

	username = MD5(username)

	o := EtcdClient.KeyObject(userKey(username), reflect.TypeOf((*AdminUser)(nil)), ReadTimeout)

	if o == nil {
		Result.Fail("登录用户名和用户密码不正确").Response(context.RequestCtx)
		return nil
	}

	adminUser := o.(*AdminUser)

	if adminUser.Password != SaltMd5(password, adminUser.Salt) {
		Result.Fail("登录用户名和用户密码不正确").Response(context.RequestCtx)
		return nil
	}

	token := GenToken(adminUser.UserId, Issuer, Issuer, TokenExpire)

	Result.Success(token, "").Response(context.RequestCtx)

	Logger.Debug("%v %v %v", adminUser.UserId, username, password)
	return nil
}

func (u *userHandler) Logout(context *routing.Context) error {
	Logger.Debug("%v", "Logout")
	Result.Success("", "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) Init(context *routing.Context) error {
	Logger.Debug("%v", "Init")
	return nil
}

func (u *userHandler) MenuTree(context *routing.Context) error {
	Logger.Debug("%v", "MenuTree")

	menus := GetMenuObj()

	Result.Success(menus["datas"], "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) SavePwd(context *routing.Context) error {
	Logger.Debug("%v", "SavePwd")

	Result.Success("", "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) SaveProfile(context *routing.Context) error {
	Logger.Debug("%v", "SaveProfile")

	Result.Success("", "OK").Response(context.RequestCtx)
	return nil
}

func userKey(username string) string {
	return fmt.Sprintf(ADMIN_USER_DATA_PATH, username)
}

func (u *userHandler) Portal(context *routing.Context) error {
	Logger.Debug("%v", "Portal")

	domainCount := DomainHandler.GetDomainCount()
	pathCount := PathHandler.GetPathCount()
	serverCount := GatewayHandler.GetGatewayCount()
	nodeCount := GatewayHandler.GetGatewayNodeCount()
	certCount := CertHandler.GetCertCount()
	clusterCount := ClusterHandler.GetClusterCount()
	endpointCount := ClusterHandler.GetEndpointCount()

	visitMetrics := GetMetrics("")

	res := make(map[string]any)

	res["domainCount"] = domainCount
	res["pathCount"] = pathCount
	res["serverCount"] = serverCount
	res["certCount"] = certCount
	res["clusterCount"] = clusterCount
	res["endpointCount"] = endpointCount
	res["nodeCount"] = nodeCount
	res["nodeVisit"] = visitMetrics.Total
	res["certTotal"] = visitMetrics.CertTotal

	var domains []*Domain
	domains = CopyArray(GetDomains(), domains)

	var totalCluster, totalEndpoint int

	for _, domain := range domains {
		if domain != nil {
			totalEndpoint += len(domain.Targets)
		}
	}

	var paths []*Path
	paths = CopyArray(GetDomainPaths(), paths)

	for _, p := range paths {

		if p != nil {
			totalEndpoint += len(p.Targets)
			var haveNode = false

			for _, t := range p.Targets {
				if t.PointerType != "Node" {
					totalCluster++
				} else {
					haveNode = true
				}
			}

			if haveNode {
				totalCluster++
			}
		}

	}

	res["clusterCount"] = totalCluster
	res["endpointCount"] = totalEndpoint

	res["visitMetrics"] = visitMetrics

	Result.Success(res, "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) LineMetric(context *routing.Context) error {
	res := MakeLineChartData("", 15, 25)

	Result.Success(res, "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) BarMetric(context *routing.Context) error {
	res := MakeLineChartData("", 60, 24)

	Result.Success(res, "OK").Response(context.RequestCtx)
	return nil
}
