package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"hello-world-temporal/app"
)

func main() {

	// 模拟Cron逻辑，每20秒启动一个主工作流
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			start(fmt.Sprintf("cron-workflow-%d", time.Now().Unix()))
		}
	}
}

func start(workflowID string) {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: app.ConsumeTaskQueue,
	}

	// Start the Workflow
	we, err := c.ExecuteWorkflow(context.Background(), options, app.MainConsumeWorkflow)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	printResults(result, we.GetID(), we.GetRunID())
}

func printResults(result string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s %s\n\n", time.Now(), result)
}
