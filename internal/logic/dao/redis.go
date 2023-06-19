package dao

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"time"
)

const (
	_prefixMidServer    = "mid_%d"
	_prefixKeyServer    = "key_%s"
	_prefixServerOnline = "ol_%s"
)

var (
	ErrStr = errors.New("redis: err string")
)

func KeyMidServer(mid int64) string {
	return fmt.Sprintf(_prefixMidServer, mid)
}

func KeyKeyServer(key string) string {
	return fmt.Sprintf(_prefixKeyServer, key)
}

func KeyServerOnline(key string) string {
	return fmt.Sprintf(_prefixServerOnline, key)
}

func (b *Base) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return BaseDao.Redis.Expire(ctx, key, expiration).Err()
}

func (b *Base) SetExpire(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return BaseDao.Redis.Set(ctx, key, value, expiration).Err()
}

func (b *Base) SetKeyExpire(ctx context.Context, key string, expiration time.Duration) error {
	return BaseDao.Redis.Expire(ctx, key, expiration).Err()
}

func (b *Base) HSet(ctx context.Context, key string, values ...interface{}) error {
	return BaseDao.Redis.HSet(ctx, key, values).Err()
}

func (b *Base) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return b.Redis.HGetAll(ctx, key).Result()
}

func (b *Base) HDel(ctx context.Context, key, field string) error {
	return b.Redis.HDel(ctx, key, field).Err()
}

func (b *Base) Del(ctx context.Context, key string) error {
	return b.Redis.Del(ctx, key).Err()
}

func (b *Base) HSetExpire(ctx context.Context, expire time.Duration, key string, values ...interface{}) error {
	if err := b.HSet(ctx, key, values...); err != nil {
		return err
	}
	if err := b.Expire(ctx, key, expire); err != nil {
		return err
	}
	return nil
}

func (b *Base) MGet(ctx context.Context, key []string) ([]interface{}, error) {
	return BaseDao.Redis.MGet(ctx, key...).Result()
}

func (b *Base) MGetString(ctx context.Context, key []string) ([]string, error) {
	values, err := b.MGet(ctx, key)
	if err != nil {
		return nil, err
	}
	strList := make([]string, 0, len(values))
	for _, v := range values {
		val, ok := v.(string)
		if !ok {
			return nil, ErrStr
		}
		strList = append(strList, val)
	}
	return strList, nil
}
