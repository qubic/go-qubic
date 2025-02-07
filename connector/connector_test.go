package connector

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

type injectPasscodeTestRequest struct {
	TickNumber uint32
	Passcode   [4]uint64
}

func (r *injectPasscodeTestRequest) AddPasscode(passcode [4]uint64) {
	r.Passcode = passcode
}

func TestInjectPasscode_WithExistingPasscode(t *testing.T) {
	req := injectPasscodeTestRequest{TickNumber: 15}
	passcodes := map[string][4]uint64{
		"127.0.0.1": {1, 2, 3, 4},
	}

	addr := net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8020}
	err := injectPasscode(&req, &addr, passcodes)
	require.NoError(t, err, "injectPasscode should not return an error")

	expected, err := binarySerialize(injectPasscodeTestRequest{TickNumber: 15, Passcode: [4]uint64{1, 2, 3, 4}})
	require.NoError(t, err, "binarySerialize the expected value should not return an error")

	serialized, err := binarySerialize(&req)
	require.NoError(t, err, "binarySerialize the actual value should not return an error")

	diff := cmp.Diff(expected, serialized)
	require.Equal(t, "", diff, "the actual value should match the expected value")
}

func TestInjectPasscode_WithNonExistingPasscode(t *testing.T) {
	req := injectPasscodeTestRequest{TickNumber: 15}
	passcodes := map[string][4]uint64{
		"127.0.0.1": {1, 2, 3, 4},
	}

	addr := net.TCPAddr{IP: net.ParseIP("192.168.0.1"), Port: 8020}
	err := injectPasscode(&req, &addr, passcodes)
	require.ErrorContains(t, err, "passcode not found for host 192.168.0.1", "injectPasscode should return an error")
}
