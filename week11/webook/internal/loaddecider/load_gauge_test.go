package loaddecider

import (
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRandomLoadGauge_Load(t *testing.T) {
	lg := ioc.InitLogger()
	rg := NewRandomLoadGauge(lg, NewNopReporter(), &IpPortIdentity{port: 1234}, time.Millisecond*20)

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
