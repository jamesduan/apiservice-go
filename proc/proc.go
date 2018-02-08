package proc

import (
	"log"

	nproc "github.com/toolkits/proc"
)

// trace
var (
	RecvDataTrace = nproc.NewDataTrace("RecvDataTrace", 3)
)

// filter
var (
	RecvDataFilter = nproc.NewDataFilter("RecvDataFilter", 5)
)

// 统计指标的整体数据
var (
	// 计数统计,正确计数,错误计数, ...
	RecvCnt       = nproc.NewSCounterQps("RecvCnt")
	RpcRecvCnt    = nproc.NewSCounterQps("RpcRecvCnt")
	HttpRecvCnt   = nproc.NewSCounterQps("HttpRecvCnt")
	SocketRecvCnt = nproc.NewSCounterQps("SocketRecvCnt")

	// SendToJudgeCnt                   = nproc.NewSCounterQps("SendToJudgeCnt")
	// SendToTsdbCnt                    = nproc.NewSCounterQps("SendToTsdbCnt")
	// SendToGraphCnt                   = nproc.NewSCounterQps("SendToGraphCnt")
	// SendToKafkaHwCnt                 = nproc.NewSCounterQps("SendToKafkaHwCnt")
	// SendToKafkaGpuCnt                = nproc.NewSCounterQps("SendToKafkaGpuCnt")
	SendToKafkaCnt = nproc.NewSCounterQps("SendToKafkaCnt")

	// SendToJudgeDropCnt                   = nproc.NewSCounterQps("SendToJudgeDropCnt")
	// SendToTsdbDropCnt                    = nproc.NewSCounterQps("SendToTsdbDropCnt")
	// SendToGraphDropCnt                   = nproc.NewSCounterQps("SendToGraphDropCnt")
	// SendToKafkaHwDropCnt                 = nproc.NewSCounterQps("SendToKafkaHwDropCnt")
	// SendToKafkaGpuDropCnt                = nproc.NewSCounterQps("SendToKafkaGpuDropCnt")
	SendToKafkaDropCnt = nproc.NewSCounterQps("SendToKafkaDropCnt")

	// SendToJudgeFailCnt    = nproc.NewSCounterQps("SendToJudgeFailCnt")
	// SendToTsdbFailCnt     = nproc.NewSCounterQps("SendToTsdbFailCnt")
	// SendToGraphFailCnt    = nproc.NewSCounterQps("SendToGraphFailCnt")
	// SendToKafkaHwFailCnt  = nproc.NewSCounterQps("SendToKafkaHwFailCnt")
	// SendToKafkaGpuFailCnt = nproc.NewSCounterQps("SendToKafkaGpuFailCnt")
	SendToKafkaFailCnt = nproc.NewSCounterQps("SendToKafkaFailCnt")

	// 发送缓存大小
	// JudgeQueuesCnt    = nproc.NewSCounterBase("JudgeSendCacheCnt")
	// TsdbQueuesCnt     = nproc.NewSCounterBase("TsdbSendCacheCnt")
	// GraphQueuesCnt    = nproc.NewSCounterBase("GraphSendCacheCnt")
	// KafkaHwQueuesCnt  = nproc.NewSCounterBase("KafkaHwSendCacheCnt")
	// KafkaGpuQueuesCnt = nproc.NewSCounterBase("KafkaGpuSendCacheCnt")
	KafkaQueuesCnt = nproc.NewSCounterBase("KafkaQueuesCnt")

	// http请求次数
	HistoryRequestCnt = nproc.NewSCounterQps("HistoryRequestCnt")
	InfoRequestCnt    = nproc.NewSCounterQps("InfoRequestCnt")
	LastRequestCnt    = nproc.NewSCounterQps("LastRequestCnt")
	LastRawRequestCnt = nproc.NewSCounterQps("LastRawRequestCnt")

	// http回执的监控数据条数
	HistoryResponseCounterCnt = nproc.NewSCounterQps("HistoryResponseCounterCnt")
	HistoryResponseItemCnt    = nproc.NewSCounterQps("HistoryResponseItemCnt")
	LastRequestItemCnt        = nproc.NewSCounterQps("LastRequestItemCnt")
	LastRawRequestItemCnt     = nproc.NewSCounterQps("LastRawRequestItemCnt")
)

//Start proc
func Start() {
	log.Println("proc.Start, ok")
}

//GetAll get all cnt
func GetAll() []interface{} {
	ret := make([]interface{}, 0)

	// recv cnt
	ret = append(ret, RecvCnt.Get())
	ret = append(ret, RpcRecvCnt.Get())
	ret = append(ret, HttpRecvCnt.Get())
	ret = append(ret, SocketRecvCnt.Get())

	// send cnt
	// ret = append(ret, SendToJudgeCnt.Get())
	// ret = append(ret, SendToTsdbCnt.Get())
	// ret = append(ret, SendToGraphCnt.Get())
	// ret = append(ret, SendToKafkaHwCnt.Get())
	// ret = append(ret, SendToKafkaGpuCnt.Get())
	ret = append(ret, SendToKafkaCnt.Get())

	// drop cnt
	// ret = append(ret, SendToJudgeDropCnt.Get())
	// ret = append(ret, SendToTsdbDropCnt.Get())
	// ret = append(ret, SendToGraphDropCnt.Get())
	// ret = append(ret, SendToKafkaHwDropCnt.Get())
	// ret = append(ret, SendToKafkaGpuDropCnt.Get())
	ret = append(ret, SendToKafkaDropCnt.Get())

	// send fail cnt
	// ret = append(ret, SendToJudgeFailCnt.Get())
	// ret = append(ret, SendToTsdbFailCnt.Get())
	// ret = append(ret, SendToGraphFailCnt.Get())
	// ret = append(ret, SendToKafkaHwFailCnt.Get())
	// ret = append(ret, SendToKafkaGpuFailCnt.Get())
	ret = append(ret, SendToKafkaFailCnt.Get())

	// cache cnt
	// ret = append(ret, JudgeQueuesCnt.Get())
	// ret = append(ret, TsdbQueuesCnt.Get())
	// ret = append(ret, GraphQueuesCnt.Get())
	// ret = append(ret, KafkaHwQueuesCnt.Get())
	// ret = append(ret, KafkaGpuQueuesCnt.Get())
	ret = append(ret, KafkaQueuesCnt.Get())

	// http request
	ret = append(ret, HistoryRequestCnt.Get())
	ret = append(ret, InfoRequestCnt.Get())
	ret = append(ret, LastRequestCnt.Get())
	ret = append(ret, LastRawRequestCnt.Get())

	// http response
	ret = append(ret, HistoryResponseCounterCnt.Get())
	ret = append(ret, HistoryResponseItemCnt.Get())
	ret = append(ret, LastRequestItemCnt.Get())
	ret = append(ret, LastRawRequestItemCnt.Get())

	return ret
}
