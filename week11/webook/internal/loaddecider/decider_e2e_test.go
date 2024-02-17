package loaddecider

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisZSetDecider_IsVictory_WhenNotExpired(t *testing.T) {
	lg := ioc.InitLogger()

	redis := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	reporter := NewRedisZSetReporter(redis)
	identity := NewIpPortIdentity(1234)

	before(t, lg, reporter, identity)

	decider := NewRedisZSetDecider(redis, lg, identity)

	success := decider.IsVictory(context.Background())
	assert.True(t, success, "选中自己")
}

func TestRedisZSetDecider_IsVictory_WhenIsExpired(t *testing.T) {
	lg := ioc.InitLogger()

	redis := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	reporter := &RedisZSetReporter{redis: redis, expiration: time.Second}
	identity := NewIpPortIdentity(1234)

	rg := before(t, lg, reporter, identity)

	decider := NewRedisZSetDecider(redis, lg, identity)

	// 测试重复 close 不能有问题.
	rg.Close()

	// 让 key 超时.
	time.Sleep(time.Second * 2)

	// 会超时后降级，被自己处理.
	success := decider.IsVictory(context.Background())
	assert.True(t, success, "选中自己")
}

func before(t *testing.T, lg logger.LoggerV1, reporter Reporter, identity Identity) Gauge {
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

	return rg
}
