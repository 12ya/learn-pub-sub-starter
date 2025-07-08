package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const amqpConn = "amqp://guest:guest@localhost:5672/"

func main() {
	conn, err := amqp.Dial(amqpConn)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("amqp connection was successful")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("couldnt open channel: %v", err)
	}

	if _, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.Durable,
	); err != nil {
		log.Fatal(err)
	}

	gamelogic.PrintServerHelp()

	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}

		switch inputs[0] {
		case routing.PauseKey:
			fmt.Println("Publishing paused game state")
			err = pubsub.PublishJSON(
				publishCh,
				string(routing.ExchangePerilDirect),
				string(routing.PauseKey),
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("error publishing JSON: %v", err)
			}
		case routing.ResumeKey:
			fmt.Println("Publishing resumed game state")
			err = pubsub.PublishJSON(
				publishCh,
				string(routing.ExchangePerilDirect),
				string(routing.PauseKey),
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("error publishing JSON: %v", err)
			}
		case routing.QuitKey:
			log.Println("goodbye")
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
