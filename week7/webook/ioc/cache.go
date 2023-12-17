package ioc

import "github.com/coocood/freecache"

const cacheSize = 1024 * 1024 * 5

func InitLocalCache() *freecache.Cache {
	return freecache.NewCache(cacheSize)
}
