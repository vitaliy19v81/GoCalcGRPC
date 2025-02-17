// Package service Реализация функции Calculate
package service

import (
	pb "GRPCADDER/pkg/api/proto"
	"context"
	"fmt"
)

// GRPCServer ...
type GRPCServer struct {
	pb.UnimplementedCalculatorServer // Встраивание gRPC-сервера с пустой реализацией
}

// Calculate ...
func (s *GRPCServer) Calculate(ctx context.Context, req *pb.CalculationRequest) (*pb.CalculationResponse, error) {

	select {
	case <-ctx.Done():
		// Клиент отменил запрос
		return nil, fmt.Errorf("запрос отменён: %v", ctx.Err())
	default:
		// Обычная обработка
		operation := req.GetOperation()
		switch operation {
		case "add":
			return &pb.CalculationResponse{Result: req.GetX() + req.GetY()}, nil
		case "subtract":
			return &pb.CalculationResponse{Result: req.GetX() - req.GetY()}, nil
		case "multiply":
			return &pb.CalculationResponse{Result: req.GetX() * req.GetY()}, nil
		case "divide":
			if req.GetY() != 0 {
				return &pb.CalculationResponse{Result: req.GetX() / req.GetY()}, nil
			}
			return &pb.CalculationResponse{Error: "деление на ноль"}, nil
		}
		return &pb.CalculationResponse{Error: "неизвестная операция"}, nil
	}
}
