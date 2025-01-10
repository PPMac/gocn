package rpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

//listen, _ := net.Listen("tcp", ":9111")
//	server := grpc.NewServer()
//	api.RegisterGoodsApiServer(server, &api.GoodsRpcService{})
//	err := server.Serve(listen)

type GocnGrpcServer struct {
	listen   net.Listener
	g        *grpc.Server
	register []func(g *grpc.Server)
	ops      []grpc.ServerOption
}

func NewGrpcServer(addr string, ops ...GocnGrpcOption) (*GocnGrpcServer, error) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	svc := &GocnGrpcServer{}
	svc.listen = listen
	for _, v := range ops {
		v.Apply(svc)
	}
	server := grpc.NewServer(svc.ops...)
	svc.g = server
	return svc, nil
}

func (s *GocnGrpcServer) Run() error {
	for _, f := range s.register {
		f(s.g)
	}
	return s.g.Serve(s.listen)
}

func (s *GocnGrpcServer) Stop() {
	s.g.Stop()
}

func (s *GocnGrpcServer) Register(f func(g *grpc.Server)) {
	s.register = append(s.register, f)
}

type GocnGrpcOption interface {
	Apply(s *GocnGrpcServer)
}

type DefaultGocnGrpcOption struct {
	f func(s *GocnGrpcServer)
}

func (d *DefaultGocnGrpcOption) Apply(s *GocnGrpcServer) {
	d.f(s)
}

func WithGrpcOptions(ops ...grpc.ServerOption) GocnGrpcOption {
	return &DefaultGocnGrpcOption{
		f: func(s *GocnGrpcServer) {
			s.ops = append(s.ops, ops...)
		},
	}
}

type GocnGrpcClient struct {
	Conn *grpc.ClientConn
}

func NewGrpcClient(config *GocnGrpcClientConfig) (*GocnGrpcClient, error) {
	var ctx = context.Background()
	var dialOptions = config.dialOptions

	if config.Block {
		//阻塞
		if config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
			defer cancel()
		}
		dialOptions = append(dialOptions, grpc.WithBlock())
	}
	if config.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.KeepAlive))
	}
	conn, err := grpc.DialContext(ctx, config.Address, dialOptions...)
	if err != nil {
		return nil, err
	}
	return &GocnGrpcClient{
		Conn: conn,
	}, nil
}

type GocnGrpcClientConfig struct {
	Address     string
	Block       bool
	DialTimeout time.Duration
	ReadTimeout time.Duration
	Direct      bool
	KeepAlive   *keepalive.ClientParameters
	dialOptions []grpc.DialOption
}

func DefaultGrpcClientConfig() *GocnGrpcClientConfig {
	return &GocnGrpcClientConfig{
		dialOptions: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
		DialTimeout: time.Second * 3,
		ReadTimeout: time.Second * 2,
		Block:       true,
	}
}
