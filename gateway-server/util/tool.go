package util

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"strconv"
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

func GetMachineData() map[string]interface{} {
	data := make(map[string]interface{})
	mem1, _ := mem.VirtualMemory()
	data["Mem"] = mem1.String()
	avg, _ := load.Avg()
	data["Avg"] = avg.String()

	cpu1, _ := cpu.Counts(false)
	data["Cpu1"] = cpu1
	cpu2, _ := cpu.Counts(true)
	data["Cpu2"] = cpu2

	infos, _ := cpu.Info()
	data["CpuInfo"] = infos

	timestamp, _ := host.BootTime()
	t := time.Unix(int64(timestamp), 0)
	data["Up"] = t.Local().Format("2006-01-02 15:04:05")
	uptime, _ := host.BootTime()
	data["Uptime"] = uptime

	version, _ := host.KernelVersion()
	data["kernel"] = version

	platform, family, version, _ := host.PlatformInformation()
	data["platform"] = platform
	data["family"] = family
	data["version"] = version

	return data
}

type Dimension struct {
	M1  string
	M15 string
	M30 string
	M60 string
	MD  string

	LM1  string
	LM15 string
	LM30 string
	LM60 string
	LMD  string
}

func Time2Dimension(t time.Time) *Dimension {
	d := &Dimension{}

	d.M1 = t.Format("2006/01/02 15:04")
	d.M15 = t.Format("2006/01/02 15:") + strconv.Itoa((t.Minute()+16)/15) + "/4"
	d.M30 = t.Format("2006/01/02 15:") + strconv.Itoa((t.Minute()+31)/30) + "/2"
	d.M60 = t.Format("2006/01/02 15")
	d.MD = t.Format("2006/01/02")

	d.LM1 = t.Add(-1 * time.Minute).Format("2006/01/02 15:04")
	t15 := t.Add(-15 * time.Minute)
	d.LM15 = t15.Format("2006/01/02 15:") + strconv.Itoa((t15.Minute()+16)/15) + "/4"
	t30 := t.Add(-30 * time.Minute)
	d.LM30 = t30.Format("2006/01/02 15:") + strconv.Itoa((t30.Minute()+31)/30) + "/2"
	t60 := t.Add(-60 * time.Minute)
	d.LM60 = t60.Format("2006/01/02 15")
	td := t.Add(-24 * time.Hour)
	d.LMD = td.Format("2006/01/02")

	return d
}
