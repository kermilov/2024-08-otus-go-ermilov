package internalgrpc

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func loggingMiddleware() grpc.UnaryServerInterceptor {
	logFile, err := os.OpenFile("access_grpc.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic(err)
	}
	logger := log.New(logFile, "", 0)
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (
		interface{}, error,
	) {
		start := time.Now()
		result, err := handler(ctx, req)
		latency := time.Since(start)

		// Получение метаданных
		md, ok := metadata.FromIncomingContext(ctx)
		userAgent := ""
		if ok {
			if ua, exists := md["user-agent"]; exists && len(ua) > 0 {
				userAgent = ua[0]
			}
		}

		// Получение IP клиента
		clientIP := "unknown"
		p, ok := peer.FromContext(ctx)
		if ok {
			clientIP = p.Addr.String()
		}

		// Логирование информации о запросе
		logger.Printf("%s [%s] %s %s %d %d %q",
			clientIP,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			"POST",           // Метод (в gRPC это всегда POST)
			info.FullMethod,  // Путь
			status.Code(err), // Код ответа
			latency.Nanoseconds()/1e6,
			userAgent, // User-Agent
		)
		return result, err
	}
}
