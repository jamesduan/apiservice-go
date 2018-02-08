package sender

import (
	"apiservice/g"

	nlist "github.com/toolkits/container/list"
)

//DefaultSendQueueMaxSize 队列的最大长度
const (
	DefaultSendQueueMaxSize = 102400 //10.24w
)

//KafkaQueues initial queue
var (
	KafkaQueues = make(map[string]*nlist.SafeListLimited)
)

func initSendQueues() {
	cfg := g.Config()

	for node := range cfg.Kafka.Cluster {
		Q := nlist.NewSafeListLimited(DefaultSendQueueMaxSize)
		KafkaQueues[node] = Q
	}
}
