package main

import (
	"context"
	"flag"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func main() {

	ctx := context.Background()

	projectid := flag.String("projectid", "", "project id")
	topicName := flag.String("topic", "", "topic name")
	msg := flag.String("msg", "", "message")

	flag.Parse()

	fmt.Printf("Project ID: %s\n", *projectid)
	fmt.Printf("Topic Name: %s\n", *topicName)
	fmt.Printf("Message: %s\n", *msg)

	client, err := pubsub.NewClient(ctx, *projectid)
	if err != nil {
		panic(err)
	}

	defer client.Close()
	topic := client.Topic(*topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		topic, err = client.CreateTopic(ctx, *topicName)
		if err != nil {
			panic(err)
		}
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(*msg),
	})

	_, err = result.Get(ctx)
	if err != nil {
		panic(err)
	}
}
