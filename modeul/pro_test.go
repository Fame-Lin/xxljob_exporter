package modeul

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

func (g *GroupData) GetGroupid() NewGroupData {
	startTime := time.Now()
	geturl := fmt.Sprintf("curl '%s/xxl-job-admin/jobgroup/pageList' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,be;q=0.6,sq;q=0.5' -H 'Cache-Control: no-cache' -H 'Connection: keep-alive' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Cookie: %v' --data 'ppname=&title=&start=0&length=100' --compressed --insecure", g.XxlUrl, g.xxlCookie)
	groupdata := AllGroupData{}
	curl := exec.Command("bash", "-c", geturl)
	out, err := curl.Output()
	if err != nil {
		fmt.Printf("GroupIderr: %v", err)
	}
	json.Unmarshal(out, &groupdata)
	newdate := NewGroupData{RecordsTotal: groupdata.RecordsTotal}

	wg := &sync.WaitGroup{}
	for _, value := range groupdata.Data {
		wg.Add(1)
		valued := value

		go func() {
			group := groupdate{}
			group.Id = valued.Id
			group.AppName = valued.Appname
			group.Title = valued.Title
			group.ALLOleTotal, group.ALlNewTotal = GetGroupTotal(valued.Id, g.XxlUrl, g.xxlCookie)
			group.OldAccessTotal, group.NewAccessTotal = GetAccessGroupData(valued.Id, g.XxlUrl, g.xxlCookie)
			group.OldErrorTotal, group.NewErrorTotal = GetErrorGroupData(valued.Id, g.XxlUrl, g.xxlCookie)
			newdate.Data = append(newdate.Data, group)
			//fmt.Println(group)
		}()
	}
	wg.Wait()
	endTime := time.Now()
	fmt.Printf("Time: %v", endTime.Sub(startTime))
	fmt.Println(newdate)
	return newdate
}
