package main

import (
	"hello-world-temporal/app/start"
	"hello-world-temporal/app/worker"
)

func main() {
	go worker.Produce()
	go worker.Consume()
	go start.Produce()
	start.Consume()
}
