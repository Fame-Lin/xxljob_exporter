package modeul

import "github.com/prometheus/client_golang/prometheus"

//定义一个指标结构体

type XxlJob struct {
	ALLOldTotal       *prometheus.Desc
	AllNewTotal       *prometheus.Desc
	AllOldErrorTotal  *prometheus.Desc
	AllOldAccessTotal *prometheus.Desc
	AllNewErrorTotal  *prometheus.Desc
	AllNewAccessTotal *prometheus.Desc
	XxlJobHeartbeat   *prometheus.Desc
	LabelVaues        []string
	XxlUrl            string
	XxlCookie         string
}

type XxlJobGroup struct {
}

// 创建指标描述符

func NewXxlJob(XXLURL, XXLCOOKIE string) *XxlJob {
	return &XxlJob{
		AllNewTotal:       prometheus.NewDesc("xxljob_old_total", "获取今天任务总数量", nil, nil),
		ALLOldTotal:       prometheus.NewDesc("xxljob_new_total", "获取昨天任务总数量", nil, nil),
		AllOldErrorTotal:  prometheus.NewDesc("xxljob_old_error_total", "获取昨天错误任务总数量", nil, nil),
		AllNewErrorTotal:  prometheus.NewDesc("xxljob_new_error_total", "获取今天错误任务总数量", nil, nil),
		AllOldAccessTotal: prometheus.NewDesc("xxljob_old_access_total", "获取昨天成功任务总数量", nil, nil),
		AllNewAccessTotal: prometheus.NewDesc("xxljob_new_access_total", "获取今天成功任务总数量", nil, nil),
		XxlJobHeartbeat:   prometheus.NewDesc("xxljob_up", "XXL-JOB存活", nil, nil),
		LabelVaues:        nil,
		XxlUrl:            XXLURL,
		XxlCookie:         XXLCOOKIE,
	}
}

// 实现prometheus的Describe接口

func (x *XxlJob) Describe(ch chan<- *prometheus.Desc) {
	ch <- x.AllNewErrorTotal
	ch <- x.AllOldErrorTotal
	ch <- x.AllNewAccessTotal
	ch <- x.AllOldAccessTotal
	ch <- x.ALLOldTotal
	ch <- x.AllNewTotal
	ch <- x.XxlJobHeartbeat
}

// 实现prometheus的Collect接口

func (x XxlJob) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(x.AllOldErrorTotal, prometheus.GaugeValue, GetOldErrorData(x.XxlUrl, x.XxlCookie), x.LabelVaues...)
	ch <- prometheus.MustNewConstMetric(x.AllNewErrorTotal, prometheus.GaugeValue, GetNewErrorData(x.XxlUrl, x.XxlCookie), x.LabelVaues...)
	ch <- prometheus.MustNewConstMetric(x.AllNewAccessTotal, prometheus.GaugeValue, GetNewAccessData(x.XxlUrl, x.XxlCookie), x.LabelVaues...)
	ch <- prometheus.MustNewConstMetric(x.AllOldAccessTotal, prometheus.GaugeValue, GetOldAccessData(x.XxlUrl, x.XxlCookie), x.LabelVaues...)
	ch <- prometheus.MustNewConstMetric(x.AllNewTotal, prometheus.GaugeValue, GetOldData(x.XxlUrl, x.XxlCookie), x.LabelVaues...)
	ch <- prometheus.MustNewConstMetric(x.ALLOldTotal, prometheus.GaugeValue, GetNewData(x.XxlUrl, x.XxlCookie), x.LabelVaues...)
	ch <- prometheus.MustNewConstMetric(x.XxlJobHeartbeat, prometheus.GaugeValue, CheckHeartbeat(x.XxlUrl), x.LabelVaues...)
}
