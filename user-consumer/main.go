package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = make(map[string]User) // Simulated database

type KafkaEvent struct {
	Action string `json:"action"`
	User   User   `json:"user"`
}

type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event KafkaEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to parse Kafka message: %v", err)
			continue
		}

		switch event.Action {
		case "create":
			users[event.User.ID] = event.User
			fmt.Printf("✅ User created: %+v\n", event.User)
		case "update":
			if _, exists := users[event.User.ID]; exists {
				users[event.User.ID] = event.User
				fmt.Printf("🔄 User updated: %+v\n", event.User)
			} else {
				fmt.Println("❌ User not found for update")
			}
		case "delete":
			delete(users, event.User.ID)
			fmt.Printf("🗑️ User deleted: %+v\n", event.User)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	brokers := []string{"localhost:9092"}
	group := "user-consumer-group"

	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}
	defer consumerGroup.Close()

	fmt.Println("User Consumer started, listening for messages...")

	for {
		ctx := context.Background()
		err := consumerGroup.Consume(ctx, []string{"users"}, ConsumerGroupHandler{})
		if err != nil {
			log.Printf("Error consuming Kafka messages: %v", err)
		}
	}
}
