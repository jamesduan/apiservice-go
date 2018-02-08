package sender

import (
	"apiservice/g"

	backend "apiservice/common/backend_pool"

	nset "github.com/toolkits/container/set"
)

var (
	KafkaConnPools *backend.KafkaConnPools
)

func initKafkaConnPools() {
	cfg := g.Config()

	kafkaInstances := nset.NewStringSet()
	for _, instance := range cfg.Kafka.Cluster {
		kafkaInstances.Add(instance)
	}
	KafkaConnPools = backend.CreateKafkaConnPools(cfg.Kafka.MaxConns, cfg.Kafka.MaxIdle,
		cfg.Kafka.ConnTimeout, cfg.Kafka.CallTimeout, kafkaInstances.ToSlice())

}

func DestroyConnPools() {
	KafkaConnPools.Destroy()
}
