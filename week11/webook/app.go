package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/events"
	"gitee.com/geekbang/basic-go/webook/internal/loaddecider"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
	cron      *cron.Cron
	guage     loaddecider.Gauge
	decider   loaddecider.Decider
}
