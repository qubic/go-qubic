package hrtypes

import (
	"fmt"

	"github.com/qubic/go-qubic/v2/common"
	"github.com/qubic/go-qubic/v2/sdk/core/nodetypes"
)

func NewSimpleTransferTransaction(sourceID, destinationID string, amount int64, targetTick uint32) (nodetypes.Transaction, error) {
	srcID := common.Identity(sourceID)
	destID := common.Identity(destinationID)
	srcPubKey, err := srcID.ToPubKey(false)
	if err != nil {
		return nodetypes.Transaction{}, fmt.Errorf("converting src id string to pubkey: %w", err)
	}
	destPubKey, err := destID.ToPubKey(false)
	if err != nil {
		return nodetypes.Transaction{}, fmt.Errorf("converting dest id string to pubkey: %w", err)
	}

	return nodetypes.Transaction{
		SourcePublicKey:      srcPubKey,
		DestinationPublicKey: destPubKey,
		Amount:               amount,
		Tick:                 targetTick,
		InputType:            0,
		InputSize:            0,
		Input:                nil,
	}, nil
}
