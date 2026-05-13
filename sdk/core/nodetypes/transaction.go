package nodetypes

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/qubic/go-qubic/v2/common"
	"github.com/qubic/go-qubic/v2/connector"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
)

const (
	TickTransactionsTypeRequest  = 29
	TickTransactionsTypeResponse = 24
	TxStatusTypeRequest          = 201
	TxStatusTypeResponse         = 202
)

type Transaction struct {
	SourcePublicKey      [32]byte
	DestinationPublicKey [32]byte
	Amount               int64
	Tick                 uint32
	InputType            uint16
	InputSize            uint16
	Input                []byte
	Signature            [64]byte
}

func (tx *Transaction) MarshallBinary() ([]byte, error) {
	var buff bytes.Buffer
	_, err := buff.Write(tx.SourcePublicKey[:])
	if err != nil {
		return nil, fmt.Errorf("writing source public key to buffer: %w", err)
	}

	_, err = buff.Write(tx.DestinationPublicKey[:])
	if err != nil {
		return nil, fmt.Errorf("writing destination public key to buffer: %w", err)
	}
	err = binary.Write(&buff, binary.LittleEndian, tx.Amount)
	if err != nil {
		return nil, fmt.Errorf("writing amount to buf: %w", err)
	}

	err = binary.Write(&buff, binary.LittleEndian, tx.Tick)
	if err != nil {
		return nil, fmt.Errorf("writing tick to buf: %w", err)
	}

	err = binary.Write(&buff, binary.LittleEndian, tx.InputType)
	if err != nil {
		return nil, fmt.Errorf("writing input type to buf: %w", err)
	}

	err = binary.Write(&buff, binary.LittleEndian, tx.InputSize)
	if err != nil {
		return nil, fmt.Errorf("writing input size to buf: %w", err)
	}

	_, err = buff.Write(tx.Input)
	if err != nil {
		return nil, fmt.Errorf("writing input to buffer: %w", err)
	}

	_, err = buff.Write(tx.Signature[:])
	if err != nil {
		return nil, fmt.Errorf("writing signature to buffer: %w", err)
	}

	return buff.Bytes(), nil
}

func (tx *Transaction) GetUnsignedDigest() ([32]byte, error) {
	serialized, err := tx.MarshallBinary()
	if err != nil {
		return [32]byte{}, fmt.Errorf("marshalling tx data: %w", err)
	}

	// create digest with data without signature
	digest, err := common.K12Hash(serialized[:len(serialized)-64])
	if err != nil {
		return [32]byte{}, fmt.Errorf("hashing tx data: %w", err)
	}

	return digest, nil
}

func (tx *Transaction) UnmarshallFromReader(r io.Reader) error {
	err := binary.Read(r, binary.LittleEndian, &tx.SourcePublicKey)
	if err != nil {
		return fmt.Errorf("reading source public key from reader: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &tx.DestinationPublicKey)
	if err != nil {
		return fmt.Errorf("reading destination public key from reader: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &tx.Amount)
	if err != nil {
		return fmt.Errorf("reading amount from reader: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &tx.Tick)
	if err != nil {
		return fmt.Errorf("reading tick from reader: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &tx.InputType)
	if err != nil {
		return fmt.Errorf("reading input type from reader: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &tx.InputSize)
	if err != nil {
		return fmt.Errorf("reading input size from reader: %w", err)
	}

	tx.Input = make([]byte, tx.InputSize)
	err = binary.Read(r, binary.LittleEndian, &tx.Input)
	if err != nil {
		return fmt.Errorf("reading input from reader: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &tx.Signature)
	if err != nil {
		return fmt.Errorf("reading signature from reader: %w", err)
	}

	return nil
}

func (tx *Transaction) Digest() ([32]byte, error) {
	serialized, err := tx.MarshallBinary()
	if err != nil {
		return [32]byte{}, fmt.Errorf("marshalling tx data: %w", err)
	}

	digest, err := common.K12Hash(serialized)
	if err != nil {
		return [32]byte{}, fmt.Errorf("hashing tx data: %w", err)
	}

	return digest, nil
}

func (tx *Transaction) EncodeToBase64() (string, error) {
	txPacket, err := tx.MarshallBinary()
	if err != nil {
		return "", fmt.Errorf("binary marshalling: %w", err)
	}

	return base64.StdEncoding.EncodeToString(txPacket[:]), nil
}

func (tx *Transaction) ToProto() (*qubicpb.Transaction, error) {
	tc := txConverter{rawTx: *tx}
	txPb, err := tc.toProto()
	if err != nil {
		return nil, fmt.Errorf("calling transaction converter to proto: %w", err)
	}

	return txPb, nil
}

type txConverter struct {
	rawTx Transaction
}

func (tc *txConverter) toProto() (*qubicpb.Transaction, error) {
	digest, err := tc.rawTx.Digest()
	if err != nil {
		return nil, fmt.Errorf("getting tx digest: %w", err)
	}

	id, err := common.DigestToTxID(digest)
	if err != nil {
		return nil, fmt.Errorf("getting tx id: %w", err)
	}

	sourceID, err := common.PubKeyToIdentity(tc.rawTx.SourcePublicKey)
	if err != nil {
		return nil, fmt.Errorf("getting tx source id: %w", err)
	}

	destID, err := common.PubKeyToIdentity(tc.rawTx.DestinationPublicKey)
	if err != nil {
		return nil, fmt.Errorf("getting tx dest id: %w", err)
	}

	return &qubicpb.Transaction{
		SourceId:  sourceID.String(),
		DestId:    destID.String(),
		Amount:    tc.rawTx.Amount,
		Tick:      tc.rawTx.Tick,
		InputType: uint32(tc.rawTx.InputType),
		InputSize: uint32(tc.rawTx.InputSize),
		Input:     base64.StdEncoding.EncodeToString(tc.rawTx.Input),
		Signature: base64.StdEncoding.EncodeToString(tc.rawTx.Signature[:]),
		TxId:      id.String(),
		Digest:    base64.StdEncoding.EncodeToString(digest[:]),
	}, nil
}

type Transactions []Transaction

func (txs *Transactions) UnmarshallFromReader(r io.Reader) error {
	for {
		var header connector.RequestResponseHeader
		err := binary.Read(r, binary.BigEndian, &header)
		if err != nil {
			return fmt.Errorf("reading header: %w", err)
		}

		if header.Type == connector.EndResponse {
			break
		}

		if header.Type != TickTransactionsTypeResponse {
			return fmt.Errorf("Invalid header type, expected %d, found %d", TickTransactionsTypeResponse, header.Type)
		}

		var tx Transaction

		err = tx.UnmarshallFromReader(r)
		if err != nil {
			return fmt.Errorf("unmarshalling transaction: %w", err)
		}

		*txs = append(*txs, tx)
	}

	return nil
}

func (txs *Transactions) ToProto() (*qubicpb.TickTransactions, error) {
	ttc := tickTxsConverter{rawTxs: *txs}
	txsPb, err := ttc.toProto()
	if err != nil {
		return nil, fmt.Errorf("calling tick transactions converter to proto: %w", err)
	}

	return txsPb, nil
}

type tickTxsConverter struct {
	rawTxs []Transaction
}

func (ttc *tickTxsConverter) toProto() (*qubicpb.TickTransactions, error) {
	convertedTxs := make([]*qubicpb.Transaction, len(ttc.rawTxs))
	for i, tx := range ttc.rawTxs {
		protoTx, err := tx.ToProto()
		if err != nil {
			return nil, fmt.Errorf("converting to proto tx index: %d: %w", i, err)
		}
		convertedTxs[i] = protoTx
	}

	return &qubicpb.TickTransactions{Transactions: convertedTxs}, nil
}

type TransactionStatus struct {
	CurrentTickOfNode  uint32
	Tick               uint32
	TxCount            uint32
	MoneyFlew          [(MaxNumberOfTransactionsPerTick + 7) / 8]byte
	TransactionDigests [][32]byte
}

func (ts *TransactionStatus) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return fmt.Errorf("reading header: %w", err)
	}

	if header.Type != TxStatusTypeResponse {
		return fmt.Errorf("Invalid header type, expected %d, found %d", TxStatusTypeResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, &ts.CurrentTickOfNode)
	if err != nil {
		return fmt.Errorf("reading current tick of node: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &ts.Tick)
	if err != nil {
		return fmt.Errorf("reading tick: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &ts.TxCount)
	if err != nil {
		return fmt.Errorf("reading tx count: %w", err)
	}

	err = binary.Read(r, binary.LittleEndian, &ts.MoneyFlew)
	if err != nil {
		return fmt.Errorf("reading reading money flew: %w", err)
	}

	ts.TransactionDigests = make([][32]byte, ts.TxCount)
	err = binary.Read(r, binary.LittleEndian, &ts.TransactionDigests)
	if err != nil {
		return fmt.Errorf("reading tx digests: %w", err)
	}

	return nil
}

func (ts *TransactionStatus) ToProto() (*qubicpb.TickTransactionsStatus, error) {
	tsc := transactionsStatusConverter{rawTxStatus: *ts}
	tsPb, err := tsc.toProto()
	if err != nil {
		return nil, fmt.Errorf("calling tick transactions status converter to proto: %w", err)
	}

	return tsPb, nil
}

type transactionsStatusConverter struct {
	rawTxStatus TransactionStatus
}

func (tsc *transactionsStatusConverter) toProto() (*qubicpb.TickTransactionsStatus, error) {
	statuses := make(map[string]bool)

	for index, digest := range tsc.rawTxStatus.TransactionDigests {
		id, err := common.DigestToTxID(digest)
		if err != nil {
			return nil, fmt.Errorf("getting tx id for tx with digest hex: %s: %w", hex.EncodeToString(digest[:]), err)
		}

		moneyFlew := tsc.getMoneyFlewFromBits(index)
		statuses[id.String()] = moneyFlew
	}

	return &qubicpb.TickTransactionsStatus{
		CurrentTickOfNode: tsc.rawTxStatus.CurrentTickOfNode,
		Tick:              tsc.rawTxStatus.Tick,
		TxCount:           tsc.rawTxStatus.TxCount,
		StatusPerTx:       statuses,
	}, nil
}

func (tsc *transactionsStatusConverter) getMoneyFlewFromBits(digestIndex int) bool {
	pos := digestIndex / 8
	bitIndex := digestIndex % 8

	return tsc.getNthBit(pos, bitIndex)
}

func (tsc *transactionsStatusConverter) getNthBit(inputPos, bitIndex int) bool {
	input := tsc.rawTxStatus.MoneyFlew[inputPos]
	// Shift the input byte to the right by the bitIndex positions
	// This isolates the bit at the bitIndex position at the least significant bit position
	shifted := input >> bitIndex

	// Extract the least significant bit using a bitwise AND operation with 1
	// If the least significant bit is 1, the result will be 1; otherwise, it will be 0
	return shifted&1 == 1
}
