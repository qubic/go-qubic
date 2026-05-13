package core

import (
	"context"
	"fmt"

	"github.com/qubic/go-qubic/v2/common"
	"github.com/qubic/go-qubic/v2/connector"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
	"github.com/qubic/go-qubic/v2/sdk/core/nodetypes"
)

type Client struct {
	connector connector.RequestPerformer
}

func NewClient(connector connector.RequestPerformer) *Client {
	return &Client{
		connector: connector,
	}
}

func (c *Client) GetTickInfo(ctx context.Context) (*qubicpb.TickInfo, error) {
	var result nodetypes.TickInfo

	err := c.connector.PerformCoreRequest(ctx, nodetypes.CurrentTickInfoTypeRequest, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("handling chainRequest: %w", err)
	}

	tickInfoPb, err := result.ToProto()
	if err != nil {
		return nil, fmt.Errorf("converting tickInfo to proto: %w", err)
	}

	return tickInfoPb, nil
}

func (c *Client) GetAddressInfo(ctx context.Context, id string) (*qubicpb.EntityInfo, error) {
	identity := common.Identity(id)
	pubKey, err := identity.ToPubKey(false)
	if err != nil {
		return nil, fmt.Errorf("converting identity to public key: %w", err)
	}

	var result nodetypes.AddressInfo
	err = c.connector.PerformCoreRequest(ctx, nodetypes.BalanceTypeRequest, pubKey, &result)
	if err != nil {
		return nil, fmt.Errorf("handling chainRequest: %w", err)
	}

	ai, err := result.ToProto()
	if err != nil {
		return nil, fmt.Errorf("converting address info to proto: %w", err)
	}

	return ai, nil
}

func (c *Client) GetComputors(ctx context.Context) (*qubicpb.Computors, error) {
	var result nodetypes.Computors

	err := c.connector.PerformCoreRequest(ctx, nodetypes.ComputorsTypeRequest, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("handling chainRequest: %w", err)
	}

	comps, err := result.ToProto()
	if err != nil {
		return nil, fmt.Errorf("converting computors to proto: %w", err)
	}

	return comps, nil
}

func (c *Client) GetTickQuorumVote(ctx context.Context, tickNumber uint32) (*qubicpb.QuorumVote, error) {
	tickInfo, err := c.GetTickInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting tick info: %w", err)
	}

	if tickInfo.Tick < tickNumber {
		return nil, fmt.Errorf("Requested tick %d is in the future. Latest tick is: %d", tickNumber, tickInfo.Tick)
	}

	request := struct {
		Tick      uint32
		VoteFlags [(nodetypes.NumberOfComputors + 7) / 8]byte
	}{Tick: tickNumber}

	var result nodetypes.QuorumVotes

	err = c.connector.PerformCoreRequest(ctx, nodetypes.QuorumTickTypeRequest, request, &result)
	if err != nil {
		return nil, fmt.Errorf("handling chainRequest: %w", err)
	}

	tqv, err := result.ToProto()
	if err != nil {
		return nil, fmt.Errorf("converting tick quorum votes to proto: %w", err)
	}

	return tqv, nil
}

func (c *Client) GetTickData(ctx context.Context, tickNumber uint32) (*qubicpb.TickData, error) {
	tickInfo, err := c.GetTickInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting tick info: %w", err)
	}

	if tickInfo.Tick < tickNumber {
		return nil, fmt.Errorf("Requested tick %d is in the future. Latest tick is: %d", tickNumber, tickInfo.Tick)
	}

	request := struct{ Tick uint32 }{Tick: tickNumber}

	var result nodetypes.TickData
	err = c.connector.PerformCoreRequest(ctx, nodetypes.TickDataTypeRequest, request, &result)
	if err != nil {
		return nil, fmt.Errorf("handling chainRequest: %w", err)
	}

	td, err := result.ToProto()
	if err != nil {
		return nil, fmt.Errorf("converting tick data to proto: %w", err)
	}

	return td, nil
}

func (c *Client) GetTickTransactions(ctx context.Context, tickNumber uint32) (*qubicpb.TickTransactions, error) {
	tickData, err := c.GetTickData(ctx, tickNumber)
	if err != nil {
		return nil, fmt.Errorf("getting tick data: %w", err)
	}

	nrTx := len(tickData.TransactionIds)
	if nrTx == 0 {
		return &qubicpb.TickTransactions{
			Transactions: []*qubicpb.Transaction{},
		}, nil
	}

	requestTickTransactions := struct {
		Tick             uint32
		TransactionFlags [nodetypes.MaxNumberOfTransactionsPerTick / 8]uint8
	}{Tick: tickNumber}

	for i := 0; i < (nrTx+7)/8; i++ {
		requestTickTransactions.TransactionFlags[i] = 0
	}
	for i := (nrTx + 7) / 8; i < nodetypes.MaxNumberOfTransactionsPerTick/8; i++ {
		requestTickTransactions.TransactionFlags[i] = 1
	}

	var result nodetypes.Transactions
	err = c.connector.PerformCoreRequest(ctx, nodetypes.TickTransactionsTypeRequest, requestTickTransactions, &result)
	if err != nil {
		return nil, fmt.Errorf("handling chainRequest: %w", err)
	}

	txs, err := result.ToProto()
	if err != nil {
		return nil, fmt.Errorf("converting tick transactions to proto: %w", err)
	}

	return txs, nil
}

func (c *Client) GetTickTransactionsStatus(ctx context.Context, tickNumber uint32) (*qubicpb.TickTransactionsStatus, error) {
	tickInfo, err := c.GetTickInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting tick info: %w", err)
	}

	if tickInfo.Tick < tickNumber {
		return nil, fmt.Errorf("Requested tick %d is in the future. Latest tick is: %d", tickNumber, tickInfo.Tick)
	}

	request := struct {
		Tick uint32
	}{
		Tick: tickNumber,
	}

	var result nodetypes.TransactionStatus
	err = c.connector.PerformCoreRequest(ctx, nodetypes.TxStatusTypeRequest, request, &result)
	if err != nil {
		return nil, fmt.Errorf("handling chainRequest: %w", err)
	}

	txStatus, err := result.ToProto()
	if err != nil {
		return nil, fmt.Errorf("converting tick transactions status to proto: %w", err)
	}

	return txStatus, nil
}
