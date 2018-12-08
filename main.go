package main

import (
	"context"
	c "context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/iampigeon/pigeon"
	"github.com/iampigeon/pigeon/backend"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type service struct {
	Client *messaging.Client
	Ctx    c.Context
}

func (s *service) Approve(content []byte) (valid bool, err error) {
	if content == nil {
		return false, errors.New("Invalid message content")
	}

	fmt.Println(string(content))
	m := new(pigeon.Push)

	err = json.Unmarshal(content, m)
	if err != nil {
		return false, err
	}

	// validate topic to avoid  breakout
	fmt.Println(m)

	return true, nil
}

func (s *service) Deliver(content []byte) error {
	m := new(pigeon.Push)
	err := json.Unmarshal(content, m)
	if err != nil {
		return err
	}

	// push, err := json.Marshal(m)
	// if err != nil {
	// 	return err
	// }

	// FCM token
	registrationToken := m.Token

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: m.Title,
			Body:  m.Body,
		},
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided
	response, err := s.Client.Send(s.Ctx, message)
	if err != nil {
		log.Fatalln(err)
	}

	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	return nil
}

func main() {
	host := flag.String("host", "", "host of the service")
	port := flag.Int("port", 9030, "port of the service")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("Serving at %s", addr)

	opt := option.WithCredentialsFile("./keys.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	// Obtain a messaging.Client from the App.
	ctx := c.Background()
	client, err := app.Messaging(ctx)

	svc := &service{
		Client: client,
		Ctx:    ctx,
	}

	if err := backend.ListenAndServe(pigeon.NetAddr(addr), svc); err != nil {
		log.Fatal(err)
	}
}
