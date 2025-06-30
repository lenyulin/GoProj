package grpc

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcd3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
)

type EtcdTest struct {
	suite.Suite
	client *etcd3.Client
}

func (s *EtcdTest) SetupSuite() {
	client, err := etcd3.NewFromURL("localhost:12379")
	require.NoError(s.T(), err)
	s.client = client
}
func (s *EtcdTest) TestClient() {
	etcdReslover, err := resolver.NewBuilder(s.client)
	require.NoError(s.T(), err)
	cc, err := grpc.Dial("etcd:///service/user",
		grpc.WithResolvers(etcdReslover),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig" : [{"round_robin": {}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(s.T(), err)
	client := NewUserServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.GetById(ctx, &GetByIdRequest{Id: 123})
	require.NoError(s.T(), err)
	fmt.Println(resp)
}
func (s *EtcdTest) TestServer() {
	t := s.T()
	em, err := endpoints.NewManager(s.client, "server/user")
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	addr := "127.0.0.1:8090"
	key := "service/user/" + addr
	l, err := net.Listen("tcp", ":8090")
	require.NoError(t, err)
	ttl := 5
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel1()
	leaseResp, err := s.client.Grant(ctx1, int64(ttl))
	require.NoError(t, err)
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, etcd3.WithLease(leaseResp.ID))
	require.NoError(t, err)
	kactx, kacancel := context.WithTimeout(context.Background(), time.Second*3)
	go func() {
		ch, e := s.client.KeepAlive(kactx, leaseResp.ID)
		require.NoError(t, e)
		for kaResp := range ch {
			fmt.Println(kaResp)
		}
	}()
	server := grpc.NewServer()
	RegisterUserServiceServer(server, &Server{})
	err = server.Serve(l)
	require.NoError(t, err)
	err = em.DeleteEndpoint(ctx, key)
	require.NoError(t, err)
	kacancel()
	server.GracefulStop()
	err = s.client.Close()
	if err != nil {
		return
	}
}
func TestEtcd(t *testing.T) {
	suite.Run(t, new(EtcdTest))
}
