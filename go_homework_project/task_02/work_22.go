package task_02

import (
	"fmt"
	"sync"
	"time"
)

/*
编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
*/

func PrintNum(wg *sync.WaitGroup, mu *sync.Mutex) {

	go func(wg *sync.WaitGroup, mu *sync.Mutex) {
		defer wg.Done()
		for i := 1; i < 11; i++ {
			if i%2 == 1 {
				mu.Lock()
				fmt.Println("奇数：", i)
				mu.Unlock()
			}
		}
	}(wg, mu)

	go func(wg *sync.WaitGroup, mu *sync.Mutex) {
		defer wg.Done()
		for i := 1; i < 11; i++ {
			if i%2 == 0 {
				mu.Lock()
				fmt.Println("偶数：", i)
				mu.Unlock()
			}
		}
	}(wg, mu)

}

/*
设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
*/

// 定义 任务 类型
type Task func()

// 定义调度器
type Scheduler struct {
	tasks []Task
}

// 创建调度器
func NewScheduler(tasks ...Task) *Scheduler {
	newScheduler := Scheduler{tasks: tasks}
	return &newScheduler
}

// 添加任务到调度器
func (s *Scheduler) Add(task Task) {
	s.tasks = append(s.tasks, task)
}

// 执行所有任务
func (s *Scheduler) Run(wg *sync.WaitGroup) {

	for _, task := range s.tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			start := time.Now() //开始时间

			t() //执行任务

			countTime := time.Since(start) //执行时间
			fmt.Printf("任务执行时间: %s\n", countTime)
		}(task)
	}

}
