package main

import (
	"fmt"
	. "gateway-manager/handle"
	. "gateway-manager/util"
	"github.com/alecthomas/kingpin"
	. "github.com/gohutool/boot4go-etcd/client"
	routing "github.com/qiangxue/fasthttp-routing"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : manager.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/23 23:20
* 修改历史 : 1. [2022/4/23 23:20] 创建文件 by LongYong
*/

const (
	SERVER_VERSION = "Gateway4go-manager-v1.0.0"
	SERVER_MAJOR   = 1
	SERVER_MINOR   = 0
	SERVER_BUILD   = 0
)

func main() {
	app := kingpin.New("Gateway-Manager", "A gateway manager with UI.")
	addr_flag := app.Flag("addr", "Addr: gateway manager listen addr.").Short('l').Default(":9998").String()
	etcd_flag := app.Flag("etcd", "Etcd: etcd server addr.").Default("192.168.56.101:32379").Short('n').String()
	username_flag := app.Flag("username", "Username: etcd username.").Short('u').Default("").String()
	password_flag := app.Flag("password", "Password: etcd password.").Short('p').Default("").String()
	issuer_flag := app.Flag("issuer", "Issuer: token's issuer.").Short('i').Default(DEFAULT_ISSUER).String()
	expired_flag := app.Flag("token_expire", "Token_expire: many hour(s) token will expire.").Short('e').Default("24").Int()

	app.HelpFlag.Short('h')
	app.Version(SERVER_VERSION)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	l, err := net.Listen("tcp", *addr_flag)
	if err != nil {
		fmt.Println("Start server error " + err.Error())
		return
	}

	if issuer_flag != nil && len(*issuer_flag) > 0 {
		Issuer = *issuer_flag
	}

	if expired_flag != nil && *expired_flag > 0 {
		TokenExpire = time.Duration(*expired_flag) * time.Hour
	}

	err = EtcdClient.Init([]string{*etcd_flag}, *username_flag, *password_flag, DialTimeout)

	if err == nil {
		fmt.Println("Etcd is connect")
	} else {
		panic("Etcd can not connect")
	}

	if err = InitAdmin(); err != nil {
		panic("Init Admin User error " + err.Error())
	}

	fmt.Println("Start manager now .... ")

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(twg *sync.WaitGroup) {
		sig := make(chan os.Signal, 2)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		fmt.Println("signal service close")
		twg.Done()
	}(wg)

	router := routing.New()

	v3Group := router.Group("/v3/api")

	UserHandler.InitRouter(router, v3Group)
	DomainHandler.InitRouter(router, v3Group)
	PathHandler.InitRouter(router, v3Group)
	CertHandler.InitRouter(router, v3Group)
	GatewayHandler.InitRouter(router, v3Group)
	ClusterHandler.InitRouter(router, v3Group)

	PrometheusHandler.InitRouter(router, v3Group)

	//router.Any("/api/*", func(context *routing.Context) error {
	//	ctx := context.RequestCtx
	//	Logger.Debug("%v", string(ctx.Path()))
	//	UpdateFSCounters(ctx)
	//	result := result{}
	//	return result.Error("Test error")
	//	//return errors.New()
	//})

	StartHttpServer(l, router)

	Logger.Debug("%v %v %v %v", *addr_flag, *etcd_flag, *username_flag, *password_flag)

	wg.Wait()

	l.Close()
	fmt.Println("Server is close")
}
