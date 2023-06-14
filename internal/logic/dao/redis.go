package dao

import (
	"fmt"
	"golang.org/x/net/context"
	"time"
)

const (
	_prefixMidServer    = "mid_%d"
	_prefixKeyServer    = "key_%s"
	_prefixServerOnline = "ol_%s"
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
	return b.redis.Expire(ctx, key, expiration).Err()
}

func (b *Base) SetExpire(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return b.redis.Set(ctx, key, value, expiration).Err()
}

func (b *Base) HSet(ctx context.Context, key string, values ...interface{}) error {
	return b.redis.HSet(ctx, key, values).Err()
}

func (b *Base) HSetExpire(ctx context.Context, expire time.Duration, key string, values ...interface{}) error {
	if err := b.HSet(ctx, key, values); err != nil {
		return err
	}
	if err := b.Expire(ctx, key, expire); err != nil {
		return err
	}
	return nil
}
