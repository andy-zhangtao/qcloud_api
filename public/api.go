package public

import (
	gg "github.com/andy-zhangtao/gogather/time"
	ga "github.com/andy-zhangtao/gogather/random"
)

type Public struct {
	Action   string `json:"action"`
	SecretId string `json:"secret_id"`
	Region   string `json:"region"`
}

var (
	PubilcField = []string{"Action", "SecretId", "Region", "Timestamp", "Nonce"}
)

// PublicParam生成公共请求数据
// 包括 Action/SecretId/Region/Timestamp/Nonce
func PublicParam(action, region, secretId string) map[string]string {
	req := make(map[string]string)
	req["Action"] = action
	req["SecretId"] = secretId
	req["Region"] = region
	req["Timestamp"] = gg.GetTimeStamp(10)
	req["Nonce"] = ga.GetRandom(6)
	return req
}

func Generate() {

}
