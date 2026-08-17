// this repository only knows the like
// how to store url in redis
// how to retrive cacked url
// how to increament hits
// how to read hits
// also it doesnt know the about anything http or business logic
package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"url-shorter/config"
)

//key naming conventions, keeps Redis keys organized
func urlKey(code string) string { return fmt.Sprintf("url:%s", code)}
func hitsKey(code string) string { return fmt.Sprintf("url:%s", code)}

//cacheURL stores a short code, original URL mapping in Redis
//expires after 24 hrs, after that , Go fetches from postgres and recaches

func CacheURL(code, original string) error {
	return config.RDB.Set(
		context.Background(),
		urlKey(code),
		original,
		24*time.Hour, //TTL- time to live
	).Err()
}

//GetCachedURL fetches original URL from Redis
//returning ("", false) if not in cache (cache miss)
func GetCachedURL(code string) (string, bool) {
	val, err := config.RDB.Get(context.Background(), urlKey(code)).Result()
	if err != nil {
		return "", false //cache miss
	}
	return val, true //cache hit
}

//increamentHits automatically increament the hit counter for a code
//INCR is atomic, even if 1000 users hit this simultaneously,
//every single increamet is counted correctly 9no face conditions)
//INCR = "Redis, increase this number by exactly 1, safely."
func IncreametHits(code string) {
	config.RDB.Incr(context.Background(), hitsKey(code))
}

//GetHits returns the urrent hit count for a short code
func GetHits(code string) int64 {
	val, err := config.RDB.Get(context.Background(), hitsKey(code)).Result()
	if err != nil {
		return 0
	}
	hits, _ := strconv.ParseInt(val, 10, 64)
	return hits
}