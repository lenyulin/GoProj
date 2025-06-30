package service

import (
	"GoProj/wedy/interactive/domain"
	service2 "GoProj/wedy/interactive/service"
	"GoProj/wedy/internal/domian"
	svcmock "GoProj/wedy/internal/service/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func Test_Ranking_Svc(t *testing.T) {
	const batchSize = 2
	now := time.Now()
	testcases := []struct {
		name string

		mock       func(ctrl *gomock.Controller) (service2.InteractiveService, VideoService)
		batchsize  int
		wantVideos []domian.Video
		wantErr    error
	}{
		{
			name: "成功获取",
			mock: func(ctrl *gomock.Controller) (service2.InteractiveService, VideoService) {
				interSvc := svcmock.NewMockInteractiveService(ctrl)
				videoSvc := svcmock.NewMockVideoService(ctrl)
				//模拟批量获取数据
				videoSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 0, 2).
					Return([]domian.Video{
						{
							Uid:   1,
							Utime: now,
						}, {
							Uid:   2,
							Utime: now,
						},
					}, nil)
				videoSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 2, 2).
					Return([]domian.Video{
						{
							Uid:   3,
							Utime: now,
						}, {
							Uid:   4,
							Utime: now,
						},
					}, nil)
				videoSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 4, 2).
					Return([]domian.Video{}, nil)
				//
				interSvc.EXPECT().GetByIds(gomock.Any(), "video", []int64{1, 2}).
					Return(map[int64]domain.Interactive{
						1: {ReadCnt: 1},
						2: {ReadCnt: 2},
					}, nil)
				interSvc.EXPECT().GetByIds(gomock.Any(), "video", []int64{3, 4}).
					Return(map[int64]domain.Interactive{
						3: {ReadCnt: 3},
						4: {ReadCnt: 4},
					}, nil)
				return interSvc, videoSvc
			},
			wantErr: nil,
			wantVideos: []domian.Video{
				{
					Uid:   4,
					Utime: now,
				},
				{
					Uid:   3,
					Utime: now,
				},
				{
					Uid:   2,
					Utime: now,
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			interSvc, videoSvc := tc.mock(ctrl)
			svc := NewBatchRankingService(
				interSvc,
				videoSvc)
			videos, err := svc.topN(context.Background())
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantVideos, videos)
		})
	}
}
