package handle

import (
	"crypto/tls"
	"fmt"
	. "gateway-manager/model"
	. "gateway-manager/util"
	. "github.com/gohutool/boot4go-etcd/client"
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

const (
	CERT_DATA_PREFIX = DB_PREFIX + "/cert-data/"
	CERT_DATA_FORMAT = CERT_DATA_PREFIX + "%s"

	CERT_BAK_DATA_PREFIX = DB_PREFIX + "/cert-data-bak/"
	CERT_BAK_DATA_FORMAT = CERT_BAK_DATA_PREFIX + "%s"
)

type certHandler struct {
}

var CertHandler = &certHandler{}

func (u *certHandler) InitRouter(router *routing.Router, routerGroup *routing.RouteGroup) {
	routerGroup.Get("/cert/list", TokenInterceptorHandler(u.ListCert))
	routerGroup.Put("/cert/<cert-id>", TokenInterceptorHandler(u.SaveCert))
	routerGroup.Get("/cert/<cert-id>", TokenInterceptorHandler(u.GetCert))
	routerGroup.Delete("/cert/<cert-id>", TokenInterceptorHandler(u.DeleteCert))

	router.Get("/cert/testlease/<cert-id>", u.TestLease)
}

func (u *certHandler) TestLease(context *routing.Context) error {

	certId := context.Param("cert-id")
	content, _ := Param(context.RequestCtx, "content")

	EtcdClient.PutKeepAliveValue("/test/release/"+certId, content, 30, 10,
		func(leaseID clientv3.LeaseID) error {
			Logger.Error("收到自动续租退出应答(%v)", leaseID)
			return nil
		})

	/*
		EtcdClient.BulkOps(func(leaseID clientv3.LeaseID) ([]clientv3.Op, error) {
			if keepRespChan, err := EtcdClient.Get().KeepAlive(c.Background(), leaseID); err == nil {
				go func() {
					for {
						select {
						case keepResp := <-keepRespChan:
							if keepRespChan == nil {
								Logger.Error("租约已经失效, 退出自动续约")
								return
							} else { //每秒会续租一次，所以就会受到一次应答
								Logger.Error("收到自动续租应答:%v", keepResp.ID)
							}
						}
					}
				}()
			} else {
				Logger.Error("lease keepalive error", err.Error())
			}

			return []clientv3.Op{
				clientv3.OpPut("/test/release/"+certId, content, clientv3.WithLease(leaseID)),
			}, nil

		}, 1, 1)
	*/
	Result.Success(content, "OK").Response(context.RequestCtx)
	return nil
}

func (u *certHandler) ListCert(context *routing.Context) error {
	Logger.Debug("%v", "ListCert")

	certs := EtcdClient.GetKeyObjectsWithPrefix(CERT_DATA_PREFIX, reflect.TypeOf((*Cert)(nil)), nil,
		0, 0, ReadTimeout)

	Result.Success(PageResultDataBuilder.New().Data(certs).Build(), "OK").Response(context.RequestCtx)
	return nil
}

func (u *certHandler) GetCert(context *routing.Context) error {
	Logger.Debug("%v", "GetCert")

	certId := context.Param("cert-id")

	if len(certId) == 0 {
		panic("没有CertId参数")
	}

	cert := EtcdClient.KeyObject(CertKey(certId), reflect.TypeOf((*Cert)(nil)), ReadTimeout)

	if cert == nil {
		Result.Fail("没有找到" + certId + "对应的证书").Response(context.RequestCtx)
	} else {
		Result.Success(cert, "OK").Response(context.RequestCtx)
	}

	return nil
}

func (u *certHandler) SaveCert(context *routing.Context) error {
	Logger.Debug("%v", "SaveCert")

	certId := context.Param("cert-id")

	if len(strings.TrimSpace(certId)) == 0 || certId == "add" {
		certId = ""
	}

	cert := &Cert{}

	//Unmarshal
	if err := JsonBodyUnmarshalObject(context.RequestCtx, cert); err != nil {
		panic("参数解析出错 " + err.Error())
	}

	if len(strings.TrimSpace(cert.SerName)) == 0 {
		panic("域名不能为空")
	}

	if len(strings.TrimSpace(cert.CertKeyBlock)) == 0 {
		panic("证书key不能为空")
	}

	if len(strings.TrimSpace(cert.CertBlock)) == 0 {
		panic("证书cert不能为空")
	}

	if _, err := tls.X509KeyPair([]byte(cert.CertBlock), []byte(cert.CertKeyBlock)); err != nil {
		panic("证书cert和key不匹配，请确认证书cert和key内容")
	}

	nt := time.Now().Format("2006/1/2 15:04:05")

	if len(certId) == 0 {
		certId = uuid.Must(uuid.NewV4(), nil).String()
		cert.Id = certId
		cert.SetTime = nt
		cert.UpdateTime = ""
	} else {
		o := EtcdClient.KeyObject(CertKey(certId), reflect.TypeOf((*Cert)(nil)), ReadTimeout)

		if o == nil {
			panic("没有找到" + certId + "对应的证书数据")
		}

		oldPath := o.(*Cert)
		cert.Id = certId
		cert.SetTime = oldPath.SetTime
		cert.UpdateTime = nt
	}

	if _, err := EtcdClient.PutValue(CertKey(certId), cert, WriteTimeout); err == nil {
		Result.Success(cert, "保存证书成功").Response(context.RequestCtx)
	} else {
		Logger.Error("保存证书失败 %+v", cert)
		panic("保存证书失败")
	}

	return nil
}

func (u *certHandler) DeleteCert(context *routing.Context) error {
	Logger.Debug("%v", "DeleteCert")

	certId := context.Param("cert-id")

	if len(certId) == 0 {
		panic("没有CertId参数")
	}

	o, err := EtcdClient.KeyValue(CertKey(certId), ReadTimeout)

	if err != nil {
		panic("没有找到" + certId + "/" + "对应的证书数据")
	}

	if _, err := EtcdClient.BulkOps(func(leaseID clientv3.LeaseID) ([]clientv3.Op, error) {
		return []clientv3.Op{
			clientv3.OpDelete(CertKey(certId)),
			clientv3.OpPut(CertBakKey(certId), o, clientv3.WithLease(leaseID)),
		}, nil
	}, BakDataTTL, WriteTimeout); err != nil {
		Logger.Error("删除证书失败 %+v", err)
		panic("删除证书失败")
	}

	Result.Success("", "删除证书成功").Response(context.RequestCtx)
	return nil
}

func CertKey(certId string) string {
	return fmt.Sprintf(CERT_DATA_FORMAT, certId)
}

func CertBakKey(certId string) string {
	return fmt.Sprintf(CERT_BAK_DATA_FORMAT, certId)
}

func (u *certHandler) GetCertCount() int {
	return EtcdClient.CountWithPrefix(CERT_DATA_PREFIX, ReadTimeout)
}
