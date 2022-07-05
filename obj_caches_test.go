package redlock

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redsync/redsync/v4"
	redsyncredigo "github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/stretchr/testify/require"
)

const _sock1 = "127.0.0.1:6379"
const _sock2 = "127.0.0.1:6378"

func TestRedsyncHappyPath(t *testing.T) {
	pool1, errNew1 := NewCache(context.Background(), _sock1, WithTTL(1))
	require.NoError(t, errNew1)

	rs1 := redsyncredigo.NewPool(&pool1.pool)

	pool2, errNew2 := NewCache(context.Background(), _sock2, WithTTL(1))
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

	require.Nil(t, caches.SetTTL(&dto))

	value, errGet := caches.Get(dto.key)
	t.Logf("errGet: %#v", errGet)
	require.Nil(t, errGet)
	require.NotNil(t, value)
	require.Equal(t, dto.value, value)

	if _, errUnlock := mutex.Unlock(); errUnlock != nil {
		panic(errUnlock)
	}
}

func TestRedsyncRace(t *testing.T) {
	pool1, errNew1 := NewCache(context.Background(), _sock1, WithTTL(5))
	require.NoError(t, errNew1)

	rs1 := redsyncredigo.NewPool(&pool1.pool)

	pool2, errNew2 := NewCache(context.Background(), _sock2, WithTTL(5))
	require.NoError(t, errNew2)

	rs2 := redsyncredigo.NewPool(&pool2.pool)

	rs := redsync.New(rs1, rs2)

	mutex := rs.NewMutex("test-redsync")

	key := []byte("xxx")

	dtoGoRoutine := DTO{
		key:   key,
		value: []byte("yyy1"),
	}

	dtoMutex := DTO{
		key:   key,
		value: []byte("yyy2"),
	}

	caches := NewCaches(pool1, pool2)

	chComm := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		<-chComm

		errSet := caches.SetTTL(&dtoGoRoutine)
		if len(errSet) != 0 {
			for _, err := range errSet {
				fmt.Println("caches.Set(&dto):", err)
			}
		}

		value, errGet := caches.Get(key)
		fmt.Println(time.Now().UnixMilli(), "errGet goroutine:", errGet, "value:", string(value))

		t.Log("mutex value: ", string(dtoMutex.value))
		t.Log("got value: ", string(value))
		require.Equal(t, dtoMutex.value, value)
	}()

	if errLock := mutex.Lock(); errLock != nil {
		panic(errLock)
	}

	errSet := caches.SetTTL(&dtoMutex)
	if len(errSet) != 0 {
		for _, err := range errSet {
			fmt.Println("caches.Set(&dto):", err)
		}
	}

	chComm <- struct{}{}

	value, errGet := caches.Get(key)
	fmt.Println(time.Now().UnixMilli(), "errGet in mutex:", errGet, "value:", string(value))

	if _, errUnlock := mutex.Unlock(); errUnlock != nil {
		panic(errUnlock)
	}

	wg.Wait()
}
