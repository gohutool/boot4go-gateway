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
	"strings"
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
	routerGroup.Get("/user/list", TokenInterceptorHandler(u.ListUser))
	routerGroup.Put("/user/pwd", TokenInterceptorHandler(u.SavePwd))
	routerGroup.Put("/user/<user-id>", TokenInterceptorHandler(u.SaveUser))
	routerGroup.Get("/user/<user-id>", TokenInterceptorHandler(u.GetUser))
	routerGroup.Delete("/user/<user-id>", TokenInterceptorHandler(u.DelUser))
	routerGroup.Put("/user/profile", TokenInterceptorHandler(u.SaveProfile))
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

func (u *userHandler) GetUser(context *routing.Context) error {
	id := context.Param("user-id")

	if IsEmpty(id) {
		panic("用户信息不正确")
	}

	o := EtcdClient.KeyObject(userKey(id), reflect.TypeOf((*AdminUser)(nil)), ReadTimeout)

	if o == nil {
		panic("用户" + id + "不存在")
	}

	Result.Success(o, "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) SaveUser(context *routing.Context) error {

	id := context.Param("user-id")

	username := GetParams(context.RequestCtx, "username", "")
	password := GetParams(context.RequestCtx, "password", "")
	password2 := GetParams(context.RequestCtx, "password2", "")

	if IsEmpty(username) || IsEmpty(password) {
		Result.Fail("请填写登录用户名和用户密码").Response(context.RequestCtx)
		return nil
	}

	if password2 != password {
		Result.Fail("密码不一致").Response(context.RequestCtx)
		return nil
	}

	username = strings.TrimSpace(username)

	userid := GetUserId(context)

	if IsEmpty(userid) {
		panic("登录信息不正确")
	}

	o := EtcdClient.KeyObject(userKey(userid), reflect.TypeOf((*AdminUser)(nil)), ReadTimeout)

	if o == nil {
		panic("登录信息不正确")
	}

	user := o.(*AdminUser)

	if user.UserName != "ginghan" {
		panic("用户权限不正确")
	}

	if !IsEmpty(id) {
		var oUser *AdminUser
		o := EtcdClient.KeyObject(userKey(id), reflect.TypeOf((*AdminUser)(nil)), ReadTimeout)
		if o == nil {
			panic("用户" + id + "不存在")
		}
		oUser = o.(*AdminUser)

		if oUser.UserName != username {
			o = EtcdClient.KeyObject(userKey(MD5(username)), reflect.TypeOf((*AdminUser)(nil)), ReadTimeout)
			if o != nil {
				panic("用户" + username + "已经存在")
			} else {
				EtcdClient.Delete(userKey(id), 0)
			}
		}
	} else {
		o := EtcdClient.KeyObject(userKey(MD5(username)), reflect.TypeOf((*AdminUser)(nil)), ReadTimeout)

		if o != nil {
			panic("用户" + username + "已经存在")
		}
	}

	err := CreateAdmin(username, password)

	if err != nil {
		panic("创建用户失败:" + err.Error())
	}

	Result.Success("", "OK").Response(context.RequestCtx)

	return nil
}

func (u *userHandler) DelUser(context *routing.Context) error {
	userid := GetUserId(context)

	if IsEmpty(userid) {
		panic("登录信息不正确")
	}

	o := EtcdClient.KeyObject(userKey(userid), reflect.TypeOf((*AdminUser)(nil)), ReadTimeout)

	if o == nil {
		panic("登录信息不正确")
	}

	user := o.(*AdminUser)

	if user.UserName != "ginghan" {
		panic("用户权限不正确")
	}

	id := context.Param("user-id")
	if IsEmpty(id) {
		panic("用户信息不存在")
	}

	o = EtcdClient.KeyObject(userKey(id), reflect.TypeOf((*AdminUser)(nil)), ReadTimeout)

	if o == nil {
		panic("用户" + id + "不存在")
	}

	user = o.(*AdminUser)

	if user.UserName == "ginghan" {
		panic("系统用户不能删除")
	}

	EtcdClient.Delete(userKey(id), 0)

	Result.Success("", "OK").Response(context.RequestCtx)

	return nil
}

func (u *userHandler) SavePwd(context *routing.Context) error {
	Logger.Debug("%v", "SavePwd")
	password := string(context.FormValue("password"))

	userid := GetUserId(context)

	if IsEmpty(userid) {
		panic("登录信息不正确")
	}
	if IsEmpty(password) {
		panic("登录密码不能为空")
	}

	o := EtcdClient.KeyObject(userKey(userid), reflect.TypeOf((*AdminUser)(nil)), ReadTimeout)

	if o == nil {
		panic("登录信息不正确")
	}

	user := o.(*AdminUser)
	user.Password = SaltMd5(password, user.Salt)

	_, err := EtcdClient.PutValue(fmt.Sprintf(ADMIN_USER_DATA_PATH, userid), user, WriteTimeout)

	if err != nil {
		panic("密码修改失败:" + err.Error())
	}

	Result.Success("", "OK").Response(context.RequestCtx)
	return nil
}

func (u *userHandler) ListUser(context *routing.Context) error {
	Logger.Debug("%v", "SavePwd")

	Result.Success(PageResultDataBuilder.New().Data(GetUsers()).Build(), "OK").Response(context.RequestCtx)
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
