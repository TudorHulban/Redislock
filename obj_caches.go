package redlock

import "errors"

type Caches []*InRedisCache

func NewCaches(cache ...*InRedisCache) *Caches {
	var res Caches

	res = append(res, cache...)

	return &res
}

func (c *Caches) SetTTL(dto *DTO) []error {
	var errs []error

	for _, cache := range *c {
		errSet := cache.SetTTL(dto)
		if errSet != nil {
			errs = append(errs, errSet)
		}
	}

	return errs
}

func (c *Caches) Get(key []byte) ([]byte, []error) {
	var errs []error
	var value []byte

	for _, cache := range *c {
		var errGet error

		value, errGet = cache.Get(key)
		if errGet != nil {
			errs = append(errs, errGet)
		}

		if len(value) == 0 {
			errs = append(errs, errors.New("item not found"))
		}
	}

	return value, errs
}
