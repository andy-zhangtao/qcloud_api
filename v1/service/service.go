package service

import (
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"strconv"
	"log"
)

const API_URL = "https://ccs.api.qcloud.com/v2/index.php?"

type Svc struct {
	Pub          public.Public `json:"pub"`
	ClusterId    string        `json:"cluster_id"`
	Namespace    string        `json:"namespace"`
	Allnamespace string        `json:"allnamespace"`
	SecretKey    string        `json:"secret_key"`
	sign         string
}

type SvcData_data_services struct {
	ServiceName     string `json:"servicename"`
	Status          string `json:"status"`
	ServiceIp       string `json:"serviceip"`
	ExternalIp      string `json:"externalip"`
	LbId            string `json:"lbid"`
	LbStatus        string `json:"lbstatus"`
	AccessType      string `json:"accesstype"`
	DesiredReplicas int    `json:"desiredreplicas"`
	CurrentReplicas int    `json:"currentreplicas"`
	CreatedAt       string `json:"createdat"`
	Namespace       string `json:"namespace"`
}
type SvcData_data struct {
	TotalCount int                     `json:"totalcount"`
	Services   []SvcData_data_services `json:"services"`
}
type SvcSMData struct {
	Code     int          `json:"code"`
	Message  string       `json:"message"`
	CodeDesc string       `json:"codedesc"`
	Data     SvcData_data `json:"data"`
}

func (this Svc) querySampleInfo() ([]string, map[string]string) {
	var field []string
	req := make(map[string]string)

	if this.ClusterId != "" {
		field = append(field, "clusterId")
		req["clusterId"] = this.ClusterId
	}

	if this.Namespace != "" {
		field = append(field, "namespace")
		req["namespace"] = this.Namespace
	}
	an, _ := strconv.Atoi(this.Allnamespace)
	if an != 0 {
		field = append(field, "allnamespace")
		req["allnamespace"] = this.Allnamespace
	}

	return field, req
}

// QueryClusters 查询集群信息
func (this Svc) QuerySampleInfo() string {
	field, reqmap := this.querySampleInfo()
	pubMap := public.PublicParam(this.Pub.Action, this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GETccs.api.qcloud.com/v2/index.php?" + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + sign

	log.Println(reqURL)
	return reqURL
}
