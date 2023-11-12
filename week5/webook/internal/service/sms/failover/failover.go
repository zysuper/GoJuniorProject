package failover

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"log"
	"sync/atomic"
)

type FailOverSmsService struct {
	svcs []sms.SmsService
	idx  uint64
}

var AllSendFailed = errors.New("所有 sms 服务商都没法提供服务")

func NewFailOverSmsService(svcs []sms.SmsService) sms.SmsService {
	return &FailOverSmsService{svcs: svcs}
}

func (f *FailOverSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		e := svc.Send(ctx, tplId, args, numbers...)
		if e == nil {
			return e
		}
	}
	return AllSendFailed
}

// 起始下标轮询
// 并且出错也轮询
func (f *FailOverSmsService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	// 我要迭代 length
	for i := idx; i < idx+length; i++ {
		// 取余数来计算下标
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			// 前者是被取消，后者是超时
			return err
		}
		log.Println(err)
	}
	return errors.New("轮询了所有的服务商，但是发送都失败了")
}
