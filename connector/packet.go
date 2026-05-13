package connector

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type ReaderUnmarshaler interface {
	UnmarshallFromReader(r io.Reader) error
}

type packetReadWriter struct {
	defaultTimeout time.Duration
}

func (prw *packetReadWriter) writePacket(ctx context.Context, conn net.Conn, packet []byte) error {
	if packet == nil {
		return nil
	}

	// context deadline overrides defaultTimeout deadline
	writeDeadline := time.Now().Add(prw.defaultTimeout)
	deadline, ok := ctx.Deadline()
	if ok {
		writeDeadline = deadline
	}
	err := conn.SetWriteDeadline(writeDeadline)
	if err != nil {
		return fmt.Errorf("setting write deadline: %w", err)
	}
	defer conn.SetWriteDeadline(time.Time{})

	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("writing serialized binary data to connection: %w", err)
	}

	return nil
}

func (prw *packetReadWriter) readPacket(ctx context.Context, conn net.Conn, dest ReaderUnmarshaler) error {
	if dest == nil {
		return nil
	}

	// context deadline overrides defaultTimeout deadline
	readDeadline := time.Now().Add(prw.defaultTimeout)
	deadline, ok := ctx.Deadline()
	if ok {
		readDeadline = deadline
	}

	err := conn.SetReadDeadline(readDeadline)
	if err != nil {
		return fmt.Errorf("setting read deadline: %w", err)
	}
	defer conn.SetReadDeadline(time.Time{})

	err = dest.UnmarshallFromReader(conn)
	if err != nil {
		return fmt.Errorf("unmarshalling response: %w", err)
	}

	return nil
}
