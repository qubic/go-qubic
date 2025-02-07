package common

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"github.com/pkg/errors"
)

func EncodePasscodeToBase64(passcode [4]uint64) (string, error) {
	buf := new(bytes.Buffer)
	for i, num := range passcode {
		err := binary.Write(buf, binary.BigEndian, num)
		if err != nil {
			return "", errors.Wrapf(err, "writing passcode index %d to buffer", i)
		}
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodePasscodeFromBase64(encodedPasscode string) ([4]uint64, error) {
	var arr [4]uint64

	data, err := base64.StdEncoding.DecodeString(encodedPasscode)
	if err != nil {
		return arr, errors.Wrap(err, "decoding base64 passcode")
	}

	buf := bytes.NewReader(data)
	for i := range arr {
		if err := binary.Read(buf, binary.BigEndian, &arr[i]); err != nil {
			return arr, errors.Wrapf(err, "reading passcode index %d from buffer", i)
		}
	}

	return arr, nil
}
