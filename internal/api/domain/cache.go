package domain

import (
	"strconv"
	"time"
)

type CacheKey string

const (
	AllProblemsCacheKey     CacheKey = "all_problems"
	CurrentProblemsCacheKey CacheKey = "current_problems"
	NextProblemsCacheKey    CacheKey = "next_problems"
)

type TeleOMACache interface {
	Get(key CacheKey) interface{}
	SetWithTTL(key CacheKey, value interface{}, timeToLive time.Duration)
	SetWithExpiration(key CacheKey, value interface{}, expirationDate time.Time)
	Delete(key CacheKey)
	Clear() int // clears everything

	GetUserAlbum(userId uint) interface{}
	ClearUserCache(userId uint) int //clears the user's cache
}

func UserAlbumCacheKey(userId uint) CacheKey {
	return CacheKey(UserCachePrefix(userId) + "_album")
}

func UserCachePrefix(userId uint) string {
	return "user_" + strconv.Itoa(int(userId))
}
