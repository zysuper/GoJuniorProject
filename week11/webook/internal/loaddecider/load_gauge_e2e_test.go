package loaddecider

import (
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRandomLoadGauge_Load_e2e(t *testing.T) {
	lg := ioc.InitLogger()
	redis := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	reporter := NewRedisZSetReporter(redis)
	identity := NewIpPortIdentity(1234)
	rg := NewRandomLoadGauge(lg, reporter, identity, time.Millisecond*20)

	time.Sleep(time.Millisecond * 22)
	load := rg.Load()

	time.Sleep(time.Millisecond * 22)
	load2 := rg.Load()

	t.Logf("load1: %v -> load2: %v", load, load2)
	assert.NotEqualf(t, load, load2, "random again not equal")

	rg.Close()
	time.Sleep(time.Millisecond * 21)
	load3 := rg.Load()
	assert.Equalf(t, load3, load2, "random again close then not change")
}
