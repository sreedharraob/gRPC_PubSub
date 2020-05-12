package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
)

func publishMessages(w io.Writer, projectID, topicID, msg string, n int) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("\n pubsub.NewClient: %v", err)
	}

	var wg sync.WaitGroup
	var totalErrors uint64
	t := client.Topic(topicID)

	for i := 0; i < n; i++ {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte(msg + "-" + strconv.Itoa(i)),
			Attributes: map[string]string{
				"origin":   "golang",
				"username": "gcp-sree",
				"action":   "row-delete-" + strconv.Itoa(i),
			},
		})		

		wg.Add(1)
		go func(i int, res *pubsub.PublishResult) {
			defer wg.Done()
			id, err := result.Get(ctx)
			if err != nil {
				fmt.Fprintf(w, "Failed to publish :%v\n", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			fmt.Fprintf(w, "Published message %d; msg ID: %v\n", i, id)					
		}(i, result)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("\n %d of %d messages did not publish successfully", totalErrors, n)
	}
	return nil
}

func main() {
	projectID := os.Getenv("PROJECT_ID")
	topicID := os.Getenv("TOPIC_ID")
	msg := os.Getenv("TOPIC_MSG")
	noOfMessages := 100

	var w bytes.Buffer
	publishMessages(&w, projectID, topicID, msg, noOfMessages)
	fmt.Println(&w)

	fmt.Println("Press Enter to close")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
