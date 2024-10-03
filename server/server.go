package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/qubic/go-qubic/connector"
	qubicpb "github.com/qubic/go-qubic/proto/v1"
	"github.com/qubic/go-qubic/sdk/core"
	"github.com/qubic/go-qubic/sdk/quottery"
	"github.com/qubic/go-qubic/sdk/qx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
)

type Server struct {
	listenAddrGRPC string
	listenAddrHTTP string
	grpcServer     *grpc.Server
}

func NewServer(listenAddrGRPC, listenAddrHTTP string, connector *connector.Connector) *Server {
	grpcServer := createGrpcServerAndRegisterServices(connector)
	server := Server{
		listenAddrGRPC: listenAddrGRPC,
		listenAddrHTTP: listenAddrHTTP,
		grpcServer:     grpcServer,
	}

	return &server
}

func createGrpcServerAndRegisterServices(connector *connector.Connector) *grpc.Server {
	srv := grpc.NewServer(
		grpc.MaxRecvMsgSize(600*1024*1024),
		grpc.MaxSendMsgSize(600*1024*1024),
	)

	coreService := NewCoreService(core.NewClient(connector))
	qubicpb.RegisterCoreServiceServer(srv, coreService)

	quotteryService := NewQuotteryService(quottery.NewClient(connector))
	qubicpb.RegisterQuotteryServiceServer(srv, quotteryService)

	qxService := NewQxService(qx.NewClient(connector))
	qubicpb.RegisterQxServiceServer(srv, qxService)

	reflection.Register(srv)

	return srv
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.listenAddrGRPC)
	if err != nil {
		return errors.Wrap(err, "grpc failed to listen to tcp port")
	}

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	if s.listenAddrHTTP != "" {
		go func() {
			mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{EmitDefaultValues: true, EmitUnpopulated: false},
			}))
			opts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithDefaultCallOptions(
					grpc.MaxCallRecvMsgSize(600*1024*1024),
					grpc.MaxCallSendMsgSize(600*1024*1024),
				),
			}

			if err := qubicpb.RegisterCoreServiceHandlerFromEndpoint(
				context.Background(),
				mux,
				s.listenAddrGRPC,
				opts,
			); err != nil {
				panic(err)
			}

			if err := qubicpb.RegisterQuotteryServiceHandlerFromEndpoint(
				context.Background(),
				mux,
				s.listenAddrGRPC,
				opts,
			); err != nil {
				panic(err)
			}

			if err := qubicpb.RegisterQxServiceHandlerFromEndpoint(
				context.Background(),
				mux,
				s.listenAddrGRPC,
				opts,
			); err != nil {
				panic(err)
			}

			if err := http.ListenAndServe(s.listenAddrHTTP, mux); err != nil {
				panic(err)
			}
		}()
	}

	return nil
}
