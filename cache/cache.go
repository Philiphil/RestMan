package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"reflect"

	"github.com/philiphil/restman/orm/entity"
	"github.com/redis/go-redis/v9"
)

type Cache[E entity.Entity] interface {
	Set(ent E) error
	Get(ent E) (E, error)
	Delete(ent E) error
}

type RedisCache[E entity.Entity] struct {
	Client       *redis.Client
	entityPrefix string
	lifetime     time.Duration
}

func NewRedisCache[E entity.Entity](addr, password string, db int, lifetime int) *RedisCache[E] {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	var example E
	entityPrefix := reflect.TypeOf(example).Name()

	return &RedisCache[E]{
		Client:       client,
		entityPrefix: entityPrefix,
		lifetime:     time.Duration(lifetime) * time.Second,
	}
}

func (r *RedisCache[E]) generateCacheKey(ent entity.Entity) string {
	return fmt.Sprintf("%s:%s", r.entityPrefix, ent.GetId().String())
}

func (r *RedisCache[E]) Set(ent E) error {
	key := r.generateCacheKey(ent)
	data, err := json.Marshal(ent)
	if err != nil {
		return err
	}

	return r.Client.Set(context.Background(), key, data, r.lifetime).Err()
}

func (r *RedisCache[E]) Get(ent E) (E, error) {
	var result E
	key := r.generateCacheKey(ent)

	data, err := r.Client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return result, fmt.Errorf("cache miss for key: %s", key)
	} else if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(data), &result)
	return result, err
}

func (r *RedisCache[E]) Delete(ent E) error {
	key := r.generateCacheKey(ent)
	return r.Client.Del(context.Background(), key).Err()
}
