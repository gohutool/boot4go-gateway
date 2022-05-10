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
	"strconv"
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

type domainHandler struct {
}

var DomainHandler = &domainHandler{}

func (u *domainHandler) InitRouter(router *routing.Router, routerGroup *routing.RouteGroup) {
	routerGroup.Get("/domain/list", TokenInterceptorHandler(u.ListDomain))
	routerGroup.Put("/domain/<domain-id>", TokenInterceptorHandler(u.SaveDomain))
	routerGroup.Put("/domain/url/<domain-id>", TokenInterceptorHandler(u.SaveDomainUrl))
	routerGroup.Get("/domain/<domain-id>", TokenInterceptorHandler(u.GetDomain))
	routerGroup.Delete("/domain/<domain-id>", TokenInterceptorHandler(u.DeleteDomain))
}

func (u *domainHandler) ListDomain(context *routing.Context) error {
	Logger.Debug("%v", "ListDomain")

	domains := GetDomains()

	Result.Success(PageResultDataBuilder.New().Data(domains).Build(), "OK").Response(context.RequestCtx)
	return nil
}

func (u *domainHandler) GetDomain(context *routing.Context) error {
	Logger.Debug("%v", "GetDomain")

	id := context.Param("domain-id")

	if len(id) == 0 {
		panic("没有domainId参数")
	}

	domain := EtcdClient.KeyObject(DomainKey(id), reflect.TypeOf((*Domain)(nil)), ReadTimeout)

	if domain == nil {
		Result.Fail("没有找到" + id + "对应的域名数据").Response(context.RequestCtx)
		// panic("没有找到" + id + "对应的域名数据")
	} else {
		Result.Success(domain, "OK").Response(context.RequestCtx)
	}

	return nil
}

func (u *domainHandler) SaveDomain(context *routing.Context) error {
	Logger.Debug("%v", "SaveDomain")
	id := context.Param("domain-id")

	if len(id) == 0 {
		panic("没有domainId参数")
	}

	domain := &Domain{}

	//Unmarshal
	if err := JsonBodyUnmarshalObject(context.RequestCtx, domain); err != nil {
		panic("参数解析出错 " + err.Error())
	}

	if len(strings.TrimSpace(domain.DomainUrl)) == 0 {
		panic("域名地址不能为空")
	}

	if len(strings.TrimSpace(domain.DomainName)) == 0 {
		panic("域名名称不能为空")
	}

	//if len(strings.TrimSpace(domain.LbType)) == 0 {
	//	panic("负载均衡不能为空")
	//}

	//urlParse, err := url.ParseRequestURI(domain.DomainUrl)
	//if err != nil {
	//	panic("域名地址格式不正确")
	//}
	//
	//domain.DomainUrl = urlParse.Host

	nt := time.Now().Format("2006/1/2 15:04:05")

	if id == "add" || len(id) == 0 {
		id = uuid.Must(uuid.NewV4(), nil).String()
		domain.Id = id
		domain.SetTime = nt
		domain.UpdateTime = ""
	} else {
		o := EtcdClient.KeyObject(DomainKey(id), reflect.TypeOf((*Domain)(nil)), ReadTimeout)

		if o == nil {
			panic("没有找到" + id + "对应的域名数据")
		}

		oldDomain := o.(*Domain)
		domain.Id = id

		domain.SslOn = oldDomain.SslOn
		domain.SslPort = oldDomain.SslPort

		domain.SetTime = oldDomain.SetTime
		domain.UpdateTime = nt
	}

	if _, err := EtcdClient.PutValue(DomainKey(id), domain, WriteTimeout); err == nil {
		Result.Success(domain, "保存域名成功").Response(context.RequestCtx)
	} else {
		Logger.Error("保存域名失败 %+v", domain)
		panic("保存域名失败")
	}

	return nil
}

func (u *domainHandler) SaveDomainUrl(context *routing.Context) error {
	Logger.Debug("%v", "SaveDomainUrl")
	id := context.Param("domain-id")
	domainName := string(context.FormValue("domain_url"))
	sslOn := string(context.FormValue("ssl_on"))
	sslPort := string(context.FormValue("ssl_port"))

	if len(id) == 0 {
		panic("没有domainId参数")
	}

	if len(strings.TrimSpace(domainName)) == 0 {
		panic("域名名称不能为空")
	}

	o := EtcdClient.KeyObject(DomainKey(id), reflect.TypeOf((*Domain)(nil)), ReadTimeout)

	if o == nil {
		panic("没有找到" + id + "对应的域名数据")
	}

	oldDomain := o.(*Domain)
	oldDomain.DomainUrl = domainName

	if sslOn == "1" {
		oldDomain.SslOn = true

		if util4go.IsEmpty(sslPort) {
			oldDomain.SslPort = 9443
		} else {
			if p, err := strconv.Atoi(sslPort); err == nil {
				oldDomain.SslPort = p
			} else {
				panic("SSLPort parse error : " + err.Error())
			}
		}
	} else {
		oldDomain.SslOn = false
	}

	oldDomain.UpdateTime = time.Now().Format("2006/1/2 15:04:05")

	if _, err := EtcdClient.PutValue(DomainKey(id), oldDomain, WriteTimeout); err == nil {
		Result.Success(oldDomain, "保存域名地址成功").Response(context.RequestCtx)
	} else {
		Logger.Error("保存域名地址失败 %+v", oldDomain)
		panic("保存域名地址失败")
	}

	return nil
}

func (u *domainHandler) DeleteDomain(context *routing.Context) error {
	Logger.Debug("%v", "DeleteDomain")
	id := context.Param("domain-id")

	if len(id) == 0 {
		panic("没有domainId参数")
	}

	o, err := EtcdClient.KeyValue(DomainKey(id), ReadTimeout)

	if err != nil {
		panic("没有找到" + id + "对应的域名数据")
	}

	if _, err := EtcdClient.BulkOps(func(leaseID clientv3.LeaseID) ([]clientv3.Op, error) {
		return []clientv3.Op{
			clientv3.OpDelete(DomainKey(id)),
			clientv3.OpPut(DomainBakKey(id), o, clientv3.WithLease(leaseID)),
		}, nil
	}, BakDataTTL, WriteTimeout); err != nil {
		Logger.Error("删除域名名称失败 %+v", err)
		panic("删除域名失败")
	}
	Result.Success("", "删除域名成功").Response(context.RequestCtx)
	return nil
}

func DomainKey(domainId string) string {
	return fmt.Sprintf(DOMAIN_DATA_DOMAIN_PATH, domainId)
}

func DomainBakKey(domainId string) string {
	return fmt.Sprintf(DOMAIN_BAK_DATA_DOMAIN_PATH, domainId)
}

func (u *domainHandler) GetDomainCount() int {
	return EtcdClient.CountWithPrefix(DOMAIN_DATA_PREFIX, ReadTimeout)
}
