package main

import (
	"context"
	"encoding/gob"
	"github.com/PPMac/gocn"
	"github.com/PPMac/gocn/register"
	"github.com/PPMac/gocn/rpc"
	"github.com/PPMac/order/api"
	"github.com/PPMac/order/model"
	"github.com/PPMac/order/service"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"log"
	"net/http"
)

func main() {
	engine := gocn.Default()
	client := rpc.NewHttpClient()
	client.RegisterHttpService("goods", &service.GoodsService{})
	group := engine.Group("order")

	group.Get("/findGrpc", func(ctx *gocn.Context) {
		config := rpc.DefaultGrpcClientConfig()
		config.Address = "localhost:9111"
		client, _ := rpc.NewGrpcClient(config)
		defer client.Conn.Close()
		goodsApiClient := api.NewGoodsApiClient(client.Conn)
		goodsResponse, _ := goodsApiClient.Find(context.Background(), &api.GoodsRequest{})
		ctx.JSON(http.StatusOK, goodsResponse)
	})

	group.Get("/findTcp", func(ctx *gocn.Context) {
		gob.Register(&model.Result{})
		gob.Register(&model.Goods{})
		option := rpc.DefaultOption
		option.SerializeType = rpc.ProtoBuff
		option.RegisterType = "nacos"
		option.RegisterOption = register.Option{
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
		}
		proxy := rpc.NewGocnTcpClientProxy(option)
		params := make([]any, 1)
		//params[0] = int64(1)
		//var Find func(id int64) any 作业
		result, err := proxy.Call(context.Background(), "goods", "Find", params)
		//Find(1)
		log.Println(err)
		ctx.JSON(http.StatusOK, result)
	})

	engine.Run(":9003")
}
