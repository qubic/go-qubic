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
		ConnectionTimeout:     time.Minute,
		HandlerRequestTimeout: time.Minute,
	}
	requestPerformer, err := connector.NewConnector("1.2.3.4", connectorConfig)
	assert.NoError(t, err)

	passcodes := map[string][4]uint64{
		"1.2.3.4": {1, 2, 3, 4},
	}

	eventClient := NewClient(requestPerformer, passcodes)

	//tickEvents, err := eventClient.GetTickEventsOneByOne(context.Background(), 20542764)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tickEvents, err := eventClient.GetTickEvents(ctx, 20542765)
	assert.NoError(t, err)
	log.Print(tickEvents.Tick)

}
