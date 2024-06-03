package main

import (
	"hello-world-temporal/app/worker"
)

func main() {
	worker.Consume()
}
