package util

import (
	"fmt"
	. "gateway-manager/model"
	"github.com/gohutool/boot4go-etcd/client"
	. "github.com/gohutool/boot4go-util"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strconv"
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
	AUTH_DATA_PATH         = DB_PREFIX + "/auth-data"
	AUTH_INIT_DATA_PATH    = AUTH_DATA_PATH + "/init"
	ADMIN_USER_DATA_PREFIX = AUTH_DATA_PATH + "/user/"
	ADMIN_USER_DATA_PATH   = ADMIN_USER_DATA_PREFIX + "%s"
)

func CreateAdmin(username, password string) error {
	user := &AdminUser{}
	user.UserName = username

	user.UserId = MD5(username)
	user.Salt = uuid.Must(uuid.NewV4(), nil).String()
	user.Password = SaltMd5(password, user.Salt)

	_, err := client.EtcdClient.PutValue(fmt.Sprintf(ADMIN_USER_DATA_PATH, user.UserId), user, WriteTimeout)
	return err
}

func IsAdminExist() bool {
	return client.EtcdClient.CountWithPrefix(ADMIN_USER_DATA_PREFIX, ReadTimeout) > 0
}

func InitAdmin() error {
	if IsAdminExist() {
		return nil
	} else {
		Logger.Info("Init Admin user with %q/%q", "ginghan", "123456")

		return CreateAdmin("ginghan", "123456")
	}
}

const (
	METRICS_DATA_PREFIX         = DB_PREFIX + "/metrics-data/gateway/"
	OVERALL_METRICS_DATA_PREFIX = DB_PREFIX + "/metrics-data/overall/"
	GATEWAY_METRICS_DATA_PREFIX = METRICS_DATA_PREFIX + "%s/"

	GATEWAY_METRICS_VISIT_CERT_PREFIX = GATEWAY_METRICS_DATA_PREFIX + "cert/%s"
	//GATEWAY_METRICS_VISIT_CERT_DOWN_PREFIX  = GATEWAY_METRICS_DATA_PREFIX + "cert/down"
	//GATEWAY_METRICS_VISIT_CERT_ERROR_PREFIX = GATEWAY_METRICS_DATA_PREFIX + "cert/fail"

	GATEWAY_METRICS_VISIT_DATA_PREFIX = GATEWAY_METRICS_DATA_PREFIX + "visit/"
	GATEWAY_METRICS_VISIT_OK_PREFIX   = GATEWAY_METRICS_VISIT_DATA_PREFIX + "ok"
	GATEWAY_METRICS_VISIT_LOST_PREFIX = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/%s"

	//GATEWAY_METRICS_VISIT_HOST_LOSE_PREFIX       = GATEWAY_METRICS_VISIT_DATA_PREFIX + "lose/host"
	//GATEWAY_METRICS_VISIT_PATH_LOSE_PREFIX       = GATEWAY_METRICS_VISIT_DATA_PREFIX + "lose/path"
	//GATEWAY_METRICS_VISIT_TARGET_LOSE_PREFIX     = GATEWAY_METRICS_VISIT_DATA_PREFIX + "lose/target"
	//GATEWAY_METRICS_VISIT_PROXYERROR_LOSE_PREFIX = GATEWAY_METRICS_VISIT_DATA_PREFIX + "lose/error"
	//GATEWAY_METRICS_VISIT_BLACKIP_LOSE_PREFIX    = GATEWAY_METRICS_VISIT_DATA_PREFIX + "lose/blackip"

	OVERALL_METRICS_VISIT_CERT_PREFIX = GATEWAY_METRICS_DATA_PREFIX + "cert/%s"
	//OVERALL_METRICS_VISIT_CERT_DOWN_PREFIX  = OVERALL_METRICS_DATA_PREFIX + "cert/down"
	//OVERALL_METRICS_VISIT_CERT_ERROR_PREFIX = OVERALL_METRICS_DATA_PREFIX + "cert/fail"

	OVERALL_METRICS_VISIT_DATA_PREFIX = OVERALL_METRICS_DATA_PREFIX + "visit/"
	OVERALL_METRICS_VISIT_OK_PREFIX   = OVERALL_METRICS_VISIT_DATA_PREFIX + "ok"
	OVERALL_METRICS_VISIT_LOST_PREFIX = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/%s"

	//OVERALL_METRICS_VISIT_HOST_LOSE_PREFIX       = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/host"
	//OVERALL_METRICS_VISIT_PATH_LOSE_PREFIX       = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/path"
	//OVERALL_METRICS_VISIT_TARGET_LOSE_PREFIX     = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/target"
	//OVERALL_METRICS_VISIT_PROXYERROR_LOSE_PREFIX = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/error"
	//OVERALL_METRICS_VISIT_BLACKIP_LOSE_PREFIX    = OVERALL_METRICS_VISIT_DATA_PREFIX + "lose/blackip"
)

func GetMetrics(sername string) VisitMetrics {
	rtn := VisitMetrics{}

	var values []client.KeyValue

	if IsEmpty(sername) {
		rtn.SerName = ""
		values = client.EtcdClient.GetKeyAndValuesWithPrefix(OVERALL_METRICS_DATA_PREFIX, nil, 0, 0, 0)
	} else {
		rtn.SerName = sername
		values = client.EtcdClient.GetKeyAndValuesWithPrefix(fmt.Sprintf(GATEWAY_METRICS_DATA_PREFIX, sername), nil, 0, 0, 0)
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
	values, _ := RegExpPool.FindStringSubmatch(key, _regEx)

	if len(values) <= 1 {
		return ""
	} else {
		return values[1]
	}
}

const (
	DOMAIN_DATA_PREFIX      = DB_PREFIX + "/domain-data/"
	DOMAIN_DATA_DOMAIN_PATH = DOMAIN_DATA_PREFIX + "%s/"

	DOMAIN_BAK_DATA_PREFIX      = DB_PREFIX + "/domain-data-bak/"
	DOMAIN_BAK_DATA_DOMAIN_PATH = DOMAIN_BAK_DATA_PREFIX + "%s/"
)

const (
	DOMAIN_PATH_PREFIX      = DB_PREFIX + "/path-data/"
	DOMAIN_PATH_DATA_PREFIX = DOMAIN_PATH_PREFIX + "%s/"
	DOMAIN_PATH_DATA_PATH   = DOMAIN_PATH_DATA_PREFIX + "%s"

	DOMAIN_BAK_PATH_DATA_PREFIX = DB_PREFIX + "/path-data-bak/%s/"
	DOMAIN_BAK_PATH_DATA_PATH   = DOMAIN_BAK_PATH_DATA_PREFIX + "%s"
)

func GetDomains() []any {
	domains := client.EtcdClient.GetKeyObjectsWithPrefix(DOMAIN_DATA_PREFIX, reflect.TypeOf((*Domain)(nil)), nil,
		0, 0, ReadTimeout)

	return domains
}

func GetDomainPaths() []any {
	domainPaths := client.EtcdClient.GetKeyObjectsWithPrefix(DOMAIN_PATH_PREFIX, reflect.TypeOf((*Path)(nil)), nil,
		0, 0, ReadTimeout)

	return domainPaths
}

func GetUsers() []any {
	users := client.EtcdClient.GetKeyObjectsWithPrefix(ADMIN_USER_DATA_PREFIX, reflect.TypeOf((*AdminUser)(nil)), nil,
		0, 0, ReadTimeout)
	return users
}

const (
	METRICS_SUMMARY_DATA_PREFIX = DB_PREFIX + "/metrics-summary-data/"

	OVERALL_METRICS_SUMMARY_DATA_PREFIX           = METRICS_SUMMARY_DATA_PREFIX + "overall/"
	OVERALL_METRICS_SUMMARY_DATA_DIMENSION_PREFIX = OVERALL_METRICS_SUMMARY_DATA_PREFIX + "%s/"
	OVERALL_METRICS_SUMMARY_DATA_DIMENSION_FORMAT = OVERALL_METRICS_SUMMARY_DATA_PREFIX + "%s/%s"

	GATEWAY_METRICS_SUMMARY_PREFIX                = METRICS_SUMMARY_DATA_PREFIX + "gateway/"
	GATEWAY_METRICS_SUMMARY_DATA_PREFIX           = GATEWAY_METRICS_SUMMARY_PREFIX + "%s/"
	GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_PREFIX = GATEWAY_METRICS_SUMMARY_DATA_PREFIX + "%s/"
	GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_FORMAT = GATEWAY_METRICS_SUMMARY_DATA_PREFIX + "%s/%s"
)

func GetVisitSummaryMetrics(sername, dimension string, pointCount int) []SummaryMetrics {
	var d []any
	sort := client.SortMode.New(client.SortByKey, client.SortDescend)
	if !IsEmpty(sername) {
		d = client.EtcdClient.GetKeyObjectsWithPrefix(fmt.Sprintf(GATEWAY_METRICS_SUMMARY_DATA_DIMENSION_PREFIX, sername, dimension),
			reflect.TypeOf((*SummaryMetrics)(nil)),
			&sort, 0, pointCount, 0)
	} else {
		d = client.EtcdClient.GetKeyObjectsWithPrefix(fmt.Sprintf(OVERALL_METRICS_SUMMARY_DATA_DIMENSION_PREFIX, dimension),
			reflect.TypeOf((*SummaryMetrics)(nil)),
			&sort, 0, pointCount, 0)
	}

	var rtn []SummaryMetrics

	for _, v := range d {
		rtn = append(rtn, *v.(*SummaryMetrics))
	}

	return rtn
}

func GetSummaryMetrics(sername string, dimension, count int) []MetricsPoint {

	var rtn []MetricsPoint

	axis := GenerateXAxis(time.Now(), dimension, count)

	var d string
	if dimension == 100 {
		d = "1d"
	} else {
		d = strconv.Itoa(dimension)
	}

	datas := GetVisitSummaryMetrics(sername, d, count)
	m := Map(datas, func(t SummaryMetrics) string {
		return t.Key
	})

	Stream(axis, func(one TimePoint) {
		point := MetricsPoint{Point: one}

		if v, ok := m[one.Value]; ok {
			point.Value = &v
		}

		rtn = append(rtn, point)
	})

	l := len(rtn)
	for idx := 0; idx < l; idx++ {
		o := rtn[idx]
		if o.Value == nil {
			if idx == 0 {
				if len(datas) > 0 {
					o.Value = &datas[len(datas)-1]
				} else {
					o.Value = &SummaryMetrics{
						Key:      o.Point.Value,
						Increase: VisitMetrics{},
						Current:  VisitMetrics{},
					}
				}
			} else {
				last := getLastValidPoint(rtn, idx)

				if last != nil {
					temp := *last
					o.Value = &SummaryMetrics{
						Key:      o.Point.Value,
						Current:  temp.Current,
						Increase: VisitMetrics{},
					}
				} else {
					next := rtn[idx-1].Value
					o.Value = &SummaryMetrics{
						Key: o.Point.Value,
						Current: VisitMetrics{
							Total:       next.Current.Total - next.Increase.Total,
							Ok:          next.Current.Ok - next.Increase.Ok,
							Lose:        next.Current.Lose - next.Increase.Lose,
							HostLose:    next.Current.HostLose - next.Increase.HostLose,
							ErrorLose:   next.Current.ErrorLose - next.Increase.ErrorLose,
							PathLose:    next.Current.PathLose - next.Increase.PathLose,
							TargetLose:  next.Current.TargetLose - next.Increase.TargetLose,
							BlackIpLose: next.Current.BlackIpLose - next.Increase.BlackIpLose,
							CertTotal:   next.Current.CertTotal - next.Increase.CertTotal,
							CertDown:    next.Current.CertDown - next.Increase.CertDown,
							CertFail:    next.Current.CertFail - next.Increase.CertFail,
						},
						Increase: VisitMetrics{},
					}
				}
			}

			rtn[idx] = o
		} else {
			fmt.Sprintf("******** %v\n", *o.Value)
		}

	}

	return rtn

}

func getLastValidPoint(datas []MetricsPoint, current int) *SummaryMetrics {

	for idx := current + 1; idx < len(datas); idx++ {
		if datas[idx].Value != nil {
			return datas[idx].Value
		}
	}

	return nil
}
