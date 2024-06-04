package start

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.temporal.io/sdk/client"

	"hello-world-temporal/app/lib"
	"hello-world-temporal/app/workflow"
)

const (
	produceDuration = 1 * time.Minute
	consumeDuration = 2 * time.Minute

	// 200 mock 20, for easy test
	subProduceWorkflowNum = 20
	subConsumeWorkflowNum = 20
)

var (
	consumeChan = make(chan string, subProduceWorkflowNum)
	produceChan = make(chan string, subConsumeWorkflowNum)
)

func Consume() {
	// 模拟Cron逻辑，每2分钟启动一个主工作流
	ticker := time.NewTicker(consumeDuration)
	for {
		select {
		case <-ticker.C:
			go consume()
		}
	}
}

func Produce() {
	// 模拟Cron逻辑，每1分钟启动一个主工作流
	ticker := time.NewTicker(produceDuration)
	for {
		select {
		case <-ticker.C:
			go produce()
		}
	}
}

func produce() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Println("unable to create Temporal client", err)
	}
	defer c.Close()

	ctx := context.Background()
	wg := new(sync.WaitGroup)
	wg.Add(subProduceWorkflowNum + 1)

	go mainProduce(ctx, c, produceChan, wg)
	go subProduce(ctx, c, produceChan, wg)

	wg.Wait()
}

func mainProduce(ctx context.Context, c client.Client, ch chan<- string, wg *sync.WaitGroup) {
	workflowID := fmt.Sprintf("main-produce-cron-workflow-%d", time.Now().Unix())
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: lib.ProduceTaskQueue,
	}

	terminateWorkflow(ctx, c, workflowID, produceChan)

	we, err := c.ExecuteWorkflow(ctx, options, workflow.MainProduceWorkflow)
	if err != nil {
		log.Println(workflowID, "is unable to complete, err:", err)
	}
	getResults(ctx, we)

	ch <- workflowID

	wg.Done()
}

func getResults(ctx context.Context, we client.WorkflowRun) {
	var result string
	err := we.Get(ctx, &result)
	if err != nil {
		log.Println(we.GetID(), "is unable to get result", err)
	}
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())
	fmt.Printf("\n%s %s\n\n", time.Now(), result)
}

func subProduce(ctx context.Context, c client.Client, ch chan<- string, wg *sync.WaitGroup) {
	var options client.StartWorkflowOptions
	for i := 0; i < subProduceWorkflowNum; i++ {
		go func(i int) {
			workflowID := fmt.Sprintf("%d-sub-produce-cron-workflow-%d", i, time.Now().Unix())
			options = client.StartWorkflowOptions{
				ID:        workflowID,
				TaskQueue: lib.ProduceTaskQueue,
			}

			we, err := c.ExecuteWorkflow(ctx, options, workflow.SubProduceWorkflow)
			if err != nil {
				log.Println(workflowID, "is unable to complete", err)
			}
			getResults(ctx, we)

			ch <- workflowID

			wg.Done()
		}(i)
	}
}

func mainConsume(ctx context.Context, c client.Client, ch chan<- string, wg *sync.WaitGroup) {
	workflowID := fmt.Sprintf("main-consume-cron-workflow-%d", time.Now().Unix())
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: lib.ConsumeTaskQueue,
	}

	terminateWorkflow(ctx, c, workflowID, consumeChan)

	we, err := c.ExecuteWorkflow(ctx, options, workflow.MainConsumeWorkflow)
	if err != nil {
		log.Println(workflowID, "is unable to complete, err:", err)
	}
	getResults(ctx, we)

	ch <- workflowID

	wg.Done()
}

func subConsume(ctx context.Context, c client.Client, ch chan<- string, wg *sync.WaitGroup) {
	var options client.StartWorkflowOptions
	for i := 0; i < subConsumeWorkflowNum; i++ {
		go func(i int) {
			workflowID := fmt.Sprintf("%d-sub-consume-cron-workflow-%d", i, time.Now().Unix())
			options = client.StartWorkflowOptions{
				ID:        workflowID,
				TaskQueue: lib.ConsumeTaskQueue,
			}

			we, err := c.ExecuteWorkflow(ctx, options, workflow.SubConsumeWorkflow)
			if err != nil {
				log.Println(workflowID, "is unable to complete", err)
			}
			getResults(ctx, we)

			ch <- workflowID

			wg.Done()
		}(i)
	}
}

func consume() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Println("unable to create Temporal client", err)
	}
	defer c.Close()

	ctx := context.Background()
	wg := new(sync.WaitGroup)
	wg.Add(subConsumeWorkflowNum + 1)

	go mainConsume(ctx, c, consumeChan, wg)
	go subConsume(ctx, c, consumeChan, wg)

	wg.Wait()
}

func terminateWorkflow(ctx context.Context, c client.Client, curID string, ch <-chan string) {
	for {
		select {
		case workflowID, ok := <-ch:
			if ok {
				desc, err := c.DescribeWorkflowExecution(ctx, workflowID, "")
				if err != nil {
					log.Println(workflowID, "describe failed", err)
				}
				if desc.WorkflowExecutionInfo.Status.String() != "Running" {
					log.Println(workflowID, "is not running while initiate", curID)
				}
				err = c.TerminateWorkflow(ctx, workflowID, "", "terminate by "+curID)
				if err != nil {
					log.Println("terminate failed", workflowID)
				}
			}
		default:
			log.Println("no running workflows")
			return
		}
	}
}
