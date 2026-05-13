package server

import (
	"context"
	qubicpb "github.com/qubic/go-qubic/v2/proto/v1"
	"github.com/qubic/go-qubic/v2/sdk/qx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type QxService struct {
	qubicpb.UnimplementedQxServiceServer
	qxClient *qx.Client
}

func NewQxService(qxClient *qx.Client) *QxService {
	return &QxService{qxClient: qxClient}
}

func (s *QxService) GetFees(ctx context.Context, _ *emptypb.Empty) (*qubicpb.Fees, error) {
	fees, err := s.qxClient.GetFees(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return fees, nil
}

func (s *QxService) GetAssetAskOrders(ctx context.Context, req *qubicpb.GetAssetOrdersRequest) (*qubicpb.AssetOrders, error) {
	aao, err := s.qxClient.GetAssetAskOrders(ctx, req.AssetName, req.IssuerId, req.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return aao, nil
}

func (s *QxService) GetAssetBidOrders(ctx context.Context, req *qubicpb.GetAssetOrdersRequest) (*qubicpb.AssetOrders, error) {
	abo, err := s.qxClient.GetAssetBidOrders(ctx, req.AssetName, req.IssuerId, req.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return abo, nil
}

func (s *QxService) GetEntityAskOrders(ctx context.Context, req *qubicpb.GetEntityOrdersRequest) (*qubicpb.EntityOrders, error) {
	eao, err := s.qxClient.GetEntityAskOrders(ctx, req.EntityId, req.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return eao, nil
}

func (s *QxService) GetEntityBidOrders(ctx context.Context, req *qubicpb.GetEntityOrdersRequest) (*qubicpb.EntityOrders, error) {
	ebo, err := s.qxClient.GetEntityBidOrders(ctx, req.EntityId, req.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return ebo, nil
}
