package code

import (
	"context"
	"github.com/coocood/freecache"
	"strconv"
	"sync"
)

const expireTime = 600
const interval = 60
const cnt = 3

type CodeMemoryCache struct {
	cache *freecache.Cache
	mux   sync.Mutex
}

func (c *CodeMemoryCache) Set(ctx context.Context, biz, phone, code string) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	k := keyBytes(biz, phone)
	kc := keyCountBytes(biz, phone)

	ttl, err := c.cache.TTL(k)

	// 完全不存在,准了.
	if err != nil && err == freecache.ErrNotFound {
		c.newCache(k, code, kc)
		return nil
	}

	// 很频繁到再次获取 code.
	if ttl >= expireTime-interval {
		return ErrCodeSendTooMany
	}

	// ttl 不是很频繁,准了.
	c.newCache(k, code, kc)

	return nil
}

func (c *CodeMemoryCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	k := keyBytes(biz, phone)
	kc := keyCountBytes(biz, phone)

	cnt, err := c.cache.Get(kc)
	expectedCode, err2 := c.cache.Get(k)

	if err != nil {
		// 没查到？
		return false, err
	}

	if err2 != nil {
		// 奇怪的点.
		return false, err2
	}

	cc, err := toInt(cnt)

	if err != nil {
		// 非法记数.
		return false, err
	}

	if cc <= 0 {
		// 记数耗尽了，亲.
		return false, nil
	}

	if code == string(expectedCode) {
		// 验证通过...
		return true, nil
	}

	// decr...
	c.cache.Set(kc, toBytes(cc-1), expireTime)

	return false, nil
}

func toInt(cnt []byte) (int, error) {
	return strconv.Atoi(string(cnt))
}

func toBytes(cnt int) []byte {
	return []byte(strconv.Itoa(cnt))
}

func keyBytes(biz string, phone string) []byte {
	return []byte(key(biz, phone))
}

func keyOfCnt(prefix string) string {
	return prefix + ":cnt"
}

func NewMemCodeCache(cache *freecache.Cache) CodeCache {
	return &CodeMemoryCache{
		cache: cache,
	}
}

func keyCountBytes(biz string, phone string) []byte {
	return []byte(keyOfCnt(key(biz, phone)))
}

func (c *CodeMemoryCache) newCache(k []byte, code string, kc []byte) {
	c.cache.Set(k, []byte(code), expireTime)
	c.cache.Set(kc, []byte(strconv.Itoa(cnt)), expireTime)
}
