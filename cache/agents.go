package cache

// 每个agent心跳上来的时候立马更新一下数据库是没必要的
// 缓存起来，每隔一个小时写一次DB
// 提供http接口查询机器信息，排查重名机器的时候比较有用

import (
	// "api_demo/common/model"
	"apiservice/common/model"
	"apiservice/db"
	"sync"
	"time"
	// "happy-hbs/common/model"
	// "happy-hbs/modules/hbs/db"
)

type SafeAgents struct {
	sync.RWMutex
	M map[string]*model.AgentUpdateInfo
}

var Agents = NewSafeAgents()

func NewSafeAgents() *SafeAgents {
	return &SafeAgents{M: make(map[string]*model.AgentUpdateInfo)}
}

func (this *SafeAgents) Put(req *model.AgentReportRequest) {
	val := &model.AgentUpdateInfo{
		LastUpdate:    time.Now().Unix(),
		ReportRequest: req,
	}

	if agentInfo, exists := this.Get(req.Hostname); !exists ||
		agentInfo.ReportRequest.AgentVersion != req.AgentVersion ||
		agentInfo.ReportRequest.IP != req.IP ||
		agentInfo.ReportRequest.PluginVersion != req.PluginVersion {

		db.UpdateAgent(val)
		this.Lock()
		this.M[req.Hostname] = val
		this.Unlock()
	} else {
		// update timestamp
		this.Lock()
		this.M[req.Hostname] = val
		this.Unlock()
	}
}

func (this *SafeAgents) Get(hostname string) (*model.AgentUpdateInfo, bool) {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[hostname]
	return val, exists
}

func (this *SafeAgents) Delete(hostname string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, hostname)
}

func (this *SafeAgents) Keys() []string {
	this.RLock()
	defer this.RUnlock()
	count := len(this.M)
	keys := make([]string, count)
	i := 0
	for hostname := range this.M {
		keys[i] = hostname
		i++
	}
	return keys
}

func DeleteStaleAgents() {
	duration := time.Minute * time.Duration(10)
	for {
		time.Sleep(duration)
		deleteStaleAgents()
	}
}

func deleteStaleAgents() {
	//10min
	before := time.Now().Unix() - 60*10
	keys := Agents.Keys()
	count := len(keys)
	if count == 0 {
		return
	}

	for i := 0; i < count; i++ {
		curr, _ := Agents.Get(keys[i])
		if curr.LastUpdate < before {
			Agents.Delete(curr.ReportRequest.Hostname)
		}
	}
}
