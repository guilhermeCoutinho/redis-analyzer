package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	redis "github.com/go-redis/redis/v8"
)

func randomKey(n int) string {
	s := randomS()
	for i := 0; i < n; i++ {
		s += ":" + randomS()
	}
	return s
}

func insertRandomKeysRoutine(rdb *redis.Client, size int) {
	wg := sync.WaitGroup{}
	routines := 100
	wg.Add(routines)
	fmt.Println("Inserting random keys into redis.")
	for i := 0; i < routines; i++ {
		go func() {
			defer wg.Done()
			insertRandomKeys(context.Background(), rdb, size/100)
		}()
	}
	wg.Wait()
}

func insertRandomKeys(ctx context.Context, rClient *redis.Client, n int) {
	for i := 0; i < n; i++ {
		_, err := rClient.Set(ctx, randomKey(rand.Intn(3)), randomS(), 0).Result()
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func randomS() string {
	s := []string{
		"wren",
		"soggy",
		"unique",
		"unaccountable",
		"lamp",
		"crash",
		"answer",
		"unsightly",
		"bone",
		"sticky",
		"tasty",
		"turkey",
		"lethal",
		"learned",
		"assorted",
		"frantic",
		"aberrant",
		"lewd",
		"brush",
		"smoke",
		"instrument",
		"bead",
		"knee",
		"bulb",
		"son",
		"haunt",
		"calendar",
		"upset",
		"homeless",
		"time",
		"mean",
		"ludicrous",
		"airport",
		"servant",
		"divergent",
		"allow",
		"guide",
		"brake",
		"toothsome",
		"extra-large",
		"file",
		"wistful",
		"quick",
		"language",
		"itch",
		"actually",
		"concerned",
		"understood",
		"apparatus",
		"befitting",
		"parched",
		"test",
		"fuel",
		"agree",
		"cooing",
		"devilish",
		"stroke",
		"heavy",
		"repair",
		"annoying",
		"approve",
		"adhesive",
		"vegetable",
		"reject",
		"command",
		"invite",
		"honorable",
		"towering",
		"wasteful",
		"walk",
		"guarded",
		"clever",
		"scissors",
		"tiresome",
		"bright",
		"doll",
		"juice",
		"statement",
		"cool",
		"confused",
		"dreary",
		"scare",
		"psychotic",
		"fork",
		"nimble",
		"consist",
		"unused",
		"loving",
		"muddle",
		"wobble",
		"quilt",
		"superficial",
		"shirt",
		"trail",
		"minister",
		"thin",
		"scream",
		"crack",
		"sloppy",
		"nice",
	}
	n := int(rand.ExpFloat64() / .1)
	if n >= len(s) {
		n = len(s) - 1
	}
	return s[n]
}
