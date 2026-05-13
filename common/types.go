package common

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"unicode"

	"github.com/cloudflare/circl/xof/k12"
)

type Identity string

// FromPubKey creates a new identity from a public key
func (i *Identity) FromPubKey(pubKey [32]byte, isLowerCase bool) error {
	letter := 'A'
	if isLowerCase {
		letter = 'a'
	}

	var identity [60]byte

	for i := 0; i < 4; i++ {
		var publicKeyFragment = binary.LittleEndian.Uint64(pubKey[i*8 : (i+1)*8])
		for j := 0; j < 14; j++ {
			identity[i*14+j] = byte((publicKeyFragment % 26) + uint64(letter))
			publicKeyFragment /= 26
		}
	}

	h := k12.NewDraft10([]byte{})
	_, err := h.Write(pubKey[:])
	if err != nil {
		return fmt.Errorf("writing msg to k12: %w", err)
	}

	var identityBytesChecksum [3]byte
	_, err = h.Read(identityBytesChecksum[:])
	if err != nil {
		return fmt.Errorf("reading hash from k12: %w", err)
	}

	var identityBytesChecksumInt uint64
	identityBytesChecksumInt = uint64(identityBytesChecksum[0]) | (uint64(identityBytesChecksum[1]) << 8) | (uint64(identityBytesChecksum[2]) << 16)
	identityBytesChecksumInt &= 0x3FFFF

	for i := 0; i < 4; i++ {
		identity[56+i] = byte((identityBytesChecksumInt % 26) + uint64(letter))
		identityBytesChecksumInt /= 26
	}

	*i = Identity(identity[:])

	return nil
}

func (i *Identity) ToPubKey(isLowerCase bool) ([32]byte, error) {
	letters := []byte{'A', 'Z'}
	if isLowerCase {
		letters = []byte{'a', 'z'}
	}

	var pubKey [32]byte

	if !i.isValidIdFormat() {
		return [32]byte{}, fmt.Errorf("invalid ID format")
	}

	idBytes := []byte(string(*i))

	if len(idBytes) != 60 {
		return [32]byte{}, fmt.Errorf("invalid ID length, expected 60, found %d", len(idBytes))
	}

	for i := 0; i < 4; i++ {
		for j := 13; j >= 0; j-- {
			if idBytes[i*14+j] < letters[0] || idBytes[i*14+j] > letters[1] {
				return [32]byte{}, errors.New("invalid conversion")
			}

			im := binary.LittleEndian.Uint64(pubKey[i*8 : (i+1)*8])
			im = im*26 + uint64(idBytes[i*14+j]-letters[0])
			imBytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(imBytes, im)

			for k := 0; k < 8; k++ {
				pubKey[i*8+k] = imBytes[k]
			}
		}
	}

	return pubKey, nil
}

func (i *Identity) String() string {
	if i == nil {
		return ""
	}
	return string(*i)
}

// isValidIdFormat checks if the provided string has a valid ID format.
func (i *Identity) isValidIdFormat() bool {
	for _, c := range *i {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}

func PubKeysToIdentities(pubKeys [][32]byte, isLowercase bool) ([]Identity, error) {
	identities := make([]Identity, 0)
	for _, identity := range pubKeys {
		if identity == [32]byte{} {
			continue
		}
		id, err := getIDFrom32Bytes(identity, isLowercase)
		if err != nil {
			return nil, fmt.Errorf("getting identity from pubKey hex %s: %w", hex.EncodeToString(identity[:]), err)
		}
		identities = append(identities, id)
	}
	return identities, nil
}

func PubKeysToIdentitiesString(pubKeys [][32]byte, isLowercase bool) ([]string, error) {
	identities := make([]string, 0)
	for _, identity := range pubKeys {
		if identity == [32]byte{} {
			continue
		}
		id, err := getIDFrom32Bytes(identity, isLowercase)
		if err != nil {
			return nil, fmt.Errorf("getting identity from pubKey hex %s: %w", hex.EncodeToString(identity[:]), err)
		}
		identities = append(identities, id.String())
	}
	return identities, nil
}

func PubKeyToIdentity(pubKey [32]byte) (Identity, error) {
	return getIDFrom32Bytes(pubKey, false)
}

func DigestToTxID(digest [32]byte) (Identity, error) {
	return getIDFrom32Bytes(digest, true)
}

func BinarySerializeLE(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	var buff bytes.Buffer
	err := binary.Write(&buff, binary.LittleEndian, data)
	if err != nil {
		return nil, fmt.Errorf("writing data to buff: %w", err)
	}

	return buff.Bytes(), nil
}

func getIDFrom32Bytes(data [32]byte, isLowercase bool) (Identity, error) {
	var id Identity
	err := id.FromPubKey(data, isLowercase)
	if err != nil {
		return "", fmt.Errorf("getting id from pubkey: %w", err)
	}

	return id, nil
}

type QubicTransfer struct {
	DestinationID Identity
	Amount        int64
}
