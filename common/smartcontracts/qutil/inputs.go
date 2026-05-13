package qutil

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/v2/common"
)

const (
	SmartContractID      = 4
	SmartContractAddress = "EAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAVWRF"
)

type InputType int

const (
	SendToManyV1InputType InputType = iota + 1
	BurnQubicInputType
)

type FunctionType int

const (
	GetSendToManyV1FeeFunctionType FunctionType = iota + 1
)

const (
	sendManyV1MaxTransfers = 25
)

type sendToManyV1Payload struct {
	addresses [sendManyV1MaxTransfers][32]byte
	amounts   [sendManyV1MaxTransfers]int64
}

func (stm *sendToManyV1Payload) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	err := binary.Read(r, binary.LittleEndian, &stm.addresses)
	if err != nil {
		return errors.Wrap(err, "reading addresses from reader")
	}

	err = binary.Read(r, binary.LittleEndian, &stm.amounts)
	if err != nil {
		return errors.Wrap(err, "reading amounts from reader")
	}

	return nil
}

func (stm *sendToManyV1Payload) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, stm.addresses)
	if err != nil {
		return nil, errors.Wrap(err, "writing addresses to buffer")
	}

	err = binary.Write(&buf, binary.LittleEndian, stm.amounts)
	if err != nil {
		return nil, errors.Wrap(err, "writing amounts to buffer")
	}

	return buf.Bytes(), nil
}

type SendToManyV1Transfers struct {
	Transfers   []common.QubicTransfer
	TotalAmount int64
}

func (stm *SendToManyV1Transfers) MarshalBinary() ([]byte, error) {
	var payload sendToManyV1Payload
	for i, transfer := range stm.Transfers {
		addressPubkey, err := transfer.DestinationID.ToPubKey(false)
		if err != nil {
			return nil, errors.Wrapf(err, "converting destination ID: %s to pubkey", transfer.DestinationID)
		}

		payload.addresses[i] = addressPubkey
		payload.amounts[i] = transfer.Amount
	}

	binaryData, err := payload.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "binary marshalling payload")
	}

	return binaryData, nil
}

func (stm *SendToManyV1Transfers) FromBase64TxInput(input string) error {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return errors.Wrap(err, "decoding base64 tx input")
	}

	err = stm.FromRawInput(data)
	if err != nil {
		return errors.Wrap(err, "parsing tx input")
	}

	return nil
}

func (stm *SendToManyV1Transfers) FromRawInput(input []byte) error {
	var payload sendToManyV1Payload
	err := payload.UnmarshalBinary(input)
	if err != nil {
		return errors.Wrap(err, "binary unmarshalling payload")
	}

	transfers := make([]common.QubicTransfer, 0, len(payload.addresses))
	var totalAmount int64
	for i, address := range payload.addresses {
		if address == [32]byte{} {
			continue
		}

		addressID, err := common.PubKeyToIdentity(address)
		if err != nil {
			return errors.Wrapf(err, "converting address: %s", hex.EncodeToString(address[:]))
		}
		amount := payload.amounts[i]
		transfers = append(transfers, common.QubicTransfer{
			DestinationID: addressID,
			Amount:        amount,
		})
		totalAmount += amount
	}

	stm.Transfers = transfers
	stm.TotalAmount = totalAmount

	return nil
}

func (stm *SendToManyV1Transfers) FromHexTxInput(input string) error {
	data, err := hex.DecodeString(input)
	if err != nil {
		return errors.Wrap(err, "decoding hex tx input")
	}

	err = stm.FromRawInput(data)
	if err != nil {
		return errors.Wrap(err, "parsing tx input")
	}

	return nil
}
