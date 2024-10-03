package nodetypes

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/common"
	"github.com/qubic/go-qubic/connector"
	qubicpb "github.com/qubic/go-qubic/proto/v1"
	"io"
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
		return errors.Wrap(err, "reading header")
	}

	if header.Type != BalanceTypeResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", BalanceTypeResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, ai)
	if err != nil {
		return errors.Wrap(err, "reading addr info data from reader")
	}

	return nil
}

func (ai *AddressInfo) ToProto() (*qubicpb.EntityInfo, error) {
	aic := addressInfoConverter{rawAddressInfo: *ai}
	aiPb, err := aic.toProto()
	if err != nil {
		return nil, errors.Wrap(err, "calling address info converter to proto")
	}

	return aiPb, nil
}

type addressInfoConverter struct {
	rawAddressInfo AddressInfo
}

func (aic addressInfoConverter) toProto() (*qubicpb.EntityInfo, error) {
	id, err := common.PubKeyToIdentity(aic.rawAddressInfo.AddressData.PublicKey)
	if err != nil {
		return nil, errors.Wrapf(err, "getting address id from pubkey hex: %s", hex.EncodeToString(aic.rawAddressInfo.AddressData.PublicKey[:]))
	}

	siblings := make([]string, 0, SpectrumDepth)
	for _, sibling := range aic.rawAddressInfo.Siblings {
		if sibling == [32]byte{} {
			continue
		}
		siblingID, err := common.PubKeyToIdentity(sibling)
		if err != nil {
			return nil, errors.Wrapf(err, "getting address id from sibling hex: %s", hex.EncodeToString(sibling[:]))
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
