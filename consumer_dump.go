package main

import (
	"context"
	"fmt"
	"os"

	_ "math/bits"
	"time"
	"pack.ag/amqp"
)

func failOnError(err error, message string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%s: %s", message, err))
		os.Exit(1)
	}
}

func main (){
/* Set custom env variable
	os.Setenv("username", "XYZ")
	os.Setenv("password", "ABC")
	os.Setenv("host", "DNS")
	os.Setenv("port", "NO")
	var username = os.Getenv("username")
	var password = os.Getenv("password")
	var host = os.Getenv("host")
	var port = os.Getenv("port")
*/
	// Create client
	client, err := amqp.Dial("amqps://"+"b-4ec445dd-e496-4f5b-95ed-416d1b391e00-1.mq.us-west-2.amazonaws.com"+":"+"5671",
		amqp.ConnSASLPlain("cfr5BQa5UXzM", "q0NkN4lxeOB8S0vNJFYh"),
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("Dialing AMQP server:", err))
	}
	defer client.Close()

	// Open a session
	session, err := client.NewSession()
	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot create session, error:", err))
	}
	ctx := context.Background()

	queue := "/test-queue"
//	Create receiver and receive mesages

	// Create a receiver
	receiver, err := session.NewReceiver(
		amqp.LinkSourceAddress(queue),
		amqp.LinkCredit(10),
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("Creating receiver link, error:", err))
	}

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 35*time.Second)
		receiver.Close(ctx)
		cancel()
	}()

	for {
		// Receive next message
		msg, err := receiver.Receive(ctx)
		if err != nil {
			fmt.Println(fmt.Sprintf("Reading message from AMQP error:", err))
		}

		// Accept message
		msg.Accept()

		fmt.Println(fmt.Sprintf("Message received: %s\n", msg.GetData()))
	}
}
