package events

import (
	"encoding/base64"
	"encoding/binary"
	qubicpb "github.com/qubic/go-qubic/proto/v1"
)

var EventConverter eventConverter

type eventConverter struct{}

func (ec *eventConverter) ToProto(event Event) *qubicpb.Event {
	return &qubicpb.Event{
		Header: &qubicpb.Event_Header{
			Epoch:       uint32(event.Header.Epoch),
			Tick:        event.Header.Tick,
			Tmp:         event.Header.Tmp,
			EventId:     event.Header.EventID,
			EventDigest: binary.LittleEndian.Uint64(event.Header.EventDigest[:]),
		},
		EventType: uint32(event.EventType),
		EventSize: event.EventSize,
		EventData: base64.StdEncoding.EncodeToString(event.Data),
	}
}
