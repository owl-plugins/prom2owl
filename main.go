package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/prometheus/common/log"
	"github.com/prometheus/prom2json"

	dto "github.com/prometheus/client_model/go"
)

func main() {
	cert := flag.String("cert", "", "certificate file")
	key := flag.String("key", "", "key file")
	labels := flag.String("labels", "", "custom label, eg:tagk1=tagv1,tagk2=tagv2")
	prefixs := flag.String("prefixs", "", "metric prefixs")
	skipServerCertCheck := flag.Bool("accept-invalid-cert", false, "Accept any certificate during TLS handshake. Insecure, use only for testing.")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatalf("Usage: %s METRICS_URL", os.Args[0])
	}
	if (*cert != "" && *key == "") || (*cert == "" && *key != "") {
		log.Fatalf("Usage: %s METRICS_URL\n with TLS client authentication: %s -cert=/path/to/certificate -key=/path/to/key METRICS_URL", os.Args[0], os.Args[0])
	}

	mfChan := make(chan *dto.MetricFamily, 1024)

	go prom2json.FetchMetricFamilies(flag.Args()[0], mfChan, *cert, *key, *skipServerCertCheck)
	metrics := []TimeSeriesData{}
	ts := time.Now().Unix()
	customTags := ParseTags(*labels)
	for mf := range mfChan {
		switch mf.GetType() {
		case dto.MetricType_SUMMARY, dto.MetricType_HISTOGRAM:
			continue
		default:
			for _, m := range mf.Metric {
				hasPrefix := false
				for _, prefix := range strings.Split(*prefixs, ",") {
					if strings.HasPrefix(mf.GetName(), prefix) {
						hasPrefix = true
						break
					}
				}
				if !hasPrefix {
					continue
				}
				metric := TimeSeriesData{
					Metric:    strings.Replace(mf.GetName(), "_", ".", -1),
					DataType:  mf.GetType().String(),
					Timestamp: ts,
					Tags:      makeLabels(m),
					Value:     getValue(m),
				}
				metric.AddTags(customTags)
				metrics = append(metrics, metric)
			}
		}
	}
	jsonBytes, _ := json.Marshal(metrics)
	if _, err := os.Stdout.Write(jsonBytes); err != nil {
		log.Fatalln("error writing to stdout:", err)
	}
	fmt.Println()
}

func getValue(m *dto.Metric) float64 {
	if m.Gauge != nil {
		return m.GetGauge().GetValue()
	}
	if m.Counter != nil {
		return m.GetCounter().GetValue()
	}
	if m.Untyped != nil {
		return m.GetUntyped().GetValue()
	}
	return 0.
}

func makeLabels(m *dto.Metric) map[string]string {
	result := map[string]string{}
	for _, lp := range m.Label {
		result[lp.GetName()] = lp.GetValue()
	}
	return result
}
