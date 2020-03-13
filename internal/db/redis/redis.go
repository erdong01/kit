package redis

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"rxt/internal/cache/uitl"
	"rxt/internal/db/redis/check"
	"rxt/internal/db/redis/drive"
)

type redisCache struct {
	Conn *redis.Client
}

func NewCache() *redisCache {
	var conn *redis.Client
	if coreReids := check.Connect(); coreReids != nil {
		conn = coreReids
	} else {
		conn = drive.New()
	}
	return &redisCache{
		Conn: conn,
	}
}

func (r *redisCache) Has(key string) bool {
	key = uitl.PrefixKey(key)
	if 1 == r.Conn.Exists(key).Val() {
		return true
	}
	return false
}

func (r *redisCache) Get(key string, v interface{}) error {
	key = uitl.PrefixKey(key)
	return r.Conn.Get(key).Scan(v)
}

func (r *redisCache) GetJSON(key string, v interface{}) error {
	key = uitl.PrefixKey(key)
	b, err := r.Conn.Get(key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (r *redisCache) Set(key string, v interface{}, expireInt int) error {
	key = uitl.PrefixKey(key)
	expire := uitl.ExpireDuration(expireInt)
	_, err := r.Conn.Set(key, v, expire).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) SetJSON(key string, v interface{}, expireInt int) error {
	key = uitl.PrefixKey(key)
	expire := uitl.ExpireDuration(expireInt)
	str, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = r.Conn.Set(key, str, expire).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) GetBytes(key string) ([]byte, error) {
	return r.Conn.Get(key).Bytes()
}

func (r *redisCache) Del(key string) error {
	return r.Conn.Del(key).Err()
}
