package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	redis "github.com/go-redis/redis/v8"
)

const (
	enableProfiler   = false
	profilerFileName = "cpu_profile"
)

func main() {
	if enableProfiler {
		f, err := os.Create(profilerFileName)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:9001",
		DB:       0,
		Password: "",
	})

	mode := os.Getenv("MODE")
	fmt.Println(mode)
	if mode == "POPULATEREDIS" {
		insertRandomKeysRoutine(rdb, 5000)
	} else {
		t := scanRedis(context.Background(), rdb)
		t.Trim(.2)
		//t.TrimLargestKeys(5)
		t.Print(10, t.mem)
	}
}

func scanRedis(ctx context.Context, rClient *redis.Client) *trie {
	t := NewTrie()
	cursor := uint64(0)
	var keys []string
	var err error
	totalKeys := 0

	for {
		keys, cursor, err = rClient.Scan(ctx, cursor, "*", 10000).Result()
		totalKeys += len(keys)
		fmt.Printf("\rScanning keys. Total memory so far:%s", ByteCountSI(t.mem))

		if err != nil {
			fmt.Println(err.Error())
			return t
		}

		addToTrie(ctx, rClient, t, keys)
		if cursor == 0 {
			fmt.Println()
			return t
		}
	}
}

func addToTrie(ctx context.Context, rClient *redis.Client, t *trie, keys []string) {
	pipe := rClient.Pipeline()
	m := map[string]*redis.IntCmd{}
	for _, key := range keys {
		m[key] = pipe.MemoryUsage(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	for key, memUsage := range m {
		prefixes := strings.Split(key, ":")
		mem, _ := memUsage.Result()
		t.Add(prefixes, mem)
	}

}
