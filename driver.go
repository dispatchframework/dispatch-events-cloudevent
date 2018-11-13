package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudevents/sdk-go/v02"
	"github.com/vmware/dispatch/pkg/events"
	"github.com/vmware/dispatch/pkg/events/driverclient"
)

type validationEvent struct {
	Data struct {
		ValidationCode string `json:"validationCode"`
		ValidationURL  string `json:"validationUrl"`
	} `json:"data"`
	EventType string `json:"eventType"`
	Topic     string `json:"topic"`
}

type validationResponse struct {
	ValidationResponse string `json:"validationResponse"`
}

// debug
var dryRun = flag.Bool("dry-run", false, "Debug, pull messages and do not send Dispatch events")
var org = flag.String("org", "default", "organization of this event driver")
var dispatchEndpoint = flag.String("dispatch-api-endpoint", "localhost:8080", "dispatch server host")
var port = flag.Int("port", 80, "Port to listen on")
var sharedSecret = flag.String("shared-secret", "", "A token or shared secret that the client should pass")

func getDriverClient() driverclient.Client {
	if *dryRun {
		return nil
	}
	token := os.Getenv(driverclient.AuthToken)
	client, err := driverclient.NewHTTPClient(driverclient.WithGateway(*dispatchEndpoint), driverclient.WithToken(token))
	if err != nil {
		log.Fatalf("Error when creating the events client: %s", err.Error())
	}
	log.Println("Event driver initialized.")
	return client
}

func main() {

	flag.Parse()

	client := getDriverClient()
	marshaller := v02.NewDefaultHTTPMarshaller()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		event, err := marshaller.FromRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		v02Event, ok := event.(*v02.Event)
		if !ok {
			log.Printf("wrong event type: %v", v02Event)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ceBytes, err := json.Marshal(v02Event.Data)
		if err != nil {
			log.Printf("failed to marshal event data: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		callbackURL := r.Header.Get("X-Callback-URL")

		dispatchEvent := &events.CloudEvent{
			EventType:          v02Event.Type,
			CloudEventsVersion: "0.1",
			ContentType:        v02Event.ContentType,
			EventID:            v02Event.ID,
			Source:             v02Event.Source,
			Extensions: map[string]interface{}{
				"callback-url": callbackURL,
			},
			Data: json.RawMessage(ceBytes),
		}
		if client != nil {
			err := client.SendOne(dispatchEvent)
			if err != nil {
				log.Printf("Error sending event: %v", err)
				return
			}
		}
		pretty, _ := json.MarshalIndent(dispatchEvent, "", "  ")
		log.Printf("Sent event successfully: %s", string(pretty))
	})

	// Create chan signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	}()

	<-done
}
