package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/silenceper/pool"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type connPool struct {
	chPool pool.Pool
}

type PoolConfig struct {
	InitialCap  int
	MaxCap      int
	MaxIdle     int
	IdleTimeout time.Duration
}

func newConnectionPool(config PoolConfig, factoryFunc func() (interface{}, error), closeFunc func(interface{}) error) (*connPool, error) {
	cfg := pool.Config{
		InitialCap: config.InitialCap,
		MaxIdle:    config.MaxIdle,
		MaxCap:     config.MaxCap,
		Factory:    factoryFunc,
		Close:      closeFunc,
		//The maximum idle time of the connection, the connection exceeding this time will be closed, which can avoid the problem of automatic failure when connecting to EOF when idle
		IdleTimeout: config.IdleTimeout,
	}
	chPool, err := pool.NewChannelPool(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "creating channel pool")
	}

	return &connPool{chPool: chPool}, nil
}

func (c *connPool) Get() (*connHandler, error) {
	v, err := c.chPool.Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting pooled conn handler")
	}

	return v.(*connHandler), nil
}

func (c *connPool) PutBack(h *connHandler, err error) {
	if err != nil {
		cErr := c.chPool.Close(h)
		if cErr != nil {
			log.Printf("closing conn handler error: %v", cErr)
		}

		return
	}

	err = c.chPool.Put(h)
	if err != nil {
		log.Printf("putting conn handler error: %v", err)
	}
}

type poolConnectionFactory struct {
	PoolFetcherURL        string
	PoolFetcherTimeout    time.Duration
	ConnectionPort        string
	ConnectionTimeout     time.Duration
	HandlerRequestTimeout time.Duration
}

func newPoolConnectionFactory(poolFetcherUrl string, poolFetcherTimeout time.Duration, connectionPort string, connectionTimeout time.Duration, handlerRequestTimeout time.Duration) *poolConnectionFactory {
	return &poolConnectionFactory{
		PoolFetcherTimeout:    poolFetcherTimeout,
		PoolFetcherURL:        poolFetcherUrl,
		ConnectionTimeout:     connectionTimeout,
		ConnectionPort:        connectionPort,
		HandlerRequestTimeout: handlerRequestTimeout,
	}
}

func (pcf *poolConnectionFactory) Connect() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pcf.PoolFetcherTimeout)
	defer cancel()

	peer, err := pcf.getNewRandomPeer(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting new random peer")
	}

	addr := net.JoinHostPort(peer, pcf.ConnectionPort)
	conn, err := net.DialTimeout("tcp", addr, pcf.ConnectionTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "dialing %s", addr)
	}

	ch, err := newConnHandler(conn, pcf.HandlerRequestTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "creating new conn handler for address %s", addr)
	}

	fmt.Printf("connected to: %s\n", peer)
	return ch, nil
}

func (pcf *poolConnectionFactory) Close(v interface{}) error { return v.(*connHandler).closeConn() }

type statusResponse struct {
	MaxTick          uint32         `json:"max_tick"`
	LastUpdate       int64          `json:"last_update"`
	ReliableNodes    []nodeResponse `json:"reliable_nodes"`
	MostReliableNode nodeResponse   `json:"most_reliable_node"`
}

type nodeResponse struct {
	Address    string `json:"address"`
	LastTick   uint32 `json:"last_tick"`
	LastUpdate int64  `json:"last_update"`
}

func (pcf *poolConnectionFactory) getNewRandomPeer(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pcf.PoolFetcherURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "creating new request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "getting peers from node fetcher")
	}
	defer res.Body.Close()

	var resp statusResponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "reading response body")
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", errors.Wrap(err, "unmarshalling response")
	}

	if len(resp.ReliableNodes) == 0 {
		return "", errors.New("no reliable nodes")
	}

	peer := resp.ReliableNodes[rand.Intn(len(resp.ReliableNodes))]

	fmt.Printf("Got %d new peers. Selected random %s\n", len(resp.ReliableNodes), peer.Address)

	return peer.Address, nil
}

type soloConnectionFactory struct {
	ConnectionIP          string
	ConnectionPort        string
	ConnectionTimeout     time.Duration
	HandlerRequestTimeout time.Duration
}

func newSoloConnectionFactory(connectionIP string, connectionPort string, connTimeout time.Duration, handlerRequestTimeout time.Duration) *soloConnectionFactory {
	return &soloConnectionFactory{
		ConnectionIP:          connectionIP,
		ConnectionPort:        connectionPort,
		ConnectionTimeout:     connTimeout,
		HandlerRequestTimeout: handlerRequestTimeout,
	}
}

func (scf *soloConnectionFactory) Connect() (interface{}, error) {
	addr := net.JoinHostPort(scf.ConnectionIP, scf.ConnectionPort)
	conn, err := net.DialTimeout("tcp", addr, scf.ConnectionTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "dialing %s", addr)
	}

	ch, err := newConnHandler(conn, scf.HandlerRequestTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "creating new conn handler for address %s", addr)
	}

	fmt.Printf("connected to: %s\n", scf.ConnectionIP)

	return ch, nil
}

func (scf *soloConnectionFactory) Close(v interface{}) error { return v.(*connHandler).closeConn() }
