package nodetypes

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/internal/connector"
	"io"
)

const (
	SignatureSize     = 64
	NumberOfComputors = 676
)

const (
	ComputorsTypeRequest  = 11
	ComputorsTypeResponse = 2
)

type Computors struct {
	Epoch     uint16
	PubKeys   [NumberOfComputors][32]byte
	Signature [SignatureSize]byte
}

func (cs *Computors) UnmarshallFromReader(r io.Reader) error {
	for {
		var header connector.RequestResponseHeader
		headerSize := binary.Size(header)
		err := binary.Read(r, binary.BigEndian, &header)
		if err != nil {
			return errors.Wrap(err, "reading header")
		}

		if header.Type != ComputorsTypeResponse {
			ignoredbytes := make([]byte, header.GetSize()-uint32(headerSize))
			_, err := r.Read(ignoredbytes)
			if err != nil {
				return errors.Wrap(err, "reading ignored bytes")
			}
			continue
		}

		err = binary.Read(r, binary.LittleEndian, cs)
		if err != nil {
			return errors.Wrap(err, "reading computors from reader")
		}

		return nil
	}
}
