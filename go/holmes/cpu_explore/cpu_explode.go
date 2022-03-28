package main

import (
	"github.com/mosn/holmes"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/cpuex", cpuex)
	go http.ListenAndServe(":10003", nil)
}

func main() {
	h, _ := holmes.New(
		holmes.WithCollectInterval("2s"),
		holmes.WithCoolDown("1m"),
		holmes.WithDumpPath("/tmp"),   // 在tmp目录下生成holmes.log和cpu.xxx.bin文件，cpu.xxx.bin文件只有在满足dump条件的时候输出
		holmes.WithCPUDump(1, 25, 80), // 当CPU使用率超过1%，波动25%，或者超过80%的CPU使用率则会Dump CPU的信息
	)
	h.EnableCPUDump()
	h.Start()
	time.Sleep(time.Hour)
}

func cpuex(wr http.ResponseWriter, req *http.Request) {
	go func() {
		for {
			time.Sleep(time.Millisecond)
		}
	}()
}
