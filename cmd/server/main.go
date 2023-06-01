package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	greetv1 "github.com/adev73/goreactstreaminggrpc/internal/gen/greet/v1"
	"github.com/adev73/goreactstreaminggrpc/internal/gen/greet/v1/greetv1connect"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	greeter := &GreetServer{}
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)
	mux.Handle(path, handler)
	log.Println("Starting server...")
	go DeadSessionMonitor(make(chan bool)) // Signal will never happen for now
	err := http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(
			newCORS().Handler(mux),
			&http2.Server{}),
	)
	if err != nil {
		log.Fatalf("server error: %s", err.Error())
	}
}

type GreetServer struct{}
type Greet struct {
	nameToGreet  string
	endOfSession bool
}

var sessions = make(map[string]chan Greet)

func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {

	sessionId := req.Msg.SessionId

	greeting := Greet{
		nameToGreet:  req.Msg.Name,
		endOfSession: req.Msg.EndSession,
	}
	if _, ok := sessions[sessionId]; !ok {
		return nil, fmt.Errorf("invalid session id %s", sessionId)
	}
	sessions[sessionId] <- greeting
	if greeting.endOfSession {
		close(sessions[sessionId])
	}
	return &connect.Response[greetv1.GreetResponse]{Msg: &greetv1.GreetResponse{Confirmed: true}}, nil
}

func (s *GreetServer) Greetings(ctx context.Context, req *connect.Request[greetv1.GreetingsRequest], stream *connect.ServerStream[greetv1.GreetingsResponse]) error {

	var err error
	var lastActivity time.Time = time.Now() // We can assume the connection is alive just now :D

	// Make a new session ID
	sessionId := uuid.New().String()

	// Send initial response to give client its session ID
	err = stream.Send(&greetv1.GreetingsResponse{
		SessionId:  sessionId,
		EndSession: false,
	})
	if err != nil {
		// Session start failed
		return err
	}

	log.Printf("Session %s started.\n", sessionId)
	sessions[sessionId] = make(chan Greet, 1)

	stopPing := make(chan bool)
	go func(signal chan bool) {
		// Ping the client every 10 seconds IF lastActivity > 10s
		for {
			select {
			case <-signal:
				// All done
				return

			case <-time.After(time.Second * 10):
				if time.Since(lastActivity) > time.Second*10 {
					log.Printf("PINGing session %s at %s\n", sessionId, time.Now().Format("04:05:06"))
					err := stream.Send(&greetv1.GreetingsResponse{
						SessionId: sessionId,
					})
					if err != nil {
						log.Printf("Session %s threw error %s\n", sessionId, err.Error())
						sessions[sessionId] <- Greet{endOfSession: true} // Put an end of session on our queue
					}

				}
			}
		}
	}(stopPing)

	for {
		if err = ctx.Err(); err != nil {
			log.Panicf("Context error: %s", err.Error())
			stopPing <- true
			delete(sessions, sessionId)
			break
		}

		nameToGreet := <-sessions[sessionId]

		if nameToGreet.endOfSession {
			v := greetv1.GreetingsResponse{
				SessionId:  sessionId,
				EndSession: true,
			}
			err = stream.Send(&v)
			stopPing <- true
			delete(sessions, sessionId)
			break
		}
		v := greetv1.GreetingsResponse{
			SessionId: sessionId,
			Greeting:  fmt.Sprintf("Hello, %s! It's %s", nameToGreet.nameToGreet, time.Now().Format("02 Jan 2006 3:04:05pm")),
		}
		err := stream.Send(&v)
		if err != nil {
			log.Printf("Session %s threw error %s\n", sessionId, err.Error())
		}
		lastActivity = time.Now()
	}

	if err != nil {
		log.Printf("Session %s ended with error: %s.\n", sessionId, err.Error())
	} else {
		log.Printf("Session %s completed.\n", sessionId)
	}

	return err
}

func DeadSessionMonitor(signal chan bool) {
	// Keep track of sessions. When we get a signal, stop
	for {
		log.Printf("There are %d sessions currently active.\n", len(sessions))
		select {
		case <-signal:
			return
		case <-time.After(10 * time.Second):
			// Continue
		}
	}
}

func newCORS() *cors.Cors {
	// To let web developers play with the demo service from browsers, we need a
	// very permissive CORS setup.
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
		// Let browsers cache CORS information for longer, which reduces the number
		// of preflight requests. Any changes to ExposedHeaders won't take effect
		// until the cached data expires. FF caps this value at 24h, and modern
		// Chrome caps it at 2h.
		MaxAge: int(2 * time.Hour / time.Second),
	})
}
