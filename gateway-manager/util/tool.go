package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : tool.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/24 14:02
* 修改历史 : 1. [2022/4/24 14:02] 创建文件 by LongYong
*/

const DEFAULT_TOKEN_EXPIRE = 24 * time.Hour
const DEFAULT_ISSUER = "GATEWAY-UIMANAGER"

var TokenExpire = DEFAULT_TOKEN_EXPIRE
var Issuer = DEFAULT_ISSUER

var menuObj = make(map[string]any)

func init() {
	menuObj = loadMenu("menus.json")
}

func loadMenu(filename string) map[string]any {
	xc := make(map[string]any)

	fd, err := os.Open(filename)
	if err != nil {
		Logger.Debug("LoadMenu: Error: Could not open %q for reading: %s", filename, err)
		return xc
	}
	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		Logger.Debug("LoadMenu: Error: Could not read %q: %s", filename, err)
		return xc
	}

	if err := json.Unmarshal(contents, &xc); err != nil {
		msg := fmt.Sprintf("LoadMenu: Error: Could not parse Json configuration in %q: %s", filename, err)
		Logger.Debug(msg)
		return xc
	}

	Logger.Debug("Menus load ok %+v", xc)

	return xc
}

func GetMenuObj() map[string]any {
	xc := make(map[string]any)

	if b, err := json.Marshal(menuObj); err != nil {
		Logger.Debug("%v", err)
		return xc
	} else {
		err = json.Unmarshal(b, &xc)
		if err == nil {
			return xc
		} else {
			Logger.Debug("%v", err)
			return xc
		}
	}

}
