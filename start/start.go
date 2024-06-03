package start

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"hello-world-temporal/app/lib"
	"hello-world-temporal/app/workflow"
)

func Consume() {
	// 模拟Cron逻辑，每20秒启动一个主工作流
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			consume(fmt.Sprintf("consume-cron-workflow-%d", time.Now().Unix()))
		}
	}
}

func Produce() {
	// 模拟Cron逻辑，每10秒启动一个主工作流
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			produce(fmt.Sprintf("produce-cron-workflow-%d", time.Now().Unix()))
		}
	}
}

func produce(workflowID string) {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: lib.ProduceTaskQueue,
	}

	// Start the Workflow
	we, err := c.ExecuteWorkflow(context.Background(), options, workflow.MainProduceWorkflow)
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

func consume(workflowID string) {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: lib.ConsumeTaskQueue,
	}

	// Start the Workflow
	we, err := c.ExecuteWorkflow(context.Background(), options, workflow.MainConsumeWorkflow)
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
