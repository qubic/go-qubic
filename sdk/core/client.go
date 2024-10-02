package core

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/common"
	"github.com/qubic/go-qubic/connector"
	qubicpb "github.com/qubic/go-qubic/proto/v1"
	"github.com/qubic/go-qubic/sdk/core/nodetypes"
)

type Client struct {
	connector *connector.Connector
}

func NewClient(connector *connector.Connector) *Client {
	return &Client{
		connector: connector,
	}
}

func (c *Client) GetTickInfo(ctx context.Context) (*qubicpb.TickInfo, error) {
	var result nodetypes.TickInfo

	err := c.connector.PerformCoreRequest(ctx, nodetypes.CurrentTickInfoTypeRequest, nil, &result)
	if err != nil {
		return nil, errors.Wrap(err, "handling chainRequest")
	}

	tickInfoPb, err := result.ToProto()
	if err != nil {
		return nil, errors.Wrap(err, "converting tickInfo to proto")
	}

	return tickInfoPb, nil
}

func (c *Client) GetAddressInfo(ctx context.Context, id string) (*qubicpb.EntityInfo, error) {
	identity := common.Identity(id)
	pubKey, err := identity.ToPubKey(false)
	if err != nil {
		return nil, errors.Wrap(err, "converting identity to public key")
	}

	var result nodetypes.AddressInfo
	err = c.connector.PerformCoreRequest(ctx, nodetypes.BalanceTypeRequest, pubKey, &result)
	if err != nil {
		return nil, errors.Wrap(err, "handling chainRequest")
	}

	ai, err := result.ToProto()
	if err != nil {
		return nil, errors.Wrap(err, "converting address info to proto")
	}

	return ai, nil
}

func (c *Client) GetComputors(ctx context.Context) (*qubicpb.Computors, error) {
	var result nodetypes.Computors

	err := c.connector.PerformCoreRequest(ctx, nodetypes.ComputorsTypeRequest, nil, &result)
	if err != nil {
		return nil, errors.Wrap(err, "handling chainRequest")
	}

	comps, err := result.ToProto()
	if err != nil {
		return nil, errors.Wrap(err, "converting computors to proto")
	}

	return comps, nil
}

func (c *Client) GetTickQuorumVote(ctx context.Context, tickNumber uint32) (*qubicpb.QuorumVote, error) {
	tickInfo, err := c.GetTickInfo(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting tick info")
	}

	if tickInfo.Tick < tickNumber {
		return nil, errors.Errorf("Requested tick %d is in the future. Latest tick is: %d", tickNumber, tickInfo.Tick)
	}

	request := struct {
		Tick      uint32
		VoteFlags [(nodetypes.NumberOfComputors + 7) / 8]byte
	}{Tick: tickNumber}

	var result nodetypes.QuorumVotes

	err = c.connector.PerformCoreRequest(ctx, nodetypes.QuorumTickTypeRequest, request, &result)
	if err != nil {
		return nil, errors.Wrap(err, "handling chainRequest")
	}

	tqv, err := result.ToProto()
	if err != nil {
		return nil, errors.Wrap(err, "converting tick quorum votes to proto")
	}

	return tqv, nil
}

func (c *Client) GetTickData(ctx context.Context, tickNumber uint32) (*qubicpb.TickData, error) {
	tickInfo, err := c.GetTickInfo(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting tick info")
	}

	if tickInfo.Tick < tickNumber {
		return nil, errors.Errorf("Requested tick %d is in the future. Latest tick is: %d", tickNumber, tickInfo.Tick)
	}

	request := struct{ Tick uint32 }{Tick: tickNumber}

	var result nodetypes.TickData
	err = c.connector.PerformCoreRequest(ctx, nodetypes.TickDataTypeRequest, request, &result)
	if err != nil {
		return nil, errors.Wrap(err, "handling chainRequest")
	}

	td, err := result.ToProto()
	if err != nil {
		return nil, errors.Wrap(err, "converting tick data to proto")
	}

	return td, nil
}

func (c *Client) GetTickTransactions(ctx context.Context, tickNumber uint32) (*qubicpb.TickTransactions, error) {
	tickData, err := c.GetTickData(ctx, tickNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting tick data")
	}

	nrTx := len(tickData.TransactionIds)
	if nrTx == 0 {
		return &qubicpb.TickTransactions{
			Transactions: []*qubicpb.Transaction{},
		}, nil
	}

	requestTickTransactions := struct {
		Tick             uint32
		TransactionFlags [nodetypes.NumberOfTransactionsPerTick / 8]uint8
	}{Tick: tickNumber}

	for i := 0; i < (nrTx+7)/8; i++ {
		requestTickTransactions.TransactionFlags[i] = 0
	}
	for i := (nrTx + 7) / 8; i < nodetypes.NumberOfTransactionsPerTick/8; i++ {
		requestTickTransactions.TransactionFlags[i] = 1
	}

	var result nodetypes.Transactions
	err = c.connector.PerformCoreRequest(ctx, nodetypes.TickTransactionsTypeRequest, requestTickTransactions, &result)
	if err != nil {
		return nil, errors.Wrap(err, "handling chainRequest")
	}

	txs, err := result.ToProto()
	if err != nil {
		return nil, errors.Wrap(err, "converting tick transactions to proto")
	}

	return txs, nil
}

func (c *Client) GetTickTransactionsStatus(ctx context.Context, tickNumber uint32) (*qubicpb.TickTransactionsStatus, error) {
	tickInfo, err := c.GetTickInfo(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting tick info")
	}

	if tickInfo.Tick < tickNumber {
		return nil, errors.Errorf("Requested tick %d is in the future. Latest tick is: %d", tickNumber, tickInfo.Tick)
	}

	request := struct {
		Tick uint32
	}{
		Tick: tickNumber,
	}

	var result nodetypes.TransactionStatus
	err = c.connector.PerformCoreRequest(ctx, nodetypes.TxStatusTypeRequest, request, &result)
	if err != nil {
		return nil, errors.Wrap(err, "handling chainRequest")
	}

	txStatus, err := result.ToProto()
	if err != nil {
		return nil, errors.Wrap(err, "converting tick transactions status to proto")
	}

	return txStatus, nil
}
