package util

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : ratelimiter.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/5 21:46
* 修改历史 : 1. [2022/5/5 21:46] 创建文件 by LongYong
*/

func Test() {
	// You can create a generic limiter for all your handlers
	// or one for each handler. Your choice.
	// This limiter basically says: allow at most 1 request per 1 second.
	lim := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{})

	// This is an example on how to limit only GET and POST requests.
	lim.SetMethods([]string{"GET", "POST"})
}
