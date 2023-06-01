package main

import (
	greetingclient "github.com/adev73/goreactstreaminggrpc/internal/greeting-client"
	"log"
	"time"
)

var greetings = []string{
	"Jane", "Jim", "Janet", "John",
}

func main() {

	g, err := greetingclient.New("localhost:8080")
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Starting greetings (session %s)...\n", g.SessionId)

	// Receiver function
	signal := make(chan bool)
	go func() {
		for {
			select {
			case <-signal:
				// Finished
				return
			case greetme := <-g.Greetings:
				log.Printf("Greeting received: %s\n", greetme.Greeting)
			case greeterr := <-g.Errors:
				log.Printf("Error received from greeter service. Panic? : %s", greeterr.Error())
			}
		}
	}()

	// Sender function
	for _, name := range greetings {
		ok, err := g.Greet(name)
		if err != nil {
			log.Printf("Failed to greet %s with error %s\n", name, err.Error())
			return
		}
		if !ok {
			log.Printf("Unable to greet %s with no error, continuing\n", name)
		}
		time.Sleep(15 * time.Second)
	}

	log.Println("Done greeting people, closing stream...")
	signal <- true // Tell goroutine to quit.
	g.Disconnect()
	log.Println("Done!")

}
