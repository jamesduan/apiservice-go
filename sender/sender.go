package sender

import (
	"apiservice/common/model"
	"apiservice/g"
	"apiservice/proc"
	"encoding/json"
	// "happy-hbs/common/model"
	// "happy-hbs/modules/hbs/g"
	// "happy-hbs/modules/hbs/proc"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	nsema "github.com/toolkits/concurrent/semaphore"
	"github.com/toolkits/container/list"
)

// 默认参数
var (
	MinStep int //最小上报周期,单位sec
)

const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

func Start() {
	MinStep = g.Config().MinStep
	if MinStep < 1 {
		MinStep = 30
	}
	initKafkaConnPools()
	log.Println("Initialize Kafka connection pools done.")
	initSendQueues()
	log.Println("Initialize send queue done.")
	initNodeRings()
	log.Println("Initialize consistent hash node ring done.")
	startSendTasks()

	startSenderCron()
	log.Println("send.Start, ok")
}

func alignTs(ts int64, period int64) int64 {
	return ts - ts%period
}

func Push2KafkaSendQueue(item *model.KafkaItem) error {
	t := item.Type
	var pk string

	switch t {
	case model.PluginSatusType:
		pk = item.Body.(*model.PluginReportStatus).PK()
	case model.AgentInfoType:
		pk = item.Body.(*model.AgentInfoRequest).PK()
	case model.CmdResultType:
		pk = item.Body.(*model.PluginCmdResultRequest).PK()
	default:
		log.Println("Invalid kafka message type:", t)
		return nil
	}

	node, err := KafkaNodeRing.GetNode(pk)
	if err != nil {
		log.Println("E:", err)
		return err
	}
	Q := KafkaQueues[node]
	isSuccess := Q.PushFront(item)
	if !isSuccess {
		proc.SendToKafkaDropCnt.Incr()
	}

	return nil
}

func startSendTasks() {
	cfg := g.Config()
	// init semaphore
	kafkaConcurrent := cfg.Kafka.MaxConns

	if kafkaConcurrent < 1 {
		kafkaConcurrent = 1
	}

	// init send go-routines
	for node := range cfg.Kafka.Cluster {
		q := KafkaQueues[node]
		go sendToKafka(q, node, kafkaConcurrent)
	}
}

func sendToKafka(Q *list.SafeListLimited, node string, concurrent int) {

	cfg := g.Config()

	batch := cfg.Kafka.Batch
	pluginTopic, exist := cfg.Kafka.Topics["plugin_exec_state"]
	if !exist {
		log.Println("kafka plugin status topic not exist.")
		return
	}

	agentTopic, exist := cfg.Kafka.Topics["agent_info"]
	if !exist {
		log.Println("kafka agent info topic not exist.")
		return
	}

	cmdTopic, exist := cfg.Kafka.Topics["cmd_result"]
	if !exist {
		log.Println("kafka cmd result topic not exist.")
		return
	}

	addr := cfg.Kafka.Cluster[node]
	sema := nsema.NewSemaphore(concurrent)
	log.Println("kafka node addr:", addr)

	for {
		items := Q.PopBackBy(batch)
		count := len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		for i := 0; i < count; i++ {
			t := items[i].(*model.KafkaItem).Type
			msg := &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &pluginTopic, Partition: kafka.PartitionAny}}
			switch t {
			case model.PluginSatusType:
				msg.TopicPartition.Topic = &pluginTopic
			case model.AgentInfoType:
				msg.TopicPartition.Topic = &agentTopic
			case model.CmdResultType:
				msg.TopicPartition.Topic = &cmdTopic
			default:
				log.Println("Invalid kafka message type:", t)
				continue
			}

			body := items[i].(*model.KafkaItem).Body
			msgContent, err := json.Marshal(body)
			if err != nil {
				log.Printf("json.Marshal() failed, body:%v, error:%v\n", body, err)
				continue
			}

			msg.Value = msgContent

			sema.Acquire()
			go func(addr string, msg *kafka.Message, count int) {
				defer sema.Release()

				var err error
				sendOk := false
				for i := 0; i < 3; i++ {
					err = KafkaConnPools.Call(addr, msg)
					if err == nil {
						sendOk = true
						break
					}
					time.Sleep(time.Millisecond * 10)
				}

				// statistics
				if !sendOk {
					log.Printf("send kafka %s:%s fail: %v", node, addr, err)
					proc.SendToKafkaFailCnt.IncrBy(int64(count))
				} else {
					proc.SendToKafkaCnt.IncrBy(int64(count))
				}
			}(addr, msg, 1)
		}

	}
}
