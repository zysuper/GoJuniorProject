package loaddecider

import (
	"fmt"
	"gitee.com/geekbang/basic-go/webook/pkg/netx"
)

type Identity interface {
	Id() string
}

type IpPortIdentity struct {
	port    int64
	idCache string
}

func (i *IpPortIdentity) Id() string {
	return fmt.Sprintf("%s:%d", i.idCache, i.port)
}

func NewIpPortIdentity(port int64) Identity {
	return &IpPortIdentity{port: port, idCache: netx.GetOutboundIP()}
}
