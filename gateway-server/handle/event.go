package handle

import (
	"crypto/tls"
	. "gateway-server/model"
	. "gateway-server/util"
	"github.com/gohutool/boot4go-etcd/client"
	fastjson "github.com/gohutool/boot4go-fastjson"
	util4go "github.com/gohutool/boot4go-util"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : event.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/2 10:33
* 修改历史 : 1. [2022/5/2 10:33] 创建文件 by LongYong
*/

var WatchEventListener = watchEventListener{}

type watchEventListener struct {
}

func (wl *watchEventListener) InitMachineDataNotify(address string) {

	serName := util4go.GuessIP(client.EtcdClient.Get().Endpoints()[0])
	if serName == nil {
		*serName = "192.168.0.1"
	}

	address = util4go.ReplaceIP(address, *serName)

	data := GetMachineData()
	data["start_time"] = UP.Format("2006-01-02 15:04:05")
	data["run_time"] = time.Now().Unix() - UP.Unix()
	data["id"] = address
	data["server_name"] = address
	SERVER_NAME = address

	var counter *int = new(int)
	*counter = 0

	registryMachine(data, counter)

	initGatewayChangeListener(address, data, counter)
}

func initGatewayChangeListener(address string, data map[string]interface{}, counter *int) {
	go WatchMachineData(address, func(event *clientv3.Event) {
		switch event.Type {

		case clientv3.EventTypePut:
			{
				Logger.Warning("网关信息变更 %s", string(event.Kv.Value))
			}
		case clientv3.EventTypeDelete:
			{
				*counter = *counter + 1
				data["register_times"] = *counter
				data["last_register_time"] = UP.Format("2006-01-02 15:04:05")
				key := GatewayPathKey(address)
				Logger.Info("=== 网关信息不被删除 %q to path %q with %+v [%v]", address, key, data, *counter)
				client.EtcdClient.PutValue(key, data, 0)

			}
		}
	})
}

func registryMachine(data map[string]interface{}, counter *int) (string, clientv3.LeaseID, error) {
	*counter = *counter + 1
	id, _ := data["id"].(string)
	data["register_times"] = *counter
	data["last_register_time"] = UP.Format("2006-01-02 15:04:05")

	Logger.Info("*** Registry Gateway %q to path %q with %+v [%v]", id, GatewayPathKey(id), data, *counter)

	return client.EtcdClient.PutKeepAliveValue(GatewayPathKey(id), data, 30, 0, func(leaseID clientv3.LeaseID) error {
		registryMachine(data, counter)
		return nil
	})
}

func (wl *watchEventListener) InitCertChangeListener() {
	go WatchCertData(func(event *clientv3.Event) {
		switch event.Type {
		case clientv3.EventTypePut:
			cert := new(Cert)
			if err := fastjson.UnmarshalObject(string(event.Kv.Value), cert); err != nil {
				Logger.Warning("证书数据解析失败 %s", string(event.Kv.Value))
			}
			certificate, err := tls.X509KeyPair([]byte(cert.CertBlock), []byte(cert.CertKeyBlock))
			if err != nil {
				Logger.Warning("证书生成失败 %s", string(event.Kv.Value))
			}
			CertPool.PutCert(cert.SerName, &certificate)
			Logger.Info("域名%s证书更新完成", cert.SerName)
		case clientv3.EventTypeDelete:
			certBak, ok := CertDataFromBak(string(event.Kv.Key))
			if !ok {
				Logger.Info("【域名证书路径%s】删除-获取备份数据失败", string(event.Kv.Key))
			} else {
				CertPool.DelCert(certBak.SerName)
			}
		default:
			Logger.Info("【域名证书路径%s】事件 %v", string(event.Kv.Key), event.Type)
		}
	})
}

func (wl *watchEventListener) InitDomainChangeListener() {
	go WatchDomainData(func(event *clientv3.Event) {
		switch event.Type {
		case clientv3.EventTypePut:
			domain := new(Domain)
			preDomain := new(Domain)
			if err := fastjson.UnmarshalObject(string(event.Kv.Value), domain); err != nil {
				Logger.Warning("域名信息解析失败 %v", string(event.Kv.Value))
			}
			if event.PrevKv == nil {
				preDomain = nil
			} else {
				if err := fastjson.UnmarshalObject(string(event.PrevKv.Value), preDomain); err != nil {
					Logger.Warning("旧域名信息解析失败 %v", string(event.PrevKv.Value))
				}
			}

			DomainPool.ReLoadDomain(domain, preDomain)
		case clientv3.EventTypeDelete:
			domainBak, ok := DomainDataFromBak(string(event.Kv.Key))
			if !ok {
				Logger.Info("【域名%s】获取备份数据失败 %v ", string(event.Kv.Key))
			} else {
				DomainPool.RemoveDomain(domainBak.DomainUrl, strconv.Itoa(domainBak.SslPort))
			}
		default:
			Logger.Info("【域名%s】事件 %v", string(event.Kv.Key), event.Type)
		}
	})
}

func (wl *watchEventListener) InitDomainPathChangeListener() {
	go WatchPathData(func(event *clientv3.Event) {
		switch event.Type {
		case clientv3.EventTypePut:
			path := new(Path)
			prePath := new(Path)
			if err := fastjson.UnmarshalObject(string(event.Kv.Value), path); err != nil {
				Logger.Warning("域名路径映射信息解析失败 %v", string(event.Kv.Value))
			}

			if event.PrevKv == nil {
				prePath = nil
			} else {
				if err := fastjson.UnmarshalObject(string(event.PrevKv.Value), prePath); err != nil {
					Logger.Warning("旧域名路径映射信息解析失败 %v", string(event.PrevKv.Value))
				}
			}

			DomainPool.ReLoadDomainPath(path, prePath)
		case clientv3.EventTypeDelete:
			pathBak, ok := DomainPathDataFromBak(string(event.Kv.Key))
			if !ok {
				Logger.Info("【域名路径映射%s】获取备份数据失败 %v ", string(event.Kv.Key))
			} else {
				DomainPool.RemovePath(pathBak)
			}
		default:
			Logger.Info("【域名路径映射%s】事件 %v", string(event.Kv.Key), event.Type)
		}
	})
}

func (wl *watchEventListener) InitMetricsScheduleJob() {
	go func() {
		for {
			func() {
				defer func() {
					if err := recover(); err != nil {
						Logger.Warning("MetricsScheduleJob error %v", err)
					}
				}()

				time.Sleep(10 * time.Second)

				SaveMetricsSummary(SERVER_NAME)

			}()
		}
	}()
}
