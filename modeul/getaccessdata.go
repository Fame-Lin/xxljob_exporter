package modeul

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

// 获取今日的执行成功的调度日志数据

func GetNewAccessData(XXLURL, XXLCOOKIE string) float64 {
	newtime := fmt.Sprintf(time.Now().Format("2006-01-02"))
	geturl := fmt.Sprintf("curl '%s/xxl-job-admin/joblog/pageList' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,be;q=0.6,sq;q=0.5' -H 'Cache-Control: no-cache' -H 'Connection: keep-alive' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Cookie: %v' --data 'jobGroup=0&jobId=0&logStatus=1&filterTime=%v 00:00:00 - %v 23:59:59&start=0&length=10' --compressed --insecure", XXLURL, XXLCOOKIE, newtime, newtime)
	errordata := UniversalData{}
	curl := exec.Command("bash", "-c", geturl)
	out, err := curl.Output()
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	json.Unmarshal(out, &errordata)
	return float64(errordata.RecordsTotal)
}

//获取昨日成功的调度日志数据

func GetOldAccessData(XXLURL, XXLCOOKIE string) float64 {
	oldtime := fmt.Sprintf(time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	geturl := fmt.Sprintf("curl '%s/xxl-job-admin/joblog/pageList' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,be;q=0.6,sq;q=0.5' -H 'Cache-Control: no-cache' -H 'Connection: keep-alive' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Cookie: %v' --data 'jobGroup=0&jobId=0&logStatus=1&filterTime=%v 00:00:00 - %v 23:59:59&start=0&length=10' --compressed --insecure", XXLURL, XXLCOOKIE, oldtime, oldtime)
	accessdata := UniversalData{}
	curl := exec.Command("bash", "-c", geturl)
	out, err := curl.Output()
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	json.Unmarshal(out, &accessdata)
	return float64(accessdata.RecordsTotal)
}
