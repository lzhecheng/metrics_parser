package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Response struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Results    []Result `json:"result"`
}

type Result struct {
	Metric Metric          `json:"metric"`
	Values [][]interface{} `json:"values"`
}

type Metric struct {
	MetricName string `json:"__name__"`
	Cluster    string `json:"cluster"`
	Container  string `json:"container"`
	ID         string `json:"id"`
	Image      string `json:"image"`
	Instance   string `json:"instance"`
	Job        string `json:"job"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Pod        string `json:"pod"`
}

func main() {
	jsonFile, err := os.Open("data/container_memory_working_set_bytes")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var resp Response

	json.Unmarshal(byteValue, &resp)

	metricNames := []string{}
	for _, result := range resp.Data.Results {
		metricNames = append(metricNames, result.Metric.Pod)
	}
	fmt.Println(metricNames)
}
