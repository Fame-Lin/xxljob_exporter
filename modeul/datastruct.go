package modeul

import (
	"time"
)

// xxl-job调度日志中成功日志的数据结构体

//type AccessData struct {
//	RecordsFiltered int `json:"recordsFiltered"`
//	Data            []struct {
//		Id                     int         `json:"id"`
//		JobGroup               int         `json:"jobGroup"`
//		JobId                  int         `json:"jobId"`
//		ExecutorAddress        interface{} `json:"executorAddress"`
//		ExecutorHandler        string      `json:"executorHandler"`
//		ExecutorParam          string      `json:"executorParam"`
//		ExecutorShardingParam  interface{} `json:"executorShardingParam"`
//		ExecutorFailRetryCount int         `json:"executorFailRetryCount"`
//		TriggerTime            time.Time   `json:"triggerTime"`
//		TriggerCode            int         `json:"triggerCode"`
//		TriggerMsg             string      `json:"triggerMsg"`
//		HandleTime             interface{} `json:"handleTime"`
//		HandleCode             int         `json:"handleCode"`
//		HandleMsg              interface{} `json:"handleMsg"`
//		AlarmStatus            int         `json:"alarmStatus"`
//	} `json:"data"`
//	RecordsTotal int `json:"recordsTotal"`
//}

// xxl-job查看执行器管理的数据结构体

type AllGroupData struct {
	RecordsFiltered int            `json:"recordsFiltered"`
	Data            []allGroupData `json:"data"`
	RecordsTotal    int            `json:"recordsTotal"`
}

type allGroupData struct {
	Id           int         `json:"id"`
	Appname      string      `json:"appname"`
	Title        string      `json:"title"`
	AddressType  int         `json:"addressType"`
	AddressList  interface{} `json:"addressList"`
	UpdateTime   time.Time   `json:"updateTime"`
	RegistryList interface{} `json:"registryList"`
}

// 所有指标的数据结构体

type AllData struct {
	NewTotal          int
	NewErrorTotal     int
	NewAccessTotal    int
	NewAvgErrorTotal  int
	NewAvgAccessTotal int
	OldTotal          int
	OldErrorTotal     int
	OldAccessTotal    int
	OldAvgErrorTotal  int
	OldAvgAccessTotal int
}

// Config config struct
type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Log      Log      `yaml:"log"`
}

// Server Config
type Server struct {
	Port  int    `yaml:"port"`
	Hosts string `yaml:"hosts"`
}

type Database struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
type Log struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}
