package events

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAssetIssuanceEvent_UnmarshalBinary(t *testing.T) {
	base64Encoded := "fBUfs37FBf00y/XqDc6kE/JNnjpN0DDl2QR/r0BhsKpAb0ABAAAAAAoAAAAAAAAAUUNBUAAAAAAAAAAAAAAA"
	data, err := base64.StdEncoding.DecodeString(base64Encoded)
	require.NoError(t, err, "decoding base64 data")
	var event AssetIssuanceEvent
	err = event.UnmarshalBinary(data)
	require.NoError(t, err, "unmarshalling binary data")
}

func TestAssetOwnershipChangeEvent_UnmarshalBinary(t *testing.T) {
	base64Encoded := "QvMt7n7vPwdDhVUbxbRVOxMpx/7trku3V9udvL77Hfm0XNyWnewpiwi3DPqGYe9p1T1ee0dgKChsGN91xWt9RAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAE1MTQAAAAAAAAAAAAAAAA=="
	data, err := base64.StdEncoding.DecodeString(base64Encoded)
	require.NoError(t, err, "decoding base64 data")
	var event AssetOwnershipChangeEvent
	err = event.UnmarshalBinary(data)
	require.NoError(t, err, "unmarshalling binary data")
}

func TestAssetPossessionChangeEvent_UnmarshalBinary(t *testing.T) {
	base64Encoded := "QvMt7n7vPwdDhVUbxbRVOxMpx/7trku3V9udvL77Hfm0XNyWnewpiwi3DPqGYe9p1T1ee0dgKChsGN91xWt9RAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAE1MTQAAAAAAAAAAAAAAAA=="
	data, err := base64.StdEncoding.DecodeString(base64Encoded)
	require.NoError(t, err, "decoding base64 data")

	var event AssetPossessionChangeEvent
	err = event.UnmarshalBinary(data)
	require.NoError(t, err, "unmarshalling binary data")
}

func TestBurningEvent_UnmarshalBinary(t *testing.T) {
	base64Encoded := "BAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKAAAAAAAAAA=="
	data, err := base64.StdEncoding.DecodeString(base64Encoded)
	require.NoError(t, err, "decoding base64 data")

	var event BurningEvent
	err = event.UnmarshalBinary(data)
	require.NoError(t, err, "unmarshalling binary data")
}

func TestContractReserveDeductionEvent_UnmarshalBinary(t *testing.T) {
	// Create test data: 8 bytes deduction + 8 bytes remaining + 4 bytes contract index + 4 bytes padding
	data := make([]byte, 24)
	// DeductionAmount: 1000 (uint64)
	data[0] = 0xE8
	data[1] = 0x03
	// RemainingAmount: 5000 (int64)
	data[8] = 0x88
	data[9] = 0x13
	// ContractIndex: 42 (uint32)
	data[16] = 0x2A

	var event ContractReserveDeductionEvent
	err := event.UnmarshalBinary(data)
	require.NoError(t, err, "unmarshalling binary data")
	require.Equal(t, uint64(1000), event.DeductionAmount)
	require.Equal(t, int64(5000), event.RemainingAmount)
	require.Equal(t, uint32(42), event.ContractIndex)
}

func TestContractReserveDeductionEvent_UnmarshalBinary_InvalidSize(t *testing.T) {
	data := make([]byte, 20) // Wrong size

	var event ContractReserveDeductionEvent
	err := event.UnmarshalBinary(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid contract reserve deduction event size")
}

func TestOracleQueryStatusChangeEvent_UnmarshalBinary(t *testing.T) {
	// Create test data: 32 bytes entity + 8 bytes queryID + 4 bytes interface + 1 byte type + 1 byte status
	data := make([]byte, 46)

	// QueryingEntity0: first 8 bytes as uint64 = 12345
	data[0] = 0x39
	data[1] = 0x30

	// QueryID at offset 32: 9999 (int64)
	data[32] = 0x0F
	data[33] = 0x27

	// InterfaceIndex at offset 40: 7 (uint32)
	data[40] = 0x07

	// Type at offset 44: 2 (OracleQueryTypeUserQuery)
	data[44] = 0x02

	// Status at offset 45: 3 (OracleQueryStatusSuccess)
	data[45] = 0x03

	var event OracleQueryStatusChangeEvent
	err := event.UnmarshalBinary(data)
	require.NoError(t, err, "unmarshalling binary data")
	require.Equal(t, int64(9999), event.QueryID)
	require.Equal(t, uint32(7), event.InterfaceIndex)
	require.Equal(t, uint8(2), event.Type)
	require.Equal(t, uint8(3), event.Status)
}

func TestOracleQueryStatusChangeEvent_UnmarshalBinary_InvalidSize(t *testing.T) {
	data := make([]byte, 40) // Wrong size

	var event OracleQueryStatusChangeEvent
	err := event.UnmarshalBinary(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid oracle query status change event size")
}

func TestOracleQueryStatusChangeEvent_StatusString(t *testing.T) {
	tests := []struct {
		status   uint8
		expected string
	}{
		{OracleQueryStatusPending, "pending"},
		{OracleQueryStatusCommitted, "committed"},
		{OracleQueryStatusSuccess, "success"},
		{OracleQueryStatusTimeout, "timeout"},
		{OracleQueryStatusUnresolvable, "unresolvable"},
		{99, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			event := OracleQueryStatusChangeEvent{Status: tt.status}
			require.Equal(t, tt.expected, event.StatusString())
		})
	}
}

func TestOracleQueryStatusChangeEvent_TypeString(t *testing.T) {
	tests := []struct {
		queryType uint8
		expected  string
	}{
		{OracleQueryTypeContractQuery, "contract_query"},
		{OracleQueryTypeContractSubscription, "contract_subscription"},
		{OracleQueryTypeUserQuery, "user_query"},
		{99, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			event := OracleQueryStatusChangeEvent{Type: tt.queryType}
			require.Equal(t, tt.expected, event.TypeString())
		})
	}
}
