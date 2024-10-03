package qx

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/connector"
	"io"
)

const (
	contractID = 1
)

type viewID int

const (
	viewFeeID viewID = iota + 1
	viewAssetAskOrder
	viewAssetBidOrder
	viewEntityAskOrder
	viewEntityBidOrder
)

/*
struct QxFees_output
{
    uint32_t assetIssuanceFee; // Amount of qus
    uint32_t transferFee; // Amount of qus
    uint32_t tradeFee; // Number of billionths
};
*/

type Fees struct {
	AssetIssuanceFee uint32
	TransferFee      uint32
	TradeFee         uint32
}

func (f *Fees) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader
	err := header.UnmarshallFromReader(r)
	if err != nil {
		return errors.Wrap(err, "reading header")
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != connector.ContractFunctionResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", connector.ContractFunctionResponse, header.Type)
	}

	err = binary.Read(r, binary.LittleEndian, f)
	if err != nil {
		return errors.Wrap(err, "reading entire struct from buffer")
	}

	return nil
}

/*
struct qxGetAssetOrder_output{
    struct AssetOrder
    {
        uint8_t entity[32];
        long long price;
        long long numberOfShares;
    };

    Order orders[256];
};
*/

type AssetOrder struct {
	Entity         [32]byte
	Price          int64
	NumberOfShares int64
}

type AssetOrders []AssetOrder

var emptyAssetOrder AssetOrder

func (ao *AssetOrders) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader
	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return errors.Wrap(err, "reading header")
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != connector.ContractFunctionResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", connector.ContractFunctionResponse, header.Type)
	}

	receivedOrders := make([]AssetOrder, 256)
	err = binary.Read(r, binary.LittleEndian, receivedOrders)
	if err != nil {
		return errors.Wrap(err, "reading bytes from buffer")
	}

	orders := make([]AssetOrder, 0, 256)

	for _, order := range receivedOrders {
		if order == emptyAssetOrder {
			continue
		}

		orders = append(orders, order)
	}

	*ao = orders

	return nil
}

/*
struct qxGetEntityOrder_output{
    struct Order
    {
        uint8_t issuer[32];
        uint64_t assetName;
        long long price;
        long long numberOfShares;
    };
    Order orders[256];
};
*/

type EntityOrder struct {
	Issuer         [32]byte
	AssetName      [8]byte
	Price          int64
	NumberOfShares int64
}

type EntityOrders []EntityOrder

var emptyEntityOrder EntityOrder

func (eo *EntityOrders) UnmarshallFromReader(r io.Reader) error {
	var header connector.RequestResponseHeader
	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return errors.Wrap(err, "reading header")
	}

	if header.Type == connector.EndResponse {
		return nil
	}

	if header.Type != connector.ContractFunctionResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", connector.ContractFunctionResponse, header.Type)
	}

	receivedOrders := make([]EntityOrder, 256)
	err = binary.Read(r, binary.LittleEndian, receivedOrders)
	if err != nil {
		return errors.Wrap(err, "reading bytes from buffer")
	}

	orders := make([]EntityOrder, 0, 256)
	for _, order := range receivedOrders {
		if order == emptyEntityOrder {
			continue
		}
		orders = append(orders, order)
	}
	*eo = orders

	return nil
}

/*
struct qxGetAssetOrder_input{
    uint8_t issuer[32];
    uint64_t assetName;
    uint64_t offset;
};
*/

type getAssetOrderInput struct {
	Issuer    [32]byte
	AssetName uint64
	Offset    uint64
}

/*
struct qxGetEntityOrder_input{
    uint8_t entity[32];
    uint64_t offset;
};
*/

type getEntityOrderInput struct {
	Entity [32]byte
	Offset uint64
}
