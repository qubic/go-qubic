package events

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/connector"
	qubicpb "github.com/qubic/go-qubic/proto/v1"
	"github.com/qubic/go-qubic/sdk/core"
	"log"
	"math"
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

func (c *Client) GetTickEventsOneByOne(ctx context.Context, passcode [4]uint64, tickNumber uint32) (*qubicpb.TickEvents, error) {
	coreClient := core.NewClient(c.connector)

	td, err := coreClient.GetTickData(ctx, tickNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting tick data")
	}

	if len(td.TransactionIds) == 0 {
		return &qubicpb.TickEvents{Tick: tickNumber, TxEvents: []*qubicpb.TransactionEvents{}}, nil
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

// GetTickEvents returns all events for a given tick number. This is not returning the special events (init_sc, begin_epoch, begin_tick, end_tick, end_epoch).
func (c *Client) GetTickEvents(ctx context.Context, passcode [4]uint64, tickNumber uint32) (*qubicpb.TickEvents, error) {
	coreClient := core.NewClient(c.connector)

	td, err := coreClient.GetTickData(ctx, tickNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting tick data")
	}

	if len(td.TransactionIds) == 0 {
		return &qubicpb.TickEvents{Tick: tickNumber, TxEvents: []*qubicpb.TransactionEvents{}}, nil
	}

	req := struct {
		Passcode   [4]uint64
		TickNumber uint32
	}{
		Passcode:   passcode,
		TickNumber: tickNumber,
	}

	var result TickTransactionEventIDs
	err = c.connector.PerformCoreRequest(ctx, TickTransactionEventsIDsTypeRequest, req, &result)
	if err != nil {
		return nil, errors.Wrap(err, "performing core request")
	}

	var startEventId int64 = math.MaxInt64
	var endEventId int64

	txForEventID := make(map[int64]string)

	// this loop do not go over special events (init_sc, begin_epoch, begin_tick, end_tick, end_epoch which are starting at pos 1024)
	for i := range len(td.TransactionIds) {
		if result.FromEventID[i] == -1 {
			continue
		}

		if result.FromEventID[i] == -2 || result.FromEventID[i] == -3 {
			return nil, errors.Errorf("From event id value %d inconsistent node", result.FromEventID[i])
		}

		addEventIDsToMap(txForEventID, result.FromEventID[i], result.Length[i], td.TransactionIds[i])

		if result.FromEventID[i] < startEventId {
			startEventId = result.FromEventID[i]
		}

		endEventId = result.FromEventID[i] + result.Length[i] - 1
	}

	if startEventId == math.MaxInt64 {
		return &qubicpb.TickEvents{Tick: tickNumber, TxEvents: []*qubicpb.TransactionEvents{}}, nil
	}

	events, err := c.GetRangeEvents(ctx, passcode, uint64(startEventId), uint64(endEventId))
	if err != nil {
		return nil, errors.Wrap(err, "getting range events")
	}

	eventsByTxID := make(map[string]*qubicpb.TransactionEvents)

	for _, ev := range events.Items {
		txID, ok := txForEventID[int64(ev.Header.EventID)]
		if !ok {
			log.Printf("Event with ID %d has no corresponding transaction ID\n", ev.Header.EventID)
		}

		txEvents, ok := eventsByTxID[txID]
		if !ok {
			txEvents = &qubicpb.TransactionEvents{
				TxId:   txID,
				Events: make([]*qubicpb.Event, 0),
			}
			eventsByTxID[txID] = txEvents
		}

		protoEvent := EventConverter.ToProto(ev)
		txEvents.Events = append(txEvents.Events, protoEvent)
		eventsByTxID[txID] = txEvents
	}

	txEvs := make([]*qubicpb.TransactionEvents, 0, len(eventsByTxID))

	for _, txEvents := range eventsByTxID {
		txEvs = append(txEvs, txEvents)
	}

	return &qubicpb.TickEvents{
		Tick:     tickNumber,
		TxEvents: txEvs,
	}, nil
}

func addEventIDsToMap(eventMap map[int64]string, fromEventID int64, length int64, txID string) {
	for i := fromEventID; i <= fromEventID+length; i++ {
		eventMap[i] = txID
	}
}
