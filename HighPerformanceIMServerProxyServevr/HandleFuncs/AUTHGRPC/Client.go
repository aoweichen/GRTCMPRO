package AUTHGRPC

import (
	"HighPerformanceIMServerProxyServevr/GlobalVars"
	"HighPerformanceIMServerProxyServevr/packages/ConsulServerFind"
	"HighPerformanceIMServerProxyServevr/packages/LoadBalanceGRPC"
	context "context"

	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetTheTargetAddress() string {
	// 发送 grpc 请求
	authServiceList, err := ConsulServerFind.ConsulServerFindHandler(GlobalVars.ConsulAddr, GlobalVars.AuthServiceName)
	if err != nil {
		zap.S().Panicln(err)
	}
	// 创建一个负载均衡器实例
	lb := LoadBalanceGRPC.NewLoadBalancer(authServiceList)
	// 使用 grpcx 库进行服务发现和负载均衡
	addr := lb.GetNextAddress()
	return addr
}

func SendeEmailHandler(email string, emailType int) (*EmailCodeResponse, error) {
	addr := GetTheTargetAddress()
	// 建立连接
	// 创建不安全的凭证
	creds := insecure.NewCredentials()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		zap.S().Panicf("无法连接到gRPC服务器: %v", err)
		panic(err)
	}
	defer conn.Close()
	// 创建gRPC客户端
	client := NewIMAuthServiceClient(conn)
	// 发送请求
	request := &EmailCodeRequest{
		Email:     email,
		EmailType: int64(emailType),
	}
	response, err := client.SendEmailCode(context.Background(), request)
	if err != nil {
		zap.S().Errorln(err)
	}
	return response, err
}

func LoginHandler(email string, password string) (*LoginResponse, error) {
	addr := GetTheTargetAddress()
	// 建立连接
	creds := insecure.NewCredentials()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		zap.S().Panicf("无法连接到gRPC服务器: %v", err)
		panic(err)
	}
	defer conn.Close()
	// 创建gRPC客户端
	client := NewIMAuthServiceClient(conn)

	// 发送请求
	request := &LoginRequest{
		Email:    email,
		Password: password,
	}
	response, err := client.Login(context.Background(), request)
	if err != nil {
		zap.S().Errorln(err)
	}
	return response, err
}

func RegisterHandler(email string, name string, emailType int64, password string,
	passwordRepest string, emailCode string) (*RegisterResponse, error) {
	addr := GetTheTargetAddress()
	// 建立连接
	creds := insecure.NewCredentials()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		zap.S().Panicf("无法连接到gRPC服务器: %v", err)
		panic(err)
	}
	defer conn.Close()
	// 创建gRPC客户端
	client := NewIMAuthServiceClient(conn)

	// 发送请求
	request := &RegisterRequest{
		Email:          email,
		Name:           name,
		EmailType:      int64(emailType),
		Password:       password,
		PasswordRepeat: passwordRepest,
		EmailCode:      emailCode,
	}
	response, err := client.Register(context.Background(), request)
	if err != nil {
		zap.S().Errorln(err)
	}
	return response, err
}

func AuthHandler(token string) (*AuthResponse, error) {
	addr := GetTheTargetAddress()
	// 建立连接
	creds := insecure.NewCredentials()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		zap.S().Panicf("无法连接到gRPC服务器: %v", err)
		panic(err)
	}
	defer conn.Close()
	// 创建gRPC客户端
	client := NewIMAuthServiceClient(conn)

	// 发送请求
	request := &AuthRequest{
		Token: token,
	}
	response, err := client.IMAuthenticateHandler(context.Background(), request)
	if err != nil {
		zap.S().Errorln(err)
	}
	return response, err
}
