package events

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/qubic/go-qubic/v2/connector"
	"github.com/qubic/go-qubic/v2/sdk/core/nodetypes"
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

	EventTypeContractReserveDeduction = 13
	EventTypeOracleQueryStatusChange  = 14
)

const EventTypeCustomMessage = 255

const (
	EventTypeRequest                     = 44
	EventTypeResponse                    = 45
	TransactionEventsRangeTypeRequest    = 48
	TransactionEventsRangeTypeResponse   = 49
	TickTransactionEventsIDsTypeRequest  = 50
	TickTransactionEventsIDsTypeResponse = 51
)

type TransactionEventsRange struct {
	FromEventID    int64
	NumberOfEvents int64
}

func (ter *TransactionEventsRange) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return fmt.Errorf("reading tick data from reader: %w", err)
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	packetSize := header.GetSize()
	_ = packetSize

	headerSize := binary.Size(header)
	_ = headerSize

	if header.Type != TransactionEventsRangeTypeResponse {
		return fmt.Errorf("Invalid header type, expected %d, found %d", TransactionEventsRangeTypeResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, ter)
	if err != nil {
		return fmt.Errorf("reading transaction events range from reader: %w", err)
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
		return fmt.Errorf("reading event header: %w", err)
	}

	ev.EventType = uint8(ev.Header.Tmp >> 24)
	ev.EventSize = (ev.Header.Tmp << 8) >> 8

	eventData := make([]byte, ev.EventSize)
	err = binary.Read(r, binary.LittleEndian, eventData)
	if err != nil {
		return fmt.Errorf("reading event data: %w", err)
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
		return fmt.Errorf("reading header: %w", err)
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != EventTypeResponse {
		return fmt.Errorf("Invalid header type, expected %d, found %d", EventTypeResponse, header.Type)
	}
	items := make([]Event, 0, evs.Count)
	for range evs.Count {
		var ev Event
		err = ev.UnmarshalFromReader(r)
		if err != nil {
			return fmt.Errorf("unmarshalling event: %w", err)
		}
		items = append(items, ev)
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
		return fmt.Errorf("reading qu transfer event: %w", err)
	}

	return nil
}

type AssetIssuanceEvent struct {
	SourceIdentityPubKey  [32]byte
	NumberOfShares        int64
	ManagingContractIndex int64
	AssetName             [7]byte
	NumberOfDecimals      uint8
	MeasurementUnit       [7]byte
}

func (e *AssetIssuanceEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return fmt.Errorf("reading asset issuance event: %w", err)
	}

	return nil
}

type AssetOwnershipChangeEvent struct {
	SourceIdentityPubKey      [32]byte
	DestinationIdentityPubKey [32]byte
	IssuerIdentityPubKey      [32]byte
	NumberOfShares            int64
	ManagingContractIndex     int64
	AssetName                 [7]byte
	NumberOfDecimals          uint8
	MeasurementUnit           [7]byte
}

func (e *AssetOwnershipChangeEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return fmt.Errorf("reading asset ownership change event: %w", err)
	}

	return nil
}

type AssetPossessionChangeEvent struct {
	SourceIdentityPubKey      [32]byte
	DestinationIdentityPubKey [32]byte
	IssuerIdentityPubKey      [32]byte
	NumberOfShares            int64
	ManagingContractIndex     int64
	AssetName                 [7]byte
	NumberOfDecimals          uint8
	MeasurementUnit           [7]byte
}

func (e *AssetPossessionChangeEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return fmt.Errorf("reading asset possession change event: %w", err)
	}

	return nil
}

type BurningEvent struct {
	SourceIdentityPubKey [32]byte
	Amount               uint64
}

func (e *BurningEvent) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return fmt.Errorf("reading burning event: %w", err)
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
		return fmt.Errorf("reading dust burning event: %w", err)
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
		return fmt.Errorf("reading spectrum stats event: %w", err)
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
		return fmt.Errorf("reading contract id: %w", err)
	}

	e.Message = make([]byte, len(data)-4)
	err = binary.Read(r, binary.LittleEndian, &e.Message)
	if err != nil {
		return fmt.Errorf("reading contract message: %w", err)
	}

	return nil
}

type ContractReserveDeductionEvent struct {
	DeductedAmount  uint64
	RemainingAmount int64
	ContractIndex   uint32
	_               uint32 // padding to match 24-byte C++ struct
}

func (e *ContractReserveDeductionEvent) UnmarshalBinary(data []byte) error {
	if len(data) != 24 {
		return fmt.Errorf("invalid contract reserve deduction event size: expected 24, got %d", len(data))
	}

	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return fmt.Errorf("reading contract reserve deduction event: %w", err)
	}

	return nil
}

type OracleQueryStatusChangeEvent struct {
	QueryingEntity [32]byte
	QueryID        int64
	InterfaceIndex uint32
	Type           uint8
	Status         uint8
}

// OracleQueryStatus constants
const (
	OracleQueryStatusPending      = 1
	OracleQueryStatusCommitted    = 2
	OracleQueryStatusSuccess      = 3
	OracleQueryStatusTimeout      = 4
	OracleQueryStatusUnresolvable = 5

	OracleQueryTypeContractQuery        = 0
	OracleQueryTypeContractSubscription = 1
	OracleQueryTypeUserQuery            = 2
)

func (e *OracleQueryStatusChangeEvent) UnmarshalBinary(data []byte) error {
	if len(data) != 46 {
		return fmt.Errorf("invalid oracle query status change event size: expected 46, got %d", len(data))
	}

	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return fmt.Errorf("reading oracle query status change event: %w", err)
	}

	return nil
}

func (e *OracleQueryStatusChangeEvent) StatusString() string {
	switch e.Status {
	case OracleQueryStatusPending:
		return "pending"
	case OracleQueryStatusCommitted:
		return "committed"
	case OracleQueryStatusSuccess:
		return "success"
	case OracleQueryStatusTimeout:
		return "timeout"
	case OracleQueryStatusUnresolvable:
		return "unresolvable"
	default:
		return "unknown"
	}
}

func (e *OracleQueryStatusChangeEvent) TypeString() string {
	switch e.Type {
	case OracleQueryTypeContractQuery:
		return "contract_query"
	case OracleQueryTypeContractSubscription:
		return "contract_subscription"
	case OracleQueryTypeUserQuery:
		return "user_query"
	default:
		return "unknown"
	}
}

const MaxNumberOfSpecialEventsPerTick = 6

type TickTransactionEventIDs struct {
	FromEventID [nodetypes.MaxNumberOfTransactionsPerTick + MaxNumberOfSpecialEventsPerTick]int64
	Length      [nodetypes.MaxNumberOfTransactionsPerTick + MaxNumberOfSpecialEventsPerTick]int64
}

func (e *TickTransactionEventIDs) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return fmt.Errorf("reading header from reader: %w", err)
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != TickTransactionEventsIDsTypeResponse {
		return fmt.Errorf("Invalid header type, expected %d, found %d", TickTransactionEventsIDsTypeResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, e)
	if err != nil {
		return fmt.Errorf("reading tick transaction event ids from reader: %w", err)
	}

	return nil
}
