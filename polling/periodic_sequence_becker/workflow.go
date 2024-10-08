package periodic_sequence_becker

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	TaskQueueName = "pollingPeriodicSequenceSampleQueueBecker"
)

func PeriodicSequencePolling(ctx workflow.Context, pollingInterval time.Duration) (string, error) {
	cwo := workflow.ChildWorkflowOptions{}
	ctx = workflow.WithChildOptions(ctx, cwo)
	params := ChildWorkflowParams{
		PollingInterval: pollingInterval,
		Index:           0,
	}
	res := workflow.ExecuteChildWorkflow(ctx, PollingChildWorkflow, params)
	var result string
	err := res.Get(ctx, &result)
	return result, err

}
