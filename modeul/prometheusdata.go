package modeul

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

// 定义常量 后续变成读取环境变量

type GroupData struct {
	Metrics   map[string]*prometheus.Desc
	mutex     sync.Mutex
	XxlUrl    string
	xxlCookie string
}

type NewGroupData struct {
	RecordsTotal int
	ALLOleTotal  float64
	ALlNewTotal  float64
	Data         []groupdate
}

type groupdate struct {
	Id             int
	AppName        string
	Title          string
	ALLOleTotal    float64
	ALlNewTotal    float64
	NewErrorTotal  float64
	NewAccessTotal float64
	OldErrorTotal  float64
	OldAccessTotal float64
}

func NewGroupMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}

func NewMetrics(namespace string, XXLURL, XXLCOOKIE string) *GroupData {
	return &GroupData{
		Metrics: map[string]*prometheus.Desc{
			"XXL-Job-UP":          NewGroupMetric(namespace, "xxl_job_up", "xxl status", nil),
			"GroupOldErrorTotal":  NewGroupMetric(namespace, "group_old_error_total", "the old error total", []string{"name"}),
			"GroupOldAccessTotal": NewGroupMetric(namespace, "group_old_access_total", "the old access total", []string{"name"}),
			"GroupNewErrorTotal":  NewGroupMetric(namespace, "group_new_error_total", "the new error total", []string{"name"}),
			"GroupNewAccessTotal": NewGroupMetric(namespace, "group_new_access_total", "the new access total", []string{"name"}),
			"GroupNewTotal":       NewGroupMetric(namespace, "group_new_total", "the new access total", []string{"name"}),
			"GroupOldTotal":       NewGroupMetric(namespace, "group_old_total", "the new access total", []string{"name"}),
		},
		XxlUrl:    XXLURL,
		xxlCookie: XXLCOOKIE,
	}
}

var GroupDataNew *NewGroupData

func (g *GroupData) Collect(ch chan<- prometheus.Metric) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	a := g.GetGroupId()
	for _, value := range a.Data {
		ch <- prometheus.MustNewConstMetric(g.Metrics["GroupOldErrorTotal"], prometheus.CounterValue, value.OldErrorTotal, value.Title)
	}
	for _, value := range a.Data {
		ch <- prometheus.MustNewConstMetric(g.Metrics["GroupNewErrorTotal"], prometheus.CounterValue, value.NewErrorTotal, value.Title)
	}
	for _, value := range a.Data {
		ch <- prometheus.MustNewConstMetric(g.Metrics["GroupNewAccessTotal"], prometheus.CounterValue, value.NewAccessTotal, value.Title)
	}
	for _, value := range a.Data {
		ch <- prometheus.MustNewConstMetric(g.Metrics["GroupOldAccessTotal"], prometheus.CounterValue, value.OldAccessTotal, value.Title)
	}
	for _, value := range a.Data {
		ch <- prometheus.MustNewConstMetric(g.Metrics["GroupNewTotal"], prometheus.CounterValue, value.ALlNewTotal, value.Title)
	}
	for _, value := range a.Data {
		ch <- prometheus.MustNewConstMetric(g.Metrics["GroupOldTotal"], prometheus.CounterValue, value.ALLOleTotal, value.Title)
	}
}

func (g *GroupData) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range g.Metrics {
		ch <- m
	}
}

// 获取组的最新/昨日的成功/失败数据

func (g *GroupData) GetGroupId() NewGroupData {
	startTime := time.Now()
	fmt.Printf("url: %v\n", g.XxlUrl)
	geturl := fmt.Sprintf("curl '%s/xxl-job-admin/jobgroup/pageList' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,be;q=0.6,sq;q=0.5' -H 'Cache-Control: no-cache' -H 'Connection: keep-alive' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Cookie: %v' --data 'appname=&title=&start=0&length=100' --compressed --insecure", g.XxlUrl, g.xxlCookie)
	var groupdata AllGroupData
	curl := exec.Command("bash", "-c", geturl)
	out, err := curl.Output()
	if err != nil {
		fmt.Printf("GroupIderr: %v", err)
	}
	json.Unmarshal(out, &groupdata)
	newdate := NewGroupData{RecordsTotal: groupdata.RecordsTotal}
	wg := sync.WaitGroup{}
	ch := make(chan groupdate)
	for _, value := range groupdata.Data {
		go func(value allGroupData) {
			group := groupdate{}
			group.Id = value.Id
			group.AppName = value.Appname
			group.Title = value.Title
			group.ALLOleTotal, group.ALlNewTotal = GetGroupTotal(value.Id, g.XxlUrl, g.xxlCookie)
			group.OldAccessTotal, group.NewAccessTotal = GetAccessGroupData(value.Id, g.XxlUrl, g.xxlCookie)
			group.OldErrorTotal, group.NewErrorTotal = GetErrorGroupData(value.Id, g.XxlUrl, g.xxlCookie)
			ch <- group
		}(value)
	}
	wg.Add(len(groupdata.Data))
	go func() {
		for i := range ch {
			newdate.Data = append(newdate.Data, i)
			wg.Done()
		}
	}()
	wg.Wait()
	close(ch)
	endTime := time.Now()
	fmt.Printf("Time: %v", endTime.Sub(startTime))
	return newdate
}

func GetAccessGroupData(groupid int, XXLURL, XXLCOOKIE string) (float64, float64) {
	newtime := fmt.Sprintf(time.Now().Format("2006-01-02"))
	oldtime := fmt.Sprintf(time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	oldurl := fmt.Sprintf("curl '%s/xxl-job-admin/joblog/pageList' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,be;q=0.6,sq;q=0.5' -H 'Cache-Control: no-cache' -H 'Connection: keep-alive' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Cookie: %v' --data 'jobGroup=%v&jobId=0&logStatus=1&filterTime=%v 00:00:00 - %v 23:59:59&start=0&length=10' --compressed --insecure", XXLURL, XXLCOOKIE, groupid, oldtime, oldtime)
	newurl := fmt.Sprintf("curl '%s/xxl-job-admin/joblog/pageList' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,be;q=0.6,sq;q=0.5' -H 'Cache-Control: no-cache' -H 'Connection: keep-alive' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Cookie: %v' --data 'jobGroup=%v&jobId=0&logStatus=1&filterTime=%v 00:00:00 - %v 23:59:59&start=0&length=10' --compressed --insecure", XXLURL, XXLCOOKIE, groupid, newtime, newtime)
	newdata := UniversalData{}
	olddata := UniversalData{}
	newcurl := exec.Command("bash", "-c", newurl)
	oldcurl := exec.Command("bash", "-c", oldurl)
	oldout, err := oldcurl.Output()
	if err != nil {
		fmt.Printf("GetAccessGroupDataError: %v", err)
	}
	newout, err := newcurl.Output()
	if err != nil {
		fmt.Printf("GetAccessGroupDataError: %v", err)
	}
	json.Unmarshal(oldout, &olddata)
	json.Unmarshal(newout, &newdata)
	return float64(olddata.RecordsTotal), float64(newdata.RecordsTotal)
}

func GetErrorGroupData(groupid int, XXLURL, XXLCOOKIE string) (float64, float64) {
	newtime := fmt.Sprintf(time.Now().Format("2006-01-02"))
	oldtime := fmt.Sprintf(time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	oldurl := fmt.Sprintf("curl '%s/xxl-job-admin/joblog/pageList' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,be;q=0.6,sq;q=0.5' -H 'Cache-Control: no-cache' -H 'Connection: keep-alive' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Cookie: %v' --data 'jobGroup=%v&jobId=0&logStatus=2&filterTime=%v 00:00:00 - %v 23:59:59&start=0&length=10' --compressed --insecure", XXLURL, XXLCOOKIE, groupid, oldtime, oldtime)
	newurl := fmt.Sprintf("curl '%s/xxl-job-admin/joblog/pageList' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,be;q=0.6,sq;q=0.5' -H 'Cache-Control: no-cache' -H 'Connection: keep-alive' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Cookie: %v' --data 'jobGroup=%v&jobId=0&logStatus=2&filterTime=%v 00:00:00 - %v 23:59:59&start=0&length=10' --compressed --insecure", XXLURL, XXLCOOKIE, groupid, newtime, newtime)
	newdata := UniversalData{}
	olddata := UniversalData{}
	newcurl := exec.Command("bash", "-c", newurl)
	oldcurl := exec.Command("bash", "-c", oldurl)
	oldout, err := oldcurl.Output()
	if err != nil {
		fmt.Printf("GetOldAccessGroupDataError: %v", err)
	}
	newout, err := newcurl.Output()
	if err != nil {
		fmt.Printf("GetNewAccessGroupDataError: %v", err)
	}
	json.Unmarshal(oldout, &olddata)
	json.Unmarshal(newout, &newdata)
	return float64(olddata.RecordsTotal), float64(newdata.RecordsTotal)
}

func GetGroupTotal(groupid int, XXLURL, XXLCOOKIE string) (float64, float64) {
	newtime := fmt.Sprintf(time.Now().Format("2006-01-02"))
	oldtime := fmt.Sprintf(time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	oldurl := fmt.Sprintf("curl '%s/xxl-job-admin/joblog/pageList' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,be;q=0.6,sq;q=0.5' -H 'Cache-Control: no-cache' -H 'Connection: keep-alive' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Cookie: %v' --data 'jobGroup=%v&jobId=0&logStatus=0&filterTime=%v 00:00:00 - %v 23:59:59&start=0&length=10' --compressed --insecure", XXLURL, XXLCOOKIE, groupid, oldtime, oldtime)
	newurl := fmt.Sprintf("curl '%s/xxl-job-admin/joblog/pageList' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,be;q=0.6,sq;q=0.5' -H 'Cache-Control: no-cache' -H 'Connection: keep-alive' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Cookie: %v' --data 'jobGroup=%v&jobId=0&logStatus=0&filterTime=%v 00:00:00 - %v 23:59:59&start=0&length=10' --compressed --insecure", XXLURL, XXLCOOKIE, groupid, newtime, newtime)
	newdata := UniversalData{}
	olddata := UniversalData{}
	newcurl := exec.Command("bash", "-c", newurl)
	oldcurl := exec.Command("bash", "-c", oldurl)
	oldout, err := oldcurl.Output()
	if err != nil {
		fmt.Printf("GetAccessGroupDataError: %v", err)
	}
	newout, err := newcurl.Output()
	if err != nil {
		fmt.Printf("GetAccessGroupDataError: %v", err)
	}
	json.Unmarshal(oldout, &olddata)
	json.Unmarshal(newout, &newdata)
	return float64(olddata.RecordsTotal), float64(newdata.RecordsTotal)
}

func CheckHeartbeat(XXLURL string) float64 {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	req, err := client.Get(XXLURL)
	if err != nil {
		return float64(0)
	}
	if req.StatusCode != 200 {
		return float64(1)
	} else {
		return float64(0)
	}
}
