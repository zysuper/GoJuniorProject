package code

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrCodeSendTooMany   = errors.New("发送太频繁")
	ErrCodeVerifyTooMany = errors.New("验证太频繁")
)

// CodeCache 1. 定义一个 CodeCache 接口，
type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

func key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
