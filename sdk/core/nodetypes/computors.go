package nodetypes

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/common"
	"github.com/qubic/go-qubic/connector"
	qubicpb "github.com/qubic/go-qubic/proto/v1"
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

func (cs *Computors) ToProto() (*qubicpb.Computors, error) {
	cc := computorsConverter{comps: *cs}
	csPb, err := cc.toProto()
	if err != nil {
		return nil, errors.Wrap(err, "calling computors converter to proto")
	}

	return csPb, nil
}

type computorsConverter struct {
	comps Computors
}

func (cc computorsConverter) toProto() (*qubicpb.Computors, error) {
	identities, err := common.PubKeysToIdentitiesString(cc.comps.PubKeys[:], false)
	if err != nil {
		return nil, errors.Wrap(err, "converting pubKeys to identities")
	}

	digest, err := cc.getDigest()
	if err != nil {
		return nil, errors.Wrap(err, "creating computors digest")
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
		return [32]byte{}, errors.Wrap(err, "serializing data")
	}

	// remove signature from computors data
	computorsData := serialized[:len(serialized)-64]
	digest, err := common.K12Hash(computorsData)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "hashing computors data")
	}

	return digest, nil
}
