package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var producer sarama.SyncProducer

func initKafkaProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	producer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
}

func sendToKafka(action string, user User) {
	message, _ := json.Marshal(map[string]interface{}{"action": action, "user": user})
	msg := &sarama.ProducerMessage{
		Topic: "users",
		Value: sarama.StringEncoder(message),
	}
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send Kafka message: %v", err)
	}
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	sendToKafka("create", user)
	c.JSON(http.StatusOK, gin.H{"message": "User event sent", "user": user})
}

func updateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	sendToKafka("update", user)
	c.JSON(http.StatusOK, gin.H{"message": "User update event sent", "user": user})
}

func deleteUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	sendToKafka("delete", user)
	c.JSON(http.StatusOK, gin.H{"message": "User delete event sent", "user": user})
}

func main() {
	initKafkaProducer()
	defer producer.Close()

	r := gin.Default()
	r.POST("/users", createUser)
	r.PUT("/users/update", updateUser)
	r.DELETE("/users/delete", deleteUser)

	fmt.Println("User API running on port 8080")
	r.Run(":8080")
}
