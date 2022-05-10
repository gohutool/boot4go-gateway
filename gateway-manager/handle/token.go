package handle

import (
	. "gateway-manager/util"
	. "github.com/gohutool/boot4go-util/http"
	. "github.com/gohutool/boot4go-util/jwt"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : token.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/24 12:58
* 修改历史 : 1. [2022/4/24 12:58] 创建文件 by LongYong
*/

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

func UserIdSubjectDecode(subject string) (any, error) {
	return subject, nil
}

func TokenInterceptorHandler(handler routing.Handler) routing.Handler {
	return func(context *routing.Context) error {
		token := GetToken(context.RequestCtx)

		if userid, err := CheckToken(Issuer, token, UserIdSubjectDecode); err == nil {
			Logger.Debug("UserId =============== %v", userid)
			SetUserId(context, userid.(string))
			if handler != nil {
				return handler(context)
			} else {
				return nil
			}
		} else {
			Logger.Debug("Token is timeout %v", token)
			context.RequestCtx.Response.Header.Set("is-session-timeout", "1")
			panic("No authorization")
		}
	}
}
