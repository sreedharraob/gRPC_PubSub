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

	fmt.Println(publish.publishMessages(w, projectID, topicID, msg, noOfMessages))
}
