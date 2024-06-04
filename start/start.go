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
	produceDuration = 1 * time.Minute
	consumeDuration = 2 * time.Minute

	// 200 mock 20, for easy test
	subProduceWorkflowNum = 20
	subConsumeWorkflowNum = 20

	produceWorkflowPrefix = "produce-cron-workflow"
	consumeWorkflowPrefix = "consume-cron-workflow"
)

var (
	consumeChan = make(chan string, subProduceWorkflowNum)
	produceChan = make(chan string, subConsumeWorkflowNum)
)

type startInfo struct {
	c              client.Client
	idPrefix       string
	id             string
	queue          string
	ch             chan string
	subWorkflowNum int
	workflow       interface{}
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

func produce() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Println("unable to create Temporal client", err)
	}
	defer c.Close()

	ctx := context.Background()
	info := &startInfo{
		c:              c,
		idPrefix:       produceWorkflowPrefix,
		queue:          lib.ProduceTaskQueue,
		ch:             produceChan,
		subWorkflowNum: subProduceWorkflowNum,
		workflow:       workflow.MainProduceWorkflow,
	}

	mainStart(ctx, info)
	info.workflow = workflow.SubProduceWorkflow
	subStart(ctx, info)
}

func consume() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Println("unable to create Temporal client", err)
	}
	defer c.Close()

	ctx := context.Background()
	info := &startInfo{
		c:              c,
		idPrefix:       consumeWorkflowPrefix,
		queue:          lib.ConsumeTaskQueue,
		ch:             consumeChan,
		subWorkflowNum: subConsumeWorkflowNum,
		workflow:       workflow.MainConsumeWorkflow,
	}

	mainStart(ctx, info)
	info.workflow = workflow.SubConsumeWorkflow
	subStart(ctx, info)
}

// main workflow to terminate
func mainStart(ctx context.Context, info *startInfo) {
	info.id = fmt.Sprintf("main-"+info.idPrefix+"-%d", time.Now().Unix())
	terminateWorkflow(ctx, info)
	start(ctx, info)
}

// sub workflow
func subStart(ctx context.Context, info *startInfo) {
	for i := 0; i < info.subWorkflowNum; i++ {
		info.id = fmt.Sprintf("%d-sub-"+info.idPrefix+"-%d", i, time.Now().Unix())
		start(ctx, info)
	}
}

func start(ctx context.Context, info *startInfo) {
	options := client.StartWorkflowOptions{
		ID:        info.id,
		TaskQueue: info.queue,
	}

	we, err := info.c.ExecuteWorkflow(ctx, options, info.workflow)
	if err != nil {
		log.Println(info.id, "is unable to complete, err:", err)
	}
	getResults(ctx, we)

	info.ch <- info.id
}

func terminateWorkflow(ctx context.Context, info *startInfo) {
	for {
		select {
		case workflowID, ok := <-info.ch:
			if ok {
				desc, err := info.c.DescribeWorkflowExecution(ctx, workflowID, "")
				if err != nil {
					log.Println(workflowID, "describe failed", err)
				}
				if desc.WorkflowExecutionInfo.Status.String() != "Running" {
					log.Println(workflowID, "is not running while initiate", info.id)
				}
				err = info.c.TerminateWorkflow(ctx, workflowID, "", "terminate by "+info.id)
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

func getResults(ctx context.Context, we client.WorkflowRun) {
	var result string
	err := we.Get(ctx, &result)
	if err != nil {
		log.Println(we.GetID(), "is unable to get result", err)
	}
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())
	fmt.Printf("\n%s %s\n\n", time.Now(), result)
}
