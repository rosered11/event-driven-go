package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	workerfail()
	// c1 := make(chan string)
	// c2 := make(chan string)

	// // Use goroutines to perform Kafka and Redis operations in parallel
	// go func() {
	// 	worker1()
	// 	c1 <- "1"
	// }()

	// go func() {
	// 	workerfail()
	// 	c2 <- "2"
	// }()

	// for i := 0; i < 2; i++ {
	// 	select {
	// 	case msg1 := <-c1:
	// 		fmt.Println("received", msg1)
	// 	case msg2 := <-c2:
	// 		fmt.Println("received", msg2)
	// 	}
	// }
}

func worker1() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:30881,localhost:30882",
		"group.id":          "worker1",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"producer"}, nil)

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

func worker2() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:30881,localhost:30882",
		"group.id":          "worker2",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"producer"}, nil)

	// A signal handler or similar could be used to set this to false to break the loop.
	run := true

	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

func workerfail() {
	// Configure the consumer
	config := &kafka.ConfigMap{
		"bootstrap.servers":  "localhost:30881,localhost:30882",
		"group.id":           "worker2",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Error creating Kafka consumer: %v\n", err)
		os.Exit(1)
	}
	defer consumer.Close()

	// Subscribe to the topic
	topic := "producer"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Error subscribing to topic: %v\n", err)
		os.Exit(1)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case sig := <-signals:
			fmt.Printf("Received signal: %v\n", sig)
			return
		default:
			ev := consumer.Poll(100) // Poll for messages

			switch e := ev.(type) {
			case *kafka.Message:
				if e.TopicPartition.Error != nil {
					fmt.Printf("Error receiving message: %v\n", e.TopicPartition.Error)
				} else {
					fmt.Printf("Received message: key=%s, value=%s, offset=%d\n", string(e.Key), string(e.Value), e.TopicPartition.Offset)

					// Process the message
					if processMessage(e) {
						// If processing is successful, manually commit the offset for the processed message
						offsets := []kafka.TopicPartition{
							{
								Topic:     e.TopicPartition.Topic,
								Partition: e.TopicPartition.Partition,
								Offset:    kafka.Offset(e.TopicPartition.Offset + 1), // Commit the next offset
							},
						}
						_, err := consumer.CommitOffsets(offsets)
						if err != nil {
							fmt.Printf("Error committing offset: %v\n", err)
						}
					} else {
						// offsets := []kafka.TopicPartition{
						// 	{
						// 		Topic:     e.TopicPartition.Topic,
						// 		Partition: e.TopicPartition.Partition,
						// 		Offset:    kafka.Offset(e.TopicPartition.Offset), // Commit the next offset
						// 	},
						// }
						// _, err := consumer.CommitOffsets(offsets)
						// if err != nil {
						// 	fmt.Printf("Error committing offset: %v\n", err)
						// }
					}
				}
			}
		}
	}
}
func processMessage(msg *kafka.Message) bool {
	// Implement your message processing logic here.
	// If processing is successful, return true; otherwise, return false.
	return false // Change to false in case of processing failure
}
