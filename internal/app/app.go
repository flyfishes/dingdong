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
	"dingdong/internal/app/service/notify"
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
	if len(conf.Bark) > 0 {
		for _, v := range conf.Bark {
			notify.Push(v, "start run dingdong helper")
		}
	}
	if len(conf.PushPlus) > 0 {
		for _, v := range conf.PushPlus {
			notify.PushPlus(v, "start run dingdong helper")
		}
	}

	err := http.ListenAndServe(conf.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe => ", err)
	}
}

func isPeak() bool {
	now := time.Now()
	if now.Hour() >= 0 && now.Hour() < 6 {
		return true
	}
	if now.Hour() == 6 && now.Minute() < 10 {
		return true
	}
	if now.Hour() == 8 && now.Minute() < 41 && now.Minute() > 20 {
		return true
	}
	if now.Hour() >= 19 {
		return true
	}
	return false
}

// Monitor 监视器 监听运力
func Monitor() {
	cartMap := service.MockCartMap()
	for {
		conf := config.Get()
		duration := conf.MonitorIntervalMin + rand.Intn(conf.MonitorIntervalMax-conf.MonitorIntervalMin)
		bMonitored := 0 // 二进制最右1位区分监听，其他几位区分高峰期
		for i := 0; i < duration/5; i++ {
			<-time.After(time.Duration(5) * time.Second)
			//<-time.After(time.Second) ////定时1s，阻塞1s,1s后产生一个事件，往channel写内容
			if !conf.MonitorNeeded && !conf.PickUpNeeded {
				continue
			}
			/* // 每5分钟第1秒运行一次
			now := time.Now()
			if now.Second() != 0 {
				continue
			} */
			if isPeak() {
				if bMonitored < 2 {
					log.Println("当前高峰期或暂未营业.")
				}
				bMonitored += 2
				continue
			}
			if bMonitored%2 == 0 {
				bMonitored += 1
				service.MonitorAndPickUp(cartMap)
				//				_ = bMonitored
			}
		}

		//<-time.After(time.Duration(duration) * time.Second)
	}
}
