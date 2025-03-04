//go:build !ci
// +build !ci

package events

import (
	"context"
	"github.com/qubic/go-qubic/connector"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestClient_GetTickEvents(t *testing.T) {

	connectorConfig := connector.Config{
		ConnectionPort:        "21841",
		ConnectionTimeout:     10 * time.Second,
		HandlerRequestTimeout: 10 * time.Second,
	}
	requestPerformer, err := connector.NewConnector("1.2.3.4", connectorConfig)
	assert.NoError(t, err)

	passcodes := map[string][4]uint64{
		"1.2.3.4": {1, 2, 3, 4},
	}

	eventClient := NewClient(requestPerformer, passcodes)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tickEvents, err := eventClient.GetTickEvents(ctx, 20577267)
	//tickEvents, err := eventClient.GetTickEvents(ctx, 20542744)
	assert.NoError(t, err)
	log.Print(tickEvents.Tick)

}
