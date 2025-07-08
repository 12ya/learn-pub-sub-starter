package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const amqpConn = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril client...")
	conn, err := amqp.Dial(amqpConn)
	if err != nil {
		log.Fatalf("error dialing the AMQP server: %v", err)
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Printf("error getting username: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.Transient)
	if err != nil {
		log.Fatalf("error declaring and binding the queue: %v", err)
	}
	fmt.Printf("Queue %v is declared and bound!\n", queue)

	gameState := gamelogic.NewGameState(username)
	for {
		inputs := gamelogic.GetInput()
		if len(inputs) <= 0 {
			continue
		}

		switch inputs[0] {
		case "spawn":
			spawnConfig := strings.Split(inputs[0], " ")
			gameState.CommandSpawn(spawnConfig)

		case "move":
			moveConfig := strings.Split(inputs[0], " ")
			armyMove, err := gameState.CommandMove(moveConfig)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("move %v 1", armyMove.ToLocation)

		case "status":
			gameState.CommandStatus()

		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			fmt.Println("Spamming not allowed yet!")

		case "quit":
			gamelogic.PrintQuit()
			return

		default:
			log.Println("unknown command")
			continue

		}
	}
}
