[![license](https://img.shields.io/github/license/RedisTimeSeries/JRedisTimeSeries.svg)](https://github.com/RedisTimeSeries/JRedisTimeSeries)
[![CircleCI](https://circleci.com/gh/RedisTimeSeries/redistimeseries-go.svg?style=svg&circle-token=022ed6c86563cbb7d19ff4fd3ca6eab9053603f2)](https://circleci.com/gh/RedisTimeSeries/redistimeseries-go)
[![GitHub issues](https://img.shields.io/github/release/RedisTimeSeries/redistimeseries-go.svg)](https://github.com/RedisTimeSeries/redistimeseries-go/releases/latest)
[![Codecov](https://codecov.io/gh/RedisTimeSeries/redistimeseries-go/branch/master/graph/badge.svg)](https://codecov.io/gh/RedisTimeSeries/redistimeseries-go)
[![GoDoc](https://godoc.org/github.com/RedisTimeSeries/redistimeseries-go?status.svg)](https://godoc.org/github.com/RedisTimeSeries/redistimeseries-go)


# redis-timeseries-go

Go client for RedisTimeSeries (https://github.com/RedisLabsModules/redis-timeseries)

Client and ConnPool based on the work of dvirsky and mnunberg on https://github.com/RediSearch/redisearch-go


## Running tests

A simple test suite is provided, and can be run with:

```sh
$ go test
```

The tests expect a Redis server with the RedisTimeSeries module loaded to be available at localhost:6379

## License

redistimeseries-go is distributed under the Apache-2 license - see [LICENSE](LICENSE)
