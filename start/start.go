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

const (
	// mock
	produceDuration = 10 * time.Second
	consumeDuration = 20 * time.Second
)

var (
	consumeChan = make(chan string, 1)
	produceChan = make(chan string, 1)
)

func Consume() {
	// 模拟Cron逻辑，每2分钟启动一个主工作流
	ticker := time.NewTicker(consumeDuration)
	for {
		select {
		case <-ticker.C:
			consume(fmt.Sprintf("consume-cron-workflow-%d", time.Now().Unix()))
		}
	}
}

func Produce() {
	// 模拟Cron逻辑，每1分钟启动一个主工作流
	ticker := time.NewTicker(produceDuration)
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
	ctx := context.Background()
	terminateWorkflow(ctx, c, workflowID, produceChan)
	we, err := c.ExecuteWorkflow(ctx, options, workflow.MainProduceWorkflow)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}
	produceChan <- workflowID

	// Get the results
	var result string
	err = we.Get(ctx, &result)
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
	ctx := context.Background()
	terminateWorkflow(ctx, c, workflowID, consumeChan)
	we, err := c.ExecuteWorkflow(ctx, options, workflow.MainConsumeWorkflow)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}
	consumeChan <- workflowID

	// Get the results
	var result string
	err = we.Get(ctx, &result)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	printResults(result, we.GetID(), we.GetRunID())
}

func terminateWorkflow(ctx context.Context, c client.Client, id string, ch <-chan string) {
	select {
	case workflowID, ok := <-ch:
		if ok {
			desc, err := c.DescribeWorkflowExecution(ctx, workflowID, "")
			if err != nil {
				log.Println("describe failed", err)
				return
			}
			if desc.WorkflowExecutionInfo.Status.String() != "Running" {
				log.Println(workflowID, "is not running")
				return
			}
			err = c.TerminateWorkflow(ctx, workflowID, "", "terminate by "+id)
			if err != nil {
				log.Println("terminate failed", workflowID)
				return
			}
		}
	default:
		log.Println("no running workflows")
		return
	}
}

func printResults(result string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s %s\n\n", time.Now(), result)
}
