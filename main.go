package main

import (	
	"fmt"
	"github.com/sreedharraob/gRPC_PubSub/topics"
	"os"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	topicID := os.Getenv("TOPIC_ID")
	msg := os.Getenv("TOPIC_MSG")
	noOfMessages := 100

	fmt.Println(publishMessages.publish(w, projectID, topicID, msg, noOfMessages))
}
