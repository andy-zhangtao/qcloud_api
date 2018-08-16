package monitor

import (
	"github.com/sirupsen/logrus"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"net/http"
	"errors"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

func QueryServiceCpu(monitor Monitor) (c QcloudMonitor, err error) {
	return queryMetric(monitor, QueryCPUUSed)
}

func QueryServiceMemory(monitor Monitor) (c QcloudMonitor, err error) {
	return queryMetric(monitor, QueryMemoryUsed)
}

func QueryServiceInNetwork(monitor Monitor) (c QcloudMonitor, err error) {
	return queryMetric(monitor, QueryInNetwork)
}

func QueryServiceOutNetwork(monitor Monitor) (c QcloudMonitor, err error) {
	return queryMetric(monitor, QueryOutNetwork)
}

func QueryServiceInBindwidth(monitor Monitor) (c QcloudMonitor, err error) {
	return queryMetric(monitor, QueryInBindwidth)
}

func QueryServiceOutBindwidth(monitor Monitor) (c QcloudMonitor, err error) {
	return queryMetric(monitor, QueryOutBindwidth)
}

func queryMetric(monitor Monitor, kind int) (c QcloudMonitor, err error) {
	signStr, sign := monitor.generatePubParam(kind)

	logrus.WithFields(logrus.Fields{"url": public.Monitor_API_URL + sign, "Key": monitor.SecretKey, "Body": signStr, "Sing": sign}).Info(ModuleName)

	resp, err := http.Get(public.Monitor_API_URL + sign)
	if err != nil {
		err = errors.New(fmt.Sprintf("Query Cpu Used Error [%s]", err.Error()))
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("Read Body Error [%s]", err.Error()))
		return
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return
	}

	return
}
