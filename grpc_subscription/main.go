package main

import (
	//"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	//"strconv"
	//"sync"
	//"runtime"
	//"sync/atomic"

	//"cloud.google.com/go/profiler"
	pubsubgrpc "cloud.google.com/go/pubsub/apiv1"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
)

func streamingPullMessages(w io.Writer, projectID, subID string) error {
	ctx := context.Background()
	client, err := pubsubgrpc.NewSubscriberClient(ctx)
	if err != nil {
		return fmt.Errorf("\n pubsub.NewClient: %v", err)
	}
	fmt.Fprintf(w, "Connected to subscriber client successfully \n")

	stream, err := client.StreamingPull(ctx)
	if err != nil {
		return fmt.Errorf("\n unable to make StreamingPull: %v", err)
	}
	fmt.Fprintf(w, "established stream with server \n")

	//var receivedAckIds []string
	go func() {
		var reqs []*pubsubpb.StreamingPullRequest
		reqs = append(reqs, &pubsubpb.StreamingPullRequest{
			Subscription: fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subID),
			//ClientId:     "go-exe-subscriber-1",
			//AckIds:       receivedAckIds,
		})
		for _, req := range reqs {
			err := stream.Send(req)
			if err != nil {
				//return fmt.Errorf("unable to send pull request: %v \n", err)
				log.Fatalf("error: %s", err)
			}
		}
		stream.CloseSend()
	}()
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("\n unable to process the request: %v", err)
		}
		for _, recvMsg := range resp.ReceivedMessages {
			log.Println(recvMsg)
			fmt.Fprintf(w, "\nGot message :%q\n", string(recvMsg.Message.Data))
			for k, v := range recvMsg.Message.Attributes {
				fmt.Fprintf(w, "%s=\"%s\"\n", k, v)
			}
			fmt.Fprintf(w, "AckId: %s \n", recvMsg.AckId)
			//fmt.Fprint(w, "DeliveryAttempt: %s \n", recvMsg.DeliveryAttempt)
			fmt.Fprintf(w, "MessageId: %s \n", recvMsg.Message.MessageId)
			fmt.Fprintf(w, "Message Published Time: %s \n", recvMsg.Message.PublishTime)
		}
	}
	return nil
}

func pullMessages(w io.Writer, projectID, subID string) error {
	ctx := context.Background()
	client, err := pubsubgrpc.NewSubscriberClient(ctx)
	if err != nil {
		return fmt.Errorf("\n pubsub.NewClient: %v", err)
	}
	fmt.Fprintf(w, "Connected to subscriber client successfully \n")

	subscriptionID := fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subID)
	req := &pubsubpb.PullRequest{
		Subscription: subscriptionID,
		MaxMessages:  400,
	}
	resp, err := client.Pull(ctx, req)
	if err != nil {
		return fmt.Errorf("\n unable to pull request: %v", err)
	}

	var receivedAckIds []string
	for _, recvMsg := range resp.ReceivedMessages {
		fmt.Fprintf(w, "\n Got message :%q\n", string(recvMsg.Message.Data))
		for k, v := range recvMsg.Message.Attributes {
			fmt.Fprintf(w, "%s=\"%s\"\n", k, v)
		}		
		fmt.Fprintf(w, "MessageId: %s \n Message Published Time: %s \n ", recvMsg.Message.MessageId, recvMsg.Message.PublishTime)		
		receivedAckIds = append(receivedAckIds, recvMsg.AckId)
	}
	ackRequest := &pubsubpb.AcknowledgeRequest{
		Subscription: subscriptionID,
		AckIds:       receivedAckIds,
	}
	client.Acknowledge(ctx, ackRequest)
	return nil
}

func main() {
	projectID := os.Getenv("PROJECT_ID")
	topicID := os.Getenv("TOPIC_ID")
	subID := os.Getenv("SUB_ID")

	if &projectID == nil || &topicID == nil || &subID == nil {
		log.Fatalln("Unable to find required args in env variables")
	}
	fmt.Printf("ProjectID: %s, TopicID: %s, SubID: %s \n", projectID, topicID, subID)

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
	start := time.Now()
	err := streamingPullMessages(&w, projectID, subID)
	//err := pullMessages(&w, projectID, subID)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(&w)
	}
	elapsed := time.Since(start)
	fmt.Printf("grpc sub pull took %s", elapsed)

	//fmt.Println("Press Enter to close")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
}
