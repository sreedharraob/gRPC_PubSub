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

	//"cloud.google.com/go/profiler"
	pubsubgrpc "cloud.google.com/go/pubsub/apiv1"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
)

func publishMessages(w io.Writer, projectID, topicID, msg string, n int) error {
	ctx := context.Background()
	client, err := pubsubgrpc.NewPublisherClient(ctx)
	if err != nil {
		return fmt.Errorf("\n pubsub.NewClient: %v", err)
	}
	fmt.Fprintf(w, "Connected to publisher client successfully \n")

	var wg sync.WaitGroup
	var totalErrors uint64
	//var pubMsgs []*pubsubpb.PubsubMessage

	for i := 0; i < n; i++ {
		var pubMsg []*pubsubpb.PubsubMessage
		pubMsg = append(pubMsg, &pubsubpb.PubsubMessage{
			Data: []byte(msg + "-" + strconv.Itoa(i)),
			Attributes: map[string]string{
				"origin":   "go-exe",
				"username": "gcp-user",
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
	topicID := os.Getenv("TOPIC_ID")
	msg := os.Getenv("TOPIC_MSG")
	noOfMessages := 200

	if &projectID == nil || &topicID == nil || &msg == nil {
		log.Fatalln("Unable to find required args in env variables")
	}
	fmt.Printf("ProjectID: %s, TopicID: %s, msg: %s, no of messages: %s \n", projectID, topicID, msg, strconv.Itoa(noOfMessages))

	//// Google Cloud profiler

	// err := profiler.Start(profiler.Config{
	// 	ProjectID:      projectID,
	// 	Service:        "grpc-pubsub-publish",
	// 	DebugLogging:   true,
	// 	MutexProfiling: true,
	// })
	// if err != nil {
	// 	log.Fatalf("failed to start the profiler: %v", err)
	// }

	var w bytes.Buffer
	err := publishMessages(&w, projectID, topicID, msg, noOfMessages)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(&w)
	}

	fmt.Println("Press Enter to close")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
