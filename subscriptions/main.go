package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"	
	"time"

	"cloud.google.com/go/pubsub"
)

func pullMessages(w io.Writer, projectID, subID string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)
	sub.ReceiveSettings.Synchronous = false
	sub.ReceiveSettings.NumGoroutines = runtime.NumCPU()
	fmt.Fprintf(w, "number of CPU in client: %q\n", strconv.Itoa(runtime.NumCPU()))

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	messagesChannel := make(chan *pubsub.Message)

	go func() {
		for {
			select {
			case msg := <-messagesChannel:
				fmt.Fprintf(w, "\n Got message :%q\n", string(msg.Data))
				for k, v := range msg.Attributes {
					fmt.Fprintf(w, "%s=\"%s\"\n", k, v)
				}
				msg.Ack()
			case <-ctx.Done():
				return
			}
		}
	}()

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		messagesChannel <- msg
	})
	if err != nil {
		return fmt.Errorf("Recieve: %v", err)
	}
	close(messagesChannel)

	return nil
}

func listTopics(client pubsub.Client, topicID string) (*pubsub.Topic, error) {
	t := client.Topic(topicID)
	if t == nil {
		return nil, fmt.Errorf("No topic found with id %d", topicID)
	}
	return t, nil
}

func main() {
	projectID := os.Getenv("PROJECT_ID")
	subID := os.Getenv("SUB_ID")

	var w bytes.Buffer
	pullMessages(&w, projectID, subID)
	fmt.Println(&w)
}