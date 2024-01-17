package article

import "gitee.com/geekbang/basic-go/webook/internal/events/article/prometheus"

var gauge = prometheus.NewGauge(
	"geektime_week10", "webook",
	"kafka_read_event", "re_gauge")
