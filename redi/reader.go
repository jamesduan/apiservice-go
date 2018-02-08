package redi

import (
	"apiservice/common/model"
	"apiservice/download"
	"apiservice/g"
	"encoding/json"
	"log"

	"github.com/garyburd/redigo/redis"
	// "happy-hbs/common/model"
	// "happy-hbs/modules/hbs/download"
	// "happy-hbs/modules/hbs/g"
)

// read command according to IP
func PopOneCmd(IP string) ([]*model.PluginCmdInfo, error) {

	var retErr error
	var Cmds []*model.PluginCmdInfo
	queue := g.Config().Redis.PluginCmdQueue + "/" + IP
	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Println("Redis RPOP failed:", err)
				retErr = err
			}

			break
		}

		if reply == "" || reply == "nil" {
			log.Println("reply is empty, queue:", queue)
			break
		}

		var cmd model.PluginCmdInfo
		err = json.Unmarshal([]byte(reply), &cmd)
		if err != nil {
			log.Println("json.Unmarshal failed:", err, ", queue:", queue, ", reply:", reply)
			retErr = err
			break
		}

		Cmds = append(Cmds, &cmd)
	}

	return Cmds, retErr
}

//GetExecuteblePluginList fetch plugin list from cache(redis) return
//PluginBasicInfo type array and error first init Redis connection
//second send HGET to redis server width ip parameter get plugin list
func GetExecuteblePluginList(ip string) ([]*model.PluginBasicInfo, error) {

	var retErr error
	var pluginList []*model.PluginBasicInfo

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	rawData, err := rc.Do("HGET", g.Config().Redis.ExecPluginList, ip)
	g.Debug("execute redis command -> HGET", g.Config().Redis.ExecPluginList, ip)
	if err != nil {
		return nil, err
	}
	if rawData == nil {
		g.Debug("redis command result -> can't find any item.")
		return pluginList, nil
	}
	reply, err := redis.String(rawData, err)
	g.Debug("redis command result  -> ", reply)
	if err != nil {
		retErr = err
	}
	err = json.Unmarshal([]byte(reply), &pluginList)
	if err != nil {
		retErr = err
	}
	return pluginList, retErr
}

// read file info according to file_name+file_version
func GetFileInfo(Name, Version string) (*download.BasicFileInfo, error) {

	hkey := g.Config().Redis.PluginInfoHK
	field := Name + "_" + Version

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	reply, err := redis.String(rc.Do("HGET", hkey, field))
	if err != nil {
		if err != redis.ErrNil {
			log.Println("Redis RPOP failed:", err)
			return nil, err
		}
	}

	if reply == "" || reply == "nil" {
		log.Printf("reply is empty, hkey:%s, field:%s\n", hkey, field)
		return nil, nil
	}

	var basicInfo download.BasicFileInfo
	err = json.Unmarshal([]byte(reply), &basicInfo)
	if err != nil {
		log.Printf("json.Unmarshal failed:%v, hkey:%s, field:%s, reply:%v\n", err, hkey, field, reply)
		return nil, err
	}

	return &basicInfo, nil
}
