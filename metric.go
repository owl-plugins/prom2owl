package main

import (
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.]+$`)

type TimeSeriesData struct {
	Metric    string            `json:"metric"`    //sys.cpu.idle
	DataType  string            `json:"data_type"` //COUNTER,GAUGE,DERIVE
	Value     float64           `json:"value"`     //99.00
	Timestamp int64             `json:"timestamp"` //unix timestamp
	Cycle     int               `json:"cycle,omitempty"`
	Tags      map[string]string `json:"tags"` //{"product":"app01", "group":"dev02"}
}

func (m *TimeSeriesData) Validate() bool {
	if !reg.MatchString(m.Metric) || m.Metric == "" {
		return false
	}
	switch strings.ToLower(m.DataType) {
	case "gauge", "counter", "derive":
	default:
		return false
	}
	return true
}

func (tsd *TimeSeriesData) AddTags(tags map[string]string) {
	if tsd.Tags == nil {
		tsd.Tags = tags
		return
	}
	for k, v := range tags {
		tsd.Tags[k] = v
	}
}

//tag1=v1,tag2=v2,tag3=v3
//{"tag1":v1,"tag2":v2,"tag3":v3}
func ParseTags(name string) map[string]string {
	res := make(map[string]string)
	kv := strings.Split(name, ",")
	for _, v := range kv {
		tmp := strings.Split(v, "=")
		if len(tmp) != 2 {
			continue
		}
		res[tmp[0]] = tmp[1]
	}
	return res
}
