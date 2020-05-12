package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	//"runtime"
	"sync/atomic"

	"cloud.google.com/go/profiler"
	pubsubgrpc "cloud.google.com/go/pubsub/apiv1"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
)

func publishMessages(w io.Writer, projectID, topicID, msg string, n int) error {
	ctx := context.Background()
	client, err := pubsubgrpc.NewPublisherClient(ctx)
	if err != nil {
		return fmt.Errorf("\n pubsub.NewClient: %v", err)
	}

	var wg sync.WaitGroup
	var totalErrors uint64
	//var pubMsgs []*pubsubpb.PubsubMessage

	for i := 100; i < n; i++ {
		var pubMsg []*pubsubpb.PubsubMessage
		pubMsg = append(pubMsg, &pubsubpb.PubsubMessage{
			Data: []byte(msg + "-" + strconv.Itoa(i)),
			Attributes: map[string]string{
				"origin":   "golang",
				"username": "gcp-sree",
				"action":   "row-delete-" + strconv.Itoa(i),
			},
		})

		req := &pubsubpb.PublishRequest{
			Topic:    fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
			Messages: pubMsg,
		}

		wg.Add(1)
		go func(i int, r *pubsubpb.PublishRequest) {
			defer wg.Done()
			resp, err := client.Publish(ctx, req)
			if err != nil {
				fmt.Fprintf(w, "Failed to publish :%v\n", err)
				atomic.AddUint64(&totalErrors, 1)
				return
			}
			fmt.Fprintf(w, "Published message %d; msg ID: %v\n", i, resp.MessageIds)
		}(i, req)
	}

	wg.Wait()

	if totalErrors > 0 {
		return fmt.Errorf("\n %d of %d messages did not publish successfully", totalErrors, n)
	}
	return nil
}

func main() {
	projectID := os.Getenv("PROJECT_ID")

	err := profiler.Start(profiler.Config{
		ProjectID:      projectID,
		Service:        "grpc-pubsub-publish",
		DebugLogging:   true,
		MutexProfiling: true,
	})
	if err != nil {
		log.Fatalf("failed to start the profiler: %v", err)
	}

	
	topicID := os.Getenv("TOPIC_ID")
	msg := os.Getenv("TOPIC_MSG")
	noOfMessages := 200

	var w bytes.Buffer
	publishMessages(&w, projectID, topicID, msg, noOfMessages)
	fmt.Println(&w)

	fmt.Println("Press Enter to close")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
