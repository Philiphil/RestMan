package cache_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	. "github.com/philiphil/restman/cache"
	"github.com/philiphil/restman/orm/entity"
)

type MockEntity struct {
	ID entity.ID `json:"id"`
}

func (m MockEntity) GetId() entity.ID {
	return m.ID
}

func (m MockEntity) SetId(test any) entity.Entity {
	panic(test)
}

func TestRedisCache_Set(t *testing.T) {
	client, mock := redismock.NewClientMock()
	c := NewRedisCache[MockEntity]("localhost:6379", "", 0, 60)
	c.Client = client

	ent := MockEntity{ID: entity.CastId(123)}
	key := fmt.Sprintf("MockEntity:%s", ent.GetId())
	data, _ := json.Marshal(ent)

	mock.ExpectSet(key, data, 60*time.Second).SetVal("OK")

	err := c.Set(ent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations were not met: %v", err)
	}
}

func TestRedisCache_Get(t *testing.T) {
	client, mock := redismock.NewClientMock()
	c := NewRedisCache[MockEntity]("localhost:6379", "", 0, 60)
	c.Client = client

	ent := MockEntity{ID: entity.CastId(123)}
	key := fmt.Sprintf("MockEntity:%s", ent.GetId())
	data, _ := json.Marshal(ent)

	mock.ExpectGet(key).SetVal(string(data))

	result, err := c.Get(ent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != ent {
		t.Fatalf("expected %v, got %v", ent, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations were not met: %v", err)
	}
}

func TestRedisCache_Get_CacheMiss(t *testing.T) {
	client, mock := redismock.NewClientMock()
	c := NewRedisCache[MockEntity]("localhost:6379", "", 0, 60)
	c.Client = client

	ent := MockEntity{ID: entity.CastId(123)}
	key := fmt.Sprintf("MockEntity:%s", ent.GetId())

	mock.ExpectGet(key).RedisNil()

	_, err := c.Get(ent)
	if err == nil || err.Error() != fmt.Sprintf("cache miss for key: %s", key) {
		t.Fatalf("expected cache miss error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations were not met: %v", err)
	}
}

func TestRedisCache_Delete(t *testing.T) {
	client, mock := redismock.NewClientMock()
	c := NewRedisCache[MockEntity]("localhost:6379", "", 0, 60)
	c.Client = client

	ent := MockEntity{ID: entity.CastId(123)}
	key := fmt.Sprintf("MockEntity:%s", ent.GetId())

	mock.ExpectDel(key).SetVal(1)

	err := c.Delete(ent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations were not met: %v", err)
	}
}
