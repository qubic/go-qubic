package qutil

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/qubic/go-qubic/v2/common"
	"github.com/qubic/go-qubic/v2/sdk/core/nodetypes"
)

const Address = "EAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAVWRF"
const SendManyInputType = 1
const SendManyFee = 10
const SendManyInputSize = 1000
const SendManyMaxTransfers = 25

type SendManyTransferPayload struct {
	addresses       [SendManyMaxTransfers][32]byte
	amounts         [SendManyMaxTransfers]int64
	filledTransfers int8
	totalAmount     int64
}

type SendManyTransfer struct {
	AddressID common.Identity
	Amount    int64
}

func (smp *SendManyTransferPayload) AddTransfer(transfer SendManyTransfer) error {
	if smp.filledTransfers == SendManyMaxTransfers {
		return fmt.Errorf("max %d send many transfers allowed", SendManyMaxTransfers)
	}

	address, err := transfer.AddressID.ToPubKey(false)
	if err != nil {
		return fmt.Errorf("converting address id to byte form: %w", err)
	}

	smp.addresses[smp.filledTransfers] = address
	smp.amounts[smp.filledTransfers] = transfer.Amount
	smp.filledTransfers += 1
	smp.totalAmount += transfer.Amount

	return nil
}

func (smp *SendManyTransferPayload) AddTransfers(transfers []SendManyTransfer) error {
	if int(smp.filledTransfers)+len(transfers) > SendManyMaxTransfers {
		return fmt.Errorf("max %d send many transfers allowed", SendManyMaxTransfers)
	}

	for _, transfer := range transfers {
		err := smp.AddTransfer(transfer)
		if err != nil {
			return fmt.Errorf("adding transfer %+v: %w", transfer, err)
		}
	}

	return nil
}

func (smp *SendManyTransferPayload) GetTransfers() ([]SendManyTransfer, error) {
	transfers := make([]SendManyTransfer, 0, SendManyMaxTransfers)
	for index, address := range smp.addresses {
		if address == [32]byte{} {
			continue
		}
		var addrID common.Identity
		err := addrID.FromPubKey(address, false)
		if err != nil {
			return nil, fmt.Errorf("getting address identity from bytes %v: %w", address, err)
		}
		transfers = append(transfers, SendManyTransfer{AddressID: addrID, Amount: smp.amounts[index]})
	}

	return transfers, nil
}

// GetTotalAmount returns total amount of transfers + SC fee
func (smp *SendManyTransferPayload) GetTotalAmount() int64 {
	return smp.totalAmount + SendManyFee
}

func (smp *SendManyTransferPayload) MarshallBinary() ([]byte, error) {
	var buff bytes.Buffer
	err := binary.Write(&buff, binary.LittleEndian, smp.addresses)
	if err != nil {
		return nil, fmt.Errorf("writing addresses to buf: %w", err)
	}

	err = binary.Write(&buff, binary.LittleEndian, smp.amounts)
	if err != nil {
		return nil, fmt.Errorf("writing amounts to buf: %w", err)
	}

	return buff.Bytes(), nil
}

func (smp *SendManyTransferPayload) UnmarshallBinary(b []byte) error {
	reader := bytes.NewReader(b)

	err := binary.Read(reader, binary.LittleEndian, &smp.addresses)
	if err != nil {
		return fmt.Errorf("reading addresses from reader: %w", err)
	}

	err = binary.Read(reader, binary.LittleEndian, &smp.amounts)
	if err != nil {
		return fmt.Errorf("reading amounts from reader: %w", err)
	}

	return nil
}

func NewSendManyTransferTransaction(sourceID string, targetTick uint32, payload SendManyTransferPayload) (nodetypes.Transaction, error) {
	srcID := common.Identity(sourceID)
	destID := common.Identity(Address)
	srcPubKey, err := srcID.ToPubKey(false)
	if err != nil {
		return nodetypes.Transaction{}, fmt.Errorf("converting src id string to pubkey: %w", err)
	}
	destPubKey, err := destID.ToPubKey(false)
	if err != nil {
		return nodetypes.Transaction{}, fmt.Errorf("converting dest id string to pubkey: %w", err)
	}

	input, err := payload.MarshallBinary()
	if err != nil {
		return nodetypes.Transaction{}, fmt.Errorf("binary marshalling payload: %w", err)
	}

	return nodetypes.Transaction{
		SourcePublicKey:      srcPubKey,
		DestinationPublicKey: destPubKey,
		Amount:               payload.GetTotalAmount(),
		Tick:                 targetTick,
		InputType:            SendManyInputType,
		InputSize:            SendManyInputSize,
		Input:                input,
	}, nil
}
