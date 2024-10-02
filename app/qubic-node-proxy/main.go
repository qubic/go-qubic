package main

import (
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/connector"
	"github.com/qubic/go-qubic/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
)

const prefix = "QUBIC_NODE_PROXY"

func main() {
	if err := run(); err != nil {
		log.Fatalf("main: exited with error: %s", err.Error())
	}
}

func run() error {
	var cfg struct {
		Server struct {
			HttpListenAddr  string        `conf:"default:0.0.0.0:8000"`
			GrpcListenAddr  string        `conf:"default:0.0.0.0:8001"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		Pool struct {
			NodeFetcherUrl     string        `conf:"default:http://127.0.0.1:8080/status"`
			NodeFetcherTimeout time.Duration `conf:"default:5s"`
			InitialCap         int           `conf:"default:5"`
			MaxIdle            int           `conf:"default:6"`
			MaxCap             int           `conf:"default:10"`
			IdleTimeout        time.Duration `conf:"default:15s"`
		}
		NodeConnector struct {
			ConnectionPort        string        `conf:"default:21841"`
			ConnectionTimeout     time.Duration `conf:"default:5s"`
			HandlerRequestTimeout time.Duration `conf:"default:5s"`
		}
	}

	if err := conf.Parse(os.Args[1:], prefix, &cfg); err != nil {
		switch {
		case errors.Is(err, conf.ErrHelpWanted):
			usage, err := conf.Usage(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			log.Println(usage)
			return nil
		case errors.Is(err, conf.ErrVersionWanted):
			version, err := conf.VersionString(prefix, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			log.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config :\n%v\n", out)

	pfConfig := connector.PoolFetcherConfig{
		URL:            cfg.Pool.NodeFetcherUrl,
		RequestTimeout: cfg.Pool.NodeFetcherTimeout,
	}
	cConfig := connector.Config{
		ConnectionPort:        cfg.NodeConnector.ConnectionPort,
		ConnectionTimeout:     cfg.NodeConnector.ConnectionTimeout,
		HandlerRequestTimeout: cfg.NodeConnector.HandlerRequestTimeout,
	}
	pConfig := connector.PoolConfig{
		InitialCap:  cfg.Pool.InitialCap,
		MaxCap:      cfg.Pool.MaxCap,
		MaxIdle:     cfg.Pool.MaxIdle,
		IdleTimeout: cfg.Pool.IdleTimeout,
	}
	conn, err := connector.NewPoolConnector(pfConfig, cConfig, pConfig)
	if err != nil {
		return errors.Wrap(err, "creating pool connector")
	}

	srv := server.NewServer(cfg.Server.GrpcListenAddr, cfg.Server.HttpListenAddr, conn)
	err = srv.Start()
	if err != nil {
		return errors.Wrap(err, "starting server")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-shutdown:
			return errors.New("shutting down")
		}
	}
}
