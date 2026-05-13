package events

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/qubic/go-qubic/v2/connector"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
	"github.com/qubic/go-qubic/v2/sdk/core"
)

type Client struct {
	connector  connector.RequestPerformer
	coreClient *core.Client
	passcodes  map[string][4]uint64
}

func NewClient(connector connector.RequestPerformer, passcodes map[string][4]uint64) *Client {
	return &Client{
		connector:  connector,
		coreClient: core.NewClient(connector),
		passcodes:  passcodes,
	}
}

type getTickTransactionsEventRangeRequest struct {
	Passcode   [4]uint64
	TickNumber uint32
	TxIndex    uint32
}

func (r *getTickTransactionsEventRangeRequest) AddPasscode(passcode [4]uint64) {
	r.Passcode = passcode
}

func (c *Client) GetTickTransactionEventsRange(ctx context.Context, tickNumber, txIndex uint32) (*TransactionEventsRange, error) {
	request := getTickTransactionsEventRangeRequest{
		TickNumber: tickNumber,
		TxIndex:    txIndex,
	}

	var result TransactionEventsRange
	err := c.connector.PerformCoreRequestWithPasscode(ctx, TransactionEventsRangeTypeRequest, c.passcodes, &request, &result)
	if err != nil {
		return nil, fmt.Errorf("performing core request: %w", err)
	}

	return &result, nil
}

type getRangeEventsRequest struct {
	Passcode    [4]uint64
	FromEventID uint64
	ToEventID   uint64
}

func (r *getRangeEventsRequest) AddPasscode(passcode [4]uint64) {
	r.Passcode = passcode
}

func (c *Client) GetRangeEvents(ctx context.Context, fromEventID, toEventID uint64) (*Events, error) {
	request := getRangeEventsRequest{
		FromEventID: fromEventID,
		ToEventID:   toEventID,
	}

	result := Events{Count: int64(toEventID-fromEventID) + 1}
	err := c.connector.PerformCoreRequestWithPasscode(ctx, EventTypeRequest, c.passcodes, &request, &result)
	if err != nil {
		return nil, fmt.Errorf("performing core request: %w", err)
	}

	return &result, nil
}

// GetTickEventsOneByOne Gets events per transaction index (one request per index).
//
// Deprecated: GetTickEventsOneByOne is deprecated as it is very inefficient. Use GetTickEvents instead.
func (c *Client) GetTickEventsOneByOne(ctx context.Context, tickNumber uint32) (*qubicpb.TickEvents, error) {
	td, err := c.coreClient.GetTickData(ctx, tickNumber)
	if err != nil {
		return nil, fmt.Errorf("getting tick data: %w", err)
	}

	if len(td.TransactionIds) == 0 {
		return &qubicpb.TickEvents{Tick: tickNumber, TxEvents: []*qubicpb.TransactionEvents{}}, nil
	}

	txEvents := make([]*qubicpb.TransactionEvents, 0, len(td.TransactionIds))

	for txIndex, txID := range td.TransactionIds {
		idRange, err := c.GetTickTransactionEventsRange(ctx, tickNumber, uint32(txIndex))
		if err != nil {
			return nil, fmt.Errorf("getting tick transaction events range for txIndex: %d: %w", txIndex, err)
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

				evs, err := c.GetRangeEvents(ctx, eventID, eventID)
				if err != nil {
					return nil, fmt.Errorf("getting events for txIndex: %d, from event id: %d, to event id: %d: %w", txIndex, from, to, err)
				}

				return evs, nil
			}(i)
			if err != nil {
				return nil, fmt.Errorf("getting events for txIndex: %d, event id: %d: %w", txIndex, i, err)
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

type getTickEventsRequest struct {
	Passcode   [4]uint64
	TickNumber uint32
}

func (r *getTickEventsRequest) AddPasscode(passcode [4]uint64) {
	r.Passcode = passcode
}

// GetTickEvents returns all events for a given tick number. This is not returning the special events (init_sc, begin_epoch, begin_tick, end_tick, end_epoch).
func (c *Client) GetTickEvents(ctx context.Context, tickNumber uint32) (*qubicpb.TickEvents, error) {
	td, err := c.coreClient.GetTickData(ctx, tickNumber)
	if err != nil {
		return nil, fmt.Errorf("getting tick data: %w", err)
	}

	if len(td.TransactionIds) == 0 {
		return &qubicpb.TickEvents{Tick: tickNumber, TxEvents: []*qubicpb.TransactionEvents{}}, nil
	}

	req := getTickEventsRequest{
		TickNumber: tickNumber,
	}

	var result TickTransactionEventIDs
	err = c.connector.PerformCoreRequestWithPasscode(ctx, TickTransactionEventsIDsTypeRequest, c.passcodes, &req, &result)
	if err != nil {
		return nil, fmt.Errorf("performing core request: %w", err)
	}

	var startEventId int64 = math.MaxInt64
	var endEventId int64

	txForEventID := make(map[int64]string)

	// this loop do not go over special events (init_sc, begin_epoch, begin_tick, end_tick, end_epoch which are starting at pos MaxNumberOfTransactionsPerTick)
	for i := range len(td.TransactionIds) {
		if result.FromEventID[i] == -1 {
			continue
		}

		if result.FromEventID[i] == -2 || result.FromEventID[i] == -3 {
			return nil, fmt.Errorf("From event id value %d inconsistent node", result.FromEventID[i])
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

	events, err := c.GetRangeEvents(ctx, uint64(startEventId), uint64(endEventId))
	if err != nil {
		return nil, fmt.Errorf("getting range events: %w", err)
	}

	eventsByTxID := make(map[string]*qubicpb.TransactionEvents)

	for _, ev := range events.Items {

		header := ev.Header
		if header.Tick != tickNumber {
			return nil, fmt.Errorf("received faulty data: event [%d] tick [%d] vs expected [%d]",
				header.EventID, header.Tick, tickNumber)
		}

		txID, ok := txForEventID[int64(header.EventID)]
		if !ok {
			log.Printf("Event with ID %d has no corresponding transaction ID\n", header.EventID)
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
