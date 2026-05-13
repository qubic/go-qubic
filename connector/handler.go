package connector

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	initialHandshakeTypeRequest = 27
)

type connHandler struct {
	prw    packetReadWriter
	conn   net.Conn
	connMu sync.Mutex
	peers  []string
}

func newConnHandler(conn net.Conn, defaultTimeout time.Duration) (*connHandler, error) {
	ch := connHandler{prw: packetReadWriter{defaultTimeout: defaultTimeout}, conn: conn}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := ch.handleInitialRequestAndSetPeers(ctx)
	if err != nil {
		return nil, fmt.Errorf("handling initial request: %w", err)
	}

	return &ch, nil
}

func (ch *connHandler) closeConn() error {
	ch.connMu.Lock()
	defer ch.connMu.Unlock()

	return ch.conn.Close()
}

func (ch *connHandler) getConn() net.Conn {
	ch.connMu.Lock()
	defer ch.connMu.Unlock()

	return ch.conn
}

func (ch *connHandler) handleCoreRequest(ctx context.Context, requestType uint8, requestData interface{}, dest ReaderUnmarshaler) error {
	conn := ch.getConn()

	req := newChainRequest(requestType, requestData)
	serializedRequest, err := req.serialize()
	if err != nil {
		return fmt.Errorf("serializing chainRequest for req type %d: %w", requestType, err)
	}

	err = ch.prw.writePacket(ctx, conn, serializedRequest)
	if err != nil {
		return fmt.Errorf("sending packet to qubic conn for req type %d: %w", requestType, err)
	}

	// if dest is nil then we don't care about the response
	if dest == nil {
		return nil
	}

	err = ch.prw.readPacket(ctx, conn, dest)
	if err != nil {
		return fmt.Errorf("reading response for req type %d: %w", requestType, err)
	}

	return nil
}

func (ch *connHandler) handleSmartContractRequest(ctx context.Context, reqContractFunction RequestContractFunction, requestData interface{}, dest ReaderUnmarshaler) error {
	conn := ch.getConn()

	req := newSmartContractRequest(reqContractFunction, requestData)
	serializedRequest, err := req.serialize()
	if err != nil {
		return fmt.Errorf(
			"serializing smart contract request for contract id: %d and input type: %d: %w",
			reqContractFunction.ContractIndex,
			reqContractFunction.InputType,
			err,
		)
	}

	err = ch.prw.writePacket(ctx, conn, serializedRequest)
	if err != nil {
		return fmt.Errorf("sending smart contract packet to qubic conn for contract id: %d and input type: %d: %w",
			reqContractFunction.ContractIndex,
			reqContractFunction.InputType,
			err,
		)
	}

	// if dest is nil then we don't care about the response
	if dest == nil {
		return nil
	}

	err = ch.prw.readPacket(ctx, conn, dest)
	if err != nil {
		return fmt.Errorf("reading smart contract response for contract id: %d and input type: %d: %w",
			reqContractFunction.ContractIndex,
			reqContractFunction.InputType,
			err,
		)
	}

	return nil
}

func (ch *connHandler) handleSmartContractRequestV2(ctx context.Context, reqContractFunction RequestContractFunction, requestData interface{}, dest ReaderUnmarshaler) error {
	conn := ch.getConn()

	req := newSmartContractRequest(reqContractFunction, requestData)
	serializedRequest, err := req.serialize()
	if err != nil {
		return fmt.Errorf(
			"serializing smart contract request for contract id: %d and input type: %d: %w",
			reqContractFunction.ContractIndex,
			reqContractFunction.InputType,
			err,
		)
	}

	err = ch.prw.writePacket(ctx, conn, serializedRequest)
	if err != nil {
		return fmt.Errorf("sending smart contract packet to qubic conn for contract id: %d and input type: %d: %w",
			reqContractFunction.ContractIndex,
			reqContractFunction.InputType,
			err,
		)
	}

	// if dest is nil then we don't care about the response
	if dest == nil {
		return nil
	}

	err = ch.prw.readPacket(ctx, conn, dest)
	if err != nil {
		return fmt.Errorf("reading smart contract response for contract id: %d and input type: %d: %w",
			reqContractFunction.ContractIndex,
			reqContractFunction.InputType,
			err,
		)
	}

	return nil
}

// this performs initial handshake with node which will return the list of known peers
func (ch *connHandler) handleInitialRequestAndSetPeers(ctx context.Context) error {
	var result PublicPeers
	err := ch.handleCoreRequest(ctx, initialHandshakeTypeRequest, nil, &result)
	if err != nil {
		return fmt.Errorf("sending req to node: %w", err)
	}
	ch.peers = result

	return nil
}
