package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"time"
)

func main() {
	// 建立与服务器的连接
	conn, err := grpc.Dial("127.0.0.1:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 创建健康检查的客户端
	healthClient := grpc_health_v1.NewHealthClient(conn)

	// 设置上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 发送健康检查请求
	response, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}

	// 打印健康检查结果
	fmt.Printf("Health check response: %v\n", response)
}
