# kafka-go-micro-service

This project demonstrates a simple Kafka producer and consumer using Go. The producer sends user data to a Kafka topic, and the consumer reads from that topic.

## Prerequisites

- Docker and Docker Compose
- Go 1.16 or later

## Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/supacheep-first/kafka-go-micro-service.git
   cd kafka-go-micro-service
   ```

2. Start Kafka and Zookeeper using Docker Compose:

   ```sh
   docker-compose up -d
   ```

3. Install Go dependencies:
   ```sh
   go mod tidy
   ```

## Running the Producer (User API)

1. Navigate to the `user-api` directory:

   ```sh
   cd user-api
   ```

2. Run the API server:

   ```sh
   go run main.go
   ```

3. The API server will be running on port 8080. You can create a user by sending a POST request to `http://localhost:8080/users` with a JSON body:
   ```json
   {
     "id": "1",
     "name": "John Doe",
     "email": "john.doe@example.com"
   }
   ```

## Running the Consumer (User Consumer)

1. Navigate to the `user-consumer` directory:

   ```sh
   cd user-consumer
   ```

2. Run the consumer:

   ```sh
   go run main.go
   ```

3. The consumer will start listening for messages on the `users` topic and print received user data to the console.

## Stopping the Services

To stop the Kafka and Zookeeper services, run:

```sh
docker-compose down
```

## License

This project is licensed under the MIT License.
