package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

const Name = "custom_weighted_round_robin"

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &PickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type PickerBuilder struct {
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	//fmt.Printf("build call -> %v\n", len(info.ReadySCs))
	conns := make([]*weightConn, 0, len(info.ReadySCs))
	for sc, sci := range info.ReadySCs {
		md, _ := sci.Address.Metadata.(map[string]any)
		weightVal, _ := md["weight"]
		weight, _ := weightVal.(float64)
		//if weight == 0 {
		//
		//}
		conns = append(conns, &weightConn{
			SubConn:       sc,
			weight:        int(weight),
			currentWeight: int(weight),
			available:     true,
		})
	}

	return &Picker{
		conns: conns,
		// 默认 10
		dynWeight: 10,
		// 默认 2 倍 total weight,
		maxWeightBoost: 2,
	}
}

type Picker struct {
	conns []*weightConn
	lock  sync.Mutex
	// 动态加权值
	dynWeight int
	// weight 最大值
	maxWeightBoost int
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	//fmt.Printf("pick call -> %v\n", info)
	p.lock.Lock()
	defer p.lock.Unlock()
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	// 总权重
	var total int
	var maxCC *weightConn
	for _, c := range p.conns {
		// 如果这个节点已经熔断了，跳过...
		if c.available == false {
			continue
		}
		total += c.weight
		c.currentWeight = c.currentWeight + c.weight
		if maxCC == nil || maxCC.currentWeight < c.currentWeight {
			maxCC = c
		}
	}

	// 全部节点都被熔断了.
	if maxCC == nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	maxCC.currentWeight = maxCC.currentWeight - total

	return balancer.PickResult{
		SubConn: maxCC.SubConn,
		Done: func(info balancer.DoneInfo) {
			err := info.Err

			if err != nil {
				if statusError, ok := status.FromError(err); ok {
					switch statusError.Code() {
					case codes.Unavailable: // 熔断了.
						// 设置成不可用.
						p.setUnAvailable(maxCC)
						// 开启健康检查.
						go p.healthCheck(maxCC)
					case codes.ResourceExhausted: // 限流了.
						fallthrough
					default:
						p.defaultWeightAdjust(maxCC)
					}
					return
				}
			}

			// 没有错误，服务正常，提高权重.
			// 提高权重,但是不大于最大权重.
			// 如果不设置最大权重，那么就会导致只有这个节点被调用了.
			maxCC.currentWeight = min(maxCC.currentWeight+p.dynWeight, p.maxWeightBoost*total)
		},
	}, nil

}

func (p *Picker) setUnAvailable(maxCC *weightConn) {
	maxCC.lock.Lock()
	defer maxCC.lock.Unlock()
	maxCC.available = false
}

// defaultWeightAdjust 默认的基于动态 weight 的权重调整.
func (p *Picker) defaultWeightAdjust(maxCC *weightConn) {
	// 降级,但是最多降到 0.
	maxCC.currentWeight = max(maxCC.currentWeight-p.dynWeight, 0)
}

func (p *Picker) healthCheck(cc *weightConn) {
	// 创建一个5秒的定时器
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	<-timer.C

	// !! 尚未找到 picker 能够获得上层 balancer 健康检查的状态,
	// 所以仅仅只设置了一个断路器的暂停时间。

	cc.lock.Lock()
	cc.available = true
	// 短路一段时间后，试图恢复，当前权重很低。
	cc.currentWeight = 0
	cc.lock.Unlock()
}

type weightConn struct {
	balancer.SubConn
	weight        int
	currentWeight int

	// 可以用来标记不可用
	available bool
	lock      sync.Mutex
}
