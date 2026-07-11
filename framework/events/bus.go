// Package events provides a Watermill-based event bus for domain events
// with GoChannel pub/sub and webhook dispatching.
package events

import (
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

var (
	// Bus is the global Watermill GoChannel pub/sub instance.
	Bus *gochannel.GoChannel
)

// Init initializes the global event bus. Call once during app startup.
func Init() {
	Bus = gochannel.NewGoChannel(
		gochannel.Config{
			OutputChannelBuffer: 1024,
		},
		watermill.NewStdLogger(false, false),
	)
	log.Println("[events] Watermill GoChannel bus initialized")
}
