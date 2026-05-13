package server

import (
	"context"
	"github.com/qubic/go-qubic/v2/common"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
	"github.com/qubic/go-qubic/v2/sdk/quottery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ qubicpb.QuotteryServiceServer = &QuotteryService{}

type QuotteryService struct {
	qubicpb.UnimplementedQuotteryServiceServer
	quotteryClient *quottery.Client
}

// NewQuotteryService creates the service and registers it to the grpc server
func NewQuotteryService(quotteryClient *quottery.Client) *QuotteryService {
	service := QuotteryService{quotteryClient: quotteryClient}

	return &service
}

func (s *QuotteryService) GetBasicInfo(ctx context.Context, _ *emptypb.Empty) (*qubicpb.BasicInfo, error) {
	bi, err := s.quotteryClient.GetBasicInfo(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return bi, nil
}

func (s *QuotteryService) GetBetInfo(ctx context.Context, req *qubicpb.GetBetInfoRequest) (*qubicpb.BetInfo, error) {
	bi, err := s.quotteryClient.GetBetInfo(ctx, req.BetId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return bi, nil
}

func (s *QuotteryService) GetActiveBets(ctx context.Context, _ *emptypb.Empty) (*qubicpb.ActiveBets, error) {
	ab, err := s.quotteryClient.GetActiveBets(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return ab, nil
}

func (s *QuotteryService) GetActiveBetsByCreator(ctx context.Context, req *qubicpb.GetActiveBetsByCreatorRequest) (*qubicpb.ActiveBets, error) {
	ab, err := s.quotteryClient.GetActiveBetsByCreator(ctx, common.Identity(req.CreatorId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return ab, nil
}

func (s *QuotteryService) GetBettorsByBetOption(ctx context.Context, req *qubicpb.GetBettorsByBetOptionRequest) (*qubicpb.BetOptionBettors, error) {
	bettors, err := s.quotteryClient.GetBettorsByBetOption(ctx, req.BetId, req.BetOption)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return bettors, nil
}
