package periodic_sequence_becker

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"

	"go.temporal.io/sdk/workflow"
)

type ChildWorkflowParams struct {
	PollingInterval time.Duration
	Index           int
}

func PollingChildWorkflow(ctx workflow.Context, params ChildWorkflowParams) (string, error) {
	defer func() {
		// 只第一次执行
		if params.Index > 1 {
			return
		}
		fmt.Printf("defer done")
		disCtx, _ := workflow.NewDisconnectedContext(ctx)
		ao1 := workflow.ActivityOptions{
			StartToCloseTimeout: 4 * time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 1,
			},
		}
		disActCtx := workflow.WithActivityOptions(disCtx, ao1)
		var b *PollingActivities
		err := workflow.ExecuteActivity(disActCtx, b.DoDeferAct).Get(disActCtx, nil)
		if err != nil {
			fmt.Println(" DoDeferAct err", "err: ", err.Error())
			return
		}
	}()
	if params.Index >= 1 {
		params.Index = params.Index + 1
		return "ok", nil
	}
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting child workflow with params", params)
	//actCtx := workflow.WithActivityOptions(ctx, ao)

	for i := 0; i < 10; i++ {
		// Here we would invoke a sequence of activities
		// For sample we just use a single one repeated several times
		workflow.Go(ctx, func(ctxInner workflow.Context) {
			// Do something here
			var pollResult string
			ao := workflow.ActivityOptions{
				StartToCloseTimeout: 4 * time.Minute,
				RetryPolicy: &temporal.RetryPolicy{
					MaximumAttempts: 1,
				},
			}
			var a *PollingActivities
			loggerInner := workflow.GetLogger(ctxInner)
			actCtx := workflow.WithActivityOptions(ctxInner, ao)
			err := workflow.ExecuteActivity(actCtx, a.DoPoll1).Get(actCtx, &pollResult)
			if err != nil {
				loggerInner.Info("ExecuteActivity err", "err: ", err.Error(), "result: ", pollResult)
				return
			}
			loggerInner.Info("ExecuteActivity done", "result: ", pollResult)
		})
	}
	params.Index = params.Index + 1

	err := workflow.Sleep(ctx, 25*time.Second)
	if err != nil {
		return "sleep err", err
	}

	logger.Info("before continue as new: %s", time.Now().String())
	err = workflow.NewContinueAsNewError(ctx, PollingChildWorkflow, params)
	logger.Info("after continue as new: %s", time.Now().String())
	return "", err
}
