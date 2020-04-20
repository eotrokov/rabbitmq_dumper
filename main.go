package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	argsWithProg := os.Args[1:]

	url := argsWithProg[0]
	name := argsWithProg[1]
	connection, err := amqp.Dial(url)
	defer connection.Close()

	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}
	channel, err := connection.Channel()

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}

	// We consume data from the queue named Test using the channel we created in go.
	msgs, err := channel.Consume(name, "", false, false, false, false, nil)
	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}
	f, err := os.Create(name + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	for msg := range msgs {
		_, err := f.WriteString(string(msg.Body) + "\n---\n")
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
