package connector

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

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

func (c *Connector) PerformCoreRequest(ctx context.Context, requestType uint8, requestData interface{}, dest ReaderUnmarshaler) error {
	var err error
	ch, err := c.conPool.Get()
	if err != nil {
		return errors.Wrap(err, "getting connection handler")
	}
	defer c.conPool.PutBack(ch, err)

	err = ch.handleCoreRequest(ctx, requestType, requestData, dest)
	if err != nil {
		return errors.Wrap(err, "handling core request")
	}

	return nil
}

func (c *Connector) PerformSmartContractRequest(ctx context.Context, reqContractFunction RequestContractFunction, requestData interface{}, dest ReaderUnmarshaler) error {
	var err error
	ch, err := c.conPool.Get()
	if err != nil {
		return errors.Wrap(err, "getting connection handler")
	}
	defer c.conPool.PutBack(ch, err)

	err = ch.handleSmartContractRequest(ctx, reqContractFunction, requestData, dest)
	if err != nil {
		return errors.Wrap(err, "handling smart contract request")
	}

	return nil
}
