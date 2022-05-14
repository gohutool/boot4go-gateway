package model

import (
	. "github.com/gohutool/boot4go-fastjson"
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
	Id                 string          `json:"id"`
	DomainName         string          `json:"domain_name"`
	DomainUrl          string          `json:"domain_url"`
	LbType             string          `json:"lb_type"`
	Targets            []*Target       `json:"targets"`
	BlackIps           map[string]bool `json:"black_ips"`
	RateLimiterNum     float64         `json:"rate_limiter_num"`
	RateLimiterMsg     string          `json:"rate_limiter_msg"`
	RateLimiterEnabled bool            `json:"rate_limiter_enabled"`
	SetTime            string          `json:"set_time"`
	UpdateTime         string          `json:"update_time"`
	Demo               string          `json:"demo"`

	SslOn   bool `json:"ssl_on"`
	SslPort int  `json:"ssl_port"`

	Path []Path
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

	d.BlackIps = make(map[string]bool)

	if list, err := UnmarshalStringList(value.Get("black_ips")); err == nil {
		for _, one := range list {
			d.BlackIps[one] = true
		}
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

	t.Schema = value.GetString("schema")
	t.Host = value.GetString("host")
	t.Query = value.GetString("query")

	return nil
}

type Path struct {
	Id                    string    `json:"id"`
	DomainId              string    `json:"domain_id"`
	DomainUrl             string    `json:"domain_url"`
	ReqMethod             string    `json:"req_method"`
	ReqPath               string    `json:"req_path"`
	SearchPath            string    `json:"search_path"`
	ReplacePath           string    `json:"replace_path"`
	CircuitBreakerRequest int       `json:"circuit_breaker_request"`
	CircuitBreakerPercent int       `json:"circuit_breaker_percent"`
	CircuitBreakerTimeout int       `json:"circuit_breaker_timeout"`
	CircuitBreakerMsg     string    `json:"circuit_breaker_msg"`
	CircuitBreakerEnabled bool      `json:"circuit_breaker_enabled"`
	CircuitBreakerForce   bool      `json:"circuit_breaker_force"`
	PrivateProxyEnabled   bool      `json:"private_proxy_enabled"`
	LbType                string    `json:"lb_type"`
	Targets               []*Target `json:"targets"`
	SetTime               string    `json:"set_time"`
	UpdateTime            string    `json:"update_time"`

	LB *LoadBalance
}

func (p *Path) Unmarshal(value *Value) error {
	p.Id = value.GetString("id")
	p.DomainId = value.GetString("domain_id")
	p.ReqMethod = value.GetString("req_method")
	p.ReqPath = value.GetString("req_path")
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

	p.DomainUrl = value.GetString("domain_url")

	p.SetTime = value.GetString("set_time")
	p.SetTime = value.GetString("update_time")

	if list, err := UnmarshalObjectList(value.Get("targets"), &Target{}); err == nil {
		p.Targets = list
	} else {
		return err
	}

	if len(p.Targets) > 1 {
		p.LB = new(LoadBalance)

		if p.LbType == "random" {
			*p.LB = NewRandom(p.Targets, time.Now().Unix())
		} else if p.LbType == "roundRobin" {
			*p.LB = NewRoundRobin(p.Targets)
		}

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
	UpTime     string `json:"up_time"`
}

func (g *Gateway) Unmarshal(value *Value) error {
	g.ServerName = value.GetString("server_name")

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

type SummaryMetrics struct {
	Key      string       `json:"key"`
	LeaseID  int64        `json:"lease_id"`
	Increase VisitMetrics `json:"increase"`
	Current  VisitMetrics `json:"current"`
}

func (s *SummaryMetrics) Unmarshal(value *Value) error {
	s.Key = value.GetString("key")
	s.LeaseID = value.GetInt64("lease_id")

	Unmarshal(value.Get("increase"), &s.Increase)
	Unmarshal(value.Get("current"), &s.Current)

	return nil
}

type VisitMetrics struct {
	Total       int64 `json:"total"`
	Ok          int64 `json:"ok"`
	Lose        int64 `json:"lose"`
	HostLose    int64 `json:"host_lose"`
	ErrorLose   int64 `json:"error_lose"`
	PathLose    int64 `json:"path_lose"`
	TargetLose  int64 `json:"target_lose"`
	BlackIpLose int64 `json:"blackip_lose"`
	CertTotal   int64 `json:"cert_total"`
	CertDown    int64 `json:"cert_down"`
	CertFail    int64 `json:"cert_fail"`
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
