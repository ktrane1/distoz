package main

import (
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/ktalele/distoz/broadcast"
	"github.com/ktalele/distoz/echo"
	"github.com/ktalele/distoz/generate"
)

func main() {
	n := maelstrom.NewNode()

	echo.Echo(n)
	generate.Generate(n)
	broadcast.Broadcast(n)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
