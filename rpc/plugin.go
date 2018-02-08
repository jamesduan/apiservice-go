package rpc

import (
	"apiservice/common/model"
	"apiservice/g"
	"apiservice/redi"
	"apiservice/sender"
	"fmt"

	// "happy-hbs/common/model"
	// "happy-hbs/modules/hbs/g"
	// "happy-hbs/modules/hbs/redi"
	// "happy-hbs/modules/hbs/sender"
	"log"
)

//GetExecPluginList plugin to exec
func (t *Plugin) GetExecPluginList(req *model.SimplePluginRequest, resp *model.PluginListResponse) error {
	g.Debug("request -> ", req)
	// get agent ip
	ipaddress := req.IP
	if ipaddress == "" {
		errorText := "ip address is empty, please check the request of agent."
		g.Debug(errorText)
		return GetExecPluginListError(errorText)
	}
	// fetch plugin list from redis cache
	plugins, err := redi.GetExecuteblePluginList(ipaddress)
	// plugins, err := redi.PopOneCmd(ipaddress)
	if err != nil {
		g.Debug("read plugins from redis error,", err)
		return err
	}
	if g.Config().Debug {
		log.Println("=================get executeble plugin list(" + ipaddress + ")=======================")
		for _, plugin := range plugins {
			log.Println("plugin: ", *plugin)
		}
		log.Println("==================================================================")
	}
	resp.Plugins = plugins
	resp.Length = len(plugins)
	return nil
}

//ReportStatus  status
func (t *Plugin) ReportStatus(req *model.PluginStatusRequest, resp *model.SimpleResponse) error {

	pluginReportStatus := &model.PluginReportStatus{
		IP:            req.IP,
		AgentVersion:  req.AgentVersion,
		PluginName:    req.Plugin.Name,
		PluginVersion: req.Plugin.Version,
		PluginStatus:  req.Plugin.Status,
		Timestamp:     req.Timestamp,
	}

	item := &model.KafkaItem{Type: model.PluginSatusType, Body: pluginReportStatus}
	sender.Push2KafkaSendQueue(item)

	resp.Code = 0
	resp.Msg = "success"
	log.Println("Plugin.ReportStatus resp:", resp)
	return nil
}

//SyncCmd sync command
func (t *Plugin) SyncCmd(req *model.SimplePluginRequest, resp *model.PluginCmdResponse) error {

	log.Println("SyncCmd ... req:", req)

	if req.IP == "" {
		log.Println("SyncCmd invalid ip:", req.IP)
		return fmt.Errorf("Invalid IP")
	}

	cmds, err := redi.PopOneCmd(req.IP)
	if err != nil {
		log.Println("redis PopOneCmd failed:", err, ", IP:", req.IP)
		return err
	}

	resp.Cmds = cmds
	resp.Length = len(cmds)

	log.Println("SyncCmd resp:", resp)
	return nil
}

// report command result
func (t *Plugin) ReportCmdResult(req *model.PluginCmdResultRequest, resp *model.SimpleResponse) error {

	log.Println("ReportCmdResult ... req:", req)

	if len(req.Results) != req.Length || len(req.Results) == 0 {
		log.Printf("Incorrect length, len:%d, length:%d\n", len(req.Results), req.Length)
		resp.Code = model.ResultParamLenErr
		resp.Msg = "parameter error"
		return fmt.Errorf("incorrect length")
	}

	item := &model.KafkaItem{Type: model.CmdResultType, Body: req}
	sender.Push2KafkaSendQueue(item)

	resp.Code = model.ResultOK
	resp.Msg = "success"

	log.Println("ReportCmdResult resp:", resp)
	return nil
}
