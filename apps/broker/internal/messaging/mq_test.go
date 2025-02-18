package messaging

import (
	"fmt"
	"github.com/nats-io/nats-server/v2/server"
	"testing"
	"time"
)

func TestNATS(t *testing.T) {
	opts := &server.Options{}

	// Initialize new server with options
	ns, err := server.NewServer(opts)

	if err != nil {
		panic(err)
	}

	// Start the server via goroutine
	go ns.Start()

	// Wait for server to be ready for connections
	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}

	// Connect to server
	nc, err := nats.Connect(ns.ClientURL())

	if err != nil {
		panic(err)
	}

	subject := "my-subject"

	// Subscribe to the subject
	nc.Subscribe(subject, func(msg *nats.Msg) {
		// Print message data
		data := string(msg.Data)
		fmt.Println(data)

		// Shutdown the server (optional)
		ns.Shutdown()
	})

	// Publish data to the subject
	nc.Publish(subject, []byte("Hello embedded NATS!"))
	testRequestReply(t, nc)
	// Wait for server shutdown
	ns.WaitForShutdown()
}

func testRequestReply(t *testing.T, nc *nats.Conn) {
	subject := "request.reply.subject"
	msg := "Request message"

	// Set up a subscriber for the request
	nc.Subscribe(subject, func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte("Reply message"))
	})

	// Send a request and wait for a response
	resp, err := nc.Request(subject, []byte(msg), 10*time.Second)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if string(resp.Data) != "Reply message" {
		t.Errorf("Expected 'Reply message', got %s", string(resp.Data))
	}
}
