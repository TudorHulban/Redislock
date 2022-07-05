package redlock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOps(t *testing.T) {
	pool, errNew := NewCache(context.Background(), _sock1, WithTTL(5))
	require.NoError(t, errNew)

	dto := DTO{
		key:   []byte("xxx"),
		value: []byte("yyy"),
	}

	require.NoError(t, pool.Set(&dto))

	value, errGet := pool.Get(dto.key)
	require.NoError(t, errGet)
	require.Equal(t, dto.value, value)
}
