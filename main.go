package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
	Phase      string `json:"phase"`
}

func process(resourceMetric, resource string) {
	metricsData, err := os.Open(fmt.Sprintf("data/%s", resourceMetric))
	if err != nil {
		fmt.Println(err)
	}
	defer metricsData.Close()

	byteValue, _ := ioutil.ReadAll(metricsData)

	var resp Response

	json.Unmarshal(byteValue, &resp)

	parsedLines := []string{}
	for _, result := range resp.Data.Results {
		for _, value := range result.Values {
			if len(value) != 2 {
				continue
			}
			parsed := ""
			if strings.Contains(resource, "bytes") {
				if result.Metric.Container == "" {
					continue
				}
				parsed = fmt.Sprintf("container-name:%s timestamp:%s %s:%s", result.Metric.Container, fmt.Sprintf("%d", int(value[0].(float64))), resource, value[1])
			} else {
				parsed = fmt.Sprintf("pod-name:%s timestamp:%s phase:%s %s:%s", result.Metric.Pod, fmt.Sprintf("%d", int(value[0].(float64))), result.Metric.Phase, resource, value[1])
			}
			parsedLines = append(parsedLines, parsed)
		}
	}

	err = os.MkdirAll("output", 0755)
	if err != nil {
		fmt.Println(err)
	}
	parsedData, err := os.Create(fmt.Sprintf("output/%s.txt", resourceMetric))
	if err != nil {
		fmt.Println(err)
	}
	defer parsedData.Close()

	for _, parsedLine := range parsedLines {
		parsedData.WriteString(parsedLine + "\n")
	}
	parsedData.Sync()
}

func main() {
	process("container_memory_working_set_bytes", "memory-usage-bytes")
	process("kube_pod_container_resource_limits", "memory-limit-bytes")
	process("kube_pod_status_phase", "pod-status")
}
