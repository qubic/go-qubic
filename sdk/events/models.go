package events

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
)

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
	AssetName            uint64
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
	AssetName                 uint64
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
	AssetName                 uint64
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
