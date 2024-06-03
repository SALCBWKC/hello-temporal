package activity

import (
	"context"
	"fmt"
	"strings"
	"time"

	"hello-world-temporal/app/lib"
)

const (
	produceMsgNum = 1
	consumeMsgNum = 1
)

func Produce(ctx context.Context) (string, error) {
	topic := "test"
	messages := make([]string, 0, produceMsgNum)

	for i := 0; i < produceMsgNum; i++ {
		messages = append(messages, fmt.Sprintf("%d gen a message %s", i, time.Now()))
	}

	err := lib.Produce(ctx, topic, messages)
	if err != nil {
		return "produce failed", err
	}
	return fmt.Sprintf("produce %d messages", produceMsgNum), nil
}

func Consume(ctx context.Context) (string, error) {
	topic := "test"

	res, err := lib.GroupConsume(ctx, topic, consumeMsgNum)
	if err != nil {
		return "consume failed", err
	}
	return fmt.Sprintf("consume messages: {%s}", strings.Join(res, "---")), nil
}
