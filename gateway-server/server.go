package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"gateway-server/handle"
	. "gateway-server/util"
	"github.com/alecthomas/kingpin"
	. "github.com/gohutool/boot4go-etcd/client"
	"github.com/valyala/fasthttp"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : server.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/1 16:35
* 修改历史 : 1. [2022/5/1 16:35] 创建文件 by LongYong
*/

const (
	SERVER_VERSION = "Gateway4go-server-v1.0.0"
	SERVER_MAJOR   = 1
	SERVER_MINOR   = 0
	SERVER_BUILD   = 0
)

func main() {
	app := kingpin.New("Gateway-Manager", "A gateway server.")
	addr_flag := app.Flag("addr", "Addr: gateway listen addr.").Short('l').Default(":9000").String()
	tlsAddr := app.Flag("tls-addr", "Tls-Addr: gateway tls listen addr").Default(":9443").String()
	etcd_flag := app.Flag("etcd", "Etcd: etcd server addr.").Default("192.168.56.101:32379").Short('n').String()
	username_flag := app.Flag("username", "Username: etcd username.").Short('u').Default("").String()
	password_flag := app.Flag("password", "Password: etcd password.").Short('p').Default("").String()
	monitor_flag := app.Flag("monitor", "Monitor: Gateway Manager can monitor.").Short('m').Default("false").Bool()
	Metrics_Enable = *monitor_flag

	app.HelpFlag.Short('h')
	app.Version(SERVER_VERSION)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	l, err := net.Listen("tcp", *addr_flag)
	if err != nil {
		panic("Start server error " + err.Error())
	}
	//
	//Logger.Info("Server is listen on %v ", *addr_flag)

	err = EtcdClient.Init([]string{*etcd_flag}, *username_flag, *password_flag, DialTimeout)

	if err == nil {
		fmt.Println("Etcd is connect")
	} else {
		panic("Etcd can not connect")
	}

	fmt.Println("Start server now .... ")

	initPool()
	registerGateWay(*addr_flag)
	startMonitor(*monitor_flag)

	requestHandler := PrometheusRequestHandler(fasthttp.RequestHandler(InitProxyHandler()))

	// Start HTTP server.
	Logger.Info("Server is listen on %v %v", *addr_flag, l.Addr().String())
	go func() {
		if err := fasthttp.Serve(l, requestHandler); err != nil {
			Logger.Critical("error in ListenAndServe: %v", err)
		}
	}()

	if *tlsAddr != "" {
		ls, err := net.Listen("tcp", *tlsAddr)

		Logger.Info("Tls Server is listen on %v %v", *tlsAddr, ls.Addr())

		if err != nil {
			panic(fmt.Sprintln("tls ", *tlsAddr, ", 监听失败", err))
		}

		tlsListener := tls.NewListener(ls,
			&tls.Config{GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {

				host, _, _ := net.SplitHostPort(info.Conn.LocalAddr().String())
				if host == "[::1]" || host == "::1" || host == "[::]" || host == "::" {
					host = "localhost"
				}

				Logger.Debug("SSL(%v) %v <-> %v ", host, info.Conn.LocalAddr(), info.Conn.RemoteAddr())

				//return nil, errors.New("tls: no certificates configured")
				c := CertPool.GetCert(host)

				if c == nil {
					MetricsCertFail()
					return nil, errors.New("tls: no certificates configured")
				} else {
					MetricsCertDown()
					return c, nil
				}
			}})

		go func() {
			if err := fasthttp.Serve(tlsListener, requestHandler); err != nil {
				Logger.Critical("error in ListenAndServe(Tls): %v", err)
			}
		}()
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(twg *sync.WaitGroup) {
		sig := make(chan os.Signal, 2)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		fmt.Println("signal service close")
		twg.Done()
	}(wg)

	wg.Wait()

	unRegisterGateWay()

	l.Close()
	fmt.Println("Server is close")
}

func initPool() {
	DomainPool.InitDomainPool()
	CertPool.InitCertPool()
}

func registerGateWay(address string) {

	handle.WatchEventListener.InitMachineDataNotify(address)

	handle.WatchEventListener.InitDomainChangeListener()
	handle.WatchEventListener.InitDomainPathChangeListener()

	handle.WatchEventListener.InitCertChangeListener()

	Logger.Info("Registry the instance into etcd server.")
}

func unRegisterGateWay() {

	Logger.Info("Unregistry the instance from etcd server.")
}

func startMonitor(monitor bool) {
	handle.WatchEventListener.InitMetricsScheduleJob(monitor)
}
