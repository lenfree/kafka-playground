package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	hosts := []string{"localhost:9092"}
	client, err := sarama.NewClient(hosts, config)
	defer client.Close()
	if err != nil {
		log.Panicln(err)
	}

	topics, err := client.Topics()
	log.Printf("Number of topics #: %+#v\n", len(topics))
	if err != nil {
		log.Panicln(err)
	}

	brokers := client.Brokers()
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Number ofbrokers #: %+#v\n", len(brokers))

	broker := sarama.NewBroker("localhost:9092")
	err = broker.Open(nil)
	defer broker.Close()

	req := sarama.MetadataRequest{Topics: topics}
	metadataResponse, err := broker.GetMetadata(&req)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("Number of res topics #: %+#v\n", len(metadataResponse.Topics))

	topicName := make(chan string)
	var filterTopic string
	filterTopic = "lo-werphat"
	for i := 1; i < 30; i++ {
		go printTopics(topicName, filterTopic, i)
	}

	for _, topic := range metadataResponse.Topics {
		topicName <- topic.Name
	}
}

func printTopics(topics chan string, filterTopic string, routine int) {
	// TODO: use regex instead
	for topic := range topics {
		if strings.Contains(topic, "_samza_checkpoint_ver_1_for_") && !strings.Contains(topic, filterTopic) {
			log.Printf("Got %s === goroutine %d\n", topic, routine)
		}

		if strings.Contains(topic, "_samza_coordinator_") && !strings.Contains(topic, filterTopic) {
			log.Printf("Got %s === goroutine %d\n", topic, routine)
		}

		if strings.Contains(topic, "--samza-store") && !strings.Contains(topic, filterTopic) {
			log.Printf("Got %s === goroutine %d\n", topic, routine)
		}
	}
}
