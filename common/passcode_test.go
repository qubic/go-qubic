package common

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func Test_LogEncodedPasscode_NoError(t *testing.T) {
	encoded, err := EncodePasscodeToBase64([4]uint64{1, 2, 3, 4})

	require.NoError(t, err)
	log.Printf("Passcode: [%s]", encoded)
	assert.Equal(t, "AAAAAAAAAAEAAAAAAAAAAgAAAAAAAAADAAAAAAAAAAQ=", encoded)
}

func Test_LogDecodedPasscode_NoError(t *testing.T) {
	validBase64 := "AAAAAAAAAAEAAAAAAAAAAgAAAAAAAAADAAAAAAAAAAQ=" // 1 2 3 4

	decoded, err := DecodePasscodeFromBase64(validBase64)
	require.NoError(t, err)
	log.Printf("Passcode: %v", decoded)
	assert.Equal(t, [4]uint64{1, 2, 3, 4}, decoded)
}

func TestEncodeDecodePasscodeBase64(t *testing.T) {
	tcs := []struct {
		name     string
		passcode [4]uint64
	}{
		{"Zero values", [4]uint64{0, 0, 0, 0}},
		{"Small values", [4]uint64{1, 2, 3, 4}},
		{"Mixed values_1", [4]uint64{123, 456, 789, 101112}},
		{"Mixed values_2", [4]uint64{123, 0, 789, 101112}},
		{"Mixed values_3", [4]uint64{123, 456, 789, 0}},
		{"Mixed values_3", [4]uint64{0, 456, 789, 101112}},
		{"Large values", [4]uint64{18446744073709551615, 9223372036854775807, 1000000000000000000, 500}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := EncodePasscodeToBase64(tc.passcode)
			require.NoError(t, err, "Failed to encode passcode to Base64")

			decoded, err := DecodePasscodeFromBase64(encoded)
			require.NoError(t, err, "Failed to decode Base64 string")

			require.Equal(t, tc.passcode, decoded, "Decoded passcode does not match original")
		})
	}
}

func TestDecodeInvalidBase64(t *testing.T) {
	invalidBase64 := "invalid@@@base64##"

	_, err := DecodePasscodeFromBase64(invalidBase64)
	require.Error(t, err, "decode passcode should return an error")
}

func TestDecodeIncorrectSize(t *testing.T) {
	// Base64-encoded string representing less than 4 * 8 = 32 bytes
	incorrectBase64 := base64.StdEncoding.EncodeToString([]byte{1, 2, 3, 4})

	_, err := DecodePasscodeFromBase64(incorrectBase64)
	require.Error(t, err, "decode passcode should return an error")
}
