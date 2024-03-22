package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
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
		})
	}

	return &Picker{
		conns: conns,
	}
}

type Picker struct {
	conns []*weightConn
	lock  sync.Mutex
	// 动态加权值
	dynWeight int
	// weight 最大值
	maxWeight int
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	// 总权重
	var total int
	var maxCC *weightConn
	for _, c := range p.conns {
		total += c.weight
		c.currentWeight = c.currentWeight + c.weight
		if maxCC == nil || maxCC.currentWeight < c.currentWeight {
			maxCC = c
		}
	}

	maxCC.currentWeight = maxCC.currentWeight - total

	return balancer.PickResult{
		SubConn: maxCC.SubConn,
		Done: func(info balancer.DoneInfo) {
			// 要在这里进一步调整weight/currentWeight
			if info.Err != nil {
				// 降级,但是最多降到 0.
				maxCC.currentWeight = max(maxCC.currentWeight-p.dynWeight, 0)
			} else {
				// 提高权重,但是不大于最大权重.
				// 如果不设置最大权重，那么就会导致只有这个节点被调用了.
				maxCC.currentWeight = min(maxCC.currentWeight+p.dynWeight, p.maxWeight)
			}
		},
	}, nil

}

type weightConn struct {
	balancer.SubConn
	weight        int
	currentWeight int

	// 可以用来标记不可用
	available bool
}
