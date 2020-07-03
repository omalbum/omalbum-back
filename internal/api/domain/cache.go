package domain

import "strconv"

type CacheKey string

const (
	AllProblemsCacheKey     CacheKey = "all_problems"
	CurrentProblemsCacheKey CacheKey = "current_problems"
	NextProblemsCacheKey    CacheKey = "next_problems"
)

func UserCachePrefix(userId uint) string {
	return "user_" + strconv.Itoa(int(userId))
}

func UserAlbumCacheKey(userId uint) CacheKey {
	return CacheKey(UserCachePrefix(userId) + "_album")
}
