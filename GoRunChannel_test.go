package goRunInChannel

import (
	"log"
	"testing"
	"time"
)

func TestGoRunChannelInt(t *testing.T) {
	// 创建并发度为1的协程池 （串行）
	goRunChannel := NewGoRunChannel[int](1)
	// 定义任务函数
	runable := func(param int) {
		log.Printf("Worker %d started\n", param)
		time.Sleep(time.Second)
		log.Printf("Worker %d done\n", param)
	}
	// 循环10个任务
	for i := 0; i < 10; i++ {
		// 执行任务
		err := goRunChannel.Run(runable, i)
		if err != nil {
			t.Errorf("run error: %v", err)
		}
	}
	// 等待所有任务完成
	goRunChannel.WaitAndClose()
}

func TestGoRunChannelAny(t *testing.T) {
	// 创建并发度为4的协程池
	goRun := NewGoRunChannel[any](4)
	// 任务函数
	runable := func(param any) {
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
		err := goRun.Run(runable, param)
		if err != nil {
			t.Errorf("run error: %v", err)
		}
	}

	// 等待所有任务完成
	goRun.WaitAndClose()
}

type Test struct {
	id int
}
