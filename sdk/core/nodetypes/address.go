package nodetypes

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/qubic/go-qubic/v2/common"
	"github.com/qubic/go-qubic/v2/connector"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
)

const (
	SpectrumDepth = 24
)

const (
	BalanceTypeRequest  = 31
	BalanceTypeResponse = 32
)

type AddressData struct {
	PublicKey                  [32]byte
	IncomingAmount             int64
	OutgoingAmount             int64
	NumberOfIncomingTransfers  uint32
	NumberOfOutgoingTransfers  uint32
	LatestIncomingTransferTick uint32
	LatestOutgoingTransferTick uint32
}

type AddressInfo struct {
	AddressData   AddressData
	Tick          uint32
	SpectrumIndex int32
	Siblings      [SpectrumDepth][32]byte
}

func (ai *AddressInfo) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return fmt.Errorf("reading header: %w", err)
	}

	if header.Type != BalanceTypeResponse {
		return fmt.Errorf("Invalid header type, expected %d, found %d", BalanceTypeResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, ai)
	if err != nil {
		return fmt.Errorf("reading addr info data from reader: %w", err)
	}

	return nil
}

func (ai *AddressInfo) ToProto() (*qubicpb.EntityInfo, error) {
	aic := addressInfoConverter{rawAddressInfo: *ai}
	aiPb, err := aic.toProto()
	if err != nil {
		return nil, fmt.Errorf("calling address info converter to proto: %w", err)
	}

	return aiPb, nil
}

type addressInfoConverter struct {
	rawAddressInfo AddressInfo
}

func (aic addressInfoConverter) toProto() (*qubicpb.EntityInfo, error) {
	id, err := common.PubKeyToIdentity(aic.rawAddressInfo.AddressData.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("getting address id from pubkey hex: %s: %w", hex.EncodeToString(aic.rawAddressInfo.AddressData.PublicKey[:]), err)
	}

	siblings := make([]string, 0, SpectrumDepth)
	for _, sibling := range aic.rawAddressInfo.Siblings {
		if sibling == [32]byte{} {
			continue
		}
		siblingID, err := common.PubKeyToIdentity(sibling)
		if err != nil {
			return nil, fmt.Errorf("getting address id from sibling hex: %s: %w", hex.EncodeToString(sibling[:]), err)
		}
		siblings = append(siblings, siblingID.String())
	}

	return &qubicpb.EntityInfo{
		Entity: &qubicpb.EntityInfo_Entity{
			Id:                         id.String(),
			IncomingAmount:             aic.rawAddressInfo.AddressData.IncomingAmount,
			OutgoingAmount:             aic.rawAddressInfo.AddressData.OutgoingAmount,
			NumberOfIncomingTransfers:  aic.rawAddressInfo.AddressData.NumberOfIncomingTransfers,
			NumberOfOutgoingTransfers:  aic.rawAddressInfo.AddressData.NumberOfOutgoingTransfers,
			LatestIncomingTransferTick: aic.rawAddressInfo.AddressData.LatestIncomingTransferTick,
			LatestOutgoingTransferTick: aic.rawAddressInfo.AddressData.LatestOutgoingTransferTick,
		},
		ValidForTick:  aic.rawAddressInfo.Tick,
		SpectrumIndex: aic.rawAddressInfo.SpectrumIndex,
		SiblingIds:    siblings,
	}, nil
}
