package goRunInChannel

import (
	"log"
	"testing"
	"time"
)

func TestGoRunChannel(t *testing.T) {
	// 创建并发度为4的协程池
	goRunChannel := NewGoRunChannel(4)
	// 定义任务函数
	runable := func(param interface{}) {
		log.Printf("Worker %d started\n", param.(Test).id)
		time.Sleep(time.Second)
		log.Printf("Worker %d done\n", param.(Test).id)
	}
	// 循环10个任务
	for i := 0; i < 10; i++ {
		// 任务参数
		param := Test{
			id: i,
		}
		// 执行任务
		goRunChannel.Run(runable, param)
	}
	// 等待所有任务完成
	goRunChannel.Wait()
}

type Test struct {
	id int
}
