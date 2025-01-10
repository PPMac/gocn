package main

import (
	"encoding/gob"
	"github.com/PPMac/gocn/register"
	"github.com/PPMac/gocn/rpc"
	"github.com/PPMac/product/model"
	"github.com/PPMac/product/service"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"log"
)

func main() {
	tcpServer, err := rpc.NewTcpServer("localhost", 9111)
	if err != nil {
		log.Fatal(err)
	}
	tcpServer.SetRegister("nacos", register.Option{
		DialTimeout: 5000,
		NacosServerConfig: []constant.ServerConfig{
			{
				IpAddr:      "127.0.0.1",
				ContextPath: "/nacos",
				Port:        8848,
				Scheme:      "http",
			},
		},
		ServiceName: "goods",
	})
	gob.Register(&model.Result{})
	gob.Register(&model.Goods{})
	tcpServer.Register("goods", &service.GoodsRpcService{})
	tcpServer.Run()
}
