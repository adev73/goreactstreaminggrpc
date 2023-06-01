package greetingclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adev73/goreactstreaminggrpc/internal/gen/greet/v1/greetv1connect"

	"github.com/bufbuild/connect-go"

	greetv1 "github.com/adev73/goreactstreaminggrpc/internal/gen/greet/v1"
)

type GreetingClient interface {
	Disconnect() error
	Greet(name string) (bool, error)
}

type greetingClient struct {
	client    greetv1connect.GreetServiceClient
	SessionId string
	stream    *connect.ServerStreamForClient[greetv1.GreetingsResponse]
	Greetings chan Greeting
	Errors    chan error
}

type Greeting struct {
	Greeting string
}

func New(address string) (*greetingClient, error) {

	gc := &greetingClient{}

	gc.client = greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)

	stream, err := gc.client.Greetings(
		context.Background(),
		&connect.Request[greetv1.GreetingsRequest]{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect stream with error: %w", err)
	}

	gc.stream = stream

	if gc.stream.Receive() {
		// First message is session ID
		gc.SessionId = stream.Msg().SessionId
	} else {
		return nil, fmt.Errorf("did not receive session id from stream, got: %w", stream.Err())
	}

	gc.Greetings = make(chan Greeting, 1) // TODO: Make greeting channel buffer configurable
	gc.Errors = make(chan error, 1)
	go gc.receive()

	return gc, nil
}
func (gc *greetingClient) Disconnect() error {

	// Ignore errors & response, just close it.
	gc.client.Greet(context.Background(), &connect.Request[greetv1.GreetRequest]{
		Msg: &greetv1.GreetRequest{
			SessionId:  gc.SessionId,
			EndSession: true,
		},
	})
	gc.stream.Close()

	return nil
}

func (gc *greetingClient) receive() {

	for {
		ok := gc.stream.Receive()
		if ok {
			if gc.stream.Msg().EndSession {
				log.Printf("End of session %s.\n", gc.SessionId)
				gc.Errors <- nil
				close(gc.Greetings)
				close(gc.Errors)
				return
			}
			// Empty string is just the service pinging us
			if gc.stream.Msg().Greeting != "" {
				gc.Greetings <- Greeting{Greeting: gc.stream.Msg().Greeting}
				log.Printf("Received greeting at %d\n", time.Now().UnixNano())
			} else {
				log.Println("Pinged...")
			}
		}
		if gc.stream.Err() != nil {
			log.Printf("Stream error: %s\n", gc.stream.Err().Error())
			gc.Errors <- gc.stream.Err()
			close(gc.Greetings)
			close(gc.Errors)
			return
		}
	}
}

func (gc *greetingClient) Greet(name string) (bool, error) {

	ok := false
	response, err := gc.client.Greet(context.Background(), &connect.Request[greetv1.GreetRequest]{
		Msg: &greetv1.GreetRequest{
			SessionId: gc.SessionId,
			Name:      name,
		},
	})

	if response != nil && response.Msg != nil {
		ok = response.Msg.Confirmed
	}

	return ok, err

}
