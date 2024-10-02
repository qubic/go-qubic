package quottery

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/connector"
	"io"
)

const contractID = 2

type viewID struct {
	BasicInfo          uint16
	BetInfo            uint16
	BetDetail          uint16
	ActiveBet          uint16
	ActiveBetByCreator uint16
}

type funcID struct {
	Issue         uint16
	Join          uint16
	CancelBet     uint16
	PublishResult uint16
}

var ViewID = viewID{
	BasicInfo:          1,
	BetInfo:            2,
	BetDetail:          3,
	ActiveBet:          4,
	ActiveBetByCreator: 5,
}

var FuncID = funcID{
	Issue:         1,
	Join:          2,
	CancelBet:     3,
	PublishResult: 4,
}

type BetInfo struct {
	ID                      uint32
	NumberOfOptions         uint32
	Creator                 [32]byte
	Description             [32]byte
	OptionDescription       [256]byte
	OracleProvider          [256]byte
	OracleFees              [8]uint32
	OpenDate                uint32
	CloseDate               uint32
	EndDate                 uint32
	_                       uint32 //padding
	MinimumBetAmount        uint64
	MaximumBetSlotPerOption uint32
	CurrentBetState         [8]uint32
	BetResultWonOption      [8]int8
	BetResultOpID           [8]int8
}

func (bi *BetInfo) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return errors.Wrap(err, "reading quottery bet info header")
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != connector.ContractFunctionResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", connector.ContractFunctionResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, bi)
	if err != nil {
		return errors.Wrap(err, "reading quottery bet info data")
	}

	return nil
}

type ActiveBets struct {
	Count  uint32
	BetIDs [1024]uint32
}

func (ab *ActiveBets) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return errors.Wrap(err, "reading quottery active bets header")
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != connector.ContractFunctionResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", connector.ContractFunctionResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, ab)
	if err != nil {
		return errors.Wrap(err, "reading quottery active bets count")
	}

	return nil
}

type BasicInfo struct {
	FeePerSlotPerDay           uint64
	GameOperatorFee            uint64
	ShareholderFee             uint64
	MinimumBetSlotAmount       uint64
	BurnFee                    uint64
	IssuedBets                 uint64
	MoneyFlow                  uint64
	MoneyFlowIssueBet          uint64
	MoneyFlowJoinBet           uint64
	MoneyFlowFinalizeBet       uint64
	EarnedAmountForShareholder uint64
	PaidAmountForShareholder   uint64
	EarnedAmountForBetWinner   uint64
	DistributedAmount          uint64
	BurnedAmount               uint64
	GameOperatorPubKey         [32]byte
}

func (bi *BasicInfo) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return errors.Wrap(err, "reading quottery basic info header")
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != connector.ContractFunctionResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", connector.ContractFunctionResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, bi)
	if err != nil {
		return errors.Wrap(err, "reading quottery basic info data")
	}

	return nil
}

type BetOptionDetail struct {
	Bettor [32 * 1024]byte
}

func (bod *BetOptionDetail) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader

	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return errors.Wrap(err, "reading quottery bet option detail header")
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != connector.ContractFunctionResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", connector.ContractFunctionResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, bod)
	if err != nil {
		return errors.Wrap(err, "reading quottery bet option detail data")
	}

	return nil
}
