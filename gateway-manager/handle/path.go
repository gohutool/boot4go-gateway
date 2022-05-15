package handle

import (
	"fmt"
	. "gateway-manager/model"
	. "gateway-manager/util"
	. "github.com/gohutool/boot4go-etcd/client"
	. "github.com/gohutool/boot4go-pathmatcher"
	util4go "github.com/gohutool/boot4go-util"
	. "github.com/gohutool/boot4go-util/http"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/satori/go.uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"reflect"
	"strings"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : domain.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/25 14:48
* 修改历史 : 1. [2022/4/25 14:48] 创建文件 by LongYong
*/

type pathHandler struct {
}

var PathHandler = &pathHandler{}

func (u *pathHandler) InitRouter(router *routing.Router, routerGroup *routing.RouteGroup) {
	routerGroup.Get("/path/<domain-id>", TokenInterceptorHandler(u.ListDomainPath))
	routerGroup.Put("/path/<domain-id>/<path-id>", TokenInterceptorHandler(u.SaveDomainPath))
	routerGroup.Get("/path/<domain-id>/<path-id>", TokenInterceptorHandler(u.GetDomainPath))
	routerGroup.Delete("/path/<domain-id>/<path-id>", TokenInterceptorHandler(u.DeleteDomainPath))

	routerGroup.Post("/path/pathpattern", TokenInterceptorHandler(u.TestDomainPathPattern))
	routerGroup.Post("/path/matchpattern", TokenInterceptorHandler(u.TestDomainMatchPattern))
}

func (u *pathHandler) ListDomainPath(context *routing.Context) error {
	Logger.Debug("%v", "ListPath")
	domainId := context.Param("domain-id")

	if len(domainId) == 0 {
		panic("没有PathId参数")
	}

	domainPaths := EtcdClient.GetKeyObjectsWithPrefix(pathDomainPrefixKey(domainId), reflect.TypeOf((*Path)(nil)), nil,
		0, 0, ReadTimeout)

	Result.Success(PageResultDataBuilder.New().Data(domainPaths).Build(), "OK").Response(context.RequestCtx)
	return nil
}

func (u *pathHandler) GetDomainPath(context *routing.Context) error {
	Logger.Debug("%v", "GetDomainPath")

	domainId := context.Param("domain-id")
	pathId := context.Param("path-id")

	if len(domainId) == 0 || len(pathId) == 0 {
		panic("没有PathId参数")
	}

	path := EtcdClient.KeyObject(domainPathKey(domainId, pathId), reflect.TypeOf((*Path)(nil)), ReadTimeout)

	if path == nil {
		Result.Fail("没有找到" + domainId + "/" + pathId + "对应的域名路径").Response(context.RequestCtx)
		// panic("没有找到" + id + "对应的域名数据")
	} else {
		Result.Success(path, "OK").Response(context.RequestCtx)
	}

	return nil
}

func (u *pathHandler) SaveDomainPath(context *routing.Context) error {
	Logger.Debug("%v", "SaveDomainPath")

	domainId := context.Param("domain-id")
	pathId := context.Param("path-id")

	if len(domainId) == 0 {
		panic("没有DomainID参数")
	}

	if len(strings.TrimSpace(pathId)) == 0 || pathId == "add" {
		pathId = ""
	}

	domainObj := EtcdClient.KeyObject(DomainKey(domainId), reflect.TypeOf((*Domain)(nil)), ReadTimeout)

	if domainObj == nil {
		panic("没有找到" + domainId + "对应的域名数据")
	}

	path := &Path{}

	//Unmarshal
	if err := JsonBodyUnmarshalObject(context.RequestCtx, path); err != nil {
		panic("参数解析出错 " + err.Error())
	}

	path.DomainId = domainId

	if len(strings.TrimSpace(path.ReqPath)) == 0 {
		panic("域名地址不能为空")
	}

	//if len(strings.TrimSpace(path.ReplacePath)) == 0 {
	//	panic("域名转换格式不能为空")
	//}

	nt := time.Now().Format("2006/1/2 15:04:05")

	if len(pathId) == 0 {
		pathId = uuid.Must(uuid.NewV4(), nil).String()
		path.Id = pathId
		path.SetTime = nt
		path.UpdateTime = ""
	} else {
		o := EtcdClient.KeyObject(domainPathKey(domainId, pathId), reflect.TypeOf((*Path)(nil)), ReadTimeout)

		if o == nil {
			panic("没有找到" + domainId + "/" + pathId + "对应的域名路径数据")
		}

		oldPath := o.(*Path)
		path.Id = pathId
		path.SetTime = oldPath.SetTime
		path.UpdateTime = nt
	}

	path.DomainUrl = domainObj.(*Domain).DomainUrl

	if _, err := EtcdClient.PutValue(domainPathKey(domainId, pathId), path, WriteTimeout); err == nil {
		Result.Success(path, "保存域名路径成功").Response(context.RequestCtx)
	} else {
		Logger.Error("保存域名路径失败 %+v", path)
		panic("保存域名路径失败")
	}

	return nil
}

func (u *pathHandler) DeleteDomainPath(context *routing.Context) error {
	Logger.Debug("%v", "DeleteDomain")
	domainId := context.Param("domain-id")
	pathId := context.Param("path-id")

	if len(domainId) == 0 {
		panic("没有domainId参数")
	}

	if len(pathId) == 0 {
		panic("没有pathId参数")
	}

	o, err := EtcdClient.KeyValue(domainPathKey(domainId, pathId), ReadTimeout)

	if err != nil {
		panic("没有找到" + domainId + "/" + pathId + "对应的域名路径数据")
	}

	if _, err := EtcdClient.BulkOps(func(leaseID clientv3.LeaseID) ([]clientv3.Op, string, error) {
		return []clientv3.Op{
			clientv3.OpDelete(domainPathKey(domainId, pathId)),
			clientv3.OpPut(domainPathBakKey(domainId, pathId), o, clientv3.WithLease(leaseID)),
		}, domainPathBakKey(domainId, pathId), nil
	}, BakDataTTL, WriteTimeout); err != nil {
		Logger.Error("删除域名名称失败 %+v", err)
		panic("删除域名失败")
	}
	Result.Success("", "删除域名成功").Response(context.RequestCtx)
	return nil
}

func domainPathKey(domainId, pathId string) string {
	return fmt.Sprintf(DOMAIN_PATH_DATA_PATH, domainId, pathId)
}

func domainPathBakKey(domainId, pathId string) string {
	return fmt.Sprintf(DOMAIN_BAK_PATH_DATA_PATH, domainId, pathId)
}

func pathDomainPrefixKey(domainId string) string {
	return fmt.Sprintf(DOMAIN_PATH_DATA_PREFIX, domainId)
}

func (u *pathHandler) TestDomainPathPattern(context *routing.Context) error {
	test, _ := Param(context.RequestCtx, "test")
	pattern, _ := Param(context.RequestCtx, "pattern")

	if util4go.IsEmpty(pattern) {
		panic("匹配路径规则不能为空")
	}

	if ok, err := Match(pattern, test); err != nil {
		panic("匹配路径规则解析不正确")
	} else {
		if ok {
			Result.Success(1, "匹配路径规则匹配").Response(context.RequestCtx)
		} else {
			Result.Success(0, "匹配路径规则不匹配").Response(context.RequestCtx)
		}
	}
	return nil
}

func (u *pathHandler) TestDomainMatchPattern(context *routing.Context) error {
	test, _ := Param(context.RequestCtx, "test")
	pattern, _ := Param(context.RequestCtx, "pattern")
	format, _ := Param(context.RequestCtx, "format")

	if util4go.IsEmpty(pattern) {
		Result.Success(test, "匹配路径规则匹配").Response(context.RequestCtx)
	}

	if util4go.IsEmpty(format) {
		Result.Success(format, "匹配路径规则匹配").Response(context.RequestCtx)
	}

	if util4go.IsEmpty(test) {
		Result.Success(test, "匹配路径规则匹配").Response(context.RequestCtx)
	}

	if r, err := util4go.RegExpPool.ConvertRegExpWithFormat(test, pattern, format); err != nil {
		panic("匹配路径规则解析不正确")
	} else {
		Result.Success(r, "匹配路径规则匹配").Response(context.RequestCtx)
	}
	return nil
}

func (u *pathHandler) GetPathCount() int {
	return EtcdClient.CountWithPrefix(DB_PREFIX+"/path-data/", ReadTimeout)
}
