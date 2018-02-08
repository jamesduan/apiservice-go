// Copyright 2017 JingDong, Inc.
// Written by Zhangyunyang 2017/10/19
//
// Kafka backends connection pool
//

package backend_pool

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	connp "github.com/toolkits/conn_pool"
)

// implement NConn
type kafkaClient struct {
	prd          *kafka.Producer
	deliveryChan chan kafka.Event
	name         string
}

// ConnPools Manager
type KafkaConnPools struct {
	sync.RWMutex
	M           map[string]*connp.ConnPool
	MaxConns    int
	MaxIdle     int
	ConnTimeout int
	CallTimeout int
}

//kafka client
func newKafkaClient(producer *kafka.Producer, delivery_chan chan kafka.Event, name string) *kafkaClient {
	return &kafkaClient{prd: producer, deliveryChan: delivery_chan, name: name}
}

func (this *kafkaClient) Name() string {
	return this.name
}

func (this *kafkaClient) Closed() bool {
	return this.prd == nil
}

func (this *kafkaClient) Close() error {
	if this.prd != nil {
		this.prd.Close()
		close(this.deliveryChan)
		this.prd = nil
		return nil
	}
	return nil
}

//KafkaConnPools
func CreateKafkaConnPools(maxConns, maxIdle, connTimeout, callTimeout int, cluster []string) *KafkaConnPools {
	cp := &KafkaConnPools{M: make(map[string]*connp.ConnPool), MaxConns: maxConns, MaxIdle: maxIdle,
		ConnTimeout: connTimeout, CallTimeout: callTimeout}

	ct := time.Duration(cp.ConnTimeout) * time.Millisecond
	for _, address := range cluster {
		if _, exist := cp.M[address]; exist {
			continue
		}
		cp.M[address] = createOneKafkaPool(address, address, ct, maxConns, maxIdle)
	}

	return cp
}

// sync
func (this *KafkaConnPools) Call(addr string, msg *kafka.Message) error {
	connPool, exists := this.Get(addr)
	if !exists {
		return fmt.Errorf("%s has no connection pool", addr)
	}

	conn, err := connPool.Fetch()
	if err != nil {
		return fmt.Errorf("%s get connection fail: conn %v, err %v. proc: %s", addr, conn, err, connPool.Proc())
	}

	kafka_client := conn.(*kafkaClient)
	//callTimeout := time.Duration(this.CallTimeout) * time.Millisecond

	done := make(chan error, 1)
	go func() {
		err := kafka_client.prd.Produce(msg, kafka_client.deliveryChan)
		if err != nil {
			log.Println("Producer send message failed:", err)
			done <- err
			return
		}

		e := <-kafka_client.deliveryChan
		m := e.(*kafka.Message)
		err = m.TopicPartition.Error

		done <- err
	}()

	select {
	// take care
	//case <-time.After(callTimeout):
	//	connPool.ForceClose(conn)
	//	return fmt.Errorf("%s, call timeout", addr)
	case err = <-done:
		if err != nil {
			connPool.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", addr, err, connPool.Proc())
		} else {
			connPool.Release(conn)
		}
		return err
	}
}

func (this *KafkaConnPools) Get(address string) (*connp.ConnPool, bool) {
	this.RLock()
	defer this.RUnlock()
	p, exists := this.M[address]
	return p, exists
}

func (this *KafkaConnPools) Destroy() {
	this.Lock()
	defer this.Unlock()
	addresses := make([]string, 0, len(this.M))
	for address := range this.M {
		addresses = append(addresses, address)
	}

	for _, address := range addresses {
		this.M[address].Destroy()
		delete(this.M, address)
	}
}

func (this *KafkaConnPools) Proc() []string {
	procs := []string{}
	for _, cp := range this.M {
		procs = append(procs, cp.Proc())
	}
	return procs
}

// define New
func createOneKafkaPool(name string, address string, connTimeout time.Duration, maxConns int, maxIdle int) *connp.ConnPool {
	p := connp.NewConnPool(name, address, int32(maxConns), int32(maxIdle))
	p.New = func(connName string) (connp.NConn, error) {

		producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": p.Address})
		if err != nil {
			log.Println("NewProducer failed:", err)
			return nil, err
		}

		delivery_chan := make(chan kafka.Event)

		return newKafkaClient(producer, delivery_chan, connName), nil
	}

	return p
}
