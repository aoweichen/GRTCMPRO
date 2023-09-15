package LoadBalanceGRPC

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type LoadBalancer struct {
	AddressList []string
	Mu          sync.Mutex
}

func NewLoadBalancer(addressList []string) *LoadBalancer {
	return &LoadBalancer{
		AddressList: addressList,
	}
}

func (LB *LoadBalancer) GetNextAddress() string {
	LB.Mu.Lock()
	defer LB.Mu.Unlock()
	// 计算每个地址的响应时间，选择最短响应时间的地址
	var shortestTime time.Duration
	var selectedAddress string

	for _, address := range LB.AddressList {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			zap.S().Fatalf("Failed to dial address %s: %v", address, err)
			continue
		}
		// 创建 gRPC 健康监测客户端
		healthClient := grpc_health_v1.NewHealthClient(conn)
		// 开始时间
		startTime := time.Now()
		// 调用健康检查服务来测量响应时间
		_, err = healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			zap.S().Fatalf("Failed to check health status for address %s: %v", address, err)
			continue
		}
		responseTime := time.Since(startTime)
		if shortestTime == 0 || responseTime < shortestTime {
			shortestTime = responseTime
			selectedAddress = address
		}
		conn.Close()
	}
	return selectedAddress
}
