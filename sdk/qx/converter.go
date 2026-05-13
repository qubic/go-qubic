package qx

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/qubic/go-qubic/v2/common"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
)

var FeesConverter feesConverter

type feesConverter struct{}

func (fc feesConverter) ToProto(fees Fees) (*qubicpb.Fees, error) {
	return &qubicpb.Fees{
		AssetIssuanceFee: fees.AssetIssuanceFee,
		TransferFee:      fees.TransferFee,
		TradeFee:         fees.TradeFee,
	}, nil
}

type assetOrdersConverter struct{}

var AssetOrdersConverter assetOrdersConverter

func (aoc assetOrdersConverter) ToProto(assetOrders AssetOrders) (*qubicpb.AssetOrders, error) {
	orders := make([]*qubicpb.AssetOrders_Order, 0, len(assetOrders))
	for _, assetOrder := range assetOrders {
		entityID, err := common.PubKeyToIdentity(assetOrder.Entity)
		if err != nil {
			return nil, fmt.Errorf("converting asset order entity pubkey: %s to id: %w", base64.StdEncoding.EncodeToString(assetOrder.Entity[:]), err)
		}
		order := &qubicpb.AssetOrders_Order{
			EntityId:       entityID.String(),
			Price:          assetOrder.Price,
			NumberOfShares: assetOrder.NumberOfShares,
		}

		orders = append(orders, order)
	}

	return &qubicpb.AssetOrders{Orders: orders}, nil
}

type entityOrdersConverter struct{}

var EntityOrdersConverter entityOrdersConverter

func (eoc entityOrdersConverter) ToProto(entityOrders EntityOrders) (*qubicpb.EntityOrders, error) {
	orders := make([]*qubicpb.EntityOrders_Order, 0, len(entityOrders))
	for _, entityOrder := range entityOrders {
		issuerID, err := common.PubKeyToIdentity(entityOrder.Issuer)
		if err != nil {
			return nil, fmt.Errorf("converting entity order issuer pubkey: %s to id: %w", base64.StdEncoding.EncodeToString(entityOrder.Issuer[:]), err)
		}

		order := &qubicpb.EntityOrders_Order{
			IssuerId:       issuerID.String(),
			Price:          entityOrder.Price,
			AssetName:      string(bytes.TrimRight(entityOrder.AssetName[:], "\x00")),
			NumberOfShares: entityOrder.NumberOfShares,
		}
		orders = append(orders, order)
	}

	return &qubicpb.EntityOrders{Orders: orders}, nil
}
