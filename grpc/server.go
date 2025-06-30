package grpc

import "context"

type Server struct {
	UnimplementedUserServiceServer
}

func (s *Server) GetByID(ctx context.Context, request *GetByIdRequest) (*GetByIdResponse, error) {
	return &GetByIdResponse{
		User: &User{
			Id:   123,
			Name: "daming",
		},
	}, nil
}
