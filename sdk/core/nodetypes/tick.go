package nodetypes

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/qubic/go-qubic/v2/common"
	"github.com/qubic/go-qubic/v2/connector"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	MaxNumberOfTransactionsPerTick = 4096
)

const (
	CurrentTickInfoTypeRequest  = 27
	CurrentTickInfoTypeResponse = 28
	TickDataTypeResponse        = 8
	TickDataTypeRequest         = 16
)

type TickData struct {
	ComputorIndex      uint16
	Epoch              uint16
	Tick               uint32
	Millisecond        uint16
	Second             uint8
	Minute             uint8
	Hour               uint8
	Day                uint8
	Month              uint8
	Year               uint8
	Timelock           [32]byte
	TransactionDigests [MaxNumberOfTransactionsPerTick][32]byte
	ContractFees       [MaxNumberOfTransactionsPerTick]int64
	Signature          [SignatureSize]byte
}

func (td *TickData) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return fmt.Errorf("reading tick data from reader: %w", err)
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != TickDataTypeResponse {
		return fmt.Errorf("Invalid header type, expected %d, found %d", TickDataTypeResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, td)
	if err != nil {
		return fmt.Errorf("reading tick data from reader: %w", err)
	}

	return nil
}

func (td *TickData) IsEmpty() bool {
	if td == nil {
		return true
	}

	return *td == TickData{}
}

func (td *TickData) ToProto() (*qubicpb.TickData, error) {
	tdc := tickDataConverter{rawTd: *td}
	tdPb, err := tdc.toProto()
	if err != nil {
		return nil, fmt.Errorf("calling tick data converter to proto: %w", err)
	}

	return tdPb, nil
}

type tickDataConverter struct {
	rawTd TickData
}

func (tdc *tickDataConverter) toProto() (*qubicpb.TickData, error) {
	if tdc.rawTd.IsEmpty() {
		return &qubicpb.TickData{}, nil
	}

	date := time.Date(2000+int(tdc.rawTd.Year), time.Month(tdc.rawTd.Month), int(tdc.rawTd.Day), int(tdc.rawTd.Hour), int(tdc.rawTd.Minute), int(tdc.rawTd.Second), 0, time.UTC)
	date.Add(time.Duration(tdc.rawTd.Millisecond) * time.Millisecond)

	transactionIds, err := common.PubKeysToIdentitiesString(tdc.rawTd.TransactionDigests[:], true)
	if err != nil {
		return nil, fmt.Errorf("getting transaction ids from digests: %w", err)
	}

	return &qubicpb.TickData{
		ComputorIndex:  uint32(tdc.rawTd.ComputorIndex),
		Epoch:          uint32(tdc.rawTd.Epoch),
		Tick:           tdc.rawTd.Tick,
		Timestamp:      timestamppb.New(date),
		TimeLock:       base64.StdEncoding.EncodeToString(tdc.rawTd.Timelock[:]),
		TransactionIds: transactionIds,
		ContractFees:   contractFeesToProto(tdc.rawTd.ContractFees),
		Signature:      base64.StdEncoding.EncodeToString(tdc.rawTd.Signature[:]),
	}, nil
}

func contractFeesToProto(contractFees [MaxNumberOfTransactionsPerTick]int64) []int64 {
	protoContractFees := make([]int64, 0, len(contractFees))
	for _, fee := range contractFees {
		if fee == 0 {
			continue
		}
		protoContractFees = append(protoContractFees, fee)
	}
	return protoContractFees
}

type TickInfo struct {
	TickDuration            uint16
	Epoch                   uint16
	Tick                    uint32
	NumberOfAlignedVotes    uint16
	NumberOfMisalignedVotes uint16
	InitialTick             uint32
}

func (ti *TickInfo) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return fmt.Errorf("reading header: %w", err)
	}

	if header.Type != CurrentTickInfoTypeResponse {
		return fmt.Errorf("Invalid header type, expected %d, found %d", CurrentTickInfoTypeResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, ti)
	if err != nil {
		return fmt.Errorf("reading tick data from reader: %w", err)
	}
	return nil
}

func (ti *TickInfo) ToProto() (*qubicpb.TickInfo, error) {
	return &qubicpb.TickInfo{
		Tick:                    ti.Tick,
		DurationInSeconds:       uint32(ti.TickDuration),
		Epoch:                   uint32(ti.Epoch),
		NumberOfAlignedVotes:    uint32(ti.NumberOfAlignedVotes),
		NumberOfMisalignedVotes: uint32(ti.NumberOfMisalignedVotes),
		InitialTickOfEpoch:      ti.InitialTick,
	}, nil
}
