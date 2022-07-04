# Redislock
Tests with redsync.

## Infrastructure
```sh
docker run --name redis-test1 -p 6378:6379  -d redis 
docker run --name redis-test2 -p 6379:6379  -d redis 
```

## Resources
```html
https://github.com/go-redsync/redsync/blob/master/examples/redigo/main.go
```
