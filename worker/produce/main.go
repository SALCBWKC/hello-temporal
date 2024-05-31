package main

import (
	"hello-world-temporal/app"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{HostPort: "", Namespace: "default"})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, app.ProduceTaskQueue, worker.Options{})
	w.RegisterWorkflow(app.MainProduceWorkflow)
	//w.RegisterWorkflow(app.ProduceWorkflow)
	w.RegisterActivity(app.ProduceActivity)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
