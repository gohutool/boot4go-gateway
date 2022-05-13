package util

import (
	"expvar"
	"fmt"

	prometheusfasthttp "github.com/gohutool/boot4go-prometheus/fasthttp"
	. "github.com/gohutool/boot4go-util/http"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
	"net"
	"net/http"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : fasthttp.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/23 23:52
* 修改历史 : 1. [2022/4/23 23:52] 创建文件 by LongYong
*/

// Various counters - see https://golang.org/pkg/expvar/ for details.
var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)

func StartHttpServer(listener net.Listener, router *routing.Router) {

	fs := &fasthttp.FS{
		Root:               "./html",
		IndexNames:         []string{"index.html", "index.hml"},
		GenerateIndexPages: true,
		Compress:           false,
		AcceptByteRange:    false,
		PathNotFound: func(ctx *fasthttp.RequestCtx) {
			ctx.Response.Header.SetContentType("application/json;charset=utf-8")
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write([]byte(Result.Fail(fmt.Sprintf("Page Not Found, %v %v", string(ctx.Method()), string(ctx.RequestURI()))).Json()))
		},
	}

	fsHandler := fs.NewRequestHandler()

	router.Get("/stats", func(ctx *routing.Context) error {
		expvarhandler.ExpvarHandler(ctx.RequestCtx)
		return nil
	})

	router.Any("/*", func(context *routing.Context) error {
		ctx := context.RequestCtx
		fsHandler(ctx)
		UpdateFSCounters(ctx)
		return nil
	})

	requestHandler := func(ctx *fasthttp.RequestCtx) {

		Logger.Debug("%v %v %v %v", string(ctx.Path()), ctx.URI().String(), string(ctx.Method()), ctx.QueryArgs().String())
		defer func() {
			if err := recover(); err != nil {
				Logger.Debug(err)
				// ctx.Error(fmt.Sprintf("%v", err), http.StatusInternalServerError)
				Error(ctx, Result.Fail(fmt.Sprintf("%v", err)).Json(), http.StatusInternalServerError)
			}

			ctx.Response.Header.Set("tick", time.Now().String())
			ctx.Response.Header.SetServer("Gateway-UIManager")

			prometheusfasthttp.RequestCounterHandler(nil)(ctx)

			Logger.Debug("router.HandleRequest is finish")

		}()

		router.HandleRequest(ctx)
	}

	// Start HTTP server.
	Logger.Info("Starting HTTP server on %v", listener.Addr().String())
	go func() {
		if err := fasthttp.Serve(listener, requestHandler); err != nil {
			Logger.Critical("error in ListenAndServe: %v", err)
		}
	}()
}

func UpdateFSCounters(ctx *fasthttp.RequestCtx) {
	// Increment the number of fsHandler calls.
	fsCalls.Add(1)

	// Update other stats counters
	resp := &ctx.Response
	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		fsOKResponses.Add(1)
		fsResponseBodyBytes.Add(int64(resp.Header.ContentLength()))
	case fasthttp.StatusNotModified:
		fsNotModifiedResponses.Add(1)
	case fasthttp.StatusNotFound:
		fsNotFoundResponses.Add(1)
	default:
		fsOtherResponses.Add(1)
	}
}

func GetUserId(context *routing.Context) string {
	return context.Get("userid").(string)
}

func SetUserId(context *routing.Context, userid string) {
	context.Set("userid", userid)
}
