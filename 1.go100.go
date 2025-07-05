package main

import (
	"fmt"
	"sync"
)

func main() {
	// 定义三个通道用于协程间通信
	a := make(chan struct{}) // 这里为什么用空结构体？因为在go中空结构体是不占用内存的
	b := make(chan struct{})
	c := make(chan struct{})

	// 使用WaitGroup确保主协程等待子协程完成
	var wg sync.WaitGroup
	wg.Add(3)

	// 打印1-100
	current := 1
	maxN := 100

	// 协程A：从通道C接收信号，打印数字，然后通知B
	go func() {
		defer wg.Done()
		for current <= maxN {
			<-c // 等待C的信号
			if current > maxN {
				close(b) // 通知协程B退出
				return
			}
			fmt.Printf("协程A: %d\n", current)
			current++
			b <- struct{}{} // 通知B
		}
	}()

	// 协程B：从通道A接收信号，打印数字，然后通知C
	go func() {
		defer wg.Done()
		for current <= maxN {
			<-b // 等待A的信号
			if current > maxN {
				close(a) // 通知协程C退出
				return
			}
			fmt.Printf("协程B: %d\n", current)
			current++
			a <- struct{}{} // 通知C
		}
	}()

	// 协程C：从通道B接收信号，打印数字，然后通知A
	go func() {
		defer wg.Done()
		for current <= maxN {
			<-a // 等待B的信号
			if current > maxN {
				close(c) // 通知协程A退出
				return
			}
			fmt.Printf("协程C: %d\n", current)
			current++
			c <- struct{}{} // 通知A
		}
	}()

	// 初始启动信号
	c <- struct{}{} // 通知协程A开始

	wg.Wait()
	fmt.Println("所有数字打印完成")
}
