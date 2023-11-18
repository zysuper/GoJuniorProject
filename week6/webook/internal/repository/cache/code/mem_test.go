package code

import (
	"context"
	"github.com/coocood/freecache"
	assert "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const biz = "login"
const phone = "123456789012"
const code = "456789"

func TestCodeMemoryCache_SetAndVerifyIn3Times(t *testing.T) {
	cache := freecache.NewCache(1024)
	codeCache := NewMemCodeCache(cache)
	cxt := context.Background()

	a := assert.New(t)

	err := codeCache.Set(cxt, biz, phone, code)

	a.Nil(err, "not has error")

	verify, err := codeCache.Verify(cxt, biz, phone, "273845")
	a.Equal(verify, false, "not verify success")

	verify, err = codeCache.Verify(cxt, biz, phone, "273845")
	a.Equal(verify, false, "not verify success")

	verify, err = codeCache.Verify(cxt, biz, phone, code)
	a.Equal(verify, true, "not verify success")

}

func TestCodeMemoryCache_SetAndVerifyMoreTime(t *testing.T) {
	cache := freecache.NewCache(1024)
	codeCache := NewMemCodeCache(cache)
	cxt := context.Background()

	a := assert.New(t)

	err := codeCache.Set(cxt, biz, phone, code)

	a.Nil(err, "not has error")

	verify, err := codeCache.Verify(cxt, biz, phone, "273845")
	a.Equal(verify, false, "not verify success")

	verify, err = codeCache.Verify(cxt, biz, phone, "273845")
	a.Equal(verify, false, "not verify success")

	verify, err = codeCache.Verify(cxt, biz, phone, "273845")
	a.Equal(verify, false, "not verify success")

	verify, err = codeCache.Verify(cxt, biz, phone, code)
	a.Equal(verify, false, "not verify success")

}

func TestCodeMemoryCache_MoreSet(t *testing.T) {
	cache := freecache.NewCache(1024)
	codeCache := NewMemCodeCache(cache)
	cxt := context.Background()

	a := assert.New(t)

	err := codeCache.Set(cxt, biz, phone, code)

	a.Nil(err, "not has error")

	err = codeCache.Set(cxt, biz, phone, code)

	a.NotNil(err, "has error of 发送太频繁")
}

func TestCodeMemoryCache_WaitLaterSet(t *testing.T) {
	cache := freecache.NewCache(1024)
	codeCache := NewMemCodeCache(cache)
	cxt := context.Background()

	a := assert.New(t)

	// 第一次...
	err := codeCache.Set(cxt, biz, phone, code)
	a.Nil(err, "not has error")

	// 第二次立刻设置，应该失败.
	err = codeCache.Set(cxt, biz, phone, code)
	a.NotNil(err, "has error of 发送太频繁")

	// 5s 再设置依旧失败.
	time.Sleep(5 * time.Second)
	err = codeCache.Set(cxt, biz, phone, code)

	// 60s 后应该成功.
	time.Sleep(56 * time.Second)
	err = codeCache.Set(cxt, biz, phone, code)

	a.Nil(err, "has error of 发送太频繁")
}
