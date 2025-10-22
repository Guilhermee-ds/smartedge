package backend

import (
	"strconv"

	"github.com/hashicorp/consul/api"
)

func DiscoverBackends() []string {
	client, _ := api.NewClient(api.DefaultConfig())
	services, _, _ := client.Health().Service("smartedge-backend", "", true, nil)
	var urls []string
	for _, s := range services {
		urls = append(urls, "http://"+s.Service.Address+":"+strconv.Itoa(s.Service.Port))
	}
	return urls
}
