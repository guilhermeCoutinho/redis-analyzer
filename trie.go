package main

import (
	"fmt"

	"github.com/guilhermeCoutinho/concurrent-generic-heap/heap"
)

type trie struct {
	subKeys  int
	mem      int64
	prefix   string
	children map[string]*trie
}

func NewTrie() *trie {
	return &trie{children: make(map[string]*trie)}
}

func (t *trie) Print(maxD int, totalMem int64) {
	t.printRecursively("", "", 0, maxD, totalMem)
}

func (t *trie) printRecursively(prefix, childrenPrefix string, currD, maxD int, totalMem int64) {
	if t == nil || currD > maxD {
		return
	}

	pattern := ""
	if t.subKeys > 1 {
		pattern = "*"
	}

	memPercent := fmt.Sprintf("%d%%", int64(100*float64(t.mem)/float64(totalMem)))
	fmt.Println(prefix, string(t.prefix)+pattern, ByteCountSI(t.mem), memPercent)

	for _, v := range t.children {
		v.printRecursively(childrenPrefix+"|─── ", childrenPrefix+"|    ", currD+1, maxD, totalMem)
	}
}

func (t *trie) Add(prefixes []string, mem int64) {
	if len(prefixes) == 0 {
		return
	}

	t.mem += mem
	t.subKeys += 1
	currentPrefix := prefixes[0]

	if t.children[currentPrefix] == nil {
		t.children[currentPrefix] = &trie{
			children: make(map[string]*trie),
			prefix:   currentPrefix,
		}
	}

	if len(prefixes) == 1 {
		t.children[currentPrefix].mem += mem
		t.children[currentPrefix].subKeys += 1
	}
	t.children[currentPrefix].Add(prefixes[1:], mem)
}

type trieInfo struct {
	t   *trie
	key string
}

func (t *trie) Trim(percentile float64) *trie {
	minHeap := buildMinHeap(t)
	keysToTrim := getKeysToTrim(t, minHeap, percentile)
	for _, key := range keysToTrim {
		t.children[key] = nil
	}
	recurseTrim(t, percentile)
	return t
}

func recurseTrim(t *trie, percentile float64) {
	for _, v := range t.children {
		if v != nil {
			v.Trim(percentile)
		}
	}
}

func buildMinHeap(t *trie) heap.Heap {
	minHeap := heap.NewHeap(100, func(a, b interface{}) bool {
		return a.(*trieInfo).t.mem < b.(*trieInfo).t.mem
	})

	for key, child := range t.children {
		minHeap.Push(&trieInfo{
			t:   child,
			key: key,
		})
	}
	return minHeap
}

func getKeysToTrim(t *trie, minHeap heap.Heap, percentile float64) []string {
	ansMemory := 0
	keysToTrim := []string{}
	percentileMemory := int((1 - percentile) * float64(t.mem))

	for ansMemory < percentileMemory {
		child, err := minHeap.Pop()
		if err != nil {
			break
		}
		ansMemory += int(child.(*trieInfo).t.mem)
		keysToTrim = append(keysToTrim, child.(*trieInfo).key)
	}
	return keysToTrim
}

func (t *trie) TrimLargestKeys(count int) *trie {
	maxHeap := buildMaxHeap(t)
	keysToTrim := getKeysToTrimLargestKeys(t, maxHeap, count)
	for _, key := range keysToTrim {
		t.children[key] = nil
	}
	recurseTrimLargestKeys(t, count)
	return t
}

func recurseTrimLargestKeys(t *trie, count int) {
	for _, v := range t.children {
		if v != nil {
			v.TrimLargestKeys(count)
		}
	}
}

func buildMaxHeap(t *trie) heap.Heap {
	minHeap := heap.NewHeap(100, func(a, b interface{}) bool {
		return a.(*trieInfo).t.mem > b.(*trieInfo).t.mem
	})

	for key, child := range t.children {
		minHeap.Push(&trieInfo{
			t:   child,
			key: key,
		})
	}
	return minHeap
}

func getKeysToTrimLargestKeys(t *trie, maxHeap heap.Heap, count int) []string {
	keysToKeep := map[string]interface{}{}
	for i := 0; i < count; i++ {
		tInfo, err := maxHeap.Pop()
		if err != nil {
			continue
		}
		keysToKeep[tInfo.(*trieInfo).key] = interface{}(nil)
	}

	keysToTrim := []string{}
	for k := range t.children {
		if _, ok := keysToKeep[k]; !ok {
			keysToTrim = append(keysToTrim, k)
		}
	}
	return keysToTrim
}
