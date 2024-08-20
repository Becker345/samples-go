package periodic_sequence_becker

import (
	"context"
	"fmt"
	"github.com/temporalio/samples-go/polling"
	"go.temporal.io/sdk/activity"
	"math/rand"
	"time"
)

type PollingActivities struct {
	TestService *polling.TestService
}

func (a *PollingActivities) DoDeferAct(ctx context.Context) error {
	sleepTime := time.Duration(rand.Intn(15))
	fmt.Printf("start DoDeferAct, sleep: %ds\n", sleepTime)
	//写一个本地文件，内容为1
	time.Sleep(sleepTime * time.Second)
	//帮写一个本地文件，内容为2
	fmt.Printf("end DoDeferAct: %s\n", time.Now())
	return nil
	//return a.TestService.GetServiceResult(ctx)
}

// DoPoll Activity.
func (a *PollingActivities) DoPoll(ctx context.Context) (string, error) {
	sleepTime := time.Duration(rand.Intn(15))
	fmt.Printf("start poll, sleep: %ds\n", sleepTime)
	//写一个本地文件，内容为1
	time.Sleep(sleepTime * time.Second)
	//帮写一个本地文件，内容为2
	fmt.Printf("end poll: %s\n", time.Now())
	return "result", nil
	//return a.TestService.GetServiceResult(ctx)
}

func (a *PollingActivities) DoPoll1(ctx context.Context) (string, error) {
	sleepTime := time.Duration(rand.Intn(35))
	fmt.Printf("start do poll1, sleep: %ds\n", sleepTime)
	for {
		select {
		// 随机等待一段时间，避免并发请求导致的API调用频率过高
		case <-time.After(sleepTime * time.Second):
			// Activity 记录心跳
			activity.RecordHeartbeat(ctx, "")
			fmt.Printf("end do poll1: %s\n", time.Now())
			return "sleep success", nil
		case <-ctx.Done():
			// 任务取消Event
			// 外部 context 被取消，停止 ticker 并退出循环
			fmt.Printf("context done: %s\n", time.Now())
			return "signal done", ctx.Err()
		}
	}
}
