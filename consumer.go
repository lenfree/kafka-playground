package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	f, err := os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	kafkaPath := os.Getenv("KAFKA_PATH")
	brokers := strings.Split(os.Getenv("BROKERS"), ",")
	brokerService := os.Getenv("BROKER_SERVICE_NAME")
	zookeeperService := os.Getenv("ZOOKEEPER_SERVICE_NAME")
	clusterName := os.Getenv("ClusterName")
	zookeeperCluster := zookeeperService + clusterName
	topic := os.Getenv("TOPIC")

	logger := log.New(f, "kafka-consumer-playground", 0)
	logger.SetPrefix(topic + " ")

	logger.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)

	logger.Printf("Set brokers: %s\n", brokers)
	logger.Printf("Set topic: %s\n", topic)

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		logger.Printf("error: %s\n", err.Error())

		logger.Printf("%s\n", "Execute some shell debugging commands")
		cmdName := "dig"
		cmdArgs := []string{"@localhost", "-p", "8600", "+tcp", "+short", brokerService}

		var cmdOut []byte
		cmdOut, err = exec.Command(cmdName, cmdArgs...).Output()
		if err != nil {
			logger.Printf("There was an error running git rev-parse command: %s\n", err.Error())
		}
		logger.Printf("Result of %s: %s\n", cmdName, string(cmdOut))

		logger.Printf("%s\n", "Try connecting with kafka cli")
		kafkaCmd := kafkaPath + "kafka-topics.sh"
		kafkaCmdArgs := []string{"--zookeeper", zookeeperCluster, "--describe", "--topic", topic}

		var kafkaCmdOut []byte
		kafkaCmdOut, err = exec.Command(kafkaCmd, kafkaCmdArgs...).Output()
		if err != nil {
			logger.Printf("There was an error running git rev-parse command: %s\n", err.Error())
		}
		logger.Printf("Result of %s: %s\n", kafkaCmd, string(kafkaCmdOut))

	}

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			logger.Printf("Consumed message offset %d\n", msg.Offset)
			msgkey := string(msg.Key[:])
			msgval := string(msg.Value[:])
			logger.Println("Consumed message: ", msgkey)
			logger.Println("Consumed val:", msgval)
			consumed++
		}
		logger.Printf("Consumed: %d\n", consumed)
	}
}
