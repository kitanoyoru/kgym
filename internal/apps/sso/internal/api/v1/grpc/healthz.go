package grpc

import (
	"context"

	pb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthzServiceServer struct {
	pb.UnimplementedHealthServer
}

func NewHealthzService() *HealthzServiceServer {
	return &HealthzServiceServer{}
}

func (s *HealthzServiceServer) Check(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}

func (s *HealthzServiceServer) Watch(in *pb.HealthCheckRequest, _ pb.Health_WatchServer) error {
	return nil
}
