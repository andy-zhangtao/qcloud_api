package monitor

import (
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/const/v1"
	"net/url"
)

type Monitor struct {
	Pub         public.Public `json:"pub"`
	SecretKey   string
	sign        string
	ClusterID   string
	ServiceName string
	Namespace   string
	StartTime   string
	EndTime     string
	Params      map[string]interface{}
}

type QcloudMonitor struct {
	Code       int       `json:"code"`
	Message    string    `json:"message"`
	CodeDesc   string    `json:"codedesc"`
	MetricName string    `json:"metricname"`
	StartTime  string    `json:"starttime"`
	EndTime    string    `json:"endtime"`
	Period     int       `json:"period"`
	DataPoints []float64 `json:"datapoints"`
}

const (
	ModuleName   = "Qcloud-Monitor-Agent"
	QueryCPUUSed = iota
	QueryMemoryUsed
	QueryInNetwork
	QueryOutNetwork
	QueryInBindwidth
	QueryOutBindwidth
)

func (this *Monitor) generatePubParam(kind int) (string, string) {
	var field []string
	reqmap := make(map[string]string)

	reqmap["namespace"] = "qce/docker"
	reqmap["dimensions.0.name"] = "clusterId"
	reqmap["dimensions.0.value"] = this.ClusterID
	reqmap["dimensions.1.name"] = "serviceName"
	reqmap["dimensions.1.value"] = this.ServiceName
	reqmap["dimensions.2.name"] = "namespace"
	reqmap["dimensions.2.value"] = this.Namespace
	reqmap["startTime"] = this.StartTime
	reqmap["endTime"] = this.EndTime
	reqmap["period"] = "60"
	optKind := "GetMonitorData"
	field = append(field, []string{"namespace", "metricName", "dimensions.0.name", "dimensions.0.value", "dimensions.1.name", "dimensions.1.value", "dimensions.2.name", "dimensions.2.value", "startTime", "endTime", "period"}...)

	switch kind {
	case QueryCPUUSed:
		reqmap["metricName"] = "service_cpu_used"
	case QueryMemoryUsed:
		reqmap["metricName"] = "service_mem_used"
	case QueryInNetwork:
		reqmap["metricName"] = "service_in_flux"
	case QueryOutNetwork:
		reqmap["metricName"] = "service_out_flux"
	case QueryInBindwidth:
		reqmap["metricName"] = "service_in_bandwidth"
	case QueryOutBindwidth:
		reqmap["metricName"] = "service_out_bandwidth"
	}

	pubMap := public.PublicParam(optKind, this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GET" + v1.QCloudMonitorEndpoint + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	return signStr, this.sign + "&Signature=" + url.QueryEscape(sign)
}
