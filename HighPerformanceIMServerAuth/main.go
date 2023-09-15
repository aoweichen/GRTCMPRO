package main

import (
	"HighPerformanceIMServerAuth/Configs"
	"HighPerformanceIMServerAuth/DAO/MySQL"
	"HighPerformanceIMServerAuth/DAO/Redis"
	"HighPerformanceIMServerAuth/GRPC/AUTH/AUTHGRPC"
	"HighPerformanceIMServerAuth/Packages/Consul"
	"HighPerformanceIMServerAuth/Packages/Logger"
	"net"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func InitConfigs() {
	Configs.InitConfigs()
	Logger.InitLogger()
	MySQL.InitMySQLDB()
	Redis.InitRedisDB()
}

func main() {
	InitConfigs()
	// 创建grpc服务端
	server := grpc.NewServer()
	defer server.Stop()
	// 注册auth的服务实现
	authService := &AUTHGRPC.IMAuthGRPCService{}
	AUTHGRPC.RegisterIMAuthServiceServer(server, authService)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	service := &Consul.Service{
		ID:   Configs.ConfigData.AuthService.Name + Configs.ConfigData.AuthService.Host + ":" + strconv.Itoa(Configs.ConfigData.AuthService.ListenPort),
		Name: Configs.ConfigData.AuthService.Name,
		Host: Configs.ConfigData.AuthService.Host,
		Port: Configs.ConfigData.AuthService.ListenPort,
	}
	err := Consul.RegisterService(service, Configs.ConfigData.CONSUL.Host+Configs.ConfigData.CONSUL.Port)
	if err != nil {
		zap.S().Errorln(err)
		return
	}

	listen, err := net.Listen("tcp", Configs.ConfigData.AuthService.Host+":"+strconv.Itoa(Configs.ConfigData.AuthService.ListenPort))
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			zap.S().Errorln(err)
			return
		}
	}(listen)
	if err != nil {
		zap.S().Fatal(err)
	}
	zap.S().Infof("服务启动成功")
	err = server.Serve(listen)
	if err != nil {
		zap.S().Errorln(err.Error())
		return
	}
}
