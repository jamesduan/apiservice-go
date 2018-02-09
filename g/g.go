package g

import (
	"log"
	"runtime"
)

// change log:
// 1.0.7: code refactor for open source
// 1.0.8: bugfix loop init cache
// 1.0.9: update host table anyway
// 1.1.0: remove Checksum when query plugins
// 1.1.1: 1.添加下载功能; 2.心跳周期修正，及时更新timestamp;
// 1.1.2: 1.增加agent信息上报接口、获取可执行plugin列表接口、同步命令接口、上报命令执行结果接口
// 1.1.3: 1.增加新版的下载功能
// 1.1.4: 支持plugin执行状态写入kafka
// 1.1.5: agent信息和命令执行结果修改为写到kafka

const (
	VERSION = "0.0.1"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
