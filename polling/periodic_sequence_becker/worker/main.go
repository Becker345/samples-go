package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/polling"
	"github.com/temporalio/samples-go/polling/periodic_sequence_becker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, periodic_sequence_becker.TaskQueueName, worker.Options{})

	w.RegisterWorkflow(periodic_sequence_becker.PeriodicSequencePolling)
	w.RegisterWorkflow(periodic_sequence_becker.PollingChildWorkflow)
	testService := polling.NewTestService(50)
	activities := &periodic_sequence_becker.PollingActivities{
		TestService: &testService,
	}
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
