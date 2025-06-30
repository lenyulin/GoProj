package grpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
)

type interceptorTestSuite struct {
	suite.Suite
}

func (s *interceptorTestSuite) TestClient() {
	cc, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(s.T(), err)
	client := NewUserServiceClient(cc)
	resp, err := client.GetById(context.Background(), &GetByIdRequest{Id: 123})
	require.NoError(s.T(), err)
	s.T().Logf("resp: %v", resp)
}
func (s *interceptorTestSuite) TestServer() {
	t := s.T()
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(NewLogInterceptor(t)))
	RegisterUserServiceServer(server, &Server{})
	l, err := net.Listen("tcp", ":8081")
	require.NoError(t, err)
	server.Serve(l)
}
func NewLogInterceptor(t *testing.T) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		t.Log("before log", req, info)
		resp, err = handler(ctx, req)
		t.Log("after log", resp, err)
		return
	}
}
func TestUnaryServerInterceptor(t *testing.T) {
	suite.Run(t, new(interceptorTestSuite))
}
