// Copyright 2017 JingDong, Inc.
// Written by Zhangyunyang 2017/09/28

package model

const (
	InvalidType     = iota
	PluginSatusType // 1
	AgentInfoType
	CmdResultType
)

type ItemType int

type KafkaItem struct {
	Type ItemType    `json:"type"`
	Body interface{} `json:"body"`
}
