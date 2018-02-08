// Copyright 2017 JingDong, Inc.
// Written by Zhangyunyang 2017/12/14

package model

import (
	"fmt"
)

type DeviceItem struct {
	Id        interface{} `json:"device_id"`
	Type      interface{} `json:"device_type"`
	DesType   interface{} `json:"device_des_type"`
	IpMan     interface{} `json:"ip_manage"`
	IdcId     interface{} `json:"idc_id"`
	IdcType   interface{} `json:"idc_type"`
	DeptMent1 interface{} `json:"department1"`
	DeptMent2 interface{} `json:"department2"`
	DeptMent3 interface{} `json:"department3"`
	DeptMent4 interface{} `json:"department4"`
}

func (this *DeviceItem) String() string {
	return fmt.Sprintf("<Id:%v, Type:%v, DesType:%v, IpMan:%v, IdcId: %v, IdcType:%v, DeptMent1:%v, DeptMent2:%v, DeptMent3:%v, DeptMent4:%v>",
		this.Id,
		this.Type,
		this.DesType,
		this.IpMan,
		this.IdcId,
		this.IdcType,
		this.DeptMent1,
		this.DeptMent2,
		this.DeptMent3,
		this.DeptMent4,
	)
}

type DeviceMetric struct {
	Ip            string      `json:"ip"`
	Metric        string      `json:"mib"`
	Value         interface{} `json:"value"`
	MetricDetl    string      `json:"mib_detail"`
	Item          string      `json:"item"`
	CreateTime    int64       `json:"create_time"`
	Unit          string      `json:"unit"`
	DeviceId      interface{} `json:"device_id"`
	DeviceType    interface{} `json:"device_type"`
	DeviceDesType interface{} `json:"device_des_type"`
	IdcId         interface{} `json:"idc_id"`
	IdcType       interface{} `json:"idc_type"`
	DeptMent1     interface{} `json:"department1_id"`
	DeptMent2     interface{} `json:"department2_id"`
	DeptMent3     interface{} `json:"department3_id"`
	DeptMent4     interface{} `json:"department4_id"`
}

func (this *DeviceMetric) String() string {
	return fmt.Sprintf("<Ip:%s, Metric:%s, Value:%v, MetricDetl:%s, Item:%s, CreateTime:%d, Unit:%s, DeviceId:%v, DeviceType:%v, DeviceDesType:%v, IdcId:%v,IdcType:%v,DeptMent1:%v,DeptMent2:%v, DeptMent3:%v, DeptMent4:%v>",
		this.Ip,
		this.Metric,
		this.Value,
		this.MetricDetl,
		this.Item,
		this.CreateTime,
		this.Unit,
		this.DeviceId,
		this.DeviceType,
		this.DeviceDesType,
		this.IdcId,
		this.IdcType,
		this.DeptMent1,
		this.DeptMent2,
		this.DeptMent3,
		this.DeptMent4,
	)
}
