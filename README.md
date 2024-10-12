# GoRunChannel

使用 channel 控制 goroutine 并发度

```go
package main

import (
	"log"
	"time"

	goRunInChannel "github.com/wooenrico/go-run-in-channel"
)

func main() {

	// 创建并发度为4的协程池
	goRun := goRunInChannel.NewGoRunChannel(4)
	// 任务函数
	Runable := func(Interface interface{}) {
		log.Printf("Worker %d started\n", Interface.(Test).id)
		time.Sleep(time.Second)
		log.Printf("Worker %d done\n", Interface.(Test).id)
	}

	// 循环10个任务
	for i := 0; i < 10; i++ {
		// 任务参数
		param := Test{
			id: i,
		}
		// 执行任务
		goRun.Run(Runable, param)
	}

	// 等待所有任务完成
	goRun.Wait()
}

// 自定义任务参数
type Test struct {
	id int
}
```

展示

```log
2024/05/22 14:12:32 Worker 1 started
2024/05/22 14:12:32 Worker 3 started
2024/05/22 14:12:32 Worker 2 started
2024/05/22 14:12:32 Worker 0 started
2024/05/22 14:12:33 Worker 3 done
2024/05/22 14:12:33 Worker 2 done
2024/05/22 14:12:33 Worker 0 done
2024/05/22 14:12:33 Worker 1 done
2024/05/22 14:12:33 Worker 4 started
2024/05/22 14:12:33 Worker 5 started
2024/05/22 14:12:33 Worker 6 started
2024/05/22 14:12:33 Worker 7 started
2024/05/22 14:12:34 Worker 7 done
2024/05/22 14:12:34 Worker 5 done
2024/05/22 14:12:34 Worker 8 started
2024/05/22 14:12:34 Worker 9 started
2024/05/22 14:12:34 Worker 4 done
2024/05/22 14:12:34 Worker 6 done
2024/05/22 14:12:35 Worker 8 done
2024/05/22 14:12:35 Worker 9 done
```