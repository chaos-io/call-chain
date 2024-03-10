package go_trace

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var tracer = opentracing.GlobalTracer()

func A1() {
	span := tracer.StartSpan("A1")
	defer fmt.Println("A1")
	defer span.Finish()
	B1(span)
}

func B1(parentSpan opentracing.Span) {
	span := tracer.StartSpan("B1", opentracing.ChildOf(parentSpan.Context()))
	defer fmt.Println("B1")
	defer span.Finish()
	// C1()
}

func clientSpan() {
	clientSpan := tracer.StartSpan("client")
	defer clientSpan.Finish()

	B1(clientSpan)

	url := "http://localhost:8082/publish"
	req, _ := http.NewRequest("GET", url, nil)

	// Set some tags on the clientSpan to annotate that it's the client span. The additional HTTP tags are useful for debugging purposes.
	ext.SpanKindRPCClient.Set(clientSpan)
	ext.HTTPUrl.Set(clientSpan, url)
	ext.HTTPMethod.Set(clientSpan, "GET")

	// Inject the client span context into the headers
	tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
}

// func C1() {
// 	defer Trace()
// 	D()
// }
//
// func D() {
// 	defer Trace()
// }

func TestTrace(t *testing.T) {
	// A1()
	clientSpan()
}

func TestPublishServer(t *testing.T) {
	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		// Extract the context from the headers
		fmt.Println("get a request")
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		serverSpan := tracer.StartSpan("server", ext.RPCServerOption(spanCtx))
		defer serverSpan.Finish()
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
}
