package cvm

import (
	"strconv"
	"github.com/andy-zhangtao/qcloud_api/public"
	gz "github.com/andy-zhangtao/gogather/zsort"
	"crypto/sha1"
	"crypto/hmac"
	"encoding/base64"
	"log"

	"net/http"
	"io/ioutil"
	"encoding/json"
)

const API_URL = "https://ccs.api.qcloud.com/v2/index.php?"

type Cluster struct {
	Pub        public.Public `json:"pub"`
	Cid        string        `json:"cid"`
	Cname      string        `json:"cname"`
	Status     string        `json:"status"`
	OrderField string        `json:"order_field"`
	OrderType  string        `json:"order_type"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
	SecretKey  string        `json:"secret_key"`
	Namespace  string        `json:"namespace"`
	sign       string
}

type ClusterNode_data_nodes struct {
	InstanceId           string `json:"instanceid"`
	InstanceName         string `json:"instancename"`
	InstanceType         string `json:"instancetype"`
	ZoneId               int    `json:"zoneid"`
	WanIp                string `json:"wanip"`
	LanIp                string `json:"lanip"`
	Cpu                  int    `json:"cpu"`
	Mem                  int    `json:"mem"`
	KernelVersion        string `json:"kernelversion"`
	OsImage              string `json:"osimage"`
	PodCidr              string `json:"podcidr"`
	IsNormal             int    `json:"isnormal"`
	AbnormalReason       string `json:"abnormalreason"`
	CvmState             int    `json:"cvmstate"`
	CvmPayMode           int    `json:"cvmpaymode"`
	NetworkPayMode       int    `json:"networkpaymode"`
	CreatedAt            string `json:"createdat"`
	InstanceCreateTime   string `json:"instancecreatetime"`
	InstanceDeadlineTime string `json:"instancedeadlinetime"`
	Unschedulable        bool   `json:"unschedulable"`
	Zone                 string `json:"zone"`
}
type ClusterNode_data struct {
	TotalCount int                      `json:"totalcount"`
	Nodes      []ClusterNode_data_nodes `json:"nodes"`
}
type ClusterNode struct {
	Code     int              `json:"code"`
	Message  string           `json:"message"`
	CodeDesc string           `json:"codedesc"`
	Data     ClusterNode_data `json:"data"`
}

// queryCluster 查询集群数据API
func (this Cluster) queryCluster() ([]string, map[string]string) {
	var field []string
	req := make(map[string]string)

	if this.Cid != "" {
		field = append(field, "clusterIds.n")
		req["clusterIds.n"] = this.Cid
	}

	if this.Cname != "" {
		field = append(field, "clusterName")
		req["clusterName"] = this.Cname
	}

	if this.Status != "" {
		field = append(field, "status")
		req["status"] = this.Status
	}

	if this.OrderField != "" {
		field = append(field, "orderField")
		req["orderField"] = this.OrderField
	}

	if this.OrderType != "" {
		field = append(field, "orderType")
		req["orderType"] = this.OrderType
	}

	if this.Offset > 0 {
		field = append(field, "offset")
		req["offset"] = strconv.Itoa(this.Offset)
	}

	if this.Limit > 0 {
		field = append(field, "limit")
		req["limit"] = strconv.Itoa(this.Limit)
	}

	return field, req
}

func (this Cluster) queryClusterNode() ([]string, map[string]string) {
	var field []string
	req := make(map[string]string)

	if this.Cid != "" {
		field = append(field, "clusterId")
		req["clusterId"] = this.Cid
	}

	if this.Offset > 0 {
		field = append(field, "offset")
		req["offset"] = strconv.Itoa(this.Offset)
	}

	if this.Limit > 0 {
		field = append(field, "limit")
		req["limit"] = strconv.Itoa(this.Limit)
	}

	if this.Namespace != "" {
		field = append(field, "namespace")
		req["namespace"] = this.Namespace
	}
	return field, req
}

// QueryClusters 查询集群信息
func (this Cluster) QueryClusters() string {
	field, reqmap := this.queryCluster()
	pubMap := public.PublicParam(this.Pub.Action, this.Pub.Region, this.Pub.SecretId)
	this.sign = generateSignatureString(field, reqmap, pubMap)
	signStr := "GETccs.api.qcloud.com/v2/index.php?" + this.sign
	sign := generateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + sign

	log.Println(reqURL)
	return reqURL
}

func (this Cluster) QueryClusterNodes() (*ClusterNode,error) {
	field, reqmap := this.queryClusterNode()
	pubMap := public.PublicParam(this.Pub.Action, this.Pub.Region, this.Pub.SecretId)
	this.sign = generateSignatureString(field, reqmap, pubMap)
	signStr := "GETccs.api.qcloud.com/v2/index.php?" + this.sign
	sign := generateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + sign

	resp, err := http.Get(API_URL + reqURL)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cn ClusterNode

	err = json.Unmarshal(data, &cn)
	if err != nil {
		return nil, err
	}

	return &cn, nil
}

// generateSignature 生成请求签名字符串
// field 请求字段集合
// reqmap 待计算的请求map
// publicmap 公共请求map,调用public.PublicParam生成
func generateSignatureString(field []string, reqmap, publicMap map[string]string) string {
	field = append(field, public.PubilcField...)
	field = gz.DictSort(field)
	//log.Println(field)

	req := ""
	for k, v := range reqmap {
		publicMap[k] = v
	}
	for i, key := range field {
		if i == 0 {
			req = key + "=" + publicMap[key]
		} else {
			req += "&" + key + "=" + publicMap[key]
		}
	}
	return req
}

// generateSignature 生成最终的请求签名,使用HMAC-SHA1
// key加密key，req请求字符串
func generateSignature(key, req string) string {
	k := []byte(key)
	mac := hmac.New(sha1.New, k)
	mac.Write([]byte(req))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
