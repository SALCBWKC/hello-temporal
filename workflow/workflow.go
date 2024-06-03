package workflow

import (
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/workflow"

	"hello-world-temporal/app/activity"
)

const (
	timeout = time.Second * 60

	subProduceWorkflowNum = 200
	subConsumeWorkflowNum = 200
)

func MainProduceWorkflow(ctx workflow.Context) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: timeout,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	for i := 0; i < subProduceWorkflowNum; i++ {
		err := workflow.ExecuteActivity(ctx, activity.Produce).Get(ctx, &result)
		if err != nil {
			return "create sub produce workflow failed", err
		}
		log.Println(i, "workflow result is", result)
	}
	return fmt.Sprintf("create %d sub produce workflows", subProduceWorkflowNum), nil
}

func MainConsumeWorkflow(ctx workflow.Context) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: timeout,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	for i := 0; i < subConsumeWorkflowNum; i++ {
		err := workflow.ExecuteActivity(ctx, activity.Consume).Get(ctx, &result)
		if err != nil {
			return "create sub Consume workflow failed", err
		}
		log.Println(i, "workflow result is", result)
	}
	return fmt.Sprintf("create %d sub consume workflows", subConsumeWorkflowNum), nil
}
