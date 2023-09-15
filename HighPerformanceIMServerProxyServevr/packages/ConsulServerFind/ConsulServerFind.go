package ConsulServerFind

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

// ConsulServerFindHandler函数用于在Consul中查找指定服务的地址列表
// 参数consulAddr表示Consul的地址，serviceName表示要查找的服务名称
func ConsulServerFindHandler(consulAddr string, serviceName string) ([]string, error) {
	// 创建Consul的默认配置
	consulConfigs := api.DefaultConfig()

	// 设置Consul的地址
	consulConfigs.Address = consulAddr

	// 创建Consul客户端
	client, err := api.NewClient(consulConfigs)
	if err != nil {
		zap.S().Fatal(err)
		return nil, err
	}

	// 使用Consul客户端查询指定服务的条目
	entries, _, err := client.Catalog().Service(serviceName, "", nil)
	if err != nil {
		zap.S().Fatal(err)
		return nil, err
	}

	var addressList []string
	// 遍历服务条目，获取地址信息并添加到地址列表中
	for _, entry := range entries {
		address := fmt.Sprintf("%s:%d", entry.ServiceAddress, entry.ServicePort)
		addressList = append(addressList, address)
	}

	// 返回地址列表和无错误
	return addressList, nil
}

// ConsulServerFindHandler函数用于在Consul中查找指定服务的地址列表
// 参数consulAddr表示Consul的地址，serviceName表示要查找的服务名称
func ConsulServerFindHandlerHTTP(consulAddr string, serviceName string) ([]string, error) {
	// 创建Consul的默认配置
	consulConfigs := api.DefaultConfig()

	// 设置Consul的地址
	consulConfigs.Address = consulAddr

	// 创建Consul客户端
	client, err := api.NewClient(consulConfigs)
	if err != nil {
		zap.S().Fatal(err)
		return nil, err
	}

	// 使用Consul客户端查询指定服务的条目
	entries, _, err := client.Catalog().Service(serviceName, "", nil)
	if err != nil {
		zap.S().Fatal(err)
		return nil, err
	}

	var addressList []string
	// 遍历服务条目，获取地址信息并添加到地址列表中
	for _, entry := range entries {
		address := fmt.Sprintf("http://%s:%d", entry.ServiceAddress, entry.ServicePort)
		addressList = append(addressList, address)
	}

	// 返回地址列表和无错误
	return addressList, nil
}
