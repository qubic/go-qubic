package connector

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"net"
)

const (
	initialHandshakeTypeResponse = 0
)

type PublicPeers []string

func (pp *PublicPeers) UnmarshallFromReader(r io.Reader) error {
	var header RequestResponseHeader
	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return errors.Wrap(err, "reading header")
	}

	if header.Type != initialHandshakeTypeResponse {
		return errors.Errorf("Invalid header type, expected %d, found %d", initialHandshakeTypeResponse, header.Type)
	}

	var peers [4][4]byte

	err = binary.Read(r, binary.LittleEndian, &peers)
	if err != nil {
		return errors.Wrap(err, "reading public peers from reader")
	}

	for _, peer := range peers {
		if peer == [4]byte{} {
			continue
		}
		ip := net.IP(peer[:])
		if ip == nil {
			continue
		}

		*pp = append(*pp, ip.String())
	}

	var nextHeader RequestResponseHeader
	err = binary.Read(r, binary.BigEndian, &nextHeader)
	if err != nil {
		return errors.Wrap(err, "reading header")
	}

	ignoredBytes := make([]byte, nextHeader.GetSize()-uint32(binary.Size(nextHeader)))
	_, err = r.Read(ignoredBytes)
	if err != nil {
		return errors.Wrap(err, "reading ignored bytes")
	}

	return nil
}
