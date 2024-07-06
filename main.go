package main

import (
	"conc/patterns"
	"context"
	"fmt"
	"sync"
)

func main() {

	// ----------------------------------FanOut -> FanIn--------------------

	sourceChan := make(chan int)
	wg := &sync.WaitGroup{}
	go func() {
		defer close(sourceChan)

		for i := 0; i < 100; i++ {
			sourceChan <- i
		}
	}()

	chans := patterns.Fanout(sourceChan, 20)
	outChan := patterns.Fanin(chans)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for val := range outChan {
			fmt.Printf("output value :%d\n", val)
		}
	}()

	wg.Wait()

	// ----------------------------------------Future----------------------

	ctx, _ := context.WithCancel(context.Background())
	future := patterns.BlockingFunction(ctx)
	// cancel()
	res, err := future.Result()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	fmt.Println(res)

	// ----------------------------------------Sharded Map-----------------

	shardedMap := patterns.NewShardedMap[int](5)

	shardedMap.Put("argentina", 132)
	shardedMap.Put("newzeland", 22)
	shardedMap.Put("poland", 3982)
	shardedMap.Put("india", 91)
	shardedMap.Put("russia", 7)

	fmt.Println(shardedMap.Get("argentina"))
	fmt.Println(shardedMap.Get("newzeland"))
	fmt.Println(shardedMap.Get("poland"))
	fmt.Println(shardedMap.Get("india"))
	fmt.Println(shardedMap.Get("russia"))

	fmt.Println("Keys: ", shardedMap.Keys())

}
