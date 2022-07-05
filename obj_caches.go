package redlock

import "errors"

type Caches []*InRedisCache

func NewCaches(cache ...*InRedisCache) *Caches {
	var res Caches

	res = append(res, cache...)

	return &res
}

func (c *Caches) Set(dto *DTO) []error {
	var errs []error

	for _, cache := range *c {
		conn := cache.pool.Get()

		_, errSet := conn.Do("SET", dto.key, dto.value)
		if errSet != nil {
			errs = append(errs, errSet)
		}

		conn.Close()
	}

	return errs
}

func (c *Caches) Get(key []byte) ([]byte, []error) {
	var errs []error
	var value interface{}

	for _, cache := range *c {
		conn := cache.pool.Get()

		var errGet error

		value, errGet = conn.Do("GET", key)
		if errGet != nil {
			errs = append(errs, errGet)
		}

		if value == nil {
			errs = append(errs, errors.New("item not found"))
		}

		conn.Close()
	}

	if value != nil {
		var buf []byte
		buf = append(buf, value.([]uint8)...)

		return buf, errs
	}

	return nil, errs
}
