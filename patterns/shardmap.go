package patterns

import (
	"crypto/sha1"
	"sync"
)

type Shard[T any] struct {
	sync.RWMutex
	m map[string]T
}

type ShardedMap[T any] []*Shard[T]

func NewShardedMap[T any](nShards int) ShardedMap[T] {
	shards := make([]*Shard[T], nShards)

	for i := range nShards {
		shardMap := make(map[string]T)
		shards[i] = &Shard[T]{m: shardMap}
	}

	return shards
}

func (sm ShardedMap[T]) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	n := int(checksum[15])
	return n % len(sm)
}

func (sm ShardedMap[T]) getShard(key string) *Shard[T] {
	idx := sm.getShardIndex(key)
	return sm[idx]
}

func (sm ShardedMap[T]) Get(key string) T {
	shard := sm.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

func (sm ShardedMap[T]) Put(key string, value T) {
	shard := sm.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = value
}

func (sm ShardedMap[T]) Keys() []string {
	keys := make([]string, 0)
	for _, shard := range sm {
		for k := range shard.m {
			keys = append(keys, k)
		}
	}
	return keys
}
