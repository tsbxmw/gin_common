package consul

import (
    "fmt"
    "github.com/hashicorp/consul/api"
    "github.com/tsbxmw/datasource/common"
)

// Register register service
func (r *ConsulRegister) Register() error {
    config := api.DefaultConfig()
    config.Address = r.Address
    client, err := api.NewClient(config)
    if err != nil {
        return err
    }
    agent := client.Agent()

    IP := common.LocalIP()

    check := api.AgentServiceCheck{
        TCP:                            fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),
        Interval:                       r.Interval.String(),
        DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(),
        Timeout:                        "1s",
        Notes:                          "Basic health checks",
    }

    reg := &api.AgentServiceRegistration{
        ID:      fmt.Sprintf("%v-%v-%v", r.Service, IP, r.Port), // 服务节点的名称
        Name:    fmt.Sprintf("%v", r.Service),                   // 服务名称
        Tags:    r.Tag,                                          // tag，可以为空
        Port:    r.Port,                                         // 服务端口
        Address: IP,                                             // 服务 IP
        Check:   &check,
    }

    if err := agent.ServiceRegister(reg); err != nil {
        return err
    }

    return nil
}

func InitConsul() {
}
