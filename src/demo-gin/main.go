package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter/kafka"
)

func main() {
	var traceId, spanId string
	for i := 0; i < 1; i++ {
		traceId, spanId = sendLog()
		time.Sleep(time.Second)
	}
	consumerLog(traceId, spanId)
	fmt.Println("Tracing completed.")
}

func consumerLog(traceId, spanId string) {
	kafkaBrokers := []string{"localhost:30881", "localhost:30882"} // Update with your Kafka brokers
	kafkaTopic := "zipkin"                                         // The Kafka topic to send traces to
	kafkaReporter, err := kafka.NewReporter(
		kafkaBrokers,
		kafka.Topic(kafkaTopic),
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka reporter: %v", err)
	}
	defer kafkaReporter.Close()

	traceID, _ := model.TraceIDFromHex(traceId)
	result, _ := strconv.ParseUint(spanId, 16, 64)
	spanID := model.ID(result)
	endpoint, _ := zipkin.NewEndpoint("consumption-service", "80")
	tracer, _ := zipkin.NewTracer(
		kafkaReporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(zipkin.AlwaysSample),
	)
	// span := tracer.StartSpan("message-consumption")
	span := tracer.StartSpan("message-consumption", zipkin.Parent(model.SpanContext{
		TraceID: traceID,
		ID:      spanID,
	}))
	defer span.Finish()

	time.Sleep(50 * time.Millisecond)

}

func sendLog() (tracId, spanId string) {
	// Configure Zipkin Kafka reporter
	kafkaBrokers := []string{"localhost:30881", "localhost:30882"} // Update with your Kafka brokers
	kafkaTopic := "zipkin"                                         // The Kafka topic to send traces to

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForLocal
	producerConfig.Producer.Return.Successes = true

	kafkaReporter, err := kafka.NewReporter(
		kafkaBrokers,
		kafka.Topic(kafkaTopic),
		// kafka.ProducerConfig(producerConfig),
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka reporter: %v", err)
	}
	defer kafkaReporter.Close()

	// Create a Zipkin Tracer
	endpoint, _ := zipkin.NewEndpoint("myservice3-service", "8099")
	kafka, _ := zipkin.NewEndpoint("kafka-service", "5444")
	redis, _ := zipkin.NewEndpoint("redis-service", "9092")
	tracer, err := zipkin.NewTracer(
		kafkaReporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(zipkin.AlwaysSample),
	)

	if err != nil {
		log.Fatalf("Failed to create Zipkin tracer: %v", err)
	}
	tracerKafka, _ := zipkin.NewTracer(
		kafkaReporter,
		zipkin.WithLocalEndpoint(kafka),
		zipkin.WithSampler(zipkin.AlwaysSample),
	)

	tracerRedis, err := zipkin.NewTracer(
		kafkaReporter,
		zipkin.WithLocalEndpoint(redis),
		zipkin.WithSampler(zipkin.AlwaysSample),
	)
	if err != nil {
		log.Fatalf("Failed to create Zipkin tracer: %v", err)
	}

	// Start a new span
	span := tracer.StartSpan("/send")

	defer span.Finish()

	// Add tags and annotations to the span if needed
	span.Tag("custom-tag", "tag-value")
	span.Annotate(time.Now(), "custom-annotation")

	producerService, _ := zipkin.NewEndpoint("producer-service", "0.0.0.0")
	tracerProducer, err := zipkin.NewTracer(
		kafkaReporter,
		zipkin.WithLocalEndpoint(producerService),
		zipkin.WithSampler(zipkin.AlwaysSample),
	)
	if err != nil {
		log.Fatalf("Failed to create Zipkin tracer: %v", err)
	}
	ctx := tracerProducer.StartSpan("ProduceMessage")
	time.Sleep(50 * time.Millisecond)
	ctx.Annotate(time.Now(), "ms")
	kafkaChild := tracer.StartSpan("kafka-child", zipkin.Parent(ctx.Context()))
	time.Sleep(50 * time.Millisecond)
	kafkaChild.Finish()
	ctx.Finish()

	childSpan := tracer.StartSpan("child-operation", zipkin.Parent(span.Context()))
	defer childSpan.Finish()

	childSpan.Annotate(time.Now(), "child-annotation")

	time.Sleep(50 * time.Millisecond)

	c1 := make(chan string)
	c2 := make(chan string)

	// Use goroutines to perform Kafka and Redis operations in parallel
	go func() {
		// Start a child span for Kafka
		kafkaSpan := tracerKafka.StartSpan("kafka-operation", zipkin.Parent(childSpan.Context()))

		// Simulate some work for the Kafka span
		time.Sleep(50 * time.Millisecond)
		kafkaSpan.Annotate(time.Now(), "kafka-annotation")
		// Finish the Kafka span
		kafkaSpan.Finish()
		c1 <- "1"
	}()

	go func() {
		// Start a child span for Redis
		redisSpan := tracerRedis.StartSpan("redis-operation", zipkin.Parent(childSpan.Context()))

		// Simulate some work for the Redis span
		time.Sleep(50 * time.Millisecond)
		redisSpan.Annotate(time.Now(), "redis-annotation")
		// Finish the Redis span
		redisSpan.Finish()
		c2 <- "2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
	tracId = childSpan.Context().TraceID.String()
	spanId = childSpan.Context().ID.String()
	return tracId, spanId
}
