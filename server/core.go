package server

import (
	"context"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
	"github.com/qubic/go-qubic/v2/sdk/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ qubicpb.CoreServiceServer = &CoreService{}

type CoreService struct {
	qubicpb.UnimplementedCoreServiceServer
	coreClient *core.Client
}

func NewCoreService(coreClient *core.Client) *CoreService {
	return &CoreService{coreClient: coreClient}
}

func (s *CoreService) GetTickInfo(ctx context.Context, _ *emptypb.Empty) (*qubicpb.TickInfo, error) {
	ti, err := s.coreClient.GetTickInfo(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return ti, nil
}

func (s *CoreService) GetEntityInfo(ctx context.Context, req *qubicpb.GetEntityInfoRequest) (*qubicpb.EntityInfo, error) {
	ei, err := s.coreClient.GetAddressInfo(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return ei, nil
}

func (s *CoreService) GetComputors(ctx context.Context, _ *emptypb.Empty) (*qubicpb.Computors, error) {
	comps, err := s.coreClient.GetComputors(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return comps, nil
}

func (s *CoreService) GetTickQuorumVote(ctx context.Context, req *qubicpb.GetTickQuorumVoteRequest) (*qubicpb.QuorumVote, error) {
	qv, err := s.coreClient.GetTickQuorumVote(ctx, req.Tick)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return qv, nil
}

func (s *CoreService) GetTickData(ctx context.Context, req *qubicpb.GetTickDataRequest) (*qubicpb.TickData, error) {
	td, err := s.coreClient.GetTickData(ctx, req.Tick)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return td, nil
}

func (s *CoreService) GetTickTransactions(ctx context.Context, req *qubicpb.GetTickTransactionsRequest) (*qubicpb.TickTransactions, error) {
	txs, err := s.coreClient.GetTickTransactions(ctx, req.Tick)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return txs, nil
}

func (s *CoreService) GetTickTransactionsStatus(ctx context.Context, req *qubicpb.GetTickTransactionsStatusRequest) (*qubicpb.TickTransactionsStatus, error) {
	tts, err := s.coreClient.GetTickTransactionsStatus(ctx, req.Tick)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return tts, nil
}
