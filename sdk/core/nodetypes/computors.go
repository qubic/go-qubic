package nodetypes

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/qubic/go-qubic/v2/common"
	"github.com/qubic/go-qubic/v2/connector"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
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
			return fmt.Errorf("reading header: %w", err)
		}

		if header.Type != ComputorsTypeResponse {
			ignoredbytes := make([]byte, header.GetSize()-uint32(headerSize))
			_, err := r.Read(ignoredbytes)
			if err != nil {
				return fmt.Errorf("reading ignored bytes: %w", err)
			}
			continue
		}

		err = binary.Read(r, binary.LittleEndian, cs)
		if err != nil {
			return fmt.Errorf("reading computors from reader: %w", err)
		}

		return nil
	}
}

func (cs *Computors) ToProto() (*qubicpb.Computors, error) {
	cc := computorsConverter{comps: *cs}
	csPb, err := cc.toProto()
	if err != nil {
		return nil, fmt.Errorf("calling computors converter to proto: %w", err)
	}

	return csPb, nil
}

type computorsConverter struct {
	comps Computors
}

func (cc computorsConverter) toProto() (*qubicpb.Computors, error) {
	identities, err := common.PubKeysToIdentitiesString(cc.comps.PubKeys[:], false)
	if err != nil {
		return nil, fmt.Errorf("converting pubKeys to identities: %w", err)
	}

	digest, err := cc.getDigest()
	if err != nil {
		return nil, fmt.Errorf("creating computors digest: %w", err)
	}

	return &qubicpb.Computors{
		Epoch:      uint32(cc.comps.Epoch),
		Identities: identities,
		Signature:  base64.StdEncoding.EncodeToString(cc.comps.Signature[:]),
		Digest:     base64.StdEncoding.EncodeToString(digest[:]),
	}, nil
}

func (cc computorsConverter) getDigest() ([32]byte, error) {
	serialized, err := common.BinarySerializeLE(cc.comps)
	if err != nil {
		return [32]byte{}, fmt.Errorf("serializing data: %w", err)
	}

	// remove signature from computors data
	computorsData := serialized[:len(serialized)-64]
	digest, err := common.K12Hash(computorsData)
	if err != nil {
		return [32]byte{}, fmt.Errorf("hashing computors data: %w", err)
	}

	return digest, nil
}
