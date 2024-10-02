package events

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/internal/connector"
)

type Client struct {
	connector *connector.Connector
}

func NewClient(connector *connector.Connector) *Client {
	return &Client{
		connector: connector,
	}
}

func (c *Client) GetTickTransactionEventsRange(ctx context.Context, passcode [4]uint64, tickNumber, txIndex uint32) (*TransactionEventsRange, error) {
	request := struct {
		Passcode   [4]uint64
		TickNumber uint32
		TxIndex    uint32
	}{
		Passcode:   passcode,
		TickNumber: tickNumber,
		TxIndex:    txIndex,
	}

	var result TransactionEventsRange
	err := c.connector.PerformCoreRequest(ctx, TransactionEventsRangeTypeRequest, request, &result)
	if err != nil {
		return nil, errors.Wrap(err, "performing core request")
	}

	return &result, nil
}

func (c *Client) GetEvents(ctx context.Context, passcode [4]uint64, fromEventID, toEventID uint64) (*Events, error) {
	request := struct {
		Passcode    [4]uint64
		FromEventID uint64
		ToEventID   uint64
	}{
		Passcode:    passcode,
		FromEventID: fromEventID,
		ToEventID:   toEventID,
	}

	result := Events{Count: int64(toEventID - fromEventID)}
	err := c.connector.PerformCoreRequest(ctx, EventTypeRequest, request, &result)
	if err != nil {
		return nil, errors.Wrap(err, "performing core request")
	}

	return &result, nil
}
