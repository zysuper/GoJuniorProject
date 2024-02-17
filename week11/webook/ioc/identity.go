package ioc

import "gitee.com/geekbang/basic-go/webook/internal/loaddecider"

func NewIdentity() loaddecider.Identity {
	return loaddecider.NewIpPortIdentity(8080)
}
