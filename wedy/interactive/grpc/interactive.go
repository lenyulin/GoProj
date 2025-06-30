package grpc

import (
	intrv1 "GoProj/wedy/api/proto/gen/wedy/api/proto/intr/v1"
	"GoProj/wedy/interactive/domain"
	"GoProj/wedy/interactive/service"
	"context"
	"google.golang.org/grpc"
)

type InteractiveServiceServer struct {
	intrv1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc}
}
func (i *InteractiveServiceServer) Register(server *grpc.Server) {
	intrv1.RegisterInteractiveServiceServer(server, i)
}
func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *intrv1.IncrReadCntRequest) (*intrv1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, request.GetBiz(), request.GetBizId())
	return &intrv1.IncrReadCntResponse{}, err
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	res, err := i.svc.GetByIds(ctx, request.GetBiz(), request.GetBizId())
	if err != nil {
		return &intrv1.GetByIdsResponse{}, err
	}
	intrs := make(map[int64]*intrv1.Interactive, len(res))
	for k, v := range res {
		intrs[k] = i.toDTO(v)
	}
	return &intrv1.GetByIdsResponse{
		Intrs: intrs,
	}, nil
}
func (i *InteractiveServiceServer) toDTO(intr domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:     intr.Biz,
		BizId:   intr.BizId,
		ReadCnt: intr.ReadCnt,
	}
}
func (i *InteractiveServiceServer) mustEmbedUnimplementedInteractiveServiceServer() {
	//TODO implement me
	panic("implement me")
}
