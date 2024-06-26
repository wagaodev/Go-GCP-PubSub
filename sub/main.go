package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func main() {
	projectid := flag.String("projectid", "", "project id")
	topicName := flag.String("topic", "", "topic name")
	subid := flag.String("subid", "", "subscription id")

	flag.Parse()

	if err := pullMsg(os.Stdout, *projectid, *subid, *topicName); err != nil {
		log.Fatalf("pullMsgs: %v", err)
	}
}

func pullMsg(w io.Writer, projectID, subID, topic string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)
	exists, err := sub.Exists(ctx)
	if err != nil {
		return fmt.Errorf("client exists got err: %v", err)
	}

	if !exists {
		sub, err = client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{Topic: client.Topic(topic)})
		if err != nil {
			return fmt.Errorf("client.CreateSubscription got err: %v", err)
		}
	}

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}
	return nil
}
