package common

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

func EncodePasscodeToBase64(passcode [4]uint64) (string, error) {
	buf := new(bytes.Buffer)
	for i, num := range passcode {
		err := binary.Write(buf, binary.BigEndian, num)
		if err != nil {
			return "", fmt.Errorf("writing passcode index %d to buffer: %w", i, err)
		}
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodePasscodeFromBase64(encodedPasscode string) ([4]uint64, error) {
	var arr [4]uint64

	data, err := base64.StdEncoding.DecodeString(encodedPasscode)
	if err != nil {
		return arr, fmt.Errorf("decoding base64 passcode: %w", err)
	}

	buf := bytes.NewReader(data)
	for i := range arr {
		if err := binary.Read(buf, binary.BigEndian, &arr[i]); err != nil {
			return arr, fmt.Errorf("reading passcode index %d from buffer: %w", i, err)
		}
	}

	return arr, nil
}
