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
	meta, err := broker.GetMetadata(&req)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("Number of res topicss #: %+#v\n", len(meta.Topics))
	for _, topic := range meta.Topics {
		// TODO: use regex instead
		if strings.Contains(topic.Name, "_samza_checkpoint_ver_1_for_") && !strings.Contains(topic.Name, "vodka") && !strings.Contains(topic.Name, "lo-werphat") {
			log.Printf("Got %s\n", topic.Name)
		}

		if strings.Contains(topic.Name, "_samza_coordinator_") && !strings.Contains(topic.Name, "vodka") && !strings.Contains(topic.Name, "lo-werphat") {
			log.Printf("Got %s\n", topic.Name)
		}

		if strings.Contains(topic.Name, "--samza-store") && !strings.Contains(topic.Name, "vodka") && !strings.Contains(topic.Name, "lo-werphat") {
			log.Printf("Got %s\n", topic.Name)
		}
	}

}
