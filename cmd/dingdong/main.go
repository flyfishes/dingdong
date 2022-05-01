package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"dingdong/internal/app"
	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/session"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	logFile, err := os.OpenFile("dingdong.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	//defer logFile.Close()
	if err != nil {
		log.Printf("create log file error")
		return
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func main() {
	defer func() {
		if v := recover(); v != nil {
			log.Printf("[严重错误]: %+v", v)
		}
	}()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Println(dir)
	config.Initialize(dir + "/config.yml")
	session.Initialize()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		s := <-c
		log.Printf("Got signal: %+v", s)
		cancel()
	}()

	app.Run(ctx)
}
