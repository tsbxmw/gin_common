package consul

import (
	"github.com/hashicorp/consul/api"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/tsbxmw/datasource/common"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func (r *ConsulRegister) NewConsulClient() *consulapi.Client {
	rand.Seed(time.Now().UTC().UnixNano())

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = r.ConsulAddress + ":" + strconv.Itoa(r.ConsulPort)
	consulClient, err := consulapi.NewClient(consulConfig)
	if err != nil {
		logrus.Error("err", err)
		os.Exit(1)
	}
	return consulClient
}

func (r *ConsulRegister) RegisterHTTP() (*consulapi.Client) {

	consulClient := r.NewConsulClient()
	portStr := strconv.Itoa(r.Port)

	check := consulapi.AgentServiceCheck{
		HTTP:     "http://" + r.Address + ":" + portStr + "/v1/health",
		Interval: r.Interval.String(),
		Timeout:  "5s",
		Notes:    "Basic health checks",
	}

	asr := api.AgentServiceRegistration{
		ID:      r.Service + "-" + common.LocalIP(),
		Name:    r.Service,
		Address: r.Address,
		Port:    r.Port,
		Tags:    r.Tag,
		Check:   &check,
	}
	err := consulClient.Agent().ServiceRegister(&asr)
	if err != nil {
		logrus.Error(err)
		panic(err)
	}

	return consulClient
}
