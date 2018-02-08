package redi

import (
	"apiservice/common/model"
	"apiservice/g"
	"encoding/json"
	"log"
	// "happy-hbs/common/model"
	// "happy-hbs/modules/hbs/g"
)

func lpush(queue, message string) error {
	rc := g.RedisConnPool.Get()
	defer rc.Close()
	_, err := rc.Do("LPUSH", queue, message)
	if err != nil {
		log.Println("LPUSH redis", queue, "fail:", err, "message:", message)
		return err
	}

	return nil
}

func WriteCmdResult(r *model.PluginCmdResultRequest) (model.ResultCode, error) {

	var (
		rCode model.ResultCode = model.ResultOK
	)

	queue := g.Config().Redis.PluginCmdResultQueue

	bs, err := json.Marshal(r)
	if err != nil {
		log.Println("json.Marshal failed:", err, ", r:", r)
		rCode = model.ResultParamFormatErr

		return rCode, err
	}

	log.Printf("write result to queue, result:%v, queue:%s\n", r, queue)
	err = lpush(queue, string(bs))
	if err != nil {
		log.Println("redis push failed:", err, ",queue:", queue)
		rCode = model.ResultFail

		return rCode, err
	}

	return rCode, nil
}

func WriteAgentInfo(info *model.AgentInfoRequest) (model.ResultCode, error) {

	var (
		rCode model.ResultCode = model.ResultOK
	)

	queue := g.Config().Redis.AgentStatusQueue

	bs, err := json.Marshal(info)
	if err != nil {
		log.Println("json.Marshal failed:", err, ", agent info:", info)
		rCode = model.ResultParamFormatErr
		return rCode, err
	}

	log.Printf("write agent info to queue, agent info:%v, queue:%s\n", info, queue)
	err = lpush(queue, string(bs))
	if err != nil {
		log.Println("redis push failed:", err, ",queue:", queue)
		rCode = model.ResultFail
		return rCode, err
	}

	return rCode, nil
}
