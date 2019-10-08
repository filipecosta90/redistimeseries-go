## Running the benchmarks
The tests expect a Redis server with the RedisTimeSeries module loaded to be available at localhost:6379
Optionally, you can alter that by setting the `REDIS_TEST_HOST` env variable.


The simple benchmark suite is provided, and can be run with:

```sh
$ go test -v -run=XXX -bench=. -benchtime=30s
```

### TS.ADD 
```
$ go test -v -run=XXX -bench=BenchmarkParallelTsAdd -benchtime=30s
goos: darwin
goarch: amd64
pkg: github.com/RedisTimeSeries/redistimeseries-go/bench
BenchmarkParallelTsAdd/radix/no_pipeline-12                                                                      2945659             12645 ns/op
BenchmarkParallelTsAdd/radix/1_pipeline,_150us_pipeline_limit,_no_max_CMDs_flush_limit-12                        8754523              4107 ns/op
BenchmarkParallelTsAdd/radix/144_pipelines,_10us_pipeline_limit,_100_max_CMDs_flush_limit-12                    16647993              2547 ns/op
BenchmarkParallelTsAdd/radix/default_-_144_pipelines,_150us_pipeline_limit,_no_max_CMDs_flush_limit-12          10329037              3201 ns/op
BenchmarkParallelTsAdd/redigo-12                                                                                 2568367             14669 ns/op
PASS
ok      github.com/RedisTimeSeries/redistimeseries-go/bench     225.206s
```



### TS.MADD 
```
$ go test -v -run=XXX -bench=BenchmarkParallelTsMadd -benchtime=60s
goos: darwin
goarch: amd64
pkg: github.com/RedisTimeSeries/redistimeseries-go/bench
BenchmarkParallelTsMadd/radix/default_144_pipelines,_150us_pipeline_limit,_no_max_CMDs_flush_limit-12            5914801             13107 ns/op
BenchmarkParallelTsMadd/redigo-12                                                                                3465058             20573 ns/op
PASS
ok      github.com/RedisTimeSeries/redistimeseries-go/bench     183.085s
```