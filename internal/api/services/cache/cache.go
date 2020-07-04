package cache

import (
	"github.com/karlseguin/ccache/v2"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"time"
)

func NewTeleOMACache() domain.TeleOMACache {
	return teleOMACache{cache: ccache.New(ccache.Configure())}
}

type teleOMACache struct {
	cache *ccache.Cache
}

func (t teleOMACache) SetWithTTL(key domain.CacheKey, value interface{}, timeToLive time.Duration) {
	t.cache.Set(string(key), value, timeToLive)
}

func (t teleOMACache) SetWithExpiration(key domain.CacheKey, value interface{}, expirationDate time.Time) {
	t.cache.Set(string(key), value, expirationDate.Sub(time.Now()))

}

func (t teleOMACache) Delete(key domain.CacheKey) {
	t.cache.Delete(string(key))
}

func (t teleOMACache) ClearUserCache(userId uint) int {
	return t.cache.DeletePrefix(domain.UserCachePrefix(userId))
}

func (t teleOMACache) Clear() int {
	return t.cache.DeletePrefix("")
}

func (t teleOMACache) Get(key domain.CacheKey) interface{} {
	item := t.cache.Get(string(key))
	if item == nil {
		return nil
	}
	if item.Expired() {
		return nil
	}
	return item.Value()
}

func (t teleOMACache) GetUserAlbum(userId uint) interface{} {
	return t.Get(domain.UserAlbumCacheKey(userId))
}
