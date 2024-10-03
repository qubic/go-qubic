package events

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/connector"
	"io"
)

const (
	EventHeaderSize = 26
)

const (
	EventTypeQuTransfer = uint8(iota)
	EventTypeAssetIssuance
	EventTypeAssetOwnershipChange
	EventTypeAssetPossessionChange
	EventTypeContractErrorMessage
	EventTypeContractWarningMessage
	EventTypeContractInformationMessage
	EventTypeContractDebugMessage
	EventTypeBurning
	EventTypeDustBurning
	EventTypeSpectrumStats
)

const EventTypeCustomMessage = 255

const (
	EventTypeRequest                   = 44
	EventTypeResponse                  = 45
	TransactionEventsRangeTypeRequest  = 48
	TransactionEventsRangeTypeResponse = 49
)

type TransactionEventsRange struct {
	FromEventID    int64
	NumberOfEvents int64
}

func (ter *TransactionEventsRange) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return errors.Wrap(err, "reading tick data from reader")
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	packetSize := header.GetSize()
	_ = packetSize

	headerSize := binary.Size(header)
	_ = headerSize

	if header.Type != TransactionEventsRangeTypeResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", TransactionEventsRangeTypeResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, ter)
	if err != nil {
		return errors.Wrap(err, "reading transaction events range from reader")
	}

	return nil
}

type Event struct {
	Header    Header
	EventType uint8
	EventSize uint32
	Data      []byte
}

type Header struct {
	Epoch       uint16
	Tick        uint32
	Tmp         uint32
	EventID     uint64
	EventDigest [8]byte
}

func (ev *Event) UnmarshalFromReader(r io.Reader) error {
	err := binary.Read(r, binary.LittleEndian, &ev.Header)
	if err != nil {
		return errors.Wrap(err, "reading event header")
	}

	ev.EventType = uint8(ev.Header.Tmp >> 24)
	ev.EventSize = (ev.Header.Tmp << 8) >> 8

	eventData := make([]byte, ev.EventSize)
	err = binary.Read(r, binary.LittleEndian, eventData)
	if err != nil {
		return errors.Wrap(err, "reading event data")
	}

	ev.Data = eventData

	return nil
}

type Events struct {
	Items []Event
	Count int64
}

func (evs *Events) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader
	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return errors.Wrap(err, "reading header")
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != EventTypeResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", EventTypeResponse, header.Type)
	}
	items := make([]Event, 0, evs.Count)
	for range evs.Count {
		var ev Event
		err = ev.UnmarshalFromReader(r)
		if err != nil {
			return errors.Wrap(err, "unmarshalling event")
		}
		items = append(items, ev)
		//fmt.Printf("Got event id: %d\n", ev.Header.EventID)
	}

	evs.Items = items

	return nil
}

type QuTransferEvent struct {
	SourceIdentityPubKey      [32]byte
	DestinationIdentityPubKey [32]byte
	Amount                    uint64
}

func (e *QuTransferEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return errors.Wrap(err, "reading qu transfer event")
	}

	return nil
}

type AssetIssuanceEvent struct {
	SourceIdentityPubKey [32]byte
	AssetName            [8]byte
	NumberOfDecimals     uint8
	MeasurementUnit      [8]byte
	NumberOfShares       int64
}

func (e *AssetIssuanceEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return errors.Wrap(err, "reading asset issuance event")
	}

	return nil
}

type AssetOwnershipChangeEvent struct {
	SourceIdentityPubKey      [32]byte
	DestinationIdentityPubKey [32]byte
	IssuerIdentityPubKey      [32]byte
	AssetName                 [8]byte
	NumberOfDecimals          uint8
	MeasurementUnit           [8]byte
	NumberOfShares            int64
}

func (e *AssetOwnershipChangeEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return errors.Wrap(err, "reading asset ownership change event")
	}

	return nil
}

type AssetPossessionChangeEvent struct {
	SourceIdentityPubKey      [32]byte
	DestinationIdentityPubKey [32]byte
	IssuerIdentityPubKey      [32]byte
	AssetName                 [8]byte
	NumberOfDecimals          uint8
	MeasurementUnit           [8]byte
	NumberOfShares            int64
}

type BurningEvent struct {
	SourceIdentityPubKey [32]byte
	Amount               uint64
}

func (e *BurningEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return errors.Wrap(err, "reading burning event")
	}

	return nil
}

type DustBurningEvent struct {
	NumberOfBurns        uint16
	SourceIdentityPubKey [32]byte
	Amount               uint64
}

func (e *DustBurningEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return errors.Wrap(err, "reading dust burning event")
	}

	return nil
}

type SpectrumStatsEvent struct {
	TotalAmount               uint64
	DustThresholdBurnAll      uint64
	DustThresholdBurnHalf     uint64
	NumberOfEntities          uint32
	EntityCategoryPopulations [48]uint32
}

func (e *SpectrumStatsEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return errors.Wrap(err, "reading spectrum stats event")
	}

	return nil
}

type ContractMessageEvent struct {
	ContractID uint32
	Message    []byte
}

func (e *ContractMessageEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, &e.ContractID)
	if err != nil {
		return errors.Wrap(err, "reading contract id")
	}

	e.Message = make([]byte, len(data)-4)
	err = binary.Read(r, binary.LittleEndian, &e.Message)
	if err != nil {
		return errors.Wrap(err, "reading contract message")
	}

	return nil
}
