package util

import (
	"encoding/json"
	"fmt"
	. "gateway-server/model"
	. "github.com/gohutool/boot4go-etcd/client"
	util4go "github.com/gohutool/boot4go-util"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : etcd.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/23 23:38
* 修改历史 : 1. [2022/4/23 23:38] 创建文件 by LongYong
*/

const (
	DialTimeout  = 3
	ReadTimeout  = 3
	WriteTimeout = 3
	BakDataTTL   = 1800
)

const (
	DB_PREFIX = "/gateway4go/database"
)

const (
	CERT_DATA_PREFIX = DB_PREFIX + "/cert-data/"
	CERT_DATA_FORMAT = CERT_DATA_PREFIX + "%s"

	CERT_BAK_DATA_PREFIX = DB_PREFIX + "/cert-data-bak/"
	CERT_BAK_DATA_FORMAT = CERT_BAK_DATA_PREFIX + "%s"
)

const (
	GATEWAY_ACTIVE_PREFIX = DB_PREFIX + "/gateway-active/"

	GATEWAY_ACTIVE_FORMAT        = GATEWAY_ACTIVE_PREFIX + "%s/"
	GATEWAY_ACTIVE_SEARCH_FORMAT = GATEWAY_ACTIVE_PREFIX + "%s"
)

const (
	DOMAIN_DATA_PREFIX      = DB_PREFIX + "/domain-data/"
	DOMAIN_DATA_DOMAIN_PATH = DOMAIN_DATA_PREFIX + "%s/"

	DOMAIN_BAK_DATA_PREFIX      = DB_PREFIX + "/domain-data-bak/"
	DOMAIN_BAK_DATA_DOMAIN_PATH = DOMAIN_BAK_DATA_PREFIX + "%s/"
)

const (
	DOMAIN_PATH_DATA_PREFIX = DB_PREFIX + "/path-data/%s/"

	DOMAIN_PATH_DATA_WATCH_PREFIX = DB_PREFIX + "/path-data/"
	DOMAIN_BAK_PATH_DATA_PREFIX   = DB_PREFIX + "/path-data-bak/"
)

func GatewayPathKey(sername string) string {
	return fmt.Sprintf(GATEWAY_ACTIVE_FORMAT, sername)
}

func CertDatas() []*Cert {
	certs := EtcdClient.GetKeyObjectsWithPrefix(CERT_DATA_PREFIX, reflect.TypeOf((*Cert)(nil)), nil,
		0, 0, ReadTimeout)

	return util4go.CopyArray(certs, make([]*Cert, 0, len(certs)))
}

func CertBakKey(key string) string {
	return strings.Replace(key, CERT_DATA_PREFIX, CERT_BAK_DATA_PREFIX, 1)
}

func CertDataFromBak(key string) (Cert, bool) {
	cert := EtcdClient.KeyObject(CertBakKey(key), reflect.TypeOf((*Cert)(nil)), ReadTimeout, nil)

	if cert != nil {
		return *cert.(*Cert), true
	} else {
		return Cert{}, false
	}
}

func WatchCertData(listener WatchChannelEventListener) {
	EtcdClient.WatchKeyWithPrefix(CERT_DATA_PREFIX, listener, clientv3.WithPrevKV())
}

func WatchMachineData(address string, listener WatchChannelEventListener) {
	EtcdClient.WatchKey(GatewayPathKey(address), listener, clientv3.WithPrevKV())
}

func DomainDataFromBak(key string) (Domain, bool) {
	domain := EtcdClient.KeyObject(DomainBakKey(key), reflect.TypeOf((*Domain)(nil)), ReadTimeout)

	if domain != nil {
		return *domain.(*Domain), true
	} else {
		return Domain{}, false
	}
}

func DomainBakKey(key string) string {
	return strings.Replace(key, DOMAIN_DATA_PREFIX, DOMAIN_BAK_DATA_PREFIX, 1)
}

func DomainPathDataFromBak(key string) (Path, bool) {
	path := EtcdClient.KeyObject(DomainPathBakKey(key), reflect.TypeOf((*Domain)(nil)), ReadTimeout)

	if path != nil {
		return *path.(*Path), true
	} else {
		return Path{}, false
	}
}

func DomainPathBakKey(key string) string {
	return strings.Replace(key, DOMAIN_PATH_DATA_PREFIX, DOMAIN_BAK_PATH_DATA_PREFIX, 1)
}

func GetDomainData() []Domain {
	domains := EtcdClient.GetKeyObjectsWithPrefix(DOMAIN_DATA_PREFIX, reflect.TypeOf((*Domain)(nil)), nil,
		0, 0, ReadTimeout)

	if domains == nil {
		return nil
	}

	rtn := make([]Domain, 0, len(domains))

	for _, domain := range domains {
		rtn = append(rtn, *domain.(*Domain))
	}

	return rtn
}

func pathDomainPrefixKey(domainId string) string {
	return fmt.Sprintf(DOMAIN_PATH_DATA_PREFIX, domainId)
}

func GetPathData(domainId string) []Path {

	domainPaths := EtcdClient.GetKeyObjectsWithPrefix(pathDomainPrefixKey(domainId), reflect.TypeOf((*Path)(nil)), nil,
		0, 0, ReadTimeout)

	if domainPaths == nil {
		return nil
	}

	rtn := make([]Path, 0, len(domainPaths))

	for _, path := range domainPaths {
		rtn = append(rtn, *path.(*Path))
	}

	return rtn
}

func WatchDomainData(listener WatchChannelEventListener) {
	EtcdClient.WatchKeyWithPrefix(DOMAIN_DATA_PREFIX, listener, clientv3.WithPrevKV())
}

func WatchPathData(listener WatchChannelEventListener) {
	EtcdClient.WatchKeyWithPrefix(DOMAIN_PATH_DATA_WATCH_PREFIX, listener, clientv3.WithPrevKV())
}

const (
	METRICS_DATA_PREFIX         = DB_PREFIX + "/metrics-data/gateway/"
	OVERALL_METRICS_DATA_PREFIX = DB_PREFIX + "/metrics-data/overall/"
	GATEWAY_METRICS_DATA_PREFIX = METRICS_DATA_PREFIX + "%s/"

	METRICS_SUMMARY_DATA_PREFIX                   = DB_PREFIX + "/metrics-summary-data/"
	OVERALL_METRICS_SUMMARY_DATA_PREFIX           = METRICS_SUMMARY_DATA_PREFIX + "overall/"
	OVERALL_METRICS_SUMMARY_DATA_DIMENSION_PREFIX = OVERALL_METRICS_SUMMARY_DATA_PREFIX + "%s/"
	OVERALL_METRICS_SUMMARY_DATA_DIMENSION_FORMAT = OVERALL_METRICS_SUMMARY_DATA_PREFIX + "%s/%s"
	//OVERALL_METRICS_SUMMARY_DATA_PREFIX_1 = OVERALL_METRICS_SUMMARY_DATA_PREFIX + "15/"
	//OVERALL_METRICS_SUMMARY_DATA_PREFIX_15 = OVERALL_METRICS_SUMMARY_DATA_PREFIX + "15/"
	//OVERALL_METRICS_SUMMARY_DATA_PREFIX_30 = OVERALL_METRICS_SUMMARY_DATA_PREFIX + "30/"
	//OVERALL_METRICS_SUMMARY_DATA_PREFIX_60 = OVERALL_METRICS_SUMMARY_DATA_PREFIX + "60/"

	GATEWAY_METRICS_SUMMARY_PREFIX                = METRICS_SUMMARY_DATA_PREFIX + "gateway/"
	GATEWAY_METRICS_SUMMARY_DATA_PREFIX           = GATEWAY_METRICS_SUMMARY_PREFIX + "%s/"
	GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_PREFIX = GATEWAY_METRICS_SUMMARY_DATA_PREFIX + "%s/"
	GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_FORMAT = GATEWAY_METRICS_SUMMARY_DATA_PREFIX + "%s/%s"
	//GATEWAY_METRICS_SUMMARY_DATA_PREFIX_1  = GATEWAY_METRICS_SUMMARY_DATA_PREFIX + "1/"
	//GATEWAY_METRICS_SUMMARY_DATA_PREFIX_15 = GATEWAY_METRICS_SUMMARY_DATA_PREFIX + "15/"
	//GATEWAY_METRICS_SUMMARY_DATA_PREFIX_30 = GATEWAY_METRICS_SUMMARY_DATA_PREFIX + "30/"
	//GATEWAY_METRICS_SUMMARY_DATA_PREFIX_60 = GATEWAY_METRICS_SUMMARY_DATA_PREFIX + "60/"

	GATEWAY_METRICS_VISIT_CERT_DOWN_PREFIX  = GATEWAY_METRICS_DATA_PREFIX + "cert/down"
	GATEWAY_METRICS_VISIT_CERT_ERROR_PREFIX = GATEWAY_METRICS_DATA_PREFIX + "cert/fail"

	GATEWAY_METRICS_VISIT_DATA_PREFIX = GATEWAY_METRICS_DATA_PREFIX + "visit/"
	GATEWAY_METRICS_VISIT_OK_PREFIX   = GATEWAY_METRICS_VISIT_DATA_PREFIX + "ok"

	GATEWAY_METRICS_VISIT_HOST_LOSE_PREFIX       = GATEWAY_METRICS_VISIT_DATA_PREFIX + "lose/host"
	GATEWAY_METRICS_VISIT_PATH_LOSE_PREFIX       = GATEWAY_METRICS_VISIT_DATA_PREFIX + "lose/path"
	GATEWAY_METRICS_VISIT_TARGET_LOSE_PREFIX     = GATEWAY_METRICS_VISIT_DATA_PREFIX + "lose/target"
	GATEWAY_METRICS_VISIT_PROXYERROR_LOSE_PREFIX = GATEWAY_METRICS_VISIT_DATA_PREFIX + "lose/error"
	GATEWAY_METRICS_VISIT_BLACKIP_LOSE_PREFIX    = GATEWAY_METRICS_VISIT_DATA_PREFIX + "lose/blackip"

	OVERALL_METRICS_VISIT_CERT_DOWN_PREFIX  = OVERALL_METRICS_DATA_PREFIX + "cert/down"
	OVERALL_METRICS_VISIT_CERT_ERROR_PREFIX = OVERALL_METRICS_DATA_PREFIX + "cert/fail"

	OVERALL_METRICS_VISIT_DATA_PREFIX = OVERALL_METRICS_DATA_PREFIX + "visit/"
	OVERALL_METRICS_VISIT_OK_PREFIX   = OVERALL_METRICS_VISIT_DATA_PREFIX + "ok"

	OVERALL_METRICS_VISIT_HOST_LOSE_PREFIX       = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/host"
	OVERALL_METRICS_VISIT_PATH_LOSE_PREFIX       = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/path"
	OVERALL_METRICS_VISIT_TARGET_LOSE_PREFIX     = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/target"
	OVERALL_METRICS_VISIT_PROXYERROR_LOSE_PREFIX = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/error"
	OVERALL_METRICS_VISIT_BLACKIP_LOSE_PREFIX    = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/blackip"
)

func IncrVisitOk(sername string) {
	go func() {
		EtcdClient.Incr(fmt.Sprintf(GATEWAY_METRICS_VISIT_OK_PREFIX, sername), 1, 0)
		EtcdClient.Incr(OVERALL_METRICS_VISIT_OK_PREFIX, 1, 0)
	}()
}

func IncrVisitERROR(sername string) {
	go func() {
		EtcdClient.Incr(fmt.Sprintf(GATEWAY_METRICS_VISIT_PROXYERROR_LOSE_PREFIX, sername), 1, 0)
		EtcdClient.Incr(OVERALL_METRICS_VISIT_PROXYERROR_LOSE_PREFIX, 1, 0)
	}()
}

func IncrVisitBlackIP(sername string) {
	go func() {
		EtcdClient.Incr(fmt.Sprintf(GATEWAY_METRICS_VISIT_BLACKIP_LOSE_PREFIX, sername), 1, 0)
		EtcdClient.Incr(OVERALL_METRICS_VISIT_BLACKIP_LOSE_PREFIX, 1, 0)
	}()
}

func IncrVisitHostLose(sername string) {
	go func() {
		EtcdClient.Incr(fmt.Sprintf(GATEWAY_METRICS_VISIT_HOST_LOSE_PREFIX, sername), 1, 0)
		EtcdClient.Incr(OVERALL_METRICS_VISIT_HOST_LOSE_PREFIX, 1, 0)
	}()
}

func IncrVisitPathLose(sername string) {
	go func() {
		EtcdClient.Incr(fmt.Sprintf(GATEWAY_METRICS_VISIT_PATH_LOSE_PREFIX, sername), 1, 0)
		EtcdClient.Incr(OVERALL_METRICS_VISIT_PATH_LOSE_PREFIX, 1, 0)
	}()
}

func IncrVisitTargetLose(sername string) {
	go func() {
		EtcdClient.Incr(fmt.Sprintf(GATEWAY_METRICS_VISIT_TARGET_LOSE_PREFIX, sername), 1, 0)
		EtcdClient.Incr(OVERALL_METRICS_VISIT_TARGET_LOSE_PREFIX, 1, 0)
	}()
}

func IncrVisitCertDown(sername string) {
	go func() {
		EtcdClient.Incr(fmt.Sprintf(GATEWAY_METRICS_VISIT_CERT_DOWN_PREFIX, sername), 1, 0)
		EtcdClient.Incr(OVERALL_METRICS_VISIT_CERT_DOWN_PREFIX, 1, 0)
	}()
}

func IncrVisitCertError(sername string) {
	go func() {
		EtcdClient.Incr(fmt.Sprintf(GATEWAY_METRICS_VISIT_CERT_ERROR_PREFIX, sername), 1, 0)
		EtcdClient.Incr(OVERALL_METRICS_VISIT_CERT_ERROR_PREFIX, 1, 0)
	}()
}

func SaveMetricsSummary(sername string) {

	EtcdClient.LockAndDo(METRICS_SUMMARY_DATA_PREFIX, 30, func() (any, error) {

		overall := GetMetrics("")
		one := GetMetrics(sername)

		t := time.Now()
		dimension := Time2Dimension(t)

		M1, M15, M30, M60, M1D := BuildSummaryMetrics(dimension, sername, one)
		OM1, OM15, OM30, OM60, OM1D := BuildSummaryMetrics(dimension, "", overall)

		SaveSummaryMetrics(sername, *dimension, M1, M15, M30, M60, M1D)
		SaveSummaryMetrics("", *dimension, OM1, OM15, OM30, OM60, OM1D)

		return nil, nil
	})
}

func BuildSummaryMetrics(dimension *Dimension, sername string, one VisitMetrics) (m1, m2, m3, m4, m5 SummaryMetrics) {

	M1, M15, M30, M60, M1D := SummaryMetrics{}, SummaryMetrics{}, SummaryMetrics{}, SummaryMetrics{}, SummaryMetrics{}
	M1.Key = dimension.M1
	M15.Key = dimension.M15
	M30.Key = dimension.M30
	M60.Key = dimension.M60
	M1D.Key = dimension.MD

	LM1 := GetSummaryLastMetrics(sername, "1", dimension.LM1)

	if LM1 == nil {
		M1.Current = one
		M1.Increase = one
	} else {
		M1.Current = one
		M1.Increase = VisitMetrics{}

		M1.Increase.Total = one.Total - LM1.Current.Total
		M1.Increase.Ok = one.Ok - LM1.Current.Ok
		M1.Increase.Lose = one.Lose - LM1.Current.Lose
		M1.Increase.HostLose = one.HostLose - LM1.Current.HostLose
		M1.Increase.ErrorLose = one.ErrorLose - LM1.Current.ErrorLose
		M1.Increase.PathLose = one.PathLose - LM1.Current.PathLose
		M1.Increase.TargetLose = one.TargetLose - LM1.Current.TargetLose
		M1.Increase.BlackIpLose = one.BlackIpLose - LM1.Current.BlackIpLose
		M1.Increase.CertTotal = one.CertTotal - LM1.Current.CertTotal
		M1.Increase.CertDown = one.CertDown - LM1.Current.CertDown
		M1.Increase.CertFail = one.CertFail - LM1.Current.CertFail

	}

	LM15 := GetSummaryLastMetrics(sername, "15", dimension.LM15)

	if LM15 == nil {
		M15.Current = one
		M15.Increase = one
	} else {
		M15.Current = one
		M15.Increase = VisitMetrics{}

		M15.Increase.Total = one.Total - LM15.Current.Total
		M15.Increase.Ok = one.Ok - LM15.Current.Ok
		M15.Increase.Lose = one.Lose - LM15.Current.Lose
		M15.Increase.HostLose = one.HostLose - LM15.Current.HostLose
		M15.Increase.ErrorLose = one.ErrorLose - LM15.Current.ErrorLose
		M15.Increase.PathLose = one.PathLose - LM15.Current.PathLose
		M15.Increase.TargetLose = one.TargetLose - LM15.Current.TargetLose
		M15.Increase.BlackIpLose = one.BlackIpLose - LM15.Current.BlackIpLose
		M15.Increase.CertTotal = one.CertTotal - LM15.Current.CertTotal
		M15.Increase.CertDown = one.CertDown - LM15.Current.CertDown
		M15.Increase.CertFail = one.CertFail - LM15.Current.CertFail
	}

	LM30 := GetSummaryLastMetrics(sername, "30", dimension.LM30)

	if LM30 == nil {
		M30.Current = one
		M30.Increase = one
	} else {
		M30.Current = one
		M30.Increase = VisitMetrics{}

		M30.Increase.Total = one.Total - LM30.Current.Total
		M30.Increase.Ok = one.Ok - LM30.Current.Ok
		M30.Increase.Lose = one.Lose - LM30.Current.Lose
		M30.Increase.HostLose = one.HostLose - LM30.Current.HostLose
		M30.Increase.ErrorLose = one.ErrorLose - LM30.Current.ErrorLose
		M30.Increase.PathLose = one.PathLose - LM30.Current.PathLose
		M30.Increase.TargetLose = one.TargetLose - LM30.Current.TargetLose
		M30.Increase.BlackIpLose = one.BlackIpLose - LM30.Current.BlackIpLose
		M30.Increase.CertTotal = one.CertTotal - LM30.Current.CertTotal
		M30.Increase.CertDown = one.CertDown - LM30.Current.CertDown
		M30.Increase.CertFail = one.CertFail - LM30.Current.CertFail
	}

	LM60 := GetSummaryLastMetrics(sername, "60", dimension.LM60)

	if LM60 == nil {
		M60.Current = one
		M60.Increase = one
	} else {
		M60.Current = one
		M60.Increase = VisitMetrics{}

		M60.Increase.Total = one.Total - LM60.Current.Total
		M60.Increase.Ok = one.Ok - LM60.Current.Ok
		M60.Increase.Lose = one.Lose - LM60.Current.Lose
		M60.Increase.HostLose = one.HostLose - LM60.Current.HostLose
		M60.Increase.ErrorLose = one.ErrorLose - LM60.Current.ErrorLose
		M60.Increase.PathLose = one.PathLose - LM60.Current.PathLose
		M60.Increase.TargetLose = one.TargetLose - LM60.Current.TargetLose
		M60.Increase.BlackIpLose = one.BlackIpLose - LM60.Current.BlackIpLose
		M60.Increase.CertTotal = one.CertTotal - LM60.Current.CertTotal
		M60.Increase.CertDown = one.CertDown - LM60.Current.CertDown
		M60.Increase.CertFail = one.CertFail - LM60.Current.CertFail
	}

	LM1D := GetSummaryLastMetrics(sername, "1d", dimension.LMD)

	if LM1D == nil {
		M1D.Current = one
		M1D.Increase = one
	} else {
		M1D.Current = one
		M1D.Increase = VisitMetrics{}

		M1D.Increase.Total = one.Total - LM1D.Current.Total
		M1D.Increase.Ok = one.Ok - LM1D.Current.Ok
		M1D.Increase.Lose = one.Lose - LM1D.Current.Lose
		M1D.Increase.HostLose = one.HostLose - LM1D.Current.HostLose
		M1D.Increase.ErrorLose = one.ErrorLose - LM1D.Current.ErrorLose
		M1D.Increase.PathLose = one.PathLose - LM1D.Current.PathLose
		M1D.Increase.TargetLose = one.TargetLose - LM1D.Current.TargetLose
		M1D.Increase.BlackIpLose = one.BlackIpLose - LM1D.Current.BlackIpLose
		M1D.Increase.CertTotal = one.CertTotal - LM1D.Current.CertTotal
		M1D.Increase.CertDown = one.CertDown - LM1D.Current.CertDown
		M1D.Increase.CertFail = one.CertFail - LM1D.Current.CertFail
	}
	return M1, M15, M30, M60, M1D
}

func summaryMetrics2String(d SummaryMetrics) string {
	b, err := json.Marshal(d)

	if err != nil {
		panic(err)
	}

	return string(b)
}

const MANY_POINT = 120

func SaveSummaryMetrics(sername string, dimension Dimension, d1, d2, d3, d4, d5 SummaryMetrics) {
	var key1, key15, key30, key60, key1D string

	if !util4go.IsEmpty(sername) {
		key1 = fmt.Sprintf(GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, sername, "1", dimension.M1)
		key15 = fmt.Sprintf(GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, sername, "15", dimension.M15)
		key30 = fmt.Sprintf(GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, sername, "30", dimension.M30)
		key60 = fmt.Sprintf(GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, sername, "60", dimension.M60)
		key1D = fmt.Sprintf(GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, sername, "1d", dimension.MD)
	} else {
		key1 = fmt.Sprintf(OVERALL_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, "1", dimension.M1)
		key15 = fmt.Sprintf(OVERALL_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, "15", dimension.M15)
		key30 = fmt.Sprintf(OVERALL_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, "30", dimension.M30)
		key60 = fmt.Sprintf(OVERALL_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, "60", dimension.M60)
		key1D = fmt.Sprintf(OVERALL_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, "1d", dimension.MD)
	}

	_, err := EtcdClient.BulkOpsPlus(func(txn clientv3.Txn, leaseID clientv3.LeaseID) (clientv3.Txn, error) {

		var l1, l15, l30, l60, l1d clientv3.LeaseID = clientv3.NoLease, clientv3.NoLease, clientv3.NoLease, clientv3.NoLease, clientv3.NoLease

		if ol := GetSummaryMetricsLeaseByKey(key1); ol <= 0 {
			if lr, err := EtcdClient.Get().Grant(context.TODO(), MANY_POINT*60); err != nil {
				panic(fmt.Sprintf("SaveSummaryMetrics error %v", err))
			} else {
				l1 = lr.ID
				d1.LeaseID = int64(l1)
			}
		} else {
			d1.LeaseID = int64(ol)
		}

		if ol := GetSummaryMetricsLeaseByKey(key15); ol <= 0 {
			if lr, err := EtcdClient.Get().Grant(context.TODO(), MANY_POINT*15*60); err != nil {
				panic(fmt.Sprintf("SaveSummaryMetrics error %v", err))
			} else {
				l15 = lr.ID
				d2.LeaseID = int64(l15)
			}
		} else {
			d2.LeaseID = int64(ol)
		}

		if ol := GetSummaryMetricsLeaseByKey(key30); ol <= 0 {
			if lr, err := EtcdClient.Get().Grant(context.TODO(), MANY_POINT*30*60); err != nil {
				panic(fmt.Sprintf("SaveSummaryMetrics error %v", err))
			} else {
				l30 = lr.ID
				d3.LeaseID = int64(l30)
			}
		} else {
			d3.LeaseID = int64(ol)
		}

		if ol := GetSummaryMetricsLeaseByKey(key60); ol <= 0 {
			if lr, err := EtcdClient.Get().Grant(context.TODO(), MANY_POINT*60*60); err != nil {
				panic(fmt.Sprintf("SaveSummaryMetrics error %v", err))
			} else {
				l60 = lr.ID
				d4.LeaseID = int64(l60)
			}
		} else {
			d4.LeaseID = int64(ol)
		}

		if ol := GetSummaryMetricsLeaseByKey(key1D); ol <= 0 {
			if lr, err := EtcdClient.Get().Grant(context.TODO(), MANY_POINT*24*60*60); err != nil {
				panic(fmt.Sprintf("SaveSummaryMetrics error %v", err))
			} else {
				l1d = lr.ID
				d5.LeaseID = int64(l1d)
			}
		} else {
			d5.LeaseID = int64(ol)
		}

		txn.Then([]clientv3.Op{
			clientv3.OpPut(key1, summaryMetrics2String(d1), clientv3.WithLease(l1)),
			clientv3.OpPut(key15, summaryMetrics2String(d2), clientv3.WithLease(l15)),
			clientv3.OpPut(key30, summaryMetrics2String(d3), clientv3.WithLease(l30)),
			clientv3.OpPut(key60, summaryMetrics2String(d4), clientv3.WithLease(l60)),
			clientv3.OpPut(key1D, summaryMetrics2String(d5), clientv3.WithLease(l1d)),
		}...)

		return txn, nil
	}, 0, 0)

	if err != nil {
		panic(fmt.Sprintf("SaveSummaryMetrics error %v", err))
	} else {
		Logger.Debug("SaveSummaryMetrics over")
	}
}

func GetSummaryMetricsLeaseByKey(key string) int64 {
	d := EtcdClient.KeyObject(key,
		reflect.TypeOf((*SummaryMetrics)(nil)), 0)

	if d == nil {
		return 0
	} else {
		return d.(*SummaryMetrics).LeaseID
	}
}

func GetSummaryMetricsLease(sername, dimension, key string) int64 {
	m := GetSummaryMetrics(sername, dimension, key)

	if m == nil {
		return 0
	} else {
		return m.LeaseID
	}
}

func GetSummaryMetrics(sername, dimension, key string) *SummaryMetrics {
	var d any
	if !util4go.IsEmpty(sername) {
		d = EtcdClient.KeyObject(fmt.Sprintf(GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, sername, dimension, key),
			reflect.TypeOf((*SummaryMetrics)(nil)), 0)
	} else {
		d = EtcdClient.KeyObject(fmt.Sprintf(OVERALL_METRICS_SUMMARY_DATA_DIMENSION_FORMAT, dimension, key),
			reflect.TypeOf((*SummaryMetrics)(nil)), 0)
	}

	if d == nil {
		return nil
	} else {
		return d.(*SummaryMetrics)
	}
}

func GetSummaryLastMetrics(sername, dimension, key string) *SummaryMetrics {
	var d []any
	sort := SortMode.New(SortByKey, SortDescend)
	if !util4go.IsEmpty(sername) {
		d = EtcdClient.GetKeyObjectsWithPrefix(fmt.Sprintf(GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_PREFIX, sername, dimension),
			reflect.TypeOf((*SummaryMetrics)(nil)),
			&sort, 0, 2, 0)
	} else {
		d = EtcdClient.GetKeyObjectsWithPrefix(fmt.Sprintf(OVERALL_METRICS_SUMMARY_DATA_DIMENSION_PREFIX, dimension),
			reflect.TypeOf((*SummaryMetrics)(nil)),
			&sort, 0, 2, 0)
	}

	if len(d) == 0 {
		return nil
	} else {
		if d[0].(*SummaryMetrics).Key > key {
			if len(d) > 1 {
				return d[1].(*SummaryMetrics)
			} else {
				return nil
			}
		} else {
			return d[0].(*SummaryMetrics)
		}
	}
}

func GetMetrics(sername string) VisitMetrics {
	rtn := VisitMetrics{}

	var values []KeyValue

	if util4go.IsEmpty(sername) {
		values = EtcdClient.GetKeyAndValuesWithPrefix(OVERALL_METRICS_DATA_PREFIX, nil, 0, 0, 0)
	} else {
		values = EtcdClient.GetKeyAndValuesWithPrefix(fmt.Sprintf(GATEWAY_METRICS_DATA_PREFIX, sername), nil, 0, 0, 0)
	}

	var total, lose int64

	for _, kv := range values {
		dimension := getDimensionFromKey(kv.Key)
		v, _ := strconv.ParseInt(kv.Value, 10, 64)
		switch dimension {
		case "ok":
			rtn.Ok = v
			total += v
		case "host":
			lose += v
			total += v
			rtn.HostLose = v
		case "path":
			lose += v
			total += v
			rtn.PathLose = v
		case "target":
			lose += v
			total += v
			rtn.TargetLose = v
		case "error":
			lose += v
			total += v
			rtn.ErrorLose = v
		case "blackip":
			lose += v
			total += v
			rtn.BlackIpLose = v
		case "down":
			rtn.CertDown = v
		case "fail":
			rtn.CertFail = v
		}
	}

	rtn.Total = total
	rtn.Lose = lose
	rtn.CertTotal = rtn.CertFail + rtn.CertDown

	return rtn
}

const _regEx = `[\s\S]*/([\s\S]*)$`

func getDimensionFromKey(key string) string {
	values, _ := util4go.RegExpPool.FindStringSubmatch(key, _regEx)

	if len(values) <= 1 {
		return ""
	} else {
		return values[1]
	}
}
