package nodetypes

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/common"
	"github.com/qubic/go-qubic/connector"
	qubicpb "github.com/qubic/go-qubic/proto/v1"
	"io"
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
		return nil, errors.Wrap(err, "writing source public key to buffer")
	}

	_, err = buff.Write(tx.DestinationPublicKey[:])
	if err != nil {
		return nil, errors.Wrap(err, "writing destination public key to buffer")
	}
	err = binary.Write(&buff, binary.LittleEndian, tx.Amount)
	if err != nil {
		return nil, errors.Wrap(err, "writing amount to buf")
	}

	err = binary.Write(&buff, binary.LittleEndian, tx.Tick)
	if err != nil {
		return nil, errors.Wrap(err, "writing tick to buf")
	}

	err = binary.Write(&buff, binary.LittleEndian, tx.InputType)
	if err != nil {
		return nil, errors.Wrap(err, "writing input type to buf")
	}

	err = binary.Write(&buff, binary.LittleEndian, tx.InputSize)
	if err != nil {
		return nil, errors.Wrap(err, "writing input size to buf")
	}

	_, err = buff.Write(tx.Input)
	if err != nil {
		return nil, errors.Wrap(err, "writing input to buffer")
	}

	_, err = buff.Write(tx.Signature[:])
	if err != nil {
		return nil, errors.Wrap(err, "writing signature to buffer")
	}

	return buff.Bytes(), nil
}

func (tx *Transaction) GetUnsignedDigest() ([32]byte, error) {
	serialized, err := tx.MarshallBinary()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "marshalling tx data")
	}

	// create digest with data without signature
	digest, err := common.K12Hash(serialized[:len(serialized)-64])
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "hashing tx data")
	}

	return digest, nil
}

func (tx *Transaction) UnmarshallFromReader(r io.Reader) error {
	err := binary.Read(r, binary.LittleEndian, &tx.SourcePublicKey)
	if err != nil {
		return errors.Wrap(err, "reading source public key from reader")
	}

	err = binary.Read(r, binary.LittleEndian, &tx.DestinationPublicKey)
	if err != nil {
		return errors.Wrap(err, "reading destination public key from reader")
	}

	err = binary.Read(r, binary.LittleEndian, &tx.Amount)
	if err != nil {
		return errors.Wrap(err, "reading amount from reader")
	}

	err = binary.Read(r, binary.LittleEndian, &tx.Tick)
	if err != nil {
		return errors.Wrap(err, "reading tick from reader")
	}

	err = binary.Read(r, binary.LittleEndian, &tx.InputType)
	if err != nil {
		return errors.Wrap(err, "reading input type from reader")
	}

	err = binary.Read(r, binary.LittleEndian, &tx.InputSize)
	if err != nil {
		return errors.Wrap(err, "reading input size from reader")
	}

	tx.Input = make([]byte, tx.InputSize)
	err = binary.Read(r, binary.LittleEndian, &tx.Input)
	if err != nil {
		return errors.Wrap(err, "reading input from reader")
	}

	err = binary.Read(r, binary.LittleEndian, &tx.Signature)
	if err != nil {
		return errors.Wrap(err, "reading signature from reader")
	}

	return nil
}

func (tx *Transaction) Digest() ([32]byte, error) {
	serialized, err := tx.MarshallBinary()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "marshalling tx data")
	}

	digest, err := common.K12Hash(serialized)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "hashing tx data")
	}

	return digest, nil
}

func (tx *Transaction) EncodeToBase64() (string, error) {
	txPacket, err := tx.MarshallBinary()
	if err != nil {
		return "", errors.Wrap(err, "binary marshalling")
	}

	return base64.StdEncoding.EncodeToString(txPacket[:]), nil
}

func (tx *Transaction) ToProto() (*qubicpb.Transaction, error) {
	tc := txConverter{rawTx: *tx}
	txPb, err := tc.toProto()
	if err != nil {
		return nil, errors.Wrap(err, "calling transaction converter to proto")
	}

	return txPb, nil
}

type txConverter struct {
	rawTx Transaction
}

func (tc *txConverter) toProto() (*qubicpb.Transaction, error) {
	digest, err := tc.rawTx.Digest()
	if err != nil {
		return nil, errors.Wrap(err, "getting tx digest")
	}

	id, err := common.DigestToTxID(digest)
	if err != nil {
		return nil, errors.Wrap(err, "getting tx id")
	}

	sourceID, err := common.PubKeyToIdentity(tc.rawTx.SourcePublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "getting tx source id")
	}

	destID, err := common.PubKeyToIdentity(tc.rawTx.DestinationPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "getting tx dest id")
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
			return errors.Wrap(err, "reading header")
		}

		if header.Type == connector.EndResponse {
			break
		}

		if header.Type != TickTransactionsTypeResponse {
			return errors.Errorf("Invalid header type, expected %d, found %d", TickTransactionsTypeResponse, header.Type)
		}

		var tx Transaction

		err = tx.UnmarshallFromReader(r)
		if err != nil {
			return errors.Wrap(err, "unmarshalling transaction")
		}

		*txs = append(*txs, tx)
	}

	return nil
}

func (txs *Transactions) ToProto() (*qubicpb.TickTransactions, error) {
	ttc := tickTxsConverter{rawTxs: *txs}
	txsPb, err := ttc.toProto()
	if err != nil {
		return nil, errors.Wrap(err, "calling tick transactions converter to proto")
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
			return nil, errors.Wrapf(err, "converting to proto tx index: %d", i)
		}
		convertedTxs[i] = protoTx
	}

	return &qubicpb.TickTransactions{Transactions: convertedTxs}, nil
}

type TransactionStatus struct {
	CurrentTickOfNode  uint32
	Tick               uint32
	TxCount            uint32
	MoneyFlew          [(NumberOfTransactionsPerTick + 7) / 8]byte
	TransactionDigests [][32]byte
}

func (ts *TransactionStatus) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return errors.Wrap(err, "reading header")
	}

	if header.Type != TxStatusTypeResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", TxStatusTypeResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, &ts.CurrentTickOfNode)
	if err != nil {
		return errors.Wrap(err, "reading current tick of node")
	}

	err = binary.Read(r, binary.LittleEndian, &ts.Tick)
	if err != nil {
		return errors.Wrap(err, "reading tick")
	}

	err = binary.Read(r, binary.LittleEndian, &ts.TxCount)
	if err != nil {
		return errors.Wrap(err, "reading tx count")
	}

	err = binary.Read(r, binary.LittleEndian, &ts.MoneyFlew)
	if err != nil {
		return errors.Wrap(err, "reading reading money flew")
	}

	ts.TransactionDigests = make([][32]byte, ts.TxCount)
	err = binary.Read(r, binary.LittleEndian, &ts.TransactionDigests)
	if err != nil {
		return errors.Wrap(err, "reading tx digests")
	}

	return nil
}

func (ts *TransactionStatus) ToProto() (*qubicpb.TickTransactionsStatus, error) {
	tsc := transactionsStatusConverter{rawTxStatus: *ts}
	tsPb, err := tsc.toProto()
	if err != nil {
		return nil, errors.Wrap(err, "calling tick transactions status converter to proto")
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
			return nil, errors.Wrapf(err, "getting tx id for tx with digest hex: %s", hex.EncodeToString(digest[:]))
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
