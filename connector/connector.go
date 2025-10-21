package connector

import (
	"context"
	"github.com/pkg/errors"
	"net"
	"time"
)

var _ RequestPerformer = &NoPoolConnector{}
var _ RequestPerformer = &Connector{}

type RequestPerformer interface {
	PerformCoreRequest(ctx context.Context, requestType uint8, requestData interface{}, dest ReaderUnmarshaler) error
	PerformSmartContractRequest(ctx context.Context, reqContractFunction RequestContractFunction, requestData interface{}, dest ReaderUnmarshaler) error
	PerformCoreRequestWithPasscode(ctx context.Context, requestType uint8, passcodes map[string][4]uint64, requestData PasscodeRequestData, dest ReaderUnmarshaler) error
}

type Connector struct {
	conPool *connPool
}

type Config struct {
	ConnectionPort        string
	ConnectionTimeout     time.Duration
	HandlerRequestTimeout time.Duration
}

func NewConnector(nodeIP string, connectorConfig Config) (*Connector, error) {
	scf := newSoloConnectionFactory(nodeIP, connectorConfig.ConnectionPort, connectorConfig.ConnectionTimeout, connectorConfig.HandlerRequestTimeout)
	pConfig := PoolConfig{
		InitialCap:  1,
		MaxCap:      5,
		MaxIdle:     3,
		IdleTimeout: 15 * time.Second,
	}
	cp, err := newConnectionPool(pConfig, scf.Connect, scf.Close)
	if err != nil {
		return nil, errors.Wrap(err, "creating new connection pool")
	}

	return &Connector{conPool: cp}, nil
}

type PoolFetcherConfig struct {
	URL            string
	RequestTimeout time.Duration
}

func NewPoolConnector(poolFetcherConfig PoolFetcherConfig, connectorConfig Config, poolConfig PoolConfig) (*Connector, error) {
	pcf := newPoolConnectionFactory(poolFetcherConfig.URL, poolFetcherConfig.RequestTimeout, connectorConfig.ConnectionPort, connectorConfig.ConnectionTimeout, connectorConfig.HandlerRequestTimeout)
	cp, err := newConnectionPool(poolConfig, pcf.Connect, pcf.Close)
	if err != nil {
		return nil, errors.Wrap(err, "creating new connection pool")
	}

	return &Connector{conPool: cp}, nil
}

func (c *Connector) WithConnection(f func(requestPerformer RequestPerformer) error) error {
	ch, err := c.conPool.Get()
	if err != nil {
		return errors.Wrap(err, "getting connection handler")
	}

	npc := NewNoPoolConnector(ch)
	err = f(npc)
	c.conPool.PutBack(ch, err)
	if err != nil {
		return errors.Wrap(err, "running function")
	}

	return nil
}

func (c *Connector) PerformCoreRequest(ctx context.Context, requestType uint8, requestData interface{}, dest ReaderUnmarshaler) error {
	var err error
	ch, err := c.conPool.Get()
	if err != nil {
		return errors.Wrap(err, "getting connection handler")
	}
	defer func() {
		c.conPool.PutBack(ch, err)
	}()

	err = ch.handleCoreRequest(ctx, requestType, requestData, dest)
	if err != nil {
		return errors.Wrap(err, "handling core request")
	}

	return nil
}

type PasscodeRequestData interface {
	AddPasscode([4]uint64)
}

type Session struct {
	ch      *connHandler
	conPool *connPool
}

func (s *Session) PerformCoreRequestWithPasscode(ctx context.Context, requestType uint8, passcodes map[string][4]uint64, requestData PasscodeRequestData, dest ReaderUnmarshaler) error {
	err := injectPasscode(requestData, s.ch.conn.RemoteAddr(), passcodes)
	if err != nil {
		return errors.Wrap(err, "injecting passcode")
	}

	err = s.ch.handleCoreRequest(ctx, requestType, requestData, dest)
	if err != nil {
		return errors.Wrap(err, "handling core request with passcode")
	}

	return nil
}

func (c *Connector) NewSession() (*Session, func(error), error) {
	ch, err := c.conPool.Get()
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting connection handler")
	}

	doneFunc := func(err error) {
		c.conPool.PutBack(ch, err)
	}

	return &Session{ch: ch, conPool: c.conPool}, doneFunc, nil
}

func (c *Connector) PerformCoreRequestWithPasscode(ctx context.Context, requestType uint8, passcodes map[string][4]uint64, requestData PasscodeRequestData, dest ReaderUnmarshaler) error {
	var err error
	ch, err := c.conPool.Get()
	if err != nil {
		return errors.Wrap(err, "getting connection handler")
	}
	defer func() {
		c.conPool.PutBack(ch, err)
	}()

	err = injectPasscode(requestData, ch.conn.RemoteAddr(), passcodes)
	if err != nil {
		return errors.Wrap(err, "injecting passcode")
	}

	err = ch.handleCoreRequest(ctx, requestType, requestData, dest)
	if err != nil {
		return errors.Wrap(err, "handling core request with passcode")
	}

	return nil
}

func (c *Connector) PerformSmartContractRequest(ctx context.Context, reqContractFunction RequestContractFunction, requestData interface{}, dest ReaderUnmarshaler) error {
	var err error
	ch, err := c.conPool.Get()
	if err != nil {
		return errors.Wrap(err, "getting connection handler")
	}
	defer func() {
		c.conPool.PutBack(ch, err)
	}()

	err = ch.handleSmartContractRequest(ctx, reqContractFunction, requestData, dest)
	if err != nil {
		return errors.Wrap(err, "handling smart contract request")
	}

	return nil
}

type NoPoolConnector struct {
	connHandler *connHandler
}

func NewNoPoolConnector(connHandler *connHandler) *NoPoolConnector {
	return &NoPoolConnector{connHandler: connHandler}
}

func (c *NoPoolConnector) PerformCoreRequestWithPasscode(ctx context.Context, requestType uint8, passcodes map[string][4]uint64, requestData PasscodeRequestData, dest ReaderUnmarshaler) error {
	err := injectPasscode(requestData, c.connHandler.conn.RemoteAddr(), passcodes)
	if err != nil {
		return errors.Wrap(err, "injecting passcode")
	}

	err = c.connHandler.handleCoreRequest(ctx, requestType, requestData, dest)
	if err != nil {
		return errors.Wrap(err, "handling core request with passcode")
	}

	return nil
}

func (c *NoPoolConnector) PerformCoreRequest(ctx context.Context, requestType uint8, requestData interface{}, dest ReaderUnmarshaler) error {
	err := c.connHandler.handleCoreRequest(ctx, requestType, requestData, dest)
	if err != nil {
		return errors.Wrap(err, "handling core request")
	}

	return nil
}

func (c *NoPoolConnector) PerformSmartContractRequest(ctx context.Context, reqContractFunction RequestContractFunction, requestData interface{}, dest ReaderUnmarshaler) error {
	err := c.connHandler.handleSmartContractRequest(ctx, reqContractFunction, requestData, dest)
	if err != nil {
		return errors.Wrap(err, "handling smart contract request")
	}

	return nil
}

func injectPasscode(requestData PasscodeRequestData, destinationNodeAddr net.Addr, passcodes map[string][4]uint64) error {
	passcode, err := getPasscodeForAddr(destinationNodeAddr, passcodes)
	if err != nil {
		return errors.Wrapf(err, "getting passcode for addr: %s", destinationNodeAddr.String())
	}

	requestData.AddPasscode(passcode)

	return nil
}

func getPasscodeForAddr(addr net.Addr, passcodes map[string][4]uint64) ([4]uint64, error) {
	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return [4]uint64{}, errors.Wrap(err, "splitting host and port")
	}

	passcode, ok := passcodes[host]
	if !ok {
		return [4]uint64{}, errors.Errorf("passcode not found for host %s", host)
	}

	return passcode, nil
}
