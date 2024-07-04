package goRunInChannel

import (
	"sync"
)

type GoRunChannel struct {
	Channel   chan bool
	WaitGroup *sync.WaitGroup
}

func NewGoRunChannel(parallel int) *GoRunChannel {
	return &GoRunChannel{
		Channel:   make(chan bool, parallel),
		WaitGroup: &sync.WaitGroup{},
	}
}

func (GoRunChannel *GoRunChannel) Run(runable func(param interface{}), param interface{}) {
	// 添加任务
	GoRunChannel.WaitGroup.Add(1)
	// 获取信号量
	GoRunChannel.Channel <- true
	// 执行任务
	go running(GoRunChannel, runable, param)
}

func running(GoRunChannel *GoRunChannel, runable func(param interface{}), param interface{}) {
	defer GoRunChannel.WaitGroup.Done()
	runable(param)
	// 释放信号量
	<-GoRunChannel.Channel
}
