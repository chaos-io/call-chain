package main

import (
	"fmt"

	"github.com/opentracing/opentracing-go/log"

	"github.com/chaos-io/go-trace/example/lib/tracing"
)

func main() {
	tracer, closer := tracing.Init("hello-world")
	defer closer.Close()

	helloTo := "example1"

	span := tracer.StartSpan("say-hello")
	span.SetTag("hello-to", helloTo)

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	println(helloStr)
	span.LogKV("event", "println")

	span.Finish()
}
