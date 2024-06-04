package worker

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"hello-world-temporal/app/activity"
	"hello-world-temporal/app/lib"
	"hello-world-temporal/app/workflow"
)

func Produce() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{HostPort: "", Namespace: "default"})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, lib.ProduceTaskQueue, worker.Options{})
	w.RegisterWorkflow(workflow.MainProduceWorkflow)
	w.RegisterWorkflow(workflow.SubProduceWorkflow)
	w.RegisterActivity(activity.Produce)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func Consume() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{HostPort: "", Namespace: "default"})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, lib.ConsumeTaskQueue, worker.Options{})
	w.RegisterWorkflow(workflow.MainConsumeWorkflow)
	w.RegisterWorkflow(workflow.SubConsumeWorkflow)
	w.RegisterActivity(activity.Consume)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
