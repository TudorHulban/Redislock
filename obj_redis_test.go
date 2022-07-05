package redlock

import (
	"context"
	"testing"

	"github.com/go-redsync/redsync/v4"
	redsyncredigo "github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/stretchr/testify/require"
)

const _sock1 = "127.0.0.1:6379"
const _sock2 = "127.0.0.1:6378"

func TestRedsync(t *testing.T) {
	pool1, errNew1 := NewCache(context.Background(), _sock1, WithTTL(5))
	require.NoError(t, errNew1)

	rs1 := redsyncredigo.NewPool(&pool1.pool)

	pool2, errNew2 := NewCache(context.Background(), _sock2, WithTTL(5))
	require.NoError(t, errNew2)

	rs2 := redsyncredigo.NewPool(&pool2.pool)

	rs := redsync.New(rs1, rs2)

	mutex := rs.NewMutex("test-redsync")

	dto := DTO{
		key:   []byte("xxx"),
		value: []byte("yyy"),
	}

	caches := NewCaches(pool1, pool2)

	if errLock := mutex.Lock(); errLock != nil {
		panic(errLock)
	}

	require.Nil(t, caches.Set(&dto))

	if _, errUnlock := mutex.Unlock(); errUnlock != nil {
		panic(errUnlock)
	}
}
