package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"

	"hello-world-temporal/app/activity"
)

const (
	timeout = time.Second * 60

	subProduceActivityNum = 1
	subConsumeActivityNum = 1
)

func SubProduceWorkflow(ctx workflow.Context) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: timeout,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	for i := 0; i < subProduceActivityNum; i++ {
		err := workflow.ExecuteActivity(ctx, activity.Produce).Get(ctx, &result)
		if err != nil {
			return "create produce activity failed", err
		}
		//log.Println(i, "workflow result is", result)
	}
	return fmt.Sprintf("create %d produce activities", subProduceActivityNum), nil
}

func MainProduceWorkflow(ctx workflow.Context) (string, error) {
	return "Main Produce workflow started", nil
}

func MainConsumeWorkflow(ctx workflow.Context) (string, error) {
	return "Main Consume workflow started", nil
}

func SubConsumeWorkflow(ctx workflow.Context) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: timeout,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	for i := 0; i < subConsumeActivityNum; i++ {
		err := workflow.ExecuteActivity(ctx, activity.Consume).Get(ctx, &result)
		if err != nil {
			return "create consume activity failed", err
		}
		//log.Println(i, "workflow result is", result)
	}
	return fmt.Sprintf("create %d consume activities", subConsumeActivityNum), nil
}
