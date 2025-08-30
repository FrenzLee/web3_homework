package task_02

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/*
编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
*/

type Counter struct {
	mu    sync.Mutex
	count int
}

func (c *Counter) increment() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func Calculate() {
	counter := Counter{}
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				counter.increment()
			}
		}()
	}

	time.Sleep(time.Second)

	fmt.Println("总计：", counter.count)
}

/*
使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
*/

func Increment() {
	var counter int64

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	time.Sleep(time.Second)

	fmt.Println("总计：", counter)

}
