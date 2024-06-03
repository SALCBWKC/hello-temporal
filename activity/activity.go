package activity

import (
	"context"
	"fmt"
	"strings"
	"time"

	"hello-world-temporal/app/lib"
)

func Produce(ctx context.Context) (string, error) {
	topic := "test"
	const msgNum = 5
	messages := make([]string, 0, msgNum)

	for i := 0; i < msgNum; i++ {
		messages = append(messages, fmt.Sprintf("%d gen a message %s", i, time.Now()))
	}

	err := lib.Produce(ctx, topic, messages)
	if err != nil {
		return "produce failed", err
	}
	return fmt.Sprintf("produce %d messages", msgNum), nil
}

func Consume(ctx context.Context) (string, error) {
	topic := "test"
	const msgNum = 5

	res, err := lib.GroupConsume(ctx, topic, msgNum)
	if err != nil {
		return "consume failed", err
	}
	return fmt.Sprintf("consume messages: {%s}", strings.Join(res, "---")), nil
}
