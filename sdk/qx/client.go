package qx

import (
	"context"
	"encoding/binary"
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

func (c *Client) GetFees(ctx context.Context) (*qubicpb.Fees, error) {
	rcf := connector.RequestContractFunction{
		ContractIndex: contractID,
		InputType:     uint16(viewFeeID),
		InputSize:     0,
	}

	var result Fees
	err := c.connector.PerformSmartContractRequest(ctx, rcf, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("performing smart contract request: %w", err)
	}

	fees, err := FeesConverter.ToProto(result)
	if err != nil {
		return nil, fmt.Errorf("converting from node type: %w", err)
	}

	return fees, nil
}

func (c *Client) GetAssetAskOrders(ctx context.Context, name string, issuerID string, offset uint64) (*qubicpb.AssetOrders, error) {
	orders, err := c.getAssetOrders(ctx, uint16(viewAssetAskOrder), name, issuerID, offset)
	if err != nil {
		return nil, fmt.Errorf("getting asset orders: %w", err)
	}

	return orders, nil
}

func (c *Client) GetAssetBidOrders(ctx context.Context, name string, issuerID string, offset uint64) (*qubicpb.AssetOrders, error) {
	orders, err := c.getAssetOrders(ctx, uint16(viewAssetBidOrder), name, issuerID, offset)
	if err != nil {
		return nil, fmt.Errorf("getting asset orders: %w", err)
	}

	return orders, nil
}

func (c *Client) getAssetOrders(ctx context.Context, assetOrderType uint16, name string, issuerID string, offset uint64) (*qubicpb.AssetOrders, error) {
	var idPubkey [32]byte
	if issuerID != "" {
		id := common.Identity(issuerID)
		issuerPubKey, err := id.ToPubKey(false)
		if err != nil {
			return nil, fmt.Errorf("converting issuer id to pubkey: %w", err)
		}
		idPubkey = issuerPubKey
	}

	var assetName [8]byte
	copy(assetName[:], name)

	request := struct {
		IssuerPubKey [32]byte
		AssetName    uint64
		Offset       uint64
	}{
		IssuerPubKey: idPubkey,
		AssetName:    binary.LittleEndian.Uint64(assetName[:]),
		Offset:       offset,
	}

	reqSize := binary.Size(request)

	rcf := connector.RequestContractFunction{
		ContractIndex: contractID,
		InputType:     assetOrderType,
		InputSize:     uint16(reqSize),
	}

	var result AssetOrders
	err := c.connector.PerformSmartContractRequest(ctx, rcf, request, &result)
	if err != nil {
		return nil, fmt.Errorf("performing smart contract request: %w", err)
	}

	aao, err := AssetOrdersConverter.ToProto(result)
	if err != nil {
		return nil, fmt.Errorf("converting from node type: %w", err)
	}

	return aao, nil
}

func (c *Client) GetEntityAskOrders(ctx context.Context, entityID string, offset uint64) (*qubicpb.EntityOrders, error) {
	orders, err := c.getEntityOrders(ctx, uint16(viewEntityAskOrder), entityID, offset)
	if err != nil {
		return nil, fmt.Errorf("getting entity orders: %w", err)
	}

	return orders, nil
}

func (c *Client) GetEntityBidOrders(ctx context.Context, entityID string, offset uint64) (*qubicpb.EntityOrders, error) {
	orders, err := c.getEntityOrders(ctx, uint16(viewEntityBidOrder), entityID, offset)
	if err != nil {
		return nil, fmt.Errorf("getting entity orders: %w", err)
	}

	return orders, nil
}

func (c *Client) getEntityOrders(ctx context.Context, entityOrderType uint16, entityID string, offset uint64) (*qubicpb.EntityOrders, error) {
	id := common.Identity(entityID)
	entityPubKey, err := id.ToPubKey(false)
	if err != nil {
		return nil, fmt.Errorf("converting entity id to pubkey: %w", err)
	}

	request := struct {
		EntityPubKey [32]byte
		Offset       uint64
	}{
		EntityPubKey: entityPubKey,
		Offset:       offset,
	}

	reqSize := binary.Size(request)

	rcf := connector.RequestContractFunction{
		ContractIndex: contractID,
		InputType:     entityOrderType,
		InputSize:     uint16(reqSize),
	}

	var result EntityOrders
	err = c.connector.PerformSmartContractRequest(ctx, rcf, request, &result)
	if err != nil {
		return nil, fmt.Errorf("performing smart contract request: %w", err)
	}

	eo, err := EntityOrdersConverter.ToProto(result)
	if err != nil {
		return nil, fmt.Errorf("converting from node type: %w", err)
	}

	return eo, nil
}
