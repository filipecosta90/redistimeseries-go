package bench

import (
	"fmt"
	redis_timeseries_go "github.com/RedisTimeSeries/redistimeseries-go"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/mediocregopher/radix/v3"
	"os"
	"runtime"
	. "testing"
	"time"
)

func getHost() (host string) {
	value, exists := os.LookupEnv("REDIS_TEST_HOST")
	host = "localhost:6379"
	if exists && value != "" {
		host = value
	}
	return
}

func getAuthPass() (authPass *string) {
	value, exists := os.LookupEnv("REDIS_TEST_AUTH")
	if exists && value != "" {
		authPass = &value
	}
	return
}

func init() {
	/* load test data */
	client := redis_timeseries_go.NewClient(getHost(), "bench_client", getAuthPass())
	conn := client.Pool.Get()
	conn.Do("FLUSHALL")
	duration, _ := time.ParseDuration("1h")
	for i := 0; i < 10; i++ {
		client.CreateKey(fmt.Sprintf("b%d", i), duration)
	}

}

func redigoTSADD(conn redigo.Conn, key, val string) error {
	_, err := redigo.Int64(conn.Do("TS.ADD", key, "*", val))
	return err
}

func redigoTSMADD(conn redigo.Conn, ternaries ...string) error {
	_, err := conn.Do("TS.MADD", redigo.Args{}.AddFlat(ternaries)...)
	return err
}

func radixTSADD(client radix.Client, key, val string) error {
	var barI int64
	if err := client.Do(radix.Cmd(&barI, "TS.ADD", key, "*", val)); err != nil {
		return err
	}
	return nil
}

func radixTSMADD(client radix.Client, ternaries ...string) error {
	if err := client.Do(radix.FlatCmd(nil, "TS.MADD", ternaries[0], ternaries[1:])); err != nil {
		return err
	}
	return nil
}

func newRedigo() redigo.Conn {
	c, err := redigo.Dial("tcp", getHost())
	if err != nil {
		panic(err)
	}
	return c
}

func BenchmarkParallelTsAdd(b *B) {

	parallel := runtime.GOMAXPROCS(0)

	// multiply parallel with GOMAXPROCS to get the actual number of goroutines and thus
	// connections needed for the benchmarks.
	poolSize := parallel * runtime.GOMAXPROCS(0)

	do := func(b *B, fn func() error) {
		b.ResetTimer()
		b.SetParallelism(parallel)
		b.RunParallel(func(pb *PB) {
			for pb.Next() {
				if err := fn(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}

	b.Run("radix", func(b *B) {
		mkRadixBench := func(opts ...radix.PoolOpt) func(b *B) {
			return func(b *B) {
				pool, err := radix.NewPool("tcp", getHost(), poolSize, opts...)
				if err != nil {
					b.Fatal(err)
				}
				defer pool.Close()

				// wait for the pool to fill up
				for {
					time.Sleep(50 * time.Millisecond)
					if pool.NumAvailConns() >= poolSize {
						break
					}
				}

				// avoid overhead of boxing the pool on each loop iteration
				client := radix.Client(pool)
				do(b, func() error {
					return radixTSADD(client, "b0", "1")
				})
			}
		}
		b.Run("no pipeline", mkRadixBench(radix.PoolPipelineWindow(0, 0)))
		b.Run(fmt.Sprintf("1 pipeline, 150us pipeline limit, no max CMDs flush limit"), mkRadixBench(radix.PoolPipelineConcurrency(1)))
		b.Run(fmt.Sprintf("%d pipelines, 10us pipeline limit, 100 max CMDs flush limit", poolSize), mkRadixBench(radix.PoolPipelineWindow(10*time.Microsecond, 100)))
		b.Run(fmt.Sprintf("default %d pipelines, 150us pipeline limit, no max CMDs flush limit", poolSize), mkRadixBench())
	})

	b.Run("redigo", func(b *B) {
		red := &redigo.Pool{MaxIdle: poolSize, Dial: func() (redigo.Conn, error) {
			return newRedigo(), nil
		}}
		defer red.Close()

		{ // make sure the pool is full
			var conns []redigo.Conn
			for red.MaxIdle > red.ActiveCount() {
				conns = append(conns, red.Get())
			}
			for _, conn := range conns {
				_ = conn.Close()
			}
		}

		do(b, func() error {
			conn := red.Get()
			if err := redigoTSADD(conn, "b0", "1"); err != nil {
				conn.Close()
				return err
			}
			return conn.Close()
		})

	})

}

func BenchmarkParallelTsMadd(b *B) {

	parallel := runtime.GOMAXPROCS(0)

	// multiply parallel with GOMAXPROCS to get the actual number of goroutines and thus
	// connections needed for the benchmarks.
	poolSize := parallel * runtime.GOMAXPROCS(0)

	do := func(b *B, fn func() error) {
		b.ResetTimer()
		b.SetParallelism(parallel)
		b.RunParallel(func(pb *PB) {
			for pb.Next() {
				if err := fn(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}

	b.Run("radix", func(b *B) {
		mkRadixBench := func(opts ...radix.PoolOpt) func(b *B) {
			return func(b *B) {
				pool, err := radix.NewPool("tcp", getHost(), poolSize, opts...)
				if err != nil {
					b.Fatal(err)
				}
				defer pool.Close()

				// wait for the pool to fill up
				for {
					time.Sleep(50 * time.Millisecond)
					if pool.NumAvailConns() >= poolSize {
						break
					}
				}

				// avoid overhead of boxing the pool on each loop iteration
				client := radix.Client(pool)
				do(b, func() error {
					return radixTSMADD(client, "b0", "*", "1", "b1", "*", "1", "b2", "*", "1", "b3", "*", "1", "b4", "*", "1", "b5", "*", "1", "b6", "*", "1", "b7", "*", "1", "b8", "*", "1", "b9", "*", "1")
				})
			}
		}

		//b.Run("no pipeline", mkRadixBench(radix.PoolPipelineWindow(0, 0)))
		//b.Run(fmt.Sprintf("1 pipeline, 150us pipeline limit, no max CMDs flush limit"), mkRadixBench(radix.PoolPipelineConcurrency(1)))
		//b.Run(fmt.Sprintf("%d pipelines, 10us pipeline limit, 100 max CMDs flush limit", poolSize), mkRadixBench(radix.PoolPipelineWindow(10*time.Microsecond, 100)))
		b.Run(fmt.Sprintf("default %d pipelines, 150us pipeline limit, no max CMDs flush limit", poolSize), mkRadixBench())
		//b.Run(fmt.Sprintf("%d pipelines, 250us pipeline limit, no max CMDs flush limit", poolSize), mkRadixBench(radix.PoolPipelineWindow(250*time.Microsecond, 0)))

	})

	b.Run("redigo", func(b *B) {
		red := &redigo.Pool{MaxIdle: poolSize, Dial: func() (redigo.Conn, error) {
			return newRedigo(), nil
		}}
		defer red.Close()

		{ // make sure the pool is full
			var conns []redigo.Conn
			for red.MaxIdle > red.ActiveCount() {
				conns = append(conns, red.Get())
			}
			for _, conn := range conns {
				_ = conn.Close()
			}
		}

		do(b, func() error {
			conn := red.Get()
			if err := redigoTSMADD(conn, "b0", "*", "1", "b1", "*", "1", "b2", "*", "1", "b3", "*", "1", "b4", "*", "1", "b5", "*", "1", "b6", "*", "1", "b7", "*", "1", "b8", "*", "1", "b9", "*", "1"); err != nil {
				conn.Close()
				return err
			}
			return conn.Close()
		})

	})

}
