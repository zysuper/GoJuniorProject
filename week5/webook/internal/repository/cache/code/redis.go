package code

import (
	"context"
	_ "embed"
	"errors"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string
)

// CodeRedisCache 2. 将现在的 CodeCache 改名为 CodeRedisCache
type CodeRedisCache struct {
	cmd redis.Cmdable
}

func (c *CodeRedisCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{key(biz, phone)}, code).Int()
	// 打印日志
	if err != nil {
		// 调用 redis 出了问题
		return err
	}
	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (c *CodeRedisCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{key(biz, phone)}, code).Int()
	// 打印日志
	if err != nil {
		// 调用 redis 出了问题
		return false, err
	}

	switch res {
	case -2:
		return false, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return true, nil
	}
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &CodeRedisCache{cmd: cmd}
}
