package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

const endpointURL = "http://localhost:9411/api/v2/spans"

func newTracer() (*zipkin.Tracer, error) {
	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(endpointURL)

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: "my_service", Port: 8081}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	return t, err
}
func main() {
	tracer, err := newTracer()
	if err != nil {
		log.Fatal(err)
	}

	// We add the instrumented transport to the defaultClient
	// that comes with the zipkin-go library
	http.DefaultClient.Transport, err = zipkinhttp.NewTransport(
		tracer,
		zipkinhttp.TransportTrace(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandlerFactory(http.DefaultClient))
	r.HandleFunc("/foo", FooHandler)
	r.Use(zipkinhttp.NewServerMiddleware(
		tracer,
		zipkinhttp.SpanName("request")), // name for request span
	)
	log.Fatal(http.ListenAndServe(":8081", r))
}

func FooHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func HomeHandlerFactory(client *http.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body := bytes.NewBufferString("")
		res, err := client.Post("http://example.com", "application/json", body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if res.StatusCode > 399 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
