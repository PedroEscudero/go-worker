package main

import (
    "fmt"
    "time"
    "log"
    "github.com/streadway/amqp"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Println("worker", id, "started  job", j)
        time.Sleep(time.Second)
        fmt.Println("worker", id, "finished job", j)
        results <- j * 2
    }
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}

func main() {

    jobs := make(chan int, 100)
    results := make(chan int, 100)

    for i := 1; i<10; {
      for w := 1; w <= 3; w++ {
          go worker(w, jobs, results)
      }

      for j := 1; j <= 5; j++ {
          jobs <- j
      }
      //close(jobs)

      for a := 1; a <= 5; a++ {
          <-results
      }

      conn, err := amqp.Dial("amqp://guest:guest@192.20.128.2:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
  }
}
