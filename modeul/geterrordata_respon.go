package modeul

import "time"

type AllErrorData struct {
	NewTotal         int
	NewErrorTotal    int
	NewAvgErrorTotal int
	OldTotal         int
	OldErrorTotal    int
	OldAvgErrorTotal int
}

type UniversalData struct {
	RecordsFiltered int `json:"recordsFiltered"`
	Data            []struct {
		Id                     int         `json:"id"`
		JobGroup               int         `json:"jobGroup"`
		JobId                  int         `json:"jobId"`
		ExecutorAddress        interface{} `json:"executorAddress"`
		ExecutorHandler        string      `json:"executorHandler"`
		ExecutorParam          string      `json:"executorParam"`
		ExecutorShardingParam  interface{} `json:"executorShardingParam"`
		ExecutorFailRetryCount int         `json:"executorFailRetryCount"`
		TriggerTime            time.Time   `json:"triggerTime"`
		TriggerCode            int         `json:"triggerCode"`
		TriggerMsg             string      `json:"triggerMsg"`
		HandleTime             interface{} `json:"handleTime"`
		HandleCode             int         `json:"handleCode"`
		HandleMsg              interface{} `json:"handleMsg"`
		AlarmStatus            int         `json:"alarmStatus"`
	} `json:"data"`
	RecordsTotal int `json:"recordsTotal"`
}

//type GetAllData struct {
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
