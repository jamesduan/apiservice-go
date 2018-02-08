package sender

import (
	"apiservice/proc"
	"time"

	// "happy-hbs/modules/hbs/proc"

	"github.com/toolkits/container/list"
)

const (
	DefaultProcCronPeriod = time.Duration(5) * time.Second //ProcCron的周期,默认1s
)

// send_cron程序入口
func startSenderCron() {
	go startProcCron()
}

func startProcCron() {
	for {
		time.Sleep(DefaultProcCronPeriod)
		refreshSendingCacheSize()
	}
}

func refreshSendingCacheSize() {
	proc.KafkaQueuesCnt.SetCnt(calcSendCacheSize(KafkaQueues))
}
func calcSendCacheSize(mapList map[string]*list.SafeListLimited) int64 {
	var cnt int64 = 0
	for _, list := range mapList {
		if list != nil {
			cnt += int64(list.Len())
		}
	}
	return cnt
}
