package connector

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"math/rand"
)

const EndResponse = 35
const contractFunctionRequest = 42
const ContractFunctionResponse = 43
const broadcastTransactionRequestType = 24

type RequestResponseHeader struct {
	Size   [3]uint8
	Type   uint8
	DejaVu uint32
}

func (h *RequestResponseHeader) GetSize() uint32 {
	// Convert the array to a 32-bit unsigned integer
	size := uint32(h.Size[0]) | (uint32(h.Size[1]) << 8) | (uint32(h.Size[2]) << 16)

	// Apply the bitwise AND operation to keep the lower 24 bits
	result := size & 0xFFFFFF

	return result
}

func (h *RequestResponseHeader) SetSize(size uint32) {
	h.Size[0] = uint8(size)
	h.Size[1] = uint8(size >> 8)
	h.Size[2] = uint8(size >> 16)
}

func (h *RequestResponseHeader) IsDejaVuZero() bool {
	return h.DejaVu == 0
}

func (h *RequestResponseHeader) ZeroDejaVu() {
	h.DejaVu = 0
}

func (h *RequestResponseHeader) RandomizeDejaVu() {
	h.DejaVu = uint32(rand.Int31())
	if h.DejaVu == 0 {
		h.DejaVu = 1
	}
}

func (h *RequestResponseHeader) UnmarshallFromReader(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, h)
	if err != nil {
		return errors.Wrap(err, "reading quottery basic info header")
	}

	return nil
}

type chainRequest struct {
	requestType uint8
	requestData interface{}
}

func newChainRequest(requestType uint8, requestData interface{}) chainRequest {
	return chainRequest{requestType: requestType, requestData: requestData}
}

func (r *chainRequest) serialize() ([]byte, error) {
	serializedReqData, err := binarySerialize(r.requestData)
	if err != nil {
		return nil, errors.Wrap(err, "serializing req data")
	}

	var header RequestResponseHeader

	packetHeaderSize := binary.Size(header)
	reqDataSize := len(serializedReqData)
	packetSize := uint32(packetHeaderSize + reqDataSize)

	header.SetSize(packetSize)
	if r.requestType == broadcastTransactionRequestType {
		header.ZeroDejaVu()
	} else {
		header.RandomizeDejaVu()
	}

	header.Type = r.requestType

	serializedHeaderData, err := binarySerialize(header)
	if err != nil {
		return nil, errors.Wrap(err, "serializing header data")
	}

	serializedPacket := make([]byte, 0, packetSize)
	serializedPacket = append(serializedPacket, serializedHeaderData...)
	serializedPacket = append(serializedPacket, serializedReqData...)

	return serializedPacket, nil
}

func binarySerialize(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	var buff bytes.Buffer
	err := binary.Write(&buff, binary.LittleEndian, data)
	if err != nil {
		return nil, errors.Wrap(err, "writing data to buff")
	}

	return buff.Bytes(), nil
}

type smartContractRequest struct {
	reqContractFunction RequestContractFunction

	requestType uint8
	requestData interface{}
}

func newSmartContractRequest(reqContractFunction RequestContractFunction, requestData interface{}) smartContractRequest {
	return smartContractRequest{reqContractFunction: reqContractFunction, requestType: contractFunctionRequest, requestData: requestData}
}

func (r *smartContractRequest) serialize() ([]byte, error) {
	serializedReqData, err := binarySerialize(r.requestData)
	if err != nil {
		return nil, errors.Wrap(err, "serializing req data")
	}

	serializedReqContractFunction, err := binarySerialize(r.reqContractFunction)
	if err != nil {
		return nil, errors.Wrap(err, "serializing req contract function")
	}

	var header RequestResponseHeader

	packetHeaderSize := binary.Size(header)
	reqDataSize := len(serializedReqData)
	reqContractFunctionSize := len(serializedReqContractFunction)
	packetSize := uint32(packetHeaderSize + reqContractFunctionSize + reqDataSize)

	header.RandomizeDejaVu()

	header.Type = r.requestType
	header.SetSize(packetSize)

	serializedHeaderData, err := binarySerialize(header)
	if err != nil {
		return nil, errors.Wrap(err, "serializing header data")
	}

	serializedPacket := make([]byte, 0, packetSize)
	serializedPacket = append(serializedPacket, serializedHeaderData...)
	serializedPacket = append(serializedPacket, serializedReqContractFunction...)
	serializedPacket = append(serializedPacket, serializedReqData...)

	return serializedPacket, nil
}

type RequestContractFunction struct {
	ContractIndex uint32
	InputType     uint16
	InputSize     uint16
}
