package main

import (
	"context"
	"github.com/temporalio/samples-go/polling/periodic_sequence_becker"
	"log"
	"time"

	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "pollingSampleQueue_" + time.Now().String(),
		TaskQueue: periodic_sequence_becker.TaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, periodic_sequence_becker.PeriodicSequencePolling, 1*time.Second)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
