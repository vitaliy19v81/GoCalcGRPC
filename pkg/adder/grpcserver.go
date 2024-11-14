package adder

import (
	"GRPCADDER/pkg/api"
	pb "GRPCADDER/pkg/api"
	"context"
)

// GRPCServer ...
type GRPCServer struct {
	pb.UnimplementedAdderServer // Встраивание gRPC-сервера с пустой реализацией
}

// Add ...
func (s *GRPCServer) Add(ctx context.Context, req *api.AddRequest) (*api.AddResponse, error) {
	operation := req.GetOperation()
	switch operation {
	case "add":
		return &api.AddResponse{Result: req.GetX() + req.GetY()}, nil
	case "subtract":
		return &api.AddResponse{Result: req.GetX() - req.GetY()}, nil
	case "multiply":
		return &api.AddResponse{Result: req.GetX() * req.GetY()}, nil
	case "divide":
		if req.GetY() != 0 {
			return &api.AddResponse{Result: req.GetX() / req.GetY()}, nil
		} else {
			return &api.AddResponse{Error: "деление на ноль"}, nil
		}
	}
	return &api.AddResponse{Error: "Что то пошло не так"}, nil
}
