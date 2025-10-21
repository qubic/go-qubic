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
