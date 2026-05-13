package quottery

import (
	"fmt"
	"time"

	"github.com/qubic/go-qubic/v2/common"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var BetInfoConverter betInfoConverter

type betInfoConverter struct{}

func (bic betInfoConverter) ToProto(bi BetInfo) (*qubicpb.BetInfo, error) {
	var creatorID common.Identity
	err := creatorID.FromPubKey(bi.Creator, false)
	if err != nil {
		return nil, fmt.Errorf("converting creator id: %w", err)
	}

	description := string(bi.Description[:])

	maxNOption := 8
	if int(bi.NumberOfOptions) < maxNOption {
		maxNOption = int(bi.NumberOfOptions)
	}

	betOptions := make([]*qubicpb.BetInfo_Option, 0, maxNOption)
	for i := 0; i < maxNOption; i++ {
		optDescription := string(bi.OptionDescription[i*32 : (i*32)+32])
		state := bi.CurrentBetState[i]
		betOptions = append(betOptions, &qubicpb.BetInfo_Option{Description: optDescription, State: state})
	}

	oraclesInfo := make([]*qubicpb.BetInfo_Oracle, 0, 8)
	oracleVotes := make([]*qubicpb.BetInfo_Vote, 0, 8)
	for i := 0; i < 8; i++ {
		var providerPubKey [32]byte
		copy(providerPubKey[:], bi.OracleProvider[i*32:(i*32)+32])
		if providerPubKey == [32]byte{} {
			continue
		}

		var providerID common.Identity
		err := providerID.FromPubKey(providerPubKey, false)
		if err != nil {
			return nil, fmt.Errorf("converting provider id with pubkey: %s: %w", providerPubKey, err)
		}
		oi := qubicpb.BetInfo_Oracle{
			Id:            providerID.String(),
			FeePercentage: float32(bi.OracleFees[i]) / 100,
		}
		oraclesInfo = append(oraclesInfo, &oi)

		if bi.BetResultWonOption[i] == -1 && bi.BetResultOpID[i] == -1 {
			continue
		}

		ov := qubicpb.BetInfo_Vote{
			OracleId:  uint32(bi.BetResultOpID[i]),
			WonOption: uint32(bi.BetResultWonOption[i]),
		}

		oracleVotes = append(oracleVotes, &ov)
	}

	openDate := quotteryParseDate(bi.OpenDate)
	closeDate := quotteryParseDate(bi.CloseDate)
	endDate := quotteryParseDate(bi.EndDate)

	return &qubicpb.BetInfo{
		Id:                      bi.ID,
		CreatorId:               creatorID.String(),
		Description:             description,
		Options:                 betOptions,
		Oracles:                 oraclesInfo,
		Votes:                   oracleVotes,
		MinimumBetAmount:        bi.MinimumBetAmount,
		MaximumBetSlotPerOption: bi.MaximumBetSlotPerOption,
		OpenTime:                openDate,
		CloseTime:               closeDate,
		EndTime:                 endDate,
	}, nil
}

func quotteryParseDate(date uint32) *timestamppb.Timestamp {
	return timestamppb.New(time.Date(quotteryGetYear(date), time.Month(quotteryGetMonth(date)), quotteryGetDay(date), quotteryGetHour(date), quotteryGetMinute(date), quotteryGetSecond(date), 0, time.UTC))
}

func quotteryGetYear(data uint32) int {
	return 2000 + int((data>>26)+24)
}

func quotteryGetMonth(data uint32) int {
	return int((data >> 22) & 0b1111)
}

func quotteryGetDay(data uint32) int {
	return int((data >> 17) & 0b11111)
}

func quotteryGetHour(data uint32) int {
	return int((data >> 12) & 0b11111)
}

func quotteryGetMinute(data uint32) int {
	return int((data >> 6) & 0b111111)
}

func quotteryGetSecond(data uint32) int {
	return int(data & 0b111111)
}

var ActiveBetsConverter activeBetsConverter

type activeBetsConverter struct{}

func (abc activeBetsConverter) ToProto(ab ActiveBets) *qubicpb.ActiveBets {
	betIDs := make([]uint32, 0, ab.Count)

	for i := 0; i < int(ab.Count); i++ {
		betIDs = append(betIDs, ab.BetIDs[i])
	}

	return &qubicpb.ActiveBets{
		ActiveBetIds: betIDs,
	}
}

var BasicInfoConverter basicInfoConverter

type basicInfoConverter struct{}

func (bic basicInfoConverter) ToProto(bi BasicInfo) (*qubicpb.BasicInfo, error) {
	gameOperatorID, err := common.PubKeyToIdentity(bi.GameOperatorPubKey)
	if err != nil {
		return nil, fmt.Errorf("converting game operator id with pubkey: %s: %w", bi.GameOperatorPubKey, err)
	}

	return &qubicpb.BasicInfo{
		Fees: &qubicpb.BasicInfo_Fees{
			SlotPerDay:   bi.FeePerSlotPerDay,
			GameOperator: bi.GameOperatorFee,
			Shareholder:  bi.ShareholderFee,
			Burn:         bi.BurnFee,
		},
		MinimumBetSlotAmount: bi.MinimumBetSlotAmount,
		IssuedBets:           bi.IssuedBets,
		MoneyFlowData: &qubicpb.BasicInfo_MoneyFlowData{
			Total:       bi.MoneyFlow,
			IssueBet:    bi.MoneyFlowIssueBet,
			JoinBet:     bi.MoneyFlowJoinBet,
			FinalizeBet: bi.MoneyFlowFinalizeBet,
		},
		EconomicsData: &qubicpb.BasicInfo_EconomicsData{
			EarnedAmountShareholder: bi.EarnedAmountForShareholder,
			PaidAmountShareholder:   bi.PaidAmountForShareholder,
			EarnedAmountBetWinner:   bi.EarnedAmountForBetWinner,
			DistributedAmount:       bi.DistributedAmount,
			BurnedAmount:            bi.BurnedAmount,
		},
		GameOperatorId: gameOperatorID.String(),
	}, nil
}

var BetOptionBettorsConverter betOptionBettorsConverter

type betOptionBettorsConverter struct{}

func (bobc betOptionBettorsConverter) ToProto(bod BetOptionDetail) (*qubicpb.BetOptionBettors, error) {
	bettorIDs := make([]string, 0)

	// declare once to avoid unnecessary memory allocation
	var idBytes [32]byte

	for i := 0; i < 1024; i++ {
		offset := i * 32
		copy(idBytes[:], bod.Bettor[offset:offset+32])
		if idBytes == [32]byte{} {
			continue
		}

		bettor, err := common.PubKeyToIdentity(idBytes)
		if err != nil {
			return nil, fmt.Errorf("converting bettor id with pubkey: %s: %w", idBytes, err)
		}

		bettorIDs = append(bettorIDs, bettor.String())
	}

	return &qubicpb.BetOptionBettors{BettorIds: bettorIDs}, nil
}
