package model

import (
	. "github.com/gohutool/boot4go-fastjson"
	util4go "github.com/gohutool/boot4go-util"
	"strconv"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : module.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/25 10:00
* 修改历史 : 1. [2022/4/25 10:00] 创建文件 by LongYong
*/

type Domain struct {
	Id                 string    `json:"id"`
	DomainName         string    `json:"domain_name"`
	DomainUrl          string    `json:"domain_url"`
	LbType             string    `json:"lb_type"`
	Targets            []*Target `json:"targets"`
	BlackIps           []string  `json:"black_ips"`
	RateLimiterNum     float64   `json:"rate_limiter_num"`
	RateLimiterMsg     string    `json:"rate_limiter_msg"`
	RateLimiterEnabled bool      `json:"rate_limiter_enabled"`
	SetTime            string    `json:"set_time"`
	UpdateTime         string    `json:"update_time"`
	Demo               string    `json:"demo"`
	SslOn              bool      `json:"ssl_on"`
	SslPort            int       `json:"ssl_port"`
}

func (d *Domain) Unmarshal(value *Value) error {
	d.Id = value.GetString("id")
	d.DomainName = value.GetString("domain_name")
	d.DomainUrl = value.GetString("domain_url")
	d.LbType = value.GetString("lb_type")
	d.RateLimiterNum = value.GetFloat64("rate_limiter_num")
	d.RateLimiterMsg = value.GetString("rate_limiter_msg")
	d.RateLimiterEnabled = value.GetBool("rate_limiter_enabled")
	d.SetTime = value.GetString("set_time")
	d.UpdateTime = value.GetString("update_time")
	d.Demo = value.GetString("demo")

	d.SslOn = value.GetBool("ssl_on")
	d.SslPort = value.GetInt("ssl_port")
	if d.SslPort <= 0 {
		d.SslPort = 9443
	}

	if list, err := UnmarshalObjectList(value.Get("targets"), &Target{}); err == nil {
		d.Targets = list
	} else {
		return err
	}

	if list, err := UnmarshalStringList(value.Get("black_ips")); err == nil {
		d.BlackIps = list
	} else {
		return err
	}

	return nil
}

type Target struct {
	Pointer       string `json:"pointer"`
	PointerType   string `json:"pointer_type"`
	Weight        int8   `json:"weight"`
	CurrentWeight int8   `json:"current_weight"`
	Schema        string `json:"schema"`
	Host          string `json:"host"`
	Query         string `json:"query"`
}

func (t *Target) Unmarshal(value *Value) error {
	t.Pointer = value.GetString("pointer")
	t.PointerType = value.GetString("pointer_type")
	t.Weight = int8(value.GetInt("weight"))
	t.CurrentWeight = int8(value.GetInt("current_weight"))

	if s, h, q, err := util4go.ParseURL(t.Pointer); err == nil {
		t.Schema = s
		t.Host = h
		t.Query = q
	} else {
		t.Schema = value.GetString("schema")
		t.Host = value.GetString("host")
		t.Query = value.GetString("query")
	}

	return nil
}

type Path struct {
	Id                    string `json:"id"`
	DomainId              string `json:"domain_id"`
	DomainUrl             string `json:"domain_url"`
	ReqMethod             string `json:"req_method"`
	ReqPath               string `json:"req_path"`
	ReqName               string `json:"req_name"`
	SearchPath            string `json:"search_path"`
	ReplacePath           string `json:"replace_path"`
	CircuitBreakerRequest int    `json:"circuit_breaker_request"`
	CircuitBreakerPercent int    `json:"circuit_breaker_percent"`
	CircuitBreakerTimeout int    `json:"circuit_breaker_timeout"`
	CircuitBreakerMsg     string `json:"circuit_breaker_msg"`
	CircuitBreakerEnabled bool   `json:"circuit_breaker_enabled"`
	CircuitBreakerForce   bool   `json:"circuit_breaker_force"`
	PrivateProxyEnabled   bool   `json:"private_proxy_enabled"`
	LbType                string `json:"lb_type"`

	BlackIps           map[string]bool `json:"black_ips"`
	RateLimiterNum     float64         `json:"rate_limiter_num"`
	RateLimiterMsg     string          `json:"rate_limiter_msg"`
	RateLimiterEnabled bool            `json:"rate_limiter_enabled"`

	Targets    []*Target `json:"targets"`
	SetTime    string    `json:"set_time"`
	UpdateTime string    `json:"update_time"`
}

func (p *Path) Unmarshal(value *Value) error {
	p.Id = value.GetString("id")
	p.DomainId = value.GetString("domain_id")
	p.ReqMethod = value.GetString("req_method")
	p.ReqPath = value.GetString("req_path")
	p.ReqName = value.GetString("req_name")
	p.SearchPath = value.GetString("search_path")
	p.ReplacePath = value.GetString("replace_path")
	p.CircuitBreakerRequest = value.GetInt("circuit_breaker_request")
	p.CircuitBreakerPercent = value.GetInt("circuit_breaker_percent")
	p.CircuitBreakerTimeout = value.GetInt("circuit_breaker_timeout")
	p.CircuitBreakerMsg = value.GetString("circuit_breaker_msg")
	p.CircuitBreakerEnabled = value.GetBool("circuit_breaker_enabled")
	p.CircuitBreakerForce = value.GetBool("circuit_breaker_force")
	p.PrivateProxyEnabled = value.GetBool("private_proxy_enabled")
	p.LbType = value.GetString("lb_type")

	p.LbType = value.GetString("lb_type")
	p.RateLimiterNum = value.GetFloat64("rate_limiter_num")
	p.RateLimiterMsg = value.GetString("rate_limiter_msg")
	p.RateLimiterEnabled = value.GetBool("rate_limiter_enabled")

	p.DomainUrl = value.GetString("domain_url")

	p.SetTime = value.GetString("set_time")
	p.SetTime = value.GetString("update_time")

	if list, err := UnmarshalObjectList(value.Get("targets"), &Target{}); err == nil {
		p.Targets = list
	} else {
		return err
	}

	//
	//if list, err := UnmarshalObjectList(value.Get("targets"), &Target{}); err == nil {
	//	d.Targets = list
	//} else {
	//	return err
	//}

	if list, err := UnmarshalBoolMap(value.Get("black_ips")); err == nil {
		p.BlackIps = list
	} else {
		return err
	}

	return nil
}

type Cert struct {
	Id           string `json:"id"`
	SerName      string `json:"ser_name"`
	CertBlock    string `json:"cert_block"`
	CertKeyBlock string `json:"cert_key_block"`
	SetTime      string `json:"set_time"`
	UpdateTime   string `json:"update_time"`
}

func (p *Cert) Unmarshal(value *Value) error {
	p.Id = value.GetString("id")
	p.SerName = value.GetString("ser_name")
	p.CertBlock = value.GetString("cert_block")
	p.CertKeyBlock = value.GetString("cert_key_block")
	p.SetTime = value.GetString("set_time")
	p.UpdateTime = value.GetString("update_time")

	return nil
}

type Gateway struct {
	ServerName string `json:"server_name"`
	UpTime     string `json:"start_time"`
	ID         string `json:"id"`
	RunTime    int    `json:"run_time"`
	RegTime    string `json:"last_register_time"`
	RegTimes   int    `json:"register_times"`
	Platform   string `json:"platform"`
}

func (g *Gateway) Unmarshal(value *Value) error {
	g.ServerName = value.GetString("server_name")
	g.UpTime = value.GetString("start_time")
	g.ID = value.GetString("id")
	g.Platform = value.GetString("platform")
	g.RunTime = value.GetInt("run_time")
	g.RegTime = value.GetString("last_register_time")
	g.RegTimes = value.GetInt("register_times")

	return nil
}

type ClusterType string

const (
	NODE    ClusterType = "Node"
	CLUSTER ClusterType = "Cluster"
	ETCD    ClusterType = "Etcd"
	NACOS   ClusterType = "Nacos"
	ZK      ClusterType = "ZK"
)

type Cluster struct {
	ClusterName string      `json:"cluster_name"`
	Id          string      `json:"id"`
	ClusterType ClusterType `json:"cluster_type"`
	NodeCount   int         `json:"node_count"`
	Endpoint    []*Endpoint `json:"endpoint"`

	SetTime    string `json:"set_time"`
	UpdateTime string `json:"update_time"`
}

func (c *Cluster) Unmarshal(value *Value) error {
	c.ClusterName = value.GetString("cluster_name")
	c.ClusterType = ClusterType(value.GetString("cluster_type"))
	c.Id = value.GetString("id")
	c.SetTime = value.GetString("set_time")
	c.UpdateTime = value.GetString("update_time")

	if list, err := UnmarshalObjectList(value.Get("endpoint"), &Endpoint{}); err == nil {
		c.Endpoint = list
		c.NodeCount = len(c.Endpoint)
	} else {
		c.NodeCount = 0
		return err
	}

	return nil
}

type Endpoint struct {
	Host string `json:"host"`
	Demo string `json:"demo"`
}

func (n *Endpoint) Unmarshal(value *Value) error {
	n.Host = value.GetString("host")
	n.Demo = value.GetString("demo")
	return nil
}

type Record struct {
	Count   float64 `json:"count"`
	Mean    float64 `json:"mean"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	TimeStr string  `json:"time_str"`
}

func (r *Record) Unmarshal(value *Value) error {
	r.Count = value.GetFloat64("count")
	r.Mean = value.GetFloat64("mean")
	r.Max = value.GetFloat64("max")
	r.Min = value.GetFloat64("min")
	r.TimeStr = value.GetString("time_str")

	return nil
}

type RecordsData struct {
	Time        int64                             `json:"time"`
	MetricsData map[string]map[string]interface{} `json:"metrics_data"`
}

type VisitMetrics struct {
	SerName     string `json:"ser_name"`
	Total       int64  `json:"total"`
	Ok          int64  `json:"ok"`
	Lose        int64  `json:"lose"`
	HostLose    int64  `json:"host_lose"`
	ErrorLose   int64  `json:"error_lose"`
	PathLose    int64  `json:"path_lose"`
	TargetLose  int64  `json:"target_lose"`
	BlackIpLose int64  `json:"blackip_lose"`
	CertTotal   int64  `json:"cert_total"`
	CertDown    int64  `json:"cert_down"`
	CertFail    int64  `json:"cert_fail"`
}

func (v *VisitMetrics) Unmarshal(value *Value) error {
	v.Total = value.GetInt64("total")
	v.Ok = value.GetInt64("ok")
	v.Lose = value.GetInt64("lose")
	v.HostLose = value.GetInt64("host_lose")
	v.ErrorLose = value.GetInt64("error_lose")
	v.PathLose = value.GetInt64("path_lose")
	v.TargetLose = value.GetInt64("target_lose")
	v.BlackIpLose = value.GetInt64("blackip_lose")
	v.CertTotal = value.GetInt64("cert_total")
	v.CertDown = value.GetInt64("cert_down")
	v.CertFail = value.GetInt64("cert_fail")

	return nil
}

type SummaryMetrics struct {
	Key      string       `json:"key"`
	Increase VisitMetrics `json:"increase"`
	Current  VisitMetrics `json:"current"`
}

func (s *SummaryMetrics) Unmarshal(value *Value) error {
	s.Key = value.GetString("key")

	Unmarshal(value.Get("increase"), &s.Increase)
	Unmarshal(value.Get("current"), &s.Current)

	return nil
}

type TimePoint struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

type MetricsPoint struct {
	Value *SummaryMetrics `json:"summary"`
	Point TimePoint       `json:"point"`
}

func GenerateXAxis(t time.Time, dimension, pointCount int) []TimePoint {

	var rtn []TimePoint

	for idx := 0; idx < pointCount; idx++ {
		var s, v string

		switch dimension {
		case 1:
			t1 := t.Add(time.Duration(-idx) * time.Minute)
			v = t1.Format("2006/01/02 15:04")
			s = v[len(v)-5:]
		case 15:
			t1 := t.Add(time.Duration(-idx) * 15 * time.Minute)
			v = t1.Format("2006/01/02 15:") + strconv.Itoa((t1.Minute()+16)/15) + "/4"
			s = v[len(v)-6:]
			//t1 = t1.Add(15*time.Minute - 1)

			m := (t1.Minute()+16)/15 - 1
			s = t1.Format("15:") + util4go.LeftPad(strconv.Itoa(m*15), 2, '0')

			//m := ((t1.Minute() + 16) / 15) * 15
			//if m == 0 {
			//	s = t1.Add(time.Hour).Format("15:00")
			//} else {
			//	s = t1.Format("15:") + util4go.LeftPad(strconv.Itoa(m*15), 2, '0')
			//}
		case 30:
			t1 := t.Add(time.Duration(-idx) * 30 * time.Minute)
			v = t1.Format("2006/01/02 15:") + strconv.Itoa((t1.Minute()+31)/30) + "/2"
			s = v[len(v)-6:]

			m := (t1.Minute()+31)/30 - 1
			s = t1.Format("15:") + util4go.LeftPad(strconv.Itoa(m*30), 2, '0')
			//
			//m := ((t1.Minute() + 31) / 30) % 2
			//if m == 0 {
			//	s = t1.Add(time.Hour).Format("15:00")
			//} else {
			//	s = t1.Format("15:") + util4go.LeftPad(strconv.Itoa(m*30), 2, '0')
			//}
		case 60:
			t1 := t.Add(time.Duration(-idx) * time.Hour)
			v = t1.Format("2006/01/02 15")
			s = v[len(v)-3:] + ":00"
		default:
			t1 := t.Add(time.Duration(-idx) * 24 * time.Hour)
			v = t1.Format("2006/01/02")
			s = v
		}

		rtn = append(rtn, TimePoint{s, v})
	}

	return rtn
}
