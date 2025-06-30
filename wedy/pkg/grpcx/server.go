package grpcx

import (
	"context"
	"fmt"
	etcd3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	Server   *grpc.Server
	client   *etcd3.Client
	Addr     string
	EtcdAddr string
	Port     int32
	KaCancel func()
}

//	func NewServer(client *etcd3.Client) *Server {
//		return &Server{
//			client: client,
//		}
//	}
func (s *Server) Serve() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	err = s.registerEtcd()
	if err != nil {
		return err
	}
	return s.Server.Serve(l)
}

func (s *Server) Stop() error {
	if s.KaCancel != nil {
		s.KaCancel()
	}
	if s.client != nil {
		s.client.Close()
	}
	s.Server.GracefulStop()
	return nil
}

func (s *Server) registerEtcd() error {
	client, err := etcd3.NewFromURL(s.EtcdAddr)
	if err != nil {
		return err
	}
	s.client = client
	em, err := endpoints.NewManager(s.client, "server/user")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	addr := "127.0.0.1:8090"
	key := "service/user/" + addr
	ttl := 5
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel1()
	leaseResp, err := s.client.Grant(ctx1, int64(ttl))
	if err != nil {
		return err
	}
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, etcd3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	kactx, kacancel := context.WithTimeout(context.Background(), time.Second*3)
	s.KaCancel = kacancel
	ch, err := s.client.KeepAlive(kactx, leaseResp.ID)
	go func() {
		for c := range ch {
			fmt.Println(c)
		}
	}()
	return err
}
