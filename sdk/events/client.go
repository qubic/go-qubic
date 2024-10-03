package events

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/connector"
	qubicpb "github.com/qubic/go-qubic/proto/v1"
	"github.com/qubic/go-qubic/sdk/core"
	"time"
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

func (c *Client) GetRangeEvents(ctx context.Context, passcode [4]uint64, fromEventID, toEventID uint64) (*Events, error) {
	request := struct {
		Passcode    [4]uint64
		FromEventID uint64
		ToEventID   uint64
	}{
		Passcode:    passcode,
		FromEventID: fromEventID,
		ToEventID:   toEventID,
	}

	result := Events{Count: int64(toEventID-fromEventID) + 1}
	err := c.connector.PerformCoreRequest(ctx, EventTypeRequest, request, &result)
	if err != nil {
		return nil, errors.Wrap(err, "performing core request")
	}

	return &result, nil
}

func (c *Client) GetTickEventsChunk(ctx context.Context, passcode [4]uint64, tickNumber uint32) (*qubicpb.TickEvents, error) {
	coreClient := core.NewClient(c.connector)

	td, err := coreClient.GetTickData(ctx, tickNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting tick data")
	}

	txEvents := make([]*qubicpb.TransactionEvents, 0, len(td.TransactionIds))

	for txIndex, txID := range td.TransactionIds {
		idRange, err := c.GetTickTransactionEventsRange(ctx, passcode, tickNumber, uint32(txIndex))
		if err != nil {
			return nil, errors.Wrapf(err, "getting tick transaction events range for txIndex: %d", txIndex)
		}

		if idRange.FromEventID == -1 || idRange.NumberOfEvents == -1 {
			continue
		}

		evs, err := c.GetRangeEvents(ctx, passcode, uint64(idRange.FromEventID), uint64(idRange.FromEventID+idRange.NumberOfEvents))
		if err != nil {
			return nil, errors.Wrap(err, "getting events")
		}

		events := make([]*qubicpb.Event, 0, len(evs.Items))
		for _, ev := range evs.Items {
			protoEvent := EventConverter.ToProto(ev)
			events = append(events, protoEvent)
		}

		txEvent := qubicpb.TransactionEvents{
			TxId:   txID,
			Events: events,
		}

		txEvents = append(txEvents, &txEvent)
	}

	return &qubicpb.TickEvents{Tick: tickNumber, TxEvents: txEvents}, nil
}

func (c *Client) GetTickEventsOneByOne(ctx context.Context, passcode [4]uint64, tickNumber uint32) (*qubicpb.TickEvents, error) {
	coreClient := core.NewClient(c.connector)

	td, err := coreClient.GetTickData(ctx, tickNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting tick data")
	}

	txEvents := make([]*qubicpb.TransactionEvents, 0, len(td.TransactionIds))

	for txIndex, txID := range td.TransactionIds {
		idRange, err := c.GetTickTransactionEventsRange(ctx, passcode, tickNumber, uint32(txIndex))
		if err != nil {
			return nil, errors.Wrapf(err, "getting tick transaction events range for txIndex: %d", txIndex)
		}

		if idRange.FromEventID == -1 || idRange.NumberOfEvents == -1 {
			continue
		}

		from := uint64(idRange.FromEventID)
		to := uint64(idRange.FromEventID + idRange.NumberOfEvents)

		events := make([]*qubicpb.Event, 0, idRange.NumberOfEvents)

		for i := from; i < to; i++ {
			evs, err := func(eventID uint64) (*Events, error) {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()

				evs, err := c.GetRangeEvents(ctx, passcode, eventID, eventID)
				if err != nil {
					return nil, errors.Wrapf(err, "getting events for txIndex: %d, from event id: %d, to event id: %d", txIndex, from, to)
				}

				return evs, nil
			}(i)
			if err != nil {
				return nil, errors.Wrapf(err, "getting events for txIndex: %d, event id: %d", txIndex, i)
			}

			for _, ev := range evs.Items {
				protoEvent := EventConverter.ToProto(ev)
				events = append(events, protoEvent)
			}
		}

		txEvent := qubicpb.TransactionEvents{
			TxId:   txID,
			Events: events,
		}

		txEvents = append(txEvents, &txEvent)
	}

	return &qubicpb.TickEvents{Tick: tickNumber, TxEvents: txEvents}, nil
}
