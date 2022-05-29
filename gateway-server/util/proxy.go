package util

import (
	"fmt"
	. "gateway-server/model"
	. "github.com/gohutool/boot4go-pathmatcher"
	util4go "github.com/gohutool/boot4go-util"
	"github.com/valyala/fasthttp"
	"net/http"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : proxy.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/4 21:19
* 修改历史 : 1. [2022/5/4 21:19] 创建文件 by LongYong
*/

func InitProxyHandler() fasthttp.RequestHandler {

	return RecoveryHandler(func(ctx *fasthttp.RequestCtx) {
		//ctx.Write([]byte("This is a test"))
		Logger.Debug("%v", ctx.URI().String())

		host := string(ctx.Host())
		if domain, ok := DomainPool.GetDomain(host); ok {

			BlackIpHandler(domain.BlackIps, nil)(ctx)

			var path *Path

			uri := string(ctx.URI().Path())

			for _, p := range domain.Path {
				ok, _ := Match(p.ReqPath, uri)
				if ok {
					path = &p
					break
				}
			}

			if path == nil {
				MetricsVisitPathLose()
				panic(uri + " Not found mapped path config")
			}

			var destUri string

			if util4go.IsEmpty(path.ReqPath) && util4go.IsEmpty(path.ReplacePath) {
				destUri = uri
			} else if util4go.IsEmpty(path.ReqPath) {
				destUri = path.ReplacePath
			} else {
				toPath, err := util4go.RegExpPool.ConvertRegExpWithFormat(uri, path.SearchPath, path.ReplacePath)

				if err != nil {
					MetricsVisitPathLose()
					panic("Not found destination path")
				}
				destUri = toPath
			}

			if len(path.Targets) == 0 {
				MetricsVisitTargetLose()
				panic(destUri + " Not found target host")
			}

			var t *Target
			if len(path.Targets) > 1 && path.LB != nil {
				if t2, err := (*path.LB).Target(); err != nil {
					MetricsVisitTargetLose()
					panic(destUri + " Load Balance not found target host")
				} else {
					t = t2
				}
			} else {
				t = path.Targets[0]
			}

			s, h, q := t.Schema, t.Host, t.Query

			if err := Proxy(h, s, destUri, q, ctx, path.CircuitBreakerTimeout); err != nil {
				Logger.Debug("Proxy error : %v", err)
				MetricsVisitError()
			} else {
				MetricsVisitOK()
			}
		} else {
			MetricsVisitHostLose()
			panic(host + " not found mapped domain config")
		}

	})
}

func RecoveryHandler(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Error(fmt.Sprintf("Gateway error %v", err), http.StatusBadGateway)
				Logger.Error("Gateway error %v", err)
			}
		}()

		if next != nil {
			next(ctx)
		}
	}
}

func BlackIpHandler(ips map[string]bool, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ip := realIP(ctx)

		if _, ok := ips[ip]; ok {
			MetricsVisitBlackIP()
			panic("Not allowed as black ip forbidden")
		}

		if next != nil {
			next(ctx)
		}
	}
}
