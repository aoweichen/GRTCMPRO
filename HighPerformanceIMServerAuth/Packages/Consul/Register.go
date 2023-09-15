package Consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type Service struct {
	ID   string
	Name string
	Host string
	Port int
}

func RegisterService(service *Service, consulAddr string) error {

	config := api.DefaultConfig()
	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Errorln(err)
		return err
	}

	registration := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Address: service.Host,
		Port:    service.Port,
		Tags:    []string{"Auth"},
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", service.Host, service.Port),
			Interval:                       "1s",
			Timeout:                        "2s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorln(err)
		return err
	}
	zap.S().Infof("Registered service '%s' with Consul", service.Name)
	return nil
}
