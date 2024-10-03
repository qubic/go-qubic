package nodetypes

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/common"
	"github.com/qubic/go-qubic/connector"
	qubicpb "github.com/qubic/go-qubic/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"time"
)

const (
	MinimumQuorumVotes = 451
)

const (
	QuorumTickTypeRequest  = 14
	QuorumTickTypeResponse = 3
)

type QuorumTickVote struct {
	ComputorIndex uint16
	Epoch         uint16
	Tick          uint32

	Millisecond uint16
	Second      uint8
	Minute      uint8
	Hour        uint8
	Day         uint8
	Month       uint8
	Year        uint8

	PreviousResourceTestingDigest uint64
	SaltedResourceTestingDigest   uint64

	PreviousSpectrumDigest [32]byte
	PreviousUniverseDigest [32]byte
	PreviousComputerDigest [32]byte

	SaltedSpectrumDigest [32]byte
	SaltedUniverseDigest [32]byte
	SaltedComputerDigest [32]byte

	TxDigest                 [32]byte
	ExpectedNextTickTxDigest [32]byte

	Signature [SignatureSize]byte
}

type QuorumVotes []QuorumTickVote

func (qv *QuorumVotes) UnmarshallFromReader(r io.Reader) error {
	for {
		var header connector.RequestResponseHeader
		err := binary.Read(r, binary.BigEndian, &header)
		if err != nil {
			return errors.Wrap(err, "reading header")
		}

		if header.Type == connector.EndResponse {
			break
		}

		var qtd QuorumTickVote
		if header.Type != QuorumTickTypeResponse {
			return errors.Errorf("Invalid header type, expected %d, found %d", QuorumTickTypeResponse, header.Type)
		}

		err = binary.Read(r, binary.LittleEndian, &qtd)
		if err != nil {
			return errors.Wrap(err, "reading quorum tick data from reader")
		}

		*qv = append(*qv, qtd)
	}

	return nil
}

func (qv *QuorumVotes) ToProto() (*qubicpb.QuorumVote, error) {
	qc := quorumConverter{quorumVotes: *qv}
	qvPb, err := qc.toProto()
	if err != nil {
		return nil, errors.Wrap(err, "calling quorum converter to proto")
	}

	return qvPb, nil
}

type quorumConverter struct {
	quorumVotes QuorumVotes
}

func (qc *quorumConverter) toProto() (*qubicpb.QuorumVote, error) {
	sharedVotes, err := qc.toSharedVotes()
	if err != nil {
		return nil, errors.Wrap(err, "to shared votes")
	}

	saltedVotes := qc.toSaltedVotes()

	return &qubicpb.QuorumVote{
		SharedVotes:                 sharedVotes,
		SaltedVotesPerComputorIndex: saltedVotes,
	}, nil
}

func (qc *quorumConverter) toSaltedVotes() map[uint32]*qubicpb.QuorumVote_SaltedVote {
	saltedVotes := make(map[uint32]*qubicpb.QuorumVote_SaltedVote)

	for _, vote := range qc.quorumVotes {
		saltedVotes[uint32(vote.ComputorIndex)] = qc.toSaltedVote(vote)
	}

	return saltedVotes
}

func (qc *quorumConverter) toSaltedVote(vote QuorumTickVote) *qubicpb.QuorumVote_SaltedVote {
	return &qubicpb.QuorumVote_SaltedVote{
		ResourceTestingDigest:    convertUint64ToBase64(vote.SaltedResourceTestingDigest),
		SpectrumDigest:           base64.StdEncoding.EncodeToString(vote.SaltedSpectrumDigest[:]),
		UniverseDigest:           base64.StdEncoding.EncodeToString(vote.SaltedUniverseDigest[:]),
		ComputerDigest:           base64.StdEncoding.EncodeToString(vote.SaltedComputerDigest[:]),
		ExpectedNextTickTxDigest: base64.StdEncoding.EncodeToString(vote.ExpectedNextTickTxDigest[:]),
		Signature:                base64.StdEncoding.EncodeToString(vote.Signature[:]),
	}
}

func (qc *quorumConverter) toSharedVotes() ([]*qubicpb.QuorumVote_GroupedSharedVotes, error) {
	votesHeatMap := make(map[[32]byte][]*qubicpb.QuorumVote_SharedVote)
	for _, qv := range qc.quorumVotes {
		hv := heatmapVote{
			Epoch:                         qv.Epoch,
			Tick:                          qv.Tick,
			Millisecond:                   qv.Millisecond,
			Second:                        qv.Second,
			Minute:                        qv.Minute,
			Hour:                          qv.Hour,
			Day:                           qv.Day,
			Month:                         qv.Month,
			Year:                          qv.Year,
			PreviousResourceTestingDigest: qv.PreviousResourceTestingDigest,
			PreviousSpectrumDigest:        qv.PreviousSpectrumDigest,
			PreviousUniverseDigest:        qv.PreviousUniverseDigest,
			PreviousComputerDigest:        qv.PreviousComputerDigest,
			TxDigest:                      qv.TxDigest,
		}

		digest, err := hv.digest()
		if err != nil {
			return nil, errors.Wrap(err, "getting digest")
		}

		sv := qc.toSharedVote(qv)
		if votes, ok := votesHeatMap[digest]; !ok {
			votesHeatMap[digest] = []*qubicpb.QuorumVote_SharedVote{sv}
		} else {
			votesHeatMap[digest] = append(votes, sv)
		}
	}

	sharedVotes := make([]*qubicpb.QuorumVote_GroupedSharedVotes, 0, len(votesHeatMap))
	for _, votes := range votesHeatMap {
		sharedVotes = append(sharedVotes, &qubicpb.QuorumVote_GroupedSharedVotes{Vote: votes[0], NumberOfVotes: int32(len(votes))})
	}

	return sharedVotes, nil
}

func (qc *quorumConverter) toSharedVote(vote QuorumTickVote) *qubicpb.QuorumVote_SharedVote {
	date := time.Date(2000+int(vote.Year), time.Month(vote.Month), int(vote.Day), int(vote.Hour), int(vote.Minute), int(vote.Second), 0, time.UTC)
	date.Add(time.Duration(vote.Millisecond) * time.Millisecond)

	return &qubicpb.QuorumVote_SharedVote{
		Epoch:                     uint32(vote.Epoch),
		Tick:                      vote.Tick,
		Timestamp:                 timestamppb.New(date),
		PrevResourceTestingDigest: convertUint64ToBase64(vote.PreviousResourceTestingDigest),
		PrevSpectrumDigest:        base64.StdEncoding.EncodeToString(vote.PreviousSpectrumDigest[:]),
		PrevUniverseDigest:        base64.StdEncoding.EncodeToString(vote.PreviousUniverseDigest[:]),
		PrevComputerDigest:        base64.StdEncoding.EncodeToString(vote.PreviousComputerDigest[:]),
		TxDigest:                  base64.StdEncoding.EncodeToString(vote.TxDigest[:]),
	}
}

func (qc *quorumConverter) getVoteDigest(vote QuorumTickVote) ([32]byte, error) {
	// xor computor index with 8
	vote.ComputorIndex ^= 3

	sData, err := common.BinarySerializeLE(vote)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "serializing data")
	}

	tickData := sData[:len(sData)-64]
	digest, err := common.K12Hash(tickData)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "hashing tick data")
	}

	return digest, nil
}

type heatmapVote struct {
	Epoch                         uint16
	Tick                          uint32
	Millisecond                   uint16
	Second                        uint8
	Minute                        uint8
	Hour                          uint8
	Day                           uint8
	Month                         uint8
	Year                          uint8
	PreviousResourceTestingDigest uint64
	PreviousSpectrumDigest        [32]byte
	PreviousUniverseDigest        [32]byte
	PreviousComputerDigest        [32]byte
	TxDigest                      [32]byte
}

func (hv *heatmapVote) digest() ([32]byte, error) {
	b, err := common.BinarySerializeLE(hv)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "serializing vote")
	}

	digest, err := common.K12Hash(b)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "hashing vote")
	}

	return digest, nil
}

func convertUint64ToBase64(value uint64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, value)
	return base64.StdEncoding.EncodeToString(b)
}
