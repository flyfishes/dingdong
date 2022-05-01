package app

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"dingdong/assets"
	"dingdong/internal/app/api"
	"dingdong/internal/app/config"
	"dingdong/internal/app/service"
)

var (
	pickUpCh         = make(chan struct{})
	dingDongNotifyCh = make(chan struct{})
	merTuanNotifyCh  = make(chan struct{})
)

func Run(ctx context.Context) {
	go service.SnapUp(ctx)
	go service.PickUp(ctx, pickUpCh)
	go service.DingDongNotify(dingDongNotifyCh)
	go service.MeiTuanNotify(merTuanNotifyCh)

	go DingDongMonitor(dingDongNotifyCh, pickUpCh)
	go MeiTuanMonitor(merTuanNotifyCh)

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

// DingDongMonitor 监视器 监听运力
func DingDongMonitor(notifyCh chan<- struct{}, pickUpCh chan<- struct{}) {
	cartMap := service.MockCartMap()
	for {
		conf := config.GetDingDong()
		duration := conf.MonitorIntervalMin + rand.Intn(conf.MonitorIntervalMax-conf.MonitorIntervalMin)
		bMonitored := 0 // 二进制最右1位区分监听，其他几位区分高峰期
		for i := 0; i < duration/5; i++ {
			<-time.After(time.Duration(5) * time.Second)
			if !conf.MonitorNeeded && !conf.PickUpNeeded {
				continue
			}

			if dingDongIsPeak() {
				if bMonitored < 2 {
					log.Println("叮咚当前高峰期或暂未营业.")
				}
				bMonitored += 2
				continue
			}
			if bMonitored%2 == 0 {
				bMonitored += 1
				service.MonitorAndPickUp(cartMap, notifyCh, pickUpCh)
				//				_ = bMonitored
			}
		}
	}
}

func dingDongIsPeak() bool {
	now := time.Now()
	if now.Hour() >= 0 && now.Hour() < 6 {
		return true
	}
	if now.Hour() == 6 && now.Minute() < 15 {
		return true
	}
	if now.Hour() == 8 && now.Minute() < 45 && now.Minute() > 20 {
		return true
	}
	if now.Hour() >= 19 {
		return true
	}
	return false
}

func MeiTuanMonitor(notifyCh chan<- struct{}) {
	for {
		conf := config.GetMeiTuan()
		if !conf.MonitorNeeded {
			<-time.After(time.Second)
			continue
		}
		now := time.Now()
		if now.Second() != 0 {
			continue
		}
		// 每分钟在第1-3秒运行一次
		random := rand.Intn(3) + 1
		<-time.After(time.Second * time.Duration(random))
		if meiTuanIsPeak() {
			log.Println("美团当前高峰期或暂未营业")
			continue
		}
		service.MeiTuanMonitorAndNotify(notifyCh)
	}
}

func meiTuanIsPeak() bool {
	now := time.Now()
	if now.Hour() >= 0 && now.Hour() < 6 {
		return true
	}
	if now.Hour() == 6 && now.Minute() < 15 {
		return true
	}
	if now.Hour() >= 22 {
		return true
	}
	return false
}
