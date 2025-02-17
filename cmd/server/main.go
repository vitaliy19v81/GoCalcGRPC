package main

import (
	pb "GRPCADDER/pkg/api/proto"
	"GRPCADDER/pkg/service"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

// secretKey - ключ для подписи JWT (замените на свой)
var secretKey = []byte("my-secret-key") // TODO перенести в секреты

// authInterceptor - проверяет JWT в метаданных gRPC-запроса
func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Извлекаем метаданные (заголовки)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["authorization"]) == 0 {
		return nil, fmt.Errorf("отсутствует токен")
	}

	// Получаем токен из заголовка "Bearer ..."
	tokenString := md["authorization"][0][7:] // Убираем "Bearer "

	// Разбираем и проверяем токен
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неизвестный метод подписи")
		}
		return secretKey, nil
	})

	// Проверяем валидность токена
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("недействительный токен")
	}

	//TODO Можно добавить проверку ролей и прав пользователя
	log.Printf("Аутентифицирован пользователь: %s", claims.Subject)

	// Вызываем обработчик gRPC (например, Calculate)
	return handler(ctx, req)
}

// Запускать на сервере
func main() {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor), // Добавляем перехватчик
	)
	srv := &service.GRPCServer{}
	pb.RegisterCalculatorServer(s, srv)

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
