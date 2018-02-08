package rpc

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"time"

	"apiservice/cache"
	"apiservice/common/model"
	"apiservice/common/utils"
	"apiservice/download"
	"apiservice/g"
	"apiservice/redi"
	"apiservice/sender"
)

func (t *Agent) MinePlugins(args model.AgentHeartbeatRequest, reply *model.AgentPluginsResponse) error {
	if args.Hostname == "" {
		return nil
	}

	reply.Plugins = cache.GetPlugins(args.Hostname)
	reply.Timestamp = time.Now().Unix()

	return nil
}

func (t *Agent) ReportStatus(req *model.AgentInfoRequest, resp *model.SimpleResponse) error {

	log.Println("Agent.ReportStatus, req:", req)
	if req.IP == "" {
		log.Printf("Invalid IP:%s\n", req.IP)
		resp.Code = model.ResultParamErr
		resp.Msg = "parameter error"
		return fmt.Errorf("Invalid IP")
	}

	item := &model.KafkaItem{Type: model.AgentInfoType, Body: req}
	sender.Push2KafkaSendQueue(item)

	resp.Code = model.ResultOK
	resp.Msg = "success"

	log.Println("Agent.ReportStatus, resp:", resp)

	return nil
}

// 需要checksum一下来减少网络开销？其实白名单通常只会有一个或者没有，无需checksum
func (t *Agent) TrustableIps(args *model.NullRpcRequest, ips *string) error {
	*ips = strings.Join(g.Config().Trustable, ",")
	return nil
}

// agent按照server端的配置，按需采集的metric，比如net.port.listen port=22 或者 proc.num name=zabbix_agentd
func (t *Agent) BuiltinMetrics(args *model.AgentHeartbeatRequest, reply *model.BuiltinMetricResponse) error {
	if args.Hostname == "" {
		return nil
	}

	metrics, err := cache.GetBuiltinMetrics(args.Hostname)
	if err != nil {
		return nil
	}

	checksum := ""
	if len(metrics) > 0 {
		checksum = DigestBuiltinMetrics(metrics)
	}

	if args.Checksum == checksum {
		reply.Metrics = []*model.BuiltinMetric{}
	} else {
		reply.Metrics = metrics
	}
	reply.Checksum = checksum
	reply.Timestamp = time.Now().Unix()

	return nil
}

func DigestBuiltinMetrics(items []*model.BuiltinMetric) string {
	sort.Sort(model.BuiltinMetricSlice(items))

	var buf bytes.Buffer
	for _, m := range items {
		buf.WriteString(m.String())
	}

	return utils.Md5(buf.String())
}

// Download
// Get basic information
func (this *Agent) GetBasicInfo(req *model.BasicInfoRequest, resp *model.BasicInfoResponse) error {
	log.Println("GetBasicInfo req:", req)

	// get file information
	fileName := strings.TrimSpace(req.FileName)
	fileVersion := strings.TrimSpace(req.FileVersion)

	f, err := redi.GetFileInfo(fileName, fileVersion)
	if err != nil {
		log.Println("redi.GetFileInfo failed:", err)
		return err
	}

	//open file
	err = download.GFileSet.Add(f)
	if err != nil {
		log.Printf("Add download file failed, error:%v, file:%v\n", err, f)
		return err
	}

	//write download status
	rt := &model.PluginCmdResultRequest{IP: req.IP, Length: 1}
	rs := &model.PluginCmdResult{PluginName: req.FileName, PluginVersion: req.FileVersion, Cmd: model.CmdDownload, Result: model.DLStart}

	rt.Results = append(rt.Results, rs)

	_, err = redi.WriteCmdResult(rt)
	if err != nil {
		log.Println("write download status failed:", err)
		return err
	}

	//basic info
	resp.Size = f.Size
	resp.CheckSum = f.CheckSum
	resp.Step = f.Step
	resp.Version = f.Version

	log.Println("GetBasicInfo resp:", resp)
	return nil
}

// Download
func (this *Agent) Download(req *model.ConcreteInfoRequest, resp *model.ConcreteInfoResponse) error {

	log.Println("Download req:", req)

	fileKey := strings.TrimSpace(req.FileName) + "_" + strings.TrimSpace(req.FileVersion)

	f, err := download.GFileSet.Get(fileKey)
	if err != nil {
		log.Printf("get file failed, req:%v, err:%v\n", req, err)
		return err
	}

	//fd
	fd, err := f.GetFD()
	if err != nil {
		log.Println("Get fd failed:", err)
		return err
	}

	buf := make([]byte, req.Length)
	n, err := fd.ReadAt(buf, req.Offset)
	if err != nil && err != io.EOF {
		log.Println("ReadAt failed:", err)
		return err
	}

	resp.Length = int64(n)
	resp.Content = buf

	log.Println("Download resp:", resp)
	return nil
}
