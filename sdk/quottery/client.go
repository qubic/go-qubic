package quottery

import (
	"context"
	"fmt"

	"github.com/qubic/go-qubic/v2/common"
	"github.com/qubic/go-qubic/v2/connector"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
)

type Client struct {
	connector *connector.Connector
}

func NewClient(connector *connector.Connector) *Client {
	return &Client{
		connector: connector,
	}
}

func (c *Client) GetBasicInfo(ctx context.Context) (*qubicpb.BasicInfo, error) {
	rcf := connector.RequestContractFunction{
		ContractIndex: contractID,
		InputType:     ViewID.BasicInfo,
		InputSize:     0,
	}

	var result BasicInfo
	err := c.connector.PerformSmartContractRequest(ctx, rcf, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("performing smart contract request: %w", err)
	}

	bi, err := BasicInfoConverter.ToProto(result)
	if err != nil {
		return nil, fmt.Errorf("converting from node type: %w", err)
	}

	return bi, nil
}

func (c *Client) GetBetInfo(ctx context.Context, betID uint32) (*qubicpb.BetInfo, error) {
	rcf := connector.RequestContractFunction{
		ContractIndex: contractID,
		InputType:     ViewID.BetInfo,
		InputSize:     4, // sizeof betID; uint32
	}

	request := struct {
		BetID uint32
	}{
		BetID: betID,
	}

	var result BetInfo
	err := c.connector.PerformSmartContractRequest(ctx, rcf, request, &result)
	if err != nil {
		return nil, fmt.Errorf("performing smart contract request: %w", err)
	}

	bi, err := BetInfoConverter.ToProto(result)
	if err != nil {
		return nil, fmt.Errorf("converting from node type: %w", err)
	}

	return bi, nil
}

func (c *Client) GetActiveBets(ctx context.Context) (*qubicpb.ActiveBets, error) {
	rcf := connector.RequestContractFunction{
		ContractIndex: contractID,
		InputType:     ViewID.ActiveBet,
		InputSize:     0,
	}

	var result ActiveBets
	err := c.connector.PerformSmartContractRequest(ctx, rcf, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("performing smart contract request: %w", err)
	}

	ab := ActiveBetsConverter.ToProto(result)

	return ab, nil
}

func (c *Client) GetActiveBetsByCreator(ctx context.Context, creatorID common.Identity) (*qubicpb.ActiveBets, error) {
	rcf := connector.RequestContractFunction{
		ContractIndex: contractID,
		InputType:     ViewID.ActiveBetByCreator,
		InputSize:     32,
	}

	creatorPubKey, err := creatorID.ToPubKey(false)
	if err != nil {
		return nil, fmt.Errorf("converting creator identity to public key: %w", err)
	}

	request := struct {
		CreatorPubKey [32]byte
	}{
		CreatorPubKey: creatorPubKey,
	}

	var result ActiveBets
	err = c.connector.PerformSmartContractRequest(ctx, rcf, request, &result)
	if err != nil {
		return nil, fmt.Errorf("performing smart contract request: %w", err)
	}

	ab := ActiveBetsConverter.ToProto(result)

	return ab, nil
}

func (c *Client) GetBettorsByBetOption(ctx context.Context, betID, betOption uint32) (*qubicpb.BetOptionBettors, error) {
	rcf := connector.RequestContractFunction{
		ContractIndex: contractID,
		InputType:     ViewID.BetDetail,
		InputSize:     8, // sizeof betID + betOptions; 2 * uint32
	}

	request := struct {
		BetID     uint32
		BetOption uint32
	}{
		BetID:     betID,
		BetOption: betOption,
	}

	var result BetOptionDetail
	err := c.connector.PerformSmartContractRequest(ctx, rcf, request, &result)
	if err != nil {
		return nil, fmt.Errorf("performing smart contract request: %w", err)
	}

	bob, err := BetOptionBettorsConverter.ToProto(result)
	if err != nil {
		return nil, fmt.Errorf("converting from node type: %w", err)
	}

	return bob, nil
}
