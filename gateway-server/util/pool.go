package util

import (
	"crypto/tls"
	. "gateway-server/model"
	util4go "github.com/gohutool/boot4go-util"
	"strconv"
	"sync"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : pool.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/1 19:48
* 修改历史 : 1. [2022/5/1 19:48] 创建文件 by LongYong
*/

var SERVER_NAME string

var UP = time.Now()

//var certMap map[string]*tls.Certificate
//var myMetrics metrics.Registry

var CertPool = certPool{}

type certPool struct {
	lock    sync.RWMutex
	certMap map[string]*tls.Certificate
}

func (c *certPool) InitCertPool() {

	c.lock.Lock()
	defer c.lock.Unlock()

	c.certMap = make(map[string]*tls.Certificate)
	certs := CertDatas()

	for _, cert := range certs {
		certificate, err := tls.X509KeyPair([]byte(cert.CertBlock), []byte(cert.CertKeyBlock))
		if err != nil {
			Logger.Warning("证书生成失败 %s", string(cert.SerName))
		}

		c.certMap[cert.SerName] = &certificate
		Logger.Warning("%v 证书生成", string(cert.SerName))
	}

	Logger.Info("所有域名证书设置完成")
}

func (c *certPool) GetCert(host string) *tls.Certificate {

	c.lock.RLock()
	defer c.lock.RUnlock()

	if v, ok := c.certMap[host]; ok {
		return v
	} else {
		return nil
	}
}

func (c *certPool) PutCert(host string, cert *tls.Certificate) {

	c.lock.Lock()
	defer c.lock.Unlock()

	c.certMap[host] = cert
}

func (c *certPool) DelCert(host string) {
	delete(c.certMap, host)
}

type domainPool struct {
	lock      sync.RWMutex
	domainMap map[string]Domain
}

var DomainPool = domainPool{}

func (d *domainPool) GetDomain(host string) (Domain, bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	// return d.domainMap[host]
	obj, ok := d.domainMap[host]
	return obj, ok
}

func (d *domainPool) InitDomainPool() {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.domainMap = make(map[string]Domain)

	domains := GetDomainData()

	for _, domain := range domains {
		d.RegistryNewDomain(domain, true)
	}
}

func (d *domainPool) RegistryNewDomain(domain Domain, loadSSL bool) {
	domain.Path = GetPathData(domain.Id)
	d.domainMap[domain.DomainUrl] = domain

	if loadSSL && domain.SslOn {
		d.EnableSSL(domain)
	}

	Logger.Info("Domain(%v) with %v is load finish %+v\n", domain.DomainUrl, domain.Id, d.domainMap[domain.DomainUrl])
}

func (d *domainPool) EnableSSL(domain Domain) {
	var nd *Domain
	nd = util4go.Copy(&domain)
	nd.DomainUrl = util4go.ReplacePort(nd.DomainUrl, strconv.Itoa(nd.SslPort))

	d.domainMap[nd.DomainUrl] = *nd

	Logger.Info("SSL Domain(%v) with %v is load finish %+v\n", nd.DomainUrl, domain.Id, d.domainMap[nd.DomainUrl])
}

func (d domainPool) ReLoadDomain(domain, preDomain *Domain) {

	defer Logger.Info("域名更新事件 Domain[" + domain.DomainUrl + "] 结束")

	d.lock.Lock()
	d.lock.Unlock()

	if preDomain == nil {
		d.RegistryNewDomain(*domain, true)
		Logger.Info("域名更新事件 Domain[" + domain.DomainUrl + "] 新增域名完成")
		return
	}

	needChangeKey := domain.DomainUrl != preDomain.DomainUrl

	var poolDomain Domain
	var ok bool

	if poolDomain, ok = d.domainMap[preDomain.DomainUrl]; ok {
		poolDomain.DomainUrl = domain.DomainUrl
		poolDomain.DomainName = domain.DomainName
		poolDomain.SslOn = domain.SslOn
		poolDomain.SslPort = domain.SslPort
		poolDomain.BlackIps = domain.BlackIps
		poolDomain.Demo = domain.Demo
	} else {
		Logger.Error("Old Domain[" + preDomain.DomainUrl + "] in Pool is lose, will load refresh ")
		d.RegistryNewDomain(*domain, false)
		poolDomain, _ = d.domainMap[domain.DomainUrl]
	}

	if needChangeKey {
		delete(d.domainMap, preDomain.DomainUrl)
		d.domainMap[poolDomain.DomainUrl] = poolDomain
		Logger.Info("域名更新事件 Old Domain[" + preDomain.DomainUrl + "] is changed to " + domain.DomainUrl)
	} else {
		d.domainMap[poolDomain.DomainUrl] = poolDomain
	}

	if !domain.SslOn && !preDomain.SslOn {
		Logger.Info("域名更新事件 Domain[" + preDomain.DomainUrl + "] SSL no changed")
		return
	}

	if domain.SslOn && !preDomain.SslOn { // Add SSL
		d.EnableSSL(poolDomain)
		Logger.Info("域名更新事件 Domain[" + preDomain.DomainUrl + "] SSL is add")

		return
	}

	if !domain.SslOn && preDomain.SslOn { // Remove SSL
		delete(d.domainMap, util4go.ReplacePort(preDomain.DomainUrl, GetSSLPort(preDomain.SslPort)))
		Logger.Info("域名更新事件 Domain[" + preDomain.DomainUrl + "] SSL is removed")
		return
	}

	// domain.SslOn && preDomain.SslOn
	// SSL have registry before

	var preObjSSL Domain
	var okSSL bool

	oldSSLHost := util4go.ReplacePort(preDomain.DomainUrl, GetSSLPort(preDomain.SslPort))

	if preObjSSL, okSSL = d.domainMap[oldSSLHost]; okSSL {
		preObjSSL.DomainUrl = util4go.ReplacePort(domain.DomainUrl, GetSSLPort(domain.SslPort))
		preObjSSL.DomainName = domain.DomainName
		preObjSSL.SslOn = domain.SslOn
		preObjSSL.SslPort = domain.SslPort
		preObjSSL.BlackIps = domain.BlackIps
		preObjSSL.Demo = domain.Demo

		isSslChange := preDomain.SslPort != domain.SslPort

		if isSslChange {
			delete(d.domainMap, oldSSLHost)
			d.EnableSSL(poolDomain)
			Logger.Info("域名更新事件 Old SSL Domain[" + oldSSLHost + "] is changed to " +
				util4go.ReplacePort(domain.DomainUrl, GetSSLPort(domain.SslPort)))
			return
		} else {
			d.domainMap[oldSSLHost] = preObjSSL
		}
	} else {
		d.EnableSSL(poolDomain)
		Logger.Error("域名更新事件 Old SSL Domain[" + oldSSLHost + "] in Pool is lose, will load refresh ")
	}
}

func GetSSLPort(port int) string {
	if port <= 0 {
		return "9443"
	}

	return strconv.Itoa(port)
}

func (d domainPool) RemoveDomain(domainUrl string, sslPort string) {
	d.lock.Lock()
	d.lock.Unlock()

	ip, _, err := util4go.SplitHostPort(domainUrl)

	if err != nil {
		Logger.Warning(domainUrl + " parse error " + err.Error())
		return
	}

	delete(d.domainMap, domainUrl)
	delete(d.domainMap, ip+sslPort)

	Logger.Info("域名删除事件[" + domainUrl + "]删除完成")
}

func (d domainPool) ReLoadDomainPath(path, prePath *Path) {

	defer Logger.Info("域名路径映射更新事件 Domain[%v %v] 结束", path.DomainUrl, path.ReqPath)

	d.lock.Lock()
	d.lock.Unlock()

	poolDomain, ok := d.domainMap[path.DomainUrl]

	if !ok {
		Logger.Error("域名路径映射更新事件 Domain[%v] is lose in pool", path.DomainUrl)
		return
	}

	var sslHost = ""
	var poolSSL *Domain

	if poolDomain.SslOn {
		sslHost = util4go.ReplacePort(poolDomain.DomainUrl, strconv.Itoa(poolDomain.SslPort))
		poolSSLDomain, okSSL := d.domainMap[sslHost]

		if okSSL {
			poolSSL = &poolSSLDomain
		}
	}

	paths := poolDomain.Path

	if prePath == nil { // add
		paths = append(paths, *path)
		poolDomain.Path = paths
		d.domainMap[path.DomainUrl] = poolDomain

		Logger.Warning("域名路径映射更新事件 Domain[%v][%v] is add in pool and append new", path.DomainUrl, path.ReqPath)

		if poolSSL != nil {
			poolSSL.Path = paths
			d.domainMap[sslHost] = *poolSSL
			Logger.Warning("域名路径映射更新事件 SSL Domain[%v][%v] is add in pool and append new", poolSSL.DomainUrl, path.ReqPath)
		}

		return
	}

	locIdx := -1
	for idx, p := range paths {
		if p.Id == path.Id {
			locIdx = idx
			break
		}
	}

	if locIdx == -1 {
		// if the pool path is losed, append the new to paths
		paths = append(paths, *path)
		poolDomain.Path = paths
		d.domainMap[path.DomainUrl] = poolDomain
		Logger.Warning("域名路径映射更新事件 Domain[%v][%v] is lose in pool and append new", path.DomainUrl, path.ReqPath)

		if poolSSL != nil {
			poolSSL.Path = paths
			d.domainMap[sslHost] = *poolSSL
			Logger.Warning("域名路径映射更新事件 SSL Domain[%v][%v] is lose in pool and append new", poolSSL.DomainUrl, path.ReqPath)
		}

	} else {
		paths = util4go.ReplaceAt(paths, locIdx, *path)
		poolDomain.Path = paths
		d.domainMap[path.DomainUrl] = poolDomain
		Logger.Warning("域名路径映射更新事件 Domain[%v][%v] is reload into pool", path.DomainUrl, path.ReqPath)

		if poolSSL != nil {
			poolSSL.Path = paths
			d.domainMap[sslHost] = *poolSSL
			Logger.Warning("域名路径映射更新事件 SSL Domain[%v][%v] is reload into pool", poolSSL.DomainUrl, path.ReqPath)
		}
	}
}

func (d domainPool) RemovePath(path Path) {
	d.lock.Lock()
	d.lock.Unlock()

	poolDomain, ok := d.domainMap[path.DomainUrl]

	if !ok {
		Logger.Error("Domain[%v] is lose in pool", path.DomainUrl)
		return
	}

	paths := poolDomain.Path

	locIdx := -1

	for idx, p := range paths {
		if p.Id == path.Id {
			locIdx = idx
			break
		}
	}

	if locIdx == -1 {
		Logger.Warning("Domain[%v][%v] is removed from pool and do nothing", path.DomainUrl, path.ReqPath)
	} else {
		paths = util4go.RemoveAt(paths, locIdx)
		poolDomain.Path = paths
		d.domainMap[path.DomainUrl] = poolDomain
		Logger.Warning("Domain[%v][%v] is removed from pool", path.DomainUrl, path.ReqPath)
	}

	Logger.Info("域名路径映射删除事件[%v][%v]删除完成", path.DomainUrl, path.ReqPath)
}

var Metrics_Enable bool

func MetricsCertDown() {
	if Metrics_Enable {
		IncrVisitCertDown(SERVER_NAME)
	}
}

func MetricsCertFail() {
	if Metrics_Enable {
		IncrVisitCertError(SERVER_NAME)
	}
}

func MetricsVisitOK() {
	if Metrics_Enable {
		IncrVisitOk(SERVER_NAME)
	}
}

func MetricsVisitHostLose() {
	if Metrics_Enable {
		IncrVisitHostLose(SERVER_NAME)
	}
}

func MetricsVisitTargetLose() {
	if Metrics_Enable {
		IncrVisitTargetLose(SERVER_NAME)
	}
}

func MetricsVisitPathLose() {
	if Metrics_Enable {
		IncrVisitPathLose(SERVER_NAME)
	}
}

func MetricsVisitError() {
	if Metrics_Enable {
		IncrVisitERROR(SERVER_NAME)
	}
}

func MetricsVisitBlackIP() {
	if Metrics_Enable {
		IncrVisitBlackIP(SERVER_NAME)
	}
}
