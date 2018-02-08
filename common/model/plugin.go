package model

import (
	"fmt"
)

type PluginBasicInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Step    int    `json:"step"` // frequency
}

type PluginStatus int

// plugin status
const (
	StatusOK          = iota // 0 success
	StatusNoFile             // 1 file not exist
	StatusNoPerm             // 2 no exec permission
	StatusTimeout            // 3 timeout
	StatusExeErr             // 4 plugin exec error
	StatusNoResult           // 5 plugin exec no result
	StatusErrorResult        // 6 plugin exec error result
)

type SimplePluginRequest struct {
	IP        string
	Timestamp int64
}

func (this *SimplePluginRequest) String() string {
	return fmt.Sprintf("<SimplePluginRequest: IP:%s, Timestamp:%d>",
		this.IP,
		this.Timestamp,
	)
}

type PluginListResponse struct {
	Plugins []*PluginBasicInfo
	Length  int
}

func (this *PluginListResponse) String() string {
	return fmt.Sprintf("<PluginListResponse: Plugins:%#v, Length:%d>",
		this.Plugins,
		this.Length,
	)
}

type PluginExecInfo struct {
	Name    string       `json:"name"`
	Version string       `json:"version"`
	Status  PluginStatus `json:"status"`
}

type PluginStatusRequest struct {
	IP           string         `json:"ip"`
	AgentVersion string         `json:"agent_version"`
	Plugin       PluginExecInfo `json:"plugin"`
	Timestamp    int64          `json:"timestamp"`
}

type PluginReportStatus struct {
	IP            string       `json:"ip"`
	AgentVersion  string       `json:"agent_version"`
	PluginName    string       `json:"plugin_name"`
	PluginVersion string       `json:"plugin_version"`
	PluginStatus  PluginStatus `json:"plugin_status"`
	Timestamp     int64        `json:"timestamp"`
}

//PK get uniq string for consistent hash ring
func (prs *PluginReportStatus) PK() string {
	return fmt.Sprintf("%s/%s/%s", prs.IP, prs.PluginName, prs.PluginVersion)
}

func (this *PluginStatusRequest) String() string {
	return fmt.Sprintf("<PluginStatusRequest: IP:%s, AgentVersion:%s, Plugin:%#v, Timestamp:%d>",
		this.IP,
		this.AgentVersion,
		this.Plugin,
		this.Timestamp,
	)
}

type ResultCode int

const (
	ResultOK             = iota // 0 OK
	ResultParamLenErr           // 1 parameter length error
	ResultParamFormatErr        // 2 parameter format error
	ResultParamErr              // 3 parameter error
	ResultFail                  // 4 fail
)

type SimpleResponse struct {
	Code ResultCode
	Msg  string
}

func (this *SimpleResponse) String() string {
	return fmt.Sprintf("<SimpleResponse: Code:%d, Msg:%s>",
		this.Code,
		this.Msg,
	)
}

type PluginCmd int

const (
	CmdInvalid = iota //Invalid
	CmdStart          // 1
	CmdStop
	CmdDownload
	CmdDelete // 4
)

type PluginCmdInfo struct {
	PluginName   string    `json:"plugin_name"`
	PluginVerson string    `json:"plugin_version"`
	PluginStep   int       `json:"plugin_step"`
	Cmd          PluginCmd `json:"cmd"`
}

type PluginCmdResponse struct {
	Cmds   []*PluginCmdInfo
	Length int
}

func (this *PluginCmdResponse) String() string {
	return fmt.Sprintf("<PluginCmdResponse: Cmds:%#v, Length:%d>",
		this.Cmds,
		this.Length,
	)
}

type CmdResultCode int

const (
	CmdOK             = iota // 0 OK
	CmdNoFile                // 1 file not exist
	CmdFileVersionErr        // 2 file version error
	CmdFileAlready           // 3 file already exist
	CmdNoPerm                // 4 file no exec perssion
	CmdFailed2Del            // 5 failed to delete plugin
	/* CmdFailed2GetFileInfo        // 6 failed to get file info
	CmdFailed2DlFile             // 7 failed to download file
	CmdFailed2WFile              // 8 failed to write file
	CmdDlErrSizeFile             // 9 download file size is bigger than base info size
	CmdFailed2DlBackup           // 10 faild to backup plugin after download plugin */
	CmdUnknown // 6 unknown cmd
)

type PluginCmdResult struct {
	PluginName    string        `json:"plugin_name"`
	PluginVersion string        `json:"plugin_version"`
	Cmd           PluginCmd     `json:"cmd"`
	Result        CmdResultCode `json:"result"`
}

type PluginCmdResultRequest struct {
	IP      string             `json:"ip"`
	Results []*PluginCmdResult `json:"results"`
	Length  int                `json:"length"`
}

func (this *PluginCmdResultRequest) String() string {
	return fmt.Sprintf("<PluginCmdResultRequest:Results:%#v, Length:%d>",
		this.Results,
		this.Length,
	)
}

func (this *PluginCmdResultRequest) PK() string {
	return fmt.Sprintf("%s", this.IP)
}
