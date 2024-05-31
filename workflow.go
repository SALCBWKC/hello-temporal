package app

import (
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	timeout = time.Second * 60
)

func MainProduceWorkflow(ctx workflow.Context) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: timeout,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	const subWorkflowNum = 2
	for i := 0; i < subWorkflowNum; i++ {
		err := workflow.ExecuteActivity(ctx, ProduceActivity).Get(ctx, &result)
		if err != nil {
			return "create sub produce workflow failed", err
		}
		log.Println(i, "workflow result is", result)
	}
	return fmt.Sprintf("create %d sub produce workflows", subWorkflowNum), nil
}

func MainConsumeWorkflow(ctx workflow.Context) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: timeout,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	const subWorkflowNum = 2
	for i := 0; i < subWorkflowNum; i++ {
		err := workflow.ExecuteActivity(ctx, ConsumeActivity).Get(ctx, &result)
		if err != nil {
			return "create sub Consume workflow failed", err
		}
		log.Println(i, "workflow result is", result)
	}
	return fmt.Sprintf("create %d sub consume workflows", subWorkflowNum), nil
}
