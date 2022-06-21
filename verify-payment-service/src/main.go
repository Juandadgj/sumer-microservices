package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

type payment struct {
	ID     int
	Client string
	Amount float32
}

type basicMessage struct {
	Message string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func publishNotification(body payment) {
	// The connection with RabbitMQ is configured
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// The queue to which the events are sent is declared
	q, err := ch.QueueDeclare(
		"notifications", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	newPayment, err := json.Marshal(body)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         newPayment,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s\n", newPayment)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	message := basicMessage{
		Message: "Welcome to Sumer Service",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func createPayment(w http.ResponseWriter, r *http.Request) {
	var newPayment payment
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "Payment invalid")
	}
	json.Unmarshal(reqBody, &newPayment)
	publishNotification(newPayment)
	time.Sleep(700 * time.Millisecond)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	message := basicMessage{
		Message: "The payment has been processed correctly. We will send a confirmation to your email and mobile phone shortly.",
	}
	json.NewEncoder(w).Encode(message)
}

func main() {

	// The server is configured
	router := mux.NewRouter()
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/payment", createPayment)
	log.Fatal(http.ListenAndServe(":3000", router))

}
