package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"xxj_exporter/modeul"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Config    modeul.Config
	namespace = flag.String("xxl-job", "xxljob", "prometheus namespace")

	//XXLURL    = os.Getenv("XXLURL")
	//XXLCOOKIE = os.Getenv("XXLCOOKIE")
	XXLURL = "http://172.26.11.106:8080"
	//XXLURL = "http://xxl-dev.ur.com.cn"
	XXLCOOKIE = "_hid=BgAfrENekWN3FbhyR1XpPgA; Hm_lvt_1e14777bd42725d12c1934650d5bced6=1670471237; AdminUserKey=5e6f389b59a2bb4a; XXL_JOB_LOGIN_IDENTITY=7b226964223a312c22757365726e616d65223a2261646d696e222c2270617373776f7264223a226531306164633339343962613539616262653536653035376632306638383365222c22726f6c65223a312c227065726d697373696f6e223a6e756c6c7d"
)

func main() {
	metrics := modeul.NewMetrics(*namespace, XXLURL, XXLCOOKIE)
	registry := prometheus.NewRegistry()
	registry.MustRegister(metrics)
	registry.MustRegister(modeul.NewXxlJob(XXLURL, XXLCOOKIE))
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	log.Fatal(http.ListenAndServe(":9393", nil))

}

func init() {
	cnf, err := os.Open("./conf/config.yaml")
	if err != nil {
		log.Fatalf("ERR: 配置文件异常: %s", err)
		err := cnf.Close()
		if err != nil {
			return
		}
	}
	cnfdec := yaml.NewDecoder(cnf)
	if err := cnfdec.Decode(&Config); err != nil {
		log.Fatalf("ERR: 读取配置文件异常: %s", err)
	}
}
