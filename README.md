# Redislock
Tests with redsync.

## Infrastructure
```sh
docker run --name redis-test1 -p 6378:6379  -d redis 
docker run --name redis-test2 -p 6379:6379  -d redis 
```

## Troubleshooting
Connect to container:
```sh
docker exec -it <container-ID> bash
```
Get all keys:
```sh
redis-cli KEYS '*'
```
Delete all keys:
```sh
redis-cli FLUSHALL / FLUSHDB
```
or per database:
```sh
 redis-cli -n <database_number> FLUSHDB
```

## Resources
```html
https://github.com/go-redsync/redsync/blob/master/examples/redigo/main.go
```
