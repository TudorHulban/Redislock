package redlock

import (
	"errors"
)

type DTO struct {
	key   []byte
	value []byte
}

func (c *InRedisCache) SetTTL(dto *DTO) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, errSet := conn.Do("SETEX", dto.key, c.secondsTTL, dto.value)
	return errSet
}

func (c *InRedisCache) Get(key []byte) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()

	value, errGet := conn.Do("GET", key)
	if errGet != nil {
		return nil, errGet
	}

	if value == nil {
		return nil, errors.New("item not found")
	}

	var buf []byte
	buf = append(buf, value.([]uint8)...)

	return buf, nil
}
