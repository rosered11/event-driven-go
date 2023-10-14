package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter/kafka"
)

func main() {
	// Configure Zipkin Kafka reporter
	kafkaBrokers := []string{"localhost:9092"} // Update with your Kafka brokers
	kafkaTopic := "zipkin"                     // The Kafka topic to send traces to

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
	tracerKafka, err := zipkin.NewTracer(
		kafkaReporter,
		zipkin.WithLocalEndpoint(kafka),
		zipkin.WithSampler(zipkin.AlwaysSample),
	)
	if err != nil {
		log.Fatalf("Failed to create Zipkin tracer: %v", err)
	}
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

	// Simulate some work
	// time.Sleep(200 * time.Millisecond)
	// endpoint1, _ := zipkin.NewEndpoint("myservice1-service", "myservice1:8081")
	// tracer1, err := zipkin.NewTracer(
	// 	kafkaReporter,
	// 	zipkin.WithLocalEndpoint(endpoint1),
	// 	zipkin.WithSampler(zipkin.AlwaysSample),
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to create Zipkin tracer: %v", err)
	// }
	// endpoint2, _ := zipkin.NewEndpoint("myservice2-service", "myservice2:9990")
	// tracer2, err := zipkin.NewTracer(
	// 	kafkaReporter,
	// 	zipkin.WithLocalEndpoint(endpoint2),
	// 	zipkin.WithSampler(zipkin.AlwaysSample),
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to create Zipkin tracer: %v", err)
	// }

	// // Create spans using each tracer
	// span1 := tracer1.StartSpan("operation-1")
	// span2 := tracer2.StartSpan("operation-2")

	// Simulate work in each span
	// time.Sleep(50 * time.Millisecond)

	// // Finish the spans
	// span1.Finish()
	// span2.Finish()

	fmt.Println("Tracing completed.")
}
