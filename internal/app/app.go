package app

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"dingdong/assets"
	"dingdong/internal/app/api"
	"dingdong/internal/app/config"
	"dingdong/internal/app/service"
)

func Run() {
	go Monitor()
	go service.SnapUp()
	go service.PickUp()
	go service.Notify()

	http.HandleFunc("/", api.ConfigView)
	http.HandleFunc("/set", api.SetConfig)
	http.HandleFunc("/config", api.ConfigView)
	http.HandleFunc("/notify", api.Notify)
	http.HandleFunc("/address", api.GetAddress)
	http.HandleFunc("/addOrder", api.AddOrder)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets.FS))))

	conf := config.Get()
	err := http.ListenAndServe(conf.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe => ", err)
	}
}

func isPeak() bool {
	now := time.Now()
	if now.Hour() >= 0 && now.Hour() < 9 {
		return true
	}
	return false
}

// Monitor 监视器 监听运力
func Monitor() {
	cartMap := service.MockCartMap()
	for {
		//<-time.After(time.Second) ////定时1s，阻塞1s,1s后产生一个事件，往channel写内容
		conf := config.Get()
		if !conf.MonitorNeeded && !conf.PickUpNeeded {
			continue
		}
		/* // 每5分钟第1秒运行一次
		now := time.Now()
		if now.Minute()%5 != 0 || now.Second() != 1 {
			continue
		} */
		if isPeak() {
			log.Println("当前高峰期或暂未营业")
			continue
		}
		service.MonitorAndPickUp(cartMap)
		duration := conf.MonitorIntervalMin + rand.Intn(conf.MonitorIntervalMax-conf.MonitorIntervalMin)
		<-time.After(time.Duration(duration) * time.Second)
	}
}
